---
id: index
title: CLI Reference
sidebar_position: 1
---

# CLI Reference

Complete reference for all `fbd` commands.

## Command Structure

```bash
fbd [global-flags] <command> [command-flags] [arguments]
```

### Global Flags

| Flag | Description |
|------|-------------|
| `--db <path>` | Use specific database file |
| `--no-daemon` | Bypass daemon, direct database access |
| `--json` | Output in JSON format |
| `--quiet` | Suppress non-essential output |
| `--verbose` | Verbose output |
| `--version` | Show version |
| `--help` | Show help |

## Command Categories

### Essential Commands

Most frequently used:

| Command | Description |
|---------|-------------|
| `fbd create` | Create new issue |
| `fbd list` | List issues with filters |
| `fbd show` | Show issue details |
| `fbd update` | Update issue fields |
| `fbd close` | Close an issue |
| `fbd ready` | Show unblocked work |
| `fbd sync` | Force sync to git |

### Issue Management

| Command | Description |
|---------|-------------|
| `fbd create` | Create issue |
| `fbd show` | Show details |
| `fbd update` | Update fields |
| `fbd close` | Close issue |
| `fbd delete` | Delete issue |
| `fbd reopen` | Reopen closed issue |

### Dependencies

| Command | Description |
|---------|-------------|
| `fbd dep add` | Add dependency |
| `fbd dep remove` | Remove dependency |
| `fbd dep tree` | Show dependency tree |
| `fbd dep cycles` | Detect circular dependencies |
| `fbd blocked` | Show blocked issues |
| `fbd ready` | Show unblocked issues |

### Labels & Comments

| Command | Description |
|---------|-------------|
| `fbd label add` | Add label to issue |
| `fbd label remove` | Remove label |
| `fbd label list` | List all labels |
| `fbd comment add` | Add comment |
| `fbd comment list` | List comments |

### Sync & Export

| Command | Description |
|---------|-------------|
| `fbd sync` | Full sync cycle |
| `fbd export` | Export to JSONL |
| `fbd import` | Import from JSONL |
| `fbd migrate` | Migrate database schema |

### System

| Command | Description |
|---------|-------------|
| `fbd init` | Initialize beads in project |
| `fbd info` | Show system info |
| `fbd version` | Show version |
| `fbd config` | Manage configuration |
| `fbd daemons` | Manage daemons |
| `fbd hooks` | Manage git hooks |

### Workflows

| Command | Description |
|---------|-------------|
| `fbd pour` | Instantiate formula as molecule |
| `fbd wisp` | Create ephemeral wisp |
| `fbd mol` | Manage molecules |
| `fbd pin` | Pin work to agent |
| `fbd hook` | Show pinned work |

## Quick Reference

### Creating Issues

```bash
# Basic
fbd create "Title" -t task -p 2

# With description
fbd create "Title" --description="Details here" -t bug -p 1

# With labels
fbd create "Title" -l "backend,urgent"

# As child of epic
fbd create "Subtask" --parent bd-42

# With discovered-from link
fbd create "Found bug" --deps discovered-from:bd-42

# JSON output
fbd create "Title" --json
```

### Querying Issues

```bash
# All open issues
fbd list --status open

# High priority bugs
fbd list --status open --priority 0,1 --type bug

# With specific labels
fbd list --label-any urgent,critical

# JSON output
fbd list --json
```

### Working with Dependencies

```bash
# Add: bd-2 depends on bd-1
fbd dep add bd-2 bd-1

# View tree
fbd dep tree bd-2

# Find cycles
fbd dep cycles

# What's ready to work?
fbd ready

# What's blocked?
fbd blocked
```

### Syncing

```bash
# Full sync (export + commit + push)
fbd sync

# Force export
fbd export

# Import from file
fbd import -i .beads/issues.jsonl
```

## See Also

- [Essential Commands](/cli-reference/essential)
- [Issue Commands](/cli-reference/issues)
- [Dependency Commands](/cli-reference/dependencies)
- [Label Commands](/cli-reference/labels)
- [Sync Commands](/cli-reference/sync)
