# Testing Guide for lirt

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

## API Tests

Some tests require access to the Linear API and will be skipped unless you provide a test API key.

### Setting up API test credentials

**Option 1: Environment variable (recommended for CI)**
```bash
export LINEAR_TEST_API_KEY="your_test_api_key_here"
go test ./internal/client -v
```

**Option 2: Create a test profile**
```bash
# Login with a test account (creates test workspace credentials)
./bin/lirt auth login --profile test

# Run tests using that profile
LINEAR_TEST_API_KEY=$(./bin/lirt auth token --profile test) go test ./internal/client -v
```

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
