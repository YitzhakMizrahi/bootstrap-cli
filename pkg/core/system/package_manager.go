package system

import (
	"fmt"
	"os/exec"
	"strings"
)

// Package represents a system package
type Package struct {
	Name         string
	Version      string
	Description  string
	Dependencies []string
}

// PackageManagerInterface defines the operations a package manager must support
type PackageManagerInterface interface {
	// Update updates the package list
	Update() error

	// Upgrade upgrades all packages
	Upgrade() error

	// Install installs a package
	Install(pkg string) error

	// InstallMany installs multiple packages
	InstallMany(pkgs []string) error

	// Remove removes a package
	Remove(pkg string) error

	// IsInstalled checks if a package is installed
	IsInstalled(pkg string) bool

	// Search searches for a package
	Search(query string) ([]Package, error)
}

// NewPackageManager creates a package manager instance for the current system
func NewPackageManager(pmType PackageManager) (PackageManagerInterface, error) {
	switch pmType {
	case Apt:
		return &aptManager{}, nil
	case Brew:
		return &brewManager{}, nil
	case Dnf:
		return &dnfManager{}, nil
	case Pacman:
		return &pacmanManager{}, nil
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", pmType)
	}
}

// Base implementation for common package manager operations
type baseManager struct{}

func (m *baseManager) runCmd(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %s\nOutput: %s", err, string(output))
	}
	return nil
}

// Apt package manager implementation
type aptManager struct {
	baseManager
}

func (m *aptManager) Update() error {
	return m.runCmd("apt", "update")
}

func (m *aptManager) Upgrade() error {
	return m.runCmd("apt", "upgrade", "-y")
}

func (m *aptManager) Install(pkg string) error {
	return m.runCmd("apt", "install", "-y", pkg)
}

func (m *aptManager) InstallMany(pkgs []string) error {
	args := append([]string{"install", "-y"}, pkgs...)
	return m.runCmd("apt", args...)
}

func (m *aptManager) Remove(pkg string) error {
	return m.runCmd("apt", "remove", "-y", pkg)
}

func (m *aptManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-l", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), pkg)
}

func (m *aptManager) Search(query string) ([]Package, error) {
	cmd := exec.Command("apt", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "/")
		if len(parts) < 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		packages = append(packages, Package{Name: name})
	}
	return packages, nil
}

// Homebrew package manager implementation
type brewManager struct {
	baseManager
}

func (m *brewManager) Update() error {
	return m.runCmd("brew", "update")
}

func (m *brewManager) Upgrade() error {
	return m.runCmd("brew", "upgrade")
}

func (m *brewManager) Install(pkg string) error {
	return m.runCmd("brew", "install", pkg)
}

func (m *brewManager) InstallMany(pkgs []string) error {
	args := append([]string{"install"}, pkgs...)
	return m.runCmd("brew", args...)
}

func (m *brewManager) Remove(pkg string) error {
	return m.runCmd("brew", "uninstall", pkg)
}

func (m *brewManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("brew", "list", pkg)
	return cmd.Run() == nil
}

func (m *brewManager) Search(query string) ([]Package, error) {
	cmd := exec.Command("brew", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if name := strings.TrimSpace(line); name != "" {
			packages = append(packages, Package{Name: name})
		}
	}
	return packages, nil
} 