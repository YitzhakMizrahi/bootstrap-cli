package implementations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// HomebrewPackageManager implements package management for macOS
type HomebrewPackageManager struct {
	brewPath string
}

// NewHomebrewPackageManager creates a new Homebrew package manager instance
func NewHomebrewPackageManager() (interfaces.PackageManager, error) {
	// Find brew executable
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, fmt.Errorf("brew is required but not found: %w", err)
	}

	return &HomebrewPackageManager{
		brewPath: brewPath,
	}, nil
}

// Name returns the name of the package manager
func (h *HomebrewPackageManager) Name() string {
	return string(interfaces.Homebrew)
}

// IsAvailable checks if the package manager is available on the system
func (h *HomebrewPackageManager) IsAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

// Install installs a package using Homebrew
func (h *HomebrewPackageManager) Install(pkg string) error {
	cmd := exec.Command("brew", "install", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Update updates the package list
func (h *HomebrewPackageManager) Update() error {
	cmd := exec.Command(h.brewPath, "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Upgrade upgrades all packages
func (h *HomebrewPackageManager) Upgrade() error {
	cmd := exec.Command("brew", "upgrade")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if a package is installed using Homebrew
func (h *HomebrewPackageManager) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command(h.brewPath, "list", "--formula", pkg)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}

	cmd = exec.Command(h.brewPath, "list", "--cask", pkg)
	err = cmd.Run()
	if err == nil {
		return true, nil
	}
	
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			return false, nil
		}
	}
	return false, fmt.Errorf("failed to check brew list status for %s: %w", pkg, err)
}

// Uninstall removes a package using Homebrew
func (h *HomebrewPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(h.brewPath, "uninstall", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetVersion returns the version of an installed package using Homebrew
func (h *HomebrewPackageManager) GetVersion(pkg string) (string, error) {
	cmd := exec.Command(h.brewPath, "info", "--json", pkg)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version for package %s: %w", pkg, err)
	}
	// Parse version from output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "\"version\":") {
			parts := strings.Split(line, "\"version\":")
			if len(parts) > 1 {
				version := strings.Trim(strings.TrimSpace(parts[1]), "\",")
				return version, nil
			}
		}
	}
	return "", fmt.Errorf("no version information found for package %s", pkg)
}

// ListInstalled returns a list of installed packages using Homebrew
func (h *HomebrewPackageManager) ListInstalled() ([]string, error) {
	cmd := exec.Command(h.brewPath, "list", "--formula")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}
	lines := strings.Split(string(output), "\n")
	var packages []string
	for _, line := range lines {
		if line != "" {
			packages = append(packages, strings.TrimSpace(line))
		}
	}
	return packages, nil
}

// GetName returns the name of the package manager
func (h *HomebrewPackageManager) GetName() string {
	return string(interfaces.Homebrew)
}

// SetupSpecialPackage sets up a special package that requires additional setup
func (h *HomebrewPackageManager) SetupSpecialPackage(pkg string) error {
	// For Homebrew, most packages don't require special setup
	// This method is kept for other packages that might need special repository setup
	return nil
}

// IsPackageAvailable checks if a package (formula or cask) is available via Homebrew
func (h *HomebrewPackageManager) IsPackageAvailable(pkg string) bool {
	cmd := exec.Command(h.brewPath, "info", pkg)
	err := cmd.Run()
	return err == nil
} 