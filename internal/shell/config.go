package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/spf13/viper"
)

// Config represents shell configuration
type Config struct {
	// Shell type (bash, zsh, fish)
	Shell string
	// Environment variables to set
	EnvVars map[string]string
	// Aliases to configure
	Aliases map[string]string
	// Functions to define
	Functions map[string]string
	// Paths to add to PATH
	Paths []string
	// Logger instance
	Logger interfaces.Logger
}

// NewConfig creates a new shell configuration
func NewConfig(shell string, logger interfaces.Logger) *Config {
	return &Config{
		Shell:     shell,
		EnvVars:   make(map[string]string),
		Aliases:   make(map[string]string),
		Functions: make(map[string]string),
		Paths:     make([]string, 0),
		Logger:    logger,
	}
}

// AddEnvVar adds an environment variable
func (c *Config) AddEnvVar(key, value string) {
	c.EnvVars[key] = value
}

// AddAlias adds a shell alias
func (c *Config) AddAlias(name, command string) {
	c.Aliases[name] = command
}

// AddFunction adds a shell function
func (c *Config) AddFunction(name, body string) {
	c.Functions[name] = body
}

// AddPath adds a path to PATH
func (c *Config) AddPath(path string) {
	c.Paths = append(c.Paths, path)
}

// GetConfigFile returns the path to the shell's config file
func (c *Config) GetConfigFile() string {
	home := os.Getenv("HOME")
	switch c.Shell {
	case "bash":
		return filepath.Join(home, ".bashrc")
	case "zsh":
		return filepath.Join(home, ".zshrc")
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return ""
	}
}

// GetTempConfigFile returns a temporary config file path
func (c *Config) GetTempConfigFile() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("bootstrap-cli-%s-config", c.Shell))
}

// GenerateConfig generates shell configuration content
func (c *Config) GenerateConfig() (string, error) {
	var config strings.Builder

	// Add environment variables
	for key, value := range c.EnvVars {
		switch c.Shell {
		case "fish":
			fmt.Fprintf(&config, "set -gx %s %s\n", key, value)
		default:
			fmt.Fprintf(&config, "export %s=%s\n", key, value)
		}
	}

	// Add PATH updates
	if len(c.Paths) > 0 {
		switch c.Shell {
		case "fish":
			for _, path := range c.Paths {
				fmt.Fprintf(&config, "fish_add_path %s\n", path)
			}
		default:
			paths := strings.Join(c.Paths, ":")
			fmt.Fprintf(&config, "export PATH=%s:$PATH\n", paths)
		}
	}

	// Add aliases
	for name, command := range c.Aliases {
		switch c.Shell {
		case "fish":
			fmt.Fprintf(&config, "alias %s='%s'\n", name, command)
		default:
			fmt.Fprintf(&config, "alias %s='%s'\n", name, command)
		}
	}

	// Add functions
	for name, body := range c.Functions {
		switch c.Shell {
		case "fish":
			fmt.Fprintf(&config, "function %s\n%s\nend\n", name, body)
		default:
			fmt.Fprintf(&config, "%s() {\n%s\n}\n", name, body)
		}
	}

	return config.String(), nil
}

// Apply writes the configuration and returns the command to source it
func (c *Config) Apply() (string, error) {
	// Generate config content
	content, err := c.GenerateConfig()
	if err != nil {
		return "", fmt.Errorf("failed to generate config: %w", err)
	}

	// Write to temp file
	tempFile := c.GetTempConfigFile()
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write config: %w", err)
	}

	// Return the appropriate source command
	switch c.Shell {
	case "fish":
		return fmt.Sprintf("source %s", tempFile), nil
	default:
		return fmt.Sprintf(". %s", tempFile), nil
	}
}

// findShellConfigDir searches for the shell config directory in multiple locations
func findShellConfigDir() (string, error) {
	// Priority order for config locations
	locations := []string{
		// 1. Environment variable
		os.Getenv("BOOTSTRAP_CLI_SHELL_CONFIG"),
		
		// 2. Current working directory
		"config/dotfiles/shell",
		
		// 3. User's config directory
		filepath.Join(os.Getenv("HOME"), ".config", "bootstrap-cli", "shell"),
		
		// 4. System-wide config directory
		filepath.Join("/etc", "bootstrap-cli", "shell"),
		
		// 5. Relative to binary location
		filepath.Join(getBinaryDir(), "config", "dotfiles", "shell"),
	}
	
	// Try each location
	for _, loc := range locations {
		if loc == "" {
			continue
		}
		
		if _, err := os.Stat(loc); err == nil {
			return loc, nil
		}
	}
	
	return "", fmt.Errorf("shell config directory not found")
}

// getBinaryDir returns the directory where the binary is located
func getBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(filepath.Dir(exe)) // Go up one level from the binary
}

// loadShellConfigs loads all shell configurations from the specified directory
func loadShellConfigs(configDir string) (map[string]*interfaces.ShellConfig, error) {
	configs := make(map[string]*interfaces.ShellConfig)

	// If configDir is empty, try to find it
	if configDir == "" {
		var err error
		configDir, err = findShellConfigDir()
		if err != nil {
			return nil, err
		}
	}

	// Read all YAML files in the config directory
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to read shell configs: %w", err)
	}

	fmt.Printf("Found config files: %v\n", files)

	for _, file := range files {
		v := viper.New()
		v.SetConfigFile(file)
		
		if err := v.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file %s: %v\n", file, err)
			return nil, fmt.Errorf("failed to read config file %s: %w", file, err)
		}

		var config interfaces.ShellConfig
		if err := v.Unmarshal(&config); err != nil {
			fmt.Printf("Error unmarshaling config file %s: %v\n", file, err)
			return nil, fmt.Errorf("failed to unmarshal config file %s: %w", file, err)
		}

		// Use the filename without extension as the shell type
		shellType := strings.TrimSuffix(filepath.Base(file), ".yaml")
		fmt.Printf("Loaded config for shell type: %s\n", shellType)
		configs[shellType] = &config
	}

	fmt.Printf("Loaded shell configs: %v\n", configs)
	return configs, nil
}

// getAvailableShells returns a list of available shells based on config and system availability
func getAvailableShells() ([]*interfaces.ShellInfo, error) {
	configs, err := loadShellConfigs("config/dotfiles/shell")
	if err != nil {
		// Fallback to default shells if config loading fails
		return nil, fmt.Errorf("failed to load shell configs: %w", err)
	}

	var available []*interfaces.ShellInfo
	for shellType, config := range configs {
		// Check if the shell binary is available
		if isShellAvailable(shellType) {
			info := &interfaces.ShellInfo{
				Current:     shellType,
				Type:       shellType,
				Path:       findShellPath(shellType),
				Version:    "unknown", // Can be populated later
				IsAvailable: true,
				ConfigFiles: config.Source, // Use Source field instead of Files
			}
			available = append(available, info)
		}
	}

	return available, nil
}

// isShellAvailable checks if a shell is available on the system
func isShellAvailable(shellType string) bool {
	_, err := exec.LookPath(shellType)
	return err == nil
}

// findShellPath returns the full path to the shell binary
func findShellPath(shellType string) string {
	path, _ := exec.LookPath(shellType)
	return path
}

// expandPath expands ~ to the user's home directory in a path
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
} 