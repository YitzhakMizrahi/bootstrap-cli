package implementations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// APTManager implements the PackageManager interface for APT-based systems
type APTManager struct {
	aptGetPath string
}

// NewAptPackageManager creates a new APT package manager instance
func NewAptPackageManager() (interfaces.PackageManager, error) {
	// Check for apt-get
	aptGetPath, err := exec.LookPath("apt-get")
	if err != nil {
		return nil, fmt.Errorf("apt-get is required but not found: %w", err)
	}

	return &APTManager{
		aptGetPath: aptGetPath,
	}, nil
}

// Name returns the name of the package manager
func (a *APTManager) Name() string {
	return "apt"
}

// IsAvailable checks if apt-get is available on the system
func (a *APTManager) IsAvailable() bool {
	_, err := exec.LookPath(a.aptGetPath)
	return err == nil
}

// Update updates the package list
func (a *APTManager) Update() error {
	cmd := exec.Command(a.aptGetPath, "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs one or more packages
func (a *APTManager) Install(packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command(a.aptGetPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Remove removes a package
func (a *APTManager) Remove(packageName string) error {
	cmd := exec.Command(a.aptGetPath, "remove", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if a package is installed
func (a *APTManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), packageName)
}

// GetVersion returns the version of an installed package
func (a *APTManager) GetVersion(packageName string) (string, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Version}", packageName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// ListInstalled returns a list of installed packages
func (a *APTManager) ListInstalled() ([]string, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package}\n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
} 