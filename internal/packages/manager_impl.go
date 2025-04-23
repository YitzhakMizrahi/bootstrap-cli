package packages

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// execCommand is a function type that matches exec.Command
type execCommand func(name string, arg ...string) *exec.Cmd

// packageManager implements the Manager interface
type packageManager struct {
	system  string
	cmd     string
	logger  *log.Logger
	execCmd execCommand // Function to execute commands, can be mocked in tests
}

// NewPackageManager creates a new package manager for the given system
func NewPackageManager(system string) (Manager, error) {
	var cmd string
	switch system {
	case "ubuntu", "debian":
		cmd = "apt-get"
	case "fedora":
		cmd = "dnf"
	case "arch":
		cmd = "pacman"
	default:
		return nil, fmt.Errorf("unsupported system: %s", system)
	}

	return &packageManager{
		system:  system,
		cmd:     cmd,
		logger:  log.New(log.InfoLevel),
		execCmd: exec.Command,
	}, nil
}

// Install installs a package
func (pm *packageManager) Install(packageName string) error {
	pm.logger.Info("Installing package: %s", packageName)
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd(pm.cmd, "install", "-y", packageName)
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "install", "-y", packageName)
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-S", "--noconfirm", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to install package %s: %v", packageName, err)
		return fmt.Errorf("command execution failed: %s %s, output: %s, error: %w", pm.cmd, packageName, output, err)
	}
	return nil
}

// Uninstall removes a package
func (pm *packageManager) Uninstall(packageName string) error {
	pm.logger.Info("Uninstalling package: %s", packageName)
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd(pm.cmd, "remove", "-y", packageName)
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "remove", "-y", packageName)
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-R", "--noconfirm", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to uninstall package %s: %v", packageName, err)
		return fmt.Errorf("command execution failed: %s %s, output: %s, error: %w", pm.cmd, packageName, output, err)
	}
	return nil
}

// Update updates a package
func (pm *packageManager) Update(packageName string) error {
	pm.logger.Info("Updating package: %s", packageName)
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd(pm.cmd, "upgrade", "-y", packageName)
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "upgrade", "-y", packageName)
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-Syu", "--noconfirm", packageName)
	default:
		return fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to update package %s: %v", packageName, err)
		return fmt.Errorf("command execution failed: %s %s, output: %s, error: %w", pm.cmd, packageName, output, err)
	}
	return nil
}

// IsInstalled checks if a package is installed
func (pm *packageManager) IsInstalled(packageName string) (bool, error) {
	pm.logger.Debug("Checking if package is installed: %s", packageName)
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd("dpkg", "-l", packageName)
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "list", "installed", packageName)
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-Q", packageName)
	default:
		return false, fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to check if package %s is installed: %v", packageName, err)
		return false, fmt.Errorf("command execution failed: %s %s, output: %s, error: %w", pm.cmd, packageName, output, err)
	}
	return strings.Contains(string(output), packageName), nil
}

// ListInstalled returns a list of installed packages
func (pm *packageManager) ListInstalled() ([]string, error) {
	pm.logger.Debug("Listing installed packages")
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd("dpkg", "-l")
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "list", "installed")
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-Q")
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to list installed packages: %v", err)
		return nil, fmt.Errorf("command execution failed: %s, output: %s, error: %w", pm.cmd, output, err)
	}

	// Parse output to get package names
	lines := strings.Split(string(output), "\n")
	var packages []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			packages = append(packages, strings.Fields(line)[0])
		}
	}
	return packages, nil
}

// GetVersion returns the version of a package
func (pm *packageManager) GetVersion(packageName string) (string, error) {
	pm.logger.Debug("Getting version for package: %s", packageName)
	var cmd *exec.Cmd
	switch pm.cmd {
	case "apt-get":
		cmd = pm.execCmd("dpkg", "-s", packageName)
	case "dnf":
		cmd = pm.execCmd(pm.cmd, "info", packageName)
	case "pacman":
		cmd = pm.execCmd(pm.cmd, "-Q", packageName)
	default:
		return "", fmt.Errorf("unsupported package manager: %s", pm.cmd)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		pm.logger.Error("Failed to get version for package %s: %v", packageName, err)
		return "", fmt.Errorf("command execution failed: %s %s, output: %s, error: %w", pm.cmd, packageName, output, err)
	}

	// Parse output to get version
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Version:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Version:")), nil
		}
	}
	return "", fmt.Errorf("version not found in output")
} 