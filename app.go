package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sahil485/memex/pkg/client"
	"github.com/sahil485/memex/pkg/indexer"
	"github.com/sahil485/memex/pkg/search"
	"github.com/sahil485/memex/pkg/types"
)

// App struct
type App struct {
	ctx              context.Context
	meilisearchCmd   *exec.Cmd
	meilisearchReady bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start MeiliSearch if not already running
	if !a.GetMeilisearchHealth() {
		go a.startMeilisearch()
	}
}

// startMeilisearch launches MeiliSearch as a background process
func (a *App) startMeilisearch() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	memexDir := filepath.Join(homeDir, ".memex")
	meilisearchBinary := filepath.Join(memexDir, "meilisearch")

	// Check if meilisearch binary exists
	if _, err := os.Stat(meilisearchBinary); os.IsNotExist(err) {
		return fmt.Errorf("meilisearch binary not found at %s", meilisearchBinary)
	}

	// Launch MeiliSearch
	cmd := exec.Command(
		meilisearchBinary,
		"--db-path", filepath.Join(memexDir, "data.ms"),
		"--http-addr", "127.0.0.1:7700",
		"--no-analytics",
	)

	// Redirect output to avoid blocking
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start meilisearch: %w", err)
	}

	a.meilisearchCmd = cmd

	// Wait for MeiliSearch to be ready
	for i := 0; i < 30; i++ {
		time.Sleep(500 * time.Millisecond)
		if a.GetMeilisearchHealth() {
			a.meilisearchReady = true
			fmt.Println("MeiliSearch started successfully")
			return nil
		}
	}

	return fmt.Errorf("meilisearch failed to start within 15 seconds")
}

// shutdown cleans up MeiliSearch process
func (a *App) shutdown(ctx context.Context) {
	if a.meilisearchCmd != nil && a.meilisearchCmd.Process != nil {
		a.meilisearchCmd.Process.Kill()
	}
}

// SearchResult represents a search result for the frontend
type SearchResult struct {
	ID            string  `json:"id"`
	Path          string  `json:"path"`
	Content       string  `json:"content"`
	Type          string  `json:"type"`
	Title         string  `json:"title"`
	RankingScore  float64 `json:"rankingScore"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	Hits               []SearchResult `json:"hits"`
	Query              string         `json:"query"`
	ProcessingTimeMs   int64          `json:"processingTimeMs"`
	EstimatedTotalHits int64          `json:"estimatedTotalHits"`
	Error              string         `json:"error,omitempty"`
}

// Search performs a search query
func (a *App) Search(query string, limit int) SearchResponse {
	if query == "" {
		return SearchResponse{
			Hits:  []SearchResult{},
			Query: query,
			Error: "",
		}
	}

	result, err := search.Search(query, int64(limit))
	if err != nil {
		return SearchResponse{
			Query: query,
			Error: fmt.Sprintf("Search failed: %v", err),
		}
	}

	// Convert meilisearch results to our frontend format
	hits := make([]SearchResult, 0, len(result.Hits))
	for _, hit := range result.Hits {
		// Decode the MeiliSearch Hit into our Document type, then convert to SearchResult
		var doc types.Document
		if err := hit.Decode(&doc); err != nil {
			continue
		}

		sr := SearchResult{
			ID:           doc.ID,
			Path:         doc.Path,
			Content:      doc.Content,
			Type:         doc.Ext,  // File extension
			Title:        doc.Name, // File name
			RankingScore: 0,        // Will be populated if available
		}

		// Try to get the ranking score from the raw hit data
		if scoreData, ok := hit["_rankingScore"]; ok {
			var score float64
			if err := json.Unmarshal(scoreData, &score); err == nil {
				sr.RankingScore = score
			}
		}

		hits = append(hits, sr)
	}

	return SearchResponse{
		Hits:               hits,
		Query:              result.Query,
		ProcessingTimeMs:   result.ProcessingTimeMs,
		EstimatedTotalHits: result.EstimatedTotalHits,
	}
}

// OpenFile opens a file in the default editor
func (a *App) OpenFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// On macOS, use 'open' command
		cmd = exec.Command("open", path)
	case "linux":
		// On Linux, use 'xdg-open'
		cmd = exec.Command("xdg-open", path)
	case "windows":
		// On Windows, use 'start'
		cmd = exec.Command("cmd", "/c", "start", path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// IndexFile indexes a single file
func (a *App) IndexFile(path string) error {
	return indexer.IndexFile(path)
}

// IndexDirectory indexes all files in a directory
func (a *App) IndexDirectory(path string) error {
	// Pass empty ignore patterns for now - could be made configurable later
	return indexer.IndexDirectory(path, []string{})
}

// GetMeilisearchHealth checks if MeiliSearch is running
func (a *App) GetMeilisearchHealth() bool {
	c := client.New()
	index := c.GetIndex()

	// Try to get index stats as a health check
	_, err := index.GetStats()
	return err == nil
}
