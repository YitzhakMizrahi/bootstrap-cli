package shell

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Shell represents a shell configuration
type Shell struct {
	Name       string
	ConfigFile string
	RCFile     string
}

// New creates a new Shell
func New(name string) (*Shell, error) {
	switch strings.ToLower(name) {
	case "bash":
		return &Shell{
			Name:       "bash",
			ConfigFile: ".bashrc",
			RCFile:     ".bash_profile",
		}, nil
	case "zsh":
		return &Shell{
			Name:       "zsh",
			ConfigFile: ".zshrc",
			RCFile:     ".zshenv",
		}, nil
	case "fish":
		return &Shell{
			Name:       "fish",
			ConfigFile: "config.fish",
			RCFile:     "config.fish",
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

// SetDefaultShell sets the default shell for the user
func (s *Shell) SetDefaultShell() error {
	// This is a simplified implementation
	// In a real implementation, we would use platform-specific commands
	// to set the default shell
	fmt.Printf("Setting default shell to %s\n", s.Name)
	return nil
} 