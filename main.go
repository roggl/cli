package main

import (
	"fmt"
	"os"
	"sort"
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

func foreachChildren(t *Trie, fn func(string, *Trie)) {
	var keys []string
	for key := range t.Children {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fn(key, t.Children[key])
	}
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

	foreachChildren(t, func(key string, child *Trie) {
		printRecursive(child, indent+"  ", key)
	})
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

		var found bool
		op, found = map[string]operation{
			"get":    get,
			"create": create,
			"start":  startTrack,
			"stop":   stopTrack,
		}[opRaw]

		if !found {
			fmt.Printf("unknown operation '%s'\n", opRaw)
			os.Exit(1)
		}
	}

	client := NewClient(url)

	switch op {
	case get:
		trie, err := client.Get()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		recording, path := trie.GetRecorded()
		if recording {
			fmt.Printf("recording: %s\n", strings.Join(path, "/"))
		}

		foreachChildren(trie, func(key string, child *Trie) {
			printRecursive(child, "", key)
		})

	case create:
		var path []string
		for len(args) > 0 {
			path = append(path, popArg())
		}
		if err := client.Create(path); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case startTrack:
		var path []string
		for len(args) > 0 {
			path = append(path, popArg())
		}
		if err := client.Start(path); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case stopTrack:
		if err := client.Stop(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	default:
		fmt.Println("unhandled operation")
		os.Exit(1)
	}
}
