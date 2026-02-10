---
id: sync
title: Sync & Export
sidebar_position: 6
---

# Sync & Export Commands

Commands for synchronizing with git.

## fbd sync

Full sync cycle: export, commit, push.

```bash
fbd sync [flags]
```

**What it does:**
1. Exports database to `.beads/issues.jsonl`
2. Stages the JSONL file
3. Commits with auto-generated message
4. Pushes to remote

**Flags:**
```bash
--json     JSON output
--dry-run  Preview without changes
```

**Examples:**
```bash
fbd sync
fbd sync --json
```

**When to use:**
- End of work session
- Before switching branches
- After significant changes

## fbd export

Export database to JSONL.

```bash
fbd export [flags]
```

**Flags:**
```bash
--output, -o    Output file (default: .beads/issues.jsonl)
--dry-run       Preview without writing
--json          JSON output
```

**Examples:**
```bash
fbd export
fbd export -o backup.jsonl
fbd export --dry-run
```

## fbd import

Import from JSONL file.

```bash
fbd import -i <file> [flags]
```

**Flags:**
```bash
--input, -i           Input file (required)
--dry-run             Preview without changes
--orphan-handling     How to handle missing parents
--dedupe-after        Run duplicate detection after import
--json                JSON output
```

**Orphan handling modes:**
| Mode | Behavior |
|------|----------|
| `allow` | Import orphans without validation (default) |
| `resurrect` | Restore deleted parents as tombstones |
| `skip` | Skip orphaned children with warning |
| `strict` | Fail if parent missing |

**Examples:**
```bash
fbd import -i .beads/issues.jsonl
fbd import -i backup.jsonl --dry-run
fbd import -i issues.jsonl --orphan-handling resurrect
fbd import -i issues.jsonl --dedupe-after --json
```

## fbd migrate

Migrate database schema.

```bash
fbd migrate [flags]
```

**Flags:**
```bash
--inspect    Show migration plan (for agents)
--dry-run    Preview without changes
--cleanup    Remove old files after migration
--yes        Skip confirmation
--json       JSON output
```

**Examples:**
```bash
fbd migrate --inspect --json
fbd migrate --dry-run
fbd migrate
fbd migrate --cleanup --yes
```

## fbd hooks

Manage git hooks.

```bash
fbd hooks <subcommand> [flags]
```

**Subcommands:**
| Command | Description |
|---------|-------------|
| `install` | Install git hooks |
| `uninstall` | Remove git hooks |
| `status` | Check hook status |

**Examples:**
```bash
fbd hooks install
fbd hooks status
fbd hooks uninstall
```

## Auto-Sync Behavior

### With Daemon (Default)

The daemon handles sync automatically:
- Exports to JSONL after changes (5s debounce)
- Imports from JSONL when newer

### Without Daemon

Use `--no-daemon` flag:
- Changes only written to SQLite
- Must manually export/sync

```bash
fbd --no-daemon create "Task"
fbd export  # Manual export needed
```

## Conflict Resolution

### Merge Driver (Recommended)

Install the beads merge driver:

```bash
fbd init  # Prompts for merge driver setup
```

The driver automatically:
- Merges non-conflicting changes
- Preserves both sides for real conflicts
- Uses latest timestamp for same-issue edits

### Manual Resolution

```bash
# After merge conflict
git checkout --ours .beads/issues.jsonl
fbd import -i .beads/issues.jsonl
fbd sync
```

## Deletion Tracking

Deletions sync via `.beads/deletions.jsonl`:

```bash
# Delete issue
fbd delete bd-42

# View deletions
fbd deleted
fbd deleted --since=30d

# Deletions propagate via git
git pull  # Imports deletions from remote
```

## Best Practices

1. **Always sync at session end** - `fbd sync`
2. **Install git hooks** - `fbd hooks install`
3. **Use merge driver** - Avoids manual conflict resolution
4. **Check sync status** - `fbd info` shows daemon/sync state
