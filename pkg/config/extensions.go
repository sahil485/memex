package config

// AllowedExtensions contains all file extensions that should be indexed
var AllowedExtensions = map[string]bool{
	// Plain text
	".txt": true,
	".md":  true,
	".rtf": true,

	// JavaScript/TypeScript
	".js":  true,
	".jsx": true,
	".ts":  true,
	".tsx": true,

	// Python
	".py": true,

	// Go
	".go": true,

	// Java/Kotlin
	".java": true,
	".kt":   true,

	// C/C++
	".c":   true,
	".cpp": true,
	".h":   true,

	// Ruby
	".rb": true,

	// PHP
	".php": true,

	// Swift
	".swift": true,

	// Rust
	".rs": true,

	// Shell
	".sh":   true,
	".bash": true,

	// Web
	".html": true,
	".css":  true,
	".scss": true,

	// Config
	".json": true,
	".yaml": true,
	".yml":  true,
	".xml":  true,

	// Database
	".sql": true,

	// Documents (metadata only)
	".pdf":  true,
	".docx": true,
	".xlsx": true,
	".pptx": true,
	".doc":  true,
	".xls":  true,
	".ppt":  true,

	// Data files (metadata only)
	".csv": true,

	// Logs (metadata only)
	".log": true,

	// Media files (metadata only)
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mp3":  true,
	".wav":  true,
	".flac": true,

	// Images (metadata only)
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".svg":  true,
	".webp": true,
}

func IsAllowedExtension(ext string) bool {
	return AllowedExtensions[ext]
}

// IgnoredContentExtensions contains file extensions that should be indexed
// (for metadata) but their content should NOT be read/indexed
var IgnoredContentExtensions = map[string]bool{
	// Documents
	".pdf":  true,
	".docx": true,
	".xlsx": true,
	".pptx": true,
	".doc":  true,
	".xls":  true,
	".ppt":  true,

	// Data files
	".csv": true,

	// Logs
	".log": true,

	// Media files
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mp3":  true,
	".wav":  true,
	".flac": true,

	// Images
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".bmp":  true,
	".svg":  true,
	".webp": true,
}

func ShouldIgnoreContent(ext string) bool {
	return IgnoredContentExtensions[ext]
}