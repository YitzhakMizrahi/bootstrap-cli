package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HomebrewPackageManager implements package management for macOS
type HomebrewPackageManager struct {
	brewPath string
}

// NewHomebrewPackageManager creates a new Homebrew package manager instance
func NewHomebrewPackageManager() (*HomebrewPackageManager, error) {
	// Find brew executable
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, fmt.Errorf("brew is required but not found: %w", err)
	}

	return &HomebrewPackageManager{
		brewPath: brewPath,
	}, nil
}

// Install installs a package using Homebrew
func (h *HomebrewPackageManager) Install(pkg string) error {
	// Update package list first
	updateCmd := exec.Command(h.brewPath, "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Install package
	cmd := exec.Command(h.brewPath, "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (h *HomebrewPackageManager) IsInstalled(pkg string) bool {
	cmd := exec.Command(h.brewPath, "list", "--formula", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), pkg)
}

// Uninstall removes a package
func (h *HomebrewPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(h.brewPath, "uninstall", "--ignore-dependencies", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", pkg, err)
	}
	return nil
} 