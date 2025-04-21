package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Shell installation functions

func installZinit(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("zinit only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zinitPath := filepath.Join(home, ".zinit")
	if _, err := os.Stat(zinitPath); err == nil {
		fmt.Println("✅ Zinit is already installed")
		return nil
	}

	fmt.Println("🔄 Installing Zinit...")
	installCmd := `sh -c "$(curl -fsSL https://git.io/zinit-install)"`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install zinit: %w", err)
	}

	return configureZinit(home)
}

func installOhMyZsh(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("oh-my-zsh only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	ohmyzshPath := filepath.Join(home, ".oh-my-zsh")
	if _, err := os.Stat(ohmyzshPath); err == nil {
		fmt.Println("✅ Oh My Zsh is already installed")
		return nil
	}

	fmt.Println("🔄 Installing Oh My Zsh...")
	installCmd := `sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func installBashIt(shell string) error {
	if shell != "bash" {
		return fmt.Errorf("bash-it only works with bash")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	bashitPath := filepath.Join(home, ".bash_it")
	if _, err := os.Stat(bashitPath); err == nil {
		fmt.Println("✅ Bash-it is already installed")
		return nil
	}

	fmt.Println("🔄 Installing Bash-it...")
	cloneCmd := fmt.Sprintf("git clone --depth=1 https://github.com/Bash-it/bash-it.git %s", bashitPath)
	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone bash-it: %w", err)
	}

	installCmd := fmt.Sprintf("%s/install.sh --silent", bashitPath)
	cmd = exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func installOhMyBash(shell string) error {
	if shell != "bash" {
		return fmt.Errorf("oh-my-bash only works with bash")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	ohMyBashPath := filepath.Join(home, ".oh-my-bash")
	if _, err := os.Stat(ohMyBashPath); err == nil {
		fmt.Println("✅ Oh My Bash is already installed")
		return nil
	}

	fmt.Println("🔄 Installing Oh My Bash...")
	installCmd := `bash -c "$(curl -fsSL https://raw.githubusercontent.com/ohmybash/oh-my-bash/master/tools/install.sh)" --unattended`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func installFisher(shell string) error {
	if shell != "fish" {
		return fmt.Errorf("fisher only works with fish")
	}

	fmt.Println("🔄 Installing Fisher...")
	installCmd := `curl -sL https://git.io/fisher | fish -c "source && fisher install jorgebucaran/fisher"`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func installOhMyFish(shell string) error {
	if shell != "fish" {
		return fmt.Errorf("oh-my-fish only works with fish")
	}

	fmt.Println("🔄 Installing Oh My Fish...")
	installCmd := `curl -L https://get.oh-my.fish | fish`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// Prompt installation and configuration functions

func installStarship(shell string) error {
	fmt.Println("🔄 Installing Starship...")
	installCmd := `curl -sS https://starship.rs/install.sh | sh -s -- -y`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func configureStarship(shell string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configPath string
	var initCmd string

	switch shell {
	case "zsh":
		configPath = filepath.Join(home, ".zshrc")
		initCmd = `eval "$(starship init zsh)"`
	case "bash":
		configPath = filepath.Join(home, ".bashrc")
		initCmd = `eval "$(starship init bash)"`
	case "fish":
		configPath = filepath.Join(home, ".config/fish/config.fish")
		initCmd = "starship init fish | source"
		os.MkdirAll(filepath.Dir(configPath), 0755)
	default:
		return fmt.Errorf("unsupported shell for starship: %s", shell)
	}

	return appendToShellConfig(configPath, "Starship", initCmd)
}

func installPowerlevel10k(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("powerlevel10k only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check for zinit or oh-my-zsh
	zinitPath := filepath.Join(home, ".zinit")
	ohmyzshPath := filepath.Join(home, ".oh-my-zsh")
	
	if _, err := os.Stat(zinitPath); err == nil {
		return configurePowerlevel10kZinit(home)
	} else if _, err := os.Stat(ohmyzshPath); err == nil {
		return configurePowerlevel10kOhMyZsh(home)
	}
	
	return fmt.Errorf("powerlevel10k requires zinit or oh-my-zsh")
}

func configurePowerlevel10k(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("powerlevel10k only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create p10k.zsh if it doesn't exist
	p10kConfigPath := filepath.Join(home, ".p10k.zsh")
	if _, err := os.Stat(p10kConfigPath); err != nil {
		fmt.Println("📝 Creating default p10k.zsh configuration...")
		downloadCmd := fmt.Sprintf("curl -o %s https://raw.githubusercontent.com/romkatv/powerlevel10k/master/config/p10k-lean.zsh", p10kConfigPath)
		cmd := exec.Command("bash", "-c", downloadCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Ignore errors
	}

	return nil
}

func configurePowerlevel10kZinit(home string) error {
	p10kConfig := `
### Added by bootstrap-cli
### Powerlevel10k
zinit ice depth=1; zinit light romkatv/powerlevel10k
### End of Powerlevel10k configuration
`
	return appendToShellConfig(filepath.Join(home, ".zshrc"), "Powerlevel10k", p10kConfig)
}

func configurePowerlevel10kOhMyZsh(home string) error {
	p10kPath := filepath.Join(home, ".oh-my-zsh/custom/themes/powerlevel10k")
	if _, err := os.Stat(p10kPath); err == nil {
		fmt.Println("✅ Powerlevel10k is already installed")
		return nil
	}

	cloneCmd := fmt.Sprintf("git clone --depth=1 https://github.com/romkatv/powerlevel10k.git %s", p10kPath)
	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone powerlevel10k: %w", err)
	}

	// Update .zshrc to use powerlevel10k theme
	zshrcPath := filepath.Join(home, ".zshrc")
	updateCmd := fmt.Sprintf(`sed -i 's/ZSH_THEME=".*"/ZSH_THEME="powerlevel10k\/powerlevel10k"/g' %s`, zshrcPath)
	cmd = exec.Command("bash", "-c", updateCmd)
	return cmd.Run()
}

func installPure(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("pure only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zinitPath := filepath.Join(home, ".zinit")
	if _, err := os.Stat(zinitPath); err != nil {
		return fmt.Errorf("pure requires zinit")
	}

	return configurePure(shell)
}

func configurePure(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("pure only works with zsh")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pureConfig := `
### Added by bootstrap-cli
### Pure prompt
zinit ice pick"async.zsh" src"pure.zsh"
zinit light sindresorhus/pure
### End of Pure configuration
`
	return appendToShellConfig(filepath.Join(home, ".zshrc"), "Pure", pureConfig)
}

func installOhMyPosh(shell string) error {
	fmt.Println("🔄 Installing Oh My Posh...")
	installCmd := `curl -s https://ohmyposh.dev/install.sh | bash -s`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func configureOhMyPosh(shell string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configPath string
	var initCmd string

	switch shell {
	case "zsh":
		configPath = filepath.Join(home, ".zshrc")
		initCmd = `eval "$(oh-my-posh init zsh)"`
	case "bash":
		configPath = filepath.Join(home, ".bashrc")
		initCmd = `eval "$(oh-my-posh init bash)"`
	case "fish":
		configPath = filepath.Join(home, ".config/fish/config.fish")
		initCmd = "oh-my-posh init fish | source"
		os.MkdirAll(filepath.Dir(configPath), 0755)
	default:
		return fmt.Errorf("unsupported shell for oh-my-posh: %s", shell)
	}

	if err := appendToShellConfig(configPath, "Oh My Posh", initCmd); err != nil {
		return err
	}

	// Create themes directory and download default theme
	themesDir := filepath.Join(home, ".poshthemes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	defaultThemePath := filepath.Join(themesDir, "default.omp.json")
	if _, err := os.Stat(defaultThemePath); os.IsNotExist(err) {
		downloadCmd := fmt.Sprintf("curl -o %s https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/main/themes/jandedobbeleer.omp.json", defaultThemePath)
		cmd := exec.Command("bash", "-c", downloadCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Ignore errors
	}

	return nil
}

// Helper functions

func configureZinit(home string) error {
	zinitConfig := `
### Added by bootstrap-cli
### Zinit configuration
source "$HOME/.zinit/bin/zinit.zsh"
zinit light zsh-users/zsh-autosuggestions
zinit light zsh-users/zsh-syntax-highlighting
zinit light zsh-users/zsh-completions
### End of Zinit configuration
`
	return appendToShellConfig(filepath.Join(home, ".zshrc"), "Zinit", zinitConfig)
}

func appendToShellConfig(configPath, name, content string) error {
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", configPath, err)
		}

		if !strings.Contains(string(data), content) {
			if err := appendToFile(configPath, content); err != nil {
				return fmt.Errorf("failed to update %s: %w", configPath, err)
			}
		}
	} else {
		if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", configPath, err)
		}
	}

	return nil
}

func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
} 