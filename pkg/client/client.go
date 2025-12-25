package client

import (
	"sync"
	"time"

	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/sahil485/memex/pkg/config"
)

type Client struct {
	ms meilisearch.ServiceManager
}

var (
	instance *Client
	once     sync.Once
)

func New() *Client {
	once.Do(func() {
		ms := meilisearch.New(config.MeilisearchURL)
		instance = &Client{ms: ms}
	})
	return instance
}

func (c *Client) GetIndex() meilisearch.IndexManager {
	return c.ms.Index(config.IndexName)
}

func (c *Client) CreateIndex(indexName string) (*meilisearch.TaskInfo, error) {
	return c.ms.CreateIndex(&meilisearch.IndexConfig{
		Uid: indexName,
	})
}

func (c *Client) WaitForTask(taskUID int64) (*meilisearch.Task, error) {
	// Wait up to 1 seconds for the task to complete
	return c.ms.WaitForTask(taskUID, 1*time.Second)
}
