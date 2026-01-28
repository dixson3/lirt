# Test Case Identifier Protocol

**Role**: Detection protocol for identifying testable scenarios
**Applies to**: ALL agents working on lirt

## Purpose

During lirt development, agents naturally encounter testable scenarios while implementing features, writing specifications, or reviewing code. This protocol ensures those scenarios are captured for comprehensive test coverage.

## When to Create Test Beads

**ALL agents** must watch for testable scenarios:

### During Specification Writing (lirt-spec-writer)
- Edge cases in command behavior
- Error conditions in CLI usage
- Complex flag combinations
- Unusual input patterns
- API error responses
- Format conversions (JSON, CSV, table)

### During Implementation (lirt-specialist)
- Edge cases discovered in code
- Error handling paths
- Complex business logic branches
- State machine transitions
- Concurrent operation scenarios
- Performance-critical code paths

### During Code Review (lirt-code-reviewer)
- Untested code paths
- Missing error handling tests
- Security-sensitive operations without tests
- Performance optimizations needing benchmarks
- Edge cases not covered by existing tests

## Creating Test Beads

When you identify a testable scenario, create a test bead with rich context:

```bash
bd create --title "Test: [specific scenario]" \
  --type task \
  --priority 2 \
  --add-label requires:testing \
  --description "[See template below]"
```

### Priority Guidelines
- **P0**: Critical path, blocks release (authentication, data loss scenarios)
- **P1**: Important functionality, needed for stability (error handling, edge cases)
- **P2**: Standard tests (default for most scenarios)
- **P3**: Nice-to-have, optimization tests
- **P4**: Future consideration

### Rich Description Template

Your test bead description MUST include:

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

### Example Test Bead

```bash
bd create --title "Test: issue list filter edge cases" \
  --type task \
  --priority 2 \
  --add-label requires:testing \
  --description "
## Scenario
The 'lirt issue list --filter' command needs comprehensive edge case testing
for invalid filter syntax and empty result handling.

## Context
Writing CLI reference documentation for issue list command. Documented the
filter syntax (field:operator:value) but realized we haven't tested invalid
inputs or empty results.

## Test Cases Needed
1. Invalid syntax: --filter 'invalid' → clear error message
2. Invalid field: --filter 'badfield:eq:foo' → 'unknown field' error
3. Invalid operator: --filter 'title:badop:foo' → 'unknown operator' error
4. Empty results: --filter 'title:eq:nonexistent' → empty table, not error
5. Multiple filters: --filter 'state:eq:done' --filter 'team:eq:ENG'
6. Conflicting filters: --filter 'state:eq:done' --filter 'state:eq:todo'

## Why This Matters
- Filters are primary way users query issues
- Poor error messages frustrate users
- Need to distinguish 'no results' from 'invalid query'
- Multiple filters have AND semantics that must be clear

## Implementation Notes
- Filter parsing: cmd/issue/filters.go:45-120
- Use table-driven tests with subtests
- Mock Linear API responses for each scenario
- Test both table and JSON output formats

## Acceptance Criteria
- All 6 test cases have passing tests
- Error messages are clear and actionable
- Empty results handled gracefully
- Test execution < 500ms total
"
```

## Labels for Test Organization

Add specific labels to help lirt-test-engineer prioritize:

```bash
# Type of testing
--add-label unit-test
--add-label integration-test
--add-label benchmark

# Component
--add-label cli
--add-label api-client
--add-label config

# Urgency
--add-label blocks-release
--add-label regression-risk
```

## Minimum Description Length

Test bead descriptions should be **150-300 words minimum** to provide adequate context.

If you can't capture test requirements richly:
- Break into multiple smaller, focused test beads
- Discuss with user to clarify requirements
- Reference specification documents for context

## Self-Check Before Creating

Before creating a test bead, verify:
- ✅ Scenario is specific and testable
- ✅ Test cases are clearly enumerated
- ✅ Expected behavior is defined
- ✅ Relevant code locations referenced
- ✅ Priority appropriate for importance

## Processing Test Beads

The lirt-test-engineer will:
1. Query: `bd ready --label requires:testing`
2. Read test bead descriptions
3. Implement tests following Go testing patterns
4. Verify all acceptance criteria met
5. Close bead: `bd close <id>`

## Success Metrics

Track protocol effectiveness:
- Number of test beads created per feature
- Test coverage improvement over time
- Percentage of edge cases caught before production
- Time from test identification to implementation
