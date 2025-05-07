package implementations

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDnfPackageManager(t *testing.T) {
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("dnf not available, skipping test")
	}

	pm, err := NewDnfPackageManager()
	if err != nil {
		t.Fatalf("NewDnfPackageManager() error = %v", err)
	}
	if pm == nil {
		t.Fatal("NewDnfPackageManager() returned nil")
	}
}

func TestDnfPackageManager_Install(t *testing.T) {
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("dnf not available, skipping test")
	}

	pm, err := NewDnfPackageManager()
	if err != nil {
		t.Fatalf("NewDnfPackageManager() error = %v", err)
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

func TestDnfPackageManager_IsInstalled(t *testing.T) {
	pm, err := NewDnfPackageManager()
	assert.NoError(t, err)

	// Check for dnf itself (should be installed if test runs)
	installed, err := pm.IsInstalled("dnf") // Handle error
	assert.NoError(t, err) // Expect no error
	assert.True(t, installed, "Expected 'dnf' to be installed")

	// Check for a non-existent package
	installed, err = pm.IsInstalled("nonexistent-package-qwezxc") // Handle error
	assert.NoError(t, err) // Expect no error, just false
	assert.False(t, installed, "Expected 'nonexistent-package-qwezxc' not to be installed")
}

func TestDnfPackageManager_Remove(t *testing.T) {
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("dnf not available, skipping test")
	}

	pm, err := NewDnfPackageManager()
	if err != nil {
		t.Fatalf("NewDnfPackageManager() error = %v", err)
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