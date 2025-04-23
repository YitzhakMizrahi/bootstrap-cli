package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DnfPackageManager implements package management for Fedora-based systems
type DnfPackageManager struct {
	sudoPath string
}

// NewDnfPackageManager creates a new DNF package manager instance
func NewDnfPackageManager() (*DnfPackageManager, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return nil, fmt.Errorf("sudo is required but not found: %w", err)
	}

	// Verify dnf is available
	if _, err := exec.LookPath("dnf"); err != nil {
		return nil, fmt.Errorf("dnf is required but not found: %w", err)
	}

	return &DnfPackageManager{
		sudoPath: sudoPath,
	}, nil
}

// Install installs a package using dnf
func (d *DnfPackageManager) Install(pkg string) error {
	// Update package list first
	updateCmd := exec.Command(d.sudoPath, "dnf", "check-update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		// DNF check-update returns 100 if updates are available, which is not an error
		if !strings.Contains(err.Error(), "exit status 100") {
			return fmt.Errorf("failed to check for updates: %w", err)
		}
	}

	// Install package
	cmd := exec.Command(d.sudoPath, "dnf", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (d *DnfPackageManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dnf", "list", "installed", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), pkg)
}

// Uninstall removes a package
func (d *DnfPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(d.sudoPath, "dnf", "remove", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", pkg, err)
	}
	return nil
} 