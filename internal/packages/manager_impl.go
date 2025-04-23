package packages

import (
	"fmt"
	"os/exec"
	"strings"
)

// packageManager implements the Manager interface
type packageManager struct {
	system string
	cmd    string
}

// NewPackageManager creates a new package manager instance
func NewPackageManager(system string) (Manager, error) {
	cmd, err := getPackageManagerCmd(system)
	if err != nil {
		return nil, err
	}
	
	return &packageManager{
		system: system,
		cmd:    cmd,
	}, nil
}

// Install implements Manager.Install
func (pm *packageManager) Install(packageName string) error {
	cmd := exec.Command(pm.cmd, "install", packageName)
	return cmd.Run()
}

// Uninstall implements Manager.Uninstall
func (pm *packageManager) Uninstall(packageName string) error {
	cmd := exec.Command(pm.cmd, "uninstall", packageName)
	return cmd.Run()
}

// Update implements Manager.Update
func (pm *packageManager) Update(packageName string) error {
	cmd := exec.Command(pm.cmd, "update", packageName)
	return cmd.Run()
}

// IsInstalled implements Manager.IsInstalled
func (pm *packageManager) IsInstalled(packageName string) (bool, error) {
	cmd := exec.Command(pm.cmd, "list", "--installed", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.Contains(string(output), packageName), nil
}

// ListInstalled implements Manager.ListInstalled
func (pm *packageManager) ListInstalled() ([]string, error) {
	cmd := exec.Command(pm.cmd, "list", "--installed")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	packages := strings.Split(string(output), "\n")
	var result []string
	for _, pkg := range packages {
		if pkg != "" {
			result = append(result, pkg)
		}
	}
	return result, nil
}

// GetVersion implements Manager.GetVersion
func (pm *packageManager) GetVersion(packageName string) (string, error) {
	cmd := exec.Command(pm.cmd, "version", packageName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
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
		return "", fmt.Errorf("unsupported system: %s", system)
	}
} 