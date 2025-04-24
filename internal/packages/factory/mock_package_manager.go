package factory

import "github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"

// MockPackageManager implements interfaces.PackageManager for testing
type MockPackageManager struct {
	name string
}

// NewMockPackageManager creates a new mock package manager
func NewMockPackageManager(maxFailures int, name string) interfaces.PackageManager {
	return &MockPackageManager{
		name: name,
	}
}

func (m *MockPackageManager) GetName() string                     { return m.name }
func (m *MockPackageManager) IsAvailable() bool                   { return true }
func (m *MockPackageManager) Install(pkg string) error           { return nil }
func (m *MockPackageManager) Update() error                      { return nil }
func (m *MockPackageManager) Upgrade() error                     { return nil }
func (m *MockPackageManager) IsInstalled(pkg string) bool        { return true }
func (m *MockPackageManager) Remove(pkg string) error            { return nil }
func (m *MockPackageManager) GetVersion(pkg string) (string, error) { return "", nil }
func (m *MockPackageManager) ListInstalled() ([]string, error)   { return nil, nil }
func (m *MockPackageManager) SetupSpecialPackage(pkg string) error { return nil } 