package implementations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// PacmanPackageManager implements package management for Arch-based systems
type PacmanPackageManager struct {
	sudoPath string
}

// NewPacmanPackageManager creates a new Pacman package manager instance
func NewPacmanPackageManager() (interfaces.PackageManager, error) {
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

// Name returns the name of the package manager
func (p *PacmanPackageManager) Name() string {
	return string(interfaces.Pacman)
}

// GetName returns the name of the package manager
func (p *PacmanPackageManager) GetName() string {
	return string(interfaces.Pacman)
}

// IsAvailable checks if pacman is available on the system
func (p *PacmanPackageManager) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

// Update updates the package list
func (p *PacmanPackageManager) Update() error {
	cmd := exec.Command(p.sudoPath, "pacman", "-Sy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update package list: %w", err)
	}
	return nil
}

// Install installs a package using pacman
func (p *PacmanPackageManager) Install(packageName string) error {
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

// Remove removes a package
func (p *PacmanPackageManager) Remove(packageName string) error {
	cmd := exec.Command(p.sudoPath, "pacman", "-R", "--noconfirm", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetVersion returns the version of an installed package
func (p *PacmanPackageManager) GetVersion(pkg string) (string, error) {
	if !p.IsInstalled(pkg) {
		return "", fmt.Errorf("package %s is not installed", pkg)
	}

	cmd := exec.Command("pacman", "-Q", pkg)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version for package %s: %w", pkg, err)
	}

	// Output format: package_name version
	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return "", fmt.Errorf("unexpected output format for package %s", pkg)
	}

	return parts[1], nil
}

// ListInstalled returns a list of installed packages
func (p *PacmanPackageManager) ListInstalled() ([]string, error) {
	cmd := exec.Command("pacman", "-Q")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed packages: %w", err)
	}

	var packages []string
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}

	return packages, nil
}

// SetupSpecialPackage sets up any special repository requirements for a package
func (p *PacmanPackageManager) SetupSpecialPackage(pkg string) error {
	// For Pacman, we might need to enable additional repositories from the AUR
	// This is a placeholder implementation that can be extended based on specific package requirements
	switch pkg {
	case "yay":
		// Install yay from AUR if not already installed
		if !p.IsInstalled("yay") {
			// First ensure base-devel is installed
			cmd := exec.Command(p.sudoPath, "pacman", "-S", "--noconfirm", "base-devel", "git")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install base-devel: %w", err)
			}

			// Clone and install yay
			tempDir, err := os.MkdirTemp("", "yay-install")
			if err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}
			defer os.RemoveAll(tempDir)

			cmd = exec.Command("git", "clone", "https://aur.archlinux.org/yay.git", tempDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to clone yay: %w", err)
			}

			cmd = exec.Command("makepkg", "-si", "--noconfirm")
			cmd.Dir = tempDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to build and install yay: %w", err)
			}
		}
		return nil
	default:
		return nil // No special setup needed for this package
	}
}

// Upgrade upgrades all packages
func (p *PacmanPackageManager) Upgrade() error {
	cmd := exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsPackageAvailable checks if a package is available in the package manager's repositories
func (p *PacmanPackageManager) IsPackageAvailable(pkg string) bool {
	cmd := exec.Command("pacman", "-Ss", "^" + pkg + "$")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
} 