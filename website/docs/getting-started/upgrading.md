---
id: upgrading
title: Upgrading
sidebar_position: 4
---

# Upgrading fbd

How to upgrade fbd and keep your projects in sync.

## Checking for Updates

```bash
# Current version
fbd version

# What's new in recent versions
fbd info --whats-new
fbd info --whats-new --json  # Machine-readable
```

## Upgrading

Use the command that matches your install method.

| Install method | Platforms | Command |
|---|---|---|
| Quick install script | macOS, Linux, FreeBSD | `curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh \| bash` |
| PowerShell installer | Windows | `irm https://raw.githubusercontent.com/steveyegge/beads/main/install.ps1 \| iex` |
| Homebrew | macOS, Linux | `brew upgrade fbd` |
| go install | macOS, Linux, FreeBSD, Windows | `go install github.com/steveyegge/fastbeads/cmd/fbd@latest` |
| npm | macOS, Linux, Windows | `npm update -g @beads/fbd` |
| bun | macOS, Linux, Windows | `bun install -g --trust @beads/fbd` |
| From source (Unix shell) | macOS, Linux, FreeBSD | `git pull && go build -o fbd ./cmd/fbd` |

### Quick install script (macOS/Linux/FreeBSD)

```bash
curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash
```

### PowerShell installer (Windows)

```pwsh
irm https://raw.githubusercontent.com/steveyegge/beads/main/install.ps1 | iex
```

### Homebrew

```bash
brew upgrade beads
```

### go install

```bash
go install github.com/steveyegge/fastbeads/cmd/fbd@latest
```

### From Source

```bash
cd beads
git pull
go build -o fbd ./cmd/fbd
sudo mv fbd /usr/local/bin/
```

## After Upgrading

**Important:** After upgrading, update your hooks and restart daemons:

```bash
# 1. Check what changed
fbd info --whats-new

# 2. Update git hooks to match new version
fbd hooks install

# 3. Restart all daemons
fbd daemons killall

# 4. Check for any outdated hooks
fbd info  # Shows warnings if hooks are outdated
```

**Why update hooks?** Git hooks are versioned with fbd. Outdated hooks may miss new auto-sync features or bug fixes.

## Database Migrations

After major upgrades, check for database migrations:

```bash
# Inspect migration plan (AI agents)
fbd migrate --inspect --json

# Preview migration changes
fbd migrate --dry-run

# Apply migrations
fbd migrate

# Migrate and clean up old files
fbd migrate --cleanup --yes
```

## Daemon Version Mismatches

If you see daemon version mismatch warnings:

```bash
# List all running daemons
fbd daemons list --json

# Check for version mismatches
fbd daemons health --json

# Restart all daemons with new version
fbd daemons killall --json
```

## Troubleshooting Upgrades

### Old daemon still running

```bash
fbd daemons killall
```

### Hooks out of date

```bash
fbd hooks install
```

### Database schema changed

```bash
fbd migrate --dry-run
fbd migrate
```

### Import errors after upgrade

Check the import configuration:

```bash
fbd config get import.orphan_handling
fbd import -i .beads/issues.jsonl --orphan-handling allow
```
