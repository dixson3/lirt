package client

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// TestAuthWithInvalidAPIKey verifies that authentication with an invalid
// API key returns a clear error message.
// Related bead: lirt-eid
//
// This test requires LINEAR_TEST_API_KEY environment variable to be set.
// To run: LINEAR_TEST_API_KEY=your_test_key go test ./internal/client
func TestAuthWithInvalidAPIKey(t *testing.T) {
	// Skip if no test API key provided
	testAPIKey := os.Getenv("LINEAR_TEST_API_KEY")
	if testAPIKey == "" {
		t.Skip("Skipping auth test: LINEAR_TEST_API_KEY not set")
	}

	tests := []struct {
		name           string
		apiKey         string
		wantErr        bool
		errContains    string
		skipIfNoTestKey bool
	}{
		{
			name:        "Empty API key",
			apiKey:      "",
			wantErr:     true,
			errContains: "API key is required",
		},
		{
			name:        "Invalid format API key",
			apiKey:      "invalid_key_123",
			wantErr:     true,
			errContains: "", // Any error is acceptable
		},
		{
			name:        "Malformed API key",
			apiKey:      "lin_api_",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "Completely bogus API key",
			apiKey:      "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			wantErr:     true,
			errContains: "",
		},
		{
			name:        "SQL injection attempt",
			apiKey:      "'; DROP TABLE users; --",
			wantErr:     true,
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create client with invalid API key
			client, err := New(tt.apiKey)

			// Some keys might fail at client creation
			if err != nil {
				if !tt.wantErr {
					t.Errorf("New() unexpected error: %v", err)
				}
				// Verify error message is clear
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}

			// Try to get viewer (authenticate)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			viewer, err := client.GetViewer(ctx)

			// We expect authentication to fail
			if !tt.wantErr && err != nil {
				t.Errorf("GetViewer() unexpected error: %v", err)
			}

			if tt.wantErr && err == nil {
				t.Errorf("GetViewer() expected error, got nil (viewer: %+v)", viewer)
			}

			// Verify error message is clear and helpful
			if err != nil {
				errMsg := err.Error()

				// Error should not expose the invalid key
				if strings.Contains(errMsg, tt.apiKey) && len(tt.apiKey) > 10 {
					t.Errorf("Error message exposes invalid API key: %q", errMsg)
				}

				// Error should be clear and actionable
				if tt.errContains != "" && !strings.Contains(errMsg, tt.errContains) {
					t.Errorf("Error message %q should contain %q", errMsg, tt.errContains)
				}

				// Error should indicate it's an authentication issue
				lowerErr := strings.ToLower(errMsg)
				hasAuthKeyword := strings.Contains(lowerErr, "auth") ||
					strings.Contains(lowerErr, "unauthorized") ||
					strings.Contains(lowerErr, "invalid") ||
					strings.Contains(lowerErr, "credential") ||
					strings.Contains(lowerErr, "api key")

				if !hasAuthKeyword {
					t.Logf("Warning: Error message may not clearly indicate auth failure: %q", errMsg)
				}
			}
		})
	}
}

// TestAuthWithValidAPIKey verifies that a valid API key can authenticate.
// This serves as a control test to ensure auth is working correctly.
//
// Requires LINEAR_TEST_API_KEY environment variable to be set to a VALID key.
func TestAuthWithValidAPIKey(t *testing.T) {
	testAPIKey := os.Getenv("LINEAR_TEST_API_KEY")
	if testAPIKey == "" {
		t.Skip("Skipping auth test: LINEAR_TEST_API_KEY not set")
	}

	// Assume the provided test key is valid
	client, err := New(testAPIKey)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	viewer, err := client.GetViewer(ctx)
	if err != nil {
		t.Fatalf("GetViewer() failed with valid API key: %v", err)
	}

	// Verify viewer has expected fields
	if viewer.ID == "" {
		t.Error("Viewer ID is empty")
	}
	if viewer.Name == "" {
		t.Error("Viewer Name is empty")
	}
	if viewer.Email == "" {
		t.Error("Viewer Email is empty")
	}
	if viewer.Organization == nil {
		t.Error("Viewer Organization is nil")
	} else {
		if viewer.Organization.ID == "" {
			t.Error("Organization ID is empty")
		}
		if viewer.Organization.Name == "" {
			t.Error("Organization Name is empty")
		}
	}
}

// TestAuthErrorFormat verifies that authentication errors follow a consistent format
func TestAuthErrorFormat(t *testing.T) {
	testAPIKey := os.Getenv("LINEAR_TEST_API_KEY")
	if testAPIKey == "" {
		t.Skip("Skipping auth test: LINEAR_TEST_API_KEY not set")
	}

	invalidKey := "definitely_not_a_valid_key_12345678"
	client, err := New(invalidKey)
	if err != nil {
		// Error at client creation is acceptable
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.GetViewer(ctx)
	if err == nil {
		t.Fatal("Expected authentication error with invalid key")
	}

	errMsg := err.Error()

	// Verify error format is helpful
	t.Logf("Auth error format: %q", errMsg)

	// Error should not be empty
	if errMsg == "" {
		t.Error("Error message is empty")
	}

	// Error should not be too verbose (< 200 chars is reasonable)
	if len(errMsg) > 200 {
		t.Errorf("Error message is too verbose (%d chars): %q", len(errMsg), errMsg)
	}
}

// TestRateLimiting verifies that the client handles rate limiting appropriately
// (This test might be slow or require special setup)
func TestRateLimiting(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping rate limit test in short mode")
	}

	testAPIKey := os.Getenv("LINEAR_TEST_API_KEY")
	if testAPIKey == "" {
		t.Skip("Skipping rate limit test: LINEAR_TEST_API_KEY not set")
	}

	// This is a placeholder for rate limiting tests
	// In practice, you'd want to make many rapid requests and verify
	// the client handles rate limit responses gracefully
	t.Skip("Rate limiting test not yet implemented")
}
