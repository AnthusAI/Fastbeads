---
id: git-integration
title: Git Integration
sidebar_position: 2
---

# Git Integration

How beads integrates with git.

## Overview

Beads uses git for:
- **JSONL sync** - Issues stored in `.beads/issues.jsonl`
- **Deletion tracking** - `.beads/deletions.jsonl`
- **Conflict resolution** - Custom merge driver
- **Hooks** - Auto-sync on git operations

## File Structure

```
.beads/
├── beads.db           # SQLite database (gitignored)
├── issues.jsonl       # Issue data (git-tracked)
├── deletions.jsonl    # Deletion manifest (git-tracked)
├── config.toml        # Project config (git-tracked)
└── fbd.sock            # Daemon socket (gitignored)
```

## Git Hooks

### Installation

```bash
fbd hooks install
```

Installs:
- **pre-commit** - Exports database to JSONL
- **post-merge** - Imports from JSONL after pull
- **pre-push** - Ensures sync before push

### Status

```bash
fbd hooks status
```

### Uninstall

```bash
fbd hooks uninstall
```

## Merge Driver

### Purpose

The beads merge driver handles JSONL conflicts automatically:
- Merges non-conflicting changes
- Uses latest timestamp for same-issue edits
- Preserves both sides for real conflicts

### Installation

```bash
fbd init  # Prompts for merge driver setup
```

Or manually add to `.gitattributes`:

```gitattributes
.beads/issues.jsonl merge=beads
.beads/deletions.jsonl merge=beads
```

And `.git/config`:

```ini
[merge "beads"]
    name = Beads JSONL merge driver
    driver = fbd merge-driver %O %A %B
```

## Protected Branches

For protected main branches:

```bash
fbd init --branch beads-sync
```

This:
- Creates a separate `beads-sync` branch
- Syncs issues to that branch
- Avoids direct commits to main

## Git Worktrees

Beads requires `--no-daemon` in git worktrees:

```bash
# In worktree
fbd --no-daemon create "Task"
fbd --no-daemon list
```

Why: Daemon uses `.beads/fbd.sock` which conflicts across worktrees.

## Branch Workflows

### Feature Branch

```bash
git checkout -b feature-x
fbd create "Feature X" -t feature
# Work...
fbd sync
git push
```

### Fork Workflow

```bash
# In fork
fbd init --contributor
# Work in separate planning repo...
fbd sync
```

### Team Workflow

```bash
fbd init --team
# All team members share issues.jsonl
git pull  # Auto-imports via hook
```

## Conflict Resolution

### With Merge Driver

Automatic - driver handles most conflicts.

### Manual Resolution

```bash
# After conflict
git checkout --ours .beads/issues.jsonl
fbd import -i .beads/issues.jsonl
fbd sync
git add .beads/
git commit
```

### Duplicate Detection

After merge:

```bash
fbd duplicates --auto-merge
```

## Best Practices

1. **Install hooks** - `fbd hooks install`
2. **Use merge driver** - Avoid manual conflict resolution
3. **Sync regularly** - `fbd sync` at session end
4. **Pull before work** - Get latest issues
5. **Use `--no-daemon` in worktrees**
