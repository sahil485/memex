package commands

import (
	"fmt"
	"github.com/sahil485/memex/pkg/client"
	"github.com/sahil485/memex/pkg/config"
)

func Init(args []string) error {
	fmt.Println("Initializing memex...")

	err := client.InitializeIndex()
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	fmt.Printf("✓ Index '%s' ready\n", config.IndexName)
	return nil
}
