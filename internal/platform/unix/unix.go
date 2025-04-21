package unix

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/platform"
)

// UnixPlatform implements the Platform interface for Unix-like systems
type UnixPlatform struct {
	packageManager string
	osInfo        *platform.OSInfo
}

// NewUnixPlatform creates a new Unix platform implementation
func NewUnixPlatform() *UnixPlatform {
	return &UnixPlatform{
		packageManager: detectPackageManager(),
		osInfo: &platform.OSInfo{
			Type:    runtime.GOOS,
			Version: getOSVersion(),
			Arch:    runtime.GOARCH,
		},
	}
}

// InstallPackage installs a package using the system package manager
func (p *UnixPlatform) InstallPackage(name string) error {
	var cmd *exec.Cmd
	switch p.packageManager {
	case "apt":
		cmd = exec.Command("sudo", "apt", "install", "-y", name)
	case "yum":
		cmd = exec.Command("sudo", "yum", "install", "-y", name)
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "install", "-y", name)
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", name)
	default:
		return fmt.Errorf("unsupported package manager: %s", p.packageManager)
	}
	return cmd.Run()
}

// UninstallPackage removes a package using the system package manager
func (p *UnixPlatform) UninstallPackage(name string) error {
	var cmd *exec.Cmd
	switch p.packageManager {
	case "apt":
		cmd = exec.Command("sudo", "apt", "remove", "-y", name)
	case "yum":
		cmd = exec.Command("sudo", "yum", "remove", "-y", name)
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "remove", "-y", name)
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-R", "--noconfirm", name)
	default:
		return fmt.Errorf("unsupported package manager: %s", p.packageManager)
	}
	return cmd.Run()
}

// IsPackageInstalled checks if a package is installed
func (p *UnixPlatform) IsPackageInstalled(name string) bool {
	var cmd *exec.Cmd
	switch p.packageManager {
	case "apt":
		cmd = exec.Command("dpkg", "-l", name)
	case "yum", "dnf":
		cmd = exec.Command("rpm", "-q", name)
	case "pacman":
		cmd = exec.Command("pacman", "-Qi", name)
	default:
		return false
	}
	return cmd.Run() == nil
}

// GetPackageManager returns the detected package manager
func (p *UnixPlatform) GetPackageManager() string {
	return p.packageManager
}

// GetDefaultShell returns the default shell for the current user
func (p *UnixPlatform) GetDefaultShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "/bin/bash"
	}
	return shell
}

// GetShellPath returns the path to the specified shell
func (p *UnixPlatform) GetShellPath(shell string) string {
	paths := []string{
		"/bin/" + shell,
		"/usr/bin/" + shell,
		"/usr/local/bin/" + shell,
	}
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// IsShellAvailable checks if a shell is available on the system
func (p *UnixPlatform) IsShellAvailable(shell string) bool {
	return p.GetShellPath(shell) != ""
}

// GetHomeDir returns the user's home directory
func (p *UnixPlatform) GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return home
}

// GetConfigDir returns the user's configuration directory
func (p *UnixPlatform) GetConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(p.GetHomeDir(), ".config")
}

// GetCacheDir returns the user's cache directory
func (p *UnixPlatform) GetCacheDir() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return xdg
	}
	return filepath.Join(p.GetHomeDir(), ".cache")
}

// IsRoot checks if the current user has root privileges
func (p *UnixPlatform) IsRoot() bool {
	return os.Geteuid() == 0
}

// ElevatePrivileges attempts to elevate privileges using sudo
func (p *UnixPlatform) ElevatePrivileges() error {
	if p.IsRoot() {
		return nil
	}
	cmd := exec.Command("sudo", "-v")
	return cmd.Run()
}

// GetOSInfo returns information about the operating system
func (p *UnixPlatform) GetOSInfo() *platform.OSInfo {
	return p.osInfo
}

// detectPackageManager detects the system package manager
func detectPackageManager() string {
	managers := []string{"apt", "yum", "dnf", "pacman"}
	paths := []string{"/usr/bin/", "/bin/"}

	for _, mgr := range managers {
		for _, path := range paths {
			if _, err := os.Stat(path + mgr); err == nil {
				return mgr
			}
		}
	}
	return ""
}

// getOSVersion returns the OS version
func getOSVersion() string {
	if _, err := os.Stat("/etc/os-release"); err == nil {
		data, err := os.ReadFile("/etc/os-release")
		if err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "VERSION_ID=") {
					return strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
				}
			}
		}
	}
	return "unknown"
}

// UpdatePackages updates the package manager's package list
func (p *UnixPlatform) UpdatePackages() error {
	var cmd *exec.Cmd
	switch p.packageManager {
	case "apt":
		cmd = exec.Command("sudo", "apt", "update")
	case "yum":
		cmd = exec.Command("sudo", "yum", "check-update")
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "check-update")
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-Sy")
	default:
		return fmt.Errorf("unsupported package manager: %s", p.packageManager)
	}
	return cmd.Run()
}

// UpgradePackages upgrades all installed packages
func (p *UnixPlatform) UpgradePackages() error {
	var cmd *exec.Cmd
	switch p.packageManager {
	case "apt":
		cmd = exec.Command("sudo", "apt", "upgrade", "-y")
	case "yum":
		cmd = exec.Command("sudo", "yum", "upgrade", "-y")
	case "dnf":
		cmd = exec.Command("sudo", "dnf", "upgrade", "-y")
	case "pacman":
		cmd = exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
	default:
		return fmt.Errorf("unsupported package manager: %s", p.packageManager)
	}
	return cmd.Run()
} 