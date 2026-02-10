# Complete Integration Guide: @beads/fbd + Claude Code for Web

This guide shows the complete end-to-end setup for using fbd (beads) in Claude Code for Web via the npm package.

## 🎯 Goal

Enable automatic issue tracking with fbd in every Claude Code for Web session with zero manual setup.

## 📋 Prerequisites

- A git repository with your project
- Claude Code for Web access

## 🚀 Quick Setup (5 Minutes)

### Step 1: Create SessionStart Hook

Create the file `.claude/hooks/session-start.sh` in your repository:

```bash
#!/bin/bash
# Auto-install fbd in every Claude Code for Web session

# Install fbd globally from npm
npm install -g @beads/fbd

# Initialize fbd if not already done
if [ ! -d .beads ]; then
  fbd init --quiet
fi

# Show current work
echo ""
echo "📋 Ready work:"
fbd ready --limit 5 || echo "No ready work found"
```

### Step 2: Make Hook Executable

```bash
chmod +x .claude/hooks/session-start.sh
```

### Step 3: Update AGENTS.md

Add fbd usage instructions to your AGENTS.md file:

```markdown
## Issue Tracking with fbd

This project uses fbd (beads) for issue tracking. It's automatically installed in each session via SessionStart hook.

### Finding Work

```bash
# See what's ready to work on
fbd ready --json | jq '.[0]'

# Get issue details
fbd show <issue-id> --json
```

### Creating Issues

```bash
# Create a new issue
fbd create "Task description" -t task -p 1 --json

# Create a bug
fbd create "Bug description" -t bug -p 0 --json
```

### Working on Issues

```bash
# Update status to in_progress
fbd update <issue-id> --status in_progress

# Add a comment
fbd comments add <issue-id> "Progress update"

# Close when done
fbd close <issue-id> --reason "Description of what was done"
```

### Managing Dependencies

```bash
# Issue A blocks issue B
fbd dep add <issue-b> <issue-a>

# Show dependency tree
fbd dep tree <issue-id>
```

### Best Practices

1. **Always use --json**: Makes output easy to parse programmatically
2. **Create issues proactively**: When you notice work, file it immediately
3. **Link discovered work**: Use `fbd dep add --type discovered-from`
4. **Close with context**: Always provide --reason when closing
5. **Commit .beads/**: The .beads/issues.jsonl file should be committed to git
```

### Step 4: Commit and Push

```bash
git add .claude/hooks/session-start.sh AGENTS.md
git commit -m "Add fbd auto-install for Claude Code for Web"
git push
```

## 🎬 How It Works

### First Session in Claude Code for Web

1. **Session starts** → Claude Code for Web creates fresh Linux VM
2. **Hook runs** → `.claude/hooks/session-start.sh` executes automatically
3. **npm install** → Downloads @beads/fbd package from npm
4. **Postinstall** → Downloads native fbd binary for platform (~17MB)
5. **fbd init** → Imports existing issues from `.beads/issues.jsonl` in git
6. **Ready** → `fbd` command is available, shows ready work

**Time: ~5-10 seconds**

### Subsequent Sessions

Same process, but:
- Git clone pulls existing `.beads/issues.jsonl`
- `fbd init --quiet` imports all existing issues
- Agent picks up right where it left off

**Time: ~5-10 seconds**

## 💡 Usage Patterns

### Pattern 1: Agent Starts Session

```
Agent: Starting new session...
System: Running SessionStart hook...
System: Installing fbd v0.21.5...
System: ✓ fbd installed successfully
System:
System: 📋 Ready work:
System: bd-a1b2 [P1] Fix authentication bug
System: bd-f14c [P1] Add user profile page
System:
Agent: I can see there are 2 ready tasks. Let me work on bd-a1b2 (Fix authentication bug).
Agent: [runs] fbd update bd-a1b2 --status in_progress
```

### Pattern 2: Discovering New Work

```
Agent: While fixing bd-a1b2, I noticed the password reset flow is also broken.
Agent: [runs] fbd create "Fix password reset flow" -t bug -p 1 --json
System: {"id": "bd-3e7a", ...}
Agent: [runs] fbd dep add bd-3e7a bd-a1b2 --type discovered-from
Agent: I've filed bd-3e7a for the password reset issue and linked it to the current work.
```

### Pattern 3: Completing Work

```
Agent: I've fixed the authentication bug. Tests are passing.
Agent: [runs] fbd close bd-a1b2 --reason "Fixed session validation logic. Added tests."
Agent: [runs] git add . && git commit -m "Fix authentication bug (bd-a1b2)"
Agent: Work completed. Ready for next task.
```

## 🔧 Advanced Configuration

### Custom Installation Path

Install to project instead of global:

```bash
#!/bin/bash
# .claude/hooks/session-start.sh

# Install as project dependency
npm install --save-dev @beads/fbd

# Use via npx
npx fbd init --quiet

echo "Use 'npx fbd' for all fbd commands"
```

### Conditional Installation

Only install if not already present:

```bash
#!/bin/bash
# .claude/hooks/session-start.sh

if ! command -v fbd &> /dev/null; then
    echo "Installing fbd..."
    npm install -g @beads/fbd
else
    echo "fbd already available: $(fbd version)"
fi

# Rest of setup...
```

### Silent Installation

Minimal output for cleaner logs:

```bash
#!/bin/bash
# .claude/hooks/session-start.sh

npm install -g @beads/fbd --silent 2>&1 | grep -v "npm WARN"
fbd init --quiet 2>&1 | grep -v "already initialized"
```

## 📊 Benefits

### For Agents

- ✅ **Persistent memory**: Issue context survives session resets
- ✅ **Structured planning**: Dependencies create clear work order
- ✅ **Automatic setup**: No manual intervention needed
- ✅ **Git-backed**: All issues are version controlled
- ✅ **Fast queries**: `fbd ready` finds work instantly

### For Humans

- ✅ **Visibility**: See what agents are working on
- ✅ **Auditability**: Full history of issue changes
- ✅ **Collaboration**: Multiple agents share same issue database
- ✅ **Portability**: Works locally and in cloud sessions
- ✅ **No servers**: Everything is git and SQLite

### vs Markdown TODOs

| Feature | fbd Issues | Markdown TODOs |
|---------|-----------|----------------|
| Dependencies | ✅ 4 types | ❌ None |
| Ready work detection | ✅ Automatic | ❌ Manual |
| Status tracking | ✅ Built-in | ❌ Manual |
| History/audit | ✅ Full trail | ❌ Git only |
| Queries | ✅ SQL-backed | ❌ Text search |
| Cross-session | ✅ Persistent | ⚠️ Markdown only |
| Agent-friendly | ✅ JSON output | ⚠️ Parsing required |

## 🐛 Troubleshooting

### "fbd: command not found"

**Cause**: SessionStart hook didn't run or installation failed

**Fix**:
```bash
# Manually install
npm install -g @beads/fbd

# Verify
fbd version
```

### "Database not found"

**Cause**: `fbd init` wasn't run

**Fix**:
```bash
fbd init
```

### "Issues.jsonl merge conflict"

**Cause**: Two sessions modified issues concurrently

**Fix**: See the main beads TROUBLESHOOTING.md for merge resolution

### Slow Installation

**Cause**: Network latency downloading binary

**Optimize**:
```bash
# Use npm cache
npm config set cache ~/.npm-cache

# Or install as dependency (cached by package-lock.json)
npm install --save-dev @beads/fbd
```

## 📚 Next Steps

1. **Read the full docs**: https://github.com/steveyegge/fastbeads
2. **Try the quickstart**: `fbd quickstart` (interactive tutorial)
3. **Set up MCP**: For local Claude Desktop integration
4. **Explore examples**: https://github.com/steveyegge/fastbeads/tree/main/examples

## 🔗 Resources

- [beads GitHub](https://github.com/steveyegge/fastbeads)
- [npm package](https://www.npmjs.com/package/@beads/fbd)
- [Claude Code docs](https://docs.claude.com/claude-code)
- [SessionStart hooks](https://docs.claude.com/claude-code/hooks)

## 💬 Example Agent Prompt

Add this to your project's system prompt or AGENTS.md:

```
You have access to fbd (beads) for issue tracking. It's automatically installed in each session.

WORKFLOW:
1. Start each session: Check `fbd ready --json` for available work
2. Choose a task: Pick highest priority with no blockers
3. Update status: `fbd update <id> --status in_progress`
4. Work on it: Implement, test, document
5. File new issues: Create issues for any work discovered
6. Link issues: Use `fbd dep add` to track relationships
7. Close when done: `fbd close <id> --reason "what you did"`
8. Commit changes: Include .beads/issues.jsonl in commits

ALWAYS:
- Use --json flags for programmatic parsing
- Create issues proactively (don't let work be forgotten)
- Link related issues with dependencies
- Close issues with descriptive reasons
- Commit .beads/issues.jsonl with code changes

NEVER:
- Use markdown TODOs (use fbd instead)
- Work on blocked issues (check `fbd show <id>` for blockers)
- Close issues without --reason
- Forget to commit .beads/issues.jsonl
```

## 🎉 Success Criteria

After setup, you should see:

✅ New sessions automatically have `fbd` available
✅ Agents use `fbd` for all issue tracking
✅ Issues persist across sessions via git
✅ Multiple agents can collaborate on same issues
✅ No manual installation required

## 🆘 Support

- [File an issue](https://github.com/steveyegge/fastbeads/issues)
- [Read the FAQ](https://github.com/steveyegge/fastbeads/blob/main/FAQ.md)
- [Check troubleshooting](https://github.com/steveyegge/fastbeads/blob/main/TROUBLESHOOTING.md)
