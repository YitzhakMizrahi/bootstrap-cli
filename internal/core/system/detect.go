package system

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// SystemInfo contains information about the current system
type SystemInfo struct {
	OS          string
	Arch        string
	Kernel      string
	Hostname    string
	PackageManager string
}

// GetSystemInfo returns information about the current system
func GetSystemInfo() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: getHostname(),
	}

	// Get kernel information
	kernel, err := getKernelInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get kernel info: %w", err)
	}
	info.Kernel = kernel

	// Determine package manager based on OS
	info.PackageManager = determinePackageManager(info.OS)

	return info, nil
}

// getHostname returns the system hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getKernelInfo returns the kernel version
func getKernelInfo() (string, error) {
	// This is a simplified implementation
	// In a real implementation, we would use platform-specific commands
	// to get the actual kernel version
	switch runtime.GOOS {
	case "linux":
		return "Linux", nil
	case "darwin":
		return "Darwin", nil
	case "windows":
		return "Windows", nil
	default:
		return runtime.GOOS, nil
	}
}

// determinePackageManager returns the appropriate package manager for the OS
func determinePackageManager(os string) string {
	switch strings.ToLower(os) {
	case "linux":
		// Check for specific Linux distributions
		// This is a simplified implementation
		return "apt" // Default to apt for now
	case "darwin":
		return "brew"
	case "windows":
		return "choco"
	default:
		return "unknown"
	}
} 