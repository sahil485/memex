package indexer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sahil485/memex/pkg/client"
	"github.com/sahil485/memex/pkg/config"
	"github.com/sahil485/memex/pkg/types"
)

func IndexFile(filePath string) error {
	if !config.IsAllowedExtension(filepath.Ext(filePath)) {
		return fmt.Errorf("file extension not allowed: %s", filepath.Ext(filePath))
	}

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

func createDocumentForFile(filePath string) (*types.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(filePath)
	var content string

	if config.ShouldIgnoreContent(ext) {
		content = ""
	} else {
		contentBytes, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		content = string(contentBytes)
	}

	doc := types.NewDocument(
		filePath,
		fileInfo.Name(),
		filepath.Dir(filePath),
		ext,
		fileInfo.Size(),
		fileInfo.ModTime().Unix(),
		content,
	)

	return doc, nil
}

func IndexDirectory(directory string, ignorePatterns []string) error {
	ms_client := client.New()
	documents := make([]types.Document, 0)
	fileCount := 0

	expandedIgnorePatterns := make([]string, 0, len(ignorePatterns))
	for _, pattern := range ignorePatterns {
		if strings.HasPrefix(pattern, "~/") {
			home, err := os.UserHomeDir()
			if err == nil {
				pattern = filepath.Join(home, pattern[2:])
			}
		}
		pattern = filepath.Clean(pattern)
		expandedIgnorePatterns = append(expandedIgnorePatterns, pattern)
	}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Skipping %s: %v\n", path, err)
			return nil
		}

		// Check if path matches any ignore patterns
		for _, pattern := range expandedIgnorePatterns {
			if filepath.IsAbs(pattern) {
				if strings.HasPrefix(path, pattern) {
					if info.IsDir() && path == pattern {
						return filepath.SkipDir
					}
					if strings.HasPrefix(path, pattern+string(filepath.Separator)) {
						if info.IsDir() {
							return filepath.SkipDir
						}
						return nil
					}
				}
			} else {
				// Pattern is a glob pattern, match against filename
				matched, err := filepath.Match(pattern, info.Name())
				if err == nil && matched {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}

		if info.IsDir() && config.ShouldIgnoreDirectory(info.Name()) {
			return filepath.SkipDir
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !config.IsAllowedExtension(filepath.Ext(path)) {
			return nil
		}

		if !info.IsDir() {
			fileCount++
			fmt.Printf("[%d] %s\n", fileCount, path)
			doc, err := createDocumentForFile(path)
			if err != nil {
				fmt.Printf("Skipping %s: %v\n", path, err)
				return nil
			}

			documents = append(documents, *doc)
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nIndexing %d files to Meilisearch...\n", len(documents))

	task, err := ms_client.GetIndex().AddDocuments(documents, nil)
	if err != nil {
		return fmt.Errorf("failed to add documents to index: %w", err)
	}

	taskInfo, err := ms_client.WaitForTask(task.TaskUID)
	if err != nil {
		return fmt.Errorf("failed to wait for indexing task: %w", err)
	}

	if taskInfo.Status == "failed" {
		return fmt.Errorf("indexing task failed: %s", taskInfo.Error.Message)
	}

	return nil
}
