---
sidebar_position: 2
title: Database Corruption
description: Recover from SQLite database corruption
---

# Database Corruption Recovery

This runbook helps you recover from SQLite database corruption in Beads.

## Symptoms

- SQLite error messages during `fbd` commands
- "database is locked" errors that persist
- Missing issues that should exist
- Inconsistent state between JSONL and database

## Diagnosis

```bash
# Check database integrity
fbd status

# Look for corruption indicators
ls -la .beads/beads.db*
```

If you see `-wal` or `-shm` files alongside `beads.db`, a transaction may have been interrupted.

## Solution

:::warning
Back up your `.beads/` directory before proceeding.
:::

**Step 1:** Stop the daemon
```bash
fbd daemon stop
```

**Step 2:** Back up current state
```bash
cp -r .beads .beads.backup
```

**Step 3:** Rebuild from JSONL (source of truth)
```bash
fbd doctor --fix
```

**Step 4:** Verify recovery
```bash
fbd status
fbd list
```

**Step 5:** Restart daemon
```bash
fbd daemon start
```

## Prevention

- Avoid interrupting `fbd sync` operations
- Let the daemon handle synchronization
- Use `fbd daemon stop` before system shutdown
