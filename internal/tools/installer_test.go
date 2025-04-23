package tools

import (
	"bytes"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// mockPackagesManager implements packages.Manager for testing
type mockPackagesManager struct {
	installedPackages map[string]bool
}

func (m *mockPackagesManager) Install(packageName string) error {
	if m.installedPackages == nil {
		m.installedPackages = make(map[string]bool)
	}
	m.installedPackages[packageName] = true
	return nil
}

func (m *mockPackagesManager) Uninstall(packageName string) error {
	delete(m.installedPackages, packageName)
	return nil
}

func (m *mockPackagesManager) Update(packageName string) error {
	return nil
}

func (m *mockPackagesManager) IsInstalled(packageName string) (bool, error) { 
	return m.installedPackages[packageName], nil 
}

func (m *mockPackagesManager) ListInstalled() ([]string, error) {
	packages := make([]string, 0, len(m.installedPackages))
	for pkg := range m.installedPackages {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (m *mockPackagesManager) GetVersion(packageName string) (string, error) { 
	return "1.0.0", nil 
}

// testInstallEssentialTools is a helper function for testing InstallEssentialTools
func testInstallEssentialTools(t *testing.T, pm packages.Manager, logger *log.Logger) error {
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
		installed, err := pm.IsInstalled(tool)
		if err != nil {
			t.Errorf("Error checking if %s is installed: %v", tool, err)
			continue
		}
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
	mockPM := &mockPackagesManager{
		installedPackages: make(map[string]bool),
	}

	// Create a logger
	logger := log.New(log.DebugLevel)

	// Run the test
	testInstallEssentialTools(t, mockPM, logger)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
} 