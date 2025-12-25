package main

import (
	"fmt"
	"os"
	"github.com/sahil485/memex/internal/commands"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: memex [command] [options]")
		os.Exit(1)
	}

	command := args[0]
	options := args[1:]

	var err error

	switch command {
	case "init":
		err = commands.Init(options)
	case "search":
		err = commands.Search(options)
	case "index":
		err = commands.Index(options)
	case "clear-index":
		err = commands.ClearIndex(options)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
