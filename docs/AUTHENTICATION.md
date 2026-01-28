# Authentication Guide

**lirt** authentication guide for managing Linear API credentials and multi-workspace access.

---

## Table of Contents

1. [Overview](#overview)
2. [Getting a Linear API Key](#getting-a-linear-api-key)
3. [Initial Setup](#initial-setup)
4. [Multi-Profile Management](#multi-profile-management)
5. [Credential Resolution](#credential-resolution)
6. [Security Best Practices](#security-best-practices)
7. [Troubleshooting](#troubleshooting)

---

## Overview

`lirt` uses Linear Personal API Keys for authentication. Each API key is scoped to exactly one Linear workspace, and lirt supports managing multiple workspaces through named profiles.

**Key Concepts:**
- **Profile**: A named configuration (e.g., `default`, `work`, `personal`) containing an API key and settings
- **Workspace**: A Linear organization/workspace determined by the API key
- **Credentials File**: `~/.config/lirt/credentials` storing API keys per profile

---

## Getting a Linear API Key

### Step 1: Access Linear Settings

1. Log in to your Linear workspace at `https://linear.app`
2. Click your profile icon (bottom left)
3. Select **Settings**
4. Navigate to **API** section

### Step 2: Generate Personal API Key

1. Click **Create API Key** or **Personal API Keys**
2. Enter a descriptive name (e.g., "lirt CLI tool")
3. Set permissions:
   - **Read**: View issues, projects, teams
   - **Write**: Create and modify issues
   - **Admin**: Full workspace access (optional)
4. Click **Create**
5. Copy the API key immediately — it won't be shown again

The key format: `lin_api_` followed by 32 alphanumeric characters

**Example**: `lin_api_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6`

---

## Initial Setup

### Interactive Login (Recommended)

```bash
# First-time setup
lirt auth login

# You'll be prompted:
# Enter your Linear API key: [input hidden]
# Validating...
# ✓ Authenticated as Alice Developer (alice@acme.com)
# ✓ Workspace: Acme Corporation
# ✓ Saved to profile 'default'
```

### Non-Interactive Login

```bash
# For scripts or CI/CD
lirt auth login --api-key lin_api_xxxxxxxxxxxxxxxxx
```

### Named Profile Login

```bash
# Set up a secondary workspace
lirt auth login --profile work

# Or non-interactive
lirt auth login --profile work --api-key lin_api_xxxxxxxxxxxxxxxxx
```

### Verify Authentication

```bash
# Check authentication status
lirt auth status

# Output:
# Profile:    default
# Workspace:  Acme Corporation
# User:       Alice Developer (alice@acme.com)
# Permissions: Read, Write
# Key prefix: lin_api_xxxx...
```

---

## Multi-Profile Management

### Use Case: Multiple Workspaces

Many developers work with multiple Linear workspaces:
- Personal projects
- Work projects
- Open source projects
- Client projects

lirt's profile system (modeled on AWS CLI) makes this simple.

### Setting Up Multiple Profiles

```bash
# Profile 1: Personal projects
lirt auth login --profile personal
# Enter API key from your personal Linear workspace

# Profile 2: Work projects
lirt auth login --profile work
# Enter API key from your company's Linear workspace

# Profile 3: Open source
lirt auth login --profile oss
# Enter API key from open source project Linear workspace
```

### Listing All Profiles

```bash
lirt auth list

# Output:
# PROFILE    WORKSPACE              KEY PREFIX
# default    Acme Corporation       lin_api_xxxx...
# personal   Alice's Projects       lin_api_yyyy...
# work       TechCorp Engineering   lin_api_zzzz...
```

### Using a Specific Profile

**Option 1: Per-command flag**
```bash
lirt issue list --profile work
lirt team list --profile personal
```

**Option 2: Environment variable**
```bash
export LIRT_PROFILE=work
lirt issue list  # Uses 'work' profile
lirt team list   # Uses 'work' profile
```

**Option 3: Shell alias with profile switching**
```bash
# Add to ~/.bashrc or ~/.zshrc
alias lirt-work='lirt --profile work'
alias lirt-personal='lirt --profile personal'

# Usage
lirt-work issue list
lirt-personal team list
```

### Switching Default Profile

```bash
# Generate export command (copy-paste to shell)
lirt auth switch work
# Output: export LIRT_PROFILE=work

# Or evaluate directly (bash/zsh)
eval "$(lirt auth switch work)"

# Now all lirt commands use 'work' profile by default
lirt issue list
```

### Removing a Profile

```bash
# Delete profile from credentials and config files
lirt auth logout --profile personal

# Confirm removal:
# Remove profile 'personal'? (y/N): y
# ✓ Profile 'personal' removed
```

---

## Credential Resolution

lirt resolves credentials with the following priority (highest to lowest):

### 1. `LIRT_API_KEY` Environment Variable

**Highest priority** — always wins if set. Useful for:
- CI/CD pipelines
- Temporary credential override
- Testing

```bash
# Override any profile
LIRT_API_KEY=lin_api_test123... lirt issue list
```

### 2. `--api-key` Flag

Per-command API key override:

```bash
lirt issue list --api-key lin_api_temp456...
```

### 3. Credentials File with Profile Selection

Default credential source. Profile selection priority:

1. `--profile <name>` flag
2. `LIRT_PROFILE` environment variable
3. `[default]` profile in credentials file

```bash
# Uses 'work' profile
lirt issue list --profile work

# Uses profile from environment
export LIRT_PROFILE=personal
lirt issue list

# Uses 'default' profile (no flag or env var)
lirt issue list
```

### 4. `LINEAR_API_KEY` Environment Variable

**Lowest priority** — fallback for compatibility with Linear's official tools:

```bash
export LINEAR_API_KEY=lin_api_fallback...
lirt issue list
```

### Resolution Examples

**Example 1**: Multiple credential sources
```bash
# credentials file has 'default' profile
# LIRT_PROFILE=work
# LIRT_API_KEY=lin_api_ci...

lirt issue list
# Uses: LIRT_API_KEY (highest priority)
```

**Example 2**: Profile flag override
```bash
# LIRT_PROFILE=work

lirt issue list --profile personal
# Uses: 'personal' profile (flag overrides env var)
```

**Example 3**: CI/CD workflow
```bash
# In CI environment, no credentials file exists
# Set LIRT_API_KEY from CI secrets

export LIRT_API_KEY=${{ secrets.LINEAR_API_KEY }}
lirt issue list --team ENG --state open
```

---

## Security Best Practices

### File Permissions

lirt automatically sets secure permissions on credential files:

```bash
# Credentials file (API keys)
~/.config/lirt/credentials → 0600 (owner read/write only)

# Config file (settings, no secrets)
~/.config/lirt/config → 0644 (owner read/write, others read)

# Cache directory
~/.config/lirt/cache/ → 0700 (owner only)
```

**Verify permissions:**
```bash
ls -l ~/.config/lirt/
# -rw------- credentials
# -rw-r--r-- config
# drwx------ cache/
```

### API Key Security

**DO:**
- ✅ Store API keys in `~/.config/lirt/credentials` (secure permissions)
- ✅ Use environment variables for CI/CD (`LIRT_API_KEY`)
- ✅ Rotate API keys periodically
- ✅ Use separate API keys per tool/purpose
- ✅ Delete unused profiles: `lirt auth logout --profile old`

**DON'T:**
- ❌ Commit credentials file to git
- ❌ Share API keys via email/Slack
- ❌ Use `--api-key` flag in scripts (visible in `ps` output)
- ❌ Log API keys to files or stderr
- ❌ Store API keys in shell history

### Revoking Compromised Keys

If an API key is compromised:

1. **Revoke in Linear**:
   - Go to Linear Settings → API
   - Find the compromised key
   - Click **Revoke**

2. **Remove from lirt**:
   ```bash
   lirt auth logout --profile compromised
   ```

3. **Generate new key**:
   - Create new API key in Linear Settings
   - Add to lirt: `lirt auth login --profile new`

### CI/CD Best Practices

**GitHub Actions Example**:
```yaml
name: Linear Integration
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup lirt
        run: |
          curl -sSL https://github.com/dixson3/lirt/releases/download/v0.1.0/lirt-linux-amd64 -o lirt
          chmod +x lirt
      - name: List open issues
        env:
          LIRT_API_KEY: ${{ secrets.LINEAR_API_KEY }}
        run: ./lirt issue list --team ENG --state open
```

**Key points**:
- Store API key in repository secrets (`secrets.LINEAR_API_KEY`)
- Pass via environment variable (not command flag)
- Use dedicated Linear workspace for CI tests

---

## Troubleshooting

### Error: "not authenticated"

**Problem**: lirt cannot find valid credentials

**Solutions**:
```bash
# Check if credentials exist
lirt auth status

# If no credentials:
lirt auth login

# If wrong profile selected:
lirt auth status --profile work
lirt auth list  # See all profiles

# Check environment variables
echo $LIRT_API_KEY
echo $LIRT_PROFILE
```

### Error: "invalid API key"

**Problem**: API key is malformed or revoked

**Solutions**:
```bash
# Verify key format (should start with lin_api_)
lirt auth status

# Test key manually
lirt api '{ viewer { id name email } }'

# If invalid, generate new key in Linear Settings and re-login
lirt auth login --profile default
```

### Error: "permission denied: ~/.config/lirt/credentials"

**Problem**: Credentials file has incorrect permissions

**Solutions**:
```bash
# Fix permissions
chmod 600 ~/.config/lirt/credentials

# Verify
ls -l ~/.config/lirt/credentials
# Should show: -rw------- (600)
```

### Error: "workspace mismatch"

**Problem**: Using wrong profile for intended workspace

**Solutions**:
```bash
# List all profiles to see workspaces
lirt auth list

# Use correct profile
lirt issue list --profile correct-profile

# Or switch default profile
eval "$(lirt auth switch correct-profile)"
```

### Can't Find `~/.config/lirt/` Directory

**Problem**: Config directory doesn't exist

**Solutions**:
```bash
# Create directory
mkdir -p ~/.config/lirt

# Run login to initialize
lirt auth login
```

### API Key Visible in Shell History

**Problem**: Used `--api-key` flag which stores key in bash history

**Solutions**:
```bash
# Clear specific entry from history
history | grep "lirt.*--api-key"
history -d <line-number>

# Or clear all history
history -c

# Better: use credentials file instead
lirt auth login
```

### Multiple lirt Installations

**Problem**: Multiple versions of lirt with different credentials

**Solutions**:
```bash
# Check which lirt is being used
which lirt

# Check version
lirt --version

# Override config directory if needed
export LIRT_CONFIG_DIR=/path/to/config
lirt auth status
```

---

## Advanced Topics

### Custom Config Directory

Override default `~/.config/lirt/`:

```bash
# Set custom directory
export LIRT_CONFIG_DIR=~/my-custom-lirt-config

# All commands use custom directory
lirt auth login
lirt issue list

# Files created in:
# ~/my-custom-lirt-config/credentials
# ~/my-custom-lirt-config/config
# ~/my-custom-lirt-config/cache/
```

### Per-File Overrides

Override specific files:

```bash
# Custom credentials file
export LIRT_CREDENTIALS_FILE=~/secure/linear-keys

# Custom config file
export LIRT_CONFIG_FILE=~/projects/lirt-config.ini

lirt auth status
```

### Reading API Key from File (CI)

```bash
# Store API key in file (0600 permissions)
echo "lin_api_..." > /secure/linear-key
chmod 600 /secure/linear-key

# Use in script
export LIRT_API_KEY=$(cat /secure/linear-key)
lirt issue list
```

### Extracting API Key for Other Tools

```bash
# Get current API key (masked by default)
lirt auth token

# Pipe to other tools
LINEAR_KEY=$(lirt auth token)
curl -H "Authorization: Bearer $LINEAR_KEY" https://api.linear.app/graphql
```

---

## Related Documentation

- [CONFIGURATION.md](./CONFIGURATION.md) - Profile configuration settings
- [SPECIFICATION.md](./SPECIFICATION.md) - Technical specification
- [COMMANDS.md](./COMMANDS.md) - CLI command reference

---

**Source**: Based on research conducted in the clio project (cl-k0a1) on 2026-01-27.
