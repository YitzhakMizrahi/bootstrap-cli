package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Data contains data for shell configuration
type Data struct {
	HomeDir     string
	Plugins     []string
	Aliases     map[string]string
	EnvVars     map[string]string
	Path        []string
	CustomPaths []string
}

// New creates a new Data instance
func New(homeDir string) *Data {
	return &Data{
		HomeDir:     homeDir,
		Plugins:     make([]string, 0),
		Aliases:     make(map[string]string),
		EnvVars:     make(map[string]string),
		Path:        make([]string, 0),
		CustomPaths: make([]string, 0),
	}
}

// GenerateZsh generates the Zsh configuration
func GenerateZsh(data *Data) error {
	configPath := filepath.Join(data.HomeDir, ".zshrc")
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// TODO: Implement Zsh config generation
	return nil
}

// GenerateBash generates the Bash configuration
func GenerateBash(data *Data) error {
	configPath := filepath.Join(data.HomeDir, ".bashrc")
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// TODO: Implement Bash config generation
	return nil
}

// GenerateFish generates the Fish configuration
func GenerateFish(data *Data) error {
	configPath := filepath.Join(data.HomeDir, ".config", "fish", "config.fish")
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// TODO: Implement Fish config generation
	return nil
}

// RestoreZsh restores the Zsh configuration
func RestoreZsh(data *Data, path string) error {
	configPath := filepath.Join(data.HomeDir, ".zshrc")
	return os.Rename(path, configPath)
}

// RestoreBash restores the Bash configuration
func RestoreBash(data *Data, path string) error {
	configPath := filepath.Join(data.HomeDir, ".bashrc")
	return os.Rename(path, configPath)
}

// RestoreFish restores the Fish configuration
func RestoreFish(data *Data, path string) error {
	configPath := filepath.Join(data.HomeDir, ".config", "fish", "config.fish")
	return os.Rename(path, configPath)
} 