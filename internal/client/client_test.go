package client

import (
	"strings"
	"testing"
)

// TestMaskAPIKey verifies that API keys are properly masked and full keys
// are never displayed.
// Related bead: lirt-3po
func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Standard Linear API key",
			input:    "lin_api_1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "lin_api_1234...",
			wantErr:  false,
		},
		{
			name:     "Long API key",
			input:    "lin_api_abcdefghijklmnopqrstuvwxyz1234567890",
			expected: "lin_api_abcd...",
			wantErr:  false,
		},
		{
			name:     "Short API key",
			input:    "short_key",
			expected: "***",
			wantErr:  false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "***",
			wantErr:  false,
		},
		{
			name:     "UUID-style key",
			input:    "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			expected: "a1b2c3d4-e5f...",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskAPIKey(tt.input)

			// Verify result matches expected
			if result != tt.expected {
				t.Errorf("MaskAPIKey(%q) = %q, want %q", tt.input, result, tt.expected)
			}

			// Verify full key is never in result (security check)
			if tt.input != "" && len(tt.input) > 12 {
				if strings.Contains(result, tt.input) {
					t.Errorf("MaskAPIKey exposed full key: %q contains %q", result, tt.input)
				}

				// Verify suffix is masked (beyond first 12 chars)
				suffix := tt.input[12:]
				if len(suffix) > 0 && strings.Contains(result, suffix) {
					t.Errorf("MaskAPIKey exposed key suffix: %q contains %q", result, suffix)
				}
			}

			// Verify result contains "..." or "***"
			if !strings.Contains(result, "...") && !strings.Contains(result, "***") {
				t.Errorf("MaskAPIKey result %q should contain masking indicator", result)
			}
		})
	}
}

// TestMaskAPIKeyLength verifies that masked keys never expose more than 12 characters
func TestMaskAPIKeyLength(t *testing.T) {
	// Generate various key lengths
	keySizes := []int{5, 10, 12, 15, 20, 50, 100}

	for _, size := range keySizes {
		t.Run(string(rune(size)), func(t *testing.T) {
			// Create key of specified size
			key := strings.Repeat("x", size)

			masked := MaskAPIKey(key)

			// Count visible characters (exclude "..." and "***")
			visible := strings.TrimSuffix(masked, "...")
			visible = strings.TrimSuffix(visible, "***")

			if len(visible) > 12 {
				t.Errorf("MaskAPIKey exposed %d chars for %d-char key (max 12 allowed): %q",
					len(visible), size, masked)
			}
		})
	}
}

// TestMaskAPIKeyNoFullKeyExposure is a fuzz-style test to ensure no full key
// is ever exposed regardless of input
func TestMaskAPIKeyNoFullKeyExposure(t *testing.T) {
	testKeys := []string{
		"lin_api_secret1234567890abcdefghijklmnop",
		"super_secret_key_do_not_expose_12345",
		strings.Repeat("a", 100),
		"lin_api_" + strings.Repeat("x", 50),
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09",
		"special!@#$%^&*()chars",
	}

	for i, key := range testKeys {
		t.Run(string(rune(i)), func(t *testing.T) {
			if len(key) < 13 {
				return // Short keys are fully masked
			}

			masked := MaskAPIKey(key)

			// The dangerous part (after first 12 chars) should never appear
			dangerousSubstring := key[12:]
			if strings.Contains(masked, dangerousSubstring) {
				t.Errorf("SECURITY ISSUE: MaskAPIKey exposed secret part of key\nFull key: %q\nMasked: %q\nExposed: %q",
					key, masked, dangerousSubstring)
			}
		})
	}
}
