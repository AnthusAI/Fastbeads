# Async Gates for Workflow Coordination

> Adapted from ACF beads skill

`fbd gate` provides async coordination primitives for cross-session and external-condition workflows. Gates are **wisps** (ephemeral issues) that block until a condition is met.

---

## Gate Types

| Type | Await Syntax | Use Case |
|------|--------------|----------|
| Human | `human:<prompt>` | Cross-session human approval |
| CI | `gh:run:<id>` | Wait for GitHub Actions completion |
| PR | `gh:pr:<id>` | Wait for PR merge/close |
| Timer | `timer:<duration>` | Deployment propagation delay |
| Mail | `mail:<pattern>` | Wait for matching email |

---

## Creating Gates

```bash
# Human approval gate
fbd gate create --await human:deploy-approval \
  --title "Approve production deploy" \
  --timeout 4h

# CI gate (GitHub Actions)
fbd gate create --await gh:run:123456789 \
  --title "Wait for CI" \
  --timeout 30m

# PR merge gate
fbd gate create --await gh:pr:42 \
  --title "Wait for PR approval" \
  --timeout 24h

# Timer gate (deployment propagation)
fbd gate create --await timer:15m \
  --title "Wait for deployment propagation"
```

**Required options**:
- `--await <spec>` — Gate condition (see types above)
- `--timeout <duration>` — Recommended: prevents forever-open gates

**Optional**:
- `--title <text>` — Human-readable description
- `--notify <recipients>` — Email/beads addresses to notify

---

## Monitoring Gates

```bash
fbd gate list              # All open gates
fbd gate list --all        # Include closed
fbd gate show <gate-id>    # Details for specific gate
fbd gate eval              # Auto-close elapsed/completed gates
fbd gate eval --dry-run    # Preview what would close
```

**Auto-close behavior** (`fbd gate eval`):
- `timer:*` — Closes when duration elapsed
- `gh:run:*` — Checks GitHub API, closes on success/failure
- `gh:pr:*` — Checks GitHub API, closes on merge/close
- `human:*` — Requires explicit `fbd gate approve`

---

## Closing Gates

```bash
# Human gates require explicit approval
fbd gate approve <gate-id>
fbd gate approve <gate-id> --comment "Reviewed and approved by Steve"

# Manual close (any gate)
fbd gate close <gate-id>
fbd gate close <gate-id> --reason "No longer needed"

# Auto-close via evaluation
fbd gate eval
```

---

## Best Practices

1. **Always set timeouts**: Prevents forever-open gates
   ```bash
   fbd gate create --await human:... --timeout 24h
   ```

2. **Clear titles**: Title should indicate what's being gated
   ```bash
   --title "Approve Phase 2: Core Implementation"
   ```

3. **Eval periodically**: Run at session start to close elapsed gates
   ```bash
   fbd gate eval
   ```

4. **Clean up obsolete gates**: Close gates that are no longer needed
   ```bash
   fbd gate close <id> --reason "superseded by new approach"
   ```

5. **Check before creating**: Avoid duplicate gates
   ```bash
   fbd gate list | grep "spec-myfeature"
   ```

---

## Gates vs Issues

| Aspect | Gates (Wisp) | Issues |
|--------|--------------|--------|
| Persistence | Ephemeral (not synced) | Permanent (synced to git) |
| Purpose | Block on external condition | Track work items |
| Lifecycle | Auto-close when condition met | Manual close |
| Visibility | `fbd gate list` | `fbd list` |
| Use case | CI, approval, timers | Tasks, bugs, features |

Gates are designed to be temporary coordination primitives—they exist only until their condition is satisfied.

---

## Troubleshooting

### Gate won't close

```bash
# Check gate details
fbd gate show <gate-id>

# For gh:run gates, verify the run exists
gh run view <run-id>

# Force close if stuck
fbd gate close <gate-id> --reason "manual override"
```

### Can't find gate ID

```bash
# List all gates (including closed)
fbd gate list --all

# Search by title pattern
fbd gate list | grep "Phase 2"
```

### CI run ID detection fails

```bash
# Check GitHub CLI auth
gh auth status

# List runs manually
gh run list --branch <branch>

# Use specific workflow
gh run list --workflow ci.yml --branch <branch>
```
