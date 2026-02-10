# Protected Branch Workflow Example

This example demonstrates how to use beads with protected branches on platforms like GitHub, GitLab, and Bitbucket.

## Scenario

You have a repository with:
- Protected `main` branch (requires pull requests)
- Multiple developers/AI agents working on issues
- Desire to track issues in git without bypassing branch protection

## Solution

Use beads' separate sync branch feature to commit issue metadata to a dedicated branch (e.g., `beads-metadata`), then periodically merge via pull request.

## Quick Demo

### 1. Setup (One Time)

```bash
# Clone this repo or create a new one
git init my-project
cd my-project

# Initialize beads with separate sync branch
fbd init --branch beads-metadata --quiet

# Verify configuration
fbd config get sync.branch
# Output: beads-metadata
```

### 2. Create Issues (Agent Workflow)

```bash
# AI agent creates issues normally
fbd create "Implement user authentication" -t feature -p 1
fbd create "Add login page" -t task -p 1
fbd create "Write auth tests" -t task -p 2

# Link tasks to parent feature
fbd link bd-XXXXX --blocks bd-YYYYY  # auth blocks login
fbd link bd-XXXXX --blocks bd-ZZZZZ  # auth blocks tests

# Start work
fbd update bd-XXXXX --status in_progress
```

**Note:** Replace `bd-XXXXX` etc. with actual issue IDs created above.

### 3. Auto-Sync (Daemon)

```bash
# Start daemon with auto-commit
fbd daemon start --auto-commit

# All issue changes are now automatically committed to beads-metadata branch
```

Check what's been committed:

```bash
# View commits on sync branch
git log beads-metadata --oneline | head -5

# View diff between main and sync branch
fbd sync --status
```

### 4. Manual Sync (Without Daemon)

If you're not using the daemon:

```bash
# Create or update issues
fbd create "Fix bug in login" -t bug -p 0
fbd update bd-XXXXX --status closed

# Manually flush to sync branch
fbd sync --flush-only

# Verify commit
git log beads-metadata -1
```

### 5. Merge to Main (Human Review)

Option 1: Via pull request (recommended):

```bash
# Push sync branch
git push origin beads-metadata

# Create PR on GitHub
gh pr create --base main --head beads-metadata \
  --title "Update issue metadata" \
  --body "Automated issue tracker updates from beads"

# After PR is approved and merged:
git checkout main
git pull
fbd import  # Import merged changes to database
```

Option 2: Direct merge (if you have push access):

```bash
# Preview merge
fbd sync --merge --dry-run

# Perform merge
fbd sync --merge

# This will:
# - Merge beads-metadata into main
# - Create merge commit
# - Push to origin
# - Import merged changes
```

### 6. Multi-Clone Sync

If you have multiple clones or agents:

```bash
# Clone 1: Create issue
fbd create "New feature" -t feature -p 1
fbd sync --flush-only  # Commit to beads-metadata
git push origin beads-metadata

# Clone 2: Pull changes
git fetch origin beads-metadata
fbd sync --no-push  # Pull from sync branch and import
fbd list  # See the new feature issue
```

## Workflow Summary

```
┌─────────────────┐
│  Agent creates  │
│  or updates     │
│  issues         │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Daemon (or     │
│  manual sync)   │
│  commits to     │
│  beads-metadata │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Periodically   │
│  merge to main  │
│  via PR         │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  All clones     │
│  pull and       │
│  import         │
└─────────────────┘
```

## Directory Structure

When using separate sync branch, your repo will have:

```
my-project/
├── .git/
│   ├── beads-worktrees/       # Hidden worktree directory
│   │   └── beads-metadata/    # Lightweight checkout of sync branch
│   │       └── .beads/
│   │           └── issues.jsonl
│   └── ...
├── .beads/                    # Main beads directory (in your workspace)
│   ├── beads.db               # SQLite database
│   ├── issues.jsonl            # JSONL export
│   └── fbd.sock                # Daemon socket (if running)
├── src/                       # Your application code
│   └── ...
└── README.md
```

**Key points:**
- `.git/beads-worktrees/` is hidden from your main workspace
- Only `.beads/` is checked out in the worktree (sparse checkout)
- Your `src/` code is never affected by beads commits
- Minimal disk overhead (~few MB for worktree)

## Tips

### For Humans

- **Review before merging:** Use `fbd sync --status` to see what changed
- **Batch merges:** Don't need to merge after every issue - merge when convenient
- **PR descriptions:** Link to specific issues in PR body for context

### For AI Agents

- **No workflow changes:** Agents use `fbd create`, `fbd update`, etc. as normal
- **Let daemon handle it:** With `--auto-commit`, agents don't think about sync
- **Session end:** Run `fbd sync` at end of session to ensure everything is committed

### Troubleshooting

**"Merge conflicts in issues.jsonl"**

JSONL is append-only and line-based, so conflicts are rare. If they occur:
1. Both versions are usually valid - keep both lines
2. If same issue updated differently, keep the line with newer `updated_at`
3. After resolving: `fbd import` to update database

**"Worktree doesn't exist"**

The daemon creates it automatically on first commit. To create manually:
```bash
fbd config get sync.branch  # Verify it's set
fbd daemon stop && fbd daemon start          # Daemon will create worktree
```

**"Changes not syncing"**

Make sure:
- `fbd config get sync.branch` returns the same value on all clones
- Daemon is running: `fbd daemon status`
- Both clones have fetched: `git fetch origin beads-metadata`

## Advanced: GitHub Actions Integration

Automate the merge process with GitHub Actions:

```yaml
name: Auto-Merge Beads Metadata
on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  workflow_dispatch:

jobs:
  merge-beads:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: Install fbd
        run: curl -fsSL https://raw.githubusercontent.com/steveyegge/beads/main/scripts/install.sh | bash

      - name: Check for changes
        id: check
        run: |
          git fetch origin beads-metadata
          if git diff --quiet main origin/beads-metadata -- .beads/; then
            echo "has_changes=false" >> $GITHUB_OUTPUT
          else
            echo "has_changes=true" >> $GITHUB_OUTPUT
          fi

      - name: Create PR
        if: steps.check.outputs.has_changes == 'true'
        run: |
          gh pr create --base main --head beads-metadata \
            --title "Update issue metadata" \
            --body "Automated issue tracker updates from beads" \
            || echo "PR already exists"
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## See Also

- [docs/PROTECTED_BRANCHES.md](../../docs/PROTECTED_BRANCHES.md) - Complete guide
- [AGENTS.md](../../AGENTS.md) - Agent integration instructions
- [commands/sync.md](../../claude-plugin/commands/sync.md) - `fbd sync` command reference
