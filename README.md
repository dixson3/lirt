# lirt

**Command-line tool for interacting with Linear from agents and developers.**

lirt is a Go-based CLI tool for Linear following `gh` CLI semantics, designed for both:
- **Claude Code + Gas Town infrastructure** (agent-driven development)
- **Human developers** (traditional CLI usage)

## Quick Start

### For Claude Code Agents

```bash
# Get full context
gt prime

# Check for work
bd ready

# Process chronicle beads
bd ready --type chronicle
```

See [AGENTS.md](./AGENTS.md) for agent protocols and workflows.

### For Human Developers

```bash
# Install (future)
go install github.com/dixson3/lirt@latest

# Configure
lirt auth login

# Use
lirt issue list
lirt issue view TEAM-123
```

See [Usage](#usage) section below for CLI reference.

---

## Project Layout

This project uses **Gas Town infrastructure** for multi-agent coordination:

```
lirt/
├── .beads/                      # Issue tracking & coordination
│   ├── config.yaml             # Beads configuration (committed)
│   └── issues.jsonl            # Issue database (on beads-sync branch)
│
├── refinery/rig/               # Main development workspace
│   ├── .claude/                # Claude Code agent definitions
│   │   ├── agents/            # Specialized agents (5 crew members)
│   │   └── CREW-MEMBERS.md    # Agent documentation
│   │
│   ├── roles/                  # Agent protocols (cross-cutting)
│   │   ├── chronicler.md      # Chronicle protocol
│   │   ├── test-case-identifier.md
│   │   └── code-quality-reporter.md
│   │
│   ├── diary/                  # Development insights chronicle
│   │   ├── _index.md          # Diary index
│   │   └── *.md               # Individual diary entries
│   │
│   ├── bin/                    # Gas Town utilities
│   │   ├── lirt-push          # Pre-push chronicle gate
│   │   ├── lirt-setup-hooks   # Install hooks
│   │   └── lirt-town-id       # Town identifier
│   │
│   ├── AGENTS.md               # Agent coordination documentation
│   ├── CLAUDE.md               # Context recovery reference
│   └── README.md               # This file
│
├── witness/                    # Monitors polecat health
├── polecats/                   # Ephemeral worker instances
└── crew/                       # Isolated workspace (gitignored)
```

### Key Directories

- **`.beads/`**: Persistent issue tracking across sessions and agents
- **`refinery/rig/`**: Primary development workspace (sparse checkout of main repo)
- **`diary/`**: Chronicle of development decisions, insights, and patterns
- **`.claude/agents/`**: Specialized agent definitions (lirt-specialist, lirt-test-engineer, etc.)
- **`roles/`**: Cross-cutting protocols that all agents follow

---

## Usage Patterns

### Claude Code + Gas Town (Agent-Driven)

This project is instrumented for multi-agent development workflows:

#### 1. Chronicle Protocol (Beads-Mediated)

Agents capture development insights via beads:

```bash
# ANY agent detects chronicle-worthy moment (decision, insight, pattern)
bd create --title "Chronicle: [topic]" \
  --type chronicle \
  --priority 3 \
  --labels decision,go-idioms \
  --description "[Rich 300-600 word context]"

# lirt-chronicler processes queue into diary entries
bd ready --type chronicle       # Check for open beads
# (lirt-chronicler agent creates diary entries)
bd ready --type chronicle       # Verify empty before push
```

See [roles/chronicler.md](./roles/chronicler.md) for full protocol.

#### 2. Multi-Agent Crew

5 specialized agents coordinate on lirt development:

- **lirt-specialist**: Lead developer (Go, GraphQL, CLI)
- **lirt-spec-writer**: Documentation engineer (specs, guides)
- **lirt-test-engineer**: Test automation (table-driven tests, benchmarks)
- **lirt-code-reviewer**: Code quality specialist (Go idioms, security)
- **lirt-chronicler**: Diary entry processor (creates diary from beads)

See [.claude/CREW-MEMBERS.md](./.claude/CREW-MEMBERS.md) for agent details.

#### 3. Beads Coordination

Persistent issue tracking survives sessions and compaction:

```bash
# Sync beads from remote
bd sync

# Find available work
bd ready

# Claim work
bd update <id> --status in_progress

# Complete work
bd close <id>

# Sync changes back
bd sync
```

Beads are stored on the `beads-sync` branch and synced via `bd sync`.

#### 4. Landing the Plane (Session Close)

Before ending a work session:

```bash
# 1. Quality gates
go test ./...
golangci-lint run

# 2. Update beads
bd close <finished-ids>
bd sync

# 3. Process chronicle beads (MANDATORY)
bd ready --type chronicle       # Must return empty
# If not empty: process with lirt-chronicler

# 4. Git workflow
git add <files>
git commit -m "..."
bd sync
git push                        # Or: bin/lirt-push (enforces chronicle gate)
```

See [AGENTS.md](./AGENTS.md) for complete Landing the Plane protocol.

---

### Human Developers (Traditional CLI)

If you're a human developer working on lirt:

#### Setup

```bash
# Clone repository
git clone git@github.com:dixson3/lirt.git
cd lirt

# Standard Go development (ignore Gas Town infrastructure)
go mod download
go build ./cmd/lirt
```

#### Development

```bash
# Run tests
go test ./...

# Lint
golangci-lint run

# Build
go build -o lirt ./cmd/lirt

# Install locally
go install ./cmd/lirt
```

#### Contributing

1. Create a feature branch: `git checkout -b feature/your-feature`
2. Make changes following Go best practices
3. Write tests (table-driven, use golden files for CLI output)
4. Run quality gates: `go test ./... && golangci-lint run`
5. Commit changes: `git commit -m "Add feature: description"`
6. Push and create PR: `git push origin feature/your-feature`

**Note**: You can safely ignore the Gas Town infrastructure (`.beads/`, `diary/`, `.claude/agents/`, etc.) - these are for agent-driven development workflows.

---

## Chronicle Protocol

The **chronicle** captures the evolution of the Overseer's thinking during lirt development. This is not a changelog—it's a record of *how* work is performed and *why* decisions are made.

### Chronicle-Worthy Moments

- **Decisions**: Choosing functional options over builder pattern
- **Insights**: Realizing cursor pagination prevents data loss
- **Patterns**: Recognizing table-driven tests as recurring theme
- **Corrections**: Pivoting GraphQL client design after performance testing
- **Lessons**: What worked or didn't in Go testing strategies

### Chronicle Workflow (Agents)

1. **Detect**: Any agent identifies chronicle-worthy moment
2. **Capture**: Create chronicle bead with 300-600 word description
3. **Process**: lirt-chronicler transforms beads into diary entries
4. **Enforce**: Pre-push gate blocks if chronicle beads remain open

### Diary Structure

```
diary/
├── _index.md                                          # Chronological index
└── {YY-MM-DD}.{HH-MM-TZ}.{town}.{topic}.md           # Entries
```

Example: `26-01-27.17-07-PST.gt.beads-mediated-coordination.md`

See [diary/_index.md](./diary/_index.md) for all entries.

---

## Beads System

**Beads** is the persistent issue tracking system used across Gas Town projects.

### What is Beads?

- SQLite database + JSONL sync
- Survives sessions, context compaction, agent changes
- Multi-user/multi-town coordination
- Git-backed via `beads-sync` branch

### Custom Types

lirt uses a custom bead type:

- **chronicle**: Chronicle-worthy moments (decisions, insights, patterns)

Standard types (bug, feature, task, epic) also available.

### Beads Commands

```bash
# List issues
bd list                         # Default view
bd list --all                   # Include closed
bd ready                        # Ready work (no blockers)
bd ready --type chronicle       # Chronicle beads only

# Create issue
bd create --title "Title" \
  --type task \
  --priority 2 \
  --description "..."

# Update issue
bd update <id> --status in_progress
bd update <id> --priority 1

# Close issue
bd close <id>

# Show details
bd show <id>

# Sync with git
bd sync                         # Export to JSONL, commit to beads-sync branch
```

### Multi-User Workflow

When collaborating across towns:

```bash
# Pull latest beads
git fetch origin
bd sync                         # Auto-imports from beads-sync branch

# Create/update beads
bd create ...
bd update ...

# Push beads
bd sync                         # Exports to JSONL
git push                        # Pushes beads-sync branch
```

---

## Gas Town Infrastructure

### What is Gas Town?

Gas Town is a coordination infrastructure for multi-agent development with persistent state across sessions.

**Key concepts:**
- **Rig**: Development workspace (lirt is a rig in the GT town)
- **Crew**: Specialized agents working on the rig
- **Witness**: Monitors agent health and progress
- **Refinery**: Processes merge queue with verification gates
- **Beads**: Persistent issue tracking across agents/sessions
- **Chronicle**: Development insights diary

### When to Use Gas Town Features

**Use Gas Town features if:**
- You're a Claude Code agent working on lirt
- You want to coordinate across multiple sessions
- You're capturing development insights (chronicle protocol)
- You're working in a multi-agent workflow

**Skip Gas Town features if:**
- You're a human developer using traditional Git workflow
- You're making a one-off contribution
- You prefer standard GitHub issues/PRs

---

## Project Goals

### Technical Requirements

- **Go 1.21+**: Modern Go idioms and patterns
- **Startup time**: < 50ms (optimized for bash scripting workflows)
- **Memory usage**: < 50MB baseline
- **Linear API**: GraphQL client with cursor pagination
- **CLI semantics**: Follows `gh` CLI patterns

### Design Principles

1. **Idiomatic Go**: Functional options, context propagation, table-driven tests
2. **AWS CLI-style profiles**: Multiple workspace support
3. **Bash scripting integration**: Fast, composable commands
4. **Chronicle-first**: Capture the "why" behind decisions
5. **Agent-friendly**: Instrumented for multi-agent coordination

---

## Testing

### For Agents

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Update golden files
go test ./... -update
```

See [.claude/agents/lirt-test-engineer.md](./.claude/agents/lirt-test-engineer.md) for testing patterns.

### For Human Developers

Standard Go testing:

```bash
go test -v ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Documentation

- **[AGENTS.md](./AGENTS.md)**: Agent protocols and coordination
- **[.claude/CREW-MEMBERS.md](./.claude/CREW-MEMBERS.md)**: Agent definitions
- **[diary/_index.md](./diary/_index.md)**: Development insights chronicle
- **[roles/chronicler.md](./roles/chronicler.md)**: Chronicle protocol
- **[roles/test-case-identifier.md](./roles/test-case-identifier.md)**: Test identification protocol
- **[roles/code-quality-reporter.md](./roles/code-quality-reporter.md)**: Code quality protocol

**Town-level documentation:**
- **Gas Town beads-mediated coordination pattern**: `../../../mayor/docs/BEADS-MEDIATED-COORDINATION.md`

---

## Contributing

### As a Claude Code Agent

Follow the protocols in [AGENTS.md](./AGENTS.md):
1. Check for work: `bd ready`
2. Follow chronicle protocol for insights
3. Create chronicle/test/quality beads as appropriate
4. Land the plane before session end

### As a Human Developer

Standard GitHub workflow:
1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Submit PR to `main` branch

**Commit guidelines:**
- Use conventional commits: `feat:`, `fix:`, `docs:`, `test:`, etc.
- Keep commits focused and atomic
- Include tests for new features

---

## License

MIT License - see [LICENSE](./LICENSE)

---

## Status

**Current Phase**: Early development
- Gas Town infrastructure established
- Chronicle protocol operational
- Beads coordination functional
- CLI implementation pending

**Next Steps**:
- Implement Linear GraphQL client
- Add core CLI commands (issue list, view, create)
- Write comprehensive tests
- Add bash completion support
