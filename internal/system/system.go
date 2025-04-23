package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// SystemInfo contains information about the current system
type SystemInfo struct {
	OS              string
	Arch            string
	Shell           string
	HomeDir         string
	Distro          string  // Linux distribution name or "macOS"
	Version         string  // OS/distro version
	Kernel          string  // Kernel version
	PackageType     string  // Package manager type (apt, dnf, pacman, brew)
	IsRoot          bool
	IsWSL           bool
	IsDocker        bool
	IsVM            bool
	IsContainer     bool
	IsHeadless      bool
	IsGUI           bool
	IsSSH           bool
	IsCI            bool
	IsDevelopment   bool
	IsProduction    bool
	IsStaging       bool
	IsTesting       bool
	IsDebug         bool
	IsVerbose       bool
	IsQuiet         bool
	IsSilent        bool
	IsInteractive   bool
	IsNonInteractive bool
	IsColor         bool
	IsNoColor       bool
	IsForce         bool
	IsDryRun        bool
}

// Detect gathers information about the current system
func Detect() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:     runtime.GOOS,
		Arch:   runtime.GOARCH,
		Shell:  os.Getenv("SHELL"),
		HomeDir: os.Getenv("HOME"),
		IsRoot: os.Geteuid() == 0,
	}

	// Get kernel version
	kernelVersion, err := getKernelVersion()
	if err == nil {
		info.Kernel = kernelVersion
	}

	// Detect OS-specific information
	switch info.OS {
	case "linux":
		if err := getLinuxDistroInfo(info); err != nil {
			return nil, fmt.Errorf("failed to get Linux distribution info: %w", err)
		}
		// Detect package manager type
		if _, err := exec.LookPath("apt"); err == nil {
			info.PackageType = "apt"
		} else if _, err := exec.LookPath("dnf"); err == nil {
			info.PackageType = "dnf"
		} else if _, err := exec.LookPath("pacman"); err == nil {
			info.PackageType = "pacman"
		}
	case "darwin":
		if err := getDarwinInfo(info); err != nil {
			return nil, fmt.Errorf("failed to get macOS info: %w", err)
		}
		// Check for Homebrew
		if _, err := exec.LookPath("brew"); err == nil {
			info.PackageType = "brew"
		}
	}

	// Detect WSL
	if _, err := os.ReadFile("/proc/version"); err == nil {
		content, err := os.ReadFile("/proc/version")
		if err == nil {
			info.IsWSL = strings.Contains(strings.ToLower(string(content)), "microsoft")
		}
	}

	// Detect Docker
	if _, err := os.ReadFile("/.dockerenv"); err == nil {
		info.IsDocker = true
	}

	// Detect VM
	if _, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		content, err := os.ReadFile("/sys/class/dmi/id/product_name")
		if err == nil {
			product := strings.ToLower(string(content))
			info.IsVM = strings.Contains(product, "virtualbox") ||
				strings.Contains(product, "vmware") ||
				strings.Contains(product, "qemu") ||
				strings.Contains(product, "xen")
		}
	}

	// Detect container
	info.IsContainer = info.IsDocker || info.IsWSL

	// Detect headless
	info.IsHeadless = os.Getenv("DISPLAY") == ""

	// Detect GUI
	info.IsGUI = !info.IsHeadless

	// Detect SSH
	info.IsSSH = os.Getenv("SSH_TTY") != ""

	// Detect CI
	info.IsCI = os.Getenv("CI") != ""

	// Detect environment
	info.IsDevelopment = os.Getenv("GO_ENV") == "development"
	info.IsProduction = os.Getenv("GO_ENV") == "production"
	info.IsStaging = os.Getenv("GO_ENV") == "staging"
	info.IsTesting = os.Getenv("GO_ENV") == "testing"

	// Detect debug mode
	info.IsDebug = os.Getenv("DEBUG") != ""
	info.IsVerbose = os.Getenv("VERBOSE") != ""
	info.IsQuiet = os.Getenv("QUIET") != ""
	info.IsSilent = os.Getenv("SILENT") != ""

	// Detect interactive mode
	info.IsInteractive = !info.IsCI && !info.IsNonInteractive
	info.IsNonInteractive = os.Getenv("NONINTERACTIVE") != ""

	// Detect color mode
	info.IsColor = !info.IsNoColor
	info.IsNoColor = os.Getenv("NO_COLOR") != ""

	// Detect force mode
	info.IsForce = os.Getenv("FORCE") != ""

	// Detect dry run mode
	info.IsDryRun = os.Getenv("DRY_RUN") != ""

	return info, nil
}

// getKernelVersion returns the kernel version string
func getKernelVersion() (string, error) {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/sys/kernel/osrelease")
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("uname", "-r")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
}

// getLinuxDistroInfo detects Linux distribution and version
func getLinuxDistroInfo(info *SystemInfo) error {
	// Try /etc/os-release first
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(data)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "ID=") {
				info.Distro = strings.Trim(line[3:], "\"")
			} else if strings.HasPrefix(line, "VERSION_ID=") {
				info.Version = strings.Trim(line[11:], "\"")
			}
		}
		if info.Distro != "" && info.Version != "" {
			return nil
		}
	}

	// Try lsb_release if available
	if path, err := exec.LookPath("lsb_release"); err == nil {
		cmd := exec.Command(path, "-a")
		out, err := cmd.Output()
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Distributor ID:") {
				info.Distro = strings.TrimSpace(line[15:])
			} else if strings.HasPrefix(line, "Release:") {
				info.Version = strings.TrimSpace(line[8:])
			}
		}
		return nil
	}

	return fmt.Errorf("could not determine Linux distribution")
}

// getDarwinInfo detects macOS version
func getDarwinInfo(info *SystemInfo) error {
	info.Distro = "macOS"
	cmd := exec.Command("sw_vers", "-productVersion")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	info.Version = strings.TrimSpace(string(out))
	return nil
} 