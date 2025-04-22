package install

import (
	"testing"
)

// MockPackageManager implements PackageManager for testing
type MockPackageManager struct {
	installedPackages map[string]bool
}

func (m *MockPackageManager) Name() string {
	return "mock"
}

func (m *MockPackageManager) IsAvailable() bool {
	return true
}

func (m *MockPackageManager) Install(packages ...string) error {
	for _, pkg := range packages {
		m.installedPackages[pkg] = true
	}
	return nil
}

func (m *MockPackageManager) Update() error {
	return nil
}

func (m *MockPackageManager) IsInstalled(pkg string) bool {
	return m.installedPackages[pkg]
}

func TestInstaller(t *testing.T) {
	// Create a mock package manager
	mockPM := &MockPackageManager{
		installedPackages: make(map[string]bool),
	}

	// Create an installer
	installer := NewInstaller(mockPM)

	// Create a test tool
	tool := &Tool{
		Name:         "test-tool",
		PackageName:  "test-package",
		Version:      "1.0.0",
		Dependencies: []string{"dep1", "dep2"},
		PostInstall:  []string{"echo 'test'"},
		// Don't set VerifyCommand for testing
	}

	// Install the tool
	err := installer.Install(tool)
	if err != nil {
		t.Errorf("Install() error = %v", err)
	}

	// Check if dependencies were installed
	for _, dep := range tool.Dependencies {
		if !mockPM.IsInstalled(dep) {
			t.Errorf("Dependency %s was not installed", dep)
		}
	}

	// Check if the tool was installed
	if !mockPM.IsInstalled(tool.PackageName) {
		t.Errorf("Tool %s was not installed", tool.PackageName)
	}
} 