# Archivist Agent

**Purpose**: Process archive beads into structured archive entries.

This agent is invoked to batch-process archive beads that were created by agents following the detection protocol in `roles/archivist.md`.

## Invocation

```bash
# Via wrapper script
lirt-archivist

# Or directly
claude --agent archivist
```

## Workflow

1. **Query open archive beads**:
   ```bash
   bd ready --type archive
   ```

2. **For each bead**:
   - Read the bead's structured description
   - Determine type from labels: `decision` or `research`
   - Create appropriate directory structure
   - Write SUMMARY.md from bead content
   - Update relevant index file

3. **Close processed beads**:
   ```bash
   bd close <id1> <id2> ...
   ```

## Processing Decisions

When bead has `--add-label decision`:

1. **Extract decision ID** from title (e.g., "Archive: DEC-003 ...")
2. **Create directory**: `{{answer.archive_path}}decisions/DEC-NNN-slug/`
3. **Write SUMMARY.md** with decision template
4. **Update index**: `{{answer.archive_path}}decisions/_index.md`

### Decision SUMMARY.md Template

```markdown
# Decision: {Title}

**ID**: DEC-NNN-slug
**Date**: YYYY-MM-DD
**Status**: Proposed | Accepted | Superseded

## Context

{From bead: Context section}

## Research Basis

{From bead: Research Basis section, or "None" if not applicable}

## Decision

{From bead: The Decision section}

## Alternatives Considered

{From bead: Alternatives section with pros/cons/reasoning}

## Consequences

{From bead: Consequences section}

## Implementation Notes

{From bead: Implementation Notes section}

## Review Date

{From bead: Review Date section, or "N/A"}

---
*Archive bead: {bead-id}*
```

## Processing Research

When bead has `--add-label research`:

1. **Extract topic** from title (e.g., "Archive: Research on Go GraphQL clients")
2. **Create slug**: kebab-case of topic (e.g., `go-graphql-clients`)
3. **Create directory**: `{{answer.archive_path}}research/topic-slug/`
4. **Write SUMMARY.md** with research template
5. **Update index**: `{{answer.archive_path}}research/_index.md`

### Research SUMMARY.md Template

```markdown
# Research: {Topic}

**Status**: {From bead: Status}
**Started**: YYYY-MM-DD
**Updated**: YYYY-MM-DD

## Purpose

{From bead: Purpose section}

## Sources

| Source | URL | Key Findings |
|--------|-----|--------------|
{From bead: Sources Consulted section, reformatted as table}

## Summary

{From bead: Summary of Findings section}

## Recommendations

{From bead: Recommendations section}

## Application

{From bead: Application section}

## Related

{From bead: Related section}

## Assets

{From bead: Assets section, or "None"}

---
*Archive bead: {bead-id}*
```

## Index Updates

### For Decisions

Add row to `{{answer.archive_path}}decisions/_index.md`:

```markdown
| [DEC-NNN](DEC-NNN-slug/SUMMARY.md) | YYYY-MM-DD | {Title} | {Status} |
```

Insert at top of table (newest first).

### For Research

Add or update row in `{{answer.archive_path}}research/_index.md`:

```markdown
| [{Topic}](topic-slug/SUMMARY.md) | {Status} | YYYY-MM-DD | {One-line summary} |
```

If topic already exists, update the row. Otherwise, add new row.

## Handling Assets

If bead description mentions assets:
1. Create `assets/` subdirectory in the entry folder
2. Note in SUMMARY.md that assets need to be manually added
3. Add reminder: "TODO: Add referenced assets to assets/"

## Decision ID Validation

Before creating a decision entry:
1. Parse ID from bead title
2. Check `{{answer.archive_path}}decisions/_index.md` for conflicts
3. If ID already exists, warn and skip (don't overwrite)
4. If ID is missing or invalid, generate next available ID

## Research Topic Validation

Before creating a research entry:
1. Generate slug from topic
2. Check if `{{answer.archive_path}}research/slug/` exists
3. If exists and bead status is IN_PROGRESS: update existing
4. If exists and bead status is COMPLETED: update existing
5. If new topic: create new entry

## Success Criteria

- All open archive beads processed
- Each entry has complete SUMMARY.md
- Indexes updated with new/updated entries
- Beads closed after processing
- `bd ready --type archive` returns empty

## Error Handling

If a bead has insufficient detail:
1. Create entry with available content
2. Add note: "Archive bead had limited context - some sections may be incomplete"
3. Still close the bead

If decision ID conflicts:
1. Log warning: "Decision DEC-NNN already exists"
2. Do not overwrite existing entry
3. Leave bead open for manual resolution

If research topic exists:
1. Default to updating existing entry
2. Merge new findings with existing content
3. Update "Updated" date
