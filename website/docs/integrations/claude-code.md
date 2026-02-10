---
id: claude-code
title: Claude Code
sidebar_position: 1
---

# Claude Code Integration

How to use beads with Claude Code.

## Setup

### Quick Setup

```bash
fbd setup claude
```

This installs:
- **SessionStart hook** - Runs `fbd prime` on session start
- **PreCompact hook** - Runs `fbd sync` before context compaction

### Manual Setup

Add to your Claude Code hooks configuration:

```json
{
  "hooks": {
    "SessionStart": ["fbd prime"],
    "PreCompact": ["fbd sync"]
  }
}
```

### Verify Setup

```bash
fbd setup claude --check
```

## How It Works

1. **Session starts** → `fbd prime` injects ~1-2k tokens of context
2. **You work** → Use `fbd` CLI commands directly
3. **Session compacts** → `fbd sync` saves work to git
4. **Session ends** → Changes synced via git

## Essential Commands for Agents

### Creating Issues

```bash
# Always include description for context
fbd create "Fix authentication bug" \
  --description="Login fails with special characters in password" \
  -t bug -p 1 --json

# Link discovered issues
fbd create "Found SQL injection" \
  --description="User input not sanitized in query builder" \
  --deps discovered-from:bd-42 --json
```

### Working on Issues

```bash
# Find ready work
fbd ready --json

# Start work
fbd update bd-42 --status in_progress --json

# Complete work
fbd close bd-42 --reason "Fixed in commit abc123" --json
```

### Querying

```bash
# List open issues
fbd list --status open --json

# Show issue details
fbd show bd-42 --json

# Check blocked issues
fbd blocked --json
```

### Syncing

```bash
# ALWAYS run at session end
fbd sync
```

## Best Practices

### Always Use `--json`

```bash
fbd list --json          # Parse programmatically
fbd create "Task" --json # Get issue ID from output
fbd show bd-42 --json    # Structured data
```

### Always Include Descriptions

```bash
# Good
fbd create "Fix auth bug" \
  --description="Login fails when password contains quotes" \
  -t bug -p 1 --json

# Bad - no context for future work
fbd create "Fix auth bug" -t bug -p 1 --json
```

### Link Related Work

```bash
# When you discover issues during work
fbd create "Found related bug" \
  --deps discovered-from:bd-current --json
```

### Sync Before Session End

```bash
# ALWAYS run before ending
fbd sync
```

## Plugin (Optional)

For enhanced UX with slash commands:

```bash
# In Claude Code
/plugin marketplace add steveyegge/beads
/plugin install beads
# Restart Claude Code
```

Adds slash commands:
- `/beads:ready` - Show ready work
- `/beads:create` - Create issue
- `/beads:show` - Show issue
- `/beads:update` - Update issue
- `/beads:close` - Close issue

## Troubleshooting

### Context not injected

```bash
# Check hook setup
fbd setup claude --check

# Manually prime
fbd prime
```

### Changes not syncing

```bash
# Force sync
fbd sync

# Check daemon
fbd info
fbd daemons health
```

### Database not found

```bash
# Initialize beads
fbd init --quiet
```

## See Also

- [MCP Server](/integrations/mcp-server) - For MCP-only environments
- [IDE Setup](/getting-started/ide-setup) - Other editors
