package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// DefaultManager implements the ShellManager interface
type DefaultManager struct {
	baseDir string
}

// NewDefaultManager creates a new shell manager
func NewDefaultManager() interfaces.ShellManager {
	return &DefaultManager{
		baseDir: filepath.Join(os.Getenv("HOME"), ".config", "bootstrap-cli"),
	}
}

// DetectCurrent detects the current shell environment
func (m *DefaultManager) DetectCurrent() (*interfaces.ShellInfo, error) {
	// Get current shell from SHELL env var
	currentShell := os.Getenv("SHELL")
	if currentShell == "" {
		return nil, fmt.Errorf("SHELL environment variable not set")
	}

	// Get shell type
	shellType := filepath.Base(currentShell)

	// Validate shell type
	if !interfaces.IsValidShell(shellType) {
		return nil, fmt.Errorf("unknown shell type: %s", shellType)
	}

	// Get shell version
	version, err := m.getShellVersion(currentShell)
	if err != nil {
		version = "unknown"
	}

	// Get config files
	configFiles := m.getConfigFiles(shellType)

	return &interfaces.ShellInfo{
		Current:     shellType,
		Type:        shellType,
		Path:        currentShell,
		Version:     version,
		IsDefault:   true,
		IsAvailable: true,
		ConfigFiles: configFiles,
	}, nil
}

// ListAvailable returns a list of available shells
func (m *DefaultManager) ListAvailable() ([]*interfaces.ShellInfo, error) {
	var shells []*interfaces.ShellInfo

	// Common shell paths
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
			shellType := filepath.Base(path)
			if !interfaces.IsValidShell(shellType) {
				continue
			}

			version, err := m.getShellVersion(path)
			if err != nil {
				version = "unknown"
			}

			shells = append(shells, &interfaces.ShellInfo{
				Current:     shellType,
				Type:        shellType,
				Path:        path,
				Version:     version,
				IsAvailable: true,
				ConfigFiles: m.getConfigFiles(shellType),
			})
		}
	}

	return shells, nil
}

// IsInstalled checks if a specific shell is installed
func (m *DefaultManager) IsInstalled(shell interfaces.ShellType) bool {
	// Common paths for each shell type
	paths := map[interfaces.ShellType][]string{
		interfaces.BashShell: {"/bin/bash", "/usr/bin/bash"},
		interfaces.ZshShell:  {"/bin/zsh", "/usr/bin/zsh"},
		interfaces.FishShell: {"/usr/bin/fish", "/usr/local/bin/fish"},
	}

	// Check if shell exists in any of its common paths
	for _, path := range paths[shell] {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// GetInfo returns detailed information about a specific shell
func (m *DefaultManager) GetInfo(shell interfaces.ShellType) (*interfaces.ShellInfo, error) {
	// Common paths for each shell type
	paths := map[interfaces.ShellType][]string{
		interfaces.BashShell: {"/bin/bash", "/usr/bin/bash"},
		interfaces.ZshShell:  {"/bin/zsh", "/usr/bin/zsh"},
		interfaces.FishShell: {"/usr/bin/fish", "/usr/local/bin/fish"},
	}

	// Find the first existing path for the shell
	var shellPath string
	for _, path := range paths[shell] {
		if _, err := os.Stat(path); err == nil {
			shellPath = path
			break
		}
	}

	if shellPath == "" {
		return nil, fmt.Errorf("shell %s not found", shell)
	}

	version, err := m.getShellVersion(shellPath)
	if err != nil {
		version = "unknown"
	}

	return &interfaces.ShellInfo{
		Current:     string(shell),
		Type:        string(shell),
		Path:        shellPath,
		Version:     version,
		IsAvailable: true,
		ConfigFiles: m.getConfigFiles(string(shell)),
	}, nil
}

// ConfigureShell configures a shell with the specified configuration
func (m *DefaultManager) ConfigureShell(config *interfaces.ShellConfig) error {
	// Get current shell
	shellInfo, err := m.DetectCurrent()
	if err != nil {
		return fmt.Errorf("failed to detect current shell: %w", err)
	}

	// Get the shell's RC file path
	rcFile := m.getDefaultRCFile(shellInfo.Type)

	// Create or append to the RC file
	file, err := os.OpenFile(rcFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open RC file: %w", err)
	}
	defer file.Close()

	// Write environment variables
	if len(config.Exports) > 0 {
		if _, err := file.WriteString("\n# Environment variables\n"); err != nil {
			return fmt.Errorf("failed to write environment header: %w", err)
		}
		for key, value := range config.Exports {
			if _, err := file.WriteString(fmt.Sprintf("export %s=%s\n", key, value)); err != nil {
				return fmt.Errorf("failed to write environment variable: %w", err)
			}
		}
	}

	// Write aliases
	if len(config.Aliases) > 0 {
		if _, err := file.WriteString("\n# Aliases\n"); err != nil {
			return fmt.Errorf("failed to write aliases header: %w", err)
		}
		for name, command := range config.Aliases {
			if _, err := file.WriteString(fmt.Sprintf("alias %s=%s\n", name, command)); err != nil {
				return fmt.Errorf("failed to write alias: %w", err)
			}
		}
	}

	// Write functions
	if len(config.Functions) > 0 {
		if _, err := file.WriteString("\n# Functions\n"); err != nil {
			return fmt.Errorf("failed to write functions header: %w", err)
		}
		for name, body := range config.Functions {
			if _, err := file.WriteString(fmt.Sprintf("%s() {\n%s\n}\n", name, body)); err != nil {
				return fmt.Errorf("failed to write function: %w", err)
			}
		}
	}

	// Add PATH entries
	if len(config.Path) > 0 {
		if _, err := file.WriteString("\n# PATH additions\n"); err != nil {
			return fmt.Errorf("failed to write PATH header: %w", err)
		}
		for _, path := range config.Path {
			if _, err := file.WriteString(fmt.Sprintf("export PATH=%s:$PATH\n", path)); err != nil {
				return fmt.Errorf("failed to write PATH addition: %w", err)
			}
		}
	}

	// Source additional files
	if len(config.Source) > 0 {
		if _, err := file.WriteString("\n# Source additional files\n"); err != nil {
			return fmt.Errorf("failed to write source header: %w", err)
		}
		for _, source := range config.Source {
			if _, err := file.WriteString(fmt.Sprintf("source %s\n", source)); err != nil {
				return fmt.Errorf("failed to write source command: %w", err)
			}
		}
	}

	return nil
}

// Helper functions

func (m *DefaultManager) getShellVersion(shellPath string) (string, error) {
	cmd := exec.Command(shellPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.Split(string(output), "\n")[0], nil
}

func (m *DefaultManager) getConfigFiles(shellType string) []string {
	home := os.Getenv("HOME")
	switch shellType {
	case "bash":
		return []string{
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
			filepath.Join(home, ".profile"),
		}
	case "zsh":
		return []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".zprofile"),
			filepath.Join(home, ".zshenv"),
		}
	case "fish":
		return []string{
			filepath.Join(home, ".config", "fish", "config.fish"),
		}
	default:
		return nil
	}
}

func (m *DefaultManager) getDefaultRCFile(shellType string) string {
	home := os.Getenv("HOME")
	switch shellType {
	case "bash":
		return filepath.Join(home, ".bashrc")
	case "zsh":
		return filepath.Join(home, ".zshrc")
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish")
	default:
		return filepath.Join(home, ".profile")
	}
} 