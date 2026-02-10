# Multi-Repo Patterns for AI Agents

This guide covers multi-repo workflow patterns specifically for AI agents working with beads.

**For humans**, see [MULTI_REPO_MIGRATION.md](MULTI_REPO_MIGRATION.md) for interactive wizards and detailed setup.

## Quick Reference

### Single MCP Server (Recommended)

AI agents should use **one MCP server instance** that automatically routes to per-project daemons:

```json
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

The MCP server automatically:
- Detects current workspace from working directory
- Routes to correct per-project daemon (`.beads/fbd.sock`)
- Auto-starts daemon if not running
- Maintains complete database isolation

**Architecture:**
```
MCP Server (one instance)
    ↓
Per-Project Daemons (one per workspace)
    ↓
SQLite Databases (complete isolation)
```

### Multi-Repo Config Options

Agents can configure multi-repo behavior via `fbd config`:

```bash
# Auto-routing (detects role: maintainer vs contributor)
fbd config set routing.mode auto
fbd config set routing.maintainer "."
fbd config set routing.contributor "~/.beads-planning"

# Explicit routing (always use default)
fbd config set routing.mode explicit
fbd config set routing.default "."

# Multi-repo aggregation (hydration)
fbd config set repos.primary "."
fbd config set repos.additional "~/repo1,~/repo2,~/repo3"
```

**Check current config:**
```bash
fbd config get routing.mode
fbd config get repos.additional
fbd info --json  # Shows all config
```

## Routing Behavior

### Auto-Routing (OSS Contributor Pattern)

When `routing.mode=auto`, beads detects user role and routes new issues automatically:

**Maintainer (SSH push access):**
```bash
# Git remote: git@github.com:user/repo.git
fbd create "Fix bug" -p 1
# → Creates in current repo (source_repo = ".")
```

**Contributor (HTTPS or no push access):**
```bash
# Git remote: https://github.com/fork/repo.git
fbd create "Fix bug" -p 1
# → Creates in planning repo (source_repo = "~/.beads-planning")
```

**Role detection priority:**
1. Explicit git config: `git config beads.role maintainer|contributor`
2. Git remote URL inspection (SSH = maintainer, HTTPS = contributor)
3. Fallback: contributor

### Explicit Override

Always available regardless of routing mode:

```bash
# Force creation in specific repo
fbd create "Issue" -p 1 --repo /path/to/repo
fbd create "Issue" -p 1 --repo ~/my-planning
```

### Discovered Issue Inheritance

Issues with `discovered-from` dependencies automatically inherit parent's `source_repo`:

```bash
# Parent in current repo
fbd create "Implement auth" -p 1
# → Created as bd-abc (source_repo = ".")

# Discovered issue inherits parent's repo
fbd create "Found bug" -p 1 --deps discovered-from:bd-abc
# → Created with source_repo = "." (same as parent)
```

**Override if needed:**
```bash
fbd create "Issue" -p 1 --deps discovered-from:bd-abc --repo /different/repo
```

## Multi-Repo Hydration

Agents working in multi-repo mode see aggregated issues from multiple repositories:

```bash
# View all issues (current + additional repos)
fbd ready --json
fbd list --json

# Filter by source repository
fbd list --json | jq '.[] | select(.source_repo == ".")'
fbd list --json | jq '.[] | select(.source_repo == "~/.beads-planning")'
```

**How it works:**
1. Beads reads JSONL from all configured repos
2. Imports into unified SQLite database
3. Maintains `source_repo` field for provenance
4. Exports route issues back to correct JSONL files

## Common Patterns

### OSS Contributor Workflow

**Setup:** Human runs `fbd init --contributor` (wizard handles config)

**Agent workflow:**
```bash
# All planning issues auto-route to separate repo
fbd create "Investigate implementation" -p 1
fbd create "Draft RFC" -p 2
# → Created in ~/.beads-planning (never appears in PRs)

# View all work (upstream + planning)
fbd ready
fbd list --json

# Complete work
fbd close plan-42 --reason "Done"

# Git commit/push - no .beads/ pollution in PR ✅
```

### Team Workflow

**Setup:** Human runs `fbd init --team` (wizard handles config)

**Agent workflow:**
```bash
# Shared team planning (committed to repo)
fbd create "Implement feature X" -p 1
# → Created in current repo (visible to team)

# Optional: Personal experiments in separate repo
fbd create "Try alternative" -p 2 --repo ~/.beads-planning-personal
# → Created in personal repo (private)

# View all
fbd ready --json
```

### Multi-Phase Development

**Setup:** Multiple repos for different phases

**Agent workflow:**
```bash
# Phase 1: Planning repo
cd ~/projects/myapp-planning
fbd create "Design auth" -p 1 -t epic

# Phase 2: Implementation repo (views planning + impl)
cd ~/projects/myapp-implementation
fbd ready  # Shows both repos
fbd create "Implement auth backend" -p 1
fbd dep add impl-42 plan-10 --type blocks  # Link across repos
```

## Troubleshooting

### Issues appearing in wrong repository

**Symptom:** `fbd create` routes to unexpected repo

**Check:**
```bash
fbd config get routing.mode
fbd config get routing.maintainer
fbd config get routing.contributor
fbd info --json | jq '.role'
```

**Fix:**
```bash
# Use explicit flag
fbd create "Issue" -p 1 --repo .

# Or reconfigure routing
fbd config set routing.mode explicit
fbd config set routing.default "."
```

### Can't see issues from other repos

**Symptom:** `fbd list` only shows current repo

**Check:**
```bash
fbd config get repos.additional
```

**Fix:**
```bash
# Add missing repos
fbd config set repos.additional "~/repo1,~/repo2"

# Force sync
fbd sync
fbd list --json
```

### Discovered issues in wrong repository

**Symptom:** Issues with `discovered-from` appear in wrong repo

**Explanation:** This is intentional - discovered issues inherit parent's `source_repo`

**Override if needed:**
```bash
fbd create "Issue" -p 1 --deps discovered-from:bd-42 --repo /different/repo
```

### Planning repo polluting PRs

**Symptom:** `~/.beads-planning` changes appear in upstream PRs

**Verify:**
```bash
# Planning repo should be separate
ls -la ~/.beads-planning/.git  # Should exist

# Fork should NOT contain planning issues
cd ~/projects/fork
fbd list --json | jq '.[] | select(.source_repo == "~/.beads-planning")'
# Should be empty

# Check routing
fbd config get routing.contributor  # Should be ~/.beads-planning
```

### Daemon routing to wrong database

**Symptom:** MCP operations affect wrong project

**Cause:** Using multiple MCP server instances (not recommended)

**Fix:**
```json
// RECOMMENDED: Single MCP server
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

The single MCP server automatically routes based on workspace directory.

### Version mismatch after upgrade

**Symptom:** Daemon operations fail after `fbd` upgrade

**Fix:**
```bash
fbd daemons health --json  # Check for mismatches
fbd daemons killall        # Restart all daemons
# Daemons auto-start with new version on next command
```

## Best Practices for Agents

### OSS Contributors
- ✅ Planning issues auto-route to `~/.beads-planning`
- ✅ Never commit `.beads/` in PRs to upstream
- ✅ Use `fbd ready` to see all work (upstream + planning)
- ❌ Don't manually override routing without good reason

### Teams
- ✅ Commit `.beads/issues.jsonl` to shared repo
- ✅ Use `fbd sync` to ensure changes are committed/pushed
- ✅ Link related issues across repos with dependencies
- ❌ Don't gitignore `.beads/` - you lose the git ledger

### Multi-Phase Projects
- ✅ Use clear repo names (`planning`, `impl`, `maint`)
- ✅ Link issues across phases with `blocks` dependencies
- ✅ Use `fbd list --json` to filter by `source_repo`
- ❌ Don't duplicate issues across repos

### General
- ✅ Always use single MCP server (per-project daemons)
- ✅ Check routing config before filing issues
- ✅ Use `fbd info --json` to verify workspace state
- ✅ Run `fbd sync` at end of session
- ❌ Don't assume routing behavior - check config

## Backward Compatibility

Multi-repo mode is fully backward compatible:

**Without multi-repo config:**
```bash
fbd create "Issue" -p 1
# → Creates in .beads/issues.jsonl (single-repo mode)
```

**With multi-repo config:**
```bash
fbd create "Issue" -p 1
# → Auto-routed based on config
# → Old issues in .beads/issues.jsonl still work
```

**Disabling multi-repo:**
```bash
fbd config unset routing.mode
fbd config unset repos.additional
# → Back to single-repo mode
```

## Configuration Reference

### Routing Config

```bash
# Auto-detect role (maintainer vs contributor)
fbd config set routing.mode auto
fbd config set routing.maintainer "."              # Where maintainer issues go
fbd config set routing.contributor "~/.beads-planning"  # Where contributor issues go

# Explicit mode (always use default)
fbd config set routing.mode explicit
fbd config set routing.default "."                 # All issues go here

# Check settings
fbd config get routing.mode
fbd config get routing.maintainer
fbd config get routing.contributor
fbd config get routing.default
```

### Multi-Repo Hydration

```bash
# Set primary repo (optional, default is current)
fbd config set repos.primary "."

# Add additional repos to aggregate
fbd config set repos.additional "~/repo1,~/repo2,~/repo3"

# Check settings
fbd config get repos.primary
fbd config get repos.additional
```

### Verify Configuration

```bash
# Show all config + database path + daemon status
fbd info --json

# Sample output:
{
  "database_path": "/Users/you/projects/myapp/.beads/beads.db",
  "config": {
    "routing": {
      "mode": "auto",
      "maintainer": ".",
      "contributor": "~/.beads-planning"
    },
    "repos": {
      "primary": ".",
      "additional": ["~/repo1", "~/repo2"]
    }
  },
  "daemon": {
    "running": true,
    "pid": 12345,
    "socket": ".beads/fbd.sock"
  }
}
```

## Related Documentation

- **[MULTI_REPO_MIGRATION.md](MULTI_REPO_MIGRATION.md)** - Complete guide for humans with interactive wizards
- **[ROUTING.md](ROUTING.md)** - Technical details of routing implementation
- **[MULTI_REPO_HYDRATION.md](MULTI_REPO_HYDRATION.md)** - Hydration layer internals
- **[AGENTS.md](../AGENTS.md)** - Main AI agent guide
