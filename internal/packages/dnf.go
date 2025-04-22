package packages

import (
	"os"
	"os/exec"
	"strings"
)

// DNFManager implements PackageManager for DNF-based systems
type DNFManager struct{}

// Name returns the name of the package manager
func (d *DNFManager) Name() string {
	return string(DNF)
}

// IsAvailable checks if dnf is available on the system
func (d *DNFManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

// Install installs the given packages using dnf
func (d *DNFManager) Install(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	args := append([]string{
		"install",
		"-y",
	}, packages...)

	cmd := exec.Command("dnf", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Update updates the package list
func (d *DNFManager) Update() error {
	cmd := exec.Command("dnf", "check-update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsInstalled checks if a package is installed
func (d *DNFManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dnf", "list", "installed", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Check if package is in the output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, pkg) && !strings.Contains(line, "Available Packages") {
			return true
		}
	}

	return false
}

// Remove removes a package using dnf
func (d *DNFManager) Remove(pkg string) error {
	cmd := exec.Command("dnf", "remove", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
} 