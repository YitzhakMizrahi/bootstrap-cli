package shell

import (
	"os"
	"path/filepath"
	"testing"
)

// testShell is a shell type specifically for testing that allows method overrides
type testShell struct {
	*Shell
	getConfigPathFunc func() (string, error)
}

func (s *testShell) GetConfigPath() (string, error) {
	if s.getConfigPathFunc != nil {
		return s.getConfigPathFunc()
	}
	return s.Shell.GetConfigPath()
}

func TestShellPluginManagement(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Test adding a plugin
	plugin := "/path/to/plugin.zsh"
	err = shell.AddPlugin(plugin)
	if err != nil {
		t.Errorf("Failed to add plugin: %v", err)
	}

	// Test checking if plugin is added
	if !shell.HasPlugin(plugin) {
		t.Error("Plugin should be added")
	}

	// Test listing plugins
	plugins := shell.ListPlugins()
	if len(plugins) != 1 || plugins[0] != plugin {
		t.Errorf("Expected plugin list to contain %s, got %v", plugin, plugins)
	}

	// Test adding the same plugin again (should fail)
	err = shell.AddPlugin(plugin)
	if err == nil {
		t.Error("Adding the same plugin again should fail")
	}

	// Test removing a plugin
	err = shell.RemovePlugin(plugin)
	if err != nil {
		t.Errorf("Failed to remove plugin: %v", err)
	}

	// Test checking if plugin is removed
	if shell.HasPlugin(plugin) {
		t.Error("Plugin should be removed")
	}

	// Test removing a non-existent plugin
	err = shell.RemovePlugin("non-existent-plugin")
	if err == nil {
		t.Error("Removing a non-existent plugin should fail")
	}
}

func TestShellPluginConfiguration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Create a temporary config file
	configPath := filepath.Join(tempDir, shell.ConfigFile)
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Add a plugin
	plugin := "/path/to/plugin.zsh"
	err = shell.AddPlugin(plugin)
	if err != nil {
		t.Errorf("Failed to add plugin: %v", err)
	}

	// Check if the plugin was added to the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expectedContent := "source " + plugin + "\n"
	if string(content) != expectedContent {
		t.Errorf("Expected config content to be %q, got %q", expectedContent, string(content))
	}
}

func TestShellPluginUnsupportedShell(t *testing.T) {
	// Create a test shell with an unsupported name
	shell := &Shell{
		Name:       "unsupported",
		ConfigFile: ".unsupportedrc",
		RCFile:     ".unsupported_profile",
		Plugins:    []string{},
	}

	// Test adding a plugin to an unsupported shell
	plugin := "/path/to/plugin"
	err := shell.AddPlugin(plugin)
	if err == nil {
		t.Error("Adding a plugin to an unsupported shell should fail")
	}
}

func TestNewShell(t *testing.T) {
	// Test creating a bash shell
	bash, err := New("bash")
	if err != nil {
		t.Errorf("Failed to create bash shell: %v", err)
	}
	if bash.Name != "bash" {
		t.Errorf("Expected shell name to be 'bash', got '%s'", bash.Name)
	}
	if bash.ConfigFile != ".bashrc" {
		t.Errorf("Expected config file to be '.bashrc', got '%s'", bash.ConfigFile)
	}
	if bash.RCFile != ".bash_profile" {
		t.Errorf("Expected RC file to be '.bash_profile', got '%s'", bash.RCFile)
	}

	// Test creating a zsh shell
	zsh, err := New("zsh")
	if err != nil {
		t.Errorf("Failed to create zsh shell: %v", err)
	}
	if zsh.Name != "zsh" {
		t.Errorf("Expected shell name to be 'zsh', got '%s'", zsh.Name)
	}
	if zsh.ConfigFile != ".zshrc" {
		t.Errorf("Expected config file to be '.zshrc', got '%s'", zsh.ConfigFile)
	}
	if zsh.RCFile != ".zshenv" {
		t.Errorf("Expected RC file to be '.zshenv', got '%s'", zsh.RCFile)
	}

	// Test creating a fish shell
	fish, err := New("fish")
	if err != nil {
		t.Errorf("Failed to create fish shell: %v", err)
	}
	if fish.Name != "fish" {
		t.Errorf("Expected shell name to be 'fish', got '%s'", fish.Name)
	}
	if fish.ConfigFile != "config.fish" {
		t.Errorf("Expected config file to be 'config.fish', got '%s'", fish.ConfigFile)
	}
	if fish.RCFile != "config.fish" {
		t.Errorf("Expected RC file to be 'config.fish', got '%s'", fish.RCFile)
	}

	// Test creating an unsupported shell
	_, err = New("unsupported")
	if err == nil {
		t.Error("Creating an unsupported shell should fail")
	}
}

func TestGetHomeDir(t *testing.T) {
	homeDir, err := GetHomeDir()
	if err != nil {
		t.Errorf("Failed to get home directory: %v", err)
	}
	if homeDir == "" {
		t.Error("Home directory should not be empty")
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Get the config path
	configPath, err := shell.GetConfigPath()
	if err != nil {
		t.Errorf("Failed to get config path: %v", err)
	}

	// Get the home directory
	homeDir, err := GetHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	// Check if the config path is correct
	expectedPath := filepath.Join(homeDir, shell.ConfigFile)
	if configPath != expectedPath {
		t.Errorf("Expected config path to be '%s', got '%s'", expectedPath, configPath)
	}
}

func TestGetRCPath(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Get the RC path
	rcPath, err := shell.GetRCPath()
	if err != nil {
		t.Errorf("Failed to get RC path: %v", err)
	}

	// Get the home directory
	homeDir, err := GetHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	// Check if the RC path is correct
	expectedPath := filepath.Join(homeDir, shell.RCFile)
	if rcPath != expectedPath {
		t.Errorf("Expected RC path to be '%s', got '%s'", expectedPath, rcPath)
	}
}

func TestCreateConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell
	baseShell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Create a temporary config file path
	configPath := filepath.Join(tempDir, baseShell.ConfigFile)
	
	// Create a test shell with a custom GetConfigPath method
	testShell := &testShell{
		Shell: baseShell,
		getConfigPathFunc: func() (string, error) {
			return configPath, nil
		},
	}

	// Create the config file
	content := "# Test config file\n"
	err = testShell.CreateConfig(content)
	if err != nil {
		t.Errorf("Failed to create config file: %v", err)
	}

	// Check if the file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check if the content is correct
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if string(fileContent) != content {
		t.Errorf("Expected config content to be '%s', got '%s'", content, string(fileContent))
	}

	// Try to create the config file again (should fail)
	err = testShell.CreateConfig(content)
	if err == nil {
		t.Error("Creating the config file again should fail")
	}
}

func TestAppendToConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell
	baseShell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Create a temporary config file path
	configPath := filepath.Join(tempDir, baseShell.ConfigFile)
	
	// Create a test shell with a custom GetConfigPath method
	testShell := &testShell{
		Shell: baseShell,
		getConfigPathFunc: func() (string, error) {
			return configPath, nil
		},
	}

	// Append to the config file
	content1 := "# Test config file\n"
	err = testShell.AppendToConfig(content1)
	if err != nil {
		t.Errorf("Failed to append to config file: %v", err)
	}

	// Check if the file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Check if the content is correct
	fileContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if string(fileContent) != content1 {
		t.Errorf("Expected config content to be '%s', got '%s'", content1, string(fileContent))
	}

	// Append more content to the config file
	content2 := "# Additional content\n"
	err = testShell.AppendToConfig(content2)
	if err != nil {
		t.Errorf("Failed to append to config file: %v", err)
	}

	// Check if the content is correct
	fileContent, err = os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	expectedContent := content1 + content2
	if string(fileContent) != expectedContent {
		t.Errorf("Expected config content to be '%s', got '%s'", expectedContent, string(fileContent))
	}
}

func TestSetDefaultShell(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Set the default shell
	err = shell.SetDefaultShell()
	if err != nil {
		t.Errorf("Failed to set default shell: %v", err)
	}

	// Note: This is a simplified test since the actual implementation
	// just prints a message and doesn't actually change the default shell
} 