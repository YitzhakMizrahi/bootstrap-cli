package packages

import (
	"testing"
)

// mockPackageManager implements Manager interface for testing
type mockPackageManager struct {
	installedPackages map[string]string // package name -> version
}

func newMockPackageManager() *mockPackageManager {
	return &mockPackageManager{
		installedPackages: make(map[string]string),
	}
}

func (m *mockPackageManager) Install(packageName string) error {
	m.installedPackages[packageName] = "1.0.0"
	return nil
}

func (m *mockPackageManager) Uninstall(packageName string) error {
	delete(m.installedPackages, packageName)
	return nil
}

func (m *mockPackageManager) Update(packageName string) error {
	if version, exists := m.installedPackages[packageName]; exists {
		m.installedPackages[packageName] = version + ".1"
		return nil
	}
	return nil
}

func (m *mockPackageManager) IsInstalled(packageName string) (bool, error) {
	_, exists := m.installedPackages[packageName]
	return exists, nil
}

func (m *mockPackageManager) ListInstalled() ([]string, error) {
	packages := make([]string, 0, len(m.installedPackages))
	for pkg := range m.installedPackages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (m *mockPackageManager) GetVersion(packageName string) (string, error) {
	if version, exists := m.installedPackages[packageName]; exists {
		return version, nil
	}
	return "", nil
}

func TestMockPackageManager(t *testing.T) {
	pm := newMockPackageManager()

	// Test Install
	t.Run("Install", func(t *testing.T) {
		err := pm.Install("test-package")
		if err != nil {
			t.Errorf("Install() error = %v", err)
		}
		installed, _ := pm.IsInstalled("test-package")
		if !installed {
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

	// Test Update
	t.Run("Update", func(t *testing.T) {
		err := pm.Update("test-package")
		if err != nil {
			t.Errorf("Update() error = %v", err)
		}
		version, _ := pm.GetVersion("test-package")
		if version != "1.0.0.1" {
			t.Errorf("GetVersion() after update = %v, want %v", version, "1.0.0.1")
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

	// Test Uninstall
	t.Run("Uninstall", func(t *testing.T) {
		err := pm.Uninstall("test-package")
		if err != nil {
			t.Errorf("Uninstall() error = %v", err)
		}
		installed, _ := pm.IsInstalled("test-package")
		if installed {
			t.Error("Package should not be installed after Uninstall()")
		}
	})
} 