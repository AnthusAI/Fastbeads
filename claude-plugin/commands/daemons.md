# fbd daemons - Daemon Management

Manage fbd daemon processes across all repositories and worktrees.

## Synopsis

```bash
fbd daemons <subcommand> [flags]
```

## Description

The `fbd daemons` command provides tools for discovering, monitoring, and managing multiple fbd daemon processes across your system. This is useful when working with multiple repositories or git worktrees.

## Subcommands

### list

List all running fbd daemons with metadata.

```bash
fbd daemons list [--search DIRS] [--json] [--no-cleanup]
```

**Flags:**
- `--search` - Directories to search for daemons (default: home, /tmp, cwd)
- `--json` - Output in JSON format
- `--no-cleanup` - Skip auto-cleanup of stale sockets

**Example:**
```bash
fbd daemons list
fbd daemons list --search /Users/me/projects --json
```

### health

Check health of all fbd daemons and report issues.

```bash
fbd daemons health [--search DIRS] [--json]
```

Reports:
- Stale sockets (dead processes)
- Version mismatches between daemon and CLI
- Unresponsive daemons

**Flags:**
- `--search` - Directories to search for daemons
- `--json` - Output in JSON format

**Example:**
```bash
fbd daemons health
fbd daemons health --json
```

### stop

Stop a specific daemon gracefully.

```bash
fbd daemons stop <workspace-path|pid> [--json]
```

**Arguments:**
- `<workspace-path|pid>` - Workspace path or PID of daemon to stop

**Flags:**
- `--json` - Output in JSON format

**Example:**
```bash
fbd daemons stop /Users/me/projects/myapp
fbd daemons stop 12345
fbd daemons stop /Users/me/projects/myapp --json
```

### restart

Restart a specific daemon gracefully.

```bash
fbd daemons restart <workspace-path|pid> [--search DIRS] [--json]
```

Stops the daemon gracefully, then starts a new one in its place. Useful after upgrading fbd or when a daemon needs to be refreshed.

**Arguments:**
- `<workspace-path|pid>` - Workspace path or PID of daemon to restart

**Flags:**
- `--search` - Directories to search for daemons
- `--json` - Output in JSON format

**Example:**
```bash
fbd daemons restart /Users/me/projects/myapp
fbd daemons restart 12345
fbd daemons restart /Users/me/projects/myapp --json
```

### logs

View logs for a specific daemon.

```bash
fbd daemons logs <workspace-path|pid> [-f] [-n LINES] [--json]
```

**Arguments:**
- `<workspace-path|pid>` - Workspace path or PID of daemon

**Flags:**
- `-f, --follow` - Follow log output (like tail -f)
- `-n, --lines INT` - Number of lines to show from end (default: 50)
- `--json` - Output in JSON format

**Example:**
```bash
fbd daemons logs /Users/me/projects/myapp
fbd daemons logs 12345 -n 100
fbd daemons logs /Users/me/projects/myapp -f
fbd daemons logs 12345 --json
```

### killall

Stop all running fbd daemons.

```bash
fbd daemons killall [--search DIRS] [--force] [--json]
```

Uses escalating shutdown strategy:
1. RPC shutdown (2 second timeout)
2. SIGTERM (3 second timeout)
3. SIGKILL (1 second timeout)

**Flags:**
- `--search` - Directories to search for daemons
- `--force` - Use SIGKILL immediately if graceful shutdown fails
- `--json` - Output in JSON format

**Example:**
```bash
fbd daemons killall
fbd daemons killall --force
fbd daemons killall --json
```

## Common Use Cases

### Version Upgrade

After upgrading fbd, restart all daemons to use the new version:

```bash
fbd daemons health  # Check for version mismatches
fbd daemons killall # Stop all old daemons
# Daemons will auto-start with new version on next fbd command

# Or restart a specific daemon
fbd daemons restart /path/to/workspace
```

### Debugging

Check daemon status and view logs:

```bash
fbd daemons list
fbd daemons health
fbd daemons logs /path/to/workspace -n 100
```

### Cleanup

Remove stale daemon sockets:

```bash
fbd daemons list  # Auto-cleanup happens by default
fbd daemons list --no-cleanup  # Skip cleanup
```

### Multi-Workspace Management

Discover daemons in specific directories:

```bash
fbd daemons list --search /Users/me/projects
fbd daemons health --search /Users/me/work
```

## Troubleshooting

### Stale Sockets

If you see stale sockets (dead process but socket file exists):

```bash
fbd daemons list  # Auto-cleanup removes stale sockets
```

### Version Mismatch

If daemon version != CLI version:

```bash
fbd daemons health  # Identify mismatched daemons
fbd daemons killall # Stop all daemons
# Next fbd command will auto-start new version
```

### Daemon Won't Stop

If graceful shutdown fails:

```bash
fbd daemons killall --force  # Force kill with SIGKILL
```

### Can't Find Daemon

If daemon isn't discovered:

```bash
fbd daemons list --search /path/to/workspace
```

Or check the socket manually:

```bash
ls -la /path/to/workspace/.beads/fbd.sock
```

## See Also

- [fbd daemon](daemon.md) - Start a daemon manually
- [AGENTS.md](../AGENTS.md) - Agent workflow guide
- [README.md](../README.md) - Main documentation
