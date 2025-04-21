package languages

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestNewPythonManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "python-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new PythonManager
	manager := NewPythonManager(tempDir)

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

func TestIsPyenvInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "python-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new PythonManager
	manager := NewPythonManager(tempDir)

	// Check if pyenv is installed (should be false initially)
	if manager.isPyenvInstalled() {
		t.Error("Expected pyenv to not be installed initially")
	}

	// Create a mock pyenv installation
	pyenvDir := filepath.Join(tempDir, ".pyenv")
	if err := os.MkdirAll(filepath.Join(pyenvDir, "bin"), 0755); err != nil {
		t.Fatalf("Failed to create mock pyenv directory: %v", err)
	}

	// Create a mock pyenv executable
	pyenvExec := filepath.Join(pyenvDir, "bin", "pyenv")
	if err := os.WriteFile(pyenvExec, []byte("#!/bin/bash\necho 'mock pyenv'"), 0755); err != nil {
		t.Fatalf("Failed to create mock pyenv executable: %v", err)
	}

	// Check if pyenv is installed (should be true now)
	if !manager.isPyenvInstalled() {
		t.Error("Expected pyenv to be installed after creating mock installation")
	}
}

func TestGetShellConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "python-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new PythonManager
	manager := NewPythonManager(tempDir)

	// Test with default shell (bash)
	configFile := manager.getShellConfig()
	expectedConfig := filepath.Join(os.Getenv("HOME"), ".bashrc")
	if configFile != expectedConfig {
		t.Errorf("Expected config file to be %s, got %s", expectedConfig, configFile)
	}

	// Test with zsh shell
	os.Setenv("SHELL", "/bin/zsh")
	configFile = manager.getShellConfig()
	expectedConfig = filepath.Join(os.Getenv("HOME"), ".zshrc")
	if configFile != expectedConfig {
		t.Errorf("Expected config file to be %s, got %s", expectedConfig, configFile)
	}
}

// Note: The following tests require pyenv to be installed and may not work in all environments.
// They are included for completeness but may be skipped if pyenv is not available.

func TestInstallPyenv(t *testing.T) {
	// Skip if git is not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Skipping test: git is not available")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "python-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new PythonManager
	manager := NewPythonManager(tempDir)

	// Install pyenv
	if err := manager.InstallPyenv(); err != nil {
		t.Errorf("Failed to install pyenv: %v", err)
	}

	// Check if pyenv is installed
	if !manager.isPyenvInstalled() {
		t.Error("Expected pyenv to be installed after installation")
	}
}

func TestInstallPython(t *testing.T) {
	// Skip if pyenv is not installed
	manager := NewPythonManager(os.Getenv("HOME"))
	if !manager.isPyenvInstalled() {
		t.Skip("Skipping test: pyenv is not installed")
	}

	// Install a specific version of Python
	version := "3.9.0"
	if err := manager.InstallPython(version); err != nil {
		t.Errorf("Failed to install Python %s: %v", version, err)
	}

	// Check if the version is installed
	installedVersions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Errorf("Failed to list installed versions: %v", err)
	}

	found := false
	for _, v := range installedVersions {
		if v == version {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected Python %s to be installed, but it was not found in the list of installed versions", version)
	}
} 