# Agent Instructions

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Project Context

**lirt** is a Go-based CLI tool for Linear following `gh` CLI semantics:
- Go 1.21+ with Cobra, Viper, and GraphQL client patterns
- Linear API integration with AWS CLI-style profile management
- Optimized for bash scripting workflows (<50ms startup, <50MB memory)
- Multiple workspace support via named profiles

---

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

---

## Crew Role Definitions

Each crew member's role definition is stored at:

```
roles/<name>.md
```

| Crew Member | Role File | Purpose |
|-------------|-----------|---------|
| chronicler | `roles/chronicler.md` | Captures insights and decision reasoning in lirt development |

### Adding New Crew Members

When adding a new crew member:
1. Create `roles/<name>.md` with their role definition
2. Add their protocol to this AGENTS.md file
3. Register them in the table above

---

## Chronicler Protocol

The **Chronicler** captures insights in the evolution of James's thinking during lirt development. This is not a changelog—it's a record of *how* work is performed and *why* decisions are made.

### Chronicle-Worthy Triggers

**ALL agents working on lirt** must watch for these patterns:

| Trigger | Example |
|---------|---------|
| **Significant decision** | Choosing an approach for GraphQL client design, rejecting an alternative |
| **Insight emerges** | Realizing a Go pattern that simplifies implementation |
| **Pattern recognized** | Noticing recurring themes in CLI UX or Linear API integration |
| **Course correction** | Changing API design based on performance measurements |
| **Lesson learned** | What worked or didn't in Go testing, and why |

### What NOT to Chronicle

- Routine task completion (that's for beads)
- Technical implementation details (that's for code/docs)
- Every conversation—only the meaningful ones

### Diary Structure

```
diary/
├── _index.md                                    # Chronological index
└── {YY-MM-DD}.{HH-MM-TZ}.{town}.{topic}.md      # Individual entries
```

**Note:** Include town ID to prevent collisions in multi-town setups. Get it with `lirt-town-id get`.

### Creating a Diary Entry

When you identify something chronicle-worthy, **invoke the lirt-chronicler agent**:

```
Ask lirt-chronicler to document [context]:

[Provide rich context:]
- What happened (decision/insight/pattern/correction/lesson)
- The reasoning behind it
- Key trade-offs or considerations
- Implications for future work
```

The lirt-chronicler agent will:
1. Create the properly formatted diary entry
2. Get town ID automatically (`lirt-town-id get`)
3. Use the template from `roles/chronicler.md`
4. Update `diary/_index.md` (newest entry first)

Your job is **detection**, not **execution**.

### When to Chronicle

**Primary triggers:**
1. **Before every push** — The Pre-Push Chronicle Gate requires reviewing all commits being pushed for chronicle-worthy items. This is the most critical trigger because it catches mid-session insights that would otherwise be missed.
2. **Completion of a coherent set of tasks** — either all tasks complete, or the Overseer decides to defer remaining work.

**Do NOT wait for:**
- Context compaction (too late - context may be lost)
- Being reminded by hooks (hooks are a safety net, not the trigger)
- Session end (chronicle mid-session if pushing mid-session work)

**When in doubt:** Ask the Overseer. It's better to ask than to miss capturing an insight.

### Prompt Template

At task completion, ask:
> "This seems like a chronicle-worthy session. Should I create a diary entry about [topic]?"

---

## Pre-Push Chronicle Gate

**Before EVERY push** — not just at session end — you MUST review uncommitted work for chronicle-worthy items. This gate applies to every `git push` or `lirt-push` invocation.

### The Gate

Before pushing, stop and ask:

1. **Chronicle-worthy?** Review all commits being pushed:
   - New capability, crew member, or infrastructure added?
   - Significant decision with reasoning worth preserving?
   - Insight about Go CLI development or Linear API integration?
   - Pattern or lesson learned?
   - GraphQL client design decision?

   If yes → Create diary entry in `diary/` before pushing

### Why This Exists

Chronicle items are easy to miss mid-session. You finish a task, the user says "commit and push," and the insight slips away uncommemorated. The gate prevents this: **no push without a conscious review of what's being shipped.**

The `lirt-push` script and PreToolUse hook both enforce this gate, but the real enforcement is the habit: before every push, scan your commits and ask "is anything here worth remembering?"

---

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds:
   ```bash
   go test ./...
   golangci-lint run
   go vet ./...
   ```
3. **Update issue status** - Close finished work, update in-progress items
4. **Run the Pre-Push Chronicle Gate** (see above) - Review ALL work from this session, not just the last commit
5. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   bd sync
   git pull --rebase
   # If conflict in .beads/issues.jsonl:
   #   bd resolve-conflicts
   #   git add .beads/issues.jsonl
   #   git rebase --continue
   lirt-push   # Use lirt-push instead of git push - enforces chronicle check
   git status  # MUST show "up to date with origin"
   ```
6. **Clean up** - Clear stashes, prune remote branches
7. **Verify** - All changes committed AND pushed
8. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
- NEVER push without running the Pre-Push Chronicle Gate first

---

## Hook Setup

When rigging lirt in a new Gas Town, install the required hooks:

```bash
lirt-setup-hooks          # Install hooks in workspace
lirt-setup-hooks --check  # Verify hooks are installed
```

This installs:
- **PreCompact** - Reminds agent to review for chronicle-worthy items before compaction
- **PreToolUse** - Enforces Pre-Push Chronicle Gate before every push
