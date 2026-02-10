---
id: issues
title: Issue Commands
sidebar_position: 3
---

# Issue Commands

Commands for managing issues.

## fbd create

Create a new issue.

```bash
fbd create <title> [flags]
```

**All flags:**
```bash
--type, -t        Issue type (bug|feature|task|epic|chore)
--priority, -p    Priority 0-4
--description, -d Detailed description
--design          Design notes
--acceptance      Acceptance criteria
--notes           Additional notes
--labels, -l      Comma-separated labels
--parent          Parent issue ID
--deps            Dependencies (type:id format)
--assignee        Assigned user
--json            JSON output
```

**Examples:**
```bash
# Bug with high priority
fbd create "Login fails with special chars" -t bug -p 1

# Feature with description
fbd create "Add export to PDF" -t feature -p 2 \
  --description="Users want to export reports as PDF files"

# Feature with design, acceptance, and notes
fbd create "Implement user authentication" -t feature -p 1 \
  --description="Add JWT-based authentication" \
  --design="Use bcrypt for password hashing, JWT for sessions" \
  --acceptance="All tests pass, security audit complete" \
  --notes="Consider rate limiting for login attempts"

# Task with labels
fbd create "Update CI config" -t task -l "ci,infrastructure"

# Epic with children
fbd create "Auth System" -t epic -p 1
fbd create "Design login UI" --parent bd-42
fbd create "Implement backend" --parent bd-42

# Discovered issue
fbd create "Found SQL injection" -t bug -p 0 \
  --deps discovered-from:bd-42 --json
```

## fbd show

Display issue details.

```bash
fbd show <id>... [flags]
```

**Flags:**
```bash
--full        Show all fields including comments
--json        JSON output
```

**Examples:**
```bash
fbd show bd-42
fbd show bd-42 --full
fbd show bd-42 bd-43 bd-44 --json
```

## fbd update

Update issue fields.

```bash
fbd update <id> [flags]
```

**All flags:**
```bash
--status          New status (open|in_progress|closed)
--priority        New priority (0-4)
--title           New title
--description     New description
--type            New type
--add-label       Add label(s)
--remove-label    Remove label(s)
--assignee        New assignee
--json            JSON output
```

**Examples:**
```bash
# Start work
fbd update bd-42 --status in_progress

# Escalate priority
fbd update bd-42 --priority 0 --add-label urgent

# Change title and description
fbd update bd-42 --title "New title" --description="Updated description"

# Multiple changes
fbd update bd-42 --status in_progress --priority 1 --add-label "in-review" --json
```

## fbd close

Close an issue.

```bash
fbd close <id> [flags]
```

**Flags:**
```bash
--reason    Closure reason (stored in comment)
--json      JSON output
```

**Examples:**
```bash
fbd close bd-42
fbd close bd-42 --reason "Fixed in commit abc123"
fbd close bd-42 --reason "Duplicate of bd-43" --json
```

## fbd reopen

Reopen a closed issue.

```bash
fbd reopen <id> [flags]
```

**Examples:**
```bash
fbd reopen bd-42
fbd reopen bd-42 --json
```

## fbd delete

Delete an issue.

```bash
fbd delete <id> [flags]
```

**Flags:**
```bash
--force, -f    Skip confirmation
--json         JSON output
```

**Examples:**
```bash
fbd delete bd-42
fbd delete bd-42 -f --json
```

**Note:** Deletions are tracked in `.beads/deletions.jsonl` for sync.

## fbd search

Search issues by text.

```bash
fbd search <query> [flags]
```

**Flags:**
```bash
--status    Filter by status
--type      Filter by type
--json      JSON output
```

**Examples:**
```bash
fbd search "authentication"
fbd search "login bug" --status open
fbd search "API" --type feature --json
```

## fbd duplicates

Find and manage duplicate issues.

```bash
fbd duplicates [flags]
```

**Flags:**
```bash
--auto-merge    Automatically merge all duplicates
--dry-run       Preview without changes
--json          JSON output
```

**Examples:**
```bash
fbd duplicates
fbd duplicates --auto-merge
fbd duplicates --dry-run --json
```

## fbd merge

Merge duplicate issues.

```bash
fbd merge <source>... --into <target> [flags]
```

**Flags:**
```bash
--into      Target issue to merge into
--dry-run   Preview without changes
--json      JSON output
```

**Examples:**
```bash
fbd merge bd-42 bd-43 --into bd-41
fbd merge bd-42 bd-43 --into bd-41 --dry-run --json
```
