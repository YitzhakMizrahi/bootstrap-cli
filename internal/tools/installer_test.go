package tools

import (
	"bytes"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// MockPackageManager for testing
type MockPackageManager struct {
	installed map[string]bool
}

func (m *MockPackageManager) Name() string                 { return "mock" }
func (m *MockPackageManager) IsAvailable() bool           { return true }
func (m *MockPackageManager) Update() error               { return nil }
func (m *MockPackageManager) Remove(pkg string) error     { return nil }
func (m *MockPackageManager) IsInstalled(pkg string) bool { return m.installed[pkg] }

func (m *MockPackageManager) Install(packages ...string) error {
	if m.installed == nil {
		m.installed = make(map[string]bool)
	}
	for _, pkg := range packages {
		m.installed[pkg] = true
	}
	return nil
}

func TestInstallCoreTools(t *testing.T) {
	// Create a buffer to capture log output
	var logBuf bytes.Buffer

	// Create a logger that writes to our buffer
	logger := log.New(log.DebugLevel)
	logger.SetOutput(&logBuf)

	// Create mock package manager
	mockPM := &MockPackageManager{}

	// Create test tools
	testTools := []*install.Tool{
		{
			Name:        "TestTool",
			PackageName: "test-tool",
			PackageNames: &install.PackageMapping{
				Default: "test-tool",
				APT:     "test-tool",
				DNF:     "test-tool",
				Pacman:  "test-tool",
				Brew:    "test-tool",
			},
			VerifyCommand: "echo 'test'",
		},
	}

	// Create install options
	opts := &InstallOptions{
		Logger:         logger,
		PackageManager: mockPM,
		Tools:         testTools,
	}

	// Test installation
	if err := InstallCoreTools(opts); err != nil {
		t.Errorf("InstallCoreTools() error = %v", err)
	}

	// Verify log output contains success messages
	logOutput := logBuf.String()
	if !contains(logOutput, "Installing core development tools") {
		t.Error("Log missing installation start message")
	}
	if !contains(logOutput, "TestTool installed successfully") {
		t.Error("Log missing tool installation success message")
	}
	if !contains(logOutput, "All core tools installed successfully") {
		t.Error("Log missing final success message")
	}

	// Verify tool was "installed"
	if !mockPM.IsInstalled("test-tool") {
		t.Error("Tool was not marked as installed")
	}

	// Test verification
	logBuf.Reset()
	if err := VerifyCoreTools(opts); err != nil {
		t.Errorf("VerifyCoreTools() error = %v", err)
	}

	// Verify verification log output
	logOutput = logBuf.String()
	if !contains(logOutput, "Verifying core tools installation") {
		t.Error("Log missing verification start message")
	}
	if !contains(logOutput, "TestTool verified") {
		t.Error("Log missing tool verification success message")
	}
	if !contains(logOutput, "All core tools verified successfully") {
		t.Error("Log missing final verification success message")
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
} 