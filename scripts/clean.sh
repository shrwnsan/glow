#!/usr/bin/env bash
# Clean up build artifacts and temporary files

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

echo "🧹 Cleaning glow repository..."

# Remove build artifacts
echo "  Removing dist/..."
rm -rf "$REPO_ROOT/dist"

# Remove test binaries (but preserve test suite files like *_test.go)
echo "  Removing glow-test binary..."
rm -f "$REPO_ROOT/glow-test" 2>/dev/null || true

# Remove temporary Go build cache
echo "  Cleaning Go cache..."
go clean -cache -modcache -testcache 2>/dev/null || true

# Remove temporary editor files
echo "  Removing temp files..."
find "$REPO_ROOT" -name "*.swp" -delete 2>/dev/null || true
find "$REPO_ROOT" -name "*~" -delete 2>/dev/null || true
find "$REPO_ROOT" -name ".DS_Store" -delete 2>/dev/null || true

echo "✨ Clean complete!"
echo ""
echo "Run 'go build' to create a new build."
