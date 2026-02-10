---
description: Reopen closed issues
argument-hint: [issue-ids...] [--reason]
---

Reopen one or more closed issues.

Sets status to 'open' and clears the closed_at timestamp. Emits a Reopened event.

## Usage

- **Reopen single**: `fbd reopen bd-42`
- **Reopen multiple**: `fbd reopen bd-42 bd-43 bd-44`
- **With reason**: `fbd reopen bd-42 --reason "Found regression"`

More explicit than `fbd update --status open` - specifically designed for reopening workflow.

Common reasons for reopening:
- Regression found
- Requirements changed
- Incomplete implementation
- New information discovered
