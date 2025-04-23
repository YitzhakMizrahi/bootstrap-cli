package system

import (
	"os/exec"
	"testing"
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
	if pm.sudoPath == "" {
		t.Error("sudoPath is empty")
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
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("dnf not available, skipping test")
	}

	pm, err := NewDnfPackageManager()
	if err != nil {
		t.Fatalf("NewDnfPackageManager() error = %v", err)
	}

	// Test with a package that should be installed
	if !pm.IsInstalled("dnf") {
		t.Error("IsInstalled() returned false for dnf package")
	}

	// Test with a non-existent package
	if pm.IsInstalled("non-existent-package-123456") {
		t.Error("IsInstalled() returned true for non-existent package")
	}
}

func TestDnfPackageManager_Uninstall(t *testing.T) {
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("dnf not available, skipping test")
	}

	pm, err := NewDnfPackageManager()
	if err != nil {
		t.Fatalf("NewDnfPackageManager() error = %v", err)
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