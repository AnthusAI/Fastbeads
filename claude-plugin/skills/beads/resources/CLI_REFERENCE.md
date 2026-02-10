# CLI Command Reference

**For:** AI agents and developers using fbd command-line interface
**Version:** 0.47.1+

## Quick Navigation

- [Health & Status](#health--status)
- [Basic Operations](#basic-operations)
- [Issue Management](#issue-management)
- [Dependencies & Labels](#dependencies--labels)
- [Filtering & Search](#filtering--search)
- [Visualization](#visualization)
- [Advanced Operations](#advanced-operations)
- [Database Management](#database-management)

## Health & Status

### Doctor (Start Here for Problems)

```bash
# Basic health check
fbd doctor                      # Check installation health
fbd doctor --json               # Machine-readable output

# Fix issues
fbd doctor --fix                # Auto-fix with confirmation
fbd doctor --fix --yes          # Auto-fix without confirmation
fbd doctor --dry-run            # Preview what --fix would do

# Deep validation
fbd doctor --deep               # Full graph integrity validation

# Performance diagnostics
fbd doctor --perf               # Run performance diagnostics
fbd doctor --output diag.json   # Export diagnostics to file

# Specific checks
fbd doctor --check=pollution              # Detect test issues
fbd doctor --check=pollution --clean      # Delete test issues

# Recovery modes
fbd doctor --fix --source=jsonl           # Rebuild DB from JSONL
fbd doctor --fix --force                  # Force repair on corrupted DB
```

### Status Overview

```bash
# Quick database snapshot (like git status for issues)
fbd status                      # Summary with activity
fbd status --json               # JSON format
fbd status --no-activity        # Skip git activity (faster)
fbd status --assigned           # Show issues assigned to you
fbd stats                       # Alias for fbd status
```

### Prime (AI Context)

```bash
# Output AI-optimized workflow context
fbd prime                       # Auto-detects MCP vs CLI mode
fbd prime --full                # Force full CLI output
fbd prime --mcp                 # Force minimal MCP output
fbd prime --stealth             # No git operations mode
fbd prime --export              # Dump default content for customization
```

**Customization:** Place `.beads/PRIME.md` to override default output.

## Basic Operations

### Check Status

```bash
# Check database path and daemon status
fbd info --json

# Example output:
# {
#   "database_path": "/path/to/.beads/beads.db",
#   "issue_prefix": "fbd",
#   "daemon_running": true
# }
```

### Find Work

```bash
# Find ready work (no blockers)
fbd ready --json
fbd list --ready --json                        # Same, integrated into list (v0.47.1+)

# Find blocked work
fbd blocked --json                             # Show all blocked issues
fbd blocked --parent bd-epic --json            # Blocked descendants of epic

# Find molecules waiting on gates for resume (v0.47.0+)
fbd ready --gated --json                       # Gate-resume discovery

# Find stale issues (not updated recently)
fbd stale --days 30 --json                    # Default: 30 days
fbd stale --days 90 --status in_progress --json  # Filter by status
fbd stale --limit 20 --json                   # Limit results
```

## Issue Management

### Create Issues

```bash
# Basic creation
# IMPORTANT: Always quote titles and descriptions with double quotes
fbd create "Issue title" -t bug|feature|task -p 0-4 -d "Description" --json

# Create with explicit ID (for parallel workers)
fbd create "Issue title" --id worker1-100 -p 1 --json

# Create with labels (--labels or --label work)
fbd create "Issue title" -t bug -p 1 -l bug,critical --json
fbd create "Issue title" -t bug -p 1 --label bug,critical --json

# Examples with special characters (all require quoting):
fbd create "Fix: auth doesn't validate tokens" -t bug -p 1 --json
fbd create "Add support for OAuth 2.0" -d "Implement RFC 6749 (OAuth 2.0 spec)" --json

# Create multiple issues from markdown file
fbd create -f feature-plan.md --json

# Create epic with hierarchical child tasks
fbd create "Auth System" -t epic -p 1 --json         # Returns: bd-a3f8e9
fbd create "Login UI" -p 1 --json                     # Auto-assigned: bd-a3f8e9.1
fbd create "Backend validation" -p 1 --json           # Auto-assigned: bd-a3f8e9.2
fbd create "Tests" -p 1 --json                        # Auto-assigned: bd-a3f8e9.3

# Create and link discovered work (one command)
fbd create "Found bug" -t bug -p 1 --deps discovered-from:<parent-id> --json

# Create with external reference (v0.9.2+)
fbd create "Fix login" -t bug -p 1 --external-ref "gh-123" --json  # Short form
fbd create "Fix login" -t bug -p 1 --external-ref "https://github.com/org/repo/issues/123" --json  # Full URL
fbd create "Jira task" -t task -p 1 --external-ref "jira-PROJ-456" --json  # Custom prefix

# Preview creation without side effects (v0.47.0+)
fbd create "Issue title" -t task -p 1 --dry-run --json  # Shows what would be created
```

### Quick Capture (q)

```bash
# Create issue and output only the ID (for scripting)
fbd q "Fix login bug"                          # Outputs: bd-a1b2
fbd q "Task" -t task -p 1                      # With type and priority
fbd q "Bug" -t bug -l critical                 # With labels

# Scripting examples
ISSUE=$(fbd q "New feature")                   # Capture ID in variable
fbd q "Task" | xargs fbd show                   # Pipe to other commands
```

### Update Issues

```bash
# Update one or more issues
fbd update <id> [<id>...] --status in_progress --json
fbd update <id> [<id>...] --priority 1 --json

# Update external reference (v0.9.2+)
fbd update <id> --external-ref "gh-456" --json           # Short form
fbd update <id> --external-ref "jira-PROJ-789" --json    # Custom prefix

# Edit issue fields in $EDITOR (HUMANS ONLY - not for agents)
# NOTE: This command is intentionally NOT exposed via the MCP server
# Agents should use 'fbd update' with field-specific parameters instead
fbd edit <id>                    # Edit description
fbd edit <id> --title            # Edit title
fbd edit <id> --design           # Edit design notes
fbd edit <id> --notes            # Edit notes
fbd edit <id> --acceptance       # Edit acceptance criteria
```

### Close/Reopen Issues

```bash
# Complete work (supports multiple IDs)
fbd close <id> [<id>...] --reason "Done" --json

# Reopen closed issues (supports multiple IDs)
fbd reopen <id> [<id>...] --reason "Reopening" --json
```

### View Issues

```bash
# Show dependency tree
fbd dep tree <id>

# Get issue details (supports multiple IDs)
fbd show <id> [<id>...] --json
```

### Comments

```bash
# List comments on an issue
fbd comments bd-123                            # Human-readable
fbd comments bd-123 --json                     # JSON format

# Add a comment
fbd comments add bd-123 "This is a comment"
fbd comments add bd-123 -f notes.txt           # From file
```

## Dependencies & Labels

### Dependencies

```bash
# Link discovered work (old way - two commands)
fbd dep add <discovered-id> <parent-id> --type discovered-from

# Create and link in one command (new way - preferred)
fbd create "Issue title" -t bug -p 1 --deps discovered-from:<parent-id> --json
```

### Labels

```bash
# Label management (supports multiple IDs)
fbd label add <id> [<id>...] <label> --json
fbd label remove <id> [<id>...] <label> --json
fbd label list <id> --json
fbd label list-all --json
```

## Filtering & Search

### Basic Filters

```bash
# Filter by status, priority, type
fbd list --status open --priority 1 --json               # Status and priority
fbd list --assignee alice --json                         # By assignee
fbd list --type bug --json                               # By issue type
fbd list --id bd-123,bd-456 --json                       # Specific IDs
```

### Label Filters

```bash
# Labels (AND: must have ALL)
fbd list --label bug,critical --json

# Labels (OR: has ANY)
fbd list --label-any frontend,backend --json
```

### Search Command

```bash
# Full-text search across title, description, and ID
fbd search "authentication bug"                          # Basic search
fbd search "login" --status open --json                  # With status filter
fbd search "database" --label backend --limit 10         # With label and limit
fbd search "bd-5q"                                       # Search by partial ID

# Find beads issue by external reference
fbd list --json | jq -r '.[] | select(.external_ref == "gh-123") | .id'

# Filtered search
fbd search "security" --priority-min 0 --priority-max 2  # Priority range
fbd search "bug" --created-after 2025-01-01              # Date filter
fbd search --query "refactor" --assignee alice           # By assignee

# Sorted results
fbd search "bug" --sort priority                         # Sort by priority
fbd search "task" --sort created --reverse               # Reverse chronological
fbd search "feature" --long                              # Detailed multi-line output
```

### Text Search (via list)

```bash
# Title search (substring)
fbd list --title "auth" --json

# Pattern matching (case-insensitive substring)
fbd list --title-contains "auth" --json                  # Search in title
fbd list --desc-contains "implement" --json              # Search in description
fbd list --notes-contains "TODO" --json                  # Search in notes
```

### Date Range Filters

```bash
# Date range filters (YYYY-MM-DD or RFC3339)
fbd list --created-after 2024-01-01 --json               # Created after date
fbd list --created-before 2024-12-31 --json              # Created before date
fbd list --updated-after 2024-06-01 --json               # Updated after date
fbd list --updated-before 2024-12-31 --json              # Updated before date
fbd list --closed-after 2024-01-01 --json                # Closed after date
fbd list --closed-before 2024-12-31 --json               # Closed before date
```

### Empty/Null Checks

```bash
# Empty/null checks
fbd list --empty-description --json                      # Issues with no description
fbd list --no-assignee --json                            # Unassigned issues
fbd list --no-labels --json                              # Issues with no labels
```

### Priority Ranges

```bash
# Priority ranges
fbd list --priority-min 0 --priority-max 1 --json        # P0 and P1 only
fbd list --priority-min 2 --json                         # P2 and below
```

### Combine Filters

```bash
# Combine multiple filters
fbd list --status open --priority 1 --label-any urgent,critical --no-assignee --json
```

## Visualization

### Graph (Dependency Visualization)

```bash
# Show dependency graph for an issue
fbd graph bd-123                               # ASCII box format (default)
fbd graph bd-123 --compact                     # Tree format, one line per issue

# Show graph for epic (includes all children)
fbd graph bd-epic

# Show all open issues grouped by component
fbd graph --all
```

**Display formats:**
- `--box` (default): ASCII boxes showing layers, more detailed
- `--compact`: Tree format, one line per issue, more scannable

**Graph interpretation:**
- Layer 0 / leftmost = no dependencies (can start immediately)
- Higher layers depend on lower layers
- Nodes in the same layer can run in parallel

**Status icons:** ○ open  ◐ in_progress  ● blocked  ✓ closed  ❄ deferred

## Global Flags

Global flags work with any fbd command and must appear **before** the subcommand.

### Sandbox Mode

**Auto-detection (v0.21.1+):** fbd automatically detects sandboxed environments and enables sandbox mode.

When detected, you'll see: `ℹ️  Sandbox detected, using direct mode`

**Manual override:**

```bash
# Explicitly enable sandbox mode
fbd --sandbox <command>

# Equivalent to combining these flags:
fbd --no-daemon --no-auto-flush --no-auto-import <command>
```

**What it does:**
- Disables daemon (uses direct SQLite mode)
- Disables auto-export to JSONL
- Disables auto-import from JSONL

**When to use:** Sandboxed environments where daemon can't be controlled (permission restrictions), or when auto-detection doesn't trigger.

### Staleness Control

```bash
# Skip staleness check (emergency escape hatch)
fbd --allow-stale <command>

# Example: access database even if out of sync with JSONL
fbd --allow-stale ready --json
fbd --allow-stale list --status open --json
```

**Shows:** `⚠️  Staleness check skipped (--allow-stale), data may be out of sync`

**⚠️ Caution:** May show stale or incomplete data. Use only when stuck and other options fail.

### Force Import

```bash
# Force metadata update even when DB appears synced
fbd import --force -i .beads/issues.jsonl
```

**When to use:** `fbd import` reports "0 created, 0 updated" but staleness errors persist.

**Shows:** `Metadata updated (database already in sync with JSONL)`

### Other Global Flags

```bash
# JSON output for programmatic use
fbd --json <command>

# Force direct mode (bypass daemon)
fbd --no-daemon <command>

# Disable auto-sync
fbd --no-auto-flush <command>    # Disable auto-export to JSONL
fbd --no-auto-import <command>   # Disable auto-import from JSONL

# Custom database path
fbd --db /path/to/.beads/beads.db <command>

# Custom actor for audit trail
fbd --actor alice <command>
```

**See also:**
- [TROUBLESHOOTING.md - Sandboxed environments](TROUBLESHOOTING.md#sandboxed-environments-codex-claude-code-etc) for detailed sandbox troubleshooting
- [DAEMON.md](DAEMON.md) for daemon mode details

## Advanced Operations

### Cleanup

```bash
# Clean up closed issues (bulk deletion)
fbd admin cleanup --force --json                                   # Delete ALL closed issues
fbd admin cleanup --older-than 30 --force --json                   # Delete closed >30 days ago
fbd admin cleanup --dry-run --json                                 # Preview what would be deleted
fbd admin cleanup --older-than 90 --cascade --force --json         # Delete old + dependents
```

### Duplicate Detection & Merging

```bash
# Find and merge duplicate issues
fbd duplicates                                          # Show all duplicates
fbd duplicates --auto-merge                             # Automatically merge all
fbd duplicates --dry-run                                # Preview merge operations

# Merge specific duplicate issues
fbd merge <source-id...> --into <target-id> --json      # Consolidate duplicates
fbd merge bd-42 bd-43 --into bd-41 --dry-run            # Preview merge
```

### Compaction (Memory Decay)

```bash
# Agent-driven compaction
fbd admin compact --analyze --json                           # Get candidates for review
fbd admin compact --analyze --tier 1 --limit 10 --json       # Limited batch
fbd admin compact --apply --id bd-42 --summary summary.txt   # Apply compaction
fbd admin compact --apply --id bd-42 --summary - < summary.txt  # From stdin
fbd admin compact --stats --json                             # Show statistics

# Legacy AI-powered compaction (requires ANTHROPIC_API_KEY)
fbd admin compact --auto --dry-run --all                     # Preview
fbd admin compact --auto --all --tier 1                      # Auto-compact tier 1

# Restore compacted issue from git history
fbd restore <id>  # View full history at time of compaction
```

### Rename Prefix

```bash
# Rename issue prefix (e.g., from 'knowledge-work-' to 'kw-')
fbd rename-prefix kw- --dry-run  # Preview changes
fbd rename-prefix kw- --json     # Apply rename
```

## Database Management

### Import/Export

```bash
# Import issues from JSONL
fbd import -i .beads/issues.jsonl --dry-run      # Preview changes
fbd import -i .beads/issues.jsonl                # Import and update issues
fbd import -i .beads/issues.jsonl --dedupe-after # Import + detect duplicates

# Handle missing parents during import
fbd import -i issues.jsonl --orphan-handling allow      # Default: import orphans without validation
fbd import -i issues.jsonl --orphan-handling resurrect  # Auto-resurrect deleted parents as tombstones
fbd import -i issues.jsonl --orphan-handling skip       # Skip orphans with warning
fbd import -i issues.jsonl --orphan-handling strict     # Fail if parent is missing

# Configure default orphan handling behavior
fbd config set import.orphan_handling "resurrect"
fbd sync  # Now uses resurrect mode by default
```

**Orphan handling modes:**

- **`allow` (default)** - Import orphaned children without parent validation. Most permissive, ensures no data loss even if hierarchy is temporarily broken.
- **`resurrect`** - Search JSONL history for deleted parents and recreate them as tombstones (Status=Closed, Priority=4). Preserves hierarchy with minimal data. Dependencies are also resurrected on best-effort basis.
- **`skip`** - Skip orphaned children with warning. Partial import succeeds but some issues are excluded.
- **`strict`** - Fail import immediately if a child's parent is missing. Use when database integrity is critical.

**When to use:**
- Use `allow` (default) for daily imports and auto-sync
- Use `resurrect` when importing from databases with deleted parents
- Use `strict` for controlled imports requiring guaranteed parent existence
- Use `skip` rarely - only for selective imports

See [CONFIG.md](CONFIG.md#example-import-orphan-handling) and [TROUBLESHOOTING.md](TROUBLESHOOTING.md#import-fails-with-missing-parent-errors) for more details.

### Migration

```bash
# Migrate databases after version upgrade
fbd migrate                                             # Detect and migrate old databases
fbd migrate --dry-run                                   # Preview migration
fbd migrate --cleanup --yes                             # Migrate and remove old files

# AI-supervised migration (check before running fbd migrate)
fbd migrate --inspect --json                            # Show migration plan for AI agents
fbd info --schema --json                                # Get schema, tables, config, sample IDs
```

**Migration workflow for AI agents:**

1. Run `--inspect` to see pending migrations and warnings
2. Check for `missing_config` (like issue_prefix)
3. Review `invariants_to_check` for safety guarantees
4. If warnings exist, fix config issues first
5. Then run `fbd migrate` safely

**Migration safety invariants:**

- **required_config_present**: Ensures issue_prefix and schema_version are set
- **foreign_keys_valid**: No orphaned dependencies or labels
- **issue_count_stable**: Issue count doesn't decrease unexpectedly

These invariants prevent data loss and would have caught issues like GH #201 (missing issue_prefix after migration).

### Daemon Management

See [docs/DAEMON.md](DAEMON.md) for complete daemon management reference.

```bash
# List all running daemons
fbd daemons list --json

# Check health (version mismatches, stale sockets)
fbd daemons health --json

# Stop/restart specific daemon
fbd daemons stop /path/to/workspace --json
fbd daemons restart 12345 --json  # By PID

# View daemon logs
fbd daemons logs /path/to/workspace -n 100
fbd daemons logs 12345 -f  # Follow mode

# Stop all daemons
fbd daemons killall --json
fbd daemons killall --force --json  # Force kill if graceful fails
```

### Sync Operations

```bash
# Manual sync (force immediate export/import/commit/push)
fbd sync

# What it does:
# 1. Export pending changes to JSONL
# 2. Commit to git
# 3. Pull from remote
# 4. Import any updates
# 5. Push to remote

# Resolve JSONL merge conflict markers (v0.47.0+)
fbd resolve-conflicts                          # Resolve in mechanical mode
fbd resolve-conflicts --dry-run --json         # Preview resolution
# Mechanical mode rules: updated_at wins, closed beats open, higher priority wins
```

## Issue Types

- `bug` - Something broken that needs fixing
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature composed of multiple issues (supports hierarchical children)
- `chore` - Maintenance work (dependencies, tooling)

**Hierarchical children:** Epics can have child issues with dotted IDs (e.g., `bd-a3f8e9.1`, `bd-a3f8e9.2`). Children are auto-numbered sequentially. Up to 3 levels of nesting supported.

## Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (nice-to-have features, minor bugs)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

## Dependency Types

- `blocks` - Hard dependency (issue X blocks issue Y)
- `related` - Soft relationship (issues are connected)
- `parent-child` - Epic/subtask relationship
- `discovered-from` - Track issues discovered during work

Only `blocks` dependencies affect the ready work queue.

**Note:** When creating an issue with a `discovered-from` dependency, the new issue automatically inherits the parent's `source_repo` field.

## External References

The `--external-ref` flag (v0.9.2+) links beads issues to external trackers:

- Supports short form (`gh-123`) or full URL (`https://github.com/...`)
- Portable via JSONL - survives sync across machines
- Custom prefixes work for any tracker (`jira-PROJ-456`, `linear-789`)

## Output Formats

### JSON Output (Recommended for Agents)

Always use `--json` flag for programmatic use:

```bash
# Single issue
fbd show bd-42 --json

# List of issues
fbd ready --json

# Operation result
fbd create "Issue" -p 1 --json
```

### Human-Readable Output

Default output without `--json`:

```bash
fbd ready
# bd-42  Fix authentication bug  [P1, bug, in_progress]
# bd-43  Add user settings page  [P2, feature, open]
```

## Common Patterns for AI Agents

### Claim and Complete Work

```bash
# 1. Find available work
fbd ready --json

# 2. Claim issue
fbd update bd-42 --status in_progress --json

# 3. Work on it...

# 4. Close when done
fbd close bd-42 --reason "Implemented and tested" --json
```

### Discover and Link Work

```bash
# While working on bd-100, discover a bug

# Old way (two commands):
fbd create "Found auth bug" -t bug -p 1 --json  # Returns bd-101
fbd dep add bd-101 bd-100 --type discovered-from

# New way (one command):
fbd create "Found auth bug" -t bug -p 1 --deps discovered-from:bd-100 --json
```

### Batch Operations

```bash
# Update multiple issues at once
fbd update bd-41 bd-42 bd-43 --priority 0 --json

# Close multiple issues
fbd close bd-41 bd-42 bd-43 --reason "Batch completion" --json

# Add label to multiple issues
fbd label add bd-41 bd-42 bd-43 urgent --json
```

### Session Workflow

```bash
# Start of session
fbd ready --json  # Find work

# During session
fbd create "..." -p 1 --json
fbd update bd-42 --status in_progress --json
# ... work ...

# End of session (IMPORTANT!)
fbd sync  # Force immediate sync, bypass debounce
```

**ALWAYS run `fbd sync` at end of agent sessions** to ensure changes are committed/pushed immediately.

## See Also

- [AGENTS.md](../AGENTS.md) - Main agent workflow guide
- [DAEMON.md](DAEMON.md) - Daemon management and event-driven mode
- [GIT_INTEGRATION.md](GIT_INTEGRATION.md) - Git workflows and merge strategies
- [LABELS.md](../LABELS.md) - Label system guide
- [README.md](../README.md) - User documentation
