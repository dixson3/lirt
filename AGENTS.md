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
| chronicler | `roles/chronicler.md` | Detection protocol for chronicle-worthy moments; execution agent processes beads |
| test-case-identifier | `roles/test-case-identifier.md` | Detection protocol for testable scenarios across all agents |
| code-quality-reporter | `roles/code-quality-reporter.md` | Detection protocol for code quality issues (primarily lirt-code-reviewer) |

### Adding New Crew Members

When adding a new crew member:
1. Create `roles/<name>.md` with their role definition
2. Add their protocol to this AGENTS.md file
3. Register them in the table above

---

## Chronicle Protocol (Beads-Mediated)

The **Chronicle Protocol** captures insights in the evolution of James's thinking during lirt development. This is not a changelog—it's a record of *how* work is performed and *why* decisions are made.

**Protocol Pattern:** Detection (ALL agents) → Beads Queue → Execution (lirt-chronicler)

### Chronicle-Worthy Triggers

**ALL agents working on lirt** must watch for these patterns and create chronicle beads IMMEDIATELY:

| Trigger | Example |
|---------|---------|
| **Decision** | Choosing functional options over builder pattern for Go client config |
| **Insight** | Realizing cursor pagination prevents data loss vs offset pagination |
| **Pattern** | Recognizing table-driven tests as recurring pattern across CLI tests |
| **Correction** | Pivoting GraphQL client design after performance testing |
| **Lesson** | What worked or didn't in Go testing strategies |

### Creating Chronicle Beads

When you identify something chronicle-worthy, create a chronicle bead **IMMEDIATELY** (while context is fresh):

```bash
bd create --title "Chronicle: [specific topic]" \
  --type chronicle \
  --priority 3 \
  --add-label [category] \
  --add-label [topic] \
  --description "[Rich 300-600 word description]"
```

**Category Labels:**
- `decision` - Architectural or design decision
- `insight` - Realization that changed understanding
- `pattern` - Recurring theme recognized
- `correction` - Direction change
- `lesson` - What worked/didn't work

**Topic Labels:**
- `go-idioms` - Go language patterns
- `graphql` - GraphQL/Linear API
- `cli-ux` - CLI user experience
- `performance` - Performance optimization
- `testing` - Testing strategies
- `security` - Security decisions

**Priority:**
- **P2**: Major architectural decision with far-reaching implications
- **P3**: Standard chronicle-worthy moments (default)
- **P4**: Minor insights, can be batched

### Rich Description Template (300-600 words minimum)

```
Type: [decision|insight|pattern|correction|lesson]

Context: [What were you working on when this occurred]

[For decisions] Alternatives Considered:
1. Option A: [description, pros/cons]
2. Option B: [description, pros/cons]
3. [Chosen] Option C: [description, why this won]

[For insights] What Changed:
[What you understood before vs what you understand now]

Reasoning: [Why this choice/insight emerged]

Trade-offs: [What you gained and what you lost with this approach]

Implementation: [File:line references, commit SHAs, code snippets]

Implications: [What this means for future work, patterns established]
```

### Example Chronicle Bead

```bash
bd create --title "Chronicle: Functional options for client configuration" \
  --type chronicle \
  --priority 3 \
  --add-label decision \
  --add-label go-idioms \
  --description "
Type: decision

Context: Implementing Linear GraphQL client configuration system. Need a way
to set optional parameters like cache TTL, timeout, retry policy...

Alternatives Considered:
1. Builder pattern: client.New().WithCache(ttl).Build()
   - Pros: Familiar from Java/C++, chainable
   - Cons: Not idiomatic Go 1.21+, mutable during construction

2. Config struct: client.New(Config{Cache: ttl, Timeout: t})
   - Pros: Simple, single function call
   - Cons: Can't distinguish zero-value from unset

3. [CHOSEN] Functional options: client.New(WithCache(ttl), WithTimeout(t))
   - Pros: Idiomatic Go, self-documenting, extensible
   - Cons: Slightly more complex implementation

Reasoning:
- Effective Go recommends functional options
- Standard library uses this pattern (grpc, net/http)
- Compile-time safety for required vs optional params

Trade-offs:
Gained: + Idiomatic Go + Better API ergonomics + Extensible
Lost: - More complex implementation (10 functions vs 1 struct)

Implementation:
- Pattern definition: pkg/client/options.go
- Client constructor: pkg/client/client.go:45-80

Implications:
- ALL future lirt packages should use functional options
- Establishes consistency pattern across codebase
"
```

### Processing Chronicle Beads

The **lirt-chronicler** agent processes chronicle beads:

```bash
# Invoke agent to process all open chronicle beads
lirt-chronicler
```

The agent will:
1. Query: `bd ready --type chronicle`
2. Analyze beads for grouping opportunities (4-hour window, shared themes)
3. Create diary entries in `diary/` with proper formatting
4. Update `diary/_index.md`
5. Close processed beads

### Diary Structure

```
diary/
├── _index.md                                    # Chronological index
└── {YY-MM-DD}.{HH-MM-TZ}.{town}.{topic}.md      # Individual entries
```

## Test Case Identification Protocol (Beads-Mediated)

**Protocol Pattern:** Detection (ALL agents) → Beads Queue → Execution (lirt-test-engineer)

### When to Create Test Beads

**ALL agents** must watch for testable scenarios:

**During Specification Writing (lirt-spec-writer):**
- Edge cases in command behavior
- Error conditions in CLI usage
- Complex flag combinations
- Format conversions (JSON, CSV, table)

**During Implementation (lirt-specialist):**
- Edge cases discovered in code
- Error handling paths
- Complex business logic branches
- Performance-critical code paths

**During Code Review (lirt-code-reviewer):**
- Untested code paths
- Missing error handling tests
- Security-sensitive operations without tests
- Performance optimizations needing benchmarks

### Creating Test Beads

```bash
bd create --title "Test: [specific scenario]" \
  --type task \
  --priority 2 \
  --add-label requires:testing \
  --description "[Rich 150-300 word description]"
```

**Priority:**
- **P0**: Critical path, blocks release
- **P1**: Important functionality, needed for stability
- **P2**: Standard tests (default)
- **P3**: Nice-to-have
- **P4**: Future consideration

### Test Bead Template

```
## Scenario
[What specific scenario needs testing]

## Context
[What were you working on when you identified this]

## Test Cases Needed
1. [Specific test case with expected behavior]
2. [Another test case]
3. [Edge case]

## Why This Matters
[Why this scenario is important to test]

## Implementation Notes
[Relevant code locations, patterns to use]

## Acceptance Criteria
[What makes this test complete]
```

### Processing Test Beads

The **lirt-test-engineer** processes test beads:

```bash
# Query test beads
bd ready --label requires:testing

# After implementing tests
bd close <bead-id>
```

## Code Quality Reporter Protocol (Beads-Mediated)

**Protocol Pattern:** Detection (lirt-code-reviewer) → Beads Queue → Execution (lirt-specialist)

**Note:** lirt-code-reviewer has read-only tools (Read, Grep, Glob) and cannot fix code directly.

### When to Create Quality Beads

Watch for these categories:

1. **Go Idiom Violations** (label: `idioms`)
   - Non-idiomatic error handling
   - Missing context propagation
   - Incorrect interface usage

2. **Performance Issues** (label: `performance`)
   - String concatenation in loops
   - Unnecessary allocations
   - Hot path optimizations needed

3. **Security Concerns** (label: `security`)
   - API keys in logs or errors
   - Insecure file permissions
   - Missing input validation

4. **Code Quality Issues** (label: `maintainability`)
   - Complex functions
   - Duplicated code
   - Magic numbers without constants

5. **Testing Gaps** (label: `testing`)
   - Untested code paths
   - Missing benchmarks
   - Missing integration tests

### Creating Quality Beads

```bash
bd create --title "[Category]: [Specific issue]" \
  --type bug \
  --priority [0-4] \
  --add-label code-quality \
  --add-label [specific-category] \
  --description "[Rich 200-400 word description]"
```

**Priority:**
- **P0**: Security vulnerability, data loss risk, crash
- **P1**: Performance regression, startup time impact >10ms
- **P2**: Go idiom violation, code smell (default)
- **P3**: Minor improvement
- **P4**: Nice-to-have refactoring

### Quality Bead Template

```
## Issue
[Specific code quality problem]

## Location
[File:line or files affected]

## Current State
[What the code does now / what's wrong]

## Why This Matters
[Impact on quality, performance, security, or maintainability]

## Suggested Fix
[How to address the issue with code example]

## References
[Links to Go best practices, benchmarks, or docs]

## Acceptance Criteria
[What makes this issue resolved]
```

### Processing Quality Beads

The **lirt-specialist** processes quality beads:

```bash
# Query quality beads
bd ready --type bug --label code-quality

# After fixing issue
bd close <bead-id>
```

---

## Pre-Push Chronicle Gate

**Before EVERY push** — not just at session end — you MUST ensure all chronicle beads are processed. This gate applies to every `git push` or `lirt-push` invocation.

### The Gate (Two Checks)

Before pushing, stop and verify:

1. **No Open Chronicle Beads:**
   ```bash
   bd ready --type chronicle
   ```
   Must return empty. If not empty:
   - Process beads with `lirt-chronicler` (creates diary entries, closes beads)
   - OR defer non-critical beads to P4 if time-constrained

2. **Review Recent Work for New Chronicle-Worthy Items:**
   Review all commits being pushed:
   - New capability, crew member, or infrastructure added?
   - Significant decision with reasoning worth preserving?
   - Insight about Go CLI development or Linear API integration?
   - Pattern or lesson learned?
   - GraphQL client design decision?

   If yes → Create chronicle bead IMMEDIATELY (while context fresh)

### Why This Exists

Chronicle items are easy to miss mid-session. The gate prevents this: **no push without ensuring chronicle beads are either processed or captured.**

The `lirt-push` script enforces this gate by checking `bd ready --type chronicle` returns empty.

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

4. **Process Chronicle Beads** - Ensure no open chronicle beads:
   ```bash
   bd ready --type chronicle
   ```
   If not empty:
   ```bash
   # Process beads (creates diary entries, closes beads)
   lirt-chronicler

   # Verify empty
   bd ready --type chronicle  # Must return empty
   ```

5. **Run the Pre-Push Chronicle Gate** (see above) - Review ALL work from this session, create new chronicle beads if needed

6. **PUSH TO REMOTE** - This is MANDATORY:
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

7. **Clean up** - Clear stashes, prune remote branches

8. **Verify** - All changes committed AND pushed

9. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
- NEVER push without processing all open chronicle beads first
- `bd ready --type chronicle` MUST return empty before push

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
