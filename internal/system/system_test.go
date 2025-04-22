package system

import (
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	info, err := Detect()
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Check OS
	if info.OS != runtime.GOOS {
		t.Errorf("Detect() OS = %v, want %v", info.OS, runtime.GOOS)
	}

	// Check Arch
	if info.Arch != runtime.GOARCH {
		t.Errorf("Detect() Arch = %v, want %v", info.Arch, runtime.GOARCH)
	}

	// Check Kernel version is not empty
	if info.Kernel == "" {
		t.Error("Detect() Kernel is empty")
	}

	// Check Distro and Version for Linux
	if info.OS == "linux" {
		if info.Distro == "" {
			t.Error("Detect() Distro is empty on Linux")
		}
		if info.Version == "" {
			t.Error("Detect() Version is empty on Linux")
		}
	}

	// Check Distro and Version for macOS
	if info.OS == "darwin" {
		if info.Distro != "macOS" {
			t.Errorf("Detect() Distro = %v, want macOS", info.Distro)
		}
		if info.Version == "" {
			t.Error("Detect() Version is empty on macOS")
		}
	}

	// Check PackageType is set appropriately
	if info.OS == "linux" || info.OS == "darwin" {
		if info.PackageType == "" {
			t.Error("Detect() PackageType is empty")
		}
	}
} 