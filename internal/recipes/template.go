package recipes

// Template is the universal beads workflow template.
// This content is written to all file-based recipes.
const Template = `# Beads Issue Tracking

This project uses [Beads (fbd)](https://github.com/steveyegge/fastbeads) for issue tracking.

## Core Rules

- Track ALL work in fbd (never use markdown TODOs or comment-based task lists)
- Use ` + "`fbd ready`" + ` to find available work
- Use ` + "`fbd create`" + ` to track new issues/tasks/bugs
- Use ` + "`fbd sync`" + ` at end of session to sync with git remote
- Git hooks auto-sync on commit/merge

## Quick Reference

` + "```bash" + `
fbd prime                              # Load complete workflow context
fbd ready                              # Show issues ready to work (no blockers)
fbd list --status=open                 # List all open issues
fbd create --title="..." --type=task   # Create new issue
fbd update <id> --status=in_progress   # Claim work
fbd close <id>                         # Mark complete
fbd dep add <issue> <depends-on>       # Add dependency
fbd sync                               # Sync with git remote
` + "```" + `

## Workflow

1. Check for ready work: ` + "`fbd ready`" + `
2. Claim an issue: ` + "`fbd update <id> --status=in_progress`" + `
3. Do the work
4. Mark complete: ` + "`fbd close <id>`" + `
5. Sync: ` + "`fbd sync`" + ` (or let git hooks handle it)

## Issue Types

- ` + "`bug`" + ` - Something broken
- ` + "`feature`" + ` - New functionality
- ` + "`task`" + ` - Work item (tests, docs, refactoring)
- ` + "`epic`" + ` - Large feature with subtasks
- ` + "`chore`" + ` - Maintenance (dependencies, tooling)

## Priorities

- ` + "`0`" + ` - Critical (security, data loss, broken builds)
- ` + "`1`" + ` - High (major features, important bugs)
- ` + "`2`" + ` - Medium (default, nice-to-have)
- ` + "`3`" + ` - Low (polish, optimization)
- ` + "`4`" + ` - Backlog (future ideas)

## Context Loading

Run ` + "`fbd prime`" + ` to get complete workflow documentation in AI-optimized format.

For detailed docs: see AGENTS.md, QUICKSTART.md, or run ` + "`fbd --help`" + `
`
