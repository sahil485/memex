package client

import (
	"fmt"
	"strings"

	"github.com/sahil485/memex/pkg/config"
)

func InitializeIndex() error {
	c := New()

	_, err := c.CreateIndex(config.IndexName)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "index_already_exists") {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	err = ConfigureIndexSettings()
	if err != nil {
		return fmt.Errorf("failed to configure index settings: %w", err)
	}

	return nil
}
