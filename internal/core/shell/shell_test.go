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
	if len(plugins) != 1 || plugins[0].Path != plugin {
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

func TestShellPluginWithMetadata(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Test adding a plugin with metadata
	pluginPath := "/path/to/plugin.zsh"
	pluginName := "test-plugin"
	pluginVersion := "1.2.3"
	pluginDescription := "A test plugin"
	pluginDependencies := []string{"dependency1", "dependency2"}

	err = shell.AddPluginWithMetadata(pluginName, pluginPath, pluginVersion, pluginDescription, pluginDependencies)
	if err != nil {
		t.Errorf("Failed to add plugin with metadata: %v", err)
	}

	// Test getting the plugin
	plugin, err := shell.GetPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to get plugin: %v", err)
	}

	// Check plugin metadata
	if plugin.Name != pluginName {
		t.Errorf("Expected plugin name to be %s, got %s", pluginName, plugin.Name)
	}
	if plugin.Path != pluginPath {
		t.Errorf("Expected plugin path to be %s, got %s", pluginPath, plugin.Path)
	}
	if plugin.Version != pluginVersion {
		t.Errorf("Expected plugin version to be %s, got %s", pluginVersion, plugin.Version)
	}
	if plugin.Description != pluginDescription {
		t.Errorf("Expected plugin description to be %s, got %s", pluginDescription, plugin.Description)
	}
	if !plugin.Enabled {
		t.Error("Plugin should be enabled")
	}
	if len(plugin.Dependencies) != len(pluginDependencies) {
		t.Errorf("Expected %d dependencies, got %d", len(pluginDependencies), len(plugin.Dependencies))
	}
}

func TestShellPluginEnableDisable(t *testing.T) {
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

	// Add a plugin
	pluginPath := "/path/to/plugin.zsh"
	err = testShell.AddPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to add plugin: %v", err)
	}

	// Disable the plugin
	err = testShell.DisablePlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to disable plugin: %v", err)
	}

	// Check if the plugin is disabled
	plugin, err := testShell.GetPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to get plugin: %v", err)
	}
	if plugin.Enabled {
		t.Error("Plugin should be disabled")
	}

	// Enable the plugin
	err = testShell.EnablePlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to enable plugin: %v", err)
	}

	// Check if the plugin is enabled
	plugin, err = testShell.GetPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to get plugin: %v", err)
	}
	if !plugin.Enabled {
		t.Error("Plugin should be enabled")
	}

	// Try to enable an already enabled plugin (should fail)
	err = testShell.EnablePlugin(pluginPath)
	if err == nil {
		t.Error("Enabling an already enabled plugin should fail")
	}

	// Try to disable an already disabled plugin (should fail)
	err = testShell.DisablePlugin(pluginPath)
	if err == nil {
		t.Error("Disabling an already disabled plugin should fail")
	}
}

func TestShellPluginConfig(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Add a plugin
	pluginPath := "/path/to/plugin.zsh"
	err = shell.AddPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to add plugin: %v", err)
	}

	// Set a plugin configuration
	configKey := "theme"
	configValue := "dark"
	err = shell.SetPluginConfig(pluginPath, configKey, configValue)
	if err != nil {
		t.Errorf("Failed to set plugin config: %v", err)
	}

	// Get the plugin configuration
	value, err := shell.GetPluginConfig(pluginPath, configKey)
	if err != nil {
		t.Errorf("Failed to get plugin config: %v", err)
	}
	if value != configValue {
		t.Errorf("Expected plugin config value to be %s, got %s", configValue, value)
	}

	// Try to get a non-existent configuration (should fail)
	_, err = shell.GetPluginConfig(pluginPath, "non-existent-key")
	if err == nil {
		t.Error("Getting a non-existent configuration should fail")
	}
}

func TestShellPluginUpdate(t *testing.T) {
	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Add a plugin
	pluginPath := "/path/to/plugin.zsh"
	err = shell.AddPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to add plugin: %v", err)
	}

	// Update the plugin version
	newVersion := "2.0.0"
	err = shell.UpdatePlugin(pluginPath, newVersion)
	if err != nil {
		t.Errorf("Failed to update plugin: %v", err)
	}

	// Check if the plugin version was updated
	plugin, err := shell.GetPlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to get plugin: %v", err)
	}
	if plugin.Version != newVersion {
		t.Errorf("Expected plugin version to be %s, got %s", newVersion, plugin.Version)
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

	// Add a plugin
	plugin := "/path/to/plugin.zsh"
	err = testShell.AddPlugin(plugin)
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
		Plugins:    []*Plugin{},
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