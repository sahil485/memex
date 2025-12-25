#!/bin/bash

# Memex Development Runner
# This script makes it easy to run the Wails desktop app

# Add Go bin to PATH if not already there
export PATH=$PATH:$(go env GOPATH)/bin

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "❌ Wails is not installed. Installing..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Check if MeiliSearch is running
echo "🔍 Checking if MeiliSearch is running..."
if curl -s http://localhost:58273/health > /dev/null 2>&1; then
    echo "✅ MeiliSearch is running"
else
    echo "⚠️  MeiliSearch is not running"
    echo "   Starting MeiliSearch..."
    if [ -f ~/.memex/meilisearch ]; then
        cd ~/.memex && ./meilisearch --http-addr localhost:58273 > /dev/null 2>&1 &
        sleep 2
        echo "✅ MeiliSearch started"
    else
        echo "❌ MeiliSearch not found. Run: go run cmd/cli/main.go init"
        exit 1
    fi
fi

# Run Wails dev mode
echo "🚀 Starting Memex in development mode..."
echo "   Press Cmd+K in the app to search"
echo ""
wails dev
