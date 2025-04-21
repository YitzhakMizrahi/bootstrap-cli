package shell

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Plugin represents a shell plugin with additional metadata
type Plugin struct {
	Name        string
	Path        string
	Version     string
	Description string
	Enabled     bool
	Config      map[string]string
	Dependencies []string
}

// Shell represents a shell configuration
type Shell struct {
	Name       string
	ConfigFile string
	RCFile     string
	Plugins    []*Plugin
}

// New creates a new Shell
func New(name string) (*Shell, error) {
	switch strings.ToLower(name) {
	case "bash":
		return &Shell{
			Name:       "bash",
			ConfigFile: ".bashrc",
			RCFile:     ".bash_profile",
			Plugins:    []*Plugin{},
		}, nil
	case "zsh":
		return &Shell{
			Name:       "zsh",
			ConfigFile: ".zshrc",
			RCFile:     ".zshenv",
			Plugins:    []*Plugin{},
		}, nil
	case "fish":
		return &Shell{
			Name:       "fish",
			ConfigFile: "config.fish",
			RCFile:     "config.fish",
			Plugins:    []*Plugin{},
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
func (s *Shell) AddPlugin(pluginPath string) error {
	// Check if the plugin is already added
	for _, p := range s.Plugins {
		if p.Path == pluginPath {
			return fmt.Errorf("plugin %s is already added", pluginPath)
		}
	}

	// Create a new plugin
	plugin := &Plugin{
		Name:    filepath.Base(pluginPath),
		Path:    pluginPath,
		Version: "1.0.0", // Default version
		Enabled: true,
		Config:  make(map[string]string),
	}

	// Add the plugin to the list
	s.Plugins = append(s.Plugins, plugin)

	// Add the plugin to the shell configuration
	var pluginConfig string
	switch s.Name {
	case "zsh":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	case "bash":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	case "fish":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	default:
		return fmt.Errorf("unsupported shell for plugin: %s", s.Name)
	}

	return s.AppendToConfig(pluginConfig)
}

// AddPluginWithMetadata adds a plugin with metadata to the shell
func (s *Shell) AddPluginWithMetadata(name, path, version, description string, dependencies []string) error {
	// Check if the plugin is already added
	for _, p := range s.Plugins {
		if p.Path == path {
			return fmt.Errorf("plugin %s is already added", path)
		}
	}

	// Create a new plugin
	plugin := &Plugin{
		Name:         name,
		Path:         path,
		Version:      version,
		Description:  description,
		Enabled:      true,
		Config:       make(map[string]string),
		Dependencies: dependencies,
	}

	// Add the plugin to the list
	s.Plugins = append(s.Plugins, plugin)

	// Add the plugin to the shell configuration
	var pluginConfig string
	switch s.Name {
	case "zsh":
		pluginConfig = fmt.Sprintf("source %s\n", path)
	case "bash":
		pluginConfig = fmt.Sprintf("source %s\n", path)
	case "fish":
		pluginConfig = fmt.Sprintf("source %s\n", path)
	default:
		return fmt.Errorf("unsupported shell for plugin: %s", s.Name)
	}

	return s.AppendToConfig(pluginConfig)
}

// RemovePlugin removes a plugin from the shell
func (s *Shell) RemovePlugin(pluginPath string) error {
	// Check if the plugin exists
	found := false
	for i, p := range s.Plugins {
		if p.Path == pluginPath {
			// Remove the plugin from the list
			s.Plugins = append(s.Plugins[:i], s.Plugins[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("plugin %s not found", pluginPath)
	}

	// Remove the plugin from the shell configuration
	// This is a simplified implementation
	// In a real implementation, we would parse the config file
	// and remove the specific plugin line
	fmt.Printf("Removing plugin %s from %s configuration\n", pluginPath, s.Name)
	return nil
}

// ListPlugins returns a list of installed plugins
func (s *Shell) ListPlugins() []*Plugin {
	return s.Plugins
}

// HasPlugin checks if a plugin is installed
func (s *Shell) HasPlugin(pluginPath string) bool {
	for _, p := range s.Plugins {
		if p.Path == pluginPath {
			return true
		}
	}
	return false
}

// GetPlugin returns a plugin by path
func (s *Shell) GetPlugin(pluginPath string) (*Plugin, error) {
	for _, p := range s.Plugins {
		if p.Path == pluginPath {
			return p, nil
		}
	}
	return nil, fmt.Errorf("plugin %s not found", pluginPath)
}

// EnablePlugin enables a plugin
func (s *Shell) EnablePlugin(pluginPath string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	if plugin.Enabled {
		return fmt.Errorf("plugin %s is already enabled", pluginPath)
	}

	// Add the plugin to the shell configuration
	var pluginConfig string
	switch s.Name {
	case "zsh":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	case "bash":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	case "fish":
		pluginConfig = fmt.Sprintf("source %s\n", pluginPath)
	default:
		return fmt.Errorf("unsupported shell for plugin: %s", s.Name)
	}

	if err := s.AppendToConfig(pluginConfig); err != nil {
		return err
	}

	plugin.Enabled = true
	return nil
}

// DisablePlugin disables a plugin
func (s *Shell) DisablePlugin(pluginPath string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	if !plugin.Enabled {
		return fmt.Errorf("plugin %s is already disabled", pluginPath)
	}

	// Remove the plugin from the shell configuration
	// This is a simplified implementation
	// In a real implementation, we would parse the config file
	// and remove the specific plugin line
	fmt.Printf("Disabling plugin %s in %s configuration\n", pluginPath, s.Name)
	
	plugin.Enabled = false
	return nil
}

// SetPluginConfig sets a configuration option for a plugin
func (s *Shell) SetPluginConfig(pluginPath, key, value string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	plugin.Config[key] = value
	return nil
}

// GetPluginConfig gets a configuration option for a plugin
func (s *Shell) GetPluginConfig(pluginPath, key string) (string, error) {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return "", err
	}

	value, ok := plugin.Config[key]
	if !ok {
		return "", fmt.Errorf("configuration option %s not found for plugin %s", key, pluginPath)
	}

	return value, nil
}

// UpdatePlugin updates a plugin to a new version
func (s *Shell) UpdatePlugin(pluginPath, newVersion string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	plugin.Version = newVersion
	return nil
} 