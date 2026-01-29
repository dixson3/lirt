# Chronicler Protocol

**Role**: Detection protocol for identifying chronicle-worthy moments
**Applies to**: ALL agents working on lirt

## Purpose

During development, insights emerge about architecture, API design, patterns, and lessons learned. This protocol ensures those insights are preserved for future sessions through diary entries.

## When to Create Chronicle Beads

**ALL agents** must watch for chronicle-worthy moments:

### Decision Points
When choosing between approaches, record the reasoning:
- Architecture patterns (client design, error handling, testing strategies)
- API integration patterns
- CLI/UX decisions (command structure, flag naming, output formats)
- Performance optimization trade-offs

### Insights
Realizations that change understanding or approach:
- Language idioms discovered or applied
- API quirks or limitations
- Workflow patterns that emerge
- Integration strategies

### Pattern Recognition
Recurring themes in how problems are solved:
- Common patterns across the codebase
- Query/request optimization patterns
- Testing and mocking strategies

### Course Corrections
When and why direction changed:
- API design pivots
- Performance optimization pivots
- UX simplification decisions

### Lessons Learned
What worked, what didn't, and why:
- Best practices validated or refuted
- Integration gotchas
- Design lessons

## Creating Chronicle Beads

When you identify something chronicle-worthy, create a chronicle bead **IMMEDIATELY** (while context is fresh):

```bash
bd create --title "Chronicle: [specific topic]" \
  --type chronicle \
  --priority 3 \
  --add-label [category] \
  --description "[See template below]"
```

### Category Labels

Add labels to help organize entries:

```bash
--add-label decision       # Architectural or design decision
--add-label insight        # Realization that changed understanding
--add-label pattern        # Recurring theme recognized
--add-label correction     # Direction change
--add-label lesson         # What worked/didn't work

--add-label architecture   # System architecture
--add-label api            # API design/integration
--add-label cli-ux         # CLI user experience
--add-label performance    # Performance optimization
--add-label testing        # Testing strategies
--add-label security       # Security decisions
```

### Priority Guidelines

Chronicle beads typically use **P3** (standard priority):
- **P2**: Major architectural decision with far-reaching implications
- **P3**: Standard chronicle-worthy moments (default)
- **P4**: Minor insights, can be batched

## Rich Description Template

Your chronicle bead description MUST capture these elements:

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

## Minimum Description Length

Chronicle bead descriptions should be **300-600 words minimum** to preserve rich context.

This ensures the chronicler can write comprehensive diary entries without losing the "why" behind decisions.

## Timing: Create Immediately

**DON'T WAIT** until end of session:
- Details fade quickly from working memory
- Context becomes less rich
- You'll forget nuances

**CREATE RIGHT NOW** when insight occurs:
- Full context is in your working memory
- Can reference code you just wrote
- Reasoning is fresh

## Example Chronicle Bead (Decision)

```bash
bd create --title "Chronicle: Functional options for client configuration" \
  --type chronicle \
  --priority 3 \
  --add-label decision \
  --add-label architecture \
  --description "
Type: decision

Context: Implementing client configuration system. Need a way to set optional
parameters like cache TTL, timeout, retry policy, and custom headers without
making the API cumbersome.

Alternatives Considered:
1. Builder pattern: client.New().WithCache(ttl).WithTimeout(t).Build()
   - Pros: Familiar from Java/C++, chainable, explicit build step
   - Cons: Not idiomatic, mutable during construction, verbose

2. Config struct: client.New(Config{Cache: ttl, Timeout: t})
   - Pros: Simple, single function call, struct is self-documenting
   - Cons: Can't distinguish zero-value from unset, poor ergonomics

3. [CHOSEN] Functional options: client.New(WithCache(ttl), WithTimeout(t))
   - Pros: Idiomatic, self-documenting, extensible without breaking changes
   - Cons: Slightly more complex implementation

Reasoning:
- Standard library uses this pattern (context, grpc, net/http)
- Provides compile-time safety: required params are args, optional are options
- Self-documenting: WithCache(5*time.Minute) is clearer than Cache: 300
- Extensible: can add new options without breaking existing callers

Trade-offs:
Gained:
+ Idiomatic code following community standards
+ Better API ergonomics for users
+ Extensibility without breaking changes

Lost:
- More complex implementation (10 option functions vs 1 config struct)
- Need to understand closure pattern for new contributors

Implementation:
- Pattern definition: pkg/client/options.go
- Client constructor: pkg/client/client.go:45-80

Implications:
- ALL future packages requiring configuration should use functional options
- Establishes consistency pattern across codebase
- This is now the project standard for configuration
"
```

## Example Chronicle Bead (Insight)

```bash
bd create --title "Chronicle: Cursor pagination advantages" \
  --type chronicle \
  --priority 3 \
  --add-label insight \
  --add-label api \
  --description "
Type: insight

Context: Implementing list pagination. API supports both offset-based and
cursor-based pagination. Initially thought offset was simpler and more
user-friendly (page numbers are intuitive).

What Changed:
Before: Assumed offset pagination was best because it's familiar and simple.
Users understand 'page 1, page 2' better than opaque cursor tokens.

After: Realized cursor pagination is critical for correctness when dataset
changes during pagination. Offset has a fundamental flaw that causes data
loss or duplication.

Reasoning:
Found warning in API docs: 'Offset pagination can miss or duplicate items
if the dataset changes between requests.'

Tested both approaches:
1. Offset pagination with concurrent creation:
   - Started paginating with offset (50 per page)
   - Added items during pagination
   - Result: Missed items from page 2 (they shifted to page 1)

2. Cursor pagination with concurrent creation:
   - Same scenario
   - Result: No missed items, cursor tracks position correctly

The insight: Offset pagination has a race condition when dataset changes.
Since the system is collaborative, this will happen frequently.

Trade-offs:
Gained:
+ Correctness: Won't miss or duplicate items
+ API best practice
+ Better for production use

Lost:
- User-facing simplicity (can't jump to arbitrary page)
- Need to explain cursor concept in docs
- Slightly more complex state management

Implementation:
- List command: cmd/list.go:125-180
- Cursor stored in result metadata for next page

Implications:
- ALL list commands should use cursor pagination
- Can't support 'jump to page N' feature
- This is a design decision: correctness > user familiarity
"
```

## Self-Check Before Creating

Before creating a chronicle bead, ask yourself:
- Could the chronicler write a comprehensive diary entry from this description alone?
- Have I captured the alternatives, reasoning, and trade-offs?
- Did I include file:line references or code snippets?
- Did I explain implications for future work?
- Is this at least 300 words with rich context?

If any answer is NO → Add more detail before creating the bead

## Processing Chronicle Beads

The lirt-chronicler agent will:
1. Query: `bd ready --type chronicle`
2. Read all open chronicle beads
3. Attempt to group related beads (same session, related topics)
4. Create diary entries (one per bead, or grouped if appropriate)
5. Update diary/_index.md
6. Close processed beads: `bd close <id1> <id2> ...`

## Grouping Related Beads

The chronicler may group chronicle beads into a single diary entry IF:
- Beads are from the same work session (created within hours)
- Beads share a common theme or feature area
- Grouping does NOT lose important context
- Grouping does NOT lose temporal ordering when it matters
- The combined entry is more coherent than separate entries

Example of GOOD grouping:
- "Chronicle: Functional options pattern" + "Chronicle: Client configuration design"
  → Single entry: "Client configuration design decisions"

Example of BAD grouping (don't group):
- "Chronicle: Pagination design (morning)" + "Chronicle: Error messages (afternoon)"
  → Different topics, keep separate

## Success Metrics

Track protocol effectiveness:
- Diary entry consistency and richness
- Number of insights captured per week
- Time from insight to diary entry
- Usefulness of entries in future sessions
