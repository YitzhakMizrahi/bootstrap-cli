package platform

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// DefaultDetector implements the Detector interface
type DefaultDetector struct{}

// NewDetector creates a new platform detector
func NewDetector() *DefaultDetector {
	return &DefaultDetector{}
}

// Detect returns information about the current platform
func (d *DefaultDetector) Detect() (Info, error) {
	info := Info{
		OS: d.DetectOS(),
	}
	
	// Detect package managers and distribution
	switch info.OS {
	case Linux:
		info.Distribution = d.detectLinuxDistribution()
		info.PackageManagers = d.detectLinuxPackageManagers()
	case Darwin:
		info.PackageManagers = d.detectDarwinPackageManagers()
	}
	
	return info, nil
}

// DetectOS returns the current operating system
func (d *DefaultDetector) DetectOS() OS {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return Darwin
	default:
		return Linux // Default to Linux for unsupported platforms
	}
}

// DetectPackageManagers returns available package managers for the current platform
func (d *DefaultDetector) DetectPackageManagers() []PackageManager {
	var managers []PackageManager
	os := d.DetectOS()

	switch os {
	case Linux:
		managers = d.detectLinuxPackageManagers()
	case Darwin:
		managers = d.detectDarwinPackageManagers()
	}

	return managers
}

// DetectShell returns the current shell
func (d *DefaultDetector) DetectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "/bin/bash" // Default shell
	}
	return shell
}

// IsCommandAvailable checks if a command is available in PATH
func (d *DefaultDetector) IsCommandAvailable(command string) bool {
	_, err := os.Stat(command)
	return err == nil
}

// GetPrimaryPackageManager returns the primary package manager for the current platform
func (d *DefaultDetector) GetPrimaryPackageManager(info Info) (PackageManager, error) {
	if len(info.PackageManagers) == 0 {
		return "", fmt.Errorf("no package managers detected")
	}

	// Return first detected package manager as primary
	return info.PackageManagers[0], nil
}

// detectLinuxDistribution tries to identify the Linux distribution
func (d *DefaultDetector) detectLinuxDistribution() string {
	// Check for /etc/os-release first
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			}
		}
	}
	
	// Check some common distribution files
	distros := map[string]string{
		"/etc/debian_version": "debian",
		"/etc/redhat-release": "redhat",
		"/etc/arch-release":   "arch",
		"/etc/SuSE-release":   "suse",
	}
	
	for file, distro := range distros {
		if _, err := os.Stat(file); err == nil {
			return distro
		}
	}
	
	return "unknown"
}

// detectLinuxPackageManagers returns available package managers on Linux
func (d *DefaultDetector) detectLinuxPackageManagers() []PackageManager {
	var managers []PackageManager

	// Check for apt
	if _, err := os.Stat("/usr/bin/apt"); err == nil {
		managers = append(managers, Apt)
	}

	// Check for yum
	if _, err := os.Stat("/usr/bin/yum"); err == nil {
		managers = append(managers, Yum)
	}

	// Check for pacman
	if _, err := os.Stat("/usr/bin/pacman"); err == nil {
		managers = append(managers, Pacman)
	}

	return managers
}

// detectDarwinPackageManagers returns available package managers on macOS
func (d *DefaultDetector) detectDarwinPackageManagers() []PackageManager {
	var managers []PackageManager

	// Check for Homebrew
	if _, err := os.Stat("/usr/local/bin/brew"); err == nil {
		managers = append(managers, Homebrew)
	}

	return managers
}

func (d *DefaultDetector) detectPackageManagers(os OS) ([]PackageManager, error) {
	var pms []PackageManager

	switch os {
	case Linux:
		// Check for Linux package managers
		if d.IsCommandAvailable("apt") || d.IsCommandAvailable("apt-get") {
			pms = append(pms, Apt)
		}
		if d.IsCommandAvailable("yum") {
			pms = append(pms, Yum)
		}
		if d.IsCommandAvailable("dnf") {
			pms = append(pms, Dnf)
		}
		if d.IsCommandAvailable("pacman") {
			pms = append(pms, Pacman)
		}
	case Darwin:
		// Check for Homebrew on macOS
		if d.IsCommandAvailable("brew") {
			pms = append(pms, Homebrew)
		}
	}

	return pms, nil
} 