package system

import (
	"runtime"
)

// SystemInfo contains information about the current system
type SystemInfo struct {
	OS      string
	Arch    string
	Distro  string
	Version string
}

// Detect returns information about the current system
func Detect() (*SystemInfo, error) {
	info := &SystemInfo{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	// TODO: Implement distro detection for Linux
	// TODO: Implement version detection

	return info, nil
} 