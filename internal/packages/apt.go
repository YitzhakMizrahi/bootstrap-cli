package packages

import (
	"os"
	"os/exec"
	"strings"
)

// APTManager implements PackageManager for APT-based systems
type APTManager struct{}

// Name returns the name of the package manager
func (a *APTManager) Name() string {
	return string(APT)
}

// IsAvailable checks if apt is available on the system
func (a *APTManager) IsAvailable() bool {
	_, err := exec.LookPath("apt")
	return err == nil
}

// Install installs the given packages using apt
func (a *APTManager) Install(packages ...string) error {
	if len(packages) == 0 {
		return nil
	}

	args := append([]string{
		"install",
		"-y",
	}, packages...)

	cmd := exec.Command("apt", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Update updates the package list
func (a *APTManager) Update() error {
	cmd := exec.Command("apt", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// IsInstalled checks if a package is installed
func (a *APTManager) IsInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-l", pkg)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Check if package is installed (ii at the start of the line)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ii  "+pkg) {
			return true
		}
	}

	return false
}

// Remove removes a package using apt
func (a *APTManager) Remove(pkg string) error {
	cmd := exec.Command("apt", "remove", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
} 