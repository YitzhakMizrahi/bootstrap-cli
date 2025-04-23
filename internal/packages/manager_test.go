package packages

import (
	"runtime"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
)

func TestPackageManagerDetection(t *testing.T) {
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()

	switch runtime.GOOS {
	case "linux":
		// On Linux, we should at least find one package manager
		if err != nil {
			t.Errorf("Failed to get package manager on Linux: %v", err)
		}
		if pm != nil && pm.GetName() == "" {
			t.Error("Package manager name is empty")
		}
	case "darwin":
		// On macOS, we might find Homebrew
		if pm != nil && pm.GetName() != string(interfaces.Homebrew) {
			t.Errorf("Expected Homebrew on macOS, got %s", pm.GetName())
		}
	default:
		// On other systems, we expect no package manager
		if pm != nil {
			t.Errorf("Expected no package manager, got %v", pm)
		}
	}
}

func TestPackageManagerFactory(t *testing.T) {
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	
	if err != nil {
		t.Errorf("Failed to get package manager: %v", err)
	}
	
	if pm == nil {
		t.Fatal("Package manager is nil")
	}
	
	if pm.GetName() == "" {
		t.Error("Package manager name is empty")
	}
} 