# Setup Instructions

You're seeing the error because the app needs to run through Wails, not as a standalone web app. Follow these steps:

## 1. Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Make sure your Go bin directory is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
# Add this to your ~/.zshrc or ~/.bashrc to make it permanent
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```

## 2. Install Frontend Dependencies

```bash
cd frontend
npm install
cd ..
```

## 3. Initialize MeiliSearch (if not already done)

```bash
go run cmd/cli/main.go init
```

This will:
- Download MeiliSearch to `~/.memex/`
- Start it on port 58273
- Create the `files` index

## 4. Index Some Files

```bash
go run cmd/cli/main.go index /path/to/your/documents
```

For example:
```bash
go run cmd/cli/main.go index ~/Documents
```

## 5. Run the Wails App

```bash
wails dev
```

This will:
- Compile the Go backend
- Start the Vite dev server
- Launch the native desktop app
- Enable hot reload for both frontend and backend

## Troubleshooting

### Error: "wails: command not found"

Make sure `$(go env GOPATH)/bin` is in your PATH:
```bash
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✓ Go bin in PATH" || echo "✗ Need to add Go bin to PATH"
```

If not found, add it:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Error: "Error performing search"

This means either:
1. You're running the frontend standalone (`npm run dev`) - **Don't do this!** Use `wails dev` instead.
2. MeiliSearch isn't running - Run `go run cmd/cli/main.go init` first
3. No files have been indexed - Run `go run cmd/cli/main.go index <directory>`

### Check if MeiliSearch is Running

```bash
curl http://localhost:58273/health
# Should return: {"status":"available"}
```

If not running, start it:
```bash
cd ~/.memex
./meilisearch --http-addr localhost:58273
```

### Frontend TypeScript Errors

If you see import errors for Wails bindings, regenerate them:
```bash
wails generate module
```

### Clean Build

If things aren't working, try a clean build:
```bash
rm -rf build frontend/dist frontend/node_modules/.vite
cd frontend && npm install && cd ..
wails build -clean
```

## Development Workflow

1. **Make Backend Changes**: Edit Go files (app.go, pkg/*)
2. **Make Frontend Changes**: Edit TypeScript files (frontend/src/*)
3. **Hot Reload**: Both will automatically reload in `wails dev`
4. **If You Add/Change Go Methods**: Run `wails generate module` to update TypeScript bindings

## Common Commands

```bash
# Run dev mode (with hot reload)
wails dev

# Build production app
wails build

# Run CLI commands
go run cmd/cli/main.go init
go run cmd/cli/main.go index <directory>
go run cmd/cli/main.go search <query>

# Type check frontend
cd frontend && npm run type-check

# Build frontend only
cd frontend && npm run build
```

## What NOT to Do

❌ **Don't run `npm run dev` in the frontend directory** - This runs a standalone web server without the Go backend

✅ **Always use `wails dev`** - This runs the full integrated app

## Quick Start (After Setup)

```bash
# Terminal 1: Make sure MeiliSearch is running (or was started by init)
# Terminal 2: Run the app
wails dev
```

Press **Cmd+K** in the app to start searching!
