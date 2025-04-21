package shell

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Shell represents a shell configuration
type Shell struct {
	Name       string
	ConfigFile string
	RCFile     string
	Plugins    []string
}

// New creates a new Shell
func New(name string) (*Shell, error) {
	switch strings.ToLower(name) {
	case "bash":
		return &Shell{
			Name:       "bash",
			ConfigFile: ".bashrc",
			RCFile:     ".bash_profile",
			Plugins:    []string{},
		}, nil
	case "zsh":
		return &Shell{
			Name:       "zsh",
			ConfigFile: ".zshrc",
			RCFile:     ".zshenv",
			Plugins:    []string{},
		}, nil
	case "fish":
		return &Shell{
			Name:       "fish",
			ConfigFile: "config.fish",
			RCFile:     "config.fish",
			Plugins:    []string{},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported shell: %s", name)
	}
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	return usr.HomeDir, nil
}

// GetConfigPath returns the path to the shell configuration file
func (s *Shell) GetConfigPath() (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, s.ConfigFile), nil
}

// GetRCPath returns the path to the shell RC file
func (s *Shell) GetRCPath() (string, error) {
	homeDir, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, s.RCFile), nil
}

// CreateConfig creates a new shell configuration file
func (s *Shell) CreateConfig(content string) error {
	configPath, err := s.GetConfigPath()
	if err != nil {
		return err
	}

	// Check if the file already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	// Create the file
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	return nil
}

// AppendToConfig appends content to the shell configuration file
func (s *Shell) AppendToConfig(content string) error {
	configPath, err := s.GetConfigPath()
	if err != nil {
		return err
	}

	// Check if the file exists
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Append the content
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to append to config file: %w", err)
	}

	return nil
}

// SetDefaultShell sets the default shell for the user
func (s *Shell) SetDefaultShell() error {
	// This is a simplified implementation
	// In a real implementation, we would use platform-specific commands
	// to set the default shell
	fmt.Printf("Setting default shell to %s\n", s.Name)
	return nil
}

// AddPlugin adds a plugin to the shell
func (s *Shell) AddPlugin(plugin string) error {
	// Check if the plugin is already added
	for _, p := range s.Plugins {
		if p == plugin {
			return fmt.Errorf("plugin %s is already added", plugin)
		}
	}

	// Add the plugin to the list
	s.Plugins = append(s.Plugins, plugin)

	// Add the plugin to the shell configuration
	var pluginConfig string
	switch s.Name {
	case "zsh":
		pluginConfig = fmt.Sprintf("source %s\n", plugin)
	case "bash":
		pluginConfig = fmt.Sprintf("source %s\n", plugin)
	case "fish":
		pluginConfig = fmt.Sprintf("source %s\n", plugin)
	default:
		return fmt.Errorf("unsupported shell for plugin: %s", s.Name)
	}

	return s.AppendToConfig(pluginConfig)
}

// RemovePlugin removes a plugin from the shell
func (s *Shell) RemovePlugin(plugin string) error {
	// Check if the plugin exists
	found := false
	for i, p := range s.Plugins {
		if p == plugin {
			// Remove the plugin from the list
			s.Plugins = append(s.Plugins[:i], s.Plugins[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("plugin %s not found", plugin)
	}

	// Remove the plugin from the shell configuration
	// This is a simplified implementation
	// In a real implementation, we would parse the config file
	// and remove the specific plugin line
	fmt.Printf("Removing plugin %s from %s configuration\n", plugin, s.Name)
	return nil
}

// ListPlugins returns a list of installed plugins
func (s *Shell) ListPlugins() []string {
	return s.Plugins
}

// HasPlugin checks if a plugin is installed
func (s *Shell) HasPlugin(plugin string) bool {
	for _, p := range s.Plugins {
		if p == plugin {
			return true
		}
	}
	return false
} 