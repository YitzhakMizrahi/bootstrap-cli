// Package factory provides factory methods for creating package managers
// and related mock implementations for testing purposes.
package factory

import "github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"

// MockPackageManager implements interfaces.PackageManager for testing
type MockPackageManager struct {
	name string
}

// NewMockPackageManager creates a new mock package manager
func NewMockPackageManager(_ int, name string) interfaces.PackageManager {
	return &MockPackageManager{
		name: name,
	}
}

// GetName returns the name of the package manager
func (m *MockPackageManager) GetName() string { return m.name }

// IsAvailable checks if the package manager is available on the system
func (m *MockPackageManager) IsAvailable() bool { return true }

// Install simulates installing a package
func (m *MockPackageManager) Install(_ string) error { return nil }

// Update simulates updating the package list
func (m *MockPackageManager) Update() error { return nil }

// Upgrade simulates upgrading all packages
func (m *MockPackageManager) Upgrade() error { return nil }

// IsInstalled checks if a package is installed
func (m *MockPackageManager) IsInstalled(_ string) bool { return true }

// Remove simulates removing a package
func (m *MockPackageManager) Remove(_ string) error { return nil }

// GetVersion returns the version of an installed package
func (m *MockPackageManager) GetVersion(_ string) (string, error) { return "", nil }

// ListInstalled returns a list of installed packages
func (m *MockPackageManager) ListInstalled() ([]string, error) { return nil, nil }

// SetupSpecialPackage simulates setting up a special package
func (m *MockPackageManager) SetupSpecialPackage(_ string) error { return nil } 