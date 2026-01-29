# Code Quality Reporter Protocol

**Role**: Detection protocol for identifying code quality issues
**Applies to**: ALL agents working on lirt

## Purpose

This protocol enables systematic detection and reporting of code quality issues via beads. Issues are captured with rich context for domain agents to address.

## When to Create Quality Beads

**ALL agents** must watch for quality issues:

### Language Idiom Violations
- Non-idiomatic patterns for the project's language
- Incorrect interface/type usage
- Missing error handling patterns
- Poor abstraction choices

### Performance Issues
- Inefficient algorithms or data structures
- Unnecessary allocations
- Missing caching opportunities
- Hot path optimizations needed

### Security Concerns
- Secrets in logs or errors
- Insecure file permissions
- Missing input validation
- Injection risks (SQL, command, etc.)
- Unhandled sensitive data

### Code Quality Issues
- High cyclomatic complexity
- Duplicated code
- Magic numbers without constants
- Missing error handling
- Unclear naming

### Testing Gaps
- Untested code paths
- Missing edge case tests
- No benchmarks for critical paths
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
- **P1**: Performance regression, critical path issue
- **P2**: Idiom violation, code smell (default for most issues)
- **P3**: Minor improvement, code cleanup
- **P4**: Nice-to-have refactoring

### Category Labels

```bash
--add-label idioms          # Language idiom violations
--add-label performance     # Performance issues
--add-label security        # Security concerns
--add-label testing         # Test coverage gaps
--add-label maintainability # Code clarity issues
--add-label documentation   # Missing docs
```

## Rich Description Template

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
[Links to best practices, benchmarks, or docs]

## Acceptance Criteria
[What makes this issue resolved]
```

## Example Quality Bead (Performance)

```bash
bd create --title "Performance: Reduce allocations in query builder" \
  --type bug \
  --priority 2 \
  --add-label code-quality \
  --add-label performance \
  --description "
## Issue
Query builder uses string concatenation in loop, causing excessive allocations.

## Location
pkg/client/query_builder.go:78-95

## Current State
Each concatenation creates new string, copying all previous content.
Benchmark shows 1.2MB allocated per query build for typical 50-field query.

## Why This Matters
- Query building is on hot path (called for every API request)
- Startup time budget is <50ms, allocations hurt this

## Suggested Fix
Use strings.Builder with capacity pre-allocation.

## References
- Effective Go: String building
- strings.Builder docs

## Acceptance Criteria
- strings.Builder used with pre-allocation
- Benchmark shows <25 allocations/op
"
```

## Example Quality Bead (Security)

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

## Current State
When auth fails, error message printed to stderr includes full API key.
User might paste error in GitHub issue, exposing key.

## Why This Matters
- CRITICAL SECURITY ISSUE
- API keys grant full access to resources
- Keys in error messages → keys in logs → potential exposure

## Suggested Fix
Redact API key in error messages, show only prefix/suffix.

## References
- OWASP: Don't leak credentials
- Security best practices

## Acceptance Criteria
- All error messages redact API keys
- Only prefix (8 chars) and suffix (4 chars) shown
- Add test verifying key redaction
"
```

## Minimum Description Length

Quality bead descriptions should be **200-400 words minimum** to provide adequate context for fixing.

## Self-Check Before Creating

Before creating a quality bead, verify:
- Issue is specific with file:line location
- Current state clearly described
- Impact/importance explained
- Suggested fix included
- Acceptance criteria defined
- Priority appropriate for severity

## Processing Quality Beads

Domain agents (specialist, developer) will:
1. Query: `bd ready --type bug --label code-quality`
2. Read quality bead descriptions
3. Implement fixes
4. Verify acceptance criteria met
5. Close bead: `bd close <id>`

## Re-Review After Fix

After a fix is applied, reviewers can:
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
