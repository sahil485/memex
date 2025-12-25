#!/bin/bash

set -e  # Exit on error

echo "Setting up memex CLI..."

# Create the ~/.memex directory if it doesn't exist
if [ ! -d "$HOME/.memex" ]; then
    mkdir -p "$HOME/.memex"
    echo "✓ Created ~/.memex directory"
fi

# Download Meilisearch binary to ~/.memex if not present
if [ ! -f "$HOME/.memex/meilisearch" ]; then
    echo "Downloading Meilisearch..."
    curl -L https://github.com/meilisearch/meilisearch/releases/latest/download/meilisearch-macos-apple-silicon -o "$HOME/.memex/meilisearch"
    chmod +x "$HOME/.memex/meilisearch"
    echo "✓ Meilisearch downloaded"
else
    echo "✓ Meilisearch already present"
fi

go mod download
go build -o memex main.go

chmod +x memex
sudo mv memex /usr/local/bin/

# Create launchd plist for auto-starting Meilisearch
echo "Setting up Meilisearch to run on login..."
PLIST_PATH="$HOME/Library/LaunchAgents/com.memex.meilisearch.plist"

cat > "$PLIST_PATH" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.memex.meilisearch</string>
    <key>ProgramArguments</key>
    <array>
        <string>$HOME/.memex/meilisearch</string>
        <string>--http-addr</string>
        <string>127.0.0.1:58273</string>
        <string>--db-path</string>
        <string>$HOME/.memex/data.ms</string>
        <string>--no-analytics</string>
    </array>
    <key>WorkingDirectory</key>
    <string>$HOME/.memex</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>$HOME/.memex/meilisearch.log</string>
    <key>StandardErrorPath</key>
    <string>$HOME/.memex/meilisearch.error.log</string>
</dict>
</plist>
EOF

launchctl unload "$PLIST_PATH" 2>/dev/null || true
launchctl load "$PLIST_PATH"
launchctl start com.memex.meilisearch

echo "✓ Meilisearch configured to run on login and started"

echo "Waiting for Meilisearch to start..."
for i in {1..10}; do
    if curl -s http://127.0.0.1:58273/health > /dev/null 2>&1; then
        echo "✓ Meilisearch is ready"
        break
    fi
    sleep 1
done

echo "Initializing index..."
go run ./pkg/client/init.go && echo "✓ Index initialized" || echo "Failed to initialize index"
echo "✓ memex built and installed!"
echo "Run with: memex search [options]"