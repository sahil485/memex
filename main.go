package main

import (
	"fmt"
	"os"
)

func main() {
	var args []string = os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: memex [command] [options]")
		os.Exit(1)
	}

	var command string = args[0]
	var options []string = args[1:]

	switch command {
	case "search":
		{
			fmt.Println("Searching for files...")
			fmt.Println(options)
		}
	}
}
