---
id: labels
title: Labels & Comments
sidebar_position: 5
---

# Labels & Comments

Commands for managing labels and comments.

## Labels

### Adding Labels

```bash
# During creation
fbd create "Task" -l "backend,urgent"

# To existing issue
fbd update bd-42 --add-label urgent
fbd update bd-42 --add-label "backend,security"
```

### Removing Labels

```bash
fbd update bd-42 --remove-label urgent
```

### Listing Labels

```bash
# All labels in use
fbd label list
fbd label list --json

# Issues with specific labels
fbd list --label-any urgent,critical
fbd list --label-all backend,security
```

### Label Conventions

Suggested label categories:

| Category | Examples | Purpose |
|----------|----------|---------|
| Type | `bug`, `feature`, `docs` | Issue classification |
| Priority | `urgent`, `critical` | Urgency markers |
| Area | `backend`, `frontend`, `api` | Code area |
| Status | `blocked`, `needs-review` | Workflow state |
| Size | `small`, `medium`, `large` | Effort estimate |

## Comments

### Adding Comments

```bash
fbd comment add bd-42 "Working on this now"
fbd comment add bd-42 --message "Found the bug in auth.go:45"
```

### Listing Comments

```bash
fbd comment list bd-42
fbd comment list bd-42 --json
```

### Viewing with Issue

```bash
fbd show bd-42 --full  # Includes comments
```

## Filtering by Labels

### Any Match (OR)

```bash
# Issues with urgent OR critical
fbd list --label-any urgent,critical
```

### All Match (AND)

```bash
# Issues with BOTH backend AND security
fbd list --label-all backend,security
```

### Combined Filters

```bash
# Open bugs with urgent label
fbd list --status open --type bug --label-any urgent --json
```

## Bulk Operations

### Add Label to Multiple Issues

```bash
# Using shell
for id in bd-42 bd-43 bd-44; do
  fbd update $id --add-label "sprint-1"
done
```

### Find and Label

```bash
# Label all open bugs as needs-triage
fbd list --status open --type bug --json | \
  jq -r '.[].id' | \
  xargs -I {} fbd update {} --add-label needs-triage
```

## Best Practices

1. **Keep labels lowercase** - `backend` not `Backend`
2. **Use hyphens for multi-word** - `needs-review` not `needs_review`
3. **Be consistent** - Establish team conventions
4. **Don't over-label** - 2-4 labels per issue is typical
5. **Review periodically** - Remove unused labels
