---
id: essential
title: Essential Commands
sidebar_position: 2
---

# Essential Commands

The most important commands for daily use.

## fbd create

Create a new issue.

```bash
fbd create <title> [flags]
```

**Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Issue type: bug, feature, task, epic, chore |
| `--priority` | `-p` | Priority: 0-4 (0=critical, 4=backlog) |
| `--description` | `-d` | Detailed description |
| `--labels` | `-l` | Comma-separated labels |
| `--parent` | | Parent issue ID (for hierarchical) |
| `--deps` | | Dependencies (e.g., `discovered-from:bd-42`) |
| `--json` | | JSON output |

**Examples:**
```bash
fbd create "Fix login bug" -t bug -p 1
fbd create "Add dark mode" -t feature -p 2 --description="User requested"
fbd create "Subtask" --parent bd-42 -p 2
fbd create "Found during work" --deps discovered-from:bd-42 --json
```

## fbd list

List issues with filters.

```bash
fbd list [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--status` | Filter by status: open, in_progress, closed |
| `--priority` | Filter by priority (comma-separated) |
| `--type` | Filter by type (comma-separated) |
| `--label-any` | Issues with any of these labels |
| `--label-all` | Issues with all of these labels |
| `--json` | JSON output |

**Examples:**
```bash
fbd list --status open
fbd list --priority 0,1 --type bug
fbd list --label-any urgent,critical --json
```

## fbd show

Show issue details.

```bash
fbd show <id> [flags]
```

**Examples:**
```bash
fbd show bd-42
fbd show bd-42 --json
fbd show bd-42 bd-43 bd-44  # Multiple issues
```

## fbd update

Update issue fields.

```bash
fbd update <id> [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--status` | New status |
| `--priority` | New priority |
| `--title` | New title |
| `--description` | New description |
| `--add-label` | Add label |
| `--remove-label` | Remove label |
| `--json` | JSON output |

**Examples:**
```bash
fbd update bd-42 --status in_progress
fbd update bd-42 --priority 0 --add-label urgent
fbd update bd-42 --title "Updated title" --json
```

## fbd close

Close an issue.

```bash
fbd close <id> [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--reason` | Closure reason |
| `--json` | JSON output |

**Examples:**
```bash
fbd close bd-42
fbd close bd-42 --reason "Fixed in PR #123"
fbd close bd-42 --json
```

## fbd ready

Show issues ready to work on (no blockers).

```bash
fbd ready [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--priority` | Filter by priority |
| `--type` | Filter by type |
| `--json` | JSON output |

**Examples:**
```bash
fbd ready
fbd ready --priority 1
fbd ready --json
```

## fbd blocked

Show blocked issues and their blockers.

```bash
fbd blocked [flags]
```

**Examples:**
```bash
fbd blocked
fbd blocked --json
```

## fbd sync

Force immediate sync to git.

```bash
fbd sync [flags]
```

Performs:
1. Export database to JSONL
2. Git add `.beads/issues.jsonl`
3. Git commit
4. Git push

**Examples:**
```bash
fbd sync
fbd sync --json
```

## fbd info

Show system information.

```bash
fbd info [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--whats-new` | Show recent version changes |
| `--schema` | Show database schema |
| `--json` | JSON output |

**Examples:**
```bash
fbd info
fbd info --whats-new
fbd info --json
```

## fbd stats

Show project statistics.

```bash
fbd stats [flags]
```

**Examples:**
```bash
fbd stats
fbd stats --json
```
