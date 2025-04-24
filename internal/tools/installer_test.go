package tools

import (
	"bytes"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// mockPackagesManager implements interfaces.PackageManager for testing
type mockPackagesManager struct {
	installedPackages map[string]bool
}

func newMockPackagesManager() *mockPackagesManager {
	return &mockPackagesManager{
		installedPackages: make(map[string]bool),
	}
}

func (m *mockPackagesManager) GetName() string {
	return "mock"
}

func (m *mockPackagesManager) IsAvailable() bool {
	return true
}

func (m *mockPackagesManager) Install(pkg string) error {
	m.installedPackages[pkg] = true
	return nil
}

func (m *mockPackagesManager) Update() error {
	return nil
}

func (m *mockPackagesManager) Upgrade() error {
	return nil
}

func (m *mockPackagesManager) IsInstalled(pkg string) bool {
	installed := m.installedPackages[pkg]
	return installed
}

func (m *mockPackagesManager) Remove(pkg string) error {
	delete(m.installedPackages, pkg)
	return nil
}

func (m *mockPackagesManager) GetVersion(packageName string) (string, error) {
	if m.installedPackages[packageName] {
		return "1.0.0", nil
	}
	return "", nil
}

func (m *mockPackagesManager) ListInstalled() ([]string, error) {
	packages := make([]string, 0, len(m.installedPackages))
	for pkg := range m.installedPackages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (m *mockPackagesManager) SetupSpecialPackage(pkg string) error {
	return nil
}

func (m *mockPackagesManager) IsPackageAvailable(pkg string) bool {
	// For testing purposes, assume all packages are available
	return true
}

// testInstallEssentialTools is a helper function for testing InstallEssentialTools
func testInstallEssentialTools(t *testing.T, pm interfaces.PackageManager, logger *log.Logger) error {
	// Create a logger that captures output
	var logOutput strings.Builder
	logger.SetOutput(&logOutput)

	// Test installing essential tools with verification skipped
	err := InstallEssentialTools(pm, logger, true)
	if err != nil {
		t.Errorf("InstallEssentialTools failed: %v", err)
		return err
	}

	// Verify that the tools were installed
	expectedTools := []string{"git", "curl", "wget"}
	for _, tool := range expectedTools {
		installed := pm.IsInstalled(tool)
		if !installed {
			t.Errorf("Expected %s to be installed", tool)
		}
	}

	// Verify log output
	logStr := logOutput.String()
	if !strings.Contains(logStr, "Installing essential development tools") {
		t.Error("Expected log to contain installation message")
	}

	return nil
}

func TestInstallEssentialTools(t *testing.T) {
	// Create a mock package manager
	mockPM := newMockPackagesManager()

	// Create a logger
	logger := log.New(log.DebugLevel)

	// Run the test
	testInstallEssentialTools(t, mockPM, logger)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
} 