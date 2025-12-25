# Wails Integration Guide

This document explains how the Wails integration works in Memex and how to work with it.

## What is Wails?

Wails is a framework for building desktop applications using Go and web technologies. It's similar to Electron, but uses Go instead of Node.js for the backend, resulting in:

- **Smaller binaries** - No need to bundle Node.js/Chromium
- **Better performance** - Go is fast and efficient
- **Type safety** - Full TypeScript support for Go↔JS communication
- **Native feel** - Uses the system's webview

## How It Works

### Backend (Go)

The Go backend ([app.go](app.go)) exposes methods that can be called from the frontend:

```go
type App struct {
    ctx context.Context
}

// This method can be called from JavaScript
func (a *App) Search(query string, limit int) SearchResponse {
    // Your implementation
}

func (a *App) OpenFile(path string) error {
    // Opens file in default editor
}
```

### Frontend (TypeScript)

The frontend calls Go methods through generated bindings:

```typescript
import { Search, OpenFile } from '../wailsjs/go/main/App';

// Call Go function
const results = await Search("my query", 20);

// Open a file
await OpenFile("/path/to/file");
```

### Auto-Generated Bindings

Wails automatically generates TypeScript bindings for your Go code:

```
frontend/src/wailsjs/
├── go/
│   ├── main/
│   │   └── App.ts          # Methods from App struct
│   └── models.ts           # Go struct types
└── runtime/
    └── runtime.d.ts        # Wails runtime API
```

**Never edit these files manually!** Regenerate them with:

```bash
wails generate module
```

## Data Flow

```
┌──────────────────────────────────────────┐
│           Frontend (TypeScript)          │
│                                          │
│  SearchModal.openSelected() calls:      │
│  await searchService.openFile(path)     │
└──────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────┐
│      Wails Generated Bindings            │
│                                          │
│  OpenFile(path) → IPC Call               │
└──────────────────────────────────────────┘
                    ↓
┌──────────────────────────────────────────┐
│          Backend (Go)                    │
│                                          │
│  func (a *App) OpenFile(path string)    │
│    → exec.Command("open", path)         │
└──────────────────────────────────────────┘
```

## Type Safety

One of Wails' best features is full type safety across the Go↔TypeScript boundary.

### Go Struct

```go
type SearchResult struct {
    ID           string  `json:"id"`
    Path         string  `json:"path"`
    Content      string  `json:"content"`
    Type         string  `json:"type"`
    RankingScore float64 `json:"rankingScore"`
}
```

### Generated TypeScript

```typescript
export namespace main {
    export class SearchResult {
        id: string;
        path: string;
        content: string;
        type: string;
        rankingScore: number;
    }
}
```

TypeScript will catch type errors at compile time!

## Development Workflow

### 1. Make Changes to Go Backend

Edit [app.go](app.go) to add/modify methods:

```go
func (a *App) MyNewMethod(arg string) string {
    return "Hello " + arg
}
```

### 2. Regenerate Bindings

```bash
wails generate module
```

### 3. Use in Frontend

```typescript
import { MyNewMethod } from '../wailsjs/go/main/App';

const result = await MyNewMethod("World");
console.log(result); // "Hello World"
```

### 4. Test with Hot Reload

```bash
wails dev
```

Both frontend and backend changes will hot reload!

## Best Practices

### 1. Keep Methods Simple

Wails methods should be thin wrappers around your business logic:

```go
// Good ✅
func (a *App) Search(query string, limit int) SearchResponse {
    result, err := search.Search(query, int64(limit))
    // Convert and return
}

// Bad ❌ - Complex logic in Wails method
func (a *App) Search(query string, limit int) SearchResponse {
    // 100 lines of search implementation
}
```

### 2. Handle Errors Properly

Return errors in your structs, not as Go errors:

```go
// Good ✅
type SearchResponse struct {
    Hits  []SearchResult `json:"hits"`
    Error string         `json:"error,omitempty"`
}

func (a *App) Search(query string) SearchResponse {
    result, err := search.Search(query)
    if err != nil {
        return SearchResponse{Error: err.Error()}
    }
    return result
}
```

### 3. Use Context for Cancellation

Access the Wails runtime context:

```go
func (a *App) LongRunningTask() {
    ctx := a.ctx // Set during startup
    // Use ctx for cancellation
}
```

### 4. Log from Both Sides

Backend logs go to terminal, frontend logs go to dev tools:

```go
// Go
log.Println("Backend: Starting search")
```

```typescript
// TypeScript
console.log("Frontend: Displaying results");
```

## Common Patterns

### Calling Go from TypeScript

```typescript
import { Search } from '../wailsjs/go/main/App';

try {
    const results = await Search("query", 20);
    if (results.error) {
        console.error("Search failed:", results.error);
    } else {
        // Handle results
    }
} catch (err) {
    console.error("Unexpected error:", err);
}
```

### Progress Updates

For long-running operations, use Wails events:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) IndexDirectory(path string) {
    files := getFiles(path)
    for i, file := range files {
        indexFile(file)
        // Emit progress event
        runtime.EventsEmit(a.ctx, "index-progress", i, len(files))
    }
}
```

```typescript
import { EventsOn } from '../wailsjs/runtime/runtime';

EventsOn("index-progress", (current: number, total: number) => {
    console.log(`Progress: ${current}/${total}`);
});
```

### File Dialogs

Use Wails runtime APIs:

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

func (a *App) ChooseDirectory() (string, error) {
    return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
        Title: "Select Directory to Index",
    })
}
```

## Debugging

### Backend

Go logs appear in the terminal where you ran `wails dev`:

```go
log.Println("Debug:", variable)
```

### Frontend

Open dev tools in the app (Cmd+Option+I on macOS):

```typescript
console.log("Debug:", variable);
```

### Breakpoints

- **Go**: Use Delve or your IDE's debugger
- **TypeScript**: Use browser dev tools

## Performance Tips

### 1. Minimize IPC Calls

Each call to a Go method has overhead. Batch operations when possible:

```typescript
// Bad ❌
for (const file of files) {
    await IndexFile(file);
}

// Good ✅
await IndexFiles(files);
```

### 2. Use Streaming for Large Data

For large result sets, consider pagination or streaming:

```go
func (a *App) SearchPaginated(query string, page int, limit int) SearchResponse {
    offset := page * limit
    return search.Search(query, limit, offset)
}
```

### 3. Cache in Frontend

Cache frequently accessed data in the frontend:

```typescript
class SearchService {
    private cache = new Map<string, SearchResponse>();

    async search(query: string) {
        if (this.cache.has(query)) {
            return this.cache.get(query)!;
        }
        const result = await Search(query, 20);
        this.cache.set(query, result);
        return result;
    }
}
```

## Troubleshooting

### Bindings Not Generated

```bash
# Make sure you're in the project root
wails generate module

# If that fails, try regenerating everything
rm -rf frontend/src/wailsjs
wails generate module
```

### Type Errors

Make sure your Go JSON tags match TypeScript expectations:

```go
type MyStruct struct {
    // This becomes "myField" in TypeScript
    MyField string `json:"myField"`
}
```

### Methods Not Appearing

Make sure:
1. Method is exported (capitalized)
2. Method is on the bound struct (App)
3. You regenerated bindings

### Hot Reload Not Working

1. Check terminal for errors
2. Restart `wails dev`
3. Clear Vite cache: `rm -rf frontend/node_modules/.vite`

## Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails Reference](https://wails.io/docs/reference/introduction)
- [Wails Discord](https://discord.gg/wails)
- [GitHub Repository](https://github.com/wailsapp/wails)
