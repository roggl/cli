package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	url = "http://localhost:8080/api"
)

func format(date *time.Time) string {
	layout := "2 Jan 2006 15:04:05"
	return date.Format(layout)
}

func printRecursive(t *Trie, indent, key string) {
	fmt.Printf("%s%s:\n", indent, key)

	for _, val := range t.Value {
		fmt.Printf("%s  %s - ", indent, format(val.Begin))
		if val.Recording() {
			fmt.Print("RECORDING ")
		} else {
			fmt.Printf("%s ", format(val.End))
		}
		if len(val.Tags) > 0 {
			fmt.Printf("tags: %s", strings.Join(val.Tags, ", "))
		}
		fmt.Println()
	}

	for key, child := range t.Children {
		printRecursive(child, indent+"  ", key)
	}
}

func main() {
	args := os.Args[1:]

	popArg := func() string {
		if len(args) == 0 {
			fmt.Println("expected argument")
			os.Exit(1)
		}

		res := args[0]
		args = args[1:]
		return res
	}

	op := get
	if len(args) > 0 {
		opRaw := popArg()

		switch opRaw {
		case "get":
			op = get
		case "create":
			op = create
		case "start":
			op = startTrack
		case "stop":
			op = stopTrack
		}
	}

	client := NewClient(url)

	switch op {
	case get:
		trie, err := client.Get()
		if err != nil {
			panic(err)
		}

		recording, path := trie.GetRecorded()
		if recording {
			fmt.Printf("recording: %s\n", strings.Join(path, "/"))
		}

		for key, child := range trie.Children {
			printRecursive(child, "", key)
		}

	case create:
		var path []string
		for len(args) > 0 {
			path = append(path, popArg())
		}
		if err := client.Create(path); err != nil {
			panic(err)
		}

	case startTrack:
		var path []string
		for len(args) > 0 {
			path = append(path, popArg())
		}
		if err := client.Start(path); err != nil {
			panic(err)
		}

	case stopTrack:
		if err := client.Stop(); err != nil {
			panic(err)
		}

	default:
		panic("unhandled operation")
	}
}
