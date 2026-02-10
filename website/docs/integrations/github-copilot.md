---
id: github-copilot
title: GitHub Copilot
sidebar_position: 4
---

# GitHub Copilot Integration

How to use beads with GitHub Copilot in VS Code.

## Setup

### Quick Setup

1. Install beads-mcp:
   ```bash
   uv tool install beads-mcp
   ```

2. Create `.vscode/mcp.json` in your project:
   ```json
   {
     "servers": {
       "beads": {
         "command": "beads-mcp"
       }
     }
   }
   ```

   **For all projects:** Add to VS Code user-level MCP config:

   | Platform | Path |
   |----------|------|
   | macOS | `~/Library/Application Support/Code/User/mcp.json` |
   | Linux | `~/.config/Code/User/mcp.json` |
   | Windows | `%APPDATA%\Code\User\mcp.json` |

   ```json
   {
     "servers": {
       "beads": {
         "command": "beads-mcp",
         "args": []
       }
     }
   }
   ```

3. Initialize beads:
   ```bash
   fbd init --quiet
   ```

4. Reload VS Code

### Verify Setup

Ask Copilot Chat: "What beads issues are ready to work on?"

## Using Natural Language

With MCP configured, interact naturally:

```
You: Create a bug for the login timeout
Copilot: Created bd-42: Login timeout bug

You: What issues are ready?
Copilot: 3 issues ready: bd-42, bd-99, bd-17

You: Close bd-42, it's fixed
Copilot: Closed bd-42
```

## MCP Tools

| Tool | Description |
|------|-------------|
| `beads_ready` | List unblocked issues |
| `beads_create` | Create new issue |
| `beads_show` | Show issue details |
| `beads_update` | Update issue |
| `beads_close` | Close issue |
| `beads_sync` | Sync to git |
| `beads_dep_add` | Add dependency |
| `beads_dep_tree` | Show dependency tree |

## Copilot Instructions

Optionally add `.github/copilot-instructions.md`:

```markdown
## Issue Tracking

This project uses **fbd (beads)** for issue tracking.
Run `fbd prime` for workflow context.

Quick reference:
- `fbd ready` - Find unblocked work
- `fbd create "Title" --type task --priority 2` - Create issue
- `fbd close <id>` - Complete work
- `fbd sync` - Sync with git
```

## Troubleshooting

### Tools not appearing

1. Check VS Code 1.96+
2. Verify mcp.json syntax is valid JSON
3. Reload VS Code window
4. Check Output panel for MCP errors

### "beads-mcp not found"

```bash
# Check installation
which beads-mcp

# Reinstall if needed
uv tool install beads-mcp --force
```

### No database found

```bash
fbd init --quiet
```

## FAQ

### Do I need to clone beads?

**No.** Beads is a system-wide CLI tool. Install once, use everywhere. The `.beads/` directory in your project only contains the issue database.

### What about git hooks?

Git hooks are optional. They auto-sync issues but you can skip them during `fbd init` and manually run `fbd sync` instead.

## See Also

- [MCP Server](/integrations/mcp-server) - Detailed MCP configuration
- [Installation](/getting-started/installation) - Full install guide
- [Detailed Copilot Guide](https://github.com/steveyegge/fastbeads/blob/main/docs/COPILOT_INTEGRATION.md) - Comprehensive documentation
