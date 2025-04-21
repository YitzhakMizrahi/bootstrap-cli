package shell

import (
	"os"
	"path/filepath"
	"testing"
)

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