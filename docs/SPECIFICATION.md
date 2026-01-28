# lirt — Linear CLI Tool Specification

**Version**: 0.1.0
**Status**: DRAFT
**Language**: Go 1.21+
**Target**: CLI tool for Linear GraphQL API
**Last Updated**: 2026-01-27

---

## Table of Contents

1. [Overview](#1-overview)
2. [Configuration & Credentials](#2-configuration--credentials)
3. [Command Structure](#3-command-structure)
4. [Commands](#4-commands)
5. [Output Formats](#5-output-formats)
6. [Caching](#6-caching)
7. [Scripting Patterns](#7-scripting-patterns)
8. [Technology Stack](#8-technology-stack)
9. [Error Handling](#9-error-handling)
10. [Pagination](#10-pagination)
11. [Testing Strategy](#11-testing-strategy)
12. [Future Considerations](#12-future-considerations)

---

## 1. Overview

`lirt` is a CLI tool for interacting with Linear, following `gh` CLI semantics. It targets bash scripting workflows alongside `jq` and `ripgrep`, and supports multiple Linear workspaces via named profiles (modeled on AWS CLI).

**Design Principles:**
- **Fast startup**: < 50ms (optimized for bash scripting)
- **Low memory**: < 50MB baseline
- **Scriptable**: JSON output, pipe-friendly, composable with Unix tools
- **Multi-workspace**: AWS CLI-style profile management
- **Idiomatic Go**: Functional options, context propagation, table-driven tests

---

## 2. Configuration & Credentials

### API Endpoint

Linear uses a **single global API endpoint** for all workspaces:

```
https://api.linear.app/graphql
```

The workspace is determined entirely by the API key — each key is scoped to exactly one workspace. There is no per-workspace URL, subdomain, or tenant identifier in the API. This means:

- Profile names are **user-chosen aliases**, not Linear workspace names
- No URL or workspace ID field is needed in configuration
- The `viewer` query returns the workspace name/org for the authenticated key
- Switching workspaces means switching API keys (profiles)

### File Layout

```
~/.config/lirt/
├── credentials       # API keys (profile-based, INI format)
├── config            # Settings per profile (INI format)
└── cache/            # Cached enumeration data (auto-managed)
    ├── <profile>/
    │   ├── teams.json
    │   ├── states.json
    │   ├── labels.json
    │   ├── users.json
    │   └── priorities.json
    └── ...
```

File permissions: `credentials` created with `0600`, `config` with `0644`, `cache/` with `0700`.

### credentials (INI format)

Modeled on `~/.aws/credentials`. Profile names in brackets, no `profile` prefix.

```ini
[default]
api_key = lin_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

[acme-corp]
api_key = lin_api_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy

[side-project]
api_key = lin_api_zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
```

### config (INI format)

Modeled on `~/.aws/config`. Profile names use `profile` prefix (except `[default]`).

Profile names are **user-chosen aliases** — they do not need to match the Linear workspace name. The `workspace` config key is a display-only label populated automatically by `lirt auth login` (via the `viewer` query) for human reference.

```ini
[default]
workspace = Acme Corporation
team = ENG
format = table

[profile side-project]
workspace = Side Project Inc
team = SP
format = table
```

### Supported Config Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `workspace` | string | (auto) | Display name of the Linear workspace (set by `lirt auth login`) |
| `team` | string | (none) | Default team key (avoids `--team` on every command) |
| `format` | string | `table` | Default output format: `table`, `json`, `csv`, `plain` |
| `cache_ttl` | duration | `5m` | How long to cache enumeration data |
| `page_size` | int | `50` | Default pagination limit |

### Credential Resolution (priority order)

1. `LIRT_API_KEY` environment variable (always wins)
2. `--api-key` flag (per-command override)
3. `~/.config/lirt/credentials` file, selected profile
4. `LINEAR_API_KEY` environment variable (fallback compatibility)

### Profile Selection (priority order)

1. `--profile <name>` flag
2. `LIRT_PROFILE` environment variable
3. `[default]` profile

### Config File Override

| Override | Purpose |
|----------|---------|
| `LIRT_CONFIG_DIR` | Override `~/.config/lirt/` entirely |
| `LIRT_CREDENTIALS_FILE` | Override credentials file path |
| `LIRT_CONFIG_FILE` | Override config file path |

---

## 3. Command Structure

```
lirt <resource> <action> [args] [flags]
```

### Global Flags

| Flag | Short | Type | Description |
|------|-------|------|-------------|
| `--profile` | `-P` | string | Named profile to use |
| `--api-key` | | string | Override API key for this invocation |
| `--team` | `-t` | string | Team key context (overrides config) |
| `--format` | `-f` | string | Output format: `table`, `json`, `csv`, `plain` |
| `--json` | | string | Output specific fields as JSON (comma-separated) |
| `--jq` | | string | Apply jq expression to JSON output |
| `--no-cache` | | bool | Bypass cached data |
| `--quiet` | `-q` | bool | Suppress non-essential output |
| `--verbose` | `-v` | bool | Debug output |
| `--help` | `-h` | bool | Help at any level |
| `--version` | `-V` | bool | Print version |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | API or runtime error |
| 2 | Usage error (bad flags, missing args) |
| 3 | Authentication error |
| 4 | Not found (entity doesn't exist) |

---

## 4. Commands

### 4.1 auth — Credential Management

```bash
lirt auth login [--profile <name>]              # Prompt for API key, store in credentials file
lirt auth login --api-key <key> [--profile <name>]  # Non-interactive
lirt auth status [--profile <name>]             # Show auth state (workspace, user, permissions)
lirt auth token [--profile <name>]              # Print API key to stdout (for piping)
lirt auth logout [--profile <name>]             # Remove profile from credentials file
lirt auth list                                  # List all configured profiles
lirt auth switch <profile>                      # Set LIRT_PROFILE in current shell (prints export command)
```

`lirt auth login` flow:
1. Prompt for API key (masked input, or `--api-key` flag)
2. Validate key by calling `viewer` query
3. Display workspace name and authenticated user
4. Write to `~/.config/lirt/credentials` under specified profile (default: `[default]`)
5. If profile already exists, confirm overwrite

`lirt auth list` output:
```
PROFILE         WORKSPACE            KEY PREFIX
default         Acme Corporation     lin_api_xxxx...
side-project    Side Project Inc     lin_api_zzzz...
```

`lirt auth status` output:
```
Profile:    default
Workspace:  Acme Corporation
User:       James Dixson (james@acme.com)
Permissions: Read, Write
Key prefix: lin_api_xxxx...
```

### 4.2 team — Team Operations

```bash
lirt team list                                  # All teams (id, key, name, member count)
lirt team view <key-or-id>                      # Team details
lirt team members <key-or-id>                   # List team members
lirt team states <key-or-id>                    # Workflow states for team
lirt team labels <key-or-id>                    # Labels for team
lirt team cycles <key-or-id>                    # Cycles for team (current, upcoming, past)
```

See [COMMANDS.md](./COMMANDS.md) for detailed command documentation.

### 4.3 issue — Issue Operations

```bash
# List / Search
lirt issue list [filters] [flags]
lirt issue search <query> [--team <key>]

# CRUD
lirt issue create --title "..." [options]
lirt issue view <id>
lirt issue edit <id> [options]

# State transitions
lirt issue close <id>
lirt issue reopen <id>
lirt issue transition <id> <state-name>

# Archive / Delete
lirt issue archive <id>
lirt issue delete <id> [--confirm]

# Labels & Assignment
lirt issue label <id> --add <name>... --remove <name>...
lirt issue assign <id> <login-or-email>
lirt issue unassign <id>

# Relations
lirt issue children <id>
lirt issue parent <id>
```

**ID resolution**: All `<id>` arguments accept both the shorthand identifier (e.g., `ENG-123`) and the UUID. The shorthand is always preferred for display.

**Priority values**: Accept either numeric (0-4) or named (`urgent`, `high`, `medium`, `low`, `none`). Display uses both: `P0 (Urgent)`.

### 4.4 project — Project Operations

```bash
lirt project list [--team <key>] [--state <name>] [--limit <n>]
lirt project view <id-or-name>
lirt project issues <id-or-name> [--state <name>] [--limit <n>]
lirt project milestones <id-or-name>
lirt project members <id-or-name>
lirt project create --title "..." [options]
lirt project edit <id> [options]
lirt project archive <id>
lirt project delete <id> [--confirm]
```

Project states: `backlog`, `planned`, `started`, `paused`, `completed`, `canceled`.

### 4.5 milestone — Project Milestone Operations

```bash
lirt milestone list --project <id-or-name>
lirt milestone view <id>
lirt milestone create --project <id-or-name> --title "..." [options]
lirt milestone edit <id> [options]
lirt milestone delete <id> [--confirm]
lirt milestone issues <id> [--state <name>] [--limit <n>]
```

### 4.6 initiative — Initiative Operations

```bash
lirt initiative list [--limit <n>]
lirt initiative view <id-or-name>
lirt initiative create --title "..." [--description "..."]
lirt initiative edit <id> [options]
lirt initiative archive <id>
lirt initiative delete <id> [--confirm]
lirt initiative projects <id-or-name>
```

### 4.7 user — User Operations

```bash
lirt user list [--limit <n>]
lirt user view <id-or-login-or-email>
lirt user me                                    # Current authenticated user
lirt user issues <id-or-login> [--state <name>] [--limit <n>]
```

### 4.8 comment — Comment Operations

```bash
lirt comment list <issue-id> [--limit <n>]
lirt comment add <issue-id> --body "..."
lirt comment add <issue-id> --body-file <path>
lirt comment edit <comment-id> --body "..."
lirt comment delete <comment-id> [--confirm]
```

### 4.9 meta — Enumeration / Reference Data

```bash
lirt meta states [--team <key>]                 # Workflow states (type, name, color)
lirt meta priorities                            # Priority levels (0=Urgent through 4=None)
lirt meta labels [--team <key>]                 # Labels (name, color, scope)
lirt meta cycles [--team <key>]                 # Cycles (name, dates, state)
lirt meta issue-types                           # Available issue types if custom types enabled
```

### 4.10 api — Raw GraphQL Access

```bash
lirt api <query-string>                         # Inline GraphQL
lirt api --input <file.graphql>                 # Query from file
lirt api --input - < query.graphql              # Query from stdin
lirt api -f field=value <query>                 # Variables via flags
```

Escape hatch for operations not covered by built-in commands. Always outputs JSON.

### 4.11 config — Configuration Management

```bash
lirt config list [--profile <name>]             # Show all config for profile
lirt config get <key> [--profile <name>]        # Get specific config value
lirt config set <key> <value> [--profile <name>] # Set config value
lirt config unset <key> [--profile <name>]      # Remove config value
```

### 4.12 completion — Shell Completions

```bash
lirt completion bash                            # Output bash completions
lirt completion zsh                             # Output zsh completions
lirt completion fish                            # Output fish completions
```

---

## 5. Output Formats

See [COMMANDS.md](./COMMANDS.md#output-formats) for detailed output format documentation.

### Supported Formats

- **table** (default): Human-readable aligned columns
- **json**: Full JSON output or field selection with `--json <fields>`
- **csv**: Comma-separated values for spreadsheet import
- **plain**: Minimal output, one value per line

### Automatic Format Detection

When stdout is not a terminal (piped), default to `json` instead of `table`. Override with explicit `--format`.

---

## 6. Caching

### What's Cached

| Data | Cache File | TTL Default | Invalidated By |
|------|-----------|-------------|----------------|
| Teams | `cache/<profile>/teams.json` | 5m | `lirt team` write ops |
| Workflow states | `cache/<profile>/states.json` | 5m | `lirt meta states --no-cache` |
| Labels | `cache/<profile>/labels.json` | 5m | Label write ops |
| Users | `cache/<profile>/users.json` | 5m | — |
| Priorities | `cache/<profile>/priorities.json` | 24h | — (static) |

### Cache Behavior

- `--no-cache` bypasses cache for the current command
- Write operations invalidate the relevant cache
- Cache files include a `fetched_at` timestamp; expired entries are refreshed transparently
- `lirt config set cache_ttl 0` disables caching entirely

---

## 7. Scripting Patterns

See [SCRIPTING.md](./SCRIPTING.md) for comprehensive bash integration patterns.

### Pipe-friendly Defaults

```bash
# When piped, output is JSON by default
lirt issue list --team ENG | jq '.[].identifier'

# Explicit fields
lirt issue list --json id,title,assignee | jq -r '.[] | [.id, .title] | @tsv'
```

### Batch Operations

```bash
# Close all issues in a milestone
lirt milestone issues <id> --json id --jq '.[].id' \
  | xargs -I{} lirt issue close {}
```

---

## 8. Technology Stack

### Core Dependencies

| Dependency | Purpose | Package |
|------------|---------|---------|
| **Cobra** | CLI framework, subcommands, flags, completions | `github.com/spf13/cobra` |
| **Viper** | Config file parsing (INI/YAML), env var binding | `github.com/spf13/viper` |
| **go-graphql-client** | GraphQL client (queries, mutations) | `github.com/hasura/go-graphql-client` |
| **go-ini** | INI file parsing for AWS-style credentials | `gopkg.in/ini.v1` |
| **tablewriter** | Table output formatting | `github.com/olekukonez/tablewriter` |
| **color** | Terminal color output | `github.com/fatih/color` |

### Project Structure

```
lirt/
├── cmd/                    # Cobra command definitions
│   ├── root.go             # Root command + global flags
│   ├── auth.go             # lirt auth *
│   ├── team.go             # lirt team *
│   ├── issue.go            # lirt issue *
│   ├── project.go          # lirt project *
│   ├── milestone.go        # lirt milestone *
│   ├── initiative.go       # lirt initiative *
│   ├── user.go             # lirt user *
│   ├── comment.go          # lirt comment *
│   ├── meta.go             # lirt meta *
│   ├── api.go              # lirt api
│   ├── config.go           # lirt config *
│   └── completion.go       # lirt completion *
├── internal/
│   ├── client/             # GraphQL client wrapper
│   │   ├── client.go       # Linear API client (auth, rate limiting, pagination)
│   │   ├── queries.go      # GraphQL query definitions
│   │   └── mutations.go    # GraphQL mutation definitions
│   ├── config/             # Configuration management
│   │   ├── credentials.go  # Credential file read/write
│   │   ├── config.go       # Config file read/write
│   │   └── profile.go      # Profile resolution logic
│   ├── cache/              # Cache management
│   │   └── cache.go        # Read/write/invalidate cached data
│   ├── output/             # Output formatting
│   │   ├── table.go        # Table format
│   │   ├── json.go         # JSON format (with --json field selection)
│   │   ├── csv.go          # CSV format
│   │   └── plain.go        # Plain format
│   └── model/              # Domain types
│       ├── issue.go
│       ├── project.go
│       ├── team.go
│       ├── initiative.go
│       ├── milestone.go
│       ├── user.go
│       ├── comment.go
│       └── state.go
├── main.go                 # Entry point
├── go.mod
├── go.sum
├── Makefile                # Build, test, lint targets
├── LICENSE                 # MIT
└── README.md
```

---

## 9. Error Handling

### GraphQL Partial Success

Linear's GraphQL API can return HTTP 200 with partial data and an `errors` array. The client must always check for errors:

```go
// Pseudo-code
resp, err := client.Query(ctx, &query, vars)
if err != nil {
    // Network or transport error
    return fmt.Errorf("API request failed: %w", err)
}
if len(resp.Errors) > 0 {
    // GraphQL-level errors (partial success)
    return fmt.Errorf("API error: %s", resp.Errors[0].Message)
}
```

### Rate Limiting

- Track remaining requests from response headers
- On 429: wait for `Retry-After` header, then retry (max 3 retries)
- Display rate limit info with `--verbose`
- With `--quiet`, rate limit warnings are suppressed unless blocking

### User-Facing Errors

```
Error: issue ENG-999 not found
Error: not authenticated — run 'lirt auth login' to set up credentials
Error: team 'INVALID' not found — run 'lirt team list' to see available teams
Error: rate limit exceeded — retry in 42 seconds
```

Always suggest the corrective action.

---

## 10. Pagination

All list commands support:

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 50 | Max items to return |
| `--all` | bool | false | Fetch all pages (ignore limit) |
| `--cursor` | string | (none) | Resume from cursor (for manual pagination) |

Implementation uses Relay-style cursor pagination:
1. First request: `first: <limit>`
2. If `pageInfo.hasNextPage` and `--all`: continue with `after: <endCursor>`
3. Collect all `nodes` arrays, concatenate
4. With `--cursor`: start from that cursor instead of the beginning

The `--all` flag streams results as they arrive (in JSON array format), so downstream `jq` processing can begin immediately.

---

## 11. Testing Strategy

### Linear Has No Sandbox API

Linear does **not** provide a sandbox, staging, or test API environment. There is no programmatic workspace creation — workspaces are created manually via the Linear web UI. The same `https://api.linear.app/graphql` endpoint serves all environments.

### Testing Approach: Pre-Provisioned Free Workspace

Linear's free plan is suitable for automated testing:

| Feature | Free Plan |
|---------|-----------|
| Price | $0 |
| Active issues | 250 |
| Members | Unlimited |
| API access | Full (same API as paid plans) |
| Rate limit | 1,500 req/hr per API key |
| All members are Admins | Yes |

**Setup (one-time manual steps)**:

1. Create a dedicated Linear workspace for testing (e.g., "lirt-test")
2. Generate a Personal API Key with full permissions (Read, Write, Admin)
3. Create at least one team (e.g., `TEST`) — required for issue creation
4. Store the API key as a CI secret

**CI Environment Configuration**:

```bash
# Environment variables for CI
LIRT_TEST_API_KEY=lin_api_...       # API key for test workspace
LIRT_TEST_TEAM=TEST                  # Team key in test workspace
```

### Test Categories

| Category | Strategy | Needs API |
|----------|----------|-----------|
| **Unit tests** | Mock the GraphQL client interface; test command parsing, output formatting, config resolution, cache logic | No |
| **Integration tests** | Hit the real Linear API with the test workspace; create/read/update/delete actual entities; clean up after each test | Yes |
| **E2E tests** | Run `lirt` binary as a subprocess; verify stdout/stderr/exit codes for complete command invocations | Yes |

### Integration Test Design

```go
// Test fixture: create entities at test start, tear down at end
func TestIssueLifecycle(t *testing.T) {
    if os.Getenv("LIRT_TEST_API_KEY") == "" {
        t.Skip("LIRT_TEST_API_KEY not set, skipping integration test")
    }

    // Create → Read → Update → Archive → Delete
    // Each step verifies the API response and entity state
}
```

**Guard against flaky tests**:
- Skip integration tests when `LIRT_TEST_API_KEY` is unset (local dev without credentials)
- Use `t.Cleanup()` to delete test entities even on failure
- Prefix test entity titles with `[lirt-test]` for easy identification and manual cleanup
- Rate limit awareness: add small delays between API calls in test loops
- 250 active issue limit: ensure tests clean up, or use `--archived` entities where possible

### Test Workspace Maintenance

The 250 active issue limit on the free plan means test cleanup is critical. A CI job or Makefile target should:

```bash
# Clean up stale test entities (run periodically or as CI pre-step)
lirt issue list --profile test --json id,title \
  | jq -r '.[] | select(.title | startswith("[lirt-test]")) | .id' \
  | xargs -I{} lirt issue delete {} --confirm --profile test
```

---

## 12. Future Considerations (Out of Scope)

These are explicitly not part of the initial implementation but noted for potential future work:

- **OAuth 2.0 support**: For multi-user/distributed team scenarios
- **Webhook listener**: `lirt watch` for real-time event streaming
- **TUI mode**: Interactive issue browser
- **Issue templates**: `lirt issue create --template <name>`
- **Aliases**: `lirt alias set bugs 'issue list --label bug --state open'`
- **Extensions**: Plugin system for custom commands

---

## Related Documentation

- [REQUIREMENTS.md](./REQUIREMENTS.md) - Requirements document
- [COMMANDS.md](./COMMANDS.md) - CLI command reference
- [AUTHENTICATION.md](./AUTHENTICATION.md) - Authentication setup guide
- [CONFIGURATION.md](./CONFIGURATION.md) - Configuration file reference
- [SCRIPTING.md](./SCRIPTING.md) - Bash integration patterns
- [GRAPHQL.md](./GRAPHQL.md) - GraphQL implementation guide

---

**Source**: Based on research conducted in the clio project (cl-k0a1) on 2026-01-27.
