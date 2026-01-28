# Chronicler Protocol

**Role**: Detection protocol for identifying chronicle-worthy moments
**Applies to**: ALL agents working on lirt

## Purpose

During lirt development, insights emerge about Go CLI patterns, Linear API integration, design decisions, and lessons learned. This protocol ensures those insights are preserved for future sessions through diary entries.

## When to Create Chronicle Beads

**ALL agents working on lirt** must watch for chronicle-worthy moments:

### Decision Points
When choosing between approaches, record the reasoning:
- Go architecture patterns (client design, error handling, testing strategies)
- Linear GraphQL API integration patterns
- CLI UX decisions (command structure, flag naming, output formats)
- Performance optimization trade-offs

### Insights
Realizations that change understanding or approach:
- Go idioms discovered or applied
- Linear API quirks or limitations
- CLI workflow patterns that emerge
- Bash scripting integration strategies

### Pattern Recognition
Recurring themes in how problems are solved:
- Common Go CLI patterns
- GraphQL query optimization patterns
- Testing and mocking strategies

### Course Corrections
When and why direction changed:
- API design pivots
- Performance optimization pivots
- UX simplification decisions

### Lessons Learned
What worked, what didn't, and why:
- Go best practices validated or refuted
- Linear API integration gotchas
- CLI design lessons

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

Add labels to help lirt-chronicler organize entries:

```bash
--add-label decision       # Architectural or design decision
--add-label insight        # Realization that changed understanding
--add-label pattern        # Recurring theme recognized
--add-label correction     # Direction change
--add-label lesson         # What worked/didn't work

--add-label go-idioms      # Go language patterns
--add-label graphql        # GraphQL/Linear API
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

This ensures lirt-chronicler can write comprehensive diary entries without losing the "why" behind decisions.

## Timing: Create Immediately

❌ **DON'T WAIT** until end of session:
- Details fade quickly from working memory
- Context becomes less rich
- You'll forget nuances

✅ **CREATE RIGHT NOW** when insight occurs:
- Full context is in your working memory
- Can reference code you just wrote
- Reasoning is fresh

## Example Chronicle Bead (Decision)

```bash
bd create --title "Chronicle: Functional options for client configuration" \
  --type chronicle \
  --priority 3 \
  --add-label decision \
  --add-label go-idioms \
  --description "
Type: decision

Context: Implementing Linear GraphQL client configuration system. Need a way
to set optional parameters like cache TTL, timeout, retry policy, and custom
headers without making the API cumbersome.

Alternatives Considered:
1. Builder pattern: client.New().WithCache(ttl).WithTimeout(t).Build()
   - Pros: Familiar from Java/C++, chainable, explicit build step
   - Cons: Not idiomatic Go 1.21+, mutable during construction, verbose
   - Example in Go: uncommon, mostly seen in legacy code

2. Config struct: client.New(Config{Cache: ttl, Timeout: t})
   - Pros: Simple, single function call, struct is self-documenting
   - Cons: Can't distinguish zero-value from unset, poor ergonomics for optional fields
   - Example: cfg := Config{Timeout: 10*time.Second} // Cache is zero-value, means no cache or default?

3. [CHOSEN] Functional options: client.New(WithCache(ttl), WithTimeout(t))
   - Pros: Idiomatic Go 1.21+, self-documenting, extensible without breaking changes
   - Cons: Slightly more complex implementation (need option function type and appliers)
   - Example: Used by grpc, net/http server options, widely accepted pattern

Reasoning:
- Effective Go recommends functional options for optional configuration
- Standard library uses this pattern (context.WithTimeout, grpc.WithInsecure)
- Provides compile-time safety: required params are function args, optional params are options
- Self-documenting: WithCache(5*time.Minute) is clearer than Cache: 300
- Extensible: can add new options without breaking existing callers
- Type-safe: option functions can validate values at construction time

Trade-offs:
Gained:
+ Idiomatic Go code following community standards
+ Better API ergonomics for users
+ Extensibility without breaking changes
+ Compile-time safety for required vs optional

Lost:
- More complex implementation (10 option functions vs 1 config struct)
- Need to understand closure pattern for new contributors
- Slightly more lines of code

Implementation:
- Pattern definition: pkg/client/options.go
  ```go
  type Option func(*Client) error

  func WithCache(ttl time.Duration) Option {
      return func(c *Client) error {
          if ttl < 0 {
              return fmt.Errorf(\"cache TTL must be positive\")
          }
          c.cache = NewCache(ttl)
          return nil
      }
  }
  ```

- Client constructor: pkg/client/client.go:45-80
  ```go
  func New(apiKey string, opts ...Option) (*Client, error) {
      c := &Client{apiKey: apiKey}
      for _, opt := range opts {
          if err := opt(c); err != nil {
              return nil, err
          }
      }
      return c, nil
  }
  ```

- Usage example: cmd/issue/list.go:25
  ```go
  client, err := client.New(
      apiKey,
      client.WithCache(5*time.Minute),
      client.WithTimeout(30*time.Second),
  )
  ```

Implications:
- ALL future lirt packages requiring configuration should use functional options
- Establishes consistency pattern across codebase
- Need to document this pattern in CONTRIBUTING.md for new contributors
- Sets expectation that required parameters are function arguments
- Optional parameters always use With* option functions
- This is now the lirt way of doing configuration
"
```

## Example Chronicle Bead (Insight)

```bash
bd create --title "Chronicle: GraphQL cursor pagination advantages" \
  --type chronicle \
  --priority 3 \
  --add-label insight \
  --add-label graphql \
  --description "
Type: insight

Context: Implementing 'lirt issue list' pagination. Linear API supports both
offset-based and cursor-based pagination. Initially thought offset was simpler
and more user-friendly (page numbers are intuitive).

What Changed:
Before: Assumed offset pagination was best because it's familiar and simple.
Users understand \"page 1, page 2\" better than opaque cursor tokens.

After: Realized cursor pagination is critical for correctness when dataset
changes during pagination. Offset has a fundamental flaw that causes data
loss or duplication.

Reasoning:
Discovered while reading Linear API docs more carefully. Found this warning:

\"Offset pagination can miss or duplicate items if the dataset changes between
requests. For production use, we recommend cursor-based pagination.\"

Tested both approaches:
1. Offset pagination with concurrent issue creation:
   - Created 100 issues
   - Started paginating with offset (50 per page)
   - Added 10 new issues during pagination
   - Result: Missed 10 issues from page 2 (they shifted to page 1)

2. Cursor pagination with concurrent issue creation:
   - Same scenario
   - Result: No missed issues, cursor tracks position correctly

The insight: Offset pagination has a race condition when the dataset changes.
Since Linear is collaborative (multiple users modifying issues), this will
happen frequently in real use.

Trade-offs:
Gained:
+ Correctness: Won't miss or duplicate items
+ Linear API recommendation (best practice)
+ Better for production use

Lost:
- User-facing simplicity (can't jump to arbitrary page)
- Need to explain cursor concept in docs
- Slightly more complex state management (track cursor token)

Implementation:
- Issue list command: cmd/issue/list.go:125-180
- Cursor stored in result metadata for next page
- Use 50-item page size (balance between API efficiency and memory)

Implications:
- ALL list commands in lirt should use cursor pagination (projects, teams, comments, labels)
- Can't support \"jump to page N\" feature users might expect
- Need to document why we chose cursor over offset in CLI reference
- Consider adding --offset flag for interactive exploration (with warning)
- This is a CLI design decision: correctness > user familiarity
"
```

## Example Chronicle Bead (Pattern)

```bash
bd create --title "Chronicle: Table-driven test pattern for CLI commands" \
  --type chronicle \
  --priority 3 \
  --add-label pattern \
  --add-label testing \
  --description "
Type: pattern

Context: Writing tests for multiple CLI commands (issue list, issue view,
issue create). Noticed I was repeating the same test structure across all
command tests.

Pattern Recognized:
All CLI command tests follow the same structure:
1. Setup mock Linear API server
2. Run CLI command with specific args
3. Capture stdout/stderr
4. Compare against golden file
5. Verify exit code

This repeated 10-15 times per command with minor variations.

Reasoning:
Recognized this is the table-driven test pattern recommended by Go community.
Each test case is a row in a table with:
- Input (command args)
- Setup (mock API responses)
- Expected output (golden file)
- Expected error state

The pattern eliminates duplication and makes adding new test cases trivial.

Implementation:
Standard structure now used across all CLI tests:

```go
func TestIssueList(t *testing.T) {
    tests := []struct {
        name       string
        args       []string
        mockSetup  func(*httptest.Server)
        wantGolden string
        wantErr    bool
        wantErrMsg string
    }{
        {
            name: \"default table format\",
            args: []string{\"issue\", \"list\"},
            mockSetup: func(srv *httptest.Server) {
                srv.Handler = mockIssueListHandler(50)
            },
            wantGolden: \"testdata/issue-list-table.golden\",
            wantErr:    false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            srv := httptest.NewServer(tt.mockSetup(srv))
            defer srv.Close()

            // Run command, compare output
            output := runCLI(t, tt.args, srv.URL)
            golden := readGolden(t, tt.wantGolden)
            assert.Equal(t, golden, output)
        })
    }
}
```

Benefits discovered:
- Adding new test case is 5-10 lines vs 50-100 lines
- Test intent is clear from table structure
- Easy to see test coverage at a glance
- Subtests provide granular failure reporting
- Can run individual test cases: go test -run TestIssueList/default

Trade-offs:
+ Much less duplication
+ Easier to maintain
+ Clear test intent
- Slightly more complex setup (table structure)
- Need helper functions for common operations

Implications:
- ALL CLI command tests should use this pattern
- Document pattern in testing guidelines
- Create helper functions for common operations:
  * mockIssueListHandler(count int) http.Handler
  * runCLI(t *testing.T, args []string, apiURL string) string
  * readGolden(t *testing.T, path string) string
- This is now the lirt testing standard
"
```

## Self-Check Before Creating

Before creating a chronicle bead, ask yourself:
- ✅ Could lirt-chronicler write a comprehensive diary entry from this description alone?
- ✅ Have I captured the alternatives, reasoning, and trade-offs?
- ✅ Did I include file:line references or code snippets?
- ✅ Did I explain implications for future work?
- ✅ Is this at least 300 words with rich context?

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

lirt-chronicler may group chronicle beads into a single diary entry IF:
- Beads are from the same work session (created within hours)
- Beads share a common theme or feature area
- Grouping does NOT lose important context
- Grouping does NOT lose temporal ordering when it matters
- The combined entry is more coherent than separate entries

Example of GOOD grouping:
- "Chronicle: Functional options pattern" + "Chronicle: Client configuration design" + "Chronicle: API key handling"
  → Single entry: "Authentication system design decisions"

Example of BAD grouping (don't group):
- "Chronicle: GraphQL pagination (morning)" + "Chronicle: CLI error messages (afternoon)"
  → Different topics, keep separate

## Success Metrics

Track protocol effectiveness:
- Diary entry consistency and richness
- Number of insights captured per week
- Time from insight to diary entry
- Usefulness of entries in future sessions
