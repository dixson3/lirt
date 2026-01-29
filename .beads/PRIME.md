# Gas Town Worker Context

> **Context Recovery**: Run `gt prime` for full context after compaction or new session.

## The Propulsion Principle (GUPP)

**If you find work on your hook, YOU RUN IT.**

No confirmation. No waiting. No announcements. The hook having work IS the assignment.
This is physics, not politeness. Gas Town is a steam engine - you are a piston.

**Failure mode we're preventing:**
- Agent starts with work on hook
- Agent announces itself and waits for human to say "ok go"
- Human is AFK / trusting the engine to run
- Work sits idle. The whole system stalls.

## Startup Protocol

1. Check your hook: `gt mol status`
2. If work is hooked → EXECUTE (no announcement, no waiting)
3. If hook empty → Check mail: `gt mail inbox`
4. Still nothing? Wait for user instructions

## Key Commands

- `gt prime` - Get full role context (run after compaction)
- `gt mol status` - Check your hooked work
- `gt mail inbox` - Check for messages
- `bd ready` - Find available work (no blockers)
- `bd sync` - Sync beads changes

## Session Close Protocol

Before signaling completion:
1. **Chronicle worthiness check** - Review session for chronicle-worthy moments
   - Decisions made? Insights gained? Patterns recognized? Lessons learned?
   - If yes → `bd create --type task --labels chronicle,<category> --title "Chronicle: ..." --description "..."`
2. **Process chronicle beads (MANDATORY)**
   - `bd list --label chronicle --status open` (check for open beads)
   - `lirt-chronicler` (process into diary entries)
   - `bd list --label chronicle --status open` (must return empty)
3. git status (check what changed)
4. git add <files> (stage code changes)
5. bd sync (commit beads changes)
6. git commit -m "..." (commit code)
7. bd sync (commit any new beads changes)
8. `lirt-push` (enforces chronicle gate, then pushes)
9. `gt done` (submit to merge queue and exit)

**Polecats MUST call `gt done` - this submits work and exits the session.**
**lirt-push will BLOCK if open chronicle beads exist.**
