package packages

import (
	"testing"
)

// mockPackageManager implements PackageManager interface for testing
type mockPackageManager struct {
	installedPackages map[string]string // package name -> version
}

func newMockPackageManager() *mockPackageManager {
	return &mockPackageManager{
		installedPackages: make(map[string]string),
	}
}

// Name returns the name of the package manager
func (m *mockPackageManager) Name() string {
	return "mock"
}

// IsAvailable checks if the package manager is available on the system
func (m *mockPackageManager) IsAvailable() bool {
	return true
}

// Install installs the given packages
func (m *mockPackageManager) Install(packages ...string) error {
	for _, packageName := range packages {
		m.installedPackages[packageName] = "1.0.0"
	}
	return nil
}

// Update updates the package list
func (m *mockPackageManager) Update() error {
	return nil
}

// IsInstalled checks if a package is installed
func (m *mockPackageManager) IsInstalled(pkg string) bool {
	_, exists := m.installedPackages[pkg]
	return exists
}

// Remove removes a package
func (m *mockPackageManager) Remove(pkg string) error {
	delete(m.installedPackages, pkg)
	return nil
}

// GetVersion returns the version of a package
func (m *mockPackageManager) GetVersion(packageName string) (string, error) {
	if version, exists := m.installedPackages[packageName]; exists {
		return version, nil
	}
	return "", nil
}

// ListInstalled returns a list of installed packages
func (m *mockPackageManager) ListInstalled() ([]string, error) {
	packages := make([]string, 0, len(m.installedPackages))
	for pkg := range m.installedPackages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func TestMockPackageManager(t *testing.T) {
	pm := newMockPackageManager()

	// Test Install
	t.Run("Install", func(t *testing.T) {
		err := pm.Install("test-package")
		if err != nil {
			t.Errorf("Install() error = %v", err)
		}
		if !pm.IsInstalled("test-package") {
			t.Error("Package should be installed after Install()")
		}
	})

	// Test GetVersion
	t.Run("GetVersion", func(t *testing.T) {
		version, err := pm.GetVersion("test-package")
		if err != nil {
			t.Errorf("GetVersion() error = %v", err)
		}
		if version != "1.0.0" {
			t.Errorf("GetVersion() = %v, want %v", version, "1.0.0")
		}
	})

	// Test ListInstalled
	t.Run("ListInstalled", func(t *testing.T) {
		packages, err := pm.ListInstalled()
		if err != nil {
			t.Errorf("ListInstalled() error = %v", err)
		}
		if len(packages) != 1 {
			t.Errorf("ListInstalled() returned %d packages, want 1", len(packages))
		}
		if packages[0] != "test-package" {
			t.Errorf("ListInstalled() returned %v, want [test-package]", packages)
		}
	})

	// Test Remove
	t.Run("Remove", func(t *testing.T) {
		err := pm.Remove("test-package")
		if err != nil {
			t.Errorf("Remove() error = %v", err)
		}
		if pm.IsInstalled("test-package") {
			t.Error("Package should not be installed after Remove()")
		}
	})
} 