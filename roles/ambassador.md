# Ambassador Protocol

**Role**: Multi-gastown coordination and capability-based work routing
**Applies to**: ALL agents working on lirt

## Purpose

The Ambassador ensures that work is routed to gastowns capable of performing it. When a task has capability requirements (e.g., requires macOS, requires GPU), the ambassador protocol determines whether this gastown can handle it or should route it elsewhere.

This enables multiple gastowns (and multiple users) to collaborate on the same project without conflicts.

## Key Concepts

### Town Identity

Each gastown has a unique identifier stored in `.ambassador/town.json`. This ID is used for:
- Diary entry filenames (attribution)
- Bead assignee values (claim ownership)
- Routing targets (explicit routing)

### Capabilities

Gastowns declare what they can do:

| Type | Examples | How Set |
|------|----------|---------|
| OS | darwin, linux, windows | Auto-detected |
| Arch | arm64, amd64 | Auto-detected |
| Tools | go, node, docker, kubectl | Auto-detected |
| Custom | gpu-access, vpn-access, admin | Manual |

### Sync Branch

All gastowns sync beads through a dedicated branch (default: `beads-sync`) rather than main. This allows:
- Parallel work without merge conflicts on main
- Atomic claiming via `bd update --claim`
- Independent bead state from code state

## When to Assess Capability

**Before starting ANY task**, the ambassador protocol applies:

### Triggers for Assessment

| Trigger | Example |
|---------|---------|
| Task has `requires:*` labels | `requires:darwin`, `requires:gpu` |
| User explicitly requests routing | "Route this to the Linux box" |
| Task mentions platform-specific work | "Fix the macOS build" |
| Task requires tools not present | "Run the GPU benchmarks" |

### Assessment Flow

```
1. Identify capability requirements from:
   - Issue labels (requires:*)
   - Task description keywords
   - Explicit user instructions

2. Compare against local capabilities:
   lirt-ambassador assess <issue-id>

3. Decision:
   - All requirements met → Proceed with work
   - Missing requirements → Route to capable gastown
```

## Routing Work

When this gastown cannot handle a task, route it appropriately:

### Route to Any Capable Gastown

```bash
# Add requirement label and release to pool
lirt-ambassador route <issue-id> --capability darwin
```

This adds `requires:darwin` label and sets assignee to empty, allowing any gastown with the `darwin` capability to claim it.

### Route to Specific Gastown

```bash
# Assign directly to another gastown
lirt-ambassador route <issue-id> --to other-town-id
```

This sets the assignee to the specific gastown.

### Create New Routed Bead

If the task isn't already a bead, create one:

```bash
bd create --title "Task: [description]" \
  --type task \
  --priority 2 \
  --add-label requires:darwin \
  --add-label requires:gpu \
  --description "[full context for the receiving gastown]"
bd sync
```

## Claiming Work

When looking for work, filter by capabilities:

```bash
# Find work this gastown can handle
bd ready --unassigned --label requires:darwin
bd ready --unassigned --label requires:arm64

# Claim atomically (double-sync pattern)
bd sync
bd update <issue-id> --claim
bd sync
```

## Capability Matching

### Issue Labels → Required Capabilities

| Label | Meaning |
|-------|---------|
| `requires:darwin` | Must run on macOS |
| `requires:linux` | Must run on Linux |
| `requires:arm64` | Must run on ARM64 architecture |
| `requires:go` | Needs Go toolchain |
| `requires:docker` | Needs Docker |
| `requires:gpu` | Needs GPU access |
| `requires:admin` | Needs admin/sudo access |

### Adding Requirements to Issues

When creating or updating issues that have specific requirements:

```bash
bd update <issue-id> --add-label requires:darwin
bd update <issue-id> --add-label requires:docker
bd sync
```

## Conflict Prevention

### The Double-Sync Pattern

Always sync before AND after claiming:

```bash
bd sync                    # Pull latest state
bd update <id> --claim     # Atomic claim
bd sync                    # Push immediately
```

This ensures:
1. You see if someone else already claimed it
2. Your claim is visible to other gastowns immediately

### Handling Claim Conflicts

If `bd update --claim` fails (someone else claimed first):

```bash
# Find other available work
bd ready --unassigned --label requires:darwin

# Or check who claimed it
bd show <issue-id>
```

## Stale Claim Detection

Periodically check for abandoned work:

```bash
# Find in-progress issues with no updates for 1+ day
bd stale --status in_progress --days 1
```

If found:
1. Check with the claiming gastown if possible
2. If unresponsive, release: `bd update <id> --status open --assignee ""`
3. Re-claim if needed

## Self-Check Before Working

Before starting any task:

1. **Check requirements**: `lirt-ambassador assess <issue-id>`
2. **If capable**: Proceed with `bd update <id> --claim && bd sync`
3. **If not capable**: Route with `lirt-ambassador route <id> --capability <needed>`

## Example Scenarios

### Scenario 1: macOS-Only Task on Linux Gastown

```
Task: "Fix the macOS notarization script"
Labels: requires:darwin

Linux gastown assessment:
  ✗ darwin (missing)

Action: Route to macOS gastown
  lirt-ambassador route issue-123 --capability darwin
```

### Scenario 2: Multi-Requirement Task

```
Task: "Run GPU benchmarks on ARM64"
Labels: requires:arm64, requires:gpu

M1 Mac gastown assessment:
  ✓ arm64
  ✗ gpu (missing)

Action: Route to gastown with both capabilities
  Already has requires:arm64, already has requires:gpu
  bd update issue-456 --assignee "" --status open
  bd sync
  # Wait for capable gastown to claim
```

### Scenario 3: No Special Requirements

```
Task: "Update README documentation"
Labels: (none with requires:)

Any gastown assessment:
  No requirements to check

Action: Claim and proceed
  bd sync && bd update issue-789 --claim && bd sync
```

## Tools Reference

| Tool | Purpose |
|------|---------|
| `lirt-ambassador init` | Initialize gastown |
| `lirt-ambassador status` | Show gastown status |
| `lirt-ambassador caps list` | List capabilities |
| `lirt-ambassador caps add <cap>` | Add capability |
| `lirt-ambassador assess <id>` | Check if can handle issue |
| `lirt-ambassador route <id>` | Route to capable gastown |
| `bd update <id> --claim` | Atomically claim issue |
| `bd ready --unassigned` | Find available work |
| `bd stale --status in_progress` | Find abandoned work |
| `bd sync` | Synchronize with remote |

## Success Metrics

Track ambassador effectiveness:
- Tasks routed vs tasks attempted locally
- Time from task creation to capable gastown claiming
- Claim conflicts (should approach zero)
- Stale claims requiring intervention
