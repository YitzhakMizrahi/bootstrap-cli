package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
)

// DefaultConfigWriter implements interfaces.ShellConfigWriter
type DefaultConfigWriter struct {
	logger *log.Logger
	shell  interfaces.Shell
	pm     interfaces.PackageManager
}

// NewConfigWriter creates a new shell config writer
func NewConfigWriter() (interfaces.ShellConfigWriter, error) {
	shellInfo, err := NewManager().DetectCurrent()
	if err != nil {
		return nil, fmt.Errorf("failed to detect shell: %w", err)
	}

	logger := log.New(log.InfoLevel)
	
	// Use the factory to get the package manager
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create package manager: %w", err)
	}

	return &DefaultConfigWriter{
		logger: logger,
		shell:  interfaces.Shell(shellInfo.Current),
		pm:     pm,
	}, nil
}

// WriteConfig writes shell configurations to the appropriate file
func (w *DefaultConfigWriter) WriteConfig(configs []string, strategy interfaces.DotfilesStrategy) error {
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
		if strategy != interfaces.ReplaceExisting && w.HasConfig(config) {
			if strategy == interfaces.SkipIfExists {
				continue
			}
			// For MergeWithExisting, we'll keep both
		}
		newConfigs = append(newConfigs, config)
	}

	// Write the new config
	var content string
	if strategy == interfaces.ReplaceExisting {
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
	return w.WriteConfig([]string{config}, interfaces.MergeWithExisting)
}

// SetEnvVar sets an environment variable
func (w *DefaultConfigWriter) SetEnvVar(name, value string) error {
	config := fmt.Sprintf("export %s=%s", name, value)
	return w.WriteConfig([]string{config}, interfaces.MergeWithExisting)
}

// AddAlias adds a shell alias
func (w *DefaultConfigWriter) AddAlias(name, command string) error {
	config := fmt.Sprintf("alias %s='%s'", name, command)
	return w.WriteConfig([]string{config}, interfaces.MergeWithExisting)
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
	case interfaces.Bash:
		return filepath.Join(home, ".bashrc")
	case interfaces.Zsh:
		return filepath.Join(home, ".zshrc")
	case interfaces.Fish:
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return ""
	}
} 