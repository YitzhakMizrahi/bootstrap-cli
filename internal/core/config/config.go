package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// General settings
	LogLevel     string `json:"log_level"`
	LogFile      string `json:"log_file"`
	
	// Plugin settings
	PluginDir    string `json:"plugin_dir"`
	MaxPlugins   int    `json:"max_plugins"`
	AutoReload   bool   `json:"auto_reload"`
	
	// Installation settings
	InstallDir   string   `json:"install_dir"`
	Tools        []string `json:"tools"`
	Languages    []string `json:"languages"`
	
	// Shell settings
	DefaultShell string   `json:"default_shell"`
	ShellConfigs []string `json:"shell_configs"`
	
	// Dotfiles settings
	DotfilesRepo string   `json:"dotfiles_repo"`
	DotfilesDir  string   `json:"dotfiles_dir"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:     "info",
		LogFile:      "bootstrap.log",
		PluginDir:    "plugins",
		MaxPlugins:   10,
		AutoReload:   true,
		InstallDir:   "~/.bootstrap",
		Tools:        []string{"git", "curl", "wget"},
		Languages:    []string{"python", "nodejs"},
		DefaultShell: "zsh",
		ShellConfigs: []string{".zshrc", ".zshenv"},
		DotfilesDir:  "~/.dotfiles",
	}
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	// If no path is provided, use default config
	if path == "" {
		return DefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the config file
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// Save saves configuration to a file
func (c *Config) Save(path string) error {
	// Marshal the config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write the config file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Validate max plugins
	if c.MaxPlugins < 0 {
		return fmt.Errorf("max plugins cannot be negative")
	}

	// Validate default shell
	validShells := map[string]bool{
		"bash": true,
		"zsh":  true,
		"fish": true,
	}
	if !validShells[c.DefaultShell] {
		return fmt.Errorf("invalid default shell: %s", c.DefaultShell)
	}

	// Validate tools list
	if len(c.Tools) == 0 {
		return fmt.Errorf("tools list cannot be empty")
	}

	// Validate languages list
	if len(c.Languages) == 0 {
		return fmt.Errorf("languages list cannot be empty")
	}

	return nil
} 