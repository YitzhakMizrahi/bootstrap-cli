package languages

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewNodeJSManager(t *testing.T) {
	manager := NewNodeJSManager("/tmp/nodejs")
	
	if manager.InstallPath != "/tmp/nodejs" {
		t.Errorf("Expected InstallPath to be '/tmp/nodejs', got '%s'", manager.InstallPath)
	}
	
	if manager.CurrentVersion != "" {
		t.Errorf("Expected CurrentVersion to be empty, got '%s'", manager.CurrentVersion)
	}
	
	if len(manager.AvailableVersions) != 0 {
		t.Errorf("Expected AvailableVersions to be empty, got %v", manager.AvailableVersions)
	}
}

func TestIsNVMInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nvm_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a manager with the temporary directory
	manager := NewNodeJSManager(tempDir)
	
	// Check if NVM is installed (should be false)
	if manager.isNVMInstalled() {
		t.Errorf("Expected NVM to not be installed, but it is")
	}
	
	// Create the NVM directory and script
	nvmDir := filepath.Join(os.Getenv("HOME"), ".nvm")
	err = os.MkdirAll(nvmDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create NVM directory: %v", err)
	}
	
	nvmScript := filepath.Join(nvmDir, "nvm.sh")
	err = os.WriteFile(nvmScript, []byte("#!/bin/bash\necho 'NVM script'"), 0755)
	if err != nil {
		t.Fatalf("Failed to create NVM script: %v", err)
	}
	
	// Check if NVM is installed (should be true)
	if !manager.isNVMInstalled() {
		t.Errorf("Expected NVM to be installed, but it is not")
	}
	
	// Clean up
	err = os.RemoveAll(nvmDir)
	if err != nil {
		t.Fatalf("Failed to clean up NVM directory: %v", err)
	}
}

func TestGetNVMInitScript(t *testing.T) {
	manager := NewNodeJSManager("/tmp/nodejs")
	
	expectedPath := filepath.Join(os.Getenv("HOME"), ".nvm", "nvm.sh")
	actualPath := manager.getNVMInitScript()
	
	if actualPath != expectedPath {
		t.Errorf("Expected NVM init script path to be '%s', got '%s'", expectedPath, actualPath)
	}
}

// Note: The following tests require NVM to be installed and may not work in all environments.
// They are included for completeness but may be skipped in environments where NVM is not available.

func TestInstallNVM(t *testing.T) {
	// Skip this test if curl is not available
	if _, err := os.Stat("/usr/bin/curl"); os.IsNotExist(err) {
		t.Skip("Curl is not available, skipping test")
	}
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nvm_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a manager with the temporary directory
	manager := NewNodeJSManager(tempDir)
	
	// Install NVM
	err = manager.InstallNVM()
	if err != nil {
		t.Fatalf("Failed to install NVM: %v", err)
	}
	
	// Check if NVM is installed
	if !manager.isNVMInstalled() {
		t.Errorf("Expected NVM to be installed, but it is not")
	}
	
	// Clean up
	nvmDir := filepath.Join(os.Getenv("HOME"), ".nvm")
	err = os.RemoveAll(nvmDir)
	if err != nil {
		t.Fatalf("Failed to clean up NVM directory: %v", err)
	}
}

func TestInstallNodeJS(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// Install a specific version of Node.js
	err := manager.InstallNodeJS("16.20.0")
	if err != nil {
		t.Fatalf("Failed to install Node.js: %v", err)
	}
	
	// Check if the version is set
	if manager.CurrentVersion != "16.20.0" {
		t.Errorf("Expected current version to be '16.20.0', got '%s'", manager.CurrentVersion)
	}
}

func TestInstallLatestLTS(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// Install the latest LTS version of Node.js
	err := manager.InstallLatestLTS()
	if err != nil {
		t.Fatalf("Failed to install latest LTS version of Node.js: %v", err)
	}
	
	// Check if the version is set
	if manager.CurrentVersion == "" {
		t.Errorf("Expected current version to be set, got empty string")
	}
}

func TestInstallLatestCurrent(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// Install the latest current version of Node.js
	err := manager.InstallLatestCurrent()
	if err != nil {
		t.Fatalf("Failed to install latest current version of Node.js: %v", err)
	}
	
	// Check if the version is set
	if manager.CurrentVersion == "" {
		t.Errorf("Expected current version to be set, got empty string")
	}
}

func TestListAvailableVersions(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// List available versions
	versions, err := manager.ListAvailableVersions()
	if err != nil {
		t.Fatalf("Failed to list available versions: %v", err)
	}
	
	// Check if versions are returned
	if len(versions) == 0 {
		t.Errorf("Expected versions to be returned, got empty slice")
	}
	
	// Check if the versions are stored in the manager
	if len(manager.AvailableVersions) == 0 {
		t.Errorf("Expected available versions to be stored in the manager, got empty slice")
	}
}

func TestListInstalledVersions(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// List installed versions
	versions, err := manager.ListInstalledVersions()
	if err != nil {
		t.Fatalf("Failed to list installed versions: %v", err)
	}
	
	// Check if versions are returned
	if len(versions) == 0 {
		t.Errorf("Expected versions to be returned, got empty slice")
	}
}

func TestUseVersion(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// Use a specific version
	err := manager.UseVersion("16.20.0")
	if err != nil {
		t.Fatalf("Failed to use version: %v", err)
	}
	
	// Check if the version is set
	if manager.CurrentVersion != "16.20.0" {
		t.Errorf("Expected current version to be '16.20.0', got '%s'", manager.CurrentVersion)
	}
}

func TestGetCurrentVersion(t *testing.T) {
	// Skip this test if NVM is not installed
	manager := NewNodeJSManager("/tmp/nodejs")
	if !manager.isNVMInstalled() {
		t.Skip("NVM is not installed, skipping test")
	}
	
	// Get the current version
	version, err := manager.GetCurrentVersion()
	if err != nil {
		t.Fatalf("Failed to get current version: %v", err)
	}
	
	// Check if the version is returned
	if version == "" {
		t.Errorf("Expected version to be returned, got empty string")
	}
	
	// Check if the version is stored in the manager
	if manager.CurrentVersion == "" {
		t.Errorf("Expected current version to be stored in the manager, got empty string")
	}
} 