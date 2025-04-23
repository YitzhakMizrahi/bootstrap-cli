package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// DefaultManager is the default implementation of interfaces.ShellManager
type DefaultManager struct{}

// NewManager creates a new shell manager
func NewManager() interfaces.ShellManager {
	return &DefaultManager{}
}

// DetectCurrent detects the current user's shell
func (m *DefaultManager) DetectCurrent() (*interfaces.ShellInfo, error) {
	// Try getting shell from SHELL environment variable first
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		return nil, fmt.Errorf("SHELL environment variable not set")
	}

	// Get shell type from path
	shellType := getShellTypeFromPath(shellPath)
	if shellType == "" {
		return nil, fmt.Errorf("unknown shell type: %s", shellPath)
	}

	// Get shell info
	return m.GetInfo(shellType)
}

// ListAvailable returns a list of available shells
func (m *DefaultManager) ListAvailable() ([]*interfaces.ShellInfo, error) {
	shells := []interfaces.Shell{interfaces.Bash, interfaces.Zsh, interfaces.Fish}
	var available []*interfaces.ShellInfo

	for _, shell := range shells {
		if m.IsInstalled(shell) {
			info, err := m.GetInfo(shell)
			if err != nil {
				continue // Skip if we can't get info
			}
			available = append(available, info)
		}
	}

	return available, nil
}

// IsInstalled checks if a specific shell is installed
func (m *DefaultManager) IsInstalled(shell interfaces.Shell) bool {
	_, err := exec.LookPath(string(shell))
	return err == nil
}

// GetInfo returns detailed information about a specific shell
func (m *DefaultManager) GetInfo(shell interfaces.Shell) (*interfaces.ShellInfo, error) {
	if !m.IsInstalled(shell) {
		return nil, fmt.Errorf("shell %s is not installed", shell)
	}

	shellPath, err := exec.LookPath(string(shell))
	if err != nil {
		return nil, fmt.Errorf("failed to find shell path: %w", err)
	}

	// Get shell version
	version, err := getShellVersion(shell, shellPath)
	if err != nil {
		version = "unknown" // Set unknown version but don't fail
	}

	// Check if it's the default shell by comparing resolved paths
	currentShell := os.Getenv("SHELL")
	currentShellPath, err := filepath.EvalSymlinks(currentShell)
	if err != nil {
		currentShellPath = currentShell
	}
	shellRealPath, err := filepath.EvalSymlinks(shellPath)
	if err != nil {
		shellRealPath = shellPath
	}
	isDefault := currentShellPath == shellRealPath

	// Get config files
	configFiles := getShellConfigFiles(shell)

	return &interfaces.ShellInfo{
		Type:        shell,
		Path:        shellPath,
		Version:     version,
		IsDefault:   isDefault,
		IsAvailable: true,
		ConfigFiles: configFiles,
	}, nil
}

// getShellTypeFromPath determines the shell type from its path
func getShellTypeFromPath(path string) interfaces.Shell {
	base := filepath.Base(path)
	switch base {
	case "bash":
		return interfaces.Bash
	case "zsh":
		return interfaces.Zsh
	case "fish":
		return interfaces.Fish
	default:
		return ""
	}
}

// getShellVersion gets the version of a shell
func getShellVersion(shell interfaces.Shell, path string) (string, error) {
	var cmd *exec.Cmd

	switch shell {
	case interfaces.Bash:
		cmd = exec.Command(path, "--version")
	case interfaces.Zsh:
		cmd = exec.Command(path, "--version")
	case interfaces.Fish:
		cmd = exec.Command(path, "--version")
	default:
		return "", fmt.Errorf("unsupported shell type: %s", shell)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Extract version from first line
	version := strings.Split(string(output), "\n")[0]
	return version, nil
}

// getShellConfigFiles returns the configuration files for a shell
func getShellConfigFiles(shell interfaces.Shell) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch shell {
	case interfaces.Bash:
		return []string{
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
			filepath.Join(home, ".profile"),
		}
	case interfaces.Zsh:
		return []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".zprofile"),
			filepath.Join(home, ".zshenv"),
		}
	case interfaces.Fish:
		return []string{
			filepath.Join(home, ".config/fish/config.fish"),
		}
	default:
		return nil
	}
} 