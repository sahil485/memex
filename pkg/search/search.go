package search

import (
	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/sahil485/memex/pkg/client"
)

func Search(query string, limit int64) (*meilisearch.SearchResponse, error) {
	c := client.New()
	return c.GetIndex().Search(query, &meilisearch.SearchRequest{
		Limit:            limit,
		ShowRankingScore: true,
	})
}
