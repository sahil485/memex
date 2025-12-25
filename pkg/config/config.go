package config

const (
	MeilisearchURL  = "http://127.0.0.1:58273"
	MeilisearchPort = 58273
	IndexName       = "files"
)

// IgnoredDirectories contains directory names that should be skipped during indexing
var IgnoredDirectories = map[string]bool{
	// macOS system directories
	"Library":      true,
	"Applications": true,
	"System":       true,
	"__MACOSX":     true,

	// JavaScript/Node
	"node_modules": true,

	// Python
	"__pycache__":   true,
	"venv":          true,
	"env":           true,
	"site-packages": true,

	// Go
	"vendor": true,

	// Rust
	"target": true,

	// .NET/C#
	"bin":      true,
	"obj":      true,
	"packages": true,

	// Build outputs
	"build":  true,
	"dist":   true,
	"out":    true,
	"_build": true,

	// Cache directories
	"cache": true,

	// Test/Coverage
	"coverage": true,
	"htmlcov":  true,

	// Logs
	"logs": true,

	// Temporary
	"tmp":  true,
	"temp": true,

	// iOS/macOS development
	"Pods":        true,
	"Carthage":    true,
	"DerivedData": true,

	// Elixir
	"deps":  true,
	"_deps": true,
}

func ShouldIgnoreDirectory(dirName string) bool {
	return IgnoredDirectories[dirName]
}
