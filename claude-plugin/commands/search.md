---
description: Search issues by text query
argument-hint: <query> [--status] [--label] [--assignee]
---

Search issues across title, description, and ID with a simple text query.

**Note:** The `search` command is optimized for quick text searches and uses less context than `list` when accessed via MCP. For advanced filtering options, use `fbd list`.

## Basic Usage

```bash
fbd search "authentication bug"
fbd search login --status open
fbd search database --label backend
fbd search "bd-5q"  # Search by partial issue ID
```

## How It Works

The search command finds issues where your query appears in **any** of:
- Issue title
- Issue description
- Issue ID (supports partial matching)

Unlike `fbd list`, which requires you to specify which field to search, `fbd search` automatically searches all text fields, making it faster and more intuitive for exploratory searches.

## Filters

- **--status, -s**: Filter by status (open, in_progress, blocked, closed)
- **--assignee, -a**: Filter by assignee
- **--type, -t**: Filter by type (bug, feature, task, epic, chore)
- **--label, -l**: Filter by labels (must have ALL specified labels)
- **--label-any**: Filter by labels (must have AT LEAST ONE)
- **--limit, -n**: Limit number of results (default: 50)
- **--sort**: Sort by field: priority, created, updated, closed, status, id, title, type, assignee
- **--reverse, -r**: Reverse sort order
- **--long**: Show detailed multi-line output for each issue
- **--json**: Output results in JSON format

## Examples

### Basic Search
```bash
# Find all issues mentioning "auth" or "authentication"
fbd search auth

# Search for performance issues
fbd search performance --status open

# Find database-related bugs
fbd search database --type bug
```

### Filtered Search
```bash
# Find open backend issues about login
fbd search login --status open --label backend

# Search Alice's tasks for "refactor"
fbd search refactor --assignee alice --type task

# Find recent bugs (limited to 10 results)
fbd search bug --status open --limit 10
```

### Sorted Output
```bash
# Search bugs sorted by priority (P0 first)
fbd search bug --sort priority

# Search features sorted by most recently updated
fbd search feature --sort updated

# Search issues sorted by priority, lowest first
fbd search refactor --sort priority --reverse
```

### JSON Output
```bash
# Get JSON results for programmatic use
fbd search "api error" --json

# Use with jq for advanced filtering
fbd search memory --json | jq '.[] | select(.priority <= 1)'
```

## Comparison with fbd list

| Command | Best For | Default Limit | Context Usage |
|---------|----------|---------------|---------------|
| `fbd search` | Quick text searches, exploratory queries | 50 | Low (efficient for LLMs) |
| `fbd list` | Advanced filtering, precise queries | None | High (all results) |

**When to use `fbd search`:**
- You want to find issues quickly by keyword
- You're exploring the issue database
- You're using an LLM/MCP and want to minimize context usage

**When to use `fbd list`:**
- You need advanced filters (date ranges, priority ranges, etc.)
- You want all results without a limit
- You need special output formats (digraph, dot)
