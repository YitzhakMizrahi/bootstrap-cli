package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PacmanPackageManager implements package management for Arch-based systems
type PacmanPackageManager struct {
	sudoPath string
}

// NewPacmanPackageManager creates a new Pacman package manager instance
func NewPacmanPackageManager() (*PacmanPackageManager, error) {
	sudoPath, err := exec.LookPath("sudo")
	if err != nil {
		return nil, fmt.Errorf("sudo is required but not found: %w", err)
	}

	// Verify pacman is available
	if _, err := exec.LookPath("pacman"); err != nil {
		return nil, fmt.Errorf("pacman is required but not found: %w", err)
	}

	return &PacmanPackageManager{
		sudoPath: sudoPath,
	}, nil
}

// Install installs a package using pacman
func (p *PacmanPackageManager) Install(pkg string) error {
	// Update package list first
	updateCmd := exec.Command(p.sudoPath, "pacman", "-Sy")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}

	// Install package
	cmd := exec.Command(p.sudoPath, "pacman", "-S", "--noconfirm", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (p *PacmanPackageManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("pacman", "-Q", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), pkg)
}

// Uninstall removes a package
func (p *PacmanPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(p.sudoPath, "pacman", "-R", "--noconfirm", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", pkg, err)
	}
	return nil
} 