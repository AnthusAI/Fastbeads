---
description: Epic management commands
argument-hint: [command]
---

Manage epics (large features composed of multiple issues).

## Available Commands

- **status**: Show epic completion status
  - Shows progress for each epic
  - Lists child issues and their states
  - Calculates completion percentage

- **close-eligible**: Close epics where all children are complete
  - Automatically closes epics when all child issues are done
  - Useful for bulk epic cleanup

## Epic Workflow

1. Create epic: `fbd create "Large Feature" -t epic -p 1`
2. Link subtasks: `fbd dep add bd-10 bd-20 --type parent-child` (epic bd-10 is parent of task bd-20)
3. Track progress: `fbd epic status`
4. Auto-close when done: `fbd epic close-eligible`

Epics use parent-child dependencies to track subtasks.
