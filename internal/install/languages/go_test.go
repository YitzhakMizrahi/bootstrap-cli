package languages

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGoManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new GoManager
	manager := NewGoManager(tempDir)

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

func TestIsGoInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new GoManager
	manager := NewGoManager(tempDir)

	// Check if Go is installed (should be false initially)
	if manager.isGoInstalled() {
		t.Error("Expected Go to not be installed initially")
	}

	// Create a mock Go installation
	goDir := filepath.Join(tempDir, "go")
	if err := os.MkdirAll(filepath.Join(goDir, "bin"), 0755); err != nil {
		t.Fatalf("Failed to create mock Go directory: %v", err)
	}

	// Create a mock Go executable
	goExec := filepath.Join(goDir, "bin", "go")
	if err := os.WriteFile(goExec, []byte("#!/bin/bash\necho 'mock go'"), 0755); err != nil {
		t.Fatalf("Failed to create mock Go executable: %v", err)
	}

	// Check if Go is installed (should be true now)
	if !manager.isGoInstalled() {
		t.Error("Expected Go to be installed after creating mock installation")
	}
}

// Note: The following tests are commented out as they are redeclared in another file.
// Uncomment them if you want to use them in this file.

/*
func TestGetCurrentVersion(t *testing.T) {
	// Skip if Go is not installed
	manager := NewGoManager(os.Getenv("HOME"))
	if !manager.isGoInstalled() {
		t.Skip("Skipping test: Go is not installed")
	}

	// Get the current version
	version, err := manager.GetCurrentVersion()
	if err != nil {
		t.Errorf("Failed to get current Go version: %v", err)
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
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "go-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new GoManager
	manager := NewGoManager(tempDir)

	// List available versions
	versions, err := manager.ListAvailableVersions()
	if err != nil {
		t.Errorf("Failed to list available versions: %v", err)
	}

	if len(versions) == 0 {
		t.Error("Expected to find available Go versions, but none were found")
	}

	// Check if the versions are stored in the manager
	if len(manager.AvailableVersions) == 0 {
		t.Error("Expected available versions to be stored in the manager, but none were found")
	}
}

func TestListInstalledVersions(t *testing.T) {
	// Skip if Go is not installed
	manager := NewGoManager(os.Getenv("HOME"))
	if !manager.isGoInstalled() {
		t.Skip("Skipping test: Go is not installed")
	}

	// List installed versions
	versions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Errorf("Failed to list installed versions: %v", err)
	}

	// Note: We don't check the length of versions here because the user may not have any Go versions installed
}

func TestUseVersion(t *testing.T) {
	// Skip if Go is not installed
	manager := NewGoManager(os.Getenv("HOME"))
	if !manager.isGoInstalled() {
		t.Skip("Skipping test: Go is not installed")
	}

	// Get the list of installed versions
	installedVersions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Errorf("Failed to list installed versions: %v", err)
	}

	if len(installedVersions) == 0 {
		t.Skip("Skipping test: no Go versions are installed")
	}

	// Use the first installed version
	version := installedVersions[0]
	if err := manager.UseVersion(version); err != nil {
		t.Errorf("Failed to use Go version %s: %v", version, err)
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