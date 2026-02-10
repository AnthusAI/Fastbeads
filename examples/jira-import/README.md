# Jira Integration for fbd

Two-way synchronization between Jira and fbd (beads).

## Scripts

| Script | Purpose |
|--------|---------|
| `jira2jsonl.py` | **Import** - Fetch Jira issues into fbd |
| `jsonl2jira.py` | **Export** - Push fbd issues to Jira |

## Overview

These tools enable bidirectional sync between Jira and fbd:

**Import (Jira → fbd):**
1. **Jira REST API** - Fetch issues directly from any Jira instance
2. **JSON Export** - Parse exported Jira issues JSON
3. **fbd config integration** - Read credentials and mappings from `fbd config`

**Export (fbd → Jira):**
1. **Create issues** - Push new fbd issues to Jira
2. **Update issues** - Sync changes to existing Jira issues
3. **Status transitions** - Handle Jira workflow transitions automatically

## Features

### Import (jira2jsonl.py)

- Fetch from Jira Cloud or Server/Data Center
- JQL query support for flexible filtering
- Configurable field mappings (status, priority, type)
- Preserve timestamps, assignees, labels
- Extract issue links as dependencies
- Set `external_ref` for re-sync capability
- Hash-based or sequential ID generation

### Export (jsonl2jira.py)

- Create new Jira issues from fbd issues
- Update existing Jira issues (matched by `external_ref`)
- Handle Jira workflow transitions for status changes
- Reverse field mappings (fbd → Jira)
- Dry-run mode for previewing changes
- Auto-update `external_ref` after creation

## Installation

No dependencies required! Uses Python 3 standard library.

## Quick Start

### Option 1: Using fbd config (Recommended)

Set up your Jira credentials once:

```bash
# Required settings
fbd config set jira.url "https://company.atlassian.net"
fbd config set jira.project "PROJ"
fbd config set jira.api_token "YOUR_API_TOKEN"

# For Jira Cloud, also set username (your email)
fbd config set jira.username "you@company.com"
```

Then import:

```bash
python jira2jsonl.py --from-config | fbd import
```

### Option 2: Using environment variables

```bash
export JIRA_API_TOKEN=your_token
export JIRA_USERNAME=you@company.com  # For Jira Cloud

python jira2jsonl.py \
  --url https://company.atlassian.net \
  --project PROJ \
  | fbd import
```

### Option 3: Command-line arguments

```bash
python jira2jsonl.py \
  --url https://company.atlassian.net \
  --project PROJ \
  --username you@company.com \
  --api-token YOUR_TOKEN \
  | fbd import
```

## Authentication

### Jira Cloud

Jira Cloud requires:
1. **Username**: Your email address
2. **API Token**: Create at https://id.atlassian.com/manage-profile/security/api-tokens

```bash
fbd config set jira.username "you@company.com"
fbd config set jira.api_token "your_api_token"
```

### Jira Server/Data Center

Jira Server/DC can use:
- **Personal Access Token (PAT)** - Just set the token, no username needed
- **Username + Password** - Set both username and password as the token

```bash
# Using PAT (recommended)
fbd config set jira.api_token "your_pat_token"

# Using username/password
fbd config set jira.username "your_username"
fbd config set jira.api_token "your_password"
```

## Usage

### Basic Usage

```bash
# Fetch all issues from a project
python jira2jsonl.py --from-config | fbd import

# Save to file first (recommended for large projects)
python jira2jsonl.py --from-config > issues.jsonl
fbd import -i issues.jsonl --dry-run  # Preview
fbd import -i issues.jsonl             # Import
```

### Filtering Issues

```bash
# Only open issues
python jira2jsonl.py --from-config --state open

# Only closed issues
python jira2jsonl.py --from-config --state closed

# Custom JQL query
python jira2jsonl.py --url https://company.atlassian.net \
  --jql "project = PROJ AND priority = High AND status != Done"
```

### ID Generation Modes

```bash
# Sequential IDs (bd-1, bd-2, ...) - default
python jira2jsonl.py --from-config

# Hash-based IDs (bd-a3f2dd, ...) - matches fbd create
python jira2jsonl.py --from-config --id-mode hash

# Custom hash length (3-8 chars)
python jira2jsonl.py --from-config --id-mode hash --hash-length 4

# Custom prefix
python jira2jsonl.py --from-config --prefix myproject
```

### From JSON File

If you have an exported JSON file:

```bash
python jira2jsonl.py --file issues.json | fbd import
```

## Field Mapping

### Default Mappings

| Jira Field | fbd Field | Notes |
|------------|----------|-------|
| `key` | (internal) | Used for dependency resolution |
| `summary` | `title` | Direct copy |
| `description` | `description` | Direct copy |
| `status.name` | `status` | Mapped via status_map |
| `priority.name` | `priority` | Mapped via priority_map |
| `issuetype.name` | `issue_type` | Mapped via type_map |
| `assignee` | `assignee` | Display name or username |
| `labels` | `labels` | Direct copy |
| `created` | `created_at` | ISO 8601 timestamp |
| `updated` | `updated_at` | ISO 8601 timestamp |
| `resolutiondate` | `closed_at` | ISO 8601 timestamp |
| (computed) | `external_ref` | URL to Jira issue |
| `issuelinks` | `dependencies` | Mapped to blocks/related |
| `parent` | `dependencies` | Mapped to parent-child |

### Status Mapping

Default status mappings (Jira status -> fbd status):

| Jira Status | fbd Status |
|-------------|-----------|
| To Do, Open, Backlog, New | `open` |
| In Progress, In Development, In Review | `in_progress` |
| Blocked, On Hold | `blocked` |
| Done, Closed, Resolved, Complete | `closed` |

Custom mappings via fbd config:

```bash
fbd config set jira.status_map.backlog "open"
fbd config set jira.status_map.in_review "in_progress"
fbd config set jira.status_map.on_hold "blocked"
```

### Priority Mapping

Default priority mappings (Jira priority -> fbd priority 0-4):

| Jira Priority | fbd Priority |
|---------------|-------------|
| Highest, Critical, Blocker | 0 (Critical) |
| High, Major | 1 (High) |
| Medium, Normal | 2 (Medium) |
| Low, Minor | 3 (Low) |
| Lowest, Trivial | 4 (Backlog) |

Custom mappings:

```bash
fbd config set jira.priority_map.urgent "0"
fbd config set jira.priority_map.nice_to_have "4"
```

### Issue Type Mapping

Default type mappings (Jira type -> fbd type):

| Jira Type | fbd Type |
|-----------|---------|
| Bug, Defect | `bug` |
| Story, Feature, Enhancement | `feature` |
| Task, Sub-task | `task` |
| Epic, Initiative | `epic` |
| Technical Task, Maintenance | `chore` |

Custom mappings:

```bash
fbd config set jira.type_map.story "feature"
fbd config set jira.type_map.spike "task"
fbd config set jira.type_map.tech_debt "chore"
```

## Issue Links & Dependencies

Jira issue links are converted to fbd dependencies:

| Jira Link Type | fbd Dependency Type |
|----------------|-------------------|
| Blocks/Is blocked by | `blocks` |
| Parent (Epic/Story) | `parent-child` |
| All others | `related` |

**Note:** Only links to issues included in the import are preserved. Links to issues outside the query results are ignored.

## Re-syncing from Jira

Each imported issue has an `external_ref` field containing the Jira issue URL. On subsequent imports:

1. Issues are matched by `external_ref` first
2. If matched, the existing fbd issue is updated (if Jira is newer)
3. If not matched, a new fbd issue is created

This enables incremental sync:

```bash
# Initial import
python jira2jsonl.py --from-config | fbd import

# Later: import only recent changes
python jira2jsonl.py --from-config \
  --jql "project = PROJ AND updated >= -7d" \
  | fbd import
```

## Examples

### Example 1: Import Active Sprint

```bash
python jira2jsonl.py --url https://company.atlassian.net \
  --jql "project = PROJ AND sprint in openSprints()" \
  | fbd import

fbd ready  # See what's ready to work on
```

### Example 2: Full Project Migration

```bash
# Export all issues
python jira2jsonl.py --from-config > all-issues.jsonl

# Preview import
fbd import -i all-issues.jsonl --dry-run

# Import
fbd import -i all-issues.jsonl

# View stats
fbd stats
```

### Example 3: Sync High Priority Bugs

```bash
python jira2jsonl.py --from-config \
  --jql "project = PROJ AND type = Bug AND priority in (Highest, High)" \
  | fbd import
```

### Example 4: Import with Hash IDs

```bash
# Use hash IDs for collision-free distributed work
python jira2jsonl.py --from-config --id-mode hash | fbd import
```

## Limitations

- **Single assignee**: Jira supports multiple assignees (watchers), fbd supports one
- **Custom fields**: Only standard fields are mapped; custom fields are ignored
- **Attachments**: Not imported
- **Comments**: Not imported (only description)
- **Worklogs**: Not imported
- **Sprints**: Sprint metadata not preserved (use labels or JQL filtering)
- **Components/Versions**: Not mapped to fbd (consider using labels)

## Troubleshooting

### "Authentication failed"

**Jira Cloud:**
- Verify you're using your email as username
- Create a fresh API token at https://id.atlassian.com/manage-profile/security/api-tokens
- Ensure the token has access to the project
- **Silent auth failure**: The Jira API may return HTTP 200 with empty results instead of 401. Check for `X-Seraph-Loginreason: AUTHENTICATED_FAILED` header in responses.

**Jira Server/DC:**
- Try using a Personal Access Token instead of password
- Check that your account has permission to access the project

### "403 Forbidden"

- Check project permissions in Jira
- Verify API token has correct scopes
- Some Jira instances restrict API access by IP

### "400 Bad Request"

- Check JQL syntax
- Verify project key exists
- Check for special characters in JQL (escape with backslash)

### Rate Limits

Jira Cloud has rate limits. For large imports:
- Add delays between requests (not implemented yet)
- Import in batches using JQL date ranges
- Use the `--file` option with a manual export

## API Rate Limits

- **Jira Cloud**: ~100 requests/minute (varies by plan)
- **Jira Server/DC**: Depends on configuration

This script fetches 100 issues per request, so a 1000-issue project requires ~10 API calls.

## Jira API v3 Notes

This script uses the Jira REST API v3 `/rest/api/3/search/jql` endpoint. The older `/rest/api/3/search` endpoint was deprecated (returns HTTP 410 Gone). Two important considerations:

### Explicit Field Selection

The v3 search endpoint returns only issue IDs by default. The script explicitly requests `fields=*all` to retrieve all fields. Without this parameter, you'll get issues with no title, description, or other metadata.

### Atlassian Document Format (ADF)

Jira API v3 returns rich text fields (like `description`) in Atlassian Document Format - a JSON structure rather than plain text or HTML. The script automatically converts ADF to markdown:

**ADF input:**
```json
{"type": "doc", "content": [{"type": "heading", "attrs": {"level": 3}, "content": [{"type": "text", "text": "Overview"}]}]}
```

**Converted output:**
```markdown
### Overview
```

Supported ADF node types: paragraph, heading, bulletList, orderedList, listItem, codeBlock, blockquote, hardBreak, rule, inlineCard, mention, and text nodes.

---

# Export: jsonl2jira.py

Push fbd issues to Jira.

## Export Quick Start

```bash
# Export all issues (create new, update existing)
fbd export | python jsonl2jira.py --from-config

# Create only (don't update existing Jira issues)
fbd export | python jsonl2jira.py --from-config --create-only

# Dry run (preview what would happen)
fbd export | python jsonl2jira.py --from-config --dry-run

# Auto-update fbd with new external_refs
fbd export | python jsonl2jira.py --from-config --update-refs
```

## Export Modes

### Create Only

Only create new Jira issues for fbd issues that don't have an `external_ref`:

```bash
fbd export | python jsonl2jira.py --from-config --create-only
```

### Create and Update

Create new issues AND update existing ones (matched by `external_ref`):

```bash
fbd export | python jsonl2jira.py --from-config
```

### Dry Run

Preview what would happen without making any changes:

```bash
fbd export | python jsonl2jira.py --from-config --dry-run
```

## Workflow Transitions

Jira often requires workflow transitions to change issue status (you can't just set `status=Done`). The export script automatically:

1. Fetches available transitions for each issue
2. Finds a transition that leads to the target status
3. Executes the transition

If no valid transition is found, the status change is skipped with a warning.

## Reverse Field Mappings

For export, you need mappings from fbd → Jira (reverse of import):

```bash
# Status: fbd status -> Jira status name
fbd config set jira.reverse_status_map.open "To Do"
fbd config set jira.reverse_status_map.in_progress "In Progress"
fbd config set jira.reverse_status_map.blocked "Blocked"
fbd config set jira.reverse_status_map.closed "Done"

# Type: fbd type -> Jira issue type name
fbd config set jira.reverse_type_map.bug "Bug"
fbd config set jira.reverse_type_map.feature "Story"
fbd config set jira.reverse_type_map.task "Task"
fbd config set jira.reverse_type_map.epic "Epic"
fbd config set jira.reverse_type_map.chore "Task"

# Priority: fbd priority (0-4) -> Jira priority name
fbd config set jira.reverse_priority_map.0 "Highest"
fbd config set jira.reverse_priority_map.1 "High"
fbd config set jira.reverse_priority_map.2 "Medium"
fbd config set jira.reverse_priority_map.3 "Low"
fbd config set jira.reverse_priority_map.4 "Lowest"
```

If not configured, sensible defaults are used.

## Updating external_ref

After creating a Jira issue, you'll want to link it back to the fbd issue:

```bash
# Option 1: Auto-update with --update-refs flag
fbd export | python jsonl2jira.py --from-config --update-refs

# Option 2: Manual update from script output
fbd export | python jsonl2jira.py --from-config | while read line; do
  bd_id=$(echo "$line" | jq -r '.bd_id')
  ext_ref=$(echo "$line" | jq -r '.external_ref')
  fbd update "$bd_id" --external-ref="$ext_ref"
done
```

## Export Examples

### Example 1: Initial Export to Jira

```bash
# First, export all open issues
fbd list --status open --json | python jsonl2jira.py --from-config --update-refs

# Now those issues have external_ref set
fbd list --status open
```

### Example 2: Sync Changes Back to Jira

```bash
# Export issues modified today
fbd list --json | python jsonl2jira.py --from-config
```

### Example 3: Preview Before Export

```bash
# See what would happen
fbd export | python jsonl2jira.py --from-config --dry-run

# If it looks good, run for real
fbd export | python jsonl2jira.py --from-config --update-refs
```

## Export Limitations

- **Assignee**: Not set (requires Jira account ID lookup)
- **Dependencies**: Not synced to Jira issue links
- **Comments**: Not exported
- **Custom fields**: design, acceptance_criteria, notes not exported
- **Attachments**: Not exported

## Bidirectional Sync Workflow

For ongoing synchronization between Jira and fbd:

```bash
# 1. Pull changes from Jira
python jira2jsonl.py --from-config --jql "project=PROJ AND updated >= -1d" | fbd import

# 2. Do local work in fbd
fbd update bd-xxx --status in_progress
# ... work ...
fbd close bd-xxx

# 3. Push changes to Jira
fbd export | python jsonl2jira.py --from-config

# 4. Repeat daily/weekly
```

## See Also

- [fbd README](../../README.md) - Main documentation
- [GitHub Import Example](../github-import/) - Similar import for GitHub Issues
- [CONFIG.md](../../docs/CONFIG.md) - Configuration documentation
- [Jira REST API docs](https://developer.atlassian.com/cloud/jira/platform/rest/v2/)
