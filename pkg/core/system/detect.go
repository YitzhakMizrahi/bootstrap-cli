package system

import (
	"os"
	"os/exec"
	"runtime"
)

// Type represents the type of Unix-like system
type Type string

const (
	Linux  Type = "linux"
	Darwin Type = "darwin"
	WSL    Type = "wsl"
)

// PackageManager represents the system's package manager
type PackageManager string

const (
	Apt    PackageManager = "apt"
	Brew   PackageManager = "brew"
	Dnf    PackageManager = "dnf"
	Pacman PackageManager = "pacman"
)

// Info contains information about the system
type Info struct {
	Type           Type
	PackageManager PackageManager
	Shell          string
	HomeDir        string
	ConfigDir      string
}

// Detect gathers system information
func Detect() (*Info, error) {
	info := &Info{}

	// Detect system type
	switch runtime.GOOS {
	case "darwin":
		info.Type = Darwin
	case "linux":
		if isWSL() {
			info.Type = WSL
		} else {
			info.Type = Linux
		}
	default:
		return nil, ErrUnsupportedOS
	}

	// Detect package manager
	pm, err := detectPackageManager()
	if err != nil {
		return nil, err
	}
	info.PackageManager = pm

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	info.HomeDir = home

	// Get current shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash" // Default to bash
	}
	info.Shell = shell

	// Set config directory
	info.ConfigDir = os.ExpandEnv("$HOME/.config/bootstrap-cli")

	return info, nil
}

// isWSL checks if running under Windows Subsystem for Linux
func isWSL() bool {
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return true
	}
	
	// Check for WSL-specific environment variable
	if os.Getenv("WSL_DISTRO_NAME") != "" {
		return true
	}

	return false
}

// detectPackageManager identifies the system's package manager
func detectPackageManager() (PackageManager, error) {
	// For macOS, prefer Homebrew
	if runtime.GOOS == "darwin" {
		if _, err := exec.LookPath("brew"); err == nil {
			return Brew, nil
		}
	}

	// Check for various package managers
	packageManagers := []struct {
		path string
		pm   PackageManager
	}{
		{"/usr/bin/apt", Apt},
		{"/usr/bin/dnf", Dnf},
		{"/usr/bin/pacman", Pacman},
	}

	for _, pm := range packageManagers {
		if _, err := os.Stat(pm.path); err == nil {
			return pm.pm, nil
		}
	}

	// Default to apt for Debian-based systems
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return Apt, nil
	}

	return "", ErrNoPackageManager
}

// HasCommand checks if a command is available in PATH
func HasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
} 