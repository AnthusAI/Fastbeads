---
id: daemon
title: Daemon Architecture
sidebar_position: 3
---

# Daemon Architecture

Beads runs a background daemon for auto-sync and performance.

## Overview

Each workspace gets its own daemon process:
- Auto-starts on first `fbd` command
- Handles database ↔ JSONL synchronization
- Listens on `.beads/fbd.sock` (Unix) or `.beads/fbd.pipe` (Windows)
- Version checking prevents mismatches after upgrades

## How It Works

```
CLI Command
    ↓
RPC to Daemon
    ↓
Daemon executes
    ↓
Auto-sync to JSONL (5s debounce)
```

Without daemon, commands access the database directly (slower, no auto-sync).

## Managing Daemons

```bash
# List all running daemons
fbd daemons list
fbd daemons list --json

# Check health and version mismatches
fbd daemons health
fbd daemons health --json

# View daemon logs
fbd daemons logs . -n 100

# Restart all daemons
fbd daemons killall
fbd daemons killall --json
```

## Daemon Info

```bash
fbd info
```

Shows:
- Daemon status (running/stopped)
- Daemon version vs CLI version
- Socket location
- Auto-sync status

## Disabling Daemon

Use `--no-daemon` flag to bypass the daemon:

```bash
fbd --no-daemon ready
fbd --no-daemon list
```

**When to disable:**
- Git worktrees (required)
- CI/CD pipelines
- Resource-constrained environments
- Debugging sync issues

## Event-Driven Mode (Experimental)

Event-driven mode replaces 5-second polling with instant reactivity:

```bash
# Enable globally
export BEADS_DAEMON_MODE=events
fbd daemons killall  # Restart to apply
```

**Benefits:**
- Less than 500ms latency (vs 5s polling)
- ~60% less CPU usage
- Instant sync after changes

**How to verify:**
```bash
fbd info | grep "daemon mode"
```

## Troubleshooting

### Daemon not starting

```bash
# Check if socket exists
ls -la .beads/fbd.sock

# Try direct mode
fbd --no-daemon info

# Restart daemon
fbd daemons killall
fbd info
```

### Version mismatch

After upgrading fbd:

```bash
fbd daemons killall
fbd info  # Should show matching versions
```

### Sync not happening

```bash
# Force sync
fbd sync

# Check daemon logs
fbd daemons logs . -n 50

# Verify git status
git status .beads/
```

### Port/socket conflicts

```bash
# Kill all daemons
fbd daemons killall

# Remove stale socket
rm -f .beads/fbd.sock

# Restart
fbd info
```

## Configuration

Daemon behavior can be configured:

```bash
# Set sync debounce interval
fbd config set daemon.sync_interval 10s

# Disable auto-start
fbd config set daemon.auto_start false

# Set log level
fbd config set daemon.log_level debug
```

See [Configuration](/reference/configuration) for all options.
