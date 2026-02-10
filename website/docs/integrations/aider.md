---
id: aider
title: Aider
sidebar_position: 3
---

# Aider Integration

How to use beads with Aider.

## Setup

### Quick Setup

```bash
fbd setup aider
```

This creates/updates `.aider.conf.yml` with beads context.

### Verify Setup

```bash
fbd setup aider --check
```

## Configuration

The setup adds to `.aider.conf.yml`:

```yaml
# Beads integration
read:
  - .beads/issues.jsonl

# Optional: Auto-run fbd prime
auto-commits: false
```

## Workflow

### Start Session

```bash
# Aider will have access to issues via .aider.conf.yml
aider

# Or manually inject context
fbd prime | aider --message-file -
```

### During Work

Use fbd commands alongside aider:

```bash
# In another terminal or after exiting aider
fbd create "Found bug during work" --deps discovered-from:bd-42 --json
fbd update bd-42 --status in_progress
fbd ready
```

### End Session

```bash
fbd sync
```

## Best Practices

1. **Keep issues visible** - Aider reads `.beads/issues.jsonl`
2. **Sync regularly** - Run `fbd sync` after significant changes
3. **Use discovered-from** - Track issues found during work
4. **Document context** - Include descriptions in issues

## Example Workflow

```bash
# 1. Check ready work
fbd ready

# 2. Start aider with issue context
aider --message "Working on bd-42: Fix auth bug"

# 3. Work in aider...

# 4. Create discovered issues
fbd create "Found related bug" --deps discovered-from:bd-42 --json

# 5. Complete and sync
fbd close bd-42 --reason "Fixed"
fbd sync
```

## Troubleshooting

### Config not loading

```bash
# Check config exists
cat .aider.conf.yml

# Regenerate
fbd setup aider
```

### Issues not visible

```bash
# Check JSONL exists
ls -la .beads/issues.jsonl

# Export if missing
fbd export
```

## See Also

- [Claude Code](/integrations/claude-code)
- [IDE Setup](/getting-started/ide-setup)
