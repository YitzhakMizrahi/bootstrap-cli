package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// manager implements the interfaces.ShellManager interface.
type manager struct {
	// Potentially add fields like a logger if needed in the future
}

// NewManager creates a new ShellManager.
func NewManager() (interfaces.ShellManager, error) {
	return &manager{}, nil
}

// DetectCurrent detects the current user's shell.
func (m *manager) DetectCurrent() (*interfaces.ShellInfo, error) {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		// Fallback or further probing if SHELL is not set
		// For now, try to find bash or zsh as a desperate measure
		probeShells := []string{"zsh", "bash"}
		for _, s := range probeShells {
			p, err := exec.LookPath(s)
			if err == nil {
				shellPath = p
				break
			}
		}
		if shellPath == "" {
			return nil, fmt.Errorf("SHELL environment variable not set and common shells not found")
		}
	}

	shellName := filepath.Base(shellPath)
	
	// Attempt to get version (simplified)
	version := "unknown"
	// This is a naive version check, real implementation needs per-shell logic
	cmd := exec.Command(shellPath, "--version")
	out, err := cmd.Output()
	if err == nil {
		// Simplistic parsing, actual version string format varies greatly
		versionOutput := string(out)
		if strings.Contains(strings.ToLower(versionOutput), shellName) { // very basic heuristic
			lines := strings.Split(versionOutput, "\n")
			if len(lines) > 0 {
				version = strings.Fields(lines[0])[0] // Highly likely to be wrong or too simple
				parts := strings.Fields(lines[0])
				for _, part := range parts {
					if _, e := exec.LookPath(part); e != nil && len(part) > 1 && (part[0] >= '0' && part[0] <= '9') { // find first numeric-like part
						version = part
						break
					}
				}
			}
		}
	}

	return &interfaces.ShellInfo{
		Current:     shellName, 
		Path:        shellPath,
		Type:        shellName, 
		Version:     version,
		IsAvailable: true, 
		IsDefault:   os.Getenv("SHELL") == shellPath, // True if $SHELL matches this detected shell
		// ConfigFiles: Determine actual config files (e.g., [~/.bashrc] for bash)
	}, nil
}

// ListAvailable returns a list of available shells on the system.
func (m *manager) ListAvailable() ([]*interfaces.ShellInfo, error) {
	available := make([]*interfaces.ShellInfo, 0)
	// Shells to check for. Could be expanded or made configurable.
	potentialShells := []interfaces.ShellType{interfaces.BashShell, interfaces.ZshShell, interfaces.FishShell}

	currentShellEnv := os.Getenv("SHELL")

	for _, shellType := range potentialShells {
		shellName := string(shellType)
		path, err := exec.LookPath(shellName)
		if err == nil { // Shell is found on PATH
			// Simplified version and config file detection
			version := "unknown"
			// Basic version detection (highly simplified)
			cmd := exec.Command(path, "--version")
			output, err := cmd.Output()
			if err == nil {
				lines := strings.Split(string(output), "\n")
				if len(lines) > 0 {
					// Crude parsing, needs to be specific per shell
					parts := strings.Fields(lines[0])
					for _, part := range parts {
						if _, e := exec.LookPath(part); e != nil && len(part) > 1 && (part[0] >= '0' && part[0] <= '9') {
							version = part
							break
						}
					}
					if version == "unknown" && len(parts) > 1 {
						version = parts[1] // fallback to second field if numeric not found
					}
				}
			}

			configFiles := []string{}
			homeDir, _ := os.UserHomeDir()
			switch shellType {
			case interfaces.BashShell:
				configFiles = append(configFiles, filepath.Join(homeDir, ".bashrc"))
				if profilePath := filepath.Join(homeDir, ".bash_profile"); pathExists(profilePath) {
				    configFiles = append(configFiles, profilePath)
				}
			case interfaces.ZshShell:
				configFiles = append(configFiles, filepath.Join(homeDir, ".zshrc"))
			case interfaces.FishShell:
				configFiles = append(configFiles, filepath.Join(homeDir, ".config", "fish", "config.fish"))
			}

			info := &interfaces.ShellInfo{
				Type:        shellName,
				Path:        path,
				Version:     version,
				IsAvailable: true,
				IsDefault:   currentShellEnv == path,
				ConfigFiles: configFiles,
				Current: shellName, // Set Current to the shellName for consistency in ShellInfo
			}
			available = append(available, info)
		}
	}
	return available, nil
}

// IsInstalled checks if a specific shell is installed.
func (m *manager) IsInstalled(shellType interfaces.ShellType) bool {
	_, err := exec.LookPath(string(shellType))
	return err == nil
}

// GetInfo returns detailed information about a specific shell.
func (m *manager) GetInfo(shellType interfaces.ShellType) (*interfaces.ShellInfo, error) {
	path, err := exec.LookPath(string(shellType))
	if err != nil {
		return nil, fmt.Errorf("%s shell not found: %w", shellType, err)
	}

	// Reuse ListAvailable's logic for populating info for a single shell
	list, err := m.ListAvailable() // This is a bit inefficient but reuses logic
	if err != nil {
		return nil, fmt.Errorf("error getting available shells while looking for %s: %w", shellType, err)
	}
	for _, info := range list {
		if info.Type == string(shellType) {
			return info, nil
		}
	}
	return nil, fmt.Errorf("could not retrieve info for %s, though it was found at %s", shellType, path)
}

// ConfigureShell configures a shell with the specified configuration.
func (m *manager) ConfigureShell(config *interfaces.ShellConfig) error {
	return fmt.Errorf("ConfigureShell not yet implemented")
}

// pathExists checks if a path exists. Helper function.
func pathExists(path string) bool {
    _, err := os.Stat(path)
    return !os.IsNotExist(err)
} 