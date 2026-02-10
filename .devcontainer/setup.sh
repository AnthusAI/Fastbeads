#!/bin/bash
set -e

echo "🔧 Building fbd from source..."
go build -o fbd ./cmd/fbd

echo "📦 Installing fbd globally..."
sudo mv fbd /usr/local/bin/fbd
sudo chmod +x /usr/local/bin/fbd

echo "✅ Verifying fbd installation..."
fbd version

echo "🎯 Initializing fbd (non-interactive)..."
if [ ! -f .beads/beads.db ]; then
  fbd init --quiet
else
  echo "fbd already initialized"
fi

echo "🪝 Installing git hooks..."
if [ -f examples/git-hooks/install.sh ]; then
  bash examples/git-hooks/install.sh
  echo "Git hooks installed successfully"
else
  echo "⚠️  Git hooks installer not found, skipping..."
fi

echo "📚 Installing Go dependencies..."
go mod download

echo "✨ Development environment ready!"
echo "Run 'fbd ready' to see available tasks"
