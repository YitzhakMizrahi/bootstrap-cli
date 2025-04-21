package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// Plugin represents a shell plugin with additional metadata
type Plugin struct {
	Name         string
	Path         string
	Version      string
	Description  string
	Enabled      bool
	Config       map[string]string
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

// AppendToRC appends content to the shell RC file
func (s *Shell) AppendToRC(content string) error {
	rcPath, err := s.GetRCPath()
	if err != nil {
		return fmt.Errorf("failed to get RC path: %w", err)
	}

	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open RC file: %w", err)
	}
	defer file.Close()

	// Append the content
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to append to RC file: %w", err)
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
	configPath, err := s.GetConfigPath()
	if err != nil {
		return err
	}

	// Read the current configuration
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a pattern to match the plugin source line
	pattern := fmt.Sprintf("source %s\n", pluginPath)
	
	// Replace the pattern with an empty string
	newContent := strings.Replace(string(content), pattern, "", 1)
	
	// Write the updated configuration back to the file
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Removed plugin %s from %s configuration\n", pluginPath, s.Name)
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

// GetPlugin returns a plugin by its path
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

	// Check dependencies
	for _, dep := range plugin.Dependencies {
		depPlugin, err := s.GetPlugin(dep)
		if err != nil {
			return fmt.Errorf("dependency %s not found for plugin %s", dep, pluginPath)
		}
		if !depPlugin.Enabled {
			return fmt.Errorf("dependency %s is not enabled for plugin %s", dep, pluginPath)
		}
	}

	// Enable the plugin
	plugin.Enabled = true

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

// DisablePlugin disables a plugin
func (s *Shell) DisablePlugin(pluginPath string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	if !plugin.Enabled {
		return fmt.Errorf("plugin %s is already disabled", pluginPath)
	}

	// Check if any other plugins depend on this one
	for _, p := range s.Plugins {
		if p.Enabled {
			for _, dep := range p.Dependencies {
				if dep == pluginPath {
					return fmt.Errorf("cannot disable plugin %s: plugin %s depends on it", pluginPath, p.Path)
				}
			}
		}
	}

	// Disable the plugin
	plugin.Enabled = false

	// Remove the plugin from the shell configuration
	configPath, err := s.GetConfigPath()
	if err != nil {
		return err
	}

	// Read the current configuration
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a pattern to match the plugin source line
	pattern := fmt.Sprintf("source %s\n", pluginPath)
	
	// Replace the pattern with a commented out version
	commentedPattern := fmt.Sprintf("# source %s\n", pluginPath)
	newContent := strings.Replace(string(content), pattern, commentedPattern, 1)
	
	// Write the updated configuration back to the file
	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Disabled plugin %s in %s configuration\n", pluginPath, s.Name)
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
		return "", fmt.Errorf("config key %s not found for plugin %s", key, pluginPath)
	}

	return value, nil
}

// UpdatePlugin updates a plugin to a new version
func (s *Shell) UpdatePlugin(pluginPath, newVersion string) error {
	plugin, err := s.GetPlugin(pluginPath)
	if err != nil {
		return err
	}

	// In a real implementation, we would download and install the new version
	plugin.Version = newVersion
	fmt.Printf("Updated plugin %s to version %s\n", pluginPath, newVersion)
	return nil
}

// ValidateConfig validates the shell configuration file content
func (s *Shell) ValidateConfig() error {
	configPath, err := s.GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Check if the file is empty
	if len(content) == 0 {
		return fmt.Errorf("config file is empty")
	}

	// Check if the file contains only comments
	lines := strings.Split(string(content), "\n")
	hasNonComment := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			hasNonComment = true
			break
		}
	}
	if !hasNonComment {
		return fmt.Errorf("config file contains no configuration, only comments")
	}

	// For zsh and bash, we can use the shell's -n option to check syntax
	if s.Name == "zsh" || s.Name == "bash" {
		cmd := exec.Command(s.Name, "-n", configPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("shell syntax validation failed: %w", err)
		}
	}

	// For fish, we can use the fish_indent command to check syntax
	if s.Name == "fish" {
		cmd := exec.Command("fish_indent", "--check", configPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("fish syntax validation failed: %w", err)
		}
	}

	return nil
}

// SetEnvVar sets an environment variable in the shell RC file
func (s *Shell) SetEnvVar(key, value string) error {
	rcPath, err := s.GetRCPath()
	if err != nil {
		return fmt.Errorf("failed to get RC path: %w", err)
	}

	// Read the current RC file content
	content, err := os.ReadFile(rcPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read RC file: %w", err)
	}

	// Create the environment variable line based on shell type
	var envLine string
	switch s.Name {
	case "zsh", "bash":
		envLine = fmt.Sprintf("export %s=%s\n", key, value)
	case "fish":
		envLine = fmt.Sprintf("set -gx %s %s\n", key, value)
	default:
		return fmt.Errorf("unsupported shell type: %s", s.Name)
	}

	// If the file doesn't exist, create it with the new environment variable
	if os.IsNotExist(err) {
		return os.WriteFile(rcPath, []byte(envLine), 0644)
	}

	// Check if the environment variable already exists
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "export "+key+"=") || 
		   strings.HasPrefix(strings.TrimSpace(line), "set -gx "+key+" ") {
			lines[i] = envLine
			return os.WriteFile(rcPath, []byte(strings.Join(lines, "\n")), 0644)
		}
	}

	// Append the new environment variable
	return s.AppendToRC(envLine)
}

// GetEnvVar gets the value of an environment variable from the shell RC file
func (s *Shell) GetEnvVar(key string) (string, error) {
	rcPath, err := s.GetRCPath()
	if err != nil {
		return "", fmt.Errorf("failed to get RC path: %w", err)
	}

	content, err := os.ReadFile(rcPath)
	if err != nil {
		return "", fmt.Errorf("failed to read RC file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch s.Name {
		case "zsh", "bash":
			if strings.HasPrefix(line, "export "+key+"=") {
				return strings.TrimPrefix(line, "export "+key+"="), nil
			}
		case "fish":
			if strings.HasPrefix(line, "set -gx "+key+" ") {
				return strings.TrimPrefix(line, "set -gx "+key+" "), nil
			}
		}
	}

	return "", fmt.Errorf("environment variable %s not found", key)
} 