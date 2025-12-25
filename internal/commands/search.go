package commands

import (
	"encoding/json"
	"fmt"

	"github.com/sahil485/memex/pkg/search"
	"github.com/sahil485/memex/pkg/types"
)

func Search(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: memex search <query>")
	}

	query := args[0]

	results, err := search.Search(query, 10)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Printf("Found %d results:\n", results.EstimatedTotalHits)
	for _, hit := range results.Hits {
		// Decode to Document struct
		var doc types.Document
		if err := hit.DecodeInto(&doc); err != nil {
			return fmt.Errorf("failed to decode document: %w", err)
		}

		rankingScore := 0.0
		if scoreRaw, ok := hit["_rankingScore"]; ok {
			var score float64
			if err := json.Unmarshal(scoreRaw, &score); err == nil {
				rankingScore = score
			}
		}

		// Print formatted result
		fmt.Printf("  - %s (score: %.2f)\n", doc.Path, rankingScore)
		fmt.Printf("    Name: %s\n", doc.Name)
		fmt.Println()
	}

	return nil
}
