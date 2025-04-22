package packages

import (
	"fmt"
	"runtime"
)

// PackageManager defines the interface for different package managers
type PackageManager interface {
	// Name returns the name of the package manager
	Name() string
	// IsAvailable checks if the package manager is available on the system
	IsAvailable() bool
	// Install installs the given packages
	Install(packages ...string) error
	// Update updates the package list
	Update() error
	// IsInstalled checks if a package is installed
	IsInstalled(pkg string) bool
	// Remove removes a package
	Remove(pkg string) error
}

// Type represents the type of package manager
type Type string

const (
	// APT package manager (Debian/Ubuntu)
	APT Type = "apt"
	// DNF package manager (Fedora)
	DNF Type = "dnf"
	// Pacman package manager (Arch)
	Pacman Type = "pacman"
	// Homebrew package manager (macOS)
	Homebrew Type = "brew"
)

// ErrPackageManagerNotFound is returned when no suitable package manager is found
var ErrPackageManagerNotFound = fmt.Errorf("no suitable package manager found")

// DetectPackageManager returns the appropriate package manager for the current system
func DetectPackageManager() (PackageManager, error) {
	// Check for macOS first
	if runtime.GOOS == "darwin" {
		brew := &HomebrewManager{}
		if brew.IsAvailable() {
			return brew, nil
		}
	}

	// Check for Linux package managers
	if runtime.GOOS == "linux" {
		// Try APT first (Debian/Ubuntu)
		apt := &APTManager{}
		if apt.IsAvailable() {
			return apt, nil
		}

		// Try DNF (Fedora)
		dnf := &DNFManager{}
		if dnf.IsAvailable() {
			return dnf, nil
		}

		// Try Pacman (Arch)
		pacman := &PacmanManager{}
		if pacman.IsAvailable() {
			return pacman, nil
		}
	}

	return nil, ErrPackageManagerNotFound
} 