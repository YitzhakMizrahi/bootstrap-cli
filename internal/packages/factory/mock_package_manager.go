// Package factory provides factory methods for creating package managers
// and related mock implementations for testing purposes.
package factory

import (
	"github.com/stretchr/testify/mock"
)

// MockPackageManager is a mock implementation of interfaces.PackageManager
type MockPackageManager struct {
	mock.Mock
}

// NewMockPackageManager creates a new mock package manager
func NewMockPackageManager() *MockPackageManager {
	return &MockPackageManager{}
}

// Install is a mock method
func (m *MockPackageManager) Install(packageName string) error {
	args := m.Called(packageName)
	return args.Error(0)
}

// IsInstalled is a mock method - Updated Signature
func (m *MockPackageManager) IsInstalled(packageName string) (bool, error) {
	args := m.Called(packageName)
	return args.Bool(0), args.Error(1) // Return bool and error
}

// GetName is a mock method
func (m *MockPackageManager) GetName() string {
	args := m.Called()
	return args.String(0)
}

// IsAvailable is a mock method
func (m *MockPackageManager) IsAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

// IsPackageAvailable is a mock method - Added
func (m *MockPackageManager) IsPackageAvailable(packageName string) bool {
	args := m.Called(packageName)
	return args.Bool(0)
}

// Update is a mock method
func (m *MockPackageManager) Update() error {
	args := m.Called()
	return args.Error(0)
}

// Upgrade is a mock method
func (m *MockPackageManager) Upgrade() error {
	args := m.Called()
	return args.Error(0)
}

// Uninstall is a mock method - Renamed from Remove
func (m *MockPackageManager) Uninstall(packageName string) error {
	args := m.Called(packageName)
	return args.Error(0)
}

// GetVersion is a mock method
func (m *MockPackageManager) GetVersion(packageName string) (string, error) {
	args := m.Called(packageName)
	return args.String(0), args.Error(1)
}

// ListInstalled is a mock method
func (m *MockPackageManager) ListInstalled() ([]string, error) {
	args := m.Called()
	// Need to handle potential nil slice if error is expected
	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}
	return result, args.Error(1)
}

// SetupSpecialPackage is a mock method
func (m *MockPackageManager) SetupSpecialPackage(packageName string) error {
	args := m.Called(packageName)
	return args.Error(0)
} 