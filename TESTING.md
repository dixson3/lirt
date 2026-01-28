# Testing Guide for lirt

## Quick Start

```bash
# Run all tests (skips API tests if LINEAR_TEST_API_KEY not set)
make test

# Set up test environment for API tests
cp /Users/james/gt/lirt/.env.test.example /Users/james/gt/lirt/.env.test
vim /Users/james/gt/lirt/.env.test  # Add your LINEAR_TEST_API_KEY

# Run API tests
make test-api
```

## Running Tests

### Run all tests
```bash
make test
# or
go test ./...
```

### Run tests with verbose output
```bash
go test -v ./...
```

### Run tests for a specific package
```bash
go test ./internal/config -v
go test ./internal/client -v
```

### Run a specific test
```bash
go test ./internal/config -run TestCredentialsFilePermissions -v
```

## Test Coverage

### Generate coverage report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### View coverage summary
```bash
go test -cover ./...
```

## Test Environment Setup

Tests that interact with the Linear API require credentials. You have multiple options for providing them:

### Option 1: Using .env.test Files (Recommended)

This is the recommended approach as it works across all worktrees and integrates with direnv.

**Location of .env.test files**:
- **Git root** (shared): `/Users/james/gt/lirt/.env.test`
  - Shared across all worktrees (refinery, witness, polecats)
  - Automatically loaded by all tests
- **Worktree-specific** (optional): `[worktree]/.env.test.local`
  - Overrides git root settings for a specific worktree
  - Useful if you need different credentials per worktree

**Setup steps**:

1. **Copy the template to git root**:
   ```bash
   cd /Users/james/gt/lirt
   cp .env.test.example .env.test
   ```

2. **Edit and add your test API key**:
   ```bash
   vim .env.test  # or your preferred editor
   # Set: LINEAR_TEST_API_KEY="lin_api_your_test_key_here"
   ```

3. **Run tests from any worktree**:
   ```bash
   cd refinery/rig
   make test-api

   # Or from any worktree:
   cd witness
   go test ./...
   ```

**The .env.test file is automatically loaded** by the test suite via `internal/testutil/env.go`.

**Worktree-specific overrides** (optional):
```bash
cd refinery/rig
cat > .env.test.local <<EOF
LINEAR_TEST_API_KEY="different_key_for_this_worktree"
LINEAR_TEST_TIMEOUT="60"
EOF
make test-api  # Uses worktree-specific settings
```

**File precedence** (later overrides earlier):
1. `/Users/james/gt/lirt/.env.test` (git root, lowest priority)
2. `[worktree]/.env.test` (worktree-level)
3. `[worktree]/.env.test.local` (worktree-specific, highest priority)

### Option 2: Using direnv

If you use [direnv](https://direnv.net/), create an `.envrc` at the git root:

```bash
cd /Users/james/gt/lirt
cat > .envrc <<EOF
# Load test environment
dotenv_if_exists .env.test
dotenv_if_exists .env.test.local
EOF
direnv allow
```

Now environment variables are automatically loaded when you `cd` into the project:
```bash
cd /Users/james/gt/lirt/refinery/rig
echo $LINEAR_TEST_API_KEY  # Automatically set by direnv
make test-api
```

### Option 3: Environment Variables (CI/CD)

For CI/CD pipelines or one-off testing:

```bash
export LINEAR_TEST_API_KEY="lin_api_your_test_key_here"
cd refinery/rig
make test-api
```

### Option 4: Using lirt Auth Profile

If you've already authenticated with lirt:

```bash
# Login with a test account
./bin/lirt auth login --profile test

# Extract API key for tests
export LINEAR_TEST_API_KEY=$(./bin/lirt auth token --profile test)
go test ./internal/client -v
```

## API Tests

Some tests require access to the Linear API and will be skipped unless you provide a test API key.

### Tests that require API access

The following tests will be **SKIPPED** if `LINEAR_TEST_API_KEY` is not set:

- `TestAuthWithInvalidAPIKey` - Verifies clear error messages for invalid keys
- `TestAuthWithValidAPIKey` - Control test to verify auth works correctly
- `TestAuthErrorFormat` - Verifies consistent error formatting
- `TestRateLimiting` - Rate limit handling (also skipped in `-short` mode)

### Tests that run without API access

These tests run without any credentials:

- `TestCredentialsFilePermissions` - File permission security (lirt-0zs)
- `TestConfigFilePermissions` - Config file permissions
- `TestMaskAPIKey` - API key masking security (lirt-3po)
- `TestMaskAPIKeyLength` - Masking length verification
- `TestMaskAPIKeyNoFullKeyExposure` - Security fuzz testing

## Test Organization

### Package: internal/config
- **config_test.go**: Configuration and credentials file security tests
  - File permissions (0600 for credentials, 0644 for config)
  - Umask handling
  - Security verification

### Package: internal/client
- **auth_test.go**: Authentication and API interaction tests
  - Invalid API key error handling (lirt-eid)
  - Valid API key authentication
  - Error message clarity and format

- **client_test.go**: Client utility function tests
  - API key masking (lirt-3po)
  - Security verification (no key exposure)
  - Edge cases and fuzzing

## Security Tests

Security-critical tests are marked with related bead IDs:

- **lirt-0zs**: Credentials file 0600 permissions
- **lirt-3po**: Masked API key never shows full key
- **lirt-eid**: Auth login with invalid key returns clear error

These tests verify:
1. Secrets are never exposed in logs, errors, or output
2. File permissions prevent unauthorized access
3. Error messages are helpful but don't leak sensitive data

## CI/CD Integration

### GitHub Actions example
```yaml
- name: Run tests
  run: go test ./...

- name: Run tests with API access
  env:
    LINEAR_TEST_API_KEY: ${{ secrets.LINEAR_TEST_API_KEY }}
  run: go test ./internal/client -v
```

### Make targets
```bash
make test          # Run all tests
make test-short    # Run tests excluding slow/API tests
make test-coverage # Generate coverage report
```

## Test Development Guidelines

### When adding new tests

1. **Unit tests**: Place in the same package as the code (e.g., `config_test.go` for `config.go`)
2. **Integration tests**: May require test fixtures or API access
3. **Security tests**: Mark with related bead ID in comments
4. **API tests**: Always check for `LINEAR_TEST_API_KEY` and skip if not present

### Test naming conventions

- `TestFunctionName` for basic functionality
- `TestFunctionNameEdgeCase` for edge cases
- `TestFunctionNameSecurity` for security-critical tests

### Skip conditions

```go
// Skip if no API key
testAPIKey := os.Getenv("LINEAR_TEST_API_KEY")
if testAPIKey == "" {
    t.Skip("Skipping test: LINEAR_TEST_API_KEY not set")
}

// Skip slow tests in short mode
if testing.Short() {
    t.Skip("Skipping slow test in short mode")
}
```

## Troubleshooting

### Tests fail with "permission denied"
- Ensure test has write access to temp directory
- Check that `t.TempDir()` is being used for test files

### API tests always skip
- Verify `LINEAR_TEST_API_KEY` is set: `echo $LINEAR_TEST_API_KEY`
- Ensure the key is valid: `./bin/lirt auth status`

### Coverage seems low
- Some packages (cmd/*) have minimal unit tests by design
- Focus on testing business logic (internal/*) and security-critical code
- Integration tests may not show in unit test coverage

## Future Test Additions

Potential areas for additional testing:

- [ ] GraphQL query generation tests
- [ ] Cache expiration and invalidation tests
- [ ] Output formatter tests (table, JSON, CSV)
- [ ] Command-line flag parsing tests
- [ ] Integration tests with real API (when API key available)
- [ ] Concurrent access tests (race detector)
