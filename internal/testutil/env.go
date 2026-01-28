package testutil

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadTestEnv loads environment variables from .env.test files.
// Search order:
// 1. .env.test.local in current directory (worktree-specific, highest priority)
// 2. .env.test in current directory (worktree-level)
// 3. .env.test in git root (project-wide, lowest priority)
//
// Later files override earlier ones if the same variable is defined.
func LoadTestEnv() error {
	// Find git root by walking up to find .repo.git or .git
	gitRoot, err := findGitRoot()
	if err != nil {
		// If we can't find git root, try current directory only
		return loadEnvFiles(".")
	}

	// Load from git root first (lowest priority)
	gitRootEnv := filepath.Join(gitRoot, ".env.test")
	_ = godotenv.Load(gitRootEnv) // Ignore error, file may not exist

	// Load from current directory (overrides git root)
	_ = godotenv.Load(".env.test") // Ignore error, file may not exist

	// Load local overrides (highest priority)
	_ = godotenv.Load(".env.test.local") // Ignore error, file may not exist

	// Always return nil - .env.test files are optional
	// Tests will skip if required env vars are not set
	return nil
}

// MustLoadTestEnv loads test environment or panics.
// Use in test init() functions.
//
// This function never panics in practice - it silently ignores missing
// .env.test files, as they are optional. Tests will skip if required
// environment variables are not set.
func MustLoadTestEnv() {
	if err := LoadTestEnv(); err != nil {
		// Silently ignore - .env.test is optional
		// Tests will skip if required env vars not set
	}
}

// findGitRoot walks up directory tree to find .git or .repo.git directory.
// This handles both standard git repos and git worktrees.
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for .repo.git (worktree parent)
		repoGitPath := filepath.Join(dir, ".repo.git")
		if stat, err := os.Stat(repoGitPath); err == nil && stat.IsDir() {
			return dir, nil
		}

		// Check for .git (standard repo or worktree file)
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			// If .git is a file, it's a worktree - read gitdir and go up
			if _, err := os.ReadFile(gitPath); err == nil {
				// .git file contains "gitdir: /path/to/worktree"
				// We want to find the repo root, so keep walking up
				parent := filepath.Dir(dir)
				if parent == dir {
					return dir, nil // At filesystem root
				}
				dir = parent
				continue
			}
			// If .git is a directory, this is the root
			return dir, nil
		}

		// Walk up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root without finding .git
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

// loadEnvFiles attempts to load .env.test files from the given directory
func loadEnvFiles(dir string) error {
	_ = godotenv.Load(filepath.Join(dir, ".env.test"))
	_ = godotenv.Load(filepath.Join(dir, ".env.test.local"))
	return nil
}

// GetTestAPIKey returns the LINEAR_TEST_API_KEY environment variable.
// This is a convenience function for tests that need the API key.
func GetTestAPIKey() string {
	return os.Getenv("LINEAR_TEST_API_KEY")
}

// HasTestAPIKey returns true if LINEAR_TEST_API_KEY is set.
// Use this to conditionally skip tests that require API access.
func HasTestAPIKey() bool {
	return GetTestAPIKey() != ""
}
