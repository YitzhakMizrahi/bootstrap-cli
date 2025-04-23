package implementations

import (
	"os/exec"
	"testing"
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
	// Skip if not on a system with pacman
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("pacman not available, skipping test")
	}

	pm, err := NewPacmanPackageManager()
	if err != nil {
		t.Fatalf("NewPacmanPackageManager() error = %v", err)
	}

	// Test with a package that should be installed
	if !pm.IsInstalled("pacman") {
		t.Error("IsInstalled() returned false for pacman package")
	}

	// Test with a non-existent package
	if pm.IsInstalled("non-existent-package-123456") {
		t.Error("IsInstalled() returned true for non-existent package")
	}
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