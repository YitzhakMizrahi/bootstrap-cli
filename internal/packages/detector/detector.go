package detector

import (
	"os/exec"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// DetectPackageManager determines the system's package manager type
func DetectPackageManager() (interfaces.PackageManagerType, error) {
	// Check for apt (Debian/Ubuntu)
	if _, err := exec.LookPath("apt"); err == nil {
		return interfaces.APT, nil
	}

	// Check for dnf (Fedora)
	if _, err := exec.LookPath("dnf"); err == nil {
		return interfaces.DNF, nil
	}

	// Check for pacman (Arch)
	if _, err := exec.LookPath("pacman"); err == nil {
		return interfaces.Pacman, nil
	}

	// Check for Homebrew
	if _, err := exec.LookPath("brew"); err == nil {
		return interfaces.Homebrew, nil
	}

	return "", nil
} 