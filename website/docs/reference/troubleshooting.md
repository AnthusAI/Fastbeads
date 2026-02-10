---
id: troubleshooting
title: Troubleshooting
sidebar_position: 4
---

# Troubleshooting

Common issues and solutions.

## Installation Issues

### `fbd: command not found`

```bash
# Check if installed
which fbd
go list -f {{.Target}} github.com/steveyegge/fastbeads/cmd/fbd

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Or reinstall
go install github.com/steveyegge/fastbeads/cmd/fbd@latest
```

### `zsh: killed fbd` on macOS

CGO/SQLite compatibility issue:

```bash
CGO_ENABLED=1 go install github.com/steveyegge/fastbeads/cmd/fbd@latest
```

### Permission denied

```bash
chmod +x $(which fbd)
```

## Database Issues

### Database not found

```bash
# Initialize beads
fbd init --quiet

# Or specify database
fbd --db .beads/beads.db list
```

### Database locked

```bash
# Stop daemon
fbd daemons killall

# Try again
fbd list
```

### Corrupted database

```bash
# Restore from JSONL
rm .beads/beads.db
fbd import -i .beads/issues.jsonl
```

## Daemon Issues

### Daemon not starting

```bash
# Check status
fbd info

# Remove stale socket
rm -f .beads/fbd.sock

# Restart
fbd daemons killall
fbd info
```

### Version mismatch

After upgrading fbd:

```bash
fbd daemons killall
fbd info
```

### High CPU usage

```bash
# Switch to event-driven mode
export BEADS_DAEMON_MODE=events
fbd daemons killall
```

## Sync Issues

### Changes not syncing

```bash
# Force sync
fbd sync

# Check daemon
fbd info | grep daemon

# Check hooks
fbd hooks status
```

### Import errors

```bash
# Allow orphans
fbd import -i .beads/issues.jsonl --orphan-handling allow

# Check for duplicates after
fbd duplicates
```

### Merge conflicts

```bash
# Use merge driver
fbd init  # Setup merge driver

# Or manual resolution
git checkout --ours .beads/issues.jsonl
fbd import -i .beads/issues.jsonl
fbd sync
```

## Git Hook Issues

### Hooks not running

```bash
# Check if installed
ls -la .git/hooks/

# Reinstall
fbd hooks install
```

### Hook errors

```bash
# Check hook script
cat .git/hooks/pre-commit

# Run manually
.git/hooks/pre-commit
```

## Dependency Issues

### Circular dependencies

```bash
# Detect cycles
fbd dep cycles

# Remove one dependency
fbd dep remove bd-A bd-B
```

### Missing dependencies

```bash
# Check orphan handling
fbd config get import.orphan_handling

# Allow orphans
fbd config set import.orphan_handling allow
```

## Performance Issues

### Slow queries

```bash
# Check database size
ls -lh .beads/beads.db

# Compact if large
fbd admin compact --analyze
```

### High memory usage

```bash
# Reduce cache
fbd config set database.cache_size 1000
```

## Getting Help

### Debug output

```bash
fbd --verbose list
```

### Logs

```bash
fbd daemons logs . -n 100
```

### System info

```bash
fbd info --json
```

### File an issue

```bash
# Include this info
fbd version
fbd info --json
uname -a
```

Report at: https://github.com/steveyegge/fastbeads/issues
