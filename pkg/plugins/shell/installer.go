package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Plugin installation functions

// InstallZinit installs the Zinit plugin manager
func InstallZinit(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("zinit is only compatible with zsh")
	}
	cmd := exec.Command("sh", "-c", "curl -fsSL https://git.io/zinit-install | sh")
	return cmd.Run()
}

// InstallOhMyZsh installs the Oh My Zsh framework
func InstallOhMyZsh(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("oh-my-zsh is only compatible with zsh")
	}
	cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh")
	return cmd.Run()
}

// InstallBashIt installs the Bash-it framework
func InstallBashIt(shell string) error {
	if shell != "bash" {
		return fmt.Errorf("bash-it is only compatible with bash")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/Bash-it/bash-it.git", filepath.Join(homeDir, ".bash_it"))
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command(filepath.Join(homeDir, ".bash_it/install.sh"), "--silent")
	return cmd.Run()
}

// InstallOhMyBash installs the Oh My Bash framework
func InstallOhMyBash(shell string) error {
	if shell != "bash" {
		return fmt.Errorf("oh-my-bash is only compatible with bash")
	}
	cmd := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmybash/oh-my-bash/master/tools/install.sh | bash")
	return cmd.Run()
}

// InstallFisher installs the Fisher plugin manager
func InstallFisher(shell string) error {
	if shell != "fish" {
		return fmt.Errorf("fisher is only compatible with fish")
	}
	cmd := exec.Command("fish", "-c", "curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher")
	return cmd.Run()
}

// InstallOhMyFish installs the Oh My Fish framework
func InstallOhMyFish(shell string) error {
	if shell != "fish" {
		return fmt.Errorf("oh-my-fish is only compatible with fish")
	}
	cmd := exec.Command("fish", "-c", "curl -L https://get.oh-my.fish | fish")
	return cmd.Run()
}

// Prompt installation and configuration functions

// InstallStarship installs the Starship prompt
func InstallStarship(shell string) error {
	cmd := exec.Command("sh", "-c", "curl -fsSL https://starship.rs/install.sh | sh")
	return cmd.Run()
}

// ConfigureStarship configures the Starship prompt
func ConfigureStarship(shell string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return os.WriteFile(filepath.Join(configDir, "starship.toml"), []byte(""), 0644)
}

// InstallPowerlevel10k installs the Powerlevel10k theme
func InstallPowerlevel10k(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("powerlevel10k is only compatible with zsh")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cmd := exec.Command("git", "clone", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", filepath.Join(homeDir, ".powerlevel10k"))
	return cmd.Run()
}

// ConfigurePowerlevel10k configures the Powerlevel10k theme
func ConfigurePowerlevel10k(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("powerlevel10k is only compatible with zsh")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	return os.WriteFile(filepath.Join(homeDir, ".p10k.zsh"), []byte(""), 0644)
}

// InstallPure installs the Pure prompt
func InstallPure(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("pure is only compatible with zsh")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	cmd := exec.Command("git", "clone", "https://github.com/sindresorhus/pure.git", filepath.Join(homeDir, ".pure"))
	return cmd.Run()
}

// ConfigurePure configures the Pure prompt
func ConfigurePure(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("pure is only compatible with zsh")
	}
	return nil // Pure doesn't require additional configuration
}

// InstallOhMyPosh installs the Oh My Posh prompt
func InstallOhMyPosh(shell string) error {
	cmd := exec.Command("sh", "-c", "curl -s https://ohmyposh.dev/install.sh | bash -s")
	return cmd.Run()
}

// ConfigureOhMyPosh configures the Oh My Posh prompt
func ConfigureOhMyPosh(shell string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config", "oh-my-posh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{}"), 0644)
} 