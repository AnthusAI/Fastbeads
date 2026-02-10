#!/bin/bash
# Interactive compaction workflow
# Run this manually when you want to compact old issues

set -e

echo "=== BD Compaction Workflow ==="
echo "Date: $(date)"
echo

# Check API key
if [ -z "$ANTHROPIC_API_KEY" ]; then
  echo "❌ Error: ANTHROPIC_API_KEY not set"
  echo
  echo "Set your API key:"
  echo "  export ANTHROPIC_API_KEY='sk-ant-...'"
  echo
  exit 1
fi

# Check fbd is installed
if ! command -v fbd &> /dev/null; then
  echo "❌ Error: fbd command not found"
  echo
  echo "Install fbd:"
  echo "  curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash"
  echo
  exit 1
fi

# Preview candidates
echo "--- Preview Tier 1 Candidates ---"
fbd admin compact --dry-run --all

echo
read -p "Proceed with Tier 1 compaction? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "--- Running Tier 1 Compaction ---"
  fbd admin compact --all
  echo "✅ Tier 1 compaction complete"
else
  echo "⏭️  Skipping Tier 1"
fi

# Preview Tier 2
echo
echo "--- Preview Tier 2 Candidates ---"
fbd admin compact --dry-run --all --tier 2

echo
read -p "Proceed with Tier 2 compaction? (y/N) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "--- Running Tier 2 Compaction ---"
  fbd admin compact --all --tier 2
  echo "✅ Tier 2 compaction complete"
else
  echo "⏭️  Skipping Tier 2"
fi

# Show stats
echo
echo "--- Final Statistics ---"
fbd admin compact --stats

echo
echo "=== Compaction Complete ==="
echo
echo "Next steps:"
echo "  1. Review compacted issues: fbd list --json | jq '.[] | select(.compaction_level > 0)'"
echo "  2. Commit changes: git add .beads/issues.jsonl issues.db && git commit -m 'Compact old issues'"
echo "  3. Push to remote: git push"
