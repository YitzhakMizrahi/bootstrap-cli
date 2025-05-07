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

// GetName returns the name of the package manager
func (d *DnfPackageManager) GetName() string {
	return string(interfaces.DNF)
}

// IsAvailable checks if the package manager is available on the system
func (d *DnfPackageManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

// Install installs a package using dnf
func (d *DnfPackageManager) Install(packageName string) error {
	cmd := exec.Command("sudo", "dnf", "install", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Update updates the package list
func (d *DnfPackageManager) Update() error {
	cmd := exec.Command(d.sudoPath, "dnf", "check-update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if a package is installed using dnf
func (d *DnfPackageManager) IsInstalled(packageName string) (bool, error) {
	cmd := exec.Command(d.sudoPath, "dnf", "list", "installed", packageName)
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to check dnf installed status for %s: %w", packageName, err)
	}
	return true, nil
}

// IsPackageAvailable checks if a specific package is available in dnf repositories
func (d *DnfPackageManager) IsPackageAvailable(packageName string) bool {
	cmd := exec.Command(d.sudoPath, "dnf", "list", "available", packageName)
	err := cmd.Run()
	return err == nil
}

// Upgrade upgrades all packages using dnf
func (d *DnfPackageManager) Upgrade() error {
	cmd := exec.Command("sudo", "dnf", "upgrade", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Uninstall removes a package using dnf (Renamed from Remove)
func (d *DnfPackageManager) Uninstall(packageName string) error {
	cmd := exec.Command(d.sudoPath, "dnf", "remove", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove package %s: %w", packageName, err)
	}
	return nil
}

// GetVersion returns the version of an installed package using dnf
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

// ListInstalled returns a list of installed packages using dnf
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

// SetupSpecialPackage for dnf (if any)
func (d *DnfPackageManager) SetupSpecialPackage(packageName string) error {
	switch packageName {
	case "docker":
		cmd := exec.Command(d.sudoPath, "dnf", "config-manager", "--add-repo", "https://download.docker.com/linux/fedora/docker-ce.repo")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to enable Docker repository: %w", err)
		}
		return nil
	default:
		return nil
	}
} 