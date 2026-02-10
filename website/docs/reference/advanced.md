---
id: advanced
title: Advanced Features
sidebar_position: 3
---

# Advanced Features

Advanced beads functionality.

## Issue Rename

Rename issues while preserving references:

```bash
fbd rename bd-42 bd-new-id
fbd rename bd-42 bd-new-id --dry-run  # Preview
```

Updates:
- All dependencies pointing to old ID
- All references in other issues
- Comments and descriptions

## Issue Merge

Merge duplicate issues:

```bash
fbd merge bd-42 bd-43 --into bd-41
fbd merge bd-42 bd-43 --into bd-41 --dry-run
```

What gets merged:
- Dependencies → target
- Text references updated across all issues
- Source issues closed with merge reason

## Database Compaction

Reduce database size by compacting old issues:

```bash
# View compaction statistics
fbd admin compact --stats

# Preview candidates (30+ days closed)
fbd admin compact --analyze --json

# Apply agent-generated summary
fbd admin compact --apply --id bd-42 --summary summary.txt

# Immediate deletion (CAUTION!)
fbd admin cleanup --force
```

**When to compact:**
- Database > 10MB with old closed issues
- After major milestones
- Before archiving project phase

## Restore from History

View deleted or compacted issues from git:

```bash
fbd restore bd-42 --show
fbd restore bd-42 --to-file issue.json
```

## Database Inspection

```bash
# Schema info
fbd info --schema --json

# Raw database query (advanced)
sqlite3 .beads/beads.db "SELECT * FROM issues LIMIT 5"
```

## Custom Tables

Extend the database with custom tables:

```go
// In Go code using beads as library
storage.UnderlyingDB().Exec(`
  CREATE TABLE IF NOT EXISTS custom_table (...)
`)
```

See [EXTENDING.md](https://github.com/steveyegge/fastbeads/blob/main/docs/EXTENDING.md).

## Event System

Subscribe to beads events:

```bash
# View recent events
fbd events list --since 1h

# Watch events in real-time
fbd events watch
```

Events:
- `issue.created`
- `issue.updated`
- `issue.closed`
- `dependency.added`
- `sync.completed`

## Batch Operations

### Create Multiple

```bash
cat issues.jsonl | fbd import -i -
```

### Update Multiple

```bash
fbd list --status open --priority 4 --json | \
  jq -r '.[].id' | \
  xargs -I {} fbd update {} --priority 3
```

### Close Multiple

```bash
fbd list --label "sprint-1" --status open --json | \
  jq -r '.[].id' | \
  xargs -I {} fbd close {} --reason "Sprint complete"
```

## API Access

Use beads as a Go library:

```go
import "github.com/steveyegge/fastbeads/internal/storage"

db, _ := storage.NewSQLite(".beads/beads.db")
issues, _ := db.ListIssues(storage.ListOptions{
    Status: "open",
})
```

## Performance Tuning

### Large Databases

```bash
# Enable WAL mode
fbd config set database.wal_mode true

# Increase cache
fbd config set database.cache_size 10000
```

### Many Concurrent Agents

```bash
# Use event-driven daemon
export BEADS_DAEMON_MODE=events
fbd daemons killall
```

### CI/CD Optimization

```bash
# Disable daemon in CI
export BEADS_NO_DAEMON=true
fbd --no-daemon list
```
