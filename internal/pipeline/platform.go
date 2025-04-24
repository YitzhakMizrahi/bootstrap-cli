package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Platform represents the current system platform
type Platform struct {
	OS             string
	PackageManager string
	Shell          string
	Arch           string
}

// DetectPlatform detects the current platform and its characteristics
func DetectPlatform() (*Platform, error) {
	platform := &Platform{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// Detect package manager
	if err := platform.detectPackageManager(); err != nil {
		return nil, fmt.Errorf("failed to detect package manager: %w", err)
	}

	// Detect shell
	if err := platform.detectShell(); err != nil {
		return nil, fmt.Errorf("failed to detect shell: %w", err)
	}

	return platform, nil
}

// detectPackageManager detects the available package manager
func (p *Platform) detectPackageManager() error {
	// Check for apt (Debian/Ubuntu)
	if _, err := exec.LookPath("apt"); err == nil {
		p.PackageManager = "apt"
		return nil
	}

	// Check for brew (macOS)
	if _, err := exec.LookPath("brew"); err == nil {
		p.PackageManager = "brew"
		return nil
	}

	// Check for pacman (Arch Linux)
	if _, err := exec.LookPath("pacman"); err == nil {
		p.PackageManager = "pacman"
		return nil
	}

	// Check for dnf (Fedora)
	if _, err := exec.LookPath("dnf"); err == nil {
		p.PackageManager = "dnf"
		return nil
	}

	// Check for yum (RHEL/CentOS)
	if _, err := exec.LookPath("yum"); err == nil {
		p.PackageManager = "yum"
		return nil
	}

	return fmt.Errorf("no supported package manager found")
}

// detectShell detects the current shell
func (p *Platform) detectShell() error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return fmt.Errorf("SHELL environment variable not set")
	}

	shellName := strings.ToLower(shell)
	switch {
	case strings.Contains(shellName, "zsh"):
		p.Shell = "zsh"
	case strings.Contains(shellName, "bash"):
		p.Shell = "bash"
	case strings.Contains(shellName, "fish"):
		p.Shell = "fish"
	default:
		p.Shell = "unknown"
	}

	return nil
}

// String returns a string representation of the platform
func (p *Platform) String() string {
	return fmt.Sprintf("OS: %s, Arch: %s, Package Manager: %s, Shell: %s",
		p.OS, p.Arch, p.PackageManager, p.Shell)
}

// IsSupported checks if the current platform is supported
func (p *Platform) IsSupported() bool {
	// Check OS
	switch p.OS {
	case "linux", "darwin":
		// These OSes are supported
	default:
		return false
	}

	// Check package manager
	switch p.PackageManager {
	case "apt", "brew", "pacman", "dnf", "yum":
		// These package managers are supported
	default:
		return false
	}

	// Check shell
	switch p.Shell {
	case "bash", "zsh", "fish":
		// These shells are supported
	default:
		return false
	}

	return true
} 