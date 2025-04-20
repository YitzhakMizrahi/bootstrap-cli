package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// OS represents the operating system type
type OS string

// PackageManager represents a system package manager
type PackageManager string

const (
	// Operating Systems
	Linux   OS = "linux"
	MacOS   OS = "darwin"
	Windows OS = "windows"
	
	// Package Managers
	Apt       PackageManager = "apt"
	Dnf       PackageManager = "dnf"
	Homebrew  PackageManager = "brew"
	Pacman    PackageManager = "pacman"
	Zypper    PackageManager = "zypper"
	Chocolatey PackageManager = "choco"
)

// Info contains information about the current platform
type Info struct {
	OS              OS
	Distribution    string
	PackageManagers []PackageManager
}

// Detect returns information about the current platform
func Detect() (Info, error) {
	info := Info{
		OS: OS(runtime.GOOS),
	}
	
	// Detect package managers and distribution
	switch info.OS {
	case Linux:
		info.Distribution = detectLinuxDistribution()
		info.PackageManagers = detectLinuxPackageManagers()
	case MacOS:
		info.PackageManagers = detectMacPackageManagers()
	case Windows:
		info.PackageManagers = detectWindowsPackageManagers()
	}
	
	return info, nil
}

// detectLinuxDistribution tries to identify the Linux distribution
func detectLinuxDistribution() string {
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
func detectLinuxPackageManagers() []PackageManager {
	var managers []PackageManager
	
	// Check for common package managers
	pmChecks := map[PackageManager][]string{
		Apt:    {"apt", "apt-get"},
		Dnf:    {"dnf", "yum"},
		Pacman: {"pacman"},
		Zypper: {"zypper"},
	}
	
	for pm, cmds := range pmChecks {
		for _, cmd := range cmds {
			if _, err := exec.LookPath(cmd); err == nil {
				managers = append(managers, pm)
				break // Found one command for this package manager
			}
		}
	}
	
	// Check for Homebrew on Linux
	if _, err := exec.LookPath("brew"); err == nil {
		managers = append(managers, Homebrew)
	}
	
	return managers
}

// detectMacPackageManagers returns available package managers on macOS
func detectMacPackageManagers() []PackageManager {
	var managers []PackageManager
	
	// Homebrew is the primary package manager for macOS
	if _, err := exec.LookPath("brew"); err == nil {
		managers = append(managers, Homebrew)
	}
	
	return managers
}

// detectWindowsPackageManagers returns available package managers on Windows
func detectWindowsPackageManagers() []PackageManager {
	var managers []PackageManager
	
	// Check for Chocolatey
	if _, err := exec.LookPath("choco"); err == nil {
		managers = append(managers, Chocolatey)
	}
	
	return managers
}

// IsCommandAvailable checks if a command is available in PATH
func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// GetPrimaryPackageManager returns the recommended package manager for the current system
func GetPrimaryPackageManager(info Info) (PackageManager, error) {
	if len(info.PackageManagers) == 0 {
		return "", fmt.Errorf("no package managers detected")
	}
	
	// Preferred order based on OS
	switch info.OS {
	case MacOS:
		return Homebrew, nil
	case Linux:
		// Prefer Homebrew if available, then the native package manager
		for _, pm := range info.PackageManagers {
			if pm == Homebrew {
				return Homebrew, nil
			}
		}
		return info.PackageManagers[0], nil
	case Windows:
		if contains(info.PackageManagers, Chocolatey) {
			return Chocolatey, nil
		}
	}
	
	// Default to first available
	return info.PackageManagers[0], nil
}

// contains checks if a slice contains a value
func contains(slice []PackageManager, item PackageManager) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}