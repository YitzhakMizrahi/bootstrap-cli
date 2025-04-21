package pkgmanager

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// MockPackageManager is a mock implementation of the PackageManager interface
type MockPackageManager struct {
	name           string
	isAvailable    bool
	installedPkgs  map[string]bool
	searchResults  map[string][]string
	installError   error
	uninstallError error
	updateError    error
	updateAllError error
}

func (m *MockPackageManager) Name() string {
	return m.name
}

func (m *MockPackageManager) IsAvailable() bool {
	return m.isAvailable
}

func (m *MockPackageManager) Install(packageName string) error {
	if m.installError != nil {
		return m.installError
	}
	m.installedPkgs[packageName] = true
	return nil
}

func (m *MockPackageManager) Uninstall(packageName string) error {
	if m.uninstallError != nil {
		return m.uninstallError
	}
	delete(m.installedPkgs, packageName)
	return nil
}

func (m *MockPackageManager) Update(packageName string) error {
	if m.updateError != nil {
		return m.updateError
	}
	return nil
}

func (m *MockPackageManager) UpdateAll() error {
	if m.updateAllError != nil {
		return m.updateAllError
	}
	return nil
}

func (m *MockPackageManager) IsInstalled(packageName string) bool {
	return m.installedPkgs[packageName]
}

func (m *MockPackageManager) Search(query string) ([]string, error) {
	if results, ok := m.searchResults[query]; ok {
		return results, nil
	}
	return []string{}, nil
}

// TestNewInstaller tests the creation of a new installer
func TestNewInstaller(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Test with valid parameters
	installer := NewInstaller(mockPM, tempDir, "state.json")
	if installer == nil {
		t.Fatal("NewInstaller() returned nil")
	}
	if installer.packageManager != mockPM {
		t.Errorf("NewInstaller() packageManager = %v, want %v", installer.packageManager, mockPM)
	}
	if installer.installDir != tempDir {
		t.Errorf("NewInstaller() installDir = %v, want %v", installer.installDir, tempDir)
	}
	if installer.stateFile != "state.json" {
		t.Errorf("NewInstaller() stateFile = %v, want %v", installer.stateFile, "state.json")
	}
	if installer.tools == nil {
		t.Error("NewInstaller() tools is nil")
	}
}

// TestRegisterTool tests the RegisterTool method
func TestRegisterTool(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   false,
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Check if the tool was registered
	if len(installer.tools) != 1 {
		t.Errorf("RegisterTool() tools count = %v, want %v", len(installer.tools), 1)
	}
	if installer.tools["test-tool"] != tool {
		t.Errorf("RegisterTool() tool = %v, want %v", installer.tools["test-tool"], tool)
	}
}

// TestInstallTool tests the InstallTool method
func TestInstallTool(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   false,
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Install the tool
	err = installer.InstallTool("test-tool")
	if err != nil {
		t.Errorf("InstallTool() error = %v", err)
	}

	// Check if the tool was installed
	if !installer.tools["test-tool"].Installed {
		t.Error("InstallTool() tool.Installed = false, want true")
	}
	if installer.tools["test-tool"].InstallTime.IsZero() {
		t.Error("InstallTool() tool.InstallTime is zero")
	}
	if !mockPM.IsInstalled("test-package") {
		t.Error("InstallTool() package not installed in mock package manager")
	}
}

// TestUninstallTool tests the UninstallTool method
func TestUninstallTool(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   true,
		InstallTime: time.Now(),
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Mark the package as installed in the mock package manager
	mockPM.installedPkgs["test-package"] = true

	// Uninstall the tool
	err = installer.UninstallTool("test-tool")
	if err != nil {
		t.Errorf("UninstallTool() error = %v", err)
	}

	// Check if the tool was uninstalled
	if installer.tools["test-tool"].Installed {
		t.Error("UninstallTool() tool.Installed = true, want false")
	}
	if !installer.tools["test-tool"].InstallTime.IsZero() {
		t.Error("UninstallTool() tool.InstallTime is not zero")
	}
	if mockPM.IsInstalled("test-package") {
		t.Error("UninstallTool() package still installed in mock package manager")
	}
}

// TestIsToolInstalled tests the IsToolInstalled method
func TestIsToolInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   true,
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Check if the tool is installed
	if !installer.IsToolInstalled("test-tool") {
		t.Error("IsToolInstalled() = false, want true")
	}

	// Create another tool
	tool2 := &Tool{
		Name:        "test-tool-2",
		Description: "Another test tool",
		PackageName: "test-package-2",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool-2",
		Installed:   false,
	}

	// Register the tool
	installer.RegisterTool(tool2)

	// Check if the tool is not installed
	if installer.IsToolInstalled("test-tool-2") {
		t.Error("IsToolInstalled() = true, want false")
	}

	// Check if a non-existent tool is not installed
	if installer.IsToolInstalled("non-existent-tool") {
		t.Error("IsToolInstalled() for non-existent tool = true, want false")
	}
}

// TestGetTool tests the GetTool method
func TestGetTool(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   true,
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Get the tool
	gotTool, err := installer.GetTool("test-tool")
	if err != nil {
		t.Errorf("GetTool() error = %v", err)
	}
	if gotTool != tool {
		t.Errorf("GetTool() = %v, want %v", gotTool, tool)
	}

	// Get a non-existent tool
	_, err = installer.GetTool("non-existent-tool")
	if err == nil {
		t.Error("GetTool() for non-existent tool did not return an error")
	}
}

// TestListTools tests the ListTools method
func TestListTools(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create tools
	tool1 := &Tool{
		Name:        "test-tool-1",
		Description: "A test tool",
		PackageName: "test-package-1",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool-1",
		Installed:   true,
	}

	tool2 := &Tool{
		Name:        "test-tool-2",
		Description: "Another test tool",
		PackageName: "test-package-2",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool-2",
		Installed:   false,
	}

	// Register the tools
	installer.RegisterTool(tool1)
	installer.RegisterTool(tool2)

	// List the tools
	tools := installer.ListTools()
	if len(tools) != 2 {
		t.Errorf("ListTools() tools count = %v, want %v", len(tools), 2)
	}

	// Check if the tools are in the list
	found1 := false
	found2 := false
	for _, t := range tools {
		if t.Name == "test-tool-1" {
			found1 = true
		}
		if t.Name == "test-tool-2" {
			found2 = true
		}
	}
	if !found1 {
		t.Error("ListTools() did not include test-tool-1")
	}
	if !found2 {
		t.Error("ListTools() did not include test-tool-2")
	}
}

// TestListInstalledTools tests the ListInstalledTools method
func TestListInstalledTools(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create tools
	tool1 := &Tool{
		Name:        "test-tool-1",
		Description: "A test tool",
		PackageName: "test-package-1",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool-1",
		Installed:   true,
	}

	tool2 := &Tool{
		Name:        "test-tool-2",
		Description: "Another test tool",
		PackageName: "test-package-2",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool-2",
		Installed:   false,
	}

	// Register the tools
	installer.RegisterTool(tool1)
	installer.RegisterTool(tool2)

	// List the installed tools
	tools := installer.ListInstalledTools()
	if len(tools) != 1 {
		t.Errorf("ListInstalledTools() tools count = %v, want %v", len(tools), 1)
	}

	// Check if the installed tool is in the list
	if tools[0].Name != "test-tool-1" {
		t.Errorf("ListInstalledTools() tool name = %v, want %v", tools[0].Name, "test-tool-1")
	}
}

// TestSaveState tests the SaveState method
func TestSaveState(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   true,
		InstallTime: time.Now(),
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Save the state
	err = installer.SaveState()
	if err != nil {
		t.Errorf("SaveState() error = %v", err)
	}

	// Check if the state file was created
	stateFile := filepath.Join(tempDir, "state.json")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		t.Error("SaveState() did not create state file")
	}
}

// TestLoadState tests the LoadState method
func TestLoadState(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "installer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock package manager
	mockPM := &MockPackageManager{
		name:          "mock",
		isAvailable:   true,
		installedPkgs: make(map[string]bool),
		searchResults: make(map[string][]string),
	}

	// Create an installer
	installer := NewInstaller(mockPM, tempDir, "state.json")

	// Create a tool
	tool := &Tool{
		Name:        "test-tool",
		Description: "A test tool",
		PackageName: "test-package",
		Version:     "1.0.0",
		Category:    "test",
		URL:         "https://example.com/test-tool",
		Installed:   true,
		InstallTime: time.Now(),
	}

	// Register the tool
	installer.RegisterTool(tool)

	// Save the state
	err = installer.SaveState()
	if err != nil {
		t.Errorf("SaveState() error = %v", err)
	}

	// Create a new installer
	installer2 := NewInstaller(mockPM, tempDir, "state.json")

	// Load the state
	err = installer2.LoadState()
	if err != nil {
		t.Errorf("LoadState() error = %v", err)
	}

	// Check if the tool was loaded
	if len(installer2.tools) != 1 {
		t.Errorf("LoadState() tools count = %v, want %v", len(installer2.tools), 1)
	}
	if installer2.tools["test-tool"] == nil {
		t.Error("LoadState() did not load test-tool")
	}
	if installer2.tools["test-tool"].Name != "test-tool" {
		t.Errorf("LoadState() tool name = %v, want %v", installer2.tools["test-tool"].Name, "test-tool")
	}
	if !installer2.tools["test-tool"].Installed {
		t.Error("LoadState() tool.Installed = false, want true")
	}
} 