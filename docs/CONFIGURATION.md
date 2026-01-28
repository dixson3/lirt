# Configuration Guide

**lirt** configuration reference for managing settings, profiles, and preferences.

---

## Table of Contents

1. [Overview](#overview)
2. [Configuration Files](#configuration-files)
3. [Configuration Keys](#configuration-keys)
4. [Profile Management](#profile-management)
5. [Environment Variables](#environment-variables)
6. [Resolution Priority](#resolution-priority)
7. [Examples](#examples)
8. [Troubleshooting](#troubleshooting)

---

## Overview

lirt uses a two-file configuration system modeled on the AWS CLI:

- **`credentials`**: Stores API keys (sensitive, permissions `0600`)
- **`config`**: Stores settings and preferences (non-sensitive, permissions `0644`)

Both files use INI format and support multiple named profiles for managing different Linear workspaces.

**Default Location**: `~/.config/lirt/`

---

## Configuration Files

### Directory Structure

```
~/.config/lirt/
├── credentials       # API keys (0600 permissions)
├── config            # Settings (0644 permissions)
└── cache/            # Cached enumeration data (0700)
    ├── default/
    │   ├── teams.json
    │   ├── states.json
    │   ├── labels.json
    │   ├── users.json
    │   └── priorities.json
    └── work/
        └── ...
```

### credentials File

**Format**: INI with profile names in brackets (no `profile` prefix)

**Permissions**: `0600` (owner read/write only)

**Contents**:
```ini
[default]
api_key = lin_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

[work]
api_key = lin_api_yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy

[personal]
api_key = lin_api_zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
```

**Security**:
- lirt automatically sets `0600` permissions when creating this file
- Never commit this file to version control
- Add to `.gitignore`: `.config/lirt/credentials`

### config File

**Format**: INI with `[profile <name>]` prefix (except `[default]`)

**Permissions**: `0644` (public readable)

**Contents**:
```ini
[default]
workspace = Acme Corporation
team = ENG
format = table
cache_ttl = 5m
page_size = 50

[profile work]
workspace = TechCorp Engineering
team = BACKEND
format = json
cache_ttl = 10m
page_size = 100

[profile personal]
workspace = Alice's Projects
team = TODO
format = table
cache_ttl = 5m
page_size = 50
```

**Note**: The `workspace` key is **display-only** and automatically populated by `lirt auth login`. Profile names are user-chosen aliases and don't need to match Linear workspace names.

---

## Configuration Keys

### Supported Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `workspace` | string | (auto) | Display name of Linear workspace (set by `lirt auth login`) |
| `team` | string | (none) | Default team key to use when `--team` flag is omitted |
| `format` | string | `table` | Default output format: `table`, `json`, `csv`, `plain` |
| `cache_ttl` | duration | `5m` | Cache lifetime for enumeration data (teams, states, labels, users) |
| `page_size` | int | `50` | Default pagination limit for list commands |

### Key Details

#### `workspace`

**Purpose**: Human-readable workspace name for display purposes

**Set by**: `lirt auth login` (via Linear's `viewer` query)

**Example**:
```ini
[default]
workspace = Acme Corporation
```

**Display**:
```bash
lirt auth status
# Workspace: Acme Corporation
```

**Cannot be set manually** — always auto-populated from Linear API.

#### `team`

**Purpose**: Default team key to avoid repeating `--team` flag

**Format**: Team key (e.g., `ENG`, `PRODUCT`, `DESIGN`)

**Example**:
```ini
[default]
team = ENG
```

**Usage**:
```bash
# Without team config
lirt issue list --team ENG

# With team config (team = ENG)
lirt issue list  # Uses ENG team automatically
```

**Override**:
```bash
# Use different team for single command
lirt issue list --team DESIGN
```

#### `format`

**Purpose**: Default output format

**Values**: `table`, `json`, `csv`, `plain`

**Example**:
```ini
[default]
format = table  # Human-readable tables

[profile ci]
format = json   # Machine-readable JSON for scripts
```

**Usage**:
```bash
# Uses format from config
lirt issue list

# Override for single command
lirt issue list --format json
```

**Auto-detection**: When stdout is not a terminal (piped), lirt defaults to `json` regardless of config.

#### `cache_ttl`

**Purpose**: How long to cache enumeration data

**Format**: Duration string (e.g., `5m`, `1h`, `30s`)

**Default**: `5m` (5 minutes)

**Example**:
```ini
[default]
cache_ttl = 5m

[profile ci]
cache_ttl = 0  # Disable caching in CI
```

**Cached Data**:
- Teams (team keys, names, member counts)
- Workflow states (state names, types, colors)
- Labels (label names, colors, scopes)
- Users (user IDs, names, emails)
- Priorities (priority levels)

**Bypass Cache**:
```bash
# Force fresh API request
lirt team list --no-cache
```

#### `page_size`

**Purpose**: Default number of results per page

**Format**: Integer (1-100, Linear's max is 100)

**Default**: `50`

**Example**:
```ini
[default]
page_size = 50

[profile bulk]
page_size = 100  # Maximum for faster bulk operations
```

**Usage**:
```bash
# Uses page_size from config
lirt issue list

# Override for single command
lirt issue list --limit 10

# Fetch all pages regardless of page_size
lirt issue list --all
```

---

## Profile Management

### Creating Profiles

Profiles are created automatically when you run `lirt auth login --profile <name>`:

```bash
# Create 'work' profile
lirt auth login --profile work
# Enter API key...
# ✓ Saved to profile 'work'
```

This creates entries in both `credentials` and `config` files.

### Editing Profiles

**Option 1: Using lirt commands** (recommended)

```bash
# Set config value
lirt config set team BACKEND --profile work

# Get config value
lirt config get team --profile work

# Unset config value
lirt config unset team --profile work

# List all config for profile
lirt config list --profile work
```

**Option 2: Direct file editing**

```bash
# Edit config file
$EDITOR ~/.config/lirt/config

# Edit credentials file (careful!)
$EDITOR ~/.config/lirt/credentials
```

### Deleting Profiles

```bash
# Remove profile from both files
lirt auth logout --profile old-profile

# Manually delete from files
$EDITOR ~/.config/lirt/credentials  # Remove [old-profile] section
$EDITOR ~/.config/lirt/config       # Remove [profile old-profile] section
```

---

## Environment Variables

### Configuration Location Overrides

| Variable | Purpose | Default |
|----------|---------|---------|
| `LIRT_CONFIG_DIR` | Override entire config directory | `~/.config/lirt/` |
| `LIRT_CREDENTIALS_FILE` | Override credentials file path | `$LIRT_CONFIG_DIR/credentials` |
| `LIRT_CONFIG_FILE` | Override config file path | `$LIRT_CONFIG_DIR/config` |

**Examples**:

```bash
# Use custom config directory
export LIRT_CONFIG_DIR=~/projects/lirt-config
lirt auth login

# Custom credentials file only
export LIRT_CREDENTIALS_FILE=/secure/linear-keys
lirt auth status

# Custom config file only
export LIRT_CONFIG_FILE=~/lirt-settings.ini
lirt issue list
```

### Profile Selection

| Variable | Purpose | Priority |
|----------|---------|----------|
| `LIRT_PROFILE` | Select active profile | Medium (overrides `[default]`, overridden by `--profile` flag) |

**Examples**:

```bash
# Use 'work' profile for all commands
export LIRT_PROFILE=work
lirt issue list
lirt team list

# Override with flag
lirt issue list --profile personal
```

### Credential Overrides

| Variable | Purpose | Priority |
|----------|---------|----------|
| `LIRT_API_KEY` | Override API key | Highest (always wins) |
| `LINEAR_API_KEY` | Fallback API key | Lowest (only if no other source) |

**Examples**:

```bash
# Temporary credential override
LIRT_API_KEY=lin_api_temp... lirt issue list

# CI/CD environment
export LIRT_API_KEY=${{ secrets.LINEAR_API_KEY }}
lirt issue list --team ENG
```

### Config Value Overrides

| Variable | Purpose | Example |
|----------|---------|---------|
| `LIRT_TEAM` | Override default team | `export LIRT_TEAM=ENG` |
| `LIRT_FORMAT` | Override output format | `export LIRT_FORMAT=json` |
| `LIRT_CACHE_TTL` | Override cache TTL | `export LIRT_CACHE_TTL=10m` |
| `LIRT_PAGE_SIZE` | Override page size | `export LIRT_PAGE_SIZE=100` |

---

## Resolution Priority

lirt resolves configuration values with the following priority (highest to lowest):

### 1. Command-line Flags

**Highest priority** — always wins

```bash
lirt issue list --team ENG --format json
# Uses: team=ENG, format=json (ignores config and env vars)
```

### 2. Environment Variables

**Second priority** — overrides config file

```bash
export LIRT_TEAM=DESIGN
export LIRT_FORMAT=csv

lirt issue list
# Uses: team=DESIGN, format=csv (ignores config file)
```

### 3. Config File (Selected Profile)

**Third priority** — default source

Profile selection priority:
1. `--profile` flag
2. `LIRT_PROFILE` env var
3. `[default]` profile

```ini
[default]
team = ENG
format = table

[profile work]
team = BACKEND
format = json
```

```bash
# Uses [default] profile
lirt issue list
# team=ENG, format=table

# Uses [profile work]
lirt issue list --profile work
# team=BACKEND, format=json
```

### 4. Built-in Defaults

**Lowest priority** — fallback if nothing else is set

| Key | Default |
|-----|---------|
| `team` | (none) |
| `format` | `table` |
| `cache_ttl` | `5m` |
| `page_size` | `50` |

### Complete Resolution Example

**Setup**:
```ini
# ~/.config/lirt/config
[default]
team = ENG
format = table
page_size = 50
```

```bash
export LIRT_TEAM=DESIGN
export LIRT_FORMAT=json
```

**Command**:
```bash
lirt issue list --format csv --limit 10
```

**Resolution**:
- `team` = `DESIGN` (from `LIRT_TEAM` env var)
- `format` = `csv` (from `--format` flag, highest priority)
- `limit` = `10` (from `--limit` flag)

---

## Examples

### Example 1: Basic Configuration

**Scenario**: Single workspace, prefer table output

```ini
# ~/.config/lirt/config
[default]
workspace = Acme Corporation
team = ENG
format = table
cache_ttl = 5m
page_size = 50
```

**Usage**:
```bash
# Uses all defaults from config
lirt issue list

# Override format for JSON export
lirt issue list --format json > issues.json
```

### Example 2: Multi-Workspace Setup

**Scenario**: Work and personal projects with different defaults

```ini
# ~/.config/lirt/config
[default]
workspace = Acme Corporation
team = ENG
format = table
page_size = 50

[profile personal]
workspace = Alice's Projects
team = TODO
format = table
page_size = 25

[profile client]
workspace = Client XYZ
team = DEV
format = json
page_size = 100
```

**Usage**:
```bash
# Work issues (default profile)
lirt issue list

# Personal issues
lirt issue list --profile personal

# Client issues
lirt issue list --profile client
```

### Example 3: CI/CD Configuration

**Scenario**: GitHub Actions workflow needing JSON output

```ini
# ~/.config/lirt/config
[profile ci]
workspace = CI Test Workspace
team = TEST
format = json
cache_ttl = 0
page_size = 100
```

**GitHub Actions**:
```yaml
- name: List open issues
  env:
    LIRT_API_KEY: ${{ secrets.LINEAR_API_KEY }}
    LIRT_PROFILE: ci
  run: lirt issue list --state open
```

### Example 4: Script-Friendly Setup

**Scenario**: Bash script processing issues

```bash
#!/bin/bash
# scripts/close-done-issues.sh

# Use JSON format for all commands
export LIRT_FORMAT=json
export LIRT_TEAM=ENG

# Get done issue IDs
done_issues=$(lirt issue list --state done | jq -r '.[].id')

# Close each issue
for id in $done_issues; do
  lirt issue close "$id"
  echo "Closed $id"
done
```

### Example 5: Developer-Specific Preferences

**Scenario**: Different developers prefer different defaults

```bash
# Developer A: Prefers tables
lirt config set format table
lirt config set page_size 50

# Developer B: Prefers JSON for piping to jq
lirt config set format json
lirt config set page_size 100
```

---

## Troubleshooting

### Config Not Loading

**Problem**: Settings from config file are ignored

**Check**:
```bash
# Verify config file location
echo $LIRT_CONFIG_FILE
# Should be empty or custom path

# Check if file exists
cat ~/.config/lirt/config

# Verify config loading
lirt config list
```

**Solutions**:
```bash
# Recreate config directory
mkdir -p ~/.config/lirt

# Re-login to regenerate config
lirt auth login
```

### Wrong Profile Selected

**Problem**: Commands use wrong workspace

**Check**:
```bash
# See which profile is active
lirt auth status

# List all profiles
lirt auth list

# Check environment variable
echo $LIRT_PROFILE
```

**Solutions**:
```bash
# Use correct profile
lirt issue list --profile correct-name

# Set environment variable
export LIRT_PROFILE=correct-name

# Switch default profile
eval "$(lirt auth switch correct-name)"
```

### Cache Issues

**Problem**: Seeing stale data (old team names, states, etc.)

**Solutions**:
```bash
# Bypass cache for single command
lirt team list --no-cache

# Disable caching for profile
lirt config set cache_ttl 0 --profile myprofile

# Clear cache manually
rm -rf ~/.config/lirt/cache/
```

### Permission Errors

**Problem**: "permission denied" on config files

**Check**:
```bash
ls -l ~/.config/lirt/
```

**Solutions**:
```bash
# Fix credentials file (must be 0600)
chmod 600 ~/.config/lirt/credentials

# Fix config file
chmod 644 ~/.config/lirt/config

# Fix cache directory
chmod 700 ~/.config/lirt/cache
```

### Config File Corruption

**Problem**: Syntax errors in INI files

**Check**:
```bash
# Validate config
lirt config list

# Check for INI syntax errors
cat ~/.config/lirt/config
```

**Solutions**:
```bash
# Backup existing config
cp ~/.config/lirt/config ~/.config/lirt/config.bak

# Recreate config
rm ~/.config/lirt/config
lirt auth login  # Regenerates config
```

---

## Advanced Configuration

### Custom Config Location

**Scenario**: Store lirt config in project directory

```bash
# Set custom config directory
export LIRT_CONFIG_DIR=/path/to/project/.lirt

# All lirt commands use custom config
lirt auth login
lirt issue list

# Files created:
# /path/to/project/.lirt/credentials
# /path/to/project/.lirt/config
# /path/to/project/.lirt/cache/
```

### Per-Project Configuration

**Scenario**: Different lirt settings per project

```bash
# Project A
cd ~/projects/project-a
export LIRT_CONFIG_DIR=.lirt
lirt config set team TEAM-A

# Project B
cd ~/projects/project-b
export LIRT_CONFIG_DIR=.lirt
lirt config set team TEAM-B
```

**Add to `.env` or `.envrc`** (with direnv):
```bash
# .envrc
export LIRT_CONFIG_DIR=.lirt
export LIRT_PROFILE=project-a
```

### Shared Team Configuration

**Scenario**: Team shares same lirt config via git

```bash
# Create shared config
mkdir -p team-config
cat > team-config/config <<EOF
[default]
team = SHARED-TEAM
format = json
page_size = 100
EOF

# Team members use shared config
export LIRT_CONFIG_FILE=~/team-config/config
lirt issue list
```

**Note**: Never commit the `credentials` file (contains API keys).

---

## Related Documentation

- [AUTHENTICATION.md](./AUTHENTICATION.md) - Credential management
- [SPECIFICATION.md](./SPECIFICATION.md) - Technical specification
- [COMMANDS.md](./COMMANDS.md) - CLI command reference

---

**Source**: Based on research conducted in the clio project (cl-k0a1) on 2026-01-27.
