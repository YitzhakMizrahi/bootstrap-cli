package packages

import (
	"runtime"
	"testing"
)

func TestDetectPackageManager(t *testing.T) {
	pm, err := DetectPackageManager()

	switch runtime.GOOS {
	case "linux":
		// On Linux, we should at least find one package manager
		if err == ErrPackageManagerNotFound {
			t.Error("No package manager found on Linux")
		}
		if pm != nil && pm.Name() == "" {
			t.Error("Package manager name is empty")
		}
	case "darwin":
		// On macOS, we might find Homebrew
		if pm != nil && pm.Name() != string(Homebrew) {
			t.Errorf("Expected Homebrew on macOS, got %s", pm.Name())
		}
	default:
		// On other systems, we expect no package manager
		if err != ErrPackageManagerNotFound {
			t.Errorf("Expected no package manager, got %v", pm)
		}
	}
}

func TestAPTManager(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skipping APT tests on non-Linux system")
	}

	apt := &APTManager{}
	if apt.Name() != string(APT) {
		t.Errorf("Expected APT name to be %s, got %s", APT, apt.Name())
	}
}

func TestHomebrewManager(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping Homebrew tests on non-macOS system")
	}

	brew := &HomebrewManager{}
	if brew.Name() != string(Homebrew) {
		t.Errorf("Expected Homebrew name to be %s, got %s", Homebrew, brew.Name())
	}
} 