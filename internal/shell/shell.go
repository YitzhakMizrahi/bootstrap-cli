package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Shell represents a shell type
type Shell string

const (
	// Supported shell types
	Bash Shell = "bash"
	Zsh  Shell = "zsh"
	Fish Shell = "fish"
)

// ShellInfo contains information about a shell
type ShellInfo struct {
	Type         Shell
	Path         string
	Version      string
	IsDefault    bool
	IsAvailable  bool
	ConfigFiles  []string
}

// Manager handles shell detection and operations
type Manager interface {
	// DetectCurrent detects the current user's shell
	DetectCurrent() (*ShellInfo, error)
	// ListAvailable returns a list of available shells
	ListAvailable() ([]*ShellInfo, error)
	// IsInstalled checks if a specific shell is installed
	IsInstalled(shell Shell) bool
	// GetInfo returns detailed information about a specific shell
	GetInfo(shell Shell) (*ShellInfo, error)
}

// DefaultManager is the default implementation of Manager
type DefaultManager struct{}

// NewManager creates a new shell manager
func NewManager() Manager {
	return &DefaultManager{}
}

// DetectCurrent detects the current user's shell
func (m *DefaultManager) DetectCurrent() (*ShellInfo, error) {
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
func (m *DefaultManager) ListAvailable() ([]*ShellInfo, error) {
	shells := []Shell{Bash, Zsh, Fish}
	var available []*ShellInfo

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
func (m *DefaultManager) IsInstalled(shell Shell) bool {
	_, err := exec.LookPath(string(shell))
	return err == nil
}

// GetInfo returns detailed information about a specific shell
func (m *DefaultManager) GetInfo(shell Shell) (*ShellInfo, error) {
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

	return &ShellInfo{
		Type:        shell,
		Path:        shellPath,
		Version:     version,
		IsDefault:   isDefault,
		IsAvailable: true,
		ConfigFiles: configFiles,
	}, nil
}

// getShellTypeFromPath determines the shell type from its path
func getShellTypeFromPath(path string) Shell {
	base := filepath.Base(path)
	switch base {
	case "bash":
		return Bash
	case "zsh":
		return Zsh
	case "fish":
		return Fish
	default:
		return ""
	}
}

// getShellVersion gets the version of a shell
func getShellVersion(shell Shell, path string) (string, error) {
	var cmd *exec.Cmd

	switch shell {
	case Bash:
		cmd = exec.Command(path, "--version")
	case Zsh:
		cmd = exec.Command(path, "--version")
	case Fish:
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
func getShellConfigFiles(shell Shell) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	switch shell {
	case Bash:
		return []string{
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
			filepath.Join(home, ".profile"),
		}
	case Zsh:
		return []string{
			filepath.Join(home, ".zshrc"),
			filepath.Join(home, ".zprofile"),
			filepath.Join(home, ".zshenv"),
		}
	case Fish:
		return []string{
			filepath.Join(home, ".config/fish/config.fish"),
		}
	default:
		return nil
	}
} 