package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// AptPackageManager implements package management for Debian-based systems
type AptPackageManager struct {
	sudoPath string
}

// NewAptPackageManager creates a new APT package manager instance
func NewAptPackageManager() (*AptPackageManager, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return nil, fmt.Errorf("sudo is required but not found: %w", err)
	}
	return &AptPackageManager{
		sudoPath: sudoPath,
	}, nil
}

// Install installs a package using apt-get
func (a *AptPackageManager) Install(pkg string) error {
	// Update package list first
	updateCmd := exec.Command(a.sudoPath, "apt-get", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Install package
	cmd := exec.Command(a.sudoPath, "apt-get", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (a *AptPackageManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-l", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), pkg)
}

// Uninstall removes a package
func (a *AptPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(a.sudoPath, "apt-get", "remove", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", pkg, err)
	}
	return nil
} 