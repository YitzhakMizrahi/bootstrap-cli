package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Test default values
	if config.LogLevel != "info" {
		t.Errorf("Expected LogLevel to be 'info', got %s", config.LogLevel)
	}
	if config.LogFile != "bootstrap.log" {
		t.Errorf("Expected LogFile to be 'bootstrap.log', got %s", config.LogFile)
	}
	if config.PluginDir != "plugins" {
		t.Errorf("Expected PluginDir to be 'plugins', got %s", config.PluginDir)
	}
	if config.MaxPlugins != 10 {
		t.Errorf("Expected MaxPlugins to be 10, got %d", config.MaxPlugins)
	}
	if !config.AutoReload {
		t.Error("Expected AutoReload to be true")
	}
	if config.InstallDir != "~/.bootstrap" {
		t.Errorf("Expected InstallDir to be '~/.bootstrap', got %s", config.InstallDir)
	}
	if config.DefaultShell != "zsh" {
		t.Errorf("Expected DefaultShell to be 'zsh', got %s", config.DefaultShell)
	}
	if config.DotfilesDir != "~/.dotfiles" {
		t.Errorf("Expected DotfilesDir to be '~/.dotfiles', got %s", config.DotfilesDir)
	}

	// Test default tools
	expectedTools := []string{"git", "curl", "wget"}
	if len(config.Tools) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(config.Tools))
	}
	for i, tool := range expectedTools {
		if config.Tools[i] != tool {
			t.Errorf("Expected tool %s at index %d, got %s", tool, i, config.Tools[i])
		}
	}

	// Test default languages
	expectedLanguages := []string{"python", "nodejs"}
	if len(config.Languages) != len(expectedLanguages) {
		t.Errorf("Expected %d languages, got %d", len(expectedLanguages), len(config.Languages))
	}
	for i, lang := range expectedLanguages {
		if config.Languages[i] != lang {
			t.Errorf("Expected language %s at index %d, got %s", lang, i, config.Languages[i])
		}
	}

	// Test default shell configs
	expectedShellConfigs := []string{".zshrc", ".zshenv"}
	if len(config.ShellConfigs) != len(expectedShellConfigs) {
		t.Errorf("Expected %d shell configs, got %d", len(expectedShellConfigs), len(config.ShellConfigs))
	}
	for i, cfg := range expectedShellConfigs {
		if config.ShellConfigs[i] != cfg {
			t.Errorf("Expected shell config %s at index %d, got %s", cfg, i, config.ShellConfigs[i])
		}
	}
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test loading with empty path (should return default config)
	config, err := Load("")
	if err != nil {
		t.Errorf("Failed to load default config: %v", err)
	}
	if config == nil {
		t.Error("Expected non-nil config")
	}

	// Create a test config file
	testConfig := &Config{
		LogLevel:     "debug",
		LogFile:      "test.log",
		PluginDir:    "test-plugins",
		MaxPlugins:   5,
		AutoReload:   false,
		InstallDir:   "/test/install",
		Tools:        []string{"test-tool"},
		Languages:    []string{"test-lang"},
		DefaultShell: "bash",
		ShellConfigs: []string{".bashrc"},
		DotfilesDir:  "/test/dotfiles",
	}

	configPath := filepath.Join(tempDir, "config.json")
	if err := testConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save test config: %v", err)
	}

	// Test loading the test config
	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Errorf("Failed to load test config: %v", err)
	}

	// Verify loaded config matches test config
	if loadedConfig.LogLevel != testConfig.LogLevel {
		t.Errorf("Expected LogLevel %s, got %s", testConfig.LogLevel, loadedConfig.LogLevel)
	}
	if loadedConfig.LogFile != testConfig.LogFile {
		t.Errorf("Expected LogFile %s, got %s", testConfig.LogFile, loadedConfig.LogFile)
	}
	if loadedConfig.PluginDir != testConfig.PluginDir {
		t.Errorf("Expected PluginDir %s, got %s", testConfig.PluginDir, loadedConfig.PluginDir)
	}
	if loadedConfig.MaxPlugins != testConfig.MaxPlugins {
		t.Errorf("Expected MaxPlugins %d, got %d", testConfig.MaxPlugins, loadedConfig.MaxPlugins)
	}
	if loadedConfig.AutoReload != testConfig.AutoReload {
		t.Errorf("Expected AutoReload %v, got %v", testConfig.AutoReload, loadedConfig.AutoReload)
	}
	if loadedConfig.InstallDir != testConfig.InstallDir {
		t.Errorf("Expected InstallDir %s, got %s", testConfig.InstallDir, loadedConfig.InstallDir)
	}
	if loadedConfig.DefaultShell != testConfig.DefaultShell {
		t.Errorf("Expected DefaultShell %s, got %s", testConfig.DefaultShell, loadedConfig.DefaultShell)
	}
	if loadedConfig.DotfilesDir != testConfig.DotfilesDir {
		t.Errorf("Expected DotfilesDir %s, got %s", testConfig.DotfilesDir, loadedConfig.DotfilesDir)
	}

	// Test loading non-existent file
	_, err = Load("non-existent.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config
	config := &Config{
		LogLevel:     "debug",
		LogFile:      "test.log",
		PluginDir:    "test-plugins",
		MaxPlugins:   5,
		AutoReload:   false,
		InstallDir:   "/test/install",
		Tools:        []string{"test-tool"},
		Languages:    []string{"test-lang"},
		DefaultShell: "bash",
		ShellConfigs: []string{".bashrc"},
		DotfilesDir:  "/test/dotfiles",
	}

	// Test saving config
	configPath := filepath.Join(tempDir, "config.json")
	if err := config.Save(configPath); err != nil {
		t.Errorf("Failed to save config: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Read the saved file and verify its contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var savedConfig Config
	if err := json.Unmarshal(data, &savedConfig); err != nil {
		t.Fatalf("Failed to parse saved config: %v", err)
	}

	// Verify saved config matches original config
	if savedConfig.LogLevel != config.LogLevel {
		t.Errorf("Expected LogLevel %s, got %s", config.LogLevel, savedConfig.LogLevel)
	}
	if savedConfig.LogFile != config.LogFile {
		t.Errorf("Expected LogFile %s, got %s", config.LogFile, savedConfig.LogFile)
	}
	if savedConfig.PluginDir != config.PluginDir {
		t.Errorf("Expected PluginDir %s, got %s", config.PluginDir, savedConfig.PluginDir)
	}
	if savedConfig.MaxPlugins != config.MaxPlugins {
		t.Errorf("Expected MaxPlugins %d, got %d", config.MaxPlugins, savedConfig.MaxPlugins)
	}
	if savedConfig.AutoReload != config.AutoReload {
		t.Errorf("Expected AutoReload %v, got %v", config.AutoReload, savedConfig.AutoReload)
	}
	if savedConfig.InstallDir != config.InstallDir {
		t.Errorf("Expected InstallDir %s, got %s", config.InstallDir, savedConfig.InstallDir)
	}
	if savedConfig.DefaultShell != config.DefaultShell {
		t.Errorf("Expected DefaultShell %s, got %s", config.DefaultShell, savedConfig.DefaultShell)
	}
	if savedConfig.DotfilesDir != config.DotfilesDir {
		t.Errorf("Expected DotfilesDir %s, got %s", config.DotfilesDir, savedConfig.DotfilesDir)
	}

	// Test saving to a directory that doesn't exist
	invalidPath := filepath.Join(tempDir, "nonexistent", "config.json")
	if err := config.Save(invalidPath); err != nil {
		t.Errorf("Expected error when saving to nonexistent directory, got %v", err)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "Valid config",
			config: &Config{
				LogLevel:     "info",
				LogFile:      "test.log",
				PluginDir:    "plugins",
				MaxPlugins:   10,
				AutoReload:   true,
				InstallDir:   "/test/install",
				Tools:        []string{"test-tool"},
				Languages:    []string{"test-lang"},
				DefaultShell: "bash",
				ShellConfigs: []string{".bashrc"},
				DotfilesDir:  "/test/dotfiles",
			},
			expectError: false,
		},
		{
			name: "Invalid log level",
			config: &Config{
				LogLevel: "invalid",
			},
			expectError: true,
		},
		{
			name: "Invalid max plugins",
			config: &Config{
				MaxPlugins: -1,
			},
			expectError: true,
		},
		{
			name: "Invalid default shell",
			config: &Config{
				DefaultShell: "invalid",
			},
			expectError: true,
		},
		{
			name: "Empty tools list",
			config: &Config{
				Tools: []string{},
			},
			expectError: true,
		},
		{
			name: "Empty languages list",
			config: &Config{
				Languages: []string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Error("Expected error for invalid config, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for valid config: %v", err)
				}
			}
		})
	}
} 