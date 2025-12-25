package indexer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sahil485/memex/pkg/client"
	"github.com/sahil485/memex/pkg/types"
)

func IndexFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	doc := types.NewDocument(
		filePath,
		fileInfo.Name(),
		filepath.Dir(filePath),
		filepath.Ext(filePath),
		fileInfo.Size(),
		fileInfo.ModTime().Unix(),
		string(content),
	)

	c := client.New()

	// Add document to Meilisearch index
	idx := c.GetIndex()
	task, err := idx.AddDocuments([]types.Document{*doc}, nil)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}

	taskInfo, err := c.WaitForTask(task.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for indexing task: %w", err)
	}

	// Check if the task failed
	if taskInfo.Status == "failed" {
		return fmt.Errorf("indexing task failed: %s", taskInfo.Error.Message)
	}

	return nil
}

func IndexFiles(directory string) error {
	_ = client.New()
	return nil
}
