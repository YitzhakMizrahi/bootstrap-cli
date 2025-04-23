package implementations

import (
	"os/exec"
	"testing"
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
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("brew not available, skipping test")
	}

	pm, err := NewHomebrewPackageManager()
	if err != nil {
		t.Fatalf("NewHomebrewPackageManager() error = %v", err)
	}

	// Test with a package that should be installed
	if !pm.IsInstalled("brew") {
		t.Error("IsInstalled() returned false for brew package")
	}

	// Test with a non-existent package
	if pm.IsInstalled("non-existent-package-123456") {
		t.Error("IsInstalled() returned true for non-existent package")
	}
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
	err = pm.Remove("curl")
	if err != nil {
		t.Errorf("Remove() error = %v", err)
	}

	// Test removing a non-existent package
	err = pm.Remove("non-existent-package-123456")
	if err == nil {
		t.Error("Remove() expected error for non-existent package, got nil")
	}
} 