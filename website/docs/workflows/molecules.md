---
id: molecules
title: Molecules
sidebar_position: 2
---

# Molecules

Molecules are work graphs created from formulas.

## What is a Molecule?

A molecule is a persistent instance of a formula:
- Contains steps with dependencies
- Tracked in `.beads/` (syncs with git)
- Steps map to issues with parent-child relationships

## Creating Molecules

### From Formula

```bash
# Pour a formula into a molecule
fbd pour <formula-name> [--var key=value]
```

**Example:**
```bash
fbd pour release --var version=1.0.0
```

This creates:
- Parent issue: `bd-xyz` (the molecule root)
- Child issues: `bd-xyz.1`, `bd-xyz.2`, etc. (the steps)

### Listing Molecules

```bash
fbd mol list
fbd mol list --json
```

### Viewing a Molecule

```bash
fbd mol show <molecule-id>
fbd dep tree <molecule-id>  # Shows full hierarchy
```

## Working with Molecules

### Step Dependencies

Steps have `needs` dependencies:

```toml
[[steps]]
id = "implement"
title = "Implement feature"
needs = ["design"]  # Must complete design first
```

The `fbd ready` command respects these:

```bash
fbd ready  # Only shows steps with completed dependencies
```

### Progressing Through Steps

```bash
# Start a step
fbd update bd-xyz.1 --status in_progress

# Complete a step
fbd close bd-xyz.1 --reason "Done"

# Check what's ready next
fbd ready
```

### Viewing Progress

```bash
# See blocked steps
fbd blocked

# See molecule stats
fbd stats
```

## Molecule Lifecycle

```
Formula (template)
    ↓ fbd pour
Molecule (instance)
    ↓ work steps
Completed Molecule
    ↓ optional cleanup
Archived
```

## Advanced Features

### Bond Points

Formulas can define bond points for composition:

```toml
[compose]
[[compose.bond_points]]
id = "entry"
step = "design"
position = "before"
```

### Hooks

Execute actions on step completion:

```toml
[[steps]]
id = "build"
title = "Build project"

[steps.on_complete]
run = "make build"
```

### Pinning Work

Assign molecules to agents:

```bash
# Pin to current agent
fbd pin bd-xyz --start

# Check what's pinned
fbd hook
```

## Example Workflow

```bash
# 1. Create molecule from formula
fbd pour feature-workflow --var name="dark-mode"

# 2. View structure
fbd dep tree bd-xyz

# 3. Start first step
fbd update bd-xyz.1 --status in_progress

# 4. Complete and progress
fbd close bd-xyz.1
fbd ready  # Shows next steps

# 5. Continue until complete
```

## See Also

- [Formulas](/workflows/formulas) - Creating templates
- [Gates](/workflows/gates) - Async coordination
- [Wisps](/workflows/wisps) - Ephemeral workflows
