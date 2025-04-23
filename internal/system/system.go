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
	OS          string
	Arch        string
	Distro      string
	Version     string
	Kernel      string
	PackageType PackageManagerType // apt, dnf, pacman, brew
}

// Detect returns information about the current system
func Detect() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	var err error

	// Get kernel version
	if info.Kernel, err = getKernelVersion(); err != nil {
		return nil, fmt.Errorf("failed to get kernel version: %w", err)
	}

	// Get distro and version info for Linux
	if info.OS == "linux" {
		if err = getLinuxDistroInfo(info); err != nil {
			return nil, fmt.Errorf("failed to get Linux distro info: %w", err)
		}
	} else if info.OS == "darwin" {
		if err = getDarwinInfo(info); err != nil {
			return nil, fmt.Errorf("failed to get macOS info: %w", err)
		}
	}

	// Detect package manager type
	if err = detectPackageManager(info); err != nil {
		return nil, fmt.Errorf("failed to detect package manager: %w", err)
	}

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

// detectPackageManager determines the system's package manager
func detectPackageManager(info *SystemInfo) error {
	if info.OS == "linux" {
		// Check for apt (Debian/Ubuntu)
		if _, err := exec.LookPath("apt"); err == nil {
			info.PackageType = PackageManagerApt
			return nil
		}

		// Check for dnf (Fedora)
		if _, err := exec.LookPath("dnf"); err == nil {
			info.PackageType = PackageManagerDnf
			return nil
		}

		// Check for pacman (Arch)
		if _, err := exec.LookPath("pacman"); err == nil {
			info.PackageType = PackageManagerPacman
			return nil
		}

		return fmt.Errorf("no supported package manager found")
	} else if info.OS == "darwin" {
		// Check for Homebrew
		if _, err := exec.LookPath("brew"); err == nil {
			info.PackageType = PackageManagerHomebrew
			return nil
		}
		return fmt.Errorf("Homebrew not found")
	}

	return fmt.Errorf("unsupported OS: %s", info.OS)
} 