package pkgmanager

import (
	"fmt"
	"os/exec"

	"github.com/YitzhakMizrahi/bootstrap-cli/pkg/platform"
)

// Manager defines the interface for package management operations
type Manager interface {
	Install(pkg string) error
	Uninstall(pkg string) error
	Update() error
	IsInstalled(pkg string) bool
}

// DefaultManager implements the Manager interface using the system's default package manager
type DefaultManager struct {
	pm platform.PackageManager
}

// NewManager creates a new package manager instance
func NewManager(pm platform.PackageManager) Manager {
	return &DefaultManager{
		pm: pm,
	}
}

// Install installs a package using the system's package manager
func (m *DefaultManager) Install(pkg string) error {
	if m.IsInstalled(pkg) {
		return nil
	}

	var cmd *exec.Cmd
	switch m.pm {
	case platform.Homebrew:
		cmd = exec.Command("brew", "install", pkg)
	case platform.Apt:
		cmd = exec.Command("sudo", "apt-get", "install", "-y", pkg)
		cmd.Env = append(cmd.Env, "DEBIAN_FRONTEND=noninteractive")
	case platform.Yum:
		cmd = exec.Command("sudo", "dnf", "install", "-y", pkg)
	case platform.Pacman:
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", pkg)
	default:
		return fmt.Errorf("unsupported package manager: %s", m.pm)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s: %w", pkg, err)
	}

	return nil
}

// Uninstall removes a package using the system's package manager
func (m *DefaultManager) Uninstall(pkg string) error {
	if !m.IsInstalled(pkg) {
		return nil
	}

	var cmd *exec.Cmd
	switch m.pm {
	case platform.Homebrew:
		cmd = exec.Command("brew", "uninstall", pkg)
	case platform.Apt:
		cmd = exec.Command("sudo", "apt-get", "remove", "-y", pkg)
	case platform.Yum:
		cmd = exec.Command("sudo", "dnf", "remove", "-y", pkg)
	case platform.Pacman:
		cmd = exec.Command("sudo", "pacman", "-R", "--noconfirm", pkg)
	default:
		return fmt.Errorf("unsupported package manager: %s", m.pm)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to uninstall package %s: %w", pkg, err)
	}

	return nil
}

// Update updates the package manager's package lists
func (m *DefaultManager) Update() error {
	var cmd *exec.Cmd
	switch m.pm {
	case platform.Homebrew:
		cmd = exec.Command("brew", "update")
	case platform.Apt:
		cmd = exec.Command("sudo", "apt-get", "update")
	case platform.Yum:
		cmd = exec.Command("sudo", "dnf", "check-update")
	case platform.Pacman:
		cmd = exec.Command("sudo", "pacman", "-Sy")
	default:
		return fmt.Errorf("unsupported package manager: %s", m.pm)
	}

	if err := cmd.Run(); err != nil {
		// DNF returns 100 when updates are available
		if m.pm == platform.Yum && cmd.ProcessState.ExitCode() == 100 {
			return nil
		}
		return fmt.Errorf("failed to update package lists: %w", err)
	}

	return nil
}

// IsInstalled checks if a package is installed
func (m *DefaultManager) IsInstalled(pkg string) bool {
	var cmd *exec.Cmd
	switch m.pm {
	case platform.Homebrew:
		cmd = exec.Command("brew", "list", pkg)
	case platform.Apt:
		cmd = exec.Command("dpkg", "-l", pkg)
	case platform.Yum:
		cmd = exec.Command("rpm", "-q", pkg)
	case platform.Pacman:
		cmd = exec.Command("pacman", "-Q", pkg)
	default:
		return false
	}

	return cmd.Run() == nil
} 