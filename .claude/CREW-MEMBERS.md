# lirt Crew Member Agents

This document describes the specialized crew member agents for the lirt project.

## Overview

lirt has 5 crew member agents: 4 adapted from [awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents) and customized for Go CLI development, plus a dedicated chronicler agent for capturing development insights.

## Crew Members

### 1. lirt-specialist
**File:** `.claude/agents/lirt-specialist.md`
**Role:** Lead Developer
**Expertise:** Go, GraphQL, CLI Design, Bash Integration

The primary development agent combining:
- Go 1.21+ idioms and best practices
- Linear GraphQL API client patterns
- AWS CLI-style profile management
- Bash scripting workflow integration
- Performance optimization (<50ms startup, <50MB memory)

**When to Use:**
- Core feature implementation
- GraphQL client development
- CLI command creation
- Configuration management
- Performance optimization

**Key Patterns:**
- Cobra for command structure
- Viper for configuration layering
- Context propagation for all API calls
- Functional options for client configuration
- Zero-allocation string building

---

### 2. lirt-spec-writer
**File:** `.claude/agents/lirt-spec-writer.md`
**Role:** Documentation Engineer
**Expertise:** Technical Writing, Specification, API Docs

Specialized in creating and maintaining:
- Technical specifications (requirements, architecture)
- CLI reference documentation
- User guides (authentication, configuration)
- API documentation (GraphQL patterns)
- Bash scripting integration guides

**When to Use:**
- Writing/updating specifications
- Creating command documentation
- Documenting API integration patterns
- Writing user guides
- Creating bash scripting examples

**Key Responsibilities:**
- Verify technical accuracy against Linear API
- Test all code examples
- Maintain documentation-as-code
- Create migration guides for breaking changes
- Ensure examples are shellcheck-clean

---

### 3. lirt-test-engineer
**File:** `.claude/agents/lirt-test-engineer.md`
**Role:** Test Automation Engineer
**Expertise:** Go Testing, Mocking, Benchmarking

Specialized in:
- Table-driven test patterns
- Golden file testing for CLI output
- GraphQL client mocking
- Integration testing
- Performance benchmarking

**When to Use:**
- Writing new tests
- Improving test coverage
- Creating mocks for Linear API
- Performance benchmarking
- CI/CD test integration

**Key Patterns:**
- `go test -update` for golden file regeneration
- httptest for mock Linear API servers
- sync.Pool for allocation testing
- Race detector verification
- Benchmark-driven optimization

---

### 4. lirt-code-reviewer
**File:** `.claude/agents/lirt-code-reviewer.md`
**Role:** Code Quality Specialist
**Expertise:** Go Idioms, Security, Performance Review

Specialized in reviewing:
- Go idioms and effective Go patterns
- CLI best practices
- GraphQL client patterns
- Performance optimization
- Security (API key handling, validation)

**When to Use:**
- Pull request reviews
- Code quality checks
- Security audits
- Performance reviews
- Ensuring consistency

**Review Checklist:**
- gofmt, golangci-lint, go vet compliance
- Context propagation
- Error handling patterns
- Startup time impact
- Memory allocation optimization
- Security best practices

---

### 5. lirt-chronicler
**File:** `.claude/agents/lirt-chronicler.md`
**Role:** Diary Entry Processor
**Expertise:** Processing Chronicle Beads, Diary Entry Creation

Processes chronicle beads into diary entries:
- Reads all open chronicle beads from queue
- Groups related beads when appropriate (4-hour window, shared themes)
- Creates formatted diary entries from bead descriptions
- Updates diary index
- Closes processed beads

**When to Use:**
- Before every push (Pre-Push Chronicle Gate)
- When ending a work session (Landing the Plane)
- Periodically to keep chronicle queue clean

**Invocation Pattern:**
```bash
# Run agent to process all open chronicle beads
lirt-chronicler
```

**Key Responsibilities:**
- Query open chronicle beads: `bd ready --type chronicle`
- Analyze beads for intelligent grouping
- Extract rich context from 300-600 word bead descriptions
- Create well-formatted diary entries
- Update diary/_index.md
- Close processed beads: `bd close <id1> <id2> ...`

**Grouping Logic:**
- Group IF: 4-hour window, shared themes, coherent narrative
- Don't group IF: Different topics, temporal ordering matters
- Default to separate entries when unclear

## Agent Coordination

### Development Workflow

```
┌──────────────────┐
│  Specification   │  lirt-spec-writer creates requirements
└────────┬─────────┘
         │
         v
┌──────────────────┐
│  Implementation  │  lirt-specialist builds features
└────────┬─────────┘
         │
         v
┌──────────────────┐
│    Testing       │  lirt-test-engineer writes tests
└────────┬─────────┘
         │
         v
┌──────────────────┐
│  Code Review     │  lirt-code-reviewer validates quality
└──────────────────┘

Cross-cutting: Any agent creates chronicle beads → lirt-chronicler processes them
```

### Example Usage

**New Feature Development:**
1. `lirt-spec-writer` - Document feature specification
2. `lirt-specialist` - Implement feature with Go code
3. `lirt-test-engineer` - Create comprehensive tests
4. `lirt-code-reviewer` - Review implementation and tests

**Bug Fix:**
1. `lirt-test-engineer` - Create failing test reproducing bug
2. `lirt-specialist` - Fix the bug
3. `lirt-code-reviewer` - Verify fix and test coverage

**Documentation Update:**
1. `lirt-spec-writer` - Update documentation
2. `lirt-test-engineer` - Verify examples work
3. `lirt-code-reviewer` - Check accuracy and completeness

## Invoking Agents

### Explicit Invocation
```
> Ask lirt-specialist to implement the issue list command
> Have lirt-test-engineer create tests for the auth module
> Get lirt-code-reviewer to review the GraphQL client code
> Ask lirt-spec-writer to document the cache strategy
```

### Automatic Invocation

Claude Code will automatically invoke agents based on context:
- Code implementation → lirt-specialist
- Test writing → lirt-test-engineer
- Documentation tasks → lirt-spec-writer
- Code review → lirt-code-reviewer

## Project Context

All agents are aware of:
- **Project:** lirt - Linear CLI tool in Go
- **Spec Location:** `/Users/james/gt/clio/mayor/rig/research/cl-k0a1/lirt-spec.md`
- **Target:** Go 1.21+, Linear GraphQL API
- **Goals:** <50ms startup, <50MB memory, bash scripting integration
- **Style:** AWS CLI patterns, gh CLI semantics

## Agent Customization

Each agent can be customized by editing the `.claude/agents/*.md` files:

```markdown
---
name: agent-name
description: When this agent should be invoked
tools: Read, Write, Edit, Bash, Glob, Grep
---

Agent instructions and patterns...
```

### Tool Permissions

Current tool assignments:
- **lirt-specialist:** Read, Write, Edit, Bash, Glob, Grep (full development)
- **lirt-spec-writer:** Read, Write, Edit, Glob, Grep, WebFetch, WebSearch (documentation + research)
- **lirt-test-engineer:** Read, Write, Edit, Bash, Glob, Grep (testing)
- **lirt-code-reviewer:** Read, Grep, Glob (read-only review)

## Success Metrics

### lirt-specialist
- Features implemented following Go idioms
- Startup time <50ms maintained
- Memory usage <50MB maintained
- All commands properly structured

### lirt-spec-writer
- 100% API coverage in docs
- All examples tested and working
- Zero broken cross-references
- Documentation passes accessibility checks

### lirt-test-engineer
- Test coverage >80%
- Fast execution (<30s full suite)
- Zero race conditions
- Benchmarks show acceptable performance

### lirt-code-reviewer
- All changes pass golangci-lint
- No performance regressions
- Security issues identified
- Consistent with codebase patterns

## Additional Resources

- **Source:** [awesome-claude-code-subagents](https://github.com/VoltAgent/awesome-claude-code-subagents)
- **lirt Spec:** `../clio/mayor/rig/research/cl-k0a1/lirt-spec.md`
- **Linear API:** https://developers.linear.app/docs/graphql/working-with-the-graphql-api
- **Effective Go:** https://go.dev/doc/effective_go
- **gh CLI:** https://cli.github.com/ (reference model)

---

Last Updated: 2026-01-27
