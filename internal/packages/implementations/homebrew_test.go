package implementations

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHomebrewPackageManager(t *testing.T) {
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("brew not available, skipping test")
	}

	pm, err := NewHomebrewPackageManager()
	if err != nil {
		t.Fatalf("NewHomebrewPackageManager() error = %v", err)
	}
	if pm == nil {
		t.Fatal("NewHomebrewPackageManager() returned nil")
	}
}

func TestHomebrewPackageManager_Install(t *testing.T) {
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("brew not available, skipping test")
	}

	pm, err := NewHomebrewPackageManager()
	if err != nil {
		t.Fatalf("NewHomebrewPackageManager() error = %v", err)
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

func TestHomebrewPackageManager_IsInstalled(t *testing.T) {
	pm, err := NewHomebrewPackageManager()
	assert.NoError(t, err)

	// Check for brew itself (should be installed if test runs)
	installed, err := pm.IsInstalled("brew") // Handle error
	assert.NoError(t, err) // Expect no error
	assert.True(t, installed, "Expected 'brew' to be installed")

	// Check for a non-existent package
	installed, err = pm.IsInstalled("nonexistent-package-xyzabc") // Handle error
	assert.NoError(t, err) // Expect no error, just false
	assert.False(t, installed, "Expected 'nonexistent-package-xyzabc' not to be installed")
}

func TestHomebrewPackageManager_Remove(t *testing.T) {
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("brew not available, skipping test")
	}

	pm, err := NewHomebrewPackageManager()
	if err != nil {
		t.Fatalf("NewHomebrewPackageManager() error = %v", err)
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