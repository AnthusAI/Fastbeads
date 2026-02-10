# CLI Command Reference

**For:** AI agents and developers using fbd command-line interface  
**Version:** 0.21.0+

## Quick Navigation

- [Basic Operations](#basic-operations)
- [Issue Management](#issue-management)
- [Dependencies & Labels](#dependencies--labels)
- [Filtering & Search](#filtering--search)
- [Advanced Operations](#advanced-operations)
- [Molecular Chemistry](#molecular-chemistry)
- [Database Management](#database-management)
- [Editor Integration](#editor-integration)

## Basic Operations

### Check Status

```bash
# Check database path and daemon status
fbd info --json

# Example output:
# {
#   "database_path": "/path/to/.beads/beads.db",
#   "issue_prefix": "fbd",
#   "daemon_running": true,
#   "agent_mail_enabled": false
# }
```

### Find Work

```bash
# Find ready work (no blockers, not already claimed)
fbd ready --json

# Atomically claim an issue from the ready queue
fbd update <id> --claim --json               # Fails if already claimed

# Find stale issues (not updated recently)
fbd stale --days 30 --json                    # Default: 30 days
fbd stale --days 90 --status in_progress --json  # Find abandoned claims
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
fbd create "Implement auth" --spec-id "docs/specs/auth.md" --json

# Create multiple issues from markdown file
fbd create -f feature-plan.md --json

# Create with description from file (avoids shell escaping issues)
fbd create "Issue title" --body-file=description.md --json
fbd create "Issue title" --body-file description.md -p 1 --json

# Read description from stdin
echo "Description text" | fbd create "Issue title" --body-file=- --json
cat description.md | fbd create "Issue title" --body-file - -p 1 --json

# Create epic with hierarchical child tasks
fbd create "Auth System" -t epic -p 1 --json                     # Returns: bd-a3f8e9
fbd create "Login UI" -p 1 --parent bd-a3f8e9 --json             # Auto-assigned: bd-a3f8e9.1
fbd create "Backend validation" -p 1 --parent bd-a3f8e9 --json   # Auto-assigned: bd-a3f8e9.2
fbd create "Tests" -p 1 --parent bd-a3f8e9 --json                # Auto-assigned: bd-a3f8e9.3

# Create and link discovered work (one command)
fbd create "Found bug" -t bug -p 1 --deps discovered-from:<parent-id> --json

# Create with external reference (v0.9.2+)
fbd create "Fix login" -t bug -p 1 --external-ref "gh-123" --json  # Short form
fbd create "Fix login" -t bug -p 1 --external-ref "https://github.com/org/repo/issues/123" --json  # Full URL
fbd create "Jira task" -t task -p 1 --external-ref "jira-PROJ-456" --json  # Custom prefix
```

### Update Issues

```bash
# Update one or more issues
fbd update <id> [<id>...] --status in_progress --json
fbd update <id> [<id>...] --priority 1 --json
fbd update <id> [<id>...] --spec-id "docs/specs/auth.md" --json

# Update external reference (v0.9.2+)
fbd update <id> --external-ref "gh-456" --json           # Short form
fbd update <id> --external-ref "jira-PROJ-789" --json    # Custom prefix

# Atomically claim an issue for work (prevents race conditions)
# Sets assignee to you and status to in_progress in one atomic operation
# Fails if already claimed (assignee is not empty)
fbd update <id> --claim --json

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

### State (Labels as Cache)

For operational state tracking on role beads. Uses `<dimension>:<value>` label convention.
See [LABELS.md](LABELS.md#operational-state-pattern-labels-as-cache) for full pattern documentation.

```bash
# Query current state value
fbd state <id> <dimension>                    # Output: value
fbd state witness-abc patrol                  # Output: active
fbd state --json witness-abc patrol           # {"issue_id": "...", "dimension": "patrol", "value": "active"}

# List all state dimensions on an issue
fbd state list <id> --json
fbd state list witness-abc                    # patrol: active, mode: normal, health: healthy

# Set state (creates event + updates label atomically)
fbd set-state <id> <dimension>=<value> --reason "explanation" --json
fbd set-state witness-abc patrol=muted --reason "Investigating stuck polecat"
fbd set-state witness-abc mode=degraded --reason "High error rate"
```

**Common dimensions:**
- `patrol`: active, muted, suspended
- `mode`: normal, degraded, maintenance
- `health`: healthy, warning, failing
- `status`: idle, working, blocked

**What `set-state` does:**
1. Creates event bead with reason (source of truth)
2. Removes old `<dimension>:*` label if exists
3. Adds new `<dimension>:<value>` label (cache)

## Filtering & Search

### Basic Filters

```bash
# Filter by status, priority, type
fbd list --status open --priority 1 --json               # Status and priority
fbd list --assignee alice --json                         # By assignee
fbd list --type bug --json                               # By issue type
fbd list --id bd-123,bd-456 --json                       # Specific IDs
fbd list --spec "docs/specs/" --json                     # Spec prefix
```

### Label Filters

```bash
# Labels (AND: must have ALL)
fbd list --label bug,critical --json

# Labels (OR: has ANY)
fbd list --label-any frontend,backend --json
```

### Text Search

```bash
# Title search (substring)
fbd list --title "auth" --json

# Pattern matching (case-insensitive substring)
fbd list --title-contains "auth" --json                  # Search in title
fbd list --desc-contains "implement" --json              # Search in description
fbd list --notes-contains "TODO" --json                  # Search in notes

# Find beads issue by external reference
fbd list --json | jq -r '.[] | select(.external_ref == "gh-123") | .id'
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

### Orphan Detection

Find issues referenced in git commits that were never closed:

```bash
# Basic usage - scan current repo
fbd orphans

# Cross-repo: scan CODE repo's commits against external BEADS database
cd ~/my-code-repo
fbd orphans --db ~/my-beads-repo/.beads/beads.db

# JSON output
fbd orphans --json
```

**Use case**: When your beads database lives in a separate repository from your code, run `fbd orphans` from the code repo and point `--db` to the external database. This scans commits in your current directory while checking issue status from the specified database.

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

### Reset

Remove all local beads data and return to uninitialized state.

```bash
# Preview what would be removed (dry-run)
fbd admin reset

# Actually perform the reset
fbd admin reset --force
```

**What gets removed:**
- `.beads/` directory (database, JSONL, config)
- Git hooks installed by fbd
- Merge driver configuration
- Sync branch worktrees (`.git/beads-worktrees/`)

**What does NOT get removed:**
- Remote sync branch (if configured)
- JSONL history in git commits
- Remote repository data

**Important:** If you want a complete clean slate (including remote data), see [Troubleshooting: Old data returns after reset](TROUBLESHOOTING.md#old-data-returns-after-reset).

**Note:** The `--hard` and `--skip-init` flags mentioned in some discussions were never implemented. Use `--force` to perform the reset.

## Molecular Chemistry

Beads uses a chemistry metaphor for template-based workflows. See [MOLECULES.md](MOLECULES.md) for full documentation.

### Phase Transitions

| Phase | State | Storage | Command |
|-------|-------|---------|---------|
| Solid | Proto | `.beads/` | `fbd formula list` |
| Liquid | Mol | `.beads/` | `fbd mol pour` |
| Vapor | Wisp | `.beads/` (Ephemeral=true, not exported) | `fbd mol wisp` |

### Proto/Template Commands

```bash
# List available formulas (templates)
fbd formula list --json

# Show proto structure and variables
fbd mol show <proto-id> --json

# Extract proto from ad-hoc epic
fbd mol distill <epic-id> --json
```

### Pour (Proto to Mol)

```bash
# Instantiate proto as persistent mol (solid → liquid)
fbd mol pour <proto-id> --var key=value --json

# Preview what would be created
fbd mol pour <proto-id> --var key=value --dry-run

# Assign root issue
fbd mol pour <proto-id> --var key=value --assignee alice --json

# Attach additional protos during pour
fbd mol pour <proto-id> --attach <other-proto> --json
```

### Wisp Commands

```bash
# Instantiate proto as ephemeral wisp (solid → vapor)
fbd mol wisp <proto-id> --var key=value --json

# List all wisps
fbd mol wisp list --json
fbd mol wisp list --all --json    # Include closed

# Garbage collect orphaned wisps
fbd mol wisp gc --json
fbd mol wisp gc --age 24h --json  # Custom age threshold
fbd mol wisp gc --dry-run         # Preview what would be cleaned
```

### Bonding (Combining Work)

```bash
# Polymorphic combine - handles proto+proto, proto+mol, mol+mol
fbd mol bond <A> <B> --json

# Bond types
fbd mol bond <A> <B> --type sequential --json   # B runs after A (default)
fbd mol bond <A> <B> --type parallel --json     # B runs alongside A
fbd mol bond <A> <B> --type conditional --json  # B runs only if A fails

# Phase control
fbd mol bond <proto> <mol> --pour --json   # Force persistent spawn
fbd mol bond <proto> <mol> --wisp --json   # Force ephemeral spawn

# Dynamic bonding (custom child IDs)
fbd mol bond <proto> <mol> --ref arm-{{name}} --var name=ace --json

# Preview bonding
fbd mol bond <A> <B> --dry-run
```

### Squash (Wisp to Digest)

```bash
# Compress wisp to permanent digest
fbd mol squash <ephemeral-id> --json

# With agent-provided summary
fbd mol squash <ephemeral-id> --summary "Work completed" --json

# Preview
fbd mol squash <ephemeral-id> --dry-run

# Keep wisp children after squash
fbd mol squash <ephemeral-id> --keep-children --json
```

### Burn (Discard Wisp)

```bash
# Delete wisp without digest (destructive)
fbd mol burn <ephemeral-id> --json

# Preview
fbd mol burn <ephemeral-id> --dry-run

# Skip confirmation
fbd mol burn <ephemeral-id> --force --json
```

**Note:** Most mol commands require `--no-daemon` flag when daemon is running.

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

### Migrate to Sync Branch

Set up a dedicated sync branch for beads data, keeping your working branches clean.

```bash
# Basic setup (creates orphan branch by default)
fbd migrate sync beads-sync                             # Create orphan sync branch
fbd migrate sync beads-sync --dry-run                   # Preview without changes

# Force reconfigure if already set up
fbd migrate sync beads-sync --force                     # Reconfigure sync branch

# Migrate existing non-orphan branch to orphan
fbd migrate sync beads-sync --orphan                    # Delete and recreate as orphan
```

**Behavior:**

| Scenario | Result |
|----------|--------|
| Branch doesn't exist | Creates orphan branch (no shared history) |
| Branch exists locally | Uses existing branch as-is |
| Branch exists + `--orphan` | Migrates: deletes and recreates as orphan |
| Remote only | Fetches from remote |
| Remote only + `--orphan` | Creates local orphan (ignores remote) |

**Why orphan branches?**

- Clean "data sync channel" mental model
- No accidental merge risk (git warns loudly)
- Smaller repository footprint (no stale source code)
- Sync branch contains only `.beads/` directory

**After setup:**

- `fbd sync` commits beads changes to the sync branch via worktree
- Your working branch stays clean of beads commits
- Essential for multi-clone setups where clones work independently

**Safety features for `--orphan` migration:**

- **Unpushed commit check**: If the branch has unpushed commits, migration fails with a helpful error. Use `--force` to override.
- **Existing worktree**: If a worktree exists for the branch, it's automatically removed before migration.
- **Non-destructive to remote**: The remote branch is not modified; use `git push --force` to update it after migration.

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
```

### Key-Value Store

Store user-defined key-value pairs that persist across sessions. Useful for feature flags, environment config, or agent memory.

```bash
# Set a value
fbd kv set <key> <value>
fbd kv set feature_flag true
fbd kv set api_endpoint https://api.example.com

# Get a value
fbd kv get <key>
fbd kv get feature_flag                 # Prints: true
fbd kv get missing_key                  # Prints: missing_key (not set), exits 1

# Delete a key
fbd kv clear <key>
fbd kv clear feature_flag

# List all key-value pairs
fbd kv list
fbd kv list --json                      # Machine-readable output
```

**Storage notes:**
- KV data is stored in the local database with a `kv.` prefix
- In `dolt-native` or `belt-and-suspenders` sync modes, KV data syncs via Dolt remotes
- In `git-portable` mode, KV data stays local (not exported to JSONL)

**Use cases:**
- Feature flags: `fbd set debug_mode true`
- Environment config: `fbd set staging_url https://staging.example.com`
- Agent memory: `fbd set last_migration 20240115_add_users.sql`
- Session state: `fbd set current_sprint 42`

## Issue Types

- `bug` - Something broken that needs fixing
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature composed of multiple issues (supports hierarchical children)
- `chore` - Maintenance work (dependencies, tooling)

**Hierarchical children:** Epics can have child issues with dotted IDs (e.g., `bd-a3f8e9.1`, `bd-a3f8e9.2`). Children are auto-numbered sequentially. Up to 3 levels of nesting supported.

## Issue Statuses

- `open` - Ready to be worked on
- `in_progress` - Currently being worked on
- `blocked` - Cannot proceed (waiting on dependencies)
- `deferred` - Deliberately put on ice for later
- `closed` - Work completed
- `tombstone` - Deleted issue (suppresses resurrections)
- `pinned` - Stays open indefinitely (used for hooks, anchors)

**Note:** The `pinned` status is used by orchestrators for hook management and persistent work items that should never be auto-closed or cleaned up.

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
# ○ bd-42 [P1] [bug] - Fix authentication bug
# ○ bd-43 [P2] [feature] - Add user settings page
```

**Dependency visibility:** When issues have blocking dependencies, they appear inline:

```bash
fbd list --parent epic-123
# ○ bd-123.1 [P1] [task] - Design API (blocks: bd-123.2, bd-123.3)
# ○ bd-123.2 [P1] [task] - Implement endpoints (blocked by: bd-123.1, blocks: bd-123.3)
# ○ bd-123.3 [P1] [task] - Add tests (blocked by: bd-123.1, bd-123.2)
```

This makes blocking relationships visible without running `fbd show` on each issue.

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

## Editor Integration

### Setup Commands

```bash
# Setup editor integration (choose based on your editor)
fbd setup factory  # Factory.ai Droid - creates/updates AGENTS.md (universal standard)
fbd setup codex    # Codex CLI - creates/updates AGENTS.md
fbd setup claude   # Claude Code - installs SessionStart/PreCompact hooks
fbd setup cursor   # Cursor IDE - creates .cursor/rules/beads.mdc
fbd setup aider    # Aider - creates .aider.conf.yml

# Check if integration is installed
fbd setup factory --check
fbd setup codex --check
fbd setup claude --check
fbd setup cursor --check
fbd setup aider --check

# Remove integration
fbd setup factory --remove
fbd setup codex --remove
fbd setup claude --remove
fbd setup cursor --remove
fbd setup aider --remove
```

**Claude Code options:**
```bash
fbd setup claude              # Install globally (~/.claude/settings.json)
fbd setup claude --project    # Install for this project only
fbd setup claude --stealth    # Use stealth mode (flush only, no git operations)
```

**What each setup does:**
- **Factory.ai** (`fbd setup factory`): Creates or updates AGENTS.md with beads workflow instructions (works with multiple AI tools using the AGENTS.md standard)
- **Codex CLI** (`fbd setup codex`): Creates or updates AGENTS.md with beads workflow instructions for Codex
- **Claude Code** (`fbd setup claude`): Adds hooks to Claude Code's settings.json that run `fbd prime` on SessionStart and PreCompact events
- **Cursor** (`fbd setup cursor`): Creates `.cursor/rules/beads.mdc` with workflow instructions
- **Aider** (`fbd setup aider`): Creates `.aider.conf.yml` with fbd workflow instructions

See also:
- [INSTALLING.md](INSTALLING.md#ide-and-editor-integrations) - Installation guide
- [AIDER_INTEGRATION.md](AIDER_INTEGRATION.md) - Detailed Aider guide
- [CLAUDE_INTEGRATION.md](CLAUDE_INTEGRATION.md) - Claude integration design

## See Also

- [AGENTS.md](../AGENTS.md) - Main agent workflow guide
- [MOLECULES.md](MOLECULES.md) - Molecular chemistry metaphor (protos, pour, bond, squash, burn)
- [DAEMON.md](DAEMON.md) - Daemon management and event-driven mode
- [GIT_INTEGRATION.md](GIT_INTEGRATION.md) - Git workflows and merge strategies
- [LABELS.md](../LABELS.md) - Label system guide
- [README.md](../README.md) - User documentation
