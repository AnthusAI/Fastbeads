#!/bin/bash
#
# bd-version-check.sh - Automatic fbd upgrade detection for AI agent sessions
#
# This script detects when fbd (beads) has been upgraded and automatically shows
# what changed, helping AI agents adapt their workflows without manual intervention.
#
# FEATURES:
# - Detects fbd version changes by comparing to last-seen version
# - Shows 'fbd info --whats-new' output when upgrade detected
# - Auto-updates git hooks if outdated
# - Persists version in .beads/metadata.json
# - Zero fbd code changes required - works today!
#
# INTEGRATION:
# Add this script to your AI environment's session startup:
#
# Claude Code:
#   Add to .claude/hooks/session-start (if supported)
#   Or manually source at beginning of work
#
# GitHub Copilot:
#   Add to your shell initialization (.bashrc, .zshrc)
#   Or manually run at session start
#
# Cursor:
#   Add to workspace settings or shell init
#
# Generic:
#   source /path/to/bd-version-check.sh
#
# USAGE:
#   # Option 1: Source it (preferred)
#   source examples/startup-hooks/bd-version-check.sh
#
#   # Option 2: Execute it
#   bash examples/startup-hooks/bd-version-check.sh
#
# REQUIREMENTS:
# - fbd (beads) installed and in PATH
# - jq for JSON manipulation
# - .beads directory exists in current project
#

# Exit early if not in a beads project
if [ ! -d ".beads" ]; then
  return 0 2>/dev/null || exit 0
fi

# Check if fbd is installed
if ! command -v fbd &> /dev/null; then
  return 0 2>/dev/null || exit 0
fi

# Check if jq is installed (required for JSON manipulation)
if ! command -v jq &> /dev/null; then
  echo "⚠️  bd-version-check: jq not found. Install jq to enable automatic upgrade detection."
  return 0 2>/dev/null || exit 0
fi

# Get current fbd version
CURRENT_VERSION=$(fbd --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)

if [ -z "$CURRENT_VERSION" ]; then
  # fbd command failed, skip
  return 0 2>/dev/null || exit 0
fi

# Path to metadata file
METADATA_FILE=".beads/metadata.json"

# Initialize metadata.json if it doesn't exist
if [ ! -f "$METADATA_FILE" ]; then
  echo '{"database": "beads.db", "jsonl_export": "beads.jsonl"}' > "$METADATA_FILE"
fi

# Read last-seen version from metadata.json
LAST_VERSION=$(jq -r '.last_bd_version // "unknown"' "$METADATA_FILE" 2>/dev/null)

# Detect version change
if [ "$CURRENT_VERSION" != "$LAST_VERSION" ] && [ "$LAST_VERSION" != "unknown" ]; then
  echo ""
  echo "🔄 fbd upgraded: $LAST_VERSION → $CURRENT_VERSION"
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  # Show what's new
  fbd info --whats-new 2>/dev/null || echo "⚠️  Could not fetch what's new (run 'fbd info --whats-new' manually)"

  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo ""
  echo "💡 Review changes above and adapt your workflow accordingly"
  echo ""
fi

# Check for outdated git hooks (works even if version didn't change)
if fbd hooks list 2>&1 | grep -q "outdated"; then
  echo "🔧 Git hooks outdated. Updating to match fbd v$CURRENT_VERSION..."
  if fbd hooks install 2>/dev/null; then
    echo "✓ Git hooks updated successfully"
  else
    echo "⚠️  Failed to update git hooks. Run 'fbd hooks install' manually."
  fi
  echo ""
fi

# Update metadata.json with current version
# Use a temp file to avoid corruption if jq fails
TEMP_FILE=$(mktemp)
if jq --arg v "$CURRENT_VERSION" '.last_bd_version = $v' "$METADATA_FILE" > "$TEMP_FILE" 2>/dev/null; then
  mv "$TEMP_FILE" "$METADATA_FILE"
else
  # jq failed, clean up temp file
  rm -f "$TEMP_FILE"
fi

# Clean exit for sourcing
return 0 2>/dev/null || exit 0
