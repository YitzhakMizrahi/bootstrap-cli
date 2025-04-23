package packages

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// packageManager implements the Manager interface
type packageManager struct {
	system string
	cmd    string
	logger *log.Logger
}

// NewPackageManager creates a new package manager instance
func NewPackageManager(system string) (Manager, error) {
	cmd, err := getPackageManagerCmd(system)
	if err != nil {
		return nil, NewSystemNotSupportedError(system)
	}
	
	return &packageManager{
		system: system,
		cmd:    cmd,
		logger: log.New(log.InfoLevel),
	}, nil
}

// Install implements Manager.Install
func (pm *packageManager) Install(packageName string) error {
	pm.logger.Info(fmt.Sprintf("Installing package: %s", packageName))
	
	cmd := exec.Command(pm.cmd, "install", "-y", packageName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to install package %s: %v", packageName, err))
		return NewCommandExecutionError(
			fmt.Sprintf("%s install %s", pm.cmd, packageName),
			stderr.String(),
			err,
		)
	}
	
	pm.logger.Info(fmt.Sprintf("Successfully installed package: %s", packageName))
	return nil
}

// Uninstall implements Manager.Uninstall
func (pm *packageManager) Uninstall(packageName string) error {
	pm.logger.Info(fmt.Sprintf("Uninstalling package: %s", packageName))
	
	cmd := exec.Command(pm.cmd, "remove", "-y", packageName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to uninstall package %s: %v", packageName, err))
		return NewCommandExecutionError(
			fmt.Sprintf("%s remove %s", pm.cmd, packageName),
			stderr.String(),
			err,
		)
	}
	
	pm.logger.Info(fmt.Sprintf("Successfully uninstalled package: %s", packageName))
	return nil
}

// Update implements Manager.Update
func (pm *packageManager) Update(packageName string) error {
	pm.logger.Info(fmt.Sprintf("Updating package: %s", packageName))
	
	cmd := exec.Command(pm.cmd, "update", "-y", packageName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to update package %s: %v", packageName, err))
		return NewCommandExecutionError(
			fmt.Sprintf("%s update %s", pm.cmd, packageName),
			stderr.String(),
			err,
		)
	}
	
	pm.logger.Info(fmt.Sprintf("Successfully updated package: %s", packageName))
	return nil
}

// IsInstalled implements Manager.IsInstalled
func (pm *packageManager) IsInstalled(packageName string) (bool, error) {
	pm.logger.Debug(fmt.Sprintf("Checking if package is installed: %s", packageName))
	
	cmd := exec.Command(pm.cmd, "list", "--installed", packageName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to check if package %s is installed: %v", packageName, err))
		return false, NewCommandExecutionError(
			fmt.Sprintf("%s list --installed %s", pm.cmd, packageName),
			stderr.String(),
			err,
		)
	}
	
	isInstalled := strings.Contains(stdout.String(), packageName)
	pm.logger.Debug(fmt.Sprintf("Package %s is installed: %v", packageName, isInstalled))
	return isInstalled, nil
}

// ListInstalled implements Manager.ListInstalled
func (pm *packageManager) ListInstalled() ([]string, error) {
	pm.logger.Debug("Listing installed packages")
	
	cmd := exec.Command(pm.cmd, "list", "--installed")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to list installed packages: %v", err))
		return nil, NewCommandExecutionError(
			fmt.Sprintf("%s list --installed", pm.cmd),
			stderr.String(),
			err,
		)
	}
	
	packages := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	var result []string
	for _, pkg := range packages {
		if pkg != "" {
			result = append(result, pkg)
		}
	}
	
	pm.logger.Debug(fmt.Sprintf("Found %d installed packages", len(result)))
	return result, nil
}

// GetVersion implements Manager.GetVersion
func (pm *packageManager) GetVersion(packageName string) (string, error) {
	pm.logger.Debug(fmt.Sprintf("Getting version for package: %s", packageName))
	
	cmd := exec.Command(pm.cmd, "show", packageName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		pm.logger.Error(fmt.Sprintf("Failed to get version for package %s: %v", packageName, err))
		return "", NewCommandExecutionError(
			fmt.Sprintf("%s show %s", pm.cmd, packageName),
			stderr.String(),
			err,
		)
	}
	
	// Extract version from output
	output := stdout.String()
	version := ""
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "Version: ") {
			version = strings.TrimPrefix(line, "Version: ")
			break
		}
	}
	
	if version == "" {
		pm.logger.Error(fmt.Sprintf("Could not find version for package %s", packageName))
		return "", NewPackageNotFoundError(packageName)
	}
	
	pm.logger.Debug(fmt.Sprintf("Version for package %s: %s", packageName, version))
	return version, nil
}

// getPackageManagerCmd returns the appropriate package manager command for the system
func getPackageManagerCmd(system string) (string, error) {
	switch system {
	case "ubuntu", "debian":
		return "apt-get", nil
	case "fedora", "rhel":
		return "dnf", nil
	case "arch":
		return "pacman", nil
	default:
		return "", NewSystemNotSupportedError(system)
	}
} 