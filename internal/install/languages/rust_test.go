package languages

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRustManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "rust-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new RustManager
	manager := NewRustManager(tempDir)

	// Check if the fields are initialized correctly
	if manager.InstallPath != tempDir {
		t.Errorf("Expected InstallPath to be %s, got %s", tempDir, manager.InstallPath)
	}

	if manager.CurrentVersion != "" {
		t.Errorf("Expected CurrentVersion to be empty, got %s", manager.CurrentVersion)
	}

	if len(manager.AvailableVersions) != 0 {
		t.Errorf("Expected AvailableVersions to be empty, got %v", manager.AvailableVersions)
	}
}

func TestIsRustupInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "rust-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new RustManager
	manager := NewRustManager(tempDir)

	// Check if rustup is installed (should be false initially)
	if manager.isRustupInstalled() {
		t.Skip("Skipping test: rustup is already installed on this system")
	}

	// Create a mock rustup installation
	rustupDir := filepath.Join(tempDir, ".cargo", "bin")
	if err := os.MkdirAll(rustupDir, 0755); err != nil {
		t.Fatalf("Failed to create mock rustup directory: %v", err)
	}

	// Create a mock rustup executable
	rustupExec := filepath.Join(rustupDir, "rustup")
	if err := os.WriteFile(rustupExec, []byte("#!/bin/bash\necho 'mock rustup'"), 0755); err != nil {
		t.Fatalf("Failed to create mock rustup executable: %v", err)
	}

	// Add the mock directory to PATH
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)
	os.Setenv("PATH", rustupDir+":"+oldPath)

	// Check if rustup is installed (should be true now)
	if !manager.isRustupInstalled() {
		t.Error("Expected rustup to be installed after creating mock installation")
	}
}

// Note: The following tests are commented out as they are redeclared in another file.
// Uncomment them if you want to use them in this file.

/*
func TestGetCurrentVersion(t *testing.T) {
	// Skip if rustup is not installed
	manager := NewRustManager(os.Getenv("HOME"))
	if !manager.isRustupInstalled() {
		t.Skip("Skipping test: rustup is not installed")
	}

	// Get the current version
	version, err := manager.GetCurrentVersion()
	if err != nil {
		t.Errorf("Failed to get current Rust version: %v", err)
	}

	if version == "" {
		t.Error("Expected to get a current version, but got an empty string")
	}

	// Check if the version is stored in the manager
	if manager.CurrentVersion != version {
		t.Errorf("Expected CurrentVersion to be %s, got %s", version, manager.CurrentVersion)
	}
}

func TestListAvailableVersions(t *testing.T) {
	// Skip if rustup is not installed
	manager := NewRustManager(os.Getenv("HOME"))
	if !manager.isRustupInstalled() {
		t.Skip("Skipping test: rustup is not installed")
	}

	// List available versions
	versions, err := manager.ListAvailableVersions()
	if err != nil {
		t.Errorf("Failed to list available versions: %v", err)
	}

	if len(versions) == 0 {
		t.Error("Expected to find available Rust versions, but none were found")
	}

	// Check if the versions are stored in the manager
	if len(manager.AvailableVersions) == 0 {
		t.Error("Expected available versions to be stored in the manager, but none were found")
	}
}

func TestListInstalledVersions(t *testing.T) {
	// Skip if rustup is not installed
	manager := NewRustManager(os.Getenv("HOME"))
	if !manager.isRustupInstalled() {
		t.Skip("Skipping test: rustup is not installed")
	}

	// List installed versions
	versions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Errorf("Failed to list installed versions: %v", err)
	}

	// Note: We don't check the length of versions here because the user may not have any Rust versions installed
}

func TestUseVersion(t *testing.T) {
	// Skip if rustup is not installed
	manager := NewRustManager(os.Getenv("HOME"))
	if !manager.isRustupInstalled() {
		t.Skip("Skipping test: rustup is not installed")
	}

	// Get the list of installed versions
	installedVersions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Errorf("Failed to list installed versions: %v", err)
	}

	if len(installedVersions) == 0 {
		t.Skip("Skipping test: no Rust versions are installed")
	}

	// Use the first installed version
	version := installedVersions[0]
	if err := manager.UseVersion(version); err != nil {
		t.Errorf("Failed to use Rust version %s: %v", version, err)
	}

	// Check if the current version is set correctly
	currentVersion, err := manager.GetCurrentVersion()
	if err != nil {
		t.Errorf("Failed to get current version: %v", err)
	}

	if currentVersion != version {
		t.Errorf("Expected current version to be %s, got %s", version, currentVersion)
	}
}
*/ 