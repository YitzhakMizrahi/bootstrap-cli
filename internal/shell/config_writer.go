package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
)

// DotfilesStrategy defines how to handle existing dotfiles
type DotfilesStrategy int

const (
	// MergeWithExisting merges new configurations with existing ones
	MergeWithExisting DotfilesStrategy = iota
	// SkipIfExists skips adding configurations if they already exist
	SkipIfExists
	// ReplaceExisting replaces existing configurations with new ones
	ReplaceExisting
)

// ConfigWriter handles shell configuration file management
type ConfigWriter interface {
	// WriteConfig writes shell configurations to the appropriate file
	WriteConfig(configs []string, strategy DotfilesStrategy) error
	// AddToPath adds a directory to the PATH environment variable
	AddToPath(path string) error
	// SetEnvVar sets an environment variable
	SetEnvVar(name, value string) error
	// AddAlias adds a shell alias
	AddAlias(name, command string) error
	// HasConfig checks if a configuration exists
	HasConfig(config string) bool
}

// DefaultConfigWriter implements ConfigWriter
type DefaultConfigWriter struct {
	logger *log.Logger
	shell  Shell
	pm     packages.Manager
}

// NewConfigWriter creates a new DefaultConfigWriter
func NewConfigWriter() (*DefaultConfigWriter, error) {
	sysInfo, err := system.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %w", err)
	}

	shellInfo, err := NewManager().DetectCurrent()
	if err != nil {
		return nil, fmt.Errorf("failed to detect shell: %w", err)
	}

	logger := log.New(log.InfoLevel)
	pm, err := packages.NewPackageManager(sysInfo.OS)
	if err != nil {
		return nil, fmt.Errorf("failed to create package manager: %w", err)
	}

	return &DefaultConfigWriter{
		logger: logger,
		shell:  shellInfo.Type,
		pm:     pm,
	}, nil
}

// WriteConfig writes shell configurations to the appropriate file
func (w *DefaultConfigWriter) WriteConfig(configs []string, strategy DotfilesStrategy) error {
	// Get shell config file path
	configFile := w.getConfigFile()
	if configFile == "" {
		return fmt.Errorf("no config file found for shell %s", w.shell)
	}

	// Read existing config if it exists
	var existingConfig string
	if _, err := os.Stat(configFile); err == nil {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		existingConfig = string(data)
	}

	// Process each config
	var newConfigs []string
	for _, config := range configs {
		if strategy != ReplaceExisting && w.HasConfig(config) {
			if strategy == SkipIfExists {
				continue
			}
			// For MergeWithExisting, we'll keep both
		}
		newConfigs = append(newConfigs, config)
	}

	// Write the new config
	var content string
	if strategy == ReplaceExisting {
		content = strings.Join(newConfigs, "\n")
		if len(newConfigs) > 0 {
			content += "\n"
		}
	} else {
		content = existingConfig
		if len(newConfigs) > 0 {
			if content != "" && !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += strings.Join(newConfigs, "\n") + "\n"
		}
	}

	// Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddToPath adds a directory to the PATH environment variable
func (w *DefaultConfigWriter) AddToPath(path string) error {
	config := fmt.Sprintf("export PATH=%s:$PATH", path)
	return w.WriteConfig([]string{config}, MergeWithExisting)
}

// SetEnvVar sets an environment variable
func (w *DefaultConfigWriter) SetEnvVar(name, value string) error {
	config := fmt.Sprintf("export %s=%s", name, value)
	return w.WriteConfig([]string{config}, MergeWithExisting)
}

// AddAlias adds a shell alias
func (w *DefaultConfigWriter) AddAlias(name, command string) error {
	config := fmt.Sprintf("alias %s='%s'", name, command)
	return w.WriteConfig([]string{config}, MergeWithExisting)
}

// HasConfig checks if a configuration exists
func (w *DefaultConfigWriter) HasConfig(config string) bool {
	configFile := w.getConfigFile()
	if configFile == "" {
		return false
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), config)
}

// getConfigFile returns the appropriate config file path for the shell
func (w *DefaultConfigWriter) getConfigFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		w.logger.Error("Failed to get user home directory: %v", err)
		return ""
	}

	switch w.shell {
	case Bash:
		return filepath.Join(home, ".bashrc")
	case Zsh:
		return filepath.Join(home, ".zshrc")
	case Fish:
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return ""
	}
} 