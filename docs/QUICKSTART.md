# Beads Quickstart

Get up and running with Beads in 2 minutes.

## Installation

```bash
cd ~/src/beads
go build -o fbd ./cmd/fbd
./fbd --help
```

## Initialize

First time in a repository:

```bash
# Basic setup (prompts for contributor mode)
fbd init

# Dolt backend (version-controlled SQL database)
fbd init --backend dolt

# OSS contributor (fork workflow with separate planning repo)
fbd init --contributor

# Team member (branch workflow for collaboration)
fbd init --team

# Protected main branch (GitHub/GitLab)
fbd init --branch beads-sync
```

The wizard will:
- Create `.beads/` directory and database
- **Prompt for your role** (maintainer or contributor) unless a flag is provided
- Import existing issues from git (if any)
- Prompt to install git hooks (recommended)
- Prompt to configure git merge driver (recommended)
- Auto-start daemon for sync (SQLite backend only)

Notes:
- SQLite backend stores data in `.beads/beads.db`.
- Dolt backend stores data in `.beads/dolt/` and records `"database": "dolt"` in `.beads/metadata.json`.
- Dolt backend runs **single-process-only**; daemon mode is disabled.
- Dolt backend **auto-commits** after each successful write command in embedded mode (`dolt.auto-commit: on`). In server mode, auto-commit defaults to OFF. Override with `fbd --dolt-auto-commit off|on ...` or config.

### Role Configuration

During `fbd init`, you'll be asked: "Contributing to someone else's repo? [y/N]"

- Answer **Y** if you're contributing to a fork (runs contributor wizard)
- Answer **N** if you're the maintainer or have push access

This sets `git config beads.role` which determines how beads routes issues:

| Role | Use Case | Issue Storage |
|------|----------|---------------|
| `maintainer` | Repo owner, team with push access | In-repo `.beads/` |
| `contributor` | Fork contributor, OSS contributor | Separate planning repo |

You can also configure manually:

```bash
# Set as contributor
git config beads.role contributor

# Set as maintainer
git config beads.role maintainer

# Check current role
git config --get beads.role
```

**Note:** If `beads.role` is not configured, beads falls back to URL-based detection (deprecated). Run `fbd doctor` to check configuration status.

## Your First Issues

```bash
# Create a few issues
./fbd create "Set up database" -p 1 -t task
./fbd create "Create API" -p 2 -t feature
./fbd create "Add authentication" -p 2 -t feature

# List them
./fbd list
```

**Note:** Issue IDs are hash-based (e.g., `bd-a1b2`, `bd-f14c`) to prevent collisions when multiple agents/branches work concurrently.

**Dependency visibility:** When issues have blocking dependencies, `fbd list` shows them inline:
```
○ bd-a1b2 [P1] [task] - Set up database
○ bd-f14c [P2] [feature] - Create API (blocked by: bd-a1b2)
○ bd-g25d [P2] [feature] - Add authentication (blocked by: bd-f14c)
```

This makes dependencies unmissable when reviewing epic subtasks.

## Hierarchical Issues (Epics)

For large features, use hierarchical IDs to organize work:

```bash
# Create epic (generates parent hash ID)
./fbd create "Auth System" -t epic -p 1
# Returns: bd-a3f8e9

# Create child tasks (automatically get .1, .2, .3 suffixes)
./fbd create "Design login UI" -p 1       # bd-a3f8e9.1
./fbd create "Backend validation" -p 1    # bd-a3f8e9.2
./fbd create "Integration tests" -p 1     # bd-a3f8e9.3

# View hierarchy
./fbd dep tree bd-a3f8e9
```

Output:
```
🌲 Dependency tree for bd-a3f8e9:

→ bd-a3f8e9: Auth System [epic] [P1] (open)
  → bd-a3f8e9.1: Design login UI [P1] (open)
  → bd-a3f8e9.2: Backend validation [P1] (open)
  → bd-a3f8e9.3: Integration tests [P1] (open)
```

## Add Dependencies

```bash
# API depends on database
./fbd dep add bd-2 bd-1

# Auth depends on API
./fbd dep add bd-3 bd-2

# View the tree
./fbd dep tree bd-3
```

Output:
```
🌲 Dependency tree for bd-3:

→ bd-3: Add authentication [P2] (open)
  → bd-2: Create API [P2] (open)
    → bd-1: Set up database [P1] (open)
```

## Find Ready Work

```bash
./fbd ready
```

Output:
```
📋 Ready work (1 issues with no blockers):

1. [P1] bd-1: Set up database
```

Only bd-1 is ready because bd-2 and bd-3 are blocked!

## Work the Queue

```bash
# Start working on bd-1
./fbd update bd-1 --status in_progress

# Complete it
./fbd close bd-1 --reason "Database setup complete"

# Check ready work again
./fbd ready
```

Now bd-2 is ready! 🎉

## Track Progress

```bash
# See blocked issues
./fbd blocked

# View statistics
./fbd stats
```

## Database Location

By default: `~/.beads/default.db`

You can use project-specific databases:

```bash
./fbd --db ./my-project.db create "Task"
```

## Migrating Databases

After upgrading fbd, use `fbd migrate` to check for and migrate old database files:

```bash
# Inspect migration plan (AI agents)
./fbd migrate --inspect --json

# Check schema and config
./fbd info --schema --json

# Preview migration changes
./fbd migrate --dry-run

# Migrate old databases to beads.db
./fbd migrate

# Migrate and clean up old files
./fbd migrate --cleanup --yes
```

**AI agents:** Use `--inspect` to analyze migration safety before running. The system verifies required config keys and data integrity invariants.

## Database Maintenance

As your project accumulates closed issues, the database grows. Manage size with these commands:

```bash
# View compaction statistics
fbd admin compact --stats

# Preview compaction candidates (30+ days closed)
fbd admin compact --analyze --json --no-daemon

# Apply agent-generated summary
fbd admin compact --apply --id bd-42 --summary summary.txt --no-daemon

# Immediately delete closed issues (CAUTION: permanent!)
fbd admin cleanup --force
```

**When to compact:**
- Database file > 10MB with many old closed issues
- After major project milestones when old issues are no longer relevant
- Before archiving a project phase

**Note:** Compaction is permanent graceful decay. Original content is discarded but viewable via `fbd restore <id>` from git history.

## Background Daemon

fbd runs a background daemon for auto-sync and performance. You rarely need to manage it directly:

```bash
# Check daemon status
fbd info | grep daemon

# List all running daemons
fbd daemons list

# Force direct mode (skip daemon)
fbd --no-daemon ready
```

**When to disable daemon:**
- Git worktrees (required: `fbd --no-daemon`)
- CI/CD pipelines
- Resource-constrained environments

See [DAEMON.md](DAEMON.md) for complete daemon management guide.

## Next Steps

- Add labels: `./fbd create "Task" -l "backend,urgent"`
- Filter ready work: `./fbd ready --priority 1`
- Search issues: `./fbd list --status open`
- Detect cycles: `./fbd dep cycles`

See [README.md](../README.md) for full documentation.
