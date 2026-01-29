# Archivist Protocol

**Role**: Detection protocol for identifying archive-worthy moments
**Applies to**: ALL agents working on lirt

## Purpose

The Archivist maintains institutional memory by capturing research findings and strategic decisions. Unlike the Chronicler (which captures insights and patterns), the Archivist focuses on:

- **Research**: External sources, technical investigations, tool evaluations
- **Decisions**: Strategic choices with context, alternatives, and reasoning

This ensures knowledge persists across sessions and provides the "why" behind choices.

## Two Types of Archive Content

### Research
Investigations and findings from external sources:
- Web searches for project work
- External documentation referenced (APIs, specs, papers)
- Third-party tools or services evaluated
- Technical investigations with conclusions

### Decisions
Strategic choices that shape project direction:
- Architecture or design choices
- Tool/technology selections
- Process or workflow decisions
- Scope changes or priority shifts

## When to Create Archive Beads

**ALL agents** must watch for archive-worthy moments:

### Research Triggers

| Trigger | Example |
|---------|---------|
| Web search conducted | Researched authentication libraries for Go |
| External docs referenced | Read Linear API GraphQL documentation |
| Tool/service evaluated | Compared Redis vs Memcached for caching |
| Technical investigation completed | Tested cursor vs offset pagination |

### Decision Triggers

Watch for language like:
| Pattern | Example |
|---------|---------|
| **Approval/rejection** | "approved", "rejected", "let's do X instead" |
| **Direction setting** | "the approach is", "we'll use", "go with" |
| **Course correction** | "don't do X, do Y", "change direction to" |
| **Scope changes** | Priority shifts, feature additions/removals |
| **Architecture choices** | "we'll structure it as", "the pattern will be" |

## What NOT to Archive

- Routine task decisions (that's for beads/issues)
- Implementation details without strategic significance
- Temporary debugging notes
- Insights and patterns (that's for Chronicler)

## Creating Archive Beads

When you identify something archive-worthy, create an archive bead **IMMEDIATELY**:

### For Research

```bash
bd create --title "Archive: Research on [topic]" \
  --type archive \
  --priority 3 \
  --add-label research \
  --add-label [topic-area] \
  --description "[See research template below]"
```

### For Decisions

```bash
bd create --title "Archive: DEC-NNN [decision title]" \
  --type archive \
  --priority 2 \
  --add-label decision \
  --add-label [area] \
  --description "[See decision template below]"
```

### Topic Labels

Add labels to categorize:

```bash
# Content type (required - pick one)
--add-label research      # External investigation
--add-label decision      # Strategic choice

# Topic area (pick relevant ones)
--add-label api           # API design/integration
--add-label architecture  # System architecture
--add-label tooling       # Tool selection
--add-label process       # Workflow/process
--add-label security      # Security decisions
--add-label performance   # Performance choices
```

### Priority Guidelines

- **P1**: Major architectural decision, critical research
- **P2**: Standard decisions (default for decisions)
- **P3**: Standard research (default for research)
- **P4**: Minor findings, can be batched

## Research Bead Template

```
Type: research
Topic: [specific topic being researched]
Status: COMPLETED | IN_PROGRESS

## Purpose
[Why this research was needed - what question are we answering]

## Sources Consulted
1. [Source name] - [URL]
   Key findings: [what we learned]

2. [Source name] - [URL]
   Key findings: [what we learned]

3. [Additional sources...]

## Summary of Findings
[Synthesized conclusions from all sources - what did we learn overall]

## Recommendations
[Based on findings, what should we do]

## Application
[How this research will be/was applied to the project]

## Related
- Decisions informed: [DEC-NNN if any]
- Related research: [other topics if any]

## Assets
[List any PDFs, screenshots, or artifacts to preserve]
```

## Decision Bead Template

```
Type: decision
ID: DEC-NNN-[slug]
Status: Proposed | Accepted | Superseded

## Context
[What problem or question prompted this decision]
[What constraints or requirements apply]

## Research Basis
[Link to research topics that informed this, if any]

## The Decision
[Clear statement of what was decided]

## Alternatives Considered

### Alternative A: [name]
- Description: [what this option entails]
- Pros: [benefits]
- Cons: [drawbacks]
- Why rejected: [reasoning]

### Alternative B: [name]
- Description: [what this option entails]
- Pros: [benefits]
- Cons: [drawbacks]
- Why rejected: [reasoning]

### Chosen: [name]
- Description: [what this option entails]
- Pros: [benefits]
- Cons: [drawbacks]
- Why chosen: [reasoning]

## Consequences
[Expected impact and implications of this decision]
[What changes as a result]

## Implementation Notes
[How this decision will be/was implemented]
[File:line references if applicable]

## Review Date
[When should this decision be revisited, if applicable]
```

## Minimum Description Length

- **Research beads**: 200-400 words minimum
- **Decision beads**: 300-500 words minimum

This ensures the archivist can create comprehensive archive entries.

## Timing: Create Immediately

**DON'T WAIT** until end of session:
- Research sources get forgotten
- Decision context fades quickly
- Alternatives considered are lost

**CREATE RIGHT NOW** when:
- Research is completed (or paused)
- Decision is made by the Overseer
- Direction changes based on new information

## Example Research Bead

```bash
bd create --title "Archive: Research on Go GraphQL clients" \
  --type archive \
  --priority 3 \
  --add-label research \
  --add-label api \
  --description "
Type: research
Topic: Go GraphQL client libraries
Status: COMPLETED

## Purpose
Need to select a GraphQL client for Linear API integration. Requirements:
query generation, type safety, good error handling.

## Sources Consulted
1. machinebox/graphql - https://github.com/machinebox/graphql
   Key findings: Simple, minimal, no code generation. Good for simple queries.

2. shurcooL/graphql - https://github.com/shurcooL/graphql
   Key findings: Struct-based queries, type-safe, maintained. Medium complexity.

3. Khan/genqlient - https://github.com/Khan/genqlient
   Key findings: Code generation from schema, full type safety, active development.

4. hasura/go-graphql-client - https://github.com/hasura/go-graphql-client
   Key findings: Fork of shurcooL with subscriptions support.

## Summary of Findings
For our use case (Linear API integration with complex queries), Khan/genqlient
offers the best balance of type safety and developer experience. Code generation
means compile-time query validation and IDE autocompletion.

## Recommendations
Use Khan/genqlient with Linear's GraphQL schema. Generate types from schema
and use strongly-typed query functions.

## Application
Will be used for all Linear API calls in the lirt CLI tool.

## Related
- Decisions informed: DEC-003-graphql-client-selection
"
```

## Example Decision Bead

```bash
bd create --title "Archive: DEC-003 GraphQL client selection" \
  --type archive \
  --priority 2 \
  --add-label decision \
  --add-label tooling \
  --description "
Type: decision
ID: DEC-003-graphql-client-selection
Status: Accepted

## Context
Need to select a Go GraphQL client library for Linear API integration.
Must support complex queries, provide type safety, and have active maintenance.

## Research Basis
See: Research on Go GraphQL clients

## The Decision
Use Khan/genqlient for all GraphQL operations.

## Alternatives Considered

### Alternative A: machinebox/graphql
- Description: Minimal GraphQL client with simple string queries
- Pros: Simple API, small dependency, easy to understand
- Cons: No type safety, queries are strings, no code generation
- Why rejected: Too error-prone for complex Linear queries

### Alternative B: shurcooL/graphql
- Description: Struct-based query construction
- Pros: Type-safe, well-tested, moderate complexity
- Cons: Verbose struct definitions, manual query construction
- Why rejected: Too verbose for our query volume

### Chosen: Khan/genqlient
- Description: Code generation from GraphQL schema
- Pros: Full type safety, IDE support, compile-time validation
- Cons: Requires schema, code generation step
- Why chosen: Best developer experience, catches errors at compile time

## Consequences
- Must download Linear GraphQL schema during build
- Generated code in gen/ directory
- All queries type-checked at compile time
- Easier refactoring when Linear API changes

## Implementation Notes
- Schema: tools/download-schema.sh
- Config: genqlient.yaml
- Generated: pkg/linear/gen/

## Review Date
Review when Linear API has breaking changes or genqlient has major update
"
```

## Decision ID Assignment

Decision IDs follow the format: `DEC-{NNN}-{slug}`

To assign a new ID:
1. Check `archive/decisions/_index.md` for latest number
2. Increment by 1
3. Add descriptive slug (kebab-case)

Example: If latest is DEC-002, next is DEC-003.

## Self-Check Before Creating

Before creating an archive bead, ask yourself:
- Is this research or a decision (not an insight/pattern)?
- Have I included all sources consulted (for research)?
- Have I documented alternatives considered (for decisions)?
- Is the reasoning clear enough for future reference?
- Is this at least 200-500 words with structured content?

If any answer is NO → Add more detail before creating the bead

## Processing Archive Beads

The lirt-archivist agent will:
1. Query: `bd ready --type archive`
2. Read all open archive beads
3. Determine type (research or decision) from labels
4. Create appropriate archive entry structure
5. Update relevant index file
6. Close processed beads: `bd close <id1> <id2> ...`

## Archive Structure

```
{{answer.archive_path}}
├── decisions/
│   ├── _index.md                    # Chronological decision index
│   └── DEC-NNN-slug/
│       ├── SUMMARY.md               # Decision document
│       └── assets/                  # Supporting artifacts
└── research/
    ├── _index.md                    # Topic index with status
    └── topic-slug/
        ├── SUMMARY.md               # Research summary
        └── assets/                  # PDFs, images, etc.
```

## Index Formats

### decisions/_index.md

```markdown
# Decisions Index

| ID | Date | Title | Status |
|----|------|-------|--------|
| [DEC-003](DEC-003-graphql-client-selection/SUMMARY.md) | 2026-01-28 | GraphQL client selection | Accepted |
| [DEC-002](DEC-002-sync-branch/SUMMARY.md) | 2026-01-26 | Sync branch for beads | Accepted |
| [DEC-001](DEC-001-architecture/SUMMARY.md) | 2026-01-20 | Initial architecture | Accepted |
```

### research/_index.md

```markdown
# Research Index

| Topic | Status | Updated | Summary |
|-------|--------|---------|---------|
| [GraphQL clients](graphql-clients/SUMMARY.md) | COMPLETED | 2026-01-28 | Evaluated Go GraphQL libraries |
| [Linear API](linear-api/SUMMARY.md) | IN_PROGRESS | 2026-01-25 | API capabilities and patterns |
```

## Success Metrics

Track protocol effectiveness:
- Research entries with complete source documentation
- Decision entries with alternatives considered
- Time from investigation/decision to archive entry
- Usefulness of entries in future sessions
