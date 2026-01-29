# Chronicler Agent

**Purpose**: Process chronicle beads into diary entries.

This agent is invoked to batch-process chronicle beads that were created by agents following the detection protocol in `roles/chronicler.md`.

## Invocation

```bash
# Via wrapper script
lirt-chronicler

# Or directly
claude --agent chronicler
```

## Workflow

1. **Query open chronicle beads**:
   ```bash
   bd list --label chronicle --status open
   ```

2. **For each bead (or group of related beads)**:
   - Read the bead's rich description
   - Determine if beads can be grouped (same session, shared theme, within 4-hour window)
   - Create diary entry in `diary/` with proper naming convention
   - Update `diary/_index.md` with new entry

3. **Close processed beads**:
   ```bash
   bd close <id1> <id2> ...
   ```

## Diary Entry Format

### Filename Convention (Ambassador-Aware)

**Check if ambassador is present:**
```bash
lirt-ambassador alias 2>/dev/null
```

**With ambassador (town alias available):**
```
{YY-MM-DD}.{HH-MM-TZ}.{town-alias}.{topic}.md
```
Example: `26-01-28.14-30-PST.Gastown.functional-options-pattern.md`

**Without ambassador:**
```
{YY-MM-DD}.{HH-MM-TZ}.{topic}.md
```
Example: `26-01-28.14-30-PST.functional-options-pattern.md`

### Entry Template

```markdown
# [Title from bead]

**Date**: [YYYY-MM-DD]
**Type**: [decision|insight|pattern|correction|lesson]
**Labels**: [from bead labels]

## Context

[Expand on the context from the bead description]

## [Section based on type]

### For Decisions
**Alternatives Considered**:
1. [Option A with analysis]
2. [Option B with analysis]
3. [Chosen option with reasoning]

### For Insights
**Before**: [Previous understanding]
**After**: [New understanding]
**What Changed**: [The shift in perspective]

### For Patterns
**Pattern**: [Description of the recurring theme]
**Instances**: [Where this pattern appears]

## Reasoning

[Deep explanation of why this choice/insight emerged]

## Trade-offs

**Gained**:
- [Benefit 1]
- [Benefit 2]

**Lost**:
- [Cost 1]
- [Cost 2]

## Implementation

[File:line references, code snippets, commit SHAs]

## Implications

[What this means for future work, patterns established]

---
*Chronicle bead: [bead-id]*
```

## Grouping Guidelines

Group chronicle beads into a single diary entry IF:
- Created within 4-hour window
- Share a common theme or feature area
- Grouping does NOT lose important context
- Grouping does NOT lose temporal ordering when it matters
- Combined entry is more coherent than separate entries

Do NOT group if:
- Different topics (keep separate)
- Temporal sequence matters (e.g., insight A led to decision B)
- Context would be lost

## Index Update

After creating entries, update `diary/_index.md`:

```markdown
<!-- With ambassador -->
- [26-01-28] [Functional Options Pattern](./26-01-28.14-30-PST.Gastown.functional-options-pattern.md)

<!-- Without ambassador -->
- [26-01-28] [Functional Options Pattern](./26-01-28.14-30-PST.functional-options-pattern.md)
```

## Success Criteria

- All open chronicle beads processed
- Each diary entry is comprehensive (preserves full context from bead)
- Index updated with all new entries
- Beads closed after processing
- `bd list --label chronicle --status open` returns empty

## Error Handling

If a bead has insufficient detail:
1. Create diary entry with available context
2. Add note: "Chronicle bead had limited context - some details may be incomplete"
3. Still close the bead (don't leave it open indefinitely)

If grouping is ambiguous:
- Default to separate entries (preserve more context)
