package system

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
	if pm.brewPath == "" {
		t.Error("brewPath is empty")
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
	if !pm.IsInstalled("tldr") {
		t.Error("IsInstalled() returned false for tldr package")
	}

	// Test with a non-existent package
	if pm.IsInstalled("non-existent-package-123456") {
		t.Error("IsInstalled() returned true for non-existent package")
	}
}

func TestHomebrewPackageManager_Uninstall(t *testing.T) {
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("brew not available, skipping test")
	}

	pm, err := NewHomebrewPackageManager()
	if err != nil {
		t.Fatalf("NewHomebrewPackageManager() error = %v", err)
	}

	// Test uninstalling a package that should exist
	err = pm.Uninstall("curl")
	if err != nil {
		t.Errorf("Uninstall() error = %v", err)
	}

	// Test uninstalling a non-existent package
	err = pm.Uninstall("non-existent-package-123456")
	if err == nil {
		t.Error("Uninstall() expected error for non-existent package, got nil")
	}
} 