# Memex Frontend

A TypeScript-based search interface for Memex, powered by Wails and MeiliSearch.

## Overview

This is the frontend for the Memex desktop application. It's built with:
- **TypeScript** - Type-safe JavaScript
- **Vite** - Fast build tool and dev server
- **Wails** - Go + Web frontend framework

The frontend communicates with the Go backend via Wails bindings, providing a native desktop experience.

## Features

- **Fast Search**: Instant search results powered by MeiliSearch
- **Keyboard Navigation**: Full keyboard support (Cmd+K, arrows, Enter, Esc)
- **Native File Opening**: Opens files in your default editor via Go backend
- **File Type Detection**: Color-coded badges for different file types
- **Relevance Scoring**: See match quality for each result
- **Modern UI**: Clean, beautiful interface with smooth animations

## Development

See [BUILD.md](../BUILD.md) in the root directory for full build instructions.

Quick start:
```bash
npm install
npm run dev  # Standalone mode
# OR
cd .. && wails dev  # Full integration with Go backend
```

## Project Structure

```
src/
├── components/SearchModal.ts   # Main search UI
├── services/search.ts          # Wails bindings wrapper
├── types/search.ts             # TypeScript types
├── wailsjs/                    # Auto-generated bindings
└── main.ts                     # Entry point
```

## Documentation

For detailed documentation, see the root [README.md](../README.md) and [BUILD.md](../BUILD.md).
