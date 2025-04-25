// Package testing provides testing utilities and helpers for the bootstrap-cli,
// including test environment setup, temporary file and directory management,
// and environment variable handling.
package testing

import (
	"os"
	"testing"
)

// TestEnv represents a test environment
type TestEnv struct {
	t       *testing.T
	HomeDir string
	cleanup func()
}

// NewTestEnv creates a new test environment
func NewTestEnv(t *testing.T, homeDir string) *TestEnv {
	// Store original home directory
	originalHome := os.Getenv("HOME")

	// Set test home directory
	err := os.Setenv("HOME", homeDir)
	if err != nil {
		t.Fatalf("Failed to set HOME environment variable: %v", err)
	}

	return &TestEnv{
		t:       t,
		HomeDir: homeDir,
		cleanup: func() {
			// Restore original home directory
			if err := os.Setenv("HOME", originalHome); err != nil {
				t.Errorf("Failed to restore HOME environment variable: %v", err)
			}
		},
	}
}

// Cleanup restores the original environment
func (e *TestEnv) Cleanup() {
	if e.cleanup != nil {
		e.cleanup()
	}
}

// CreateTempFile creates a temporary file in the test environment
func (e *TestEnv) CreateTempFile(name string, content []byte) string {
	path := e.HomeDir + "/" + name
	if err := os.MkdirAll(e.HomeDir, 0755); err != nil {
		e.t.Fatalf("Failed to create home directory: %v", err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		e.t.Fatalf("Failed to write file %s: %v", name, err)
	}
	return path
}

// CreateTempDir creates a temporary directory in the test environment
func (e *TestEnv) CreateTempDir(name string) string {
	path := e.HomeDir + "/" + name
	if err := os.MkdirAll(path, 0755); err != nil {
		e.t.Fatalf("Failed to create directory %s: %v", name, err)
	}
	return path
} 