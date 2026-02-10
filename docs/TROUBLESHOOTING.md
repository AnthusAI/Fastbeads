# Troubleshooting fbd

Common issues and solutions for fbd users.

## Table of Contents

- [Debug Environment Variables](#debug-environment-variables)
- [Installation Issues](#installation-issues)
- [Antivirus False Positives](#antivirus-false-positives)
- [Database Issues](#database-issues)
- [Git and Sync Issues](#git-and-sync-issues)
- [Ready Work and Dependencies](#ready-work-and-dependencies)
- [Performance Issues](#performance-issues)
- [Agent-Specific Issues](#agent-specific-issues)
- [Platform-Specific Issues](#platform-specific-issues)

## Debug Environment Variables

fbd supports several environment variables for debugging specific subsystems. Enable these when troubleshooting issues or when requested by maintainers.

### Available Debug Variables

| Variable | Purpose | Output Location | Usage |
|----------|---------|----------------|-------|
| `BD_DEBUG` | General debug logging | stderr | Set to any value to enable |
| `BD_DEBUG_RPC` | RPC communication between CLI and daemon | stderr | Set to `1` or `true` |
| `BD_DEBUG_SYNC` | Sync and import timestamp protection | stderr | Set to any value to enable |
| `BD_DEBUG_ROUTING` | Issue routing and multi-repo resolution | stderr | Set to any value to enable |
| `BD_DEBUG_FRESHNESS` | Database file replacement detection | daemon logs | Set to any value to enable |

### Usage Examples

**General debugging:**
```bash
# Enable all general debug logging
export BD_DEBUG=1
fbd ready
```

**RPC communication issues:**
```bash
# Debug daemon communication
export BD_DEBUG_RPC=1
fbd list

# Example output:
# [RPC DEBUG] Connecting to daemon at .beads/fbd.sock
# [RPC DEBUG] Sent request: list (correlation_id=abc123)
# [RPC DEBUG] Received response: 200 OK
```

**Sync conflicts:**
```bash
# Debug timestamp protection during sync
export BD_DEBUG_SYNC=1
fbd sync

# Example output:
# [debug] Protected bd-123: local=2024-01-20T10:00:00Z >= incoming=2024-01-20T09:55:00Z
```

**Routing issues:**
```bash
# Debug issue routing in multi-repo setups
export BD_DEBUG_ROUTING=1
fbd create "Test issue" --rig=planning

# Example output:
# [routing] Rig "planning" -> prefix plan, path /path/to/planning-repo (townRoot=/path/to/town)
# [routing] ID plan-123 matched prefix plan -> /path/to/planning-repo/beads
```

**Database reconnection issues:**
```bash
# Debug database file replacement detection
export BD_DEBUG_FRESHNESS=1
fbd daemon start --foreground

# Example output:
# [freshness] FreshnessChecker: inode changed 27548143 -> 7945906
# [freshness] FreshnessChecker: triggering reconnection
# [freshness] Database file replaced, reconnection triggered

# Or check daemon logs
BD_DEBUG_FRESHNESS=1 fbd daemon restart
fbd daemons logs . -n 100 | grep freshness
```

**Multiple debug flags:**
```bash
# Enable multiple subsystems
export BD_DEBUG=1
export BD_DEBUG_RPC=1
export BD_DEBUG_FRESHNESS=1
fbd daemon start --foreground
```

### Tips

- **Disable after debugging**: Debug logging can be verbose. Disable by unsetting the variable:
  ```bash
  unset BD_DEBUG
  unset BD_DEBUG_RPC
  # etc.
  ```

- **Capture debug output**: Redirect stderr to a file for analysis:
  ```bash
  BD_DEBUG=1 fbd sync 2> debug.log
  ```

- **Daemon logs**: `BD_DEBUG_FRESHNESS` output goes to daemon logs, not stderr:
  ```bash
  # View daemon logs
  fbd daemons logs . -n 200

  # Or directly:
  tail -f .beads/daemon.log
  ```

- **When filing bug reports**: Include relevant debug output to help maintainers diagnose issues faster.

### Related Documentation

- [DAEMON.md](DAEMON.md) - Daemon management and troubleshooting
- [SYNC.md](SYNC.md) - Git sync behavior and conflict resolution
- [ROUTING.md](ROUTING.md) - Multi-repo routing configuration

## Installation Issues

### `fbd: command not found`

fbd is not in your PATH. Either:

```bash
# Check if installed
go list -f {{.Target}} github.com/steveyegge/fastbeads/cmd/fbd

# Add Go bin to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH="$PATH:$(go env GOPATH)/bin"

# Or reinstall
go install github.com/steveyegge/fastbeads/cmd/fbd@latest
```

### Wrong version of fbd running / Multiple fbd binaries in PATH

If `fbd version` shows an unexpected version (e.g., older than what you just installed), you likely have multiple `fbd` binaries in your PATH.

**Diagnosis:**
```bash
# Check all fbd binaries in PATH
which -a fbd

# Example output showing conflict:
# /Users/you/go/bin/fbd        <- From go install (older)
# /opt/homebrew/bin/fbd        <- From Homebrew (newer)
```

**Solution:**
```bash
# Remove old go install version
rm ~/go/bin/fbd

# Or remove mise-managed Go installs
rm ~/.local/share/mise/installs/go/*/bin/fbd

# Verify you're using the correct version
which fbd        # Should show /opt/homebrew/bin/fbd or your package manager path
fbd version      # Should show the expected version
```

**Why this happens:** If you previously installed fbd via `go install`, the binary was placed in `~/go/bin/`. When you later install via Homebrew or another package manager, the old `~/go/bin/fbd` may appear earlier in your PATH, causing the wrong version to run.

**Recommendation:** Choose one installation method (Homebrew recommended) and stick with it. Avoid mixing `go install` with package managers.

### `zsh: killed fbd` or crashes on macOS

Some users report crashes when running `fbd init` or other commands on macOS. This is typically caused by CGO/SQLite compatibility issues.

**Workaround:**
```bash
# Build with CGO enabled
CGO_ENABLED=1 go install github.com/steveyegge/fastbeads/cmd/fbd@latest

# Or if building from source
git clone https://github.com/steveyegge/fastbeads
cd beads
CGO_ENABLED=1 go build -o fbd ./cmd/fbd
sudo mv fbd /usr/local/bin/
```

If you installed via Homebrew, this shouldn't be necessary as the formula already enables CGO. If you're still seeing crashes with the Homebrew version, please [file an issue](https://github.com/steveyegge/fastbeads/issues).

## Antivirus False Positives

### Antivirus software flags fbd as malware

**Symptom**: Kaspersky, Windows Defender, or other antivirus software detects `fbd` or `fbd.exe` as a trojan or malicious software and removes it.

**Common detections**:
- Kaspersky: `PDM:Trojan.Win32.Generic`
- Windows Defender: Various generic trojan detections

**Cause**: This is a **false positive**. Go binaries are commonly flagged by antivirus heuristics because some malware is written in Go. This is a known industry-wide issue affecting many legitimate Go projects.

**Solutions**:

1. **Add fbd to antivirus exclusions** (recommended):
   - Add the fbd installation directory to your antivirus exclusion list
   - This is safe - beads is open source and checksums are provided

2. **Verify file integrity before excluding**:
   ```bash
   # Windows PowerShell
   Get-FileHash fbd.exe -Algorithm SHA256

   # macOS/Linux
   shasum -a 256 fbd
   ```
   Compare with checksums from the [GitHub release page](https://github.com/steveyegge/fastbeads/releases)

3. **Report the false positive**:
   - Help improve detection by reporting to your antivirus vendor
   - Most vendors have false positive submission forms

**Detailed guide**: See [docs/ANTIVIRUS.md](ANTIVIRUS.md) for complete instructions including:
- How to add exclusions for specific antivirus software
- How to report false positives to vendors
- Why Go binaries trigger these detections
- Future plans for code signing

## Database Issues

### `database is locked`

Another fbd process is accessing the database, or SQLite didn't close properly. Solutions:

```bash
# Find and kill hanging processes
ps aux | grep fbd
kill <pid>

# Remove lock files (safe if no fbd processes running)
rm .beads/*.db-journal .beads/*.db-wal .beads/*.db-shm
```

**Note**: fbd uses a pure Go SQLite driver (`modernc.org/sqlite`) for better portability. Under extreme concurrent load (100+ simultaneous operations), you may see "database is locked" errors. This is a known limitation of the pure Go implementation and does not affect normal usage. For very high concurrency scenarios, consider using the CGO-enabled driver or PostgreSQL (planned for future release).

### `fbd init` fails with "directory not empty"

`.beads/` already exists. Options:

```bash
# Use existing database
fbd list  # Should work if already initialized

# Or remove and reinitialize (DESTROYS DATA!)
rm -rf .beads/
fbd init
```

### `failed to import: issue already exists`

You're trying to import issues that conflict with existing ones. Options:

```bash
# Skip existing issues (only import new ones)
fbd import -i issues.jsonl --skip-existing

# Or clear database and re-import everything
rm .beads/*.db
fbd import -i .beads/issues.jsonl
```

### Import fails with missing parent errors

If you see errors like `parent issue bd-abc does not exist` when importing hierarchical issues (e.g., `bd-abc.1`, `bd-abc.2`), this means the parent issue was deleted but children still reference it.

**Quick fix using resurrection:**

```bash
# Auto-resurrect deleted parents from JSONL history
fbd import -i issues.jsonl --orphan-handling resurrect

# Or set as default behavior
fbd config set import.orphan_handling "resurrect"
fbd sync  # Now uses resurrect mode
```

**What resurrection does:**

1. Searches the full JSONL file for the missing parent issue
2. Recreates it as a tombstone (Status=Closed, Priority=4)
3. Preserves the parent's original title and description
4. Maintains referential integrity for hierarchical children
5. Also resurrects dependencies on best-effort basis

**Other handling modes:**

```bash
# Allow orphans (default) - import without validation
fbd config set import.orphan_handling "allow"

# Skip orphans - partial import with warnings
fbd config set import.orphan_handling "skip"

# Strict - fail fast on missing parents
fbd config set import.orphan_handling "strict"
```

**When this happens:**

- Parent issue was deleted using `fbd delete`
- Branch merge where one side deleted the parent
- Manual JSONL editing that removed parent entries
- Database corruption or incomplete import

**Prevention:**

- Use `fbd delete --cascade` to also delete children
- Check for orphans before cleanup: `fbd list --id bd-abc.*`
- Review impact before deleting epic/parent issues

See [CONFIG.md](CONFIG.md#example-import-orphan-handling) for complete configuration documentation.

### Old data returns after reset

**Symptom:** After running `fbd admin reset --force` and `fbd init`, old issues reappear.

**Cause:** `fbd admin reset --force` only removes **local** beads data. Old data can return from:

1. **Remote sync branch** - If you configured a sync branch (via `fbd init --branch` or `fbd config set sync.branch`), old JSONL data may exist on the remote
2. **Git history** - JSONL files committed to git are preserved in history
3. **Other machines** - Other clones may push old data after you reset

**Solution for complete clean slate:**

```bash
# 1. Reset local beads
fbd admin reset --force

# 2. Delete remote sync branch (if configured)
# Check your sync branch name first:
fbd config get sync.branch
# Then delete it from remote:
git push origin --delete <sync-branch-name>
# Common names: beads-sync, beads-metadata

# 3. Remove JSONL from git history (optional, destructive)
# Only do this if you want to completely erase beads history
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch .beads/issues.jsonl' \
  --prune-empty -- --all
git push origin --force --all

# 4. Re-initialize
fbd init
```

**Less destructive alternatives:**

```bash
# Option A: Just delete the sync branch and reinit
fbd admin reset --force
git push origin --delete beads-sync  # or your sync branch name
fbd init

# Option B: Start fresh without sync branch
fbd admin reset --force
fbd init
fbd config set sync.branch ""  # Disable sync branch feature
```

**Note:** The `--hard` and `--skip-init` flags mentioned in [GH#479](https://github.com/steveyegge/fastbeads/issues/479) were never implemented. Use the workarounds above for a complete reset.

**Related:** [GH#922](https://github.com/steveyegge/fastbeads/issues/922)

### Database corruption

**Important**: Distinguish between **logical consistency issues** (ID collisions, wrong prefixes) and **physical SQLite corruption**.

For **physical database corruption** (disk failures, power loss, filesystem errors):

```bash
# Check database integrity
sqlite3 .beads/*.db "PRAGMA integrity_check;"

# If corrupted, reimport from JSONL (source of truth in git)
mv .beads/*.db .beads/*.db.backup
fbd init
fbd import -i .beads/issues.jsonl
```

For **logical consistency issues** (ID collisions from branch merges, parallel workers):

```bash
# This is NOT corruption - use collision resolution instead
fbd import -i .beads/issues.jsonl
```

See [FAQ](FAQ.md#whats-the-difference-between-sqlite-corruption-and-id-collisions) for the distinction.

### Multiple databases detected warning

If you see a warning about multiple `.beads` databases in the directory hierarchy:

```
╔══════════════════════════════════════════════════════════════════════════╗
║ WARNING: 2 beads databases detected in directory hierarchy             ║
╠══════════════════════════════════════════════════════════════════════════╣
║ Multiple databases can cause confusion and database pollution.          ║
║                                                                          ║
║ ▶ /path/to/project/.beads (15 issues)                                   ║
║   /path/to/parent/.beads (32 issues)                                    ║
║                                                                          ║
║ Currently using the closest database (▶). This is usually correct.      ║
║                                                                          ║
║ RECOMMENDED: Consolidate or remove unused databases to avoid confusion. ║
╚══════════════════════════════════════════════════════════════════════════╝
```

This means fbd found multiple `.beads` directories in your directory hierarchy. The `▶` marker shows which database is actively being used (usually the closest one to your current directory).

**Why this matters:**
- Can cause confusion about which database contains your work
- Easy to accidentally work in the wrong database
- May lead to duplicate tracking of the same work

**Solutions:**

1. **If you have nested projects** (intentional):
   - This is fine! fbd is designed to support this
   - Just be aware which database you're using
   - Set `BEADS_DIR` environment variable to point to your `.beads` directory if you want to override the default selection
   - Or use `BEADS_DB` (deprecated) to point directly to the database file

2. **If you have accidental duplicates** (unintentional):
   - Decide which database to keep
   - Export issues from the unwanted database: `cd <unwanted-dir> && fbd export -o backup.jsonl`
   - Remove the unwanted `.beads` directory: `rm -rf <unwanted-dir>/.beads`
   - Optionally import issues into the main database if needed

3. **Override database selection**:
   ```bash
   # Temporarily use specific .beads directory (recommended)
   BEADS_DIR=/path/to/.beads fbd list

   # Or add to shell config for permanent override
   export BEADS_DIR=/path/to/.beads

   # Legacy method (deprecated, points to database file directly)
   BEADS_DB=/path/to/.beads/issues.db fbd list
   export BEADS_DB=/path/to/.beads/issues.db
   ```

**Note**: The warning only appears when fbd detects multiple databases. If you see this consistently and want to suppress it, you're using the correct database (marked with `▶`).

## Git and Sync Issues

### Git merge conflict in `issues.jsonl`

When both sides add issues, you'll get conflicts. Resolution:

1. Open `.beads/issues.jsonl`
2. Look for `<<<<<<< HEAD` markers
3. Most conflicts can be resolved by **keeping both sides**
4. Each line is independent unless IDs conflict
5. For same-ID conflicts, keep the newest (check `updated_at`)

Example resolution:
```bash
# After resolving conflicts manually
git add .beads/issues.jsonl
git commit
fbd import -i .beads/issues.jsonl  # Sync to SQLite
```

See [ADVANCED.md](ADVANCED.md) for detailed merge strategies.

### Git merge conflicts in JSONL

**With hash-based IDs (v0.20.1+), ID collisions don't occur.** Different issues get different hash IDs.

If git shows a conflict in `.beads/issues.jsonl`, it's because the same issue was modified on both branches:

```bash
# Preview what will be updated
fbd import -i .beads/issues.jsonl --dry-run

# Resolve git conflict (keep newer version or manually merge)
git checkout --theirs .beads/issues.jsonl  # Or --ours, or edit manually

# Import updates the database
fbd import -i .beads/issues.jsonl
```

See [ADVANCED.md#handling-git-merge-conflicts](ADVANCED.md#handling-git-merge-conflicts) for details.

### Permission denied on git hooks

Git hooks need execute permissions:

```bash
chmod +x .git/hooks/pre-commit
chmod +x .git/hooks/post-merge
chmod +x .git/hooks/post-checkout
```

### "Branch already checked out" when switching branches

**Symptom:**
```bash
$ git checkout main
fatal: 'main' is already checked out at '/path/to/.git/beads-worktrees/beads-sync'
```

**Cause:** Beads creates git worktrees internally when using the sync-branch feature (configured via `fbd init --branch` or `fbd config set sync.branch`). These worktrees lock the branches they're checked out to.

**Solution:**
```bash
# Remove beads-created worktrees
rm -rf .git/beads-worktrees
rm -rf .git/worktrees/beads-*
git worktree prune

# Now you can checkout the branch
git checkout main
```

**Permanent fix (disable sync-branch):**
```bash
fbd config set sync.branch ""
```

See [WORKTREES.md#beads-created-worktrees-sync-branch](WORKTREES.md#beads-created-worktrees-sync-branch) for details.

### Unexpected worktree directories in .git/

**Symptom:** You notice `.git/beads-worktrees/` or `.git/worktrees/beads-*` directories you didn't create.

**Explanation:** Beads automatically creates these worktrees when using the sync-branch feature to commit issue updates to a separate branch without switching your working directory.

**If you don't want these:**
```bash
# Disable sync-branch feature
fbd config set sync.branch ""

# Clean up existing worktrees
rm -rf .git/beads-worktrees
rm -rf .git/worktrees/beads-*
git worktree prune
```

See [WORKTREES.md](WORKTREES.md) for details on how beads uses worktrees.

### Auto-sync not working

Check if auto-sync is enabled:

```bash
# Check if daemon is running
ps aux | grep "fbd daemon"

# Manually export/import
fbd export -o .beads/issues.jsonl
fbd import -i .beads/issues.jsonl

# Install git hooks for guaranteed sync
fbd hooks install
```

If you disabled auto-sync with `--no-auto-flush` or `--no-auto-import`, remove those flags or use `fbd sync` manually.

## Ready Work and Dependencies

### `fbd ready` shows nothing but I have open issues

Those issues probably have open blockers. Check:

```bash
# See blocked issues
fbd blocked

# Show dependency tree (default max depth: 50)
fbd dep tree <issue-id>

# Limit tree depth to prevent deep traversals
fbd dep tree <issue-id> --max-depth 10

# Remove blocking dependency if needed
fbd dep remove <from-id> <to-id>
```

Remember: Only `blocks` dependencies affect ready work.

### Circular dependency errors

fbd prevents dependency cycles, which break ready work detection. To fix:

```bash
# Detect all cycles
fbd dep cycles

# Remove the dependency causing the cycle
fbd dep remove <from-id> <to-id>

# Or redesign your dependency structure
```

### Dependencies not showing up

Check the dependency type:

```bash
# Show full issue details including dependencies
fbd show <issue-id>

# Visualize the dependency tree
fbd dep tree <issue-id>
```

Remember: Different dependency types have different meanings:
- `blocks` - Hard blocker, affects ready work
- `related` - Soft relationship, doesn't block
- `parent-child` - Hierarchical (child depends on parent)
- `discovered-from` - Work discovered during another issue

## Performance Issues

### Export/import is slow

For large databases (10k+ issues):

```bash
# Export only open issues
fbd export --format=jsonl --status=open -o .beads/issues.jsonl

# Or filter by priority
fbd export --format=jsonl --priority=0 --priority=1 -o critical.jsonl
```

Consider splitting large projects into multiple databases.

### Commands are slow

Check database size and consider compaction:

```bash
# Check database stats
fbd stats

# Preview compaction candidates
fbd admin compact --dry-run --all

# Compact old closed issues
fbd admin compact --days 90
```

### Large JSONL files

If `.beads/issues.jsonl` is very large:

```bash
# Check file size
ls -lh .beads/issues.jsonl

# Remove old closed issues
fbd admin compact --days 90

# Or split into multiple projects
cd ~/project/component1 && fbd init --prefix comp1
cd ~/project/component2 && fbd init --prefix comp2
```

## Agent-Specific Issues

### Agent creates duplicate issues

Agents may not realize an issue already exists. Prevention strategies:

- Have agents search first: `fbd list --json | grep "title"`
- Use labels to mark auto-created issues: `fbd create "..." -l auto-generated`
- Review and deduplicate periodically: `fbd list | sort`
- Use `fbd merge` to consolidate duplicates: `fbd merge bd-2 --into bd-1`

### Agent gets confused by complex dependencies

Simplify the dependency structure:

```bash
# Check for overly complex trees
fbd dep tree <issue-id>

# Remove unnecessary dependencies
fbd dep remove <from-id> <to-id>

# Use labels instead of dependencies for loose relationships
fbd label add <issue-id> related-to-feature-X
```

### Agent can't find ready work

Check if issues are blocked:

```bash
# See what's blocked
fbd blocked

# See what's actually ready
fbd ready --json

# Check specific issue
fbd show <issue-id>
fbd dep tree <issue-id>
```

### MCP server not working

Check installation and configuration:

```bash
# Verify MCP server is installed
pip list | grep beads-mcp

# Check MCP configuration
cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

# Test CLI works
fbd version
fbd ready

# Check for daemon
ps aux | grep "fbd daemon"
```

See [integrations/beads-mcp/README.md](../integrations/beads-mcp/README.md) for MCP-specific troubleshooting.

### Sandboxed environments (Codex, Claude Code, etc.)

**Issue:** Sandboxed environments restrict permissions, preventing daemon control and causing "out of sync" errors.

**Common symptoms:**
- "Database out of sync with JSONL" errors that persist after running `fbd import`
- `fbd daemon stop` fails with "operation not permitted"
- Cannot kill daemon process with `kill <pid>`
- JSONL hash mismatch warnings (bd-160)
- Commands intermittently fail with staleness errors

**Root cause:** The sandbox can't signal/kill the existing daemon process, so the DB stays stale and refuses to import.

---

#### Quick fix: Sandbox mode (auto-detected)

**As of v0.21.1+**, fbd automatically detects sandboxed environments and enables sandbox mode.

When auto-detected, you'll see: `ℹ️  Sandbox detected, using direct mode`

**Manual override** (if auto-detection fails):

```bash
# Explicitly enable sandbox mode
fbd --sandbox ready
fbd --sandbox create "Fix bug" -p 1
fbd --sandbox update bd-42 --status in_progress

# Equivalent to:
fbd --no-daemon --no-auto-flush --no-auto-import <command>
```

**What sandbox mode does:**
- Disables daemon (uses direct SQLite mode)
- Disables auto-export to JSONL
- Disables auto-import from JSONL
- Allows fbd to work in network-restricted environments

**Note:** You'll need to manually sync when outside the sandbox:
```bash
# After leaving sandbox, sync manually
fbd sync
```

---

#### Escape hatches for stuck states

If you're stuck in a "database out of sync" loop with a running daemon you can't stop, use these flags:

**1. Force metadata update (`--force` flag on import)**

When `fbd import` reports "0 created, 0 updated" but staleness persists:

```bash
# Force metadata refresh even when DB appears synced
fbd import --force

# This updates internal metadata tracking without changing issues
# Fixes: stuck state caused by stale daemon cache
```

**Shows:** `Metadata updated (database already in sync with JSONL)`

**2. Skip staleness check (`--allow-stale` global flag)**

Emergency escape hatch to bypass staleness validation:

```bash
# Allow operations on potentially stale data
fbd --allow-stale ready
fbd --allow-stale list --status open

# Shows warning:
# ⚠️  Staleness check skipped (--allow-stale), data may be out of sync
```

**⚠️ Caution:** Use sparingly - you may see incomplete or outdated data.

**3. Use sandbox mode (preferred)**

```bash
# Most reliable for sandboxed environments
fbd --sandbox ready
fbd --sandbox import -i .beads/issues.jsonl
```

---

#### Troubleshooting workflow

If stuck in a sandboxed environment:

```bash
# Step 1: Try sandbox mode (cleanest solution)
fbd --sandbox ready

# Step 2: If you get staleness errors, force import
fbd import --force -i .beads/issues.jsonl

# Step 3: If still blocked, use allow-stale (emergency only)
fbd --allow-stale ready

# Step 4: When back outside sandbox, sync normally
fbd sync
```

---

#### Understanding the flags

| Flag | Purpose | When to use | Risk |
|------|---------|-------------|------|
| `--sandbox` | Disable daemon and auto-sync | Sandboxed environments (Codex, containers) | Low - safe for sandboxes |
| `--force` (import) | Force metadata update | Stuck "0 created, 0 updated" loop | Low - updates metadata only |
| `--allow-stale` | Skip staleness validation | Emergency access to database | **High** - may show stale data |

**Related:**
- See [DAEMON.md](DAEMON.md) for daemon troubleshooting
- See [Claude Code sandboxing documentation](https://www.anthropic.com/engineering/claude-code-sandboxing) for more about sandbox restrictions
- GitHub issue [#353](https://github.com/steveyegge/fastbeads/issues/353) for background

## Platform-Specific Issues

### Windows: Path issues

```pwsh
# Check if fbd.exe is in PATH
where.exe fbd

# Add Go bin to PATH (permanently)
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:USERPROFILE\go\bin",
    [EnvironmentVariableTarget]::User
)

# Reload PATH in current session
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")
```

### Windows: Firewall blocking daemon

The daemon listens on loopback TCP. Allow `fbd.exe` through Windows Firewall:

1. Open Windows Security → Firewall & network protection
2. Click "Allow an app through firewall"
3. Add `fbd.exe` and enable for Private networks
4. Or disable firewall temporarily for testing

### Windows: Controlled Folder Access blocks fbd init

**Symptom:** `fbd init` hangs indefinitely with high CPU usage, and CTRL+C doesn't work.

**Cause:** Windows Controlled Folder Access is blocking `fbd.exe` from creating the `.beads` directory.

**Diagnosis:** Run with verbose flag to see the actual error:
```pwsh
fbd init -v
# Error: failed to create .beads directory: mkdir .beads: The system cannot find the file specified
```

**Solution:** Add `fbd.exe` to the Controlled Folder Access whitelist:

1. Open Windows Security → Virus & threat protection
2. Click "Ransomware protection" → "Manage ransomware protection"
3. Under "Controlled folder access", click "Allow an app through Controlled folder access"
4. Click "Add an allowed app" → "Browse all apps"
5. Navigate to and select `fbd.exe` (typically in `%USERPROFILE%\go\bin\fbd.exe`)
6. Retry `fbd init` - it should work instantly

**Note:** Unlike typical blocked apps, Controlled Folder Access may not show a notification when blocking `fbd init`, making this issue hard to diagnose without the `-v` flag.

### macOS: Gatekeeper blocking execution

If macOS blocks fbd:

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/fbd

# Or allow in System Preferences
# System Preferences → Security & Privacy → General → "Allow anyway"
```

### Linux: Permission denied

If you get permission errors:

```bash
# Make fbd executable
chmod +x /usr/local/bin/fbd

# Or install to user directory
mkdir -p ~/.local/bin
mv fbd ~/.local/bin/
export PATH="$HOME/.local/bin:$PATH"
```

## Getting Help

If none of these solutions work:

1. **Check existing issues**: [GitHub Issues](https://github.com/steveyegge/fastbeads/issues)
2. **Enable debug logging**: `fbd --verbose <command>`
3. **File a bug report**: Include:
   - fbd version: `fbd version`
   - OS and architecture: `uname -a`
   - Error message and full command
   - Steps to reproduce
4. **Join discussions**: [GitHub Discussions](https://github.com/steveyegge/fastbeads/discussions)

## Related Documentation

- **[README.md](../README.md)** - Core features and quick start
- **[ADVANCED.md](ADVANCED.md)** - Advanced features
- **[FAQ.md](FAQ.md)** - Frequently asked questions
- **[INSTALLING.md](INSTALLING.md)** - Installation guide
