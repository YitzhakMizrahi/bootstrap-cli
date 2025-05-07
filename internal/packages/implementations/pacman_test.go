package implementations

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPacmanPackageManager(t *testing.T) {
	// Skip if not on a system with pacman
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not available, skipping test")
	}

	pm, err := NewPacmanPackageManager()
	if err != nil {
		t.Fatalf("NewPacmanPackageManager() error = %v", err)
	}
	if pm == nil {
		t.Fatal("NewPacmanPackageManager() returned nil")
	}
}

func TestPacmanPackageManager_Install(t *testing.T) {
	// Skip if not on a system with pacman
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not available, skipping test")
	}

	pm, err := NewPacmanPackageManager()
	if err != nil {
		t.Fatalf("NewPacmanPackageManager() error = %v", err)
	}

	// Test installing a package that should exist
	err = pm.Install("curl")
	if err != nil {
		t.Errorf("Install() error = %v", err)
	}

	// Test installing a non-existent package
	err = pm.Install("non-existent-package-123456")
	if err == nil {
		t.Error("Install() expected error for non-existent package, got nil")
	}
}

func TestPacmanPackageManager_IsInstalled(t *testing.T) {
	pm, err := NewPacmanPackageManager()
	assert.NoError(t, err)

	// Check for pacman itself (should be installed if test runs)
	installed, err := pm.IsInstalled("pacman") // Handle error
	assert.NoError(t, err) // Expect no error
	assert.True(t, installed, "Expected 'pacman' to be installed")

	// Check for a non-existent package
	installed, err = pm.IsInstalled("nonexistent-package-qwertyuiop") // Handle error
	assert.NoError(t, err) // Expect no error, just false
	assert.False(t, installed, "Expected 'nonexistent-package-qwertyuiop' not to be installed")
}

func TestPacmanPackageManager_Remove(t *testing.T) {
	// Skip if not on a system with pacman
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not available, skipping test")
	}

	pm, err := NewPacmanPackageManager()
	if err != nil {
		t.Fatalf("NewPacmanPackageManager() error = %v", err)
	}

	// Test removing a package that should exist
	err = pm.Uninstall("curl")
	if err != nil {
		t.Errorf("Uninstall() error = %v", err)
	}

	// Test removing a non-existent package
	err = pm.Uninstall("non-existent-package-123456")
	if err == nil {
		t.Error("Uninstall() expected error for non-existent package, got nil")
	}
} 