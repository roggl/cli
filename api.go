package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{
		url: url,
	}
}

func (c *Client) Get() (*Trie, error) {
	resp, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res *Trie
	if err := json.Unmarshal(bytes, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) Create(path []string) error {
	joined := strings.Join(path, "+")
	resp, err := http.Post(c.url+"/projects/"+joined, "", nil)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		return errors.New("invalid status code")
	}
	return nil
}

func (c *Client) Start(path []string) error {
	joined := strings.Join(path, "+")
	resp, err := http.Post(c.url+"/projects/"+joined+"/start", "", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return errors.New("project doesn't exist")
	case http.StatusBadRequest:
		return errors.New("already recording")

	case http.StatusCreated:
		return nil

	default:
		return errors.New("invalid status code")
	}
}

func (c *Client) Stop() error {
	resp, err := http.Post(c.url+"/stop", "", nil)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return errors.New("not recording")

	case http.StatusOK:
		return nil

	default:
		return errors.New("invalid status code")
	}
}
