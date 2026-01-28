# Code Quality Reporter Protocol

**Role**: Detection protocol for identifying code quality issues
**Applies to**: lirt-code-reviewer (primarily), but all agents can report issues

## Purpose

The lirt-code-reviewer has **read-only tools** (Read, Grep, Glob) and cannot fix code. This protocol enables systematic reporting of quality issues via beads for lirt-specialist to address.

## When to Create Quality Beads

**lirt-code-reviewer** (and other agents during code review) must watch for:

### Go Idiom Violations
- Non-idiomatic error handling
- Missing context propagation
- Incorrect interface usage
- Ineffective Go patterns
- Exported items without documentation

### Performance Issues
- String concatenation in loops
- Unnecessary allocations
- Missing sync.Pool usage
- Inefficient data structures
- Hot path optimizations

### Security Concerns
- API keys in logs or errors
- Insecure file permissions
- Missing input validation
- Command/SQL injection risks
- Unhandled sensitive data

### Code Quality Issues
- Complex functions (high cyclomatic complexity)
- Duplicated code
- Magic numbers without constants
- Missing error handling
- Unclear variable naming

### Testing Gaps
- Untested code paths
- Missing race detector tests
- No benchmarks for hot paths
- Golden files not updated
- Integration tests missing

## Creating Quality Beads

When you identify a code quality issue, create a bead with rich context:

```bash
bd create --title "[Category]: [Specific issue]" \
  --type bug \
  --priority [0-4] \
  --add-label code-quality \
  --add-label [specific-category] \
  --description "[See template below]"
```

### Priority Guidelines
- **P0**: Security vulnerability, data loss risk, crash
- **P1**: Performance regression, startup time impact >10ms
- **P2**: Go idiom violation, code smell (default for most issues)
- **P3**: Minor improvement, code cleanup
- **P4**: Nice-to-have refactoring

### Category Labels

```bash
--add-label idioms          # Go idiom violations
--add-label performance     # Performance issues
--add-label security        # Security concerns
--add-label testing         # Test coverage gaps
--add-label maintainability # Code clarity issues
--add-label documentation   # Missing docs
```

### Rich Description Template

Your quality bead description MUST include:

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
[How to address the issue]

## References
[Links to Go best practices, benchmarks, or docs]

## Acceptance Criteria
[What makes this issue resolved]
```

### Example Quality Bead (Performance)

```bash
bd create --title "Performance: Reduce allocations in query builder" \
  --type bug \
  --priority 2 \
  --add-label code-quality \
  --add-label performance \
  --description "
## Issue
GraphQL query builder uses string concatenation in loop, causing excessive
allocations during query construction.

## Location
pkg/client/query_builder.go:78-95

Current code:
```go
func (qb *QueryBuilder) Build() string {
    query := \"\"
    for _, field := range qb.fields {
        query += field + \"\\n\"  // Allocates new string each iteration
    }
    return query
}
```

## Current State
Benchmark shows 1.2MB allocated per query build for typical 50-field query.
Each concatenation creates new string, copying all previous content.

## Why This Matters
- Query building is on hot path (called for every API request)
- Startup time budget is <50ms, allocations hurt this
- Benchmark: BenchmarkQueryBuilder-8  1000  1247 ns/op  1234 B/op  52 allocs/op

## Suggested Fix
Use strings.Builder with capacity pre-allocation:

```go
func (qb *QueryBuilder) Build() string {
    // Pre-allocate capacity: ~30 bytes per field
    var sb strings.Builder
    sb.Grow(len(qb.fields) * 30)

    for _, field := range qb.fields {
        sb.WriteString(field)
        sb.WriteByte('\\n')
    }
    return sb.String()
}
```

Expected improvement: ~50x fewer allocations, ~10x faster

## References
- Effective Go: String building: https://go.dev/doc/effective_go#building_strings
- strings.Builder docs: https://pkg.go.dev/strings#Builder
- Zero-allocation patterns: lirt performance requirements

## Acceptance Criteria
- strings.Builder used with pre-allocation
- Benchmark shows <25 allocations/op
- Query building time <100ns/op
- No startup time regression
"
```

### Example Quality Bead (Security)

```bash
bd create --title "Security: API key exposed in error messages" \
  --type bug \
  --priority 0 \
  --add-label code-quality \
  --add-label security \
  --description "
## Issue
Error messages include full API key when authentication fails.

## Location
pkg/client/auth.go:34

Current code:
```go
if err := c.authenticate(apiKey); err != nil {
    return fmt.Errorf(\"authentication failed with key %s: %w\", apiKey, err)
}
```

## Current State
When authentication fails, error message printed to stderr includes full API key.
User might paste error message in GitHub issue or Slack, exposing key.

## Why This Matters
- CRITICAL SECURITY ISSUE
- API keys grant full access to Linear workspace
- Keys in error messages → keys in logs → potential exposure
- Violates API key handling best practices

## Suggested Fix
Redact API key in error messages, show only prefix:

```go
func redactKey(key string) string {
    if len(key) < 12 {
        return \"***\"
    }
    return key[:8] + \"...\" + key[len(key)-4:]
}

if err := c.authenticate(apiKey); err != nil {
    return fmt.Errorf(\"authentication failed with key %s: %w\", redactKey(apiKey), err)
}
```

## References
- OWASP API Security: Don't leak credentials
- Linear API best practices
- lirt security requirements

## Acceptance Criteria
- All error messages redact API keys
- Only prefix (8 chars) and suffix (4 chars) shown
- Audit all fmt.Errorf calls for API key exposure
- Add test verifying key redaction
"
```

### Example Quality Bead (Go Idioms)

```bash
bd create --title "Idioms: Missing context propagation in API calls" \
  --type bug \
  --priority 1 \
  --add-label code-quality \
  --add-label idioms \
  --description "
## Issue
API client methods don't accept context.Context, preventing timeout/cancellation
control.

## Location
pkg/client/client.go:125-180 (all public methods)

Current signature:
```go
func (c *Client) ListIssues(opts *ListOptions) ([]*Issue, error)
```

## Current State
All API methods use background context internally:
- No timeout control
- Can't cancel long-running requests
- Can't propagate request context from CLI commands

## Why This Matters
- Not idiomatic Go (all I/O should accept context)
- User can't ctrl-c to cancel slow requests
- Can't set per-request timeouts
- Blocks compliance with Go best practices

## Suggested Fix
Add context.Context as first parameter to all public methods:

```go
func (c *Client) ListIssues(ctx context.Context, opts *ListOptions) ([]*Issue, error) {
    req, err := c.buildRequest(ctx, \"ListIssues\", opts)
    if err != nil {
        return nil, err
    }
    // Use ctx in http.NewRequestWithContext
}
```

Update all callers to pass context from command execution.

## References
- Effective Go: Contexts
- Go blog: Context patterns
- Standard library: net/http uses context

## Acceptance Criteria
- All public client methods accept context.Context as first param
- Context used in HTTP requests
- Context passed from CLI commands
- Tests verify cancellation works
- Breaking change documented in CHANGELOG
"
```

## Minimum Description Length

Quality bead descriptions should be **200-400 words minimum** to provide adequate context for fixing.

## Self-Check Before Creating

Before creating a quality bead, verify:
- ✅ Issue is specific with file:line location
- ✅ Current state clearly described
- ✅ Impact/importance explained
- ✅ Suggested fix included
- ✅ Acceptance criteria defined
- ✅ Priority appropriate for severity

## Processing Quality Beads

The lirt-specialist will:
1. Query: `bd ready --type bug --label code-quality`
2. Read quality bead descriptions
3. Implement fixes
4. Verify acceptance criteria met
5. Close bead: `bd close <id>`

## Re-Review After Fix

After lirt-specialist closes a quality bead, lirt-code-reviewer can:
1. Review the fix
2. Verify issue resolved
3. Confirm bead can stay closed
4. Or reopen if fix incomplete: `bd reopen <id>`

## Success Metrics

Track protocol effectiveness:
- Number of quality issues found per review
- Resolution time for quality beads
- Reduction in quality issues over time
- Code quality trend (fewer P0/P1 issues)
