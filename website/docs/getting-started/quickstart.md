---
id: quickstart
title: Quick Start
sidebar_position: 2
---

# Beads Quick Start

Get up and running with Beads in 2 minutes.

## Initialize

First time in a repository:

```bash
# Basic setup
fbd init

# Dolt backend (version-controlled SQL database)
fbd init --backend dolt

# For AI agents (non-interactive)
fbd init --quiet

# OSS contributor (fork workflow)
fbd init --contributor

# Team member (branch workflow)
fbd init --team

# Protected main branch (GitHub/GitLab)
fbd init --branch beads-sync
```

The wizard will:
- Create `.beads/` directory and database
- Import existing issues from git (if any)
- Prompt to install git hooks (recommended)
- Prompt to configure git merge driver (recommended)
- Auto-start daemon for sync (SQLite backend only)

Notes:
- SQLite backend stores data in `.beads/beads.db`.
- Dolt backend stores data in `.beads/dolt/` and records `"database": "dolt"` in `.beads/metadata.json`.
- Dolt backend runs **single-process-only**; daemon mode is disabled.

Notes:
- SQLite backend stores data in `.beads/beads.db`.
- Dolt backend stores data in `.beads/dolt/` and records `"database": "dolt"` in `.beads/metadata.json`.

## Your First Issues

```bash
# Create a few issues
fbd create "Set up database" -p 1 -t task
fbd create "Create API" -p 2 -t feature
fbd create "Add authentication" -p 2 -t feature

# List them
fbd list
```

**Note:** Issue IDs are hash-based (e.g., `bd-a1b2`, `bd-f14c`) to prevent collisions when multiple agents/branches work concurrently.

## Hierarchical Issues (Epics)

For large features, use hierarchical IDs to organize work:

```bash
# Create epic (generates parent hash ID)
fbd create "Auth System" -t epic -p 1
# Returns: bd-a3f8e9

# Create child tasks (automatically get .1, .2, .3 suffixes)
fbd create "Design login UI" -p 1       # bd-a3f8e9.1
fbd create "Backend validation" -p 1    # bd-a3f8e9.2
fbd create "Integration tests" -p 1     # bd-a3f8e9.3

# View hierarchy
fbd dep tree bd-a3f8e9
```

Output:
```
Dependency tree for bd-a3f8e9:

> bd-a3f8e9: Auth System [epic] [P1] (open)
  > bd-a3f8e9.1: Design login UI [P1] (open)
  > bd-a3f8e9.2: Backend validation [P1] (open)
  > bd-a3f8e9.3: Integration tests [P1] (open)
```

## Add Dependencies

```bash
# API depends on database
fbd dep add bd-2 bd-1

# Auth depends on API
fbd dep add bd-3 bd-2

# View the tree
fbd dep tree bd-3
```

Output:
```
Dependency tree for bd-3:

> bd-3: Add authentication [P2] (open)
  > bd-2: Create API [P2] (open)
    > bd-1: Set up database [P1] (open)
```

## Find Ready Work

```bash
fbd ready
```

Output:
```
Ready work (1 issues with no blockers):

1. [P1] bd-1: Set up database
```

Only bd-1 is ready because bd-2 and bd-3 are blocked!

## Work the Queue

```bash
# Start working on bd-1
fbd update bd-1 --status in_progress

# Complete it
fbd close bd-1 --reason "Database setup complete"

# Check ready work again
fbd ready
```

Now bd-2 is ready!

## Track Progress

```bash
# See blocked issues
fbd blocked

# View statistics
fbd stats
```

## Database Location

By default, the database is in `.beads/beads.db` (gitignored).

The JSONL file `.beads/issues.jsonl` is git-tracked and syncs automatically.

## Next Steps

- Add labels: `fbd create "Task" -l "backend,urgent"`
- Filter ready work: `fbd ready --priority 1`
- Search issues: `fbd list --status open`
- Detect cycles: `fbd dep cycles`
- See [CLI Reference](/cli-reference) for all commands
