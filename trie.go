package main

import (
	"strings"
	"time"
)

type Value struct {
	Begin *time.Time `json:"begin"`
	End   *time.Time `json:"end"`
	Tags  []string   `json:"tags"`
}

func (v Value) Recording() bool {
	return v.Begin != nil && v.End == nil
}

type Trie struct {
	Children map[string]*Trie `json:"children"`
	Value    []Value          `json:"value"`
}

func getRecorded(t *Trie, path []string) (bool, []string) {
	currentRecording := func(leaf *Trie) bool {
		for _, val := range leaf.Value {
			if val.Recording() {
				return true
			}
		}
		return false
	}

	if currentRecording(t) {
		return true, path
	}

	for key, child := range t.Children {
		found, newPath := getRecorded(child, append(path[:], key))
		if found {
			return true, newPath
		}
	}

	return false, path
}

func (t *Trie) GetRecorded() (recording bool, path []string) {
	return getRecorded(t, []string{})
}

func gather(t *Trie, path []string, m map[string][]Value) {
	pathString := strings.Join(path, "/")

	m[pathString] = t.Value

	for key, child := range t.Children {
		newPath := append(path[:], key)
		gather(child, newPath, m)
	}
}

func (t *Trie) Gather() map[string][]Value {
	m := make(map[string][]Value)
	gather(t, []string{}, m)
	return m
}
