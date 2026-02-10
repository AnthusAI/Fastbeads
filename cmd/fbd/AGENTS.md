# Agent Instructions

This project uses **fbd** (beads) for issue tracking. Run `fbd onboard` to get started.

## Quick Reference

```bash
fbd ready              # Find available work
fbd show <id>          # View issue details
fbd update <id> --status in_progress  # Claim work
fbd close <id>         # Complete work
fbd sync               # Sync with git
```

## Agent Warning: Interactive Commands

**DO NOT use `fbd edit`** - it opens an interactive editor ($EDITOR) which AI agents cannot use.

Use `fbd update` with flags instead:
```bash
fbd update <id> --description "new description"
fbd update <id> --title "new title"
fbd update <id> --design "design notes"
fbd update <id> --notes "additional notes"
fbd update <id> --acceptance "acceptance criteria"
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   fbd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds

