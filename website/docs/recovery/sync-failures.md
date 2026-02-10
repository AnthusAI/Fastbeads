---
sidebar_position: 5
title: Sync Failures
description: Recover from fbd sync failures
---

# Sync Failures Recovery

This runbook helps you recover from `fbd sync` failures.

## Symptoms

- `fbd sync` hangs or times out
- Network-related error messages
- "failed to push" or "failed to pull" errors
- Daemon not responding

## Diagnosis

```bash
# Check daemon status
fbd daemon status

# Check sync state
fbd status

# View daemon logs
cat .beads/daemon.log | tail -50
```

## Solution

**Step 1:** Stop the daemon
```bash
fbd daemon stop
```

**Step 2:** Check for lock files
```bash
ls -la .beads/*.lock
# Remove stale locks if daemon is definitely stopped
rm -f .beads/*.lock
```

**Step 3:** Force a fresh sync
```bash
fbd doctor --fix
```

**Step 4:** Restart daemon
```bash
fbd daemon start
```

**Step 5:** Verify sync works
```bash
fbd sync
fbd status
```

## Common Causes

| Cause | Solution |
|-------|----------|
| Network timeout | Retry with better connection |
| Stale lock file | Remove lock after stopping daemon |
| Corrupted state | Use `fbd doctor --fix` |
| Git conflicts | See [Merge Conflicts](/recovery/merge-conflicts) |

## Prevention

- Ensure stable network before sync
- Let sync complete before closing terminal
- Use `fbd daemon stop` before system shutdown
