package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// Manager handles shell operations
type Manager struct {
	configLoader *config.ConfigLoader
}

// NewManager creates a new shell manager
func NewManager() *Manager {
	return &Manager{
		configLoader: config.NewConfigLoader("config"),
	}
}

// DetectCurrent detects the current shell environment
func (m *Manager) DetectCurrent() (*interfaces.ShellInfo, error) {
	// Get current shell from SHELL env var
	currentShell := os.Getenv("SHELL")
	if currentShell == "" {
		return nil, fmt.Errorf("SHELL environment variable not set")
	}

	// Get available shells
	available, err := m.detectAvailableShells()
	if err != nil {
		return nil, fmt.Errorf("failed to detect available shells: %w", err)
	}

	// Get default shell
	defaultShell, err := m.getDefaultShell()
	if err != nil {
		return nil, fmt.Errorf("failed to get default shell: %w", err)
	}

	return &interfaces.ShellInfo{
		Current:     filepath.Base(currentShell),
		Available:   available,
		DefaultPath: defaultShell,
	}, nil
}

// detectAvailableShells detects which shells are available on the system
func (m *Manager) detectAvailableShells() ([]string, error) {
	var available []string

	// Check for common shell locations
	shellPaths := []string{
		"/bin/bash",
		"/usr/bin/bash",
		"/bin/zsh",
		"/usr/bin/zsh",
		"/usr/bin/fish",
		"/usr/local/bin/fish",
	}

	for _, path := range shellPaths {
		if _, err := os.Stat(path); err == nil {
			shell := filepath.Base(path)
			if interfaces.IsValidShell(shell) {
				available = append(available, shell)
			}
		}
	}

	return available, nil
}

// getDefaultShell gets the system default shell
func (m *Manager) getDefaultShell() (string, error) {
	// Try to get from /etc/passwd
	output, err := exec.Command("getent", "passwd", os.Getenv("USER")).Output()
	if err == nil {
		fields := strings.Split(string(output), ":")
		if len(fields) >= 7 {
			return fields[6], nil
		}
	}

	// Fallback to $SHELL
	return os.Getenv("SHELL"), nil
}

// ConfigureShell configures a shell with the specified type
func (m *Manager) ConfigureShell(shellType string) error {
	// Load shell-specific dotfile configuration
	dotfiles, err := m.configLoader.LoadDotfiles()
	if err != nil {
		return fmt.Errorf("failed to load dotfiles: %w", err)
	}

	// Find the configuration for the selected shell
	var shellConfig *install.Dotfile
	for _, dotfile := range dotfiles {
		if strings.Contains(dotfile.Name, shellType) {
			shellConfig = dotfile
			break
		}
	}

	if shellConfig == nil {
		return fmt.Errorf("no configuration found for shell: %s", shellType)
	}

	// Convert install.Dotfile to interfaces.Dotfile
	interfaceDotfile := &interfaces.Dotfile{
		Name:        shellConfig.Name,
		Description: shellConfig.Description,
		Category:    shellConfig.Category,
		Tags:        shellConfig.Tags,
		Files: []interfaces.FileConfig{
			{
				Source:      shellConfig.Source,
				Destination: shellConfig.Target,
				Type:        "content", // Default to content type
				Backup:      true,      // Default to backing up
			},
		},
	}

	// For MVP: Log the dotfile configuration that will be applied
	fmt.Printf("Shell configuration to be applied: %+v\n", interfaceDotfile)

	// TODO: Use interfaceDotfile with the dotfiles manager to apply the configuration
	// This will be implemented when the dotfiles manager is integrated with the shell manager

	return nil
} 