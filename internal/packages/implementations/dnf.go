package implementations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// DnfPackageManager implements package management for Fedora-based systems
type DnfPackageManager struct {
	sudoPath string
}

// NewDnfPackageManager creates a new DNF package manager instance
func NewDnfPackageManager() (interfaces.PackageManager, error) {
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

// Name returns the name of the package manager
func (d *DnfPackageManager) Name() string {
	return string(interfaces.DNF)
}

// IsAvailable checks if the package manager is available on the system
func (d *DnfPackageManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

// Install installs packages using dnf
func (d *DnfPackageManager) Install(packages ...string) error {
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

	// Install packages
	args := append([]string{"dnf", "install", "-y"}, packages...)
	cmd := exec.Command(d.sudoPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packages: %w", err)
	}
	return nil
}

// Update updates the package list
func (d *DnfPackageManager) Update() error {
	cmd := exec.Command(d.sudoPath, "dnf", "check-update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

// Remove removes a package
func (d *DnfPackageManager) Remove(pkg string) error {
	cmd := exec.Command(d.sudoPath, "dnf", "remove", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove package %s: %w", pkg, err)
	}
	return nil
}

// GetVersion returns the version of a package
func (d *DnfPackageManager) GetVersion(packageName string) (string, error) {
	cmd := exec.Command("dnf", "list", "installed", packageName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version for package %s: %w", packageName, err)
	}
	// Parse version from output
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("no version information found for package %s", packageName)
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid version format for package %s", packageName)
	}
	return fields[1], nil
}

// ListInstalled returns a list of installed packages
func (d *DnfPackageManager) ListInstalled() ([]string, error) {
	cmd := exec.Command("dnf", "list", "installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("no packages found")
	}
	var packages []string
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			packages = append(packages, fields[0])
		}
	}
	return packages, nil
} 