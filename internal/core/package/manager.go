package pkgmanager

import (
	"fmt"
	"os/exec"
	"strings"
)

// PackageManager handles package installation
type PackageManager struct {
	Type string
}

// New creates a new PackageManager
func New(pmType string) *PackageManager {
	return &PackageManager{
		Type: pmType,
	}
}

// Install installs a package
func (pm *PackageManager) Install(packageName string) error {
	var cmd *exec.Cmd

	switch strings.ToLower(pm.Type) {
	case "apt":
		cmd = exec.Command("apt-get", "install", "-y", packageName)
	case "dnf":
		cmd = exec.Command("dnf", "install", "-y", packageName)
	case "pacman":
		cmd = exec.Command("pacman", "-S", "--noconfirm", packageName)
	case "brew":
		cmd = exec.Command("brew", "install", packageName)
	case "choco":
		cmd = exec.Command("choco", "install", "-y", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm.Type)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package %s: %w\nOutput: %s", packageName, err, string(output))
	}

	return nil
}

// InstallMultiple installs multiple packages
func (pm *PackageManager) InstallMultiple(packages []string) error {
	for _, pkg := range packages {
		if err := pm.Install(pkg); err != nil {
			return fmt.Errorf("failed to install package %s: %w", pkg, err)
		}
	}
	return nil
}

// Update updates the package manager
func (pm *PackageManager) Update() error {
	var cmd *exec.Cmd

	switch strings.ToLower(pm.Type) {
	case "apt":
		cmd = exec.Command("apt-get", "update")
	case "dnf":
		cmd = exec.Command("dnf", "check-update")
	case "pacman":
		cmd = exec.Command("pacman", "-Sy")
	case "brew":
		cmd = exec.Command("brew", "update")
	case "choco":
		cmd = exec.Command("choco", "upgrade", "all", "-y")
	default:
		return fmt.Errorf("unsupported package manager: %s", pm.Type)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update package manager: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// IsInstalled checks if a package is installed
func (pm *PackageManager) IsInstalled(packageName string) (bool, error) {
	var cmd *exec.Cmd

	switch strings.ToLower(pm.Type) {
	case "apt":
		cmd = exec.Command("dpkg", "-l", packageName)
	case "dnf":
		cmd = exec.Command("dnf", "list", "installed", packageName)
	case "pacman":
		cmd = exec.Command("pacman", "-Q", packageName)
	case "brew":
		cmd = exec.Command("brew", "list", packageName)
	case "choco":
		cmd = exec.Command("choco", "list", "--local-only", packageName)
	default:
		return false, fmt.Errorf("unsupported package manager: %s", pm.Type)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil // Package is not installed
	}

	return strings.Contains(string(output), packageName), nil
} 