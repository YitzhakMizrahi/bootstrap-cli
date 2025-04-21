package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestPluginDependencies(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "shell-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}
	shell.ConfigFile = filepath.Join(tmpDir, ".zshrc")

	// Add plugins with dependencies
	err = shell.AddPluginWithMetadata("plugin1", "/path/to/plugin1", "1.0.0", "Test plugin 1", nil)
	if err != nil {
		t.Fatalf("Failed to add plugin1: %v", err)
	}

	err = shell.AddPluginWithMetadata("plugin2", "/path/to/plugin2", "1.0.0", "Test plugin 2", []string{"plugin1"})
	if err != nil {
		t.Fatalf("Failed to add plugin2: %v", err)
	}

	// Try to enable plugin2 without enabling plugin1 (should fail)
	err = shell.EnablePlugin("plugin2")
	if err == nil {
		t.Error("Expected error when enabling plugin2 without plugin1 enabled")
	}

	// Enable plugin1 first
	if err := shell.EnablePlugin("plugin1"); err != nil {
		t.Fatalf("Failed to enable plugin1: %v", err)
	}

	// Now enable plugin2 (should succeed)
	if err := shell.EnablePlugin("plugin2"); err != nil {
		t.Errorf("Failed to enable plugin2: %v", err)
	}

	// Try to disable plugin1 while plugin2 is enabled (should fail)
	err = shell.DisablePlugin("plugin1")
	if err == nil {
		t.Error("Expected error when disabling plugin1 while plugin2 is enabled")
	}

	// Disable plugin2 first
	if err := shell.DisablePlugin("plugin2"); err != nil {
		t.Fatalf("Failed to disable plugin2: %v", err)
	}

	// Now disable plugin1 (should succeed)
	if err := shell.DisablePlugin("plugin1"); err != nil {
		t.Errorf("Failed to disable plugin1: %v", err)
	}
}

func TestPluginConfiguration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "shell-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}
	shell.ConfigFile = filepath.Join(tmpDir, ".zshrc")

	// Add a plugin
	err = shell.AddPluginWithMetadata("test-plugin", "/path/to/plugin", "1.0.0", "Test plugin", nil)
	if err != nil {
		t.Fatalf("Failed to add plugin: %v", err)
	}

	// Set plugin configuration
	if err := shell.SetPluginConfig("test-plugin", "option1", "value1"); err != nil {
		t.Fatalf("Failed to set plugin config: %v", err)
	}

	// Get plugin configuration
	value, err := shell.GetPluginConfig("test-plugin", "option1")
	if err != nil {
		t.Fatalf("Failed to get plugin config: %v", err)
	}
	if value != "value1" {
		t.Errorf("Expected config value 'value1', got '%s'", value)
	}

	// Try to get non-existent config
	_, err = shell.GetPluginConfig("test-plugin", "nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent config")
	}
}

func TestPluginVersionManagement(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "shell-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test shell
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}
	shell.ConfigFile = filepath.Join(tmpDir, ".zshrc")

	// Add a plugin
	err = shell.AddPluginWithMetadata("test-plugin", "/path/to/plugin", "1.0.0", "Test plugin", nil)
	if err != nil {
		t.Fatalf("Failed to add plugin: %v", err)
	}

	// Update plugin version
	if err := shell.UpdatePlugin("test-plugin", "2.0.0"); err != nil {
		t.Fatalf("Failed to update plugin version: %v", err)
	}

	// Verify version update
	updatedPlugin, err := shell.GetPlugin("test-plugin")
	if err != nil {
		t.Fatalf("Failed to get updated plugin: %v", err)
	}
	if updatedPlugin.Version != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got '%s'", updatedPlugin.Version)
	}
}

func TestPluginRemoveAndDisable(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell with a custom config path
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Override the GetConfigPath method to use our temp directory
	testShell := &testShell{
		Shell: shell,
		getConfigPathFunc: func() (string, error) {
			return filepath.Join(tempDir, ".zshrc"), nil
		},
	}

	// Create an initial config file with some content
	initialConfig := `# Initial config
export PATH=$HOME/bin:$PATH
source /path/to/other-plugin.zsh
`
	if err := os.WriteFile(filepath.Join(tempDir, ".zshrc"), []byte(initialConfig), 0644); err != nil {
		t.Fatalf("Failed to create initial config file: %v", err)
	}

	// Add a plugin
	pluginPath := "/path/to/test-plugin.zsh"
	err = testShell.AddPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed to add plugin: %v", err)
	}

	// Verify the plugin was added to the config file
	configContent, err := os.ReadFile(filepath.Join(tempDir, ".zshrc"))
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expectedSource := fmt.Sprintf("source %s\n", pluginPath)
	if !strings.Contains(string(configContent), expectedSource) {
		t.Errorf("Config file does not contain the expected source line: %s", expectedSource)
	}

	// Test disabling the plugin
	err = testShell.DisablePlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to disable plugin: %v", err)
	}

	// Verify the plugin was commented out in the config file
	configContent, err = os.ReadFile(filepath.Join(tempDir, ".zshrc"))
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	expectedCommented := fmt.Sprintf("# source %s\n", pluginPath)
	if !strings.Contains(string(configContent), expectedCommented) {
		t.Errorf("Config file does not contain the expected commented source line: %s", expectedCommented)
	}

	// Test removing the plugin
	err = testShell.RemovePlugin(pluginPath)
	if err != nil {
		t.Errorf("Failed to remove plugin: %v", err)
	}

	// Verify the plugin was removed from the config file
	configContent, err = os.ReadFile(filepath.Join(tempDir, ".zshrc"))
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if strings.Contains(string(configContent), expectedSource) || strings.Contains(string(configContent), expectedCommented) {
		t.Errorf("Config file still contains the plugin source line after removal")
	}
}

func TestPluginRemoveAndDisableErrors(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test shell with a custom config path
	shell, err := New("zsh")
	if err != nil {
		t.Fatalf("Failed to create shell: %v", err)
	}

	// Override the GetConfigPath method to use our temp directory
	testShell := &testShell{
		Shell: shell,
		getConfigPathFunc: func() (string, error) {
			return filepath.Join(tempDir, ".zshrc"), nil
		},
	}

	// Test removing a non-existent plugin
	err = testShell.RemovePlugin("/path/to/non-existent-plugin.zsh")
	if err == nil {
		t.Error("Expected error when removing non-existent plugin, got nil")
	}

	// Test disabling a non-existent plugin
	err = testShell.DisablePlugin("/path/to/non-existent-plugin.zsh")
	if err == nil {
		t.Error("Expected error when disabling non-existent plugin, got nil")
	}

	// Add a plugin
	pluginPath := "/path/to/test-plugin.zsh"
	err = testShell.AddPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed to add plugin: %v", err)
	}

	// Test disabling an already disabled plugin
	err = testShell.DisablePlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed to disable plugin: %v", err)
	}
	err = testShell.DisablePlugin(pluginPath)
	if err == nil {
		t.Error("Expected error when disabling already disabled plugin, got nil")
	}

	// Test removing an already removed plugin
	err = testShell.RemovePlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed to remove plugin: %v", err)
	}
	err = testShell.RemovePlugin(pluginPath)
	if err == nil {
		t.Error("Expected error when removing already removed plugin, got nil")
	}
}

func TestShellTypeValidation(t *testing.T) {
	tests := []struct {
		name        string
		shellType   string
		expectError bool
		configFile  string
		rcFile      string
	}{
		{
			name:        "Valid zsh shell",
			shellType:   "zsh",
			expectError: false,
			configFile:  ".zshrc",
			rcFile:      ".zshenv",
		},
		{
			name:        "Valid bash shell",
			shellType:   "bash",
			expectError: false,
			configFile:  ".bashrc",
			rcFile:      ".bash_profile",
		},
		{
			name:        "Valid fish shell",
			shellType:   "fish",
			expectError: false,
			configFile:  "config.fish",
			rcFile:      "config.fish",
		},
		{
			name:        "Invalid shell type",
			shellType:   "invalid",
			expectError: true,
		},
		{
			name:        "Empty shell type",
			shellType:   "",
			expectError: true,
		},
		{
			name:        "Case insensitive zsh",
			shellType:   "ZSH",
			expectError: false,
			configFile:  ".zshrc",
			rcFile:      ".zshenv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shell, err := New(tt.shellType)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error for invalid shell type, got nil")
				}
				if shell != nil {
					t.Error("Expected nil shell for invalid shell type")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for valid shell type: %v", err)
			}

			if shell.Name != strings.ToLower(tt.shellType) {
				t.Errorf("Expected shell name %s, got %s", strings.ToLower(tt.shellType), shell.Name)
			}

			if shell.ConfigFile != tt.configFile {
				t.Errorf("Expected config file %s, got %s", tt.configFile, shell.ConfigFile)
			}

			if shell.RCFile != tt.rcFile {
				t.Errorf("Expected RC file %s, got %s", tt.rcFile, shell.RCFile)
			}
		})
	}
}

func TestShellConfigContentValidation(t *testing.T) {
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

	// Override the GetConfigPath method to use our temp directory
	testShell := &testShell{
		Shell: shell,
		getConfigPathFunc: func() (string, error) {
			return filepath.Join(tempDir, ".zshrc"), nil
		},
	}

	tests := []struct {
		name          string
		content       string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid config content",
			content:     "# Test config\nexport PATH=$HOME/bin:$PATH\n",
			expectError: false,
		},
		{
			name:          "Invalid shell syntax",
			content:       "export PATH=$HOME/bin:$PATH\nif [ true ]; then\n",
			expectError:   true,
			errorContains: "unclosed",
		},
		{
			name:          "Empty content",
			content:       "",
			expectError:   true,
			errorContains: "empty",
		},
		{
			name:          "Only comments",
			content:       "# This is a comment\n# Another comment\n",
			expectError:   true,
			errorContains: "no configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the config file with the test content
			configPath := filepath.Join(tempDir, ".zshrc")
			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			// Try to validate the config
			err := testShell.ValidateConfig()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error for invalid config, got nil")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for valid config: %v", err)
				}
			}
		})
	}
}

func TestShellRCFileManagement(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		shellType   string
		rcContent   string
		expectError bool
	}{
		{
			name:      "Zsh RC file",
			shellType: "zsh",
			rcContent: `# Zsh environment variables
export EDITOR=vim
export LANG=en_US.UTF-8
`,
			expectError: false,
		},
		{
			name:      "Bash RC file",
			shellType: "bash",
			rcContent: `# Bash environment variables
export EDITOR=vim
export LANG=en_US.UTF-8
`,
			expectError: false,
		},
		{
			name:      "Fish RC file",
			shellType: "fish",
			rcContent: `# Fish environment variables
set -gx EDITOR vim
set -gx LANG en_US.UTF-8
`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test shell
			shell, err := New(tt.shellType)
			if err != nil {
				t.Fatalf("Failed to create shell: %v", err)
			}

			// Override the GetRCPath method to use our temp directory
			testShell := &testShell{
				Shell: shell,
				getConfigPathFunc: func() (string, error) {
					return filepath.Join(tempDir, shell.ConfigFile), nil
				},
			}

			// Create the RC file
			rcPath := filepath.Join(tempDir, shell.RCFile)
			if err := os.WriteFile(rcPath, []byte(tt.rcContent), 0644); err != nil {
				t.Fatalf("Failed to write RC file: %v", err)
			}

			// Test reading the RC file
			content, err := os.ReadFile(rcPath)
			if err != nil {
				t.Fatalf("Failed to read RC file: %v", err)
			}

			if string(content) != tt.rcContent {
				t.Errorf("Expected RC content %q, got %q", tt.rcContent, string(content))
			}

			// Test appending to the RC file
			additionalContent := "\n# Additional settings\n"
			if err := testShell.AppendToRC(additionalContent); err != nil {
				if !tt.expectError {
					t.Errorf("Failed to append to RC file: %v", err)
				}
				return
			}

			// Verify the appended content
			content, err = os.ReadFile(rcPath)
			if err != nil {
				t.Fatalf("Failed to read RC file after append: %v", err)
			}

			expectedContent := tt.rcContent + additionalContent
			if string(content) != expectedContent {
				t.Errorf("Expected RC content after append %q, got %q", expectedContent, string(content))
			}
		})
	}
}

func TestShellEnvironmentVariables(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "shell-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		shellType   string
		envVars     map[string]string
		expectError bool
	}{
		{
			name:      "Zsh environment variables",
			shellType: "zsh",
			envVars: map[string]string{
				"EDITOR": "vim",
				"LANG":   "en_US.UTF-8",
				"PATH":   "$HOME/bin:$PATH",
			},
			expectError: false,
		},
		{
			name:      "Bash environment variables",
			shellType: "bash",
			envVars: map[string]string{
				"EDITOR": "vim",
				"LANG":   "en_US.UTF-8",
				"PATH":   "$HOME/bin:$PATH",
			},
			expectError: false,
		},
		{
			name:      "Fish environment variables",
			shellType: "fish",
			envVars: map[string]string{
				"EDITOR": "vim",
				"LANG":   "en_US.UTF-8",
				"PATH":   "$HOME/bin $PATH",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test shell
			shell, err := New(tt.shellType)
			if err != nil {
				t.Fatalf("Failed to create shell: %v", err)
			}

			// Override the GetRCPath method to use our temp directory
			testShell := &testShell{
				Shell: shell,
				getConfigPathFunc: func() (string, error) {
					return filepath.Join(tempDir, shell.ConfigFile), nil
				},
			}

			// Set environment variables
			for key, value := range tt.envVars {
				if err := testShell.SetEnvVar(key, value); err != nil {
					if !tt.expectError {
						t.Errorf("Failed to set environment variable %s: %v", key, err)
					}
					return
				}
			}

			// Read the RC file to verify the environment variables
			rcPath := filepath.Join(tempDir, shell.RCFile)
			content, err := os.ReadFile(rcPath)
			if err != nil {
				t.Fatalf("Failed to read RC file: %v", err)
			}

			// Verify each environment variable is set correctly
			for key, value := range tt.envVars {
				var expectedLine string
				switch shell.Name {
				case "zsh", "bash":
					expectedLine = fmt.Sprintf("export %s=%s\n", key, value)
				case "fish":
					expectedLine = fmt.Sprintf("set -gx %s %s\n", key, value)
				}

				if !strings.Contains(string(content), expectedLine) {
					t.Errorf("Expected RC file to contain %q, got %q", expectedLine, string(content))
				}
			}

			// Test getting environment variables
			for key, value := range tt.envVars {
				gotValue, err := testShell.GetEnvVar(key)
				if err != nil {
					if !tt.expectError {
						t.Errorf("Failed to get environment variable %s: %v", key, err)
					}
					continue
				}

				if gotValue != value {
					t.Errorf("Expected environment variable %s to be %q, got %q", key, value, gotValue)
				}
			}
		})
	}
} 