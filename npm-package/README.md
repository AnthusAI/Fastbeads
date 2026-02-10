# @beads/fbd - Beads Issue Tracker

[![npm version](https://img.shields.io/npm/v/@beads/fbd)](https://www.npmjs.com/package/@beads/fbd)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Give your coding agent a memory upgrade**

Beads is a lightweight memory system for coding agents, using a graph-based issue tracker. This npm package provides easy installation of the native fbd binary for Node.js environments, including Claude Code for Web.

## Installation

```bash
npm install -g @beads/fbd
```

Or as a project dependency:

```bash
npm install --save-dev @beads/fbd
```

## What is Beads?

Beads is an issue tracker designed specifically for AI coding agents. It provides:

- ✨ **Zero setup** - `fbd init` creates project-local database
- 🔗 **Dependency tracking** - Four dependency types (blocks, related, parent-child, discovered-from)
- 📋 **Ready work detection** - Automatically finds issues with no open blockers
- 🤖 **Agent-friendly** - `--json` flags for programmatic integration
- 📦 **Git-versioned** - JSONL records stored in git, synced across machines
- 🌍 **Distributed by design** - Share one logical database via git

## Quick Start

After installation, initialize beads in your project:

```bash
fbd init
```

Then tell your AI agent to use fbd for task tracking instead of markdown:

```bash
echo "Use 'fbd' commands for issue tracking instead of markdown TODOs" >> AGENTS.md
```

Your agent will automatically:
- Create and track issues during work
- Manage dependencies between tasks
- Find ready work with `fbd ready`
- Keep long-term context across sessions

## Common Commands

```bash
# Find ready work
fbd ready --json

# Create an issue
fbd create "Fix bug" -t bug -p 1

# Show issue details
fbd show bd-a1b2

# List all issues
fbd list --json

# Update status
fbd update bd-a1b2 --status in_progress

# Add dependency
fbd dep add bd-f14c bd-a1b2

# Close issue
fbd close bd-a1b2 --reason "Fixed"
```

## Claude Code for Web Integration

To auto-install fbd in Claude Code for Web sessions, add to your SessionStart hook:

```bash
# .claude/hooks/session-start.sh
npm install -g @beads/fbd
fbd init --quiet
```

This ensures fbd is available in every new session without manual setup.

## Platform Support

This package downloads the appropriate native binary for your platform:

- **macOS**: darwin-amd64, darwin-arm64
- **Linux**: linux-amd64, linux-arm64
- **Windows**: windows-amd64

## Full Documentation

For complete documentation, see the [beads GitHub repository](https://github.com/steveyegge/fastbeads):

- [Complete README](https://github.com/steveyegge/fastbeads#readme)
- [Quick Start Guide](https://github.com/steveyegge/fastbeads/blob/main/docs/QUICKSTART.md)
- [Installation Guide](https://github.com/steveyegge/fastbeads/blob/main/docs/INSTALLING.md)
- [FAQ](https://github.com/steveyegge/fastbeads/blob/main/docs/FAQ.md)
- [Troubleshooting](https://github.com/steveyegge/fastbeads/blob/main/docs/TROUBLESHOOTING.md)

## Why npm Package vs WASM?

This npm package wraps the native fbd binary rather than using WebAssembly because:

- ✅ Full SQLite support (no custom VFS needed)
- ✅ All features work identically to native fbd
- ✅ Better performance (native vs WASM overhead)
- ✅ Simpler maintenance

## License

MIT - See [LICENSE](LICENSE) for details.

## Support

- [GitHub Issues](https://github.com/steveyegge/fastbeads/issues)
- [Documentation](https://github.com/steveyegge/fastbeads)
