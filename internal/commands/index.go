package commands

import (
	"fmt"

	"github.com/sahil485/memex/pkg/indexer"
)

func Index(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: memex index <file_path>")
	}

	filePath := args[0]
	fmt.Printf("Indexing file %s...\n", filePath)

	err := indexer.IndexFile(filePath)
	if err != nil {
		return fmt.Errorf("indexing failed: %w", err)
	}

	fmt.Println("✓ Indexing complete")
	return nil
}
