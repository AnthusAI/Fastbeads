# TODO Command

The `fbd todo` command provides a lightweight interface for managing TODO items as task-type issues.

## Philosophy

TODOs in fbd are not a separate tracking system - they are regular task-type issues with convenient shortcuts. This means:

- **No parallel systems**: TODOs use the same storage and sync as all other issues
- **Promotable**: Easy to convert a TODO to a bug/feature when needed
- **Full featured**: TODOs support all fbd features (dependencies, labels, routing)
- **Simple interface**: Quick commands for common TODO workflows

## Quick Start

```bash
# Add a TODO
fbd todo add "Fix the login bug" -p 1

# List TODOs
fbd todo

# Mark TODO as done
fbd todo done <id>
```

## Commands

### `fbd todo` (or `fbd todo list`)

List all open task-type issues.

```bash
fbd todo                  # List open TODOs
fbd todo list            # Same as above
fbd todo list --all      # Show completed TODOs too
fbd todo list --json     # JSON output
```

**Output:**
```
  ○ test-yxg  Fix the login bug                         ● P1  open
  ○ test-ryl  Update documentation                      ● P3  open

Total: 2 TODOs
```

### `fbd todo add <title>`

Create a new TODO item (task-type issue).

```bash
fbd todo add "Fix the login bug"                                # Default P2
fbd todo add "Update docs" -p 3 -d "Add examples"              # With priority and description
fbd todo add "Critical fix" --priority 0 --description "ASAP"  # P0 task
```

**Flags:**
- `-p, --priority <0-4>`: Priority (default: 2)
- `-d, --description <text>`: Description

### `fbd todo done <id> [<id>...]`

Mark one or more TODOs as complete.

```bash
fbd todo done test-abc              # Close one TODO
fbd todo done test-abc test-def     # Close multiple
fbd todo done test-abc --reason "Fixed in PR #42"  # With reason
```

**Flags:**
- `--reason <text>`: Reason for closing (default: "Completed")

## Converting TODOs

TODOs are regular task issues, so you can convert them:

```bash
# Promote TODO to bug
fbd update test-abc --type bug --priority 0

# Add dependencies
fbd dep add test-abc test-def

# Add labels
fbd update test-abc --labels "urgent,frontend"
```

## Viewing TODO Details

Use regular fbd commands:

```bash
fbd show test-abc        # View TODO details
fbd list --type task     # List all tasks (including TODOs)
fbd ready               # See ready TODOs in work queue
```

## Examples

### Daily TODO workflow

```bash
# Morning: add your tasks
fbd todo add "Review PRs"
fbd todo add "Fix CI pipeline" -p 1
fbd todo add "Update changelog" -p 3

# Check what's on your plate
fbd todo

# Complete work
fbd todo done <id>
fbd todo done <id>

# End of day: see what's left
fbd todo
```

### Converting TODO to full issue

```bash
# Start with a quick TODO
fbd todo add "Login is broken"

# Later, realize it's more serious
fbd update <id> --type bug --priority 0 --description "Users can't login, multiple reports"
fbd update <id> --acceptance "Login works for all user types"

# Now it's a full-fledged bug with proper tracking
fbd show <id>
```

## FAQ

**Q: Are TODOs different from tasks?**
A: No, TODOs are just task-type issues. The `fbd todo` command provides shortcuts for common task operations.

**Q: Can TODOs have dependencies?**
A: Yes! Use `fbd dep add <todo-id> <blocks-id>` like any other issue.

**Q: Do TODOs sync with git?**
A: Yes, they're exported to `.beads/issues.jsonl` like all other issues.

**Q: Can I use TODOs with fbd ready?**
A: Yes! `fbd ready` shows all unblocked issues, including task-type TODOs.

**Q: Should I use TODOs or regular tasks?**
A: Use `fbd todo` for quick, informal tasks. Use `fbd create -t task` for tasks that need more context or are part of larger planning.

## Design Rationale

The TODO command follows beads' philosophy of **minimal surface area**:

1. **No new types**: TODOs are task-type issues
2. **No special storage**: Same database and JSONL as everything else
3. **Convenience layer**: Just shortcuts for common operations
4. **Fully compatible**: Works with all fbd features and commands

This ensures:
- No duplicate tracking systems
- No migration needed between TODOs and tasks
- Works with all existing fbd tooling (federation, compaction, routing)
- Simple to understand and maintain
