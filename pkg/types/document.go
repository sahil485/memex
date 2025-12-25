package types

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Document struct {
	ID string `json:"id"`

	Path    string `json:"path"`
	Name    string `json:"name"`
	Dir     string `json:"dir"`
	Ext     string `json:"ext"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`

	Content     string `json:"content"`
	ContentHash string `json:"content_hash"`

	Description string   `json:"description,omitempty"`

	IndexedAt int64 `json:"indexed_at"`
}

func NewDocument(path, name, dir, ext string, size, modTime int64, content string) *Document {
	hash := sha256.Sum256([]byte(content))
	contentHash := hex.EncodeToString(hash[:])

	// Generate a valid Meilisearch ID using path hash
	// Meilisearch IDs can only contain alphanumeric, hyphens, and underscores
	pathHash := sha256.Sum256([]byte(path))
	id := hex.EncodeToString(pathHash[:])

	return &Document{
		ID:          id,
		Path:        path,
		Name:        name,
		Dir:         dir,
		Ext:         ext,
		Size:        size, // bytes
		ModTime:     modTime,
		Content:     content,
		ContentHash: contentHash,
		IndexedAt:   time.Now().Unix(),
	}
}
