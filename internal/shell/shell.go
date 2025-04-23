package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// detectCurrentShell returns the current shell type based on the SHELL environment variable
func detectCurrentShell() interfaces.ShellType {
	shell := os.Getenv("SHELL")
	switch {
	case strings.Contains(shell, "bash"):
		return interfaces.BashShell
	case strings.Contains(shell, "zsh"):
		return interfaces.ZshShell
	case strings.Contains(shell, "fish"):
		return interfaces.FishShell
	default:
		return interfaces.BashShell // Default to bash if unknown
	}
}

// isShellInstalled checks if a specific shell is installed on the system
func isShellInstalled(shell interfaces.ShellType) bool {
	shellPath := os.Getenv("SHELL")
	switch shell {
	case interfaces.BashShell:
		return strings.Contains(shellPath, "bash")
	case interfaces.ZshShell:
		return strings.Contains(shellPath, "zsh")
	case interfaces.FishShell:
		return strings.Contains(shellPath, "fish")
	default:
		return false
	}
}

// getShellConfigFiles returns the list of configuration files for a given shell
func getShellConfigFiles(shell interfaces.ShellType) ([]string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, interfaces.ErrHomeDirNotFound
	}

	switch shell {
	case interfaces.BashShell:
		return []string{
			filepath.Join(home, ".bashrc"),
			filepath.Join(home, ".bash_profile"),
		}, nil
	case interfaces.ZshShell:
		return []string{
			filepath.Join(home, ".zshrc"),
		}, nil
	case interfaces.FishShell:
		return []string{
			filepath.Join(home, ".config/fish/config.fish"),
		}, nil
	default:
		return nil, interfaces.ErrUnsupportedShell
	}
}

// getShellTypeFromPath determines the shell type from a path
func getShellTypeFromPath(path string) interfaces.ShellType {
	switch {
	case strings.Contains(path, "bash"):
		return interfaces.BashShell
	case strings.Contains(path, "zsh"):
		return interfaces.ZshShell
	case strings.Contains(path, "fish"):
		return interfaces.FishShell
	default:
		return ""
	}
} 