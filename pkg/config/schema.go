package config

import (
	"os"
	"path/filepath"
)

// SystemType represents the type of Unix-like system
type SystemType string

const (
	Linux  SystemType = "linux"
	Darwin SystemType = "darwin"
	WSL    SystemType = "wsl"
)

// ShellType represents the user's preferred shell
type ShellType string

const (
	Bash ShellType = "bash"
	Zsh  ShellType = "zsh"
	Fish ShellType = "fish"
)

// PackageManager represents the system's package manager
type PackageManager string

const (
	Apt   PackageManager = "apt"
	Brew  PackageManager = "brew"
	Dnf   PackageManager = "dnf"
	Pacman PackageManager = "pacman"
)

// DotfilesConfig represents the dotfiles configuration
type DotfilesConfig struct {
	// Repository URL if using existing dotfiles
	RepoURL string `json:"repo_url,omitempty"`
	
	// Base directory for dotfiles
	BaseDir string `json:"base_dir"`
	
	// Whether to backup existing dotfiles
	Backup bool `json:"backup"`
	
	// Map of source files to target locations
	Symlinks map[string]string `json:"symlinks"`
}

// ToolConfig represents a development tool configuration
type ToolConfig struct {
	Name         string            `json:"name"`
	Enabled      bool             `json:"enabled"`
	Version      string           `json:"version,omitempty"`
	EnvVars      map[string]string `json:"env_vars,omitempty"`
	Dependencies []string         `json:"dependencies,omitempty"`
}

// Config represents the complete user configuration
type Config struct {
	// System information
	System struct {
		Type           SystemType     `json:"type"`
		PackageManager PackageManager `json:"package_manager"`
	} `json:"system"`

	// Shell configuration
	Shell struct {
		Type    ShellType `json:"type"`
		Plugins []string  `json:"plugins,omitempty"`
	} `json:"shell"`

	// Development tools
	Tools map[string]ToolConfig `json:"tools"`

	// Dotfiles configuration
	Dotfiles DotfilesConfig `json:"dotfiles"`
}

// DetectSystem determines the system type and package manager
func DetectSystem() (SystemType, PackageManager) {
	// Check for WSL
	if _, err := os.Stat("/proc/sys/fs/binfmt_misc/WSLInterop"); err == nil {
		return WSL, Apt
	}

	// Check for Darwin (macOS)
	if _, err := os.Stat("/usr/local/bin/brew"); err == nil {
		return Darwin, Brew
	}

	// Check for Linux package managers
	for _, pm := range []struct {
		path string
		pm   PackageManager
	}{
		{"/usr/bin/apt", Apt},
		{"/usr/bin/dnf", Dnf},
		{"/usr/bin/pacman", Pacman},
	} {
		if _, err := os.Stat(pm.path); err == nil {
			return Linux, pm.pm
		}
	}

	// Default to Linux with apt
	return Linux, Apt
}

// GetDefaultConfig returns a default configuration based on the system
func GetDefaultConfig() (*Config, error) {
	sysType, pkgManager := DetectSystem()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	cfg.System.Type = sysType
	cfg.System.PackageManager = pkgManager

	// Set default shell
	if _, err := os.Stat("/bin/zsh"); err == nil {
		cfg.Shell.Type = Zsh
	} else {
		cfg.Shell.Type = Bash
	}

	// Set default dotfiles configuration
	cfg.Dotfiles = DotfilesConfig{
		BaseDir: filepath.Join(homeDir, ".dotfiles"),
		Backup:  true,
		Symlinks: map[string]string{
			".zshrc":     filepath.Join(homeDir, ".zshrc"),
			".bashrc":    filepath.Join(homeDir, ".bashrc"),
			".gitconfig": filepath.Join(homeDir, ".gitconfig"),
			".tmux.conf": filepath.Join(homeDir, ".tmux.conf"),
		},
	}

	// Initialize tools map
	cfg.Tools = make(map[string]ToolConfig)

	return cfg, nil
} 