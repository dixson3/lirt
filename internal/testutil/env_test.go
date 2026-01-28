package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadTestEnv verifies that .env.test files are loaded correctly
func TestLoadTestEnv(t *testing.T) {
	// Save original env
	originalKey := os.Getenv("LINEAR_TEST_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("LINEAR_TEST_API_KEY", originalKey)
		} else {
			os.Unsetenv("LINEAR_TEST_API_KEY")
		}
	}()

	// Clear env var
	os.Unsetenv("LINEAR_TEST_API_KEY")

	// LoadTestEnv should not error even if files don't exist
	err := LoadTestEnv()
	if err != nil {
		t.Errorf("LoadTestEnv() should not return error: %v", err)
	}

	// If .env.test exists at git root, verify it was loaded
	gitRoot, _ := findGitRoot()
	if gitRoot != "" {
		envFile := filepath.Join(gitRoot, ".env.test")
		if _, err := os.Stat(envFile); err == nil {
			// File exists, check if env var was loaded
			key := os.Getenv("LINEAR_TEST_API_KEY")
			if key == "" {
				t.Log("Note: .env.test exists but LINEAR_TEST_API_KEY not set (may be empty in file)")
			} else {
				t.Logf("Successfully loaded LINEAR_TEST_API_KEY from .env.test: %s", MaskKey(key))
			}
		}
	}
}

// TestFindGitRoot verifies that findGitRoot works in the lirt project
func TestFindGitRoot(t *testing.T) {
	root, err := findGitRoot()
	if err != nil {
		t.Errorf("findGitRoot() error = %v", err)
		return
	}

	if root == "" {
		t.Error("findGitRoot() returned empty string")
		return
	}

	// Verify that the root contains either .git or .repo.git
	gitPath := filepath.Join(root, ".git")
	repoGitPath := filepath.Join(root, ".repo.git")

	hasGit := false
	if _, err := os.Stat(gitPath); err == nil {
		hasGit = true
	}
	if _, err := os.Stat(repoGitPath); err == nil {
		hasGit = true
	}

	if !hasGit {
		t.Errorf("findGitRoot() returned %s which has no .git or .repo.git", root)
	}

	t.Logf("Git root: %s", root)
}

// TestHasTestAPIKey verifies the helper function
func TestHasTestAPIKey(t *testing.T) {
	// Save original
	originalKey := os.Getenv("LINEAR_TEST_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("LINEAR_TEST_API_KEY", originalKey)
		} else {
			os.Unsetenv("LINEAR_TEST_API_KEY")
		}
	}()

	// Test with key set
	os.Setenv("LINEAR_TEST_API_KEY", "test_key")
	if !HasTestAPIKey() {
		t.Error("HasTestAPIKey() should return true when key is set")
	}
	if GetTestAPIKey() != "test_key" {
		t.Error("GetTestAPIKey() should return the set key")
	}

	// Test with key unset
	os.Unsetenv("LINEAR_TEST_API_KEY")
	if HasTestAPIKey() {
		t.Error("HasTestAPIKey() should return false when key is not set")
	}
	if GetTestAPIKey() != "" {
		t.Error("GetTestAPIKey() should return empty string when key is not set")
	}
}

// MaskKey masks an API key for safe logging (first 12 chars + "...")
func MaskKey(key string) string {
	if len(key) <= 12 {
		return "***"
	}
	return key[:12] + "..."
}
