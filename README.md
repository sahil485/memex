# Memex

Accurate, efficient file search for MacOS with a beautiful native desktop interface.

## Features

- ⚡ **Lightning-fast search** powered by MeiliSearch
- 🖥️ **Native desktop app** built with Wails (Go + TypeScript)
- ⌨️ **Keyboard-first** navigation (Cmd+K to search)
- 📁 **Smart file detection** with color-coded type indicators
- 🎯 **Relevance scoring** to find what you need quickly
- 🎨 **Beautiful UI** with smooth animations

## Quick Start

### 1. Install Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install frontend dependencies
cd frontend
npm install
cd ..
```

### 3. Initialize MeiliSearch

```bash
go run cmd/cli/main.go init
```

### 4. Index Your Files

```bash
go run cmd/cli/main.go index /path/to/your/files
```

### 5. Run the App

**Important:** Do NOT run `npm run dev` in the frontend folder. Always use Wails:

```bash
# Easy way - use the helper script
./run.sh

# Or manually
export PATH=$PATH:$(go env GOPATH)/bin
wails dev
```

Or build for production:

```bash
wails build
open build/bin/Memex.app
```

## Usage

### Desktop App

1. Launch Memex
2. Press **Cmd+K** (or **Ctrl+K**) to open search
3. Type your query
4. Use **↑/↓** arrows to navigate results
5. Press **Enter** to open the file in your default editor
6. Press **Esc** to close

### CLI

```bash
# Initialize MeiliSearch
memex-cli init

# Index files
memex-cli index /path/to/directory

# Search from command line
memex-cli search "your query"

# Clear the index
memex-cli clear-index
```

## Architecture

```
┌─────────────────────────────────────┐
│         Memex Desktop App           │
│  ┌──────────────────────────────┐  │
│  │   TypeScript Frontend (Vite) │  │
│  │  - Search Modal UI           │  │
│  │  - Keyboard Navigation       │  │
│  └──────────────────────────────┘  │
│              ↕ Wails Bindings       │
│  ┌──────────────────────────────┐  │
│  │     Go Backend               │  │
│  │  - Search API                │  │
│  │  - File Opening              │  │
│  │  - MeiliSearch Client        │  │
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
              ↕
    ┌──────────────────┐
    │   MeiliSearch    │
    │   Search Engine  │
    └──────────────────┘
```

## Technology Stack

- **Backend**: Go with MeiliSearch client
- **Frontend**: TypeScript + Vite
- **Desktop Framework**: Wails v2
- **Search Engine**: MeiliSearch
- **UI**: Native web technologies with custom styling

## Project Structure

```
memex/
├── main.go              # Wails app entry point
├── app.go               # Go backend with Wails bindings
├── cmd/cli/             # CLI tool
│   └── main.go
├── frontend/            # TypeScript frontend
│   ├── src/
│   │   ├── components/  # UI components
│   │   ├── services/    # API wrappers
│   │   └── wailsjs/     # Generated bindings
│   └── package.json
├── pkg/                 # Go packages
│   ├── client/          # MeiliSearch client
│   ├── search/          # Search logic
│   ├── indexer/         # File indexing
│   └── config/          # Configuration
└── wails.json           # Wails configuration
```

## Building

See [BUILD.md](BUILD.md) for detailed build instructions.

Quick build:
```bash
wails build
```

## Development

```bash
# Run with hot reload
wails dev

# Type check frontend
cd frontend && npm run type-check

# Run tests (when added)
go test ./...
```

## Configuration

MeiliSearch runs on `localhost:58273` by default. Configuration is stored in:
- **MeiliSearch data**: `~/.memex/`
- **Index name**: `files`

## Contributing

Contributions welcome! Please feel free to submit issues and pull requests.

## License

MIT 
