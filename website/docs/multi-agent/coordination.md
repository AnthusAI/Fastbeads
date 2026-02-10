---
id: coordination
title: Agent Coordination
sidebar_position: 3
---

# Agent Coordination

Patterns for coordinating work between multiple AI agents.

## Work Assignment

### Pinning Work

Assign work to a specific agent:

```bash
# Pin issue to agent
fbd pin bd-42 --for agent-1

# Pin and start work
fbd pin bd-42 --for agent-1 --start

# Unpin work
fbd unpin bd-42
```

### Checking Pinned Work

```bash
# What's on my hook?
fbd hook

# What's on agent-1's hook?
fbd hook --agent agent-1

# JSON output
fbd hook --json
```

## Handoff Patterns

### Sequential Handoff

Agent A completes work, hands off to Agent B:

```bash
# Agent A
fbd close bd-42 --reason "Ready for review"
fbd pin bd-42 --for agent-b

# Agent B picks up
fbd hook  # Sees bd-42
fbd update bd-42 --status in_progress
```

### Parallel Work

Multiple agents work on different issues:

```bash
# Coordinator
fbd pin bd-42 --for agent-a --start
fbd pin bd-43 --for agent-b --start
fbd pin bd-44 --for agent-c --start

# Each agent works independently
# Coordinator monitors progress
fbd list --status in_progress --json
```

### Fan-Out / Fan-In

Split work, then merge:

```bash
# Fan-out
fbd create "Part A" --parent bd-epic
fbd create "Part B" --parent bd-epic
fbd create "Part C" --parent bd-epic

fbd pin bd-epic.1 --for agent-a
fbd pin bd-epic.2 --for agent-b
fbd pin bd-epic.3 --for agent-c

# Fan-in: wait for all parts
fbd dep add bd-merge bd-epic.1 bd-epic.2 bd-epic.3
```

## Agent Discovery

Find available agents:

```bash
# List known agents (if using agent registry)
fbd agents list

# Check agent status
fbd agents status agent-1
```

## Conflict Prevention

### File Reservations

Prevent concurrent edits:

```bash
# Reserve files before editing
fbd reserve auth.go --for agent-1

# Check reservations
fbd reservations list

# Release when done
fbd reserve --release auth.go
```

### Issue Locking

```bash
# Lock issue for exclusive work
fbd lock bd-42 --for agent-1

# Unlock when done
fbd unlock bd-42
```

## Communication Patterns

### Via Comments

```bash
# Agent A leaves note
fbd comment add bd-42 "Completed API, needs frontend integration"

# Agent B reads
fbd show bd-42 --full
```

### Via Labels

```bash
# Mark for review
fbd update bd-42 --add-label "needs-review"

# Agent B filters
fbd list --label-any needs-review
```

## Best Practices

1. **Clear ownership** - Always pin work to specific agent
2. **Document handoffs** - Use comments to explain context
3. **Use labels for status** - `needs-review`, `blocked`, `ready`
4. **Avoid conflicts** - Use reservations for shared files
5. **Monitor progress** - Regular status checks
