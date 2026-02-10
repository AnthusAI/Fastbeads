---
id: dependencies
title: Dependency Commands
sidebar_position: 4
---

# Dependency Commands

Commands for managing issue dependencies.

## fbd dep add

Add a dependency between issues.

```bash
fbd dep add <dependent> <dependency> [flags]
```

**Semantics:** `<dependent>` depends on `<dependency>` (dependency blocks dependent).

**Flags:**
```bash
--type    Dependency type (blocks|related|discovered-from)
--json    JSON output
```

**Examples:**
```bash
# bd-2 depends on bd-1 (bd-1 blocks bd-2)
fbd dep add bd-2 bd-1

# Soft relationship
fbd dep add bd-2 bd-1 --type related

# JSON output
fbd dep add bd-2 bd-1 --json
```

## fbd dep remove

Remove a dependency.

```bash
fbd dep remove <dependent> <dependency> [flags]
```

**Examples:**
```bash
fbd dep remove bd-2 bd-1
fbd dep remove bd-2 bd-1 --json
```

## fbd dep tree

Display dependency tree.

```bash
fbd dep tree <id> [flags]
```

**Flags:**
```bash
--depth    Maximum depth to display
--json     JSON output
```

**Examples:**
```bash
fbd dep tree bd-42
fbd dep tree bd-42 --depth 3
fbd dep tree bd-42 --json
```

**Output:**
```
Dependency tree for bd-42:

> bd-42: Add authentication [P2] (open)
  > bd-41: Create API [P2] (open)
    > bd-40: Set up database [P1] (closed)
```

## fbd dep cycles

Detect circular dependencies.

```bash
fbd dep cycles [flags]
```

**Flags:**
```bash
--json    JSON output
```

**Examples:**
```bash
fbd dep cycles
fbd dep cycles --json
```

## fbd ready

Show issues with no blockers.

```bash
fbd ready [flags]
```

**Flags:**
```bash
--priority    Filter by priority
--type        Filter by type
--label       Filter by label
--json        JSON output
```

**Examples:**
```bash
fbd ready
fbd ready --priority 1
fbd ready --type bug
fbd ready --json
```

**Output:**
```
Ready work (3 issues with no blockers):

1. [P1] bd-40: Set up database
2. [P2] bd-45: Write tests
3. [P3] bd-46: Update docs
```

## fbd blocked

Show blocked issues and their blockers.

```bash
fbd blocked [flags]
```

**Flags:**
```bash
--json    JSON output
```

**Examples:**
```bash
fbd blocked
fbd blocked --json
```

**Output:**
```
Blocked issues (2 issues):

bd-42: Add authentication
  Blocked by: bd-41 (open)

bd-41: Create API
  Blocked by: bd-40 (in_progress)
```

## fbd relate

Create a soft relationship between issues.

```bash
fbd relate <id1> <id2> [flags]
```

**Examples:**
```bash
fbd relate bd-42 bd-43
fbd relate bd-42 bd-43 --json
```

## fbd duplicate

Mark an issue as duplicate.

```bash
fbd duplicate <id> --of <canonical> [flags]
```

**Examples:**
```bash
fbd duplicate bd-43 --of bd-42
fbd duplicate bd-43 --of bd-42 --json
```

## fbd supersede

Mark an issue as superseding another.

```bash
fbd supersede <old> --with <new> [flags]
```

**Examples:**
```bash
fbd supersede bd-42 --with bd-50
fbd supersede bd-42 --with bd-50 --json
```

## Understanding Dependencies

### Blocking vs Non-blocking

| Type | Blocks Ready Queue | Use Case |
|------|-------------------|----------|
| `blocks` | Yes | Hard dependency |
| `parent-child` | No | Epic/subtask hierarchy |
| `discovered-from` | No | Track origin |
| `related` | No | Soft link |
| `duplicates` | No | Mark duplicate |
| `supersedes` | No | Version chain |

### Dependency Direction

```bash
# bd-2 depends on bd-1
# Meaning: bd-1 must complete before bd-2 can start
fbd dep add bd-2 bd-1

# After bd-1 closes:
fbd close bd-1
fbd ready  # bd-2 now appears
```

### Avoiding Cycles

```bash
# Check before adding complex dependencies
fbd dep cycles

# If cycle detected, remove one dependency
fbd dep remove bd-A bd-B
```
