# Environment Setup for lirt Tests

This document explains how to set up test credentials for the lirt project.

## Overview

The lirt test suite requires Linear API credentials for integration tests. We support multiple approaches for providing credentials, all of which work seamlessly with the Gas Town multi-worktree setup.

## File Locations

```
refinery/rig/              # Project root (repository)
├── .env.test.example      # Template (committed to git)
├── .env.test              # Your actual credentials (gitignored)
├── .env.test.local        # Optional worktree-specific overrides (gitignored)
└── ...
```

Note: In the Gas Town multi-worktree setup, each worktree (`refinery/rig`, `witness`, etc.) can have its own `.env.test` file.

## Setup Methods

### Method 1: .env.test Files (Recommended)

**Best for**: Local development, works with direnv, shared across all worktrees

**Setup**:
```bash
# 1. Copy template to git root
cd .
cp .env.test.example .env.test

# 2. Edit with your test API key
vim .env.test
# Set: LINEAR_TEST_API_KEY="lin_api_your_test_key_here"

# 3. Run tests from any worktree
cd refinery/rig && make test-api
cd witness && go test ./...
cd polecats/polecat-xyz && make test
```

**How it works**:
- Tests automatically load `.env.test` from git root via `internal/testutil`
- Variables are loaded before tests run
- Works across all worktrees without per-worktree setup

**Worktree-specific overrides** (optional):
```bash
cd refinery/rig
echo 'LINEAR_TEST_API_KEY="different_key"' > .env.test.local
make test-api  # Uses worktree-specific key
```

**Precedence** (later overrides earlier):
1. `./.env.test` (project-wide)
2. `[worktree]/.env.test` (worktree-level)
3. `[worktree]/.env.test.local` (worktree-specific overrides)

### Method 2: direnv Integration

**Best for**: Users already using direnv for environment management

**Setup**:
```bash
cd .

# Create .envrc
cat > .envrc <<'EOF'
# Load test environment
dotenv_if_exists .env.test
dotenv_if_exists .env.test.local
EOF

# Allow direnv
direnv allow

# Now environment is automatically loaded when entering directory
cd refinery/rig
echo $LINEAR_TEST_API_KEY  # Automatically set
```

**How it works**:
- direnv automatically loads `.env.test` when you `cd` into the project
- Works with the same `.env.test` file as Method 1
- Environment variables available to shell commands and tests

### Method 3: Manual Environment Variables

**Best for**: CI/CD pipelines, one-off testing, explicit control

**Setup**:
```bash
export LINEAR_TEST_API_KEY="lin_api_your_test_key_here"
cd ./refinery/rig
make test-api
```

**How it works**:
- Standard environment variable approach
- Tests check for `LINEAR_TEST_API_KEY` and skip if not set
- No files needed

### Method 4: Claude Code Hook (Future)

**Best for**: Automated credential loading when Claude Code is active in a project

**Status**: Not yet implemented (requires Claude Code hook configuration)

**Planned approach**:
```bash
# Pre-tool-call hook would automatically load .env.test
# when Claude executes tools in the lirt project
```

## Getting a Test API Key

1. **Create a test workspace** in Linear (recommended):
   - Go to https://linear.app
   - Create a new workspace (free tier is fine)
   - Use this workspace exclusively for testing

2. **Generate an API key**:
   - Go to Settings → API
   - Create a Personal API key
   - Copy the key (starts with `lin_api_`)

3. **Add to .env.test**:
   ```bash
   cd .
   echo 'LINEAR_TEST_API_KEY="lin_api_your_key_here"' >> .env.test
   ```

## Security Considerations

### What's Protected

✅ `.env.test` is gitignored (won't be committed)
✅ `.env.test.local` is gitignored (won't be committed)
✅ API keys are masked in test output
✅ Credentials file has 0600 permissions

### What to Watch For

⚠️ Don't commit `.env.test` or `.env.test.local`
⚠️ Don't use production API keys for testing
⚠️ Don't share `.env.test` between developers (each creates their own)
⚠️ Don't echo API keys in scripts or logs

### Verifying Gitignore

```bash
# Check that .env.test won't be committed
cd .
git status  # Should not show .env.test as untracked

# If it shows up, check .gitignore
cat .gitignore | grep env.test
```

## Testing Without API Key

Most tests work without an API key and verify code behavior:

```bash
# These tests always run (no API key needed)
make test  # Runs all tests, skips API tests if no key

# Tests that run without API key:
- TestCredentialsFilePermissions (file security)
- TestConfigFilePermissions (config security)
- TestMaskAPIKey (key masking)
- TestMaskAPIKeyLength (masking verification)
- TestMaskAPIKeyNoFullKeyExposure (security fuzzing)
```

Tests requiring API access gracefully skip:
```
=== RUN   TestAuthWithInvalidAPIKey
    auth_test.go:21: Skipping auth test: LINEAR_TEST_API_KEY not set
--- SKIP: TestAuthWithInvalidAPIKey (0.00s)
```

## Troubleshooting

### Tests skip even though I set LINEAR_TEST_API_KEY

**Check**:
```bash
echo $LINEAR_TEST_API_KEY  # Is it set in your shell?
cd .
cat .env.test | grep LINEAR_TEST_API_KEY  # Is it in the file?
```

**Fix**:
```bash
# Reload environment
cd .
source .env.test  # Manual load
# or
direnv reload      # If using direnv
```

### .env.test not loading automatically

**Check git root location**:
```bash
cd ./refinery/rig
go test ./internal/testutil -v -run TestFindGitRoot
# Should output: Git root: .
```

**Check file exists**:
```bash
ls -la ./.env.test
```

### API tests fail with "authentication failed"

**Possible causes**:
1. Invalid or expired API key
2. API key from wrong workspace
3. Network connectivity issues

**Verify API key**:
```bash
cd ./refinery/rig
./bin/lirt auth login --api-key "$LINEAR_TEST_API_KEY"
./bin/lirt auth status  # Should show workspace info
```

## Advanced Usage

### Different keys per worktree

```bash
# Main key in git root (default for most worktrees)
echo 'LINEAR_TEST_API_KEY="main_workspace_key"' > ./.env.test

# Override in specific worktree
cd ./refinery/rig
echo 'LINEAR_TEST_API_KEY="refinery_specific_key"' > .env.test.local
```

### Temporary key for one test run

```bash
# Doesn't modify .env.test
LINEAR_TEST_API_KEY="temporary_key" go test ./internal/client -v
```

### CI/CD Integration

```yaml
# GitHub Actions example
- name: Run API tests
  env:
    LINEAR_TEST_API_KEY: ${{ secrets.LINEAR_TEST_API_KEY }}
  run: make test-api
```

## Related Documentation

- [TESTING.md](refinery/rig/TESTING.md) - Full testing guide
- [.env.test.example](.env.test.example) - Template file
- [internal/testutil/env.go](refinery/rig/internal/testutil/env.go) - Implementation

## Support

If you have issues with environment setup:
1. Check this document
2. Verify `.env.test` exists at git root
3. Check that `LINEAR_TEST_API_KEY` is set correctly
4. Verify API key works with `./bin/lirt auth login`
