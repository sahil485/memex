# Building Memex

This guide explains how to build and run Memex, a native desktop search application built with Wails, Go, and TypeScript.

## Prerequisites

### 1. Install Go (1.21 or later)
```bash
brew install go
```

### 2. Install Node.js (18 or later)
```bash
brew install node
```

### 3. Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### 4. Install Wails Dependencies (macOS)
```bash
# Wails requires Xcode Command Line Tools
xcode-select --install
```

## Building the Application

### Development Mode

1. **Install frontend dependencies**:
```bash
cd frontend
npm install
cd ..
```

2. **Run in development mode** (with live reload):
```bash
wails dev
```

This will:
- Start the Go backend
- Start the Vite dev server for the frontend
- Open the application window
- Enable hot reload for both frontend and backend changes

### Production Build

Build a production-ready application:

```bash
wails build
```

The compiled application will be in the `build/bin/` directory.

For a production build with optimizations:
```bash
wails build -clean -ldflags "-s -w"
```

Flags:
- `-clean`: Clean the build directory before building
- `-ldflags "-s -w"`: Strip debug information to reduce binary size

## Running the Application

### From Development Build
```bash
wails dev
```

### From Production Build
```bash
./build/bin/memex
```

Or on macOS:
```bash
open build/bin/Memex.app
```

## Project Structure

```
memex/
├── main.go              # Wails application entry point
├── app.go               # Application logic and Go↔JS bindings
├── wails.json           # Wails configuration
├── go.mod               # Go dependencies
├── frontend/            # Frontend TypeScript/Vite project
│   ├── src/
│   │   ├── main.ts                    # Frontend entry point
│   │   ├── components/
│   │   │   └── SearchModal.ts         # Search UI component
│   │   ├── services/
│   │   │   └── search.ts              # Wails bindings wrapper
│   │   ├── types/
│   │   │   └── search.ts              # TypeScript types
│   │   ├── wailsjs/                   # Auto-generated Wails bindings
│   │   └── styles/
│   │       └── global.css
│   ├── package.json
│   ├── vite.config.ts
│   └── tsconfig.json
├── pkg/                 # Go packages
│   ├── client/          # MeiliSearch client
│   ├── search/          # Search functionality
│   ├── indexer/         # File indexing
│   ├── config/          # Configuration
│   └── types/           # Go types
├── cmd/cli/             # CLI version (separate from GUI)
│   └── main.go
└── build/               # Build output (generated)
```

## CLI Tool

The CLI tool is separate from the GUI application. To build and use it:

```bash
# Build CLI
go build -o memex-cli ./cmd/cli

# Run CLI commands
./memex-cli init           # Initialize MeiliSearch
./memex-cli index <path>   # Index files
./memex-cli search <query> # Search from command line
```

## Regenerating Wails Bindings

If you modify the Go structs or methods in `app.go`, regenerate the TypeScript bindings:

```bash
wails generate module
```

This updates the files in `frontend/src/wailsjs/`.

## Troubleshooting

### MeiliSearch Not Running

The application requires MeiliSearch to be running. Start it with:
```bash
./memex-cli init
```

Or manually:
```bash
cd ~/.memex
./meilisearch --http-addr localhost:58273
```

### Frontend Build Errors

Clear and reinstall dependencies:
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
cd ..
```

### Wails Build Errors

Clean the build cache:
```bash
wails build -clean
```

### Port Already in Use

If port 58273 is in use, stop the existing MeiliSearch instance:
```bash
pkill meilisearch
```

## Distribution

### macOS

Create a DMG for distribution:
```bash
wails build -clean
# The .app bundle will be in build/bin/
# You can create a DMG using hdiutil or create-dmg
```

### Code Signing (macOS)

For distribution outside the App Store:
```bash
wails build -clean -ldflags "-s -w" -codesign
```

## Development Tips

1. **Use Wails Dev Mode**: The hot reload is very fast and makes development smooth

2. **Check Console Logs**: Open the Dev Tools in the app (macOS: Cmd+Option+I)

3. **Backend Logs**: Go logs appear in the terminal where you ran `wails dev`

4. **Frontend Logs**: JavaScript logs appear in the browser dev tools

5. **Type Safety**: The Wails bindings are fully typed - TypeScript will catch errors

## Next Steps

- **Add Tests**: Consider adding Go tests for backend and Vitest for frontend
- **CI/CD**: Set up GitHub Actions to build releases automatically
- **Icons**: Customize the app icon in `build/appicon.png`
- **Menu**: Add native menus using Wails menu API
