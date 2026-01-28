# lirt Requirements

**Version**: 0.1.0
**Status**: DRAFT
**Last Updated**: 2026-01-27

---

## Table of Contents

1. [Overview](#1-overview)
2. [Goals](#2-goals)
3. [Non-Goals](#3-non-goals)
4. [Functional Requirements](#4-functional-requirements)
5. [Non-Functional Requirements](#5-non-functional-requirements)
6. [Technical Constraints](#6-technical-constraints)
7. [Dependencies](#7-dependencies)
8. [Security Considerations](#8-security-considerations)

---

## 1. Overview

**lirt** is a command-line interface tool for interacting with Linear's issue tracking system. Following the design principles of the `gh` CLI (GitHub's official CLI tool), lirt provides a fast, scriptable, and intuitive interface to Linear's GraphQL API.

**Target Users:**
- Developers working with Linear in terminal-based workflows
- DevOps engineers automating issue management via bash scripts
- Teams needing multi-workspace Linear access from a single tool
- Claude Code agents performing automated issue tracking and coordination

**Core Value Proposition:**
- Fast startup (< 50ms) optimized for script execution
- AWS CLI-style multi-profile support for managing multiple Linear workspaces
- Pipe-friendly output formats (JSON, CSV, plain text, table)
- First-class bash integration with jq, ripgrep, and xargs

---

## 2. Goals

### G-1: Fast, Scriptable CLI Tool
Provide a CLI tool optimized for bash scripting workflows with sub-50ms startup time and composability with standard Unix tools (jq, ripgrep, xargs, awk).

### G-2: Multi-Workspace Support
Enable users to manage multiple Linear workspaces from a single tool using AWS CLI-style named profiles.

### G-3: Complete Linear API Coverage
Support all core Linear entities: issues, projects, milestones, initiatives, teams, users, comments, and workflow metadata.

### G-4: Consistent Command Structure
Follow `gh` CLI semantics for command structure, flag naming, and output formatting to leverage existing user mental models.

### G-5: Flexible Output Formats
Support multiple output formats (table, JSON, CSV, plain) with automatic format detection based on terminal vs pipe context.

### G-6: Efficient API Usage
Implement smart caching for enumeration data (teams, states, labels, users, priorities) to minimize API calls and improve response time.

### G-7: Developer-Friendly Errors
Provide clear error messages with suggested corrective actions and appropriate exit codes.

### G-8: Testable Implementation
Enable automated testing through unit tests, integration tests against a real Linear workspace, and E2E tests of CLI behavior.

---

## 3. Non-Goals

### NG-1: TUI (Terminal UI) Mode
No interactive terminal UI browser. lirt is strictly a CLI tool optimized for scripting and command-line workflows.

### NG-2: OAuth 2.0 Authentication
Initial version supports only Personal API Keys. OAuth 2.0 support deferred to future versions.

### NG-3: Webhook/Event Streaming
No real-time event streaming or webhook listener. Users should use Linear's built-in webhook functionality.

### NG-4: Offline Mode
lirt requires network connectivity to function. No offline read-only mode for cached data.

### NG-5: Issue Templates
No built-in issue template system. Users can script template behavior using bash and `--body-file`.

### NG-6: Alias System
No built-in command aliasing. Users can use shell aliases for custom shortcuts.

### NG-7: Plugin System
No extensibility via plugins. All functionality is built-in.

### NG-8: Custom GraphQL Client
Use existing Go GraphQL client libraries rather than building a custom client.

---

## 4. Functional Requirements

### FR-1: Authentication & Credential Management

**FR-1.1**: Support Personal API Key authentication
Users must be able to authenticate using Linear Personal API Keys.

**FR-1.2**: Store credentials securely
API keys stored in `~/.config/lirt/credentials` with file permissions `0600`.

**FR-1.3**: Multi-profile credential management
Support multiple named profiles (e.g., `default`, `work`, `personal`) using AWS CLI-style INI format.

**FR-1.4**: Credential resolution priority
Resolve credentials in order: `LIRT_API_KEY` env var → `--api-key` flag → credentials file → `LINEAR_API_KEY` env var.

**FR-1.5**: Auth validation
Validate API keys by calling Linear's `viewer` query and display workspace/user information.

**FR-1.6**: Profile listing
Display all configured profiles with workspace names and masked API key prefixes.

### FR-2: Configuration Management

**FR-2.1**: Profile-based configuration
Store per-profile settings in `~/.config/lirt/config` using AWS CLI-style INI format.

**FR-2.2**: Configuration keys
Support: `workspace` (display name), `team` (default team key), `format` (output format), `cache_ttl` (cache lifetime), `page_size` (pagination limit).

**FR-2.3**: Configuration resolution
Resolve config values with priority: command flags → environment variables → config file → built-in defaults.

**FR-2.4**: Config file overrides
Support `LIRT_CONFIG_DIR`, `LIRT_CREDENTIALS_FILE`, `LIRT_CONFIG_FILE` environment variables for custom paths.

**FR-2.5**: Config management commands
Provide `lirt config list|get|set|unset` commands for configuration management.

### FR-3: Issue Management

**FR-3.1**: List issues
`lirt issue list` with filters: team, state, assignee, label, project, priority, milestone.

**FR-3.2**: Search issues
`lirt issue search <query>` with full-text search across titles and descriptions.

**FR-3.3**: View issue details
`lirt issue view <id>` displaying all issue metadata, description, and recent comments.

**FR-3.4**: Create issues
`lirt issue create` with required title and optional: description, team, assignee, labels, priority, state, project, milestone, parent.

**FR-3.5**: Edit issues
`lirt issue edit <id>` to update any mutable field (title, description, priority, assignee, state, project, milestone).

**FR-3.6**: State transitions
Dedicated commands: `lirt issue close <id>`, `lirt issue reopen <id>`, `lirt issue transition <id> <state>`.

**FR-3.7**: Archive and delete
`lirt issue archive <id>` and `lirt issue delete <id> --confirm` with confirmation prompt.

**FR-3.8**: Label management
`lirt issue label <id> --add <name>... --remove <name>...` for batch label operations.

**FR-3.9**: Assignment management
`lirt issue assign <id> <user>` and `lirt issue unassign <id>`.

**FR-3.10**: Issue relations
`lirt issue parent <id>` and `lirt issue children <id>` to navigate hierarchies.

**FR-3.11**: ID resolution
Accept both shorthand (e.g., `ENG-123`) and UUID identifiers for all commands.

### FR-4: Project Management

**FR-4.1**: List projects
`lirt project list` with filters: team, state, limit.

**FR-4.2**: View project details
`lirt project view <id-or-name>` with full metadata.

**FR-4.3**: Project issues
`lirt project issues <id-or-name>` to list all issues in a project.

**FR-4.4**: Project milestones
`lirt project milestones <id-or-name>` to list project milestones.

**FR-4.5**: Project members
`lirt project members <id-or-name>` to list team members on the project.

**FR-4.6**: Create and edit projects
`lirt project create` and `lirt project edit <id>` with title, description, teams, lead.

**FR-4.7**: Archive and delete projects
`lirt project archive <id>` and `lirt project delete <id> --confirm`.

### FR-5: Milestone Management

**FR-5.1**: List milestones
`lirt milestone list --project <id>` to list all milestones in a project.

**FR-5.2**: View milestone details
`lirt milestone view <id>` with metadata and progress.

**FR-5.3**: Create and edit milestones
`lirt milestone create --project <id> --title "..."` and `lirt milestone edit <id>`.

**FR-5.4**: Milestone issues
`lirt milestone issues <id>` to list all issues in a milestone.

**FR-5.5**: Delete milestones
`lirt milestone delete <id> --confirm`.

### FR-6: Initiative Management

**FR-6.1**: List initiatives
`lirt initiative list` with pagination.

**FR-6.2**: View initiative details
`lirt initiative view <id-or-name>` with metadata.

**FR-6.3**: Create and edit initiatives
`lirt initiative create --title "..."` and `lirt initiative edit <id>`.

**FR-6.4**: Initiative projects
`lirt initiative projects <id>` to list projects under an initiative.

**FR-6.5**: Archive and delete initiatives
`lirt initiative archive <id>` and `lirt initiative delete <id> --confirm`.

### FR-7: Team Management

**FR-7.1**: List teams
`lirt team list` with team key, name, member count, issue count.

**FR-7.2**: View team details
`lirt team view <key-or-id>` with full metadata.

**FR-7.3**: Team members
`lirt team members <key-or-id>` listing all team members.

**FR-7.4**: Team workflow states
`lirt team states <key-or-id>` listing workflow states with type, name, color, position.

**FR-7.5**: Team labels
`lirt team labels <key-or-id>` listing team-scoped labels.

**FR-7.6**: Team cycles
`lirt team cycles <key-or-id>` listing current, upcoming, and past cycles.

### FR-8: User Management

**FR-8.1**: List users
`lirt user list` with pagination.

**FR-8.2**: View user details
`lirt user view <id-or-login-or-email>` with profile information.

**FR-8.3**: Current user
`lirt user me` to display authenticated user information.

**FR-8.4**: User issues
`lirt user issues <id-or-login>` to list issues assigned to a user.

### FR-9: Comment Management

**FR-9.1**: List comments
`lirt comment list <issue-id>` to display all comments on an issue.

**FR-9.2**: Add comments
`lirt comment add <issue-id> --body "..."` or `--body-file <path>` for markdown content.

**FR-9.3**: Edit comments
`lirt comment edit <comment-id> --body "..."`.

**FR-9.4**: Delete comments
`lirt comment delete <comment-id> --confirm`.

**FR-9.5**: Comments on projects and initiatives
Support `--project <id>` and `--initiative <id>` flags for comments on non-issue entities.

### FR-10: Metadata Operations

**FR-10.1**: Workflow states
`lirt meta states [--team <key>]` to list available workflow states.

**FR-10.2**: Priorities
`lirt meta priorities` to list priority levels (0=Urgent through 4=None).

**FR-10.3**: Labels
`lirt meta labels [--team <key>]` to list available labels with colors.

**FR-10.4**: Cycles
`lirt meta cycles [--team <key>]` to list cycles with dates and states.

**FR-10.5**: Issue types
`lirt meta issue-types` to list available issue types if custom types are enabled.

### FR-11: Raw API Access

**FR-11.1**: GraphQL query execution
`lirt api <query>` to execute arbitrary GraphQL queries.

**FR-11.2**: Query from file
`lirt api --input <file.graphql>` or `--input -` for stdin.

**FR-11.3**: Query variables
`lirt api -f field=value <query>` to pass variables to queries.

**FR-11.4**: JSON output
All `lirt api` commands output JSON for downstream processing.

### FR-12: Output Formats

**FR-12.1**: Table format
Human-readable aligned columns with headers (default for terminal output).

**FR-12.2**: JSON format
Full JSON output with `--format json` or field selection with `--json <fields>`.

**FR-12.3**: CSV format
Comma-separated values with headers for spreadsheet import.

**FR-12.4**: Plain format
Minimal output, one value per line, for piping to other commands.

**FR-12.5**: jq integration
`--jq <expression>` flag to apply jq expressions to JSON output inline.

**FR-12.6**: Automatic format detection
Default to JSON when stdout is not a terminal (piped).

### FR-13: Caching

**FR-13.1**: Cache enumeration data
Cache teams, states, labels, users, and priorities with configurable TTL.

**FR-13.2**: Cache invalidation
Automatically invalidate cache on write operations.

**FR-13.3**: Manual cache bypass
`--no-cache` flag to force fresh API requests.

**FR-13.4**: Cache configuration
`cache_ttl` config key to control cache lifetime (default: 5 minutes).

**FR-13.5**: Cache location
Store cache in `~/.config/lirt/cache/<profile>/` with `0700` permissions.

### FR-14: Pagination

**FR-14.1**: Limit control
`--limit <n>` flag to control page size (default: 50).

**FR-14.2**: Fetch all pages
`--all` flag to recursively fetch all pages.

**FR-14.3**: Manual pagination
`--cursor <cursor>` flag to resume from a specific cursor.

**FR-14.4**: Cursor pagination
Use Relay-style cursor pagination with `pageInfo.hasNextPage` and `endCursor`.

**FR-14.5**: Streaming output
Stream results with `--all` to enable immediate downstream processing.

### FR-15: Shell Completions

**FR-15.1**: Bash completions
`lirt completion bash` to generate bash completion script.

**FR-15.2**: Zsh completions
`lirt completion zsh` to generate zsh completion script.

**FR-15.3**: Fish completions
`lirt completion fish` to generate fish completion script.

**FR-15.4**: Completion generation
Use Cobra's built-in completion generation.

---

## 5. Non-Functional Requirements

### NFR-1: Performance

**NFR-1.1**: Startup time < 50ms
Tool must start and begin execution in under 50 milliseconds for script-friendly usage.

**NFR-1.2**: Memory usage < 50MB
Baseline memory usage (before large data operations) must stay under 50MB.

**NFR-1.3**: Response time
Cached operations should respond in < 10ms. Uncached API calls should complete within Linear's API response time + minimal overhead.

**NFR-1.4**: Binary size
Compiled binary size should be < 20MB for easy distribution.

### NFR-2: Compatibility

**NFR-2.1**: Go version
Require Go 1.21 or later for development and building.

**NFR-2.2**: Platform support
Support macOS (arm64, amd64), Linux (amd64, arm64).

**NFR-2.3**: Linear API version
Track Linear GraphQL API version and document compatibility.

**NFR-2.4**: No CGo
Pure Go implementation with no CGo dependencies for easy cross-compilation.

### NFR-3: Usability

**NFR-3.1**: Consistent command structure
All commands follow `lirt <resource> <action>` pattern.

**NFR-3.2**: Help text
Every command and flag must have clear help text accessible via `--help`.

**NFR-3.3**: Error messages
Errors must include suggested corrective actions.

**NFR-3.4**: Sensible defaults
Provide reasonable defaults for all optional parameters.

### NFR-4: Reliability

**NFR-4.1**: Error handling
Gracefully handle network errors, API errors, rate limiting, and partial success responses.

**NFR-4.2**: Retry logic
Automatically retry rate-limited requests after `Retry-After` delay (max 3 retries).

**NFR-4.3**: Data validation
Validate user input before making API calls.

**NFR-4.4**: Exit codes
Use appropriate exit codes (0=success, 1=error, 2=usage error, 3=auth error, 4=not found).

### NFR-5: Security

**NFR-5.1**: Secure credential storage
Credentials file must use `0600` permissions (owner read/write only).

**NFR-5.2**: No API keys in logs
Never log or display full API keys. Show only masked prefixes (e.g., `lin_api_xxxx...`).

**NFR-5.3**: No API keys in error messages
Ensure API keys are not leaked in error messages or stack traces.

**NFR-5.4**: HTTPS only
All API communication over HTTPS.

### NFR-6: Maintainability

**NFR-6.1**: Idiomatic Go
Follow Go best practices and idioms (effective Go guidelines).

**NFR-6.2**: Code organization
Clear separation of concerns: commands, client, config, output, models.

**NFR-6.3**: Documentation
All exported functions and types must have godoc comments.

**NFR-6.4**: Testing
Maintain > 70% test coverage with unit, integration, and E2E tests.

---

## 6. Technical Constraints

### TC-1: Linear API Limitations

**TC-1.1**: Single global endpoint
Linear uses a single GraphQL endpoint (`https://api.linear.app/graphql`) for all workspaces. The API key determines the workspace.

**TC-1.2**: No sandbox API
Linear does not provide a sandbox or staging API environment. Testing must use a real workspace.

**TC-1.3**: Rate limiting
Linear enforces rate limits (1,500 requests/hour per API key on free plan). Tool must respect `Retry-After` headers.

**TC-1.4**: GraphQL-only
Linear has no REST API. All operations use GraphQL queries and mutations.

**TC-1.5**: Cursor pagination
Linear uses Relay-style cursor pagination. Offset pagination is not supported.

**TC-1.6**: Authentication
Personal API Keys only. OAuth 2.0 requires manual app registration and user consent flow.

### TC-2: Implementation Constraints

**TC-2.1**: CLI framework
Use Cobra for command structure and flag parsing.

**TC-2.2**: Config management
Use Viper for configuration file parsing and environment variable binding.

**TC-2.3**: GraphQL client
Use existing Go GraphQL client library (e.g., `github.com/hasura/go-graphql-client`).

**TC-2.4**: Table rendering
Use tablewriter library for aligned table output.

**TC-2.5**: INI parsing
Use `gopkg.in/ini.v1` for AWS-style INI credential files.

### TC-3: Testing Constraints

**TC-3.1**: Pre-provisioned workspace
Integration tests require a manually created Linear free workspace with a dedicated API key.

**TC-3.2**: CI secrets
CI pipeline must have `LIRT_TEST_API_KEY` and `LIRT_TEST_TEAM` environment variables configured.

**TC-3.3**: Active issue limit
Free Linear workspace has 250 active issue limit. Tests must clean up or use archived issues.

**TC-3.4**: Rate limit awareness
Tests must respect rate limits and include delays between API calls in test loops.

---

## 7. Dependencies

### Core Libraries

| Dependency | Purpose | License |
|------------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | Apache 2.0 |
| `github.com/spf13/viper` | Config management | MIT |
| `github.com/hasura/go-graphql-client` | GraphQL client | MIT |
| `gopkg.in/ini.v1` | INI file parsing | Apache 2.0 |
| `github.com/olekukonez/tablewriter` | Table rendering | MIT |
| `github.com/fatih/color` | Terminal colors | MIT |

### Development Dependencies

| Dependency | Purpose | License |
|------------|---------|---------|
| `github.com/stretchr/testify` | Testing assertions | MIT |
| `golang.org/x/tools/cmd/goimports` | Import formatting | BSD |
| `github.com/golangci/golangci-lint` | Linting | GPL 3.0 |

### External Services

| Service | Purpose | Required |
|---------|---------|----------|
| Linear GraphQL API | Issue tracking operations | Yes |
| Linear Free Workspace | Integration testing | Yes (CI only) |

---

## 8. Security Considerations

### SC-1: Credential Protection

**SC-1.1**: File permissions
Credentials file (`~/.config/lirt/credentials`) must be created with `0600` permissions.

**SC-1.2**: Credential validation
Validate credential file permissions on every read. Warn if permissions are too permissive.

**SC-1.3**: Environment variable fallback
Support `LIRT_API_KEY` and `LINEAR_API_KEY` environment variables for CI/CD workflows.

**SC-1.4**: No credential logging
Never log API keys in verbose output, errors, or debug messages.

### SC-2: API Key Exposure Prevention

**SC-2.1**: Masked display
Display only first 8 characters + ellipsis when showing API keys (e.g., `lin_api_xxxx...`).

**SC-2.2**: Error sanitization
Strip API keys from error messages and stack traces before display.

**SC-2.3**: Process listing
API keys passed via `--api-key` flag are visible in `ps` output. Document this risk.

### SC-3: Network Security

**SC-3.1**: HTTPS only
All Linear API communication must use HTTPS.

**SC-3.2**: Certificate validation
Validate SSL/TLS certificates. Do not allow insecure connections.

**SC-3.3**: Timeout configuration
Enforce reasonable timeouts to prevent hung connections.

### SC-4: Input Validation

**SC-4.1**: Command injection prevention
Sanitize all user input before passing to shell commands or GraphQL queries.

**SC-4.2**: Path traversal prevention
Validate file paths for `--body-file` and config file operations.

**SC-4.3**: GraphQL injection prevention
Use parameterized GraphQL queries. Never interpolate user input directly into query strings.

### SC-5: Data Privacy

**SC-5.1**: Cache security
Cache directory (`~/.config/lirt/cache/`) must have `0700` permissions.

**SC-5.2**: Sensitive data in cache
Cache contains issue titles and metadata. Document this in security documentation.

**SC-5.3**: Cache cleanup
Provide mechanism to clear cache for security-conscious users.

---

## Acceptance Criteria

This requirements document is considered complete and ready for implementation when:

1. All stakeholders have reviewed and approved the requirements
2. All functional requirements (FR-*) are clear and testable
3. All non-functional requirements (NFR-*) have measurable acceptance criteria
4. All technical constraints (TC-*) are documented and understood
5. Security considerations (SC-*) are reviewed by security team
6. Dependencies are identified and license-compatible
7. Testing strategy addresses the "no sandbox API" constraint

---

## Related Documentation

- [SPECIFICATION.md](./SPECIFICATION.md) - Technical specification
- [COMMANDS.md](./COMMANDS.md) - CLI command reference
- [AUTHENTICATION.md](./AUTHENTICATION.md) - Authentication setup guide
- [CONFIGURATION.md](./CONFIGURATION.md) - Configuration file reference

---

**Source**: Based on research conducted in the clio project (cl-k0a1) on 2026-01-27.
