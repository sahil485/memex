package commands

import (
	"fmt"
	"strings"

	"github.com/sahil485/memex/pkg/indexer"
)

func Index(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: memex index <directory> [--ignore pattern1,pattern2,...]")
	}

	directory := args[0]
	var ignorePatterns []string

	// Parse --ignore flags (can be specified multiple times)
	for i := 1; i < len(args); i++ {
		if args[i] == "--ignore" {
			if i+1 >= len(args) {
				return fmt.Errorf("--ignore requires a pattern argument")
			}
			// Get the next argument
			patterns := args[i+1]
			// Split comma-separated patterns
			parts := strings.Split(patterns, ",")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					ignorePatterns = append(ignorePatterns, p)
				}
			}
			i++ // Skip the pattern argument
		}
	}

	fmt.Printf("Indexing directory %s...\n", directory)
	if len(ignorePatterns) > 0 {
		fmt.Printf("Ignoring %d patterns:\n", len(ignorePatterns))
		for _, p := range ignorePatterns {
			fmt.Printf("  - %s\n", p)
		}
	}

	err := indexer.IndexDirectory(directory, ignorePatterns)
	if err != nil {
		return fmt.Errorf("indexing failed: %w", err)
	}

	fmt.Println("✓ Indexing complete")
	return nil
}
