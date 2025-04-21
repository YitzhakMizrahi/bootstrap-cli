package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Manager handles shell plugin management operations
type Manager struct {
	homeDir string
}

// NewPluginManager creates a new plugin manager instance
func NewPluginManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	return &Manager{homeDir: homeDir}, nil
}

// installOhMyZsh installs Oh My Zsh
func (m *Manager) installOhMyZsh() error {
	ohmyzshPath := filepath.Join(m.homeDir, ".oh-my-zsh")
	if _, err := os.Stat(ohmyzshPath); err == nil {
		return nil // Already installed
	}

	// Install Oh My Zsh
	cmd := exec.Command("sh", "-c", `curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh`)
	cmd.Env = append(os.Environ(), "RUNZSH=no") // Don't run zsh after installation
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Oh My Zsh: %w", err)
	}

	return nil
}

// installAntigen installs Antigen
func (m *Manager) installAntigen() error {
	antigenPath := filepath.Join(m.homeDir, ".antigen.zsh")
	if _, err := os.Stat(antigenPath); err == nil {
		return nil // Already installed
	}

	// Download antigen.zsh
	cmd := exec.Command("curl", "-L", "git.io/antigen", "-o", antigenPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Antigen: %w", err)
	}

	return nil
}

// installZinit installs Zinit
func (m *Manager) installZinit() error {
	zinitPath := filepath.Join(m.homeDir, ".zinit")
	if _, err := os.Stat(zinitPath); err == nil {
		return nil // Already installed
	}

	// Install Zinit
	cmd := exec.Command("sh", "-c", `mkdir -p "${HOME}/.zinit" && git clone https://github.com/zdharma-continuum/zinit.git "${HOME}/.zinit/bin"`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Zinit: %w", err)
	}

	return nil
}

// installFisherman installs Fisherman
func (m *Manager) installFisherman() error {
	fisherPath := filepath.Join(m.homeDir, ".config/fish/functions/fisher.fish")
	if _, err := os.Stat(fisherPath); err == nil {
		return nil // Already installed
	}

	// Install Fisherman
	cmd := exec.Command("fish", "-c", "curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Fisherman: %w", err)
	}

	return nil
} 