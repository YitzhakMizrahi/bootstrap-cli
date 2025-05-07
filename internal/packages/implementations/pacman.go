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
func (p *PacmanPackageManager) Install(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if a package is installed using Pacman
func (p *PacmanPackageManager) IsInstalled(pkg string) (bool, error) {
	cmd := exec.Command(p.sudoPath, "pacman", "-Q", pkg)
	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 { // Pacman exits 1 if package not found
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to check pacman installed status for %s: %w", pkg, err)
	}
	return true, nil // Exit code 0 means installed
}

// IsPackageAvailable checks if a specific package is available in Pacman repositories
func (p *PacmanPackageManager) IsPackageAvailable(pkg string) bool {
	cmd := exec.Command(p.sudoPath, "pacman", "-Si", pkg)
	err := cmd.Run()
	return err == nil
}

// Uninstall removes a package using Pacman (Renamed from Remove)
func (p *PacmanPackageManager) Uninstall(pkg string) error {
	cmd := exec.Command(p.sudoPath, "pacman", "-Rns", "--noconfirm", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetVersion returns the version of an installed package
func (p *PacmanPackageManager) GetVersion(pkg string) (string, error) {
	installed, err := p.IsInstalled(pkg)
	if err != nil {
		return "", fmt.Errorf("failed to check if package %s is installed: %w", pkg, err)
	}
	if !installed {
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
		installed, err := p.IsInstalled("yay")
		if err != nil {
			return fmt.Errorf("failed to check if yay is installed: %w", err)
		}
		if !installed {
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