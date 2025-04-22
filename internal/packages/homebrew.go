package packages

import (
	"os"
	"os/exec"
	"strings"
)

// HomebrewManager implements PackageManager for Homebrew
type HomebrewManager struct{}

// Name returns the name of the package manager
func (h *HomebrewManager) Name() string {
	return string(Homebrew)
}

// IsAvailable checks if brew is available on the system
func (h *HomebrewManager) IsAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

// Install installs the given packages using brew
func (h *HomebrewManager) Install(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	args := append([]string{
		"install",
	}, packages...)

	cmd := exec.Command("brew", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Update updates the package list and upgrades all packages
func (h *HomebrewManager) Update() error {
	cmd := exec.Command("brew", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// IsInstalled checks if a package is installed
func (h *HomebrewManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("brew", "list", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), pkg)
} 