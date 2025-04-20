package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/platform"
)

// Shell defines installation methods for different shells
type Shell struct {
	Name         string
	Description  string
	BrewPackage  string
	AptPackage   string
	DnfPackage   string
	PacmanPackage string
	ZypperPackage string
	ChocoPackage string
	InstallCmd   func() error
}

// ShellOptions maps shell names to their installation info
var ShellOptions = map[string]Shell{
	"zsh": {
		Name:         "Zsh",
		Description:  "The Z Shell",
		BrewPackage:  "zsh",
		AptPackage:   "zsh",
		DnfPackage:   "zsh",
		PacmanPackage: "zsh",
		ZypperPackage: "zsh",
		ChocoPackage: "zsh",
	},
	"bash": {
		Name:         "Bash",
		Description:  "Bourne Again Shell",
		BrewPackage:  "bash",
		AptPackage:   "bash",
		DnfPackage:   "bash",
		PacmanPackage: "bash",
		ZypperPackage: "bash",
		ChocoPackage: "bash",
	},
	"fish": {
		Name:         "Fish",
		Description:  "Friendly Interactive Shell",
		BrewPackage:  "fish",
		AptPackage:   "fish",
		DnfPackage:   "fish",
		PacmanPackage: "fish",
		ZypperPackage: "fish",
		ChocoPackage: "fish",
	},
}

// PluginManager defines installation methods for different plugin managers
type PluginManager struct {
	Name        string
	Description string
	CompatibleShells []string
	InstallCmd  func(shell string) error
}

// Prompt defines installation methods for different prompts
type Prompt struct {
	Name            string
	Description     string
	CompatibleShells []string
	BrewPackage     string
	AptPackage      string
	DnfPackage      string
	PacmanPackage   string
	ZypperPackage   string
	ChocoPackage    string
	InstallCmd      func(shell string) error
}

// Declare variables without initializing
var PluginManagerOptions map[string]PluginManager
var PromptOptions map[string]Prompt

// Initialize in init function
func init() {
	PluginManagerOptions = map[string]PluginManager{
		"zinit": {
			Name:        "Zinit",
			Description: "Fast and feature-rich plugin manager for Zsh",
			CompatibleShells: []string{"zsh"},
			InstallCmd:  installZinit,
		},
		"oh-my-zsh": {
			Name:        "Oh My Zsh",
			Description: "Community-driven framework for Zsh",
			CompatibleShells: []string{"zsh"},
			InstallCmd:  installOhMyZsh,
		},
		"bash-it": {
			Name:        "Bash-it",
			Description: "Community framework for Bash",
			CompatibleShells: []string{"bash"},
			InstallCmd:  installBashIt,
		},
		"oh-my-bash": {
			Name:        "Oh My Bash",
			Description: "Community-driven framework for Bash",
			CompatibleShells: []string{"bash"},
			InstallCmd:  installOhMyBash,
		},
		"fisher": {
			Name:        "Fisher",
			Description: "Plugin manager for Fish",
			CompatibleShells: []string{"fish"},
			InstallCmd:  installFisher,
		},
		"oh-my-fish": {
			Name:        "Oh My Fish",
			Description: "Community Fish framework",
			CompatibleShells: []string{"fish"},
			InstallCmd:  installOhMyFish,
		},
	}

	PromptOptions = map[string]Prompt{
		"starship": {
			Name:            "Starship",
			Description:     "Cross-shell customizable prompt",
			CompatibleShells: []string{"zsh", "bash", "fish"},
			BrewPackage:     "starship",
			AptPackage:      "starship",
			DnfPackage:      "starship",
			PacmanPackage:   "starship",
			ZypperPackage:   "starship",
			ChocoPackage:    "starship",
			InstallCmd:      installStarship,
		},
		"powerlevel10k": {
			Name:            "Powerlevel10k",
			Description:     "Fast and customizable Zsh theme",
			CompatibleShells: []string{"zsh"},
			InstallCmd:      installPowerlevel10k,
		},
		"pure": {
			Name:            "Pure",
			Description:     "Pretty, minimal and fast Zsh prompt",
			CompatibleShells: []string{"zsh"},
			InstallCmd:      installPure,
		},
		"oh-my-posh": {
			Name:            "Oh My Posh",
			Description:     "Cross-platform, customizable prompt",
			CompatibleShells: []string{"zsh", "bash", "fish"},
			BrewPackage:     "oh-my-posh",
			AptPackage:      "oh-my-posh",
			DnfPackage:      "oh-my-posh",
			PacmanPackage:   "oh-my-posh",
			ChocoPackage:    "oh-my-posh",
			InstallCmd:      installOhMyPosh,
		},
	}
}

// InstallShell installs the selected shell
func InstallShell(shellName string) error {
	shell, exists := ShellOptions[shellName]
	if !exists {
		return fmt.Errorf("unknown shell: %s", shellName)
	}

	// Check if the shell is already installed
	if isCommandAvailable(shellName) {
		fmt.Printf("✅ %s is already installed\n", shell.Name)
		return nil
	}

	// Try to install with the appropriate package manager
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}

	// Install the shell
	var installErr error
	switch primaryPM {
	case platform.Homebrew:
		installErr = brewInstall(shell.BrewPackage)
	case platform.Apt:
		installErr = aptInstall(shell.AptPackage)
	case platform.Dnf:
		installErr = dnfInstall(shell.DnfPackage)
	case platform.Pacman:
		installErr = pacmanInstall(shell.PacmanPackage)
	case platform.Zypper:
		installErr = zypperInstall(shell.ZypperPackage)
	case platform.Chocolatey:
		installErr = chocoInstall(shell.ChocoPackage)
	default:
		installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
	}

	if installErr != nil {
		return fmt.Errorf("failed to install %s: %w", shell.Name, installErr)
	}

	// Set as default shell if on Unix-like systems
	if runtime.GOOS != "windows" {
		fmt.Printf("🔄 Setting %s as default shell...\n", shell.Name)
		shellPath, err := exec.Command("which", shellName).Output()
		if err != nil {
			fmt.Printf("⚠️ Could not find path to %s: %v\n", shell.Name, err)
			return nil
		}

		shellPathStr := strings.TrimSpace(string(shellPath))
		
		// Check if shell is in /etc/shells
		checkCmd := fmt.Sprintf("grep -q '^%s$' /etc/shells || sudo sh -c 'echo %s >> /etc/shells'", 
			shellPathStr, shellPathStr)
		checkShellCmd := exec.Command("bash", "-c", checkCmd)
		checkShellCmd.Run() // Ignore errors
		
		// Change default shell
		chshCmd := exec.Command("chsh", "-s", shellPathStr)
		chshCmd.Stdout = os.Stdout
		chshCmd.Stderr = os.Stderr
		
		if err := chshCmd.Run(); err != nil {
			fmt.Printf("⚠️ Could not change default shell: %v\n", err)
			fmt.Printf("💡 You can manually change your shell with: chsh -s %s\n", shellPathStr)
		} else {
			fmt.Printf("✅ Default shell changed to %s\n", shell.Name)
		}
	}

	return nil
}

// InstallPluginManager installs the selected plugin manager
func InstallPluginManager(managerName, shellName string) error {
	if managerName == "none" {
		return nil
	}

	manager, exists := PluginManagerOptions[managerName]
	if !exists {
		return fmt.Errorf("unknown plugin manager: %s", managerName)
	}

	// Check if the plugin manager is compatible with the shell
	compatible := false
	for _, shell := range manager.CompatibleShells {
		if shell == shellName {
			compatible = true
			break
		}
	}

	if !compatible {
		return fmt.Errorf("%s is not compatible with %s", manager.Name, shellName)
	}

	// Install the plugin manager
	return manager.InstallCmd(shellName)
}

// InstallPrompt installs the selected prompt
func InstallPrompt(promptName, shellName string) error {
	if promptName == "none" {
		return nil
	}

	prompt, exists := PromptOptions[promptName]
	if !exists {
		return fmt.Errorf("unknown prompt: %s", promptName)
	}

	// Check if the prompt is compatible with the shell
	compatible := false
	for _, shell := range prompt.CompatibleShells {
		if shell == shellName {
			compatible = true
			break
		}
	}

	if !compatible {
		return fmt.Errorf("%s is not compatible with %s", prompt.Name, shellName)
	}

	// If there's a custom install command, use it
	if prompt.InstallCmd != nil {
		return prompt.InstallCmd(shellName)
	}

	// Otherwise, try to install with the package manager
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}

	// Install the prompt
	var installErr error
	switch primaryPM {
	case platform.Homebrew:
		installErr = brewInstall(prompt.BrewPackage)
	case platform.Apt:
		installErr = aptInstall(prompt.AptPackage)
	case platform.Dnf:
		installErr = dnfInstall(prompt.DnfPackage)
	case platform.Pacman:
		installErr = pacmanInstall(prompt.PacmanPackage)
	case platform.Zypper:
		installErr = zypperInstall(prompt.ZypperPackage)
	case platform.Chocolatey:
		installErr = chocoInstall(prompt.ChocoPackage)
	default:
		installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
	}

	return installErr
}

// Plugin manager installation functions

func installZinit(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("zinit only works with zsh")
	}

	// Check if zinit is already installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zinitPath := fmt.Sprintf("%s/.zinit", home)
	if _, err := os.Stat(zinitPath); err == nil {
		fmt.Println("✅ Zinit is already installed")
		return nil
	}

	// Install zinit
	fmt.Println("🔄 Installing Zinit...")
	installCmd := `sh -c "$(curl -fsSL https://git.io/zinit-install)"`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install zinit: %w", err)
	}

	// Add to zshrc if it doesn't exist
	zshrcPath := fmt.Sprintf("%s/.zshrc", home)
	if _, err := os.Stat(zshrcPath); err == nil {
		// Check if zinit is already in zshrc
		data, err := os.ReadFile(zshrcPath)
		if err != nil {
			return fmt.Errorf("failed to read .zshrc: %w", err)
		}

		if !strings.Contains(string(data), "zinit") {
			// Add zinit to zshrc
			zinitConfig := `
### Added by bootstrap-cli
### Zinit configuration
source "$HOME/.zinit/bin/zinit.zsh"
zinit light zsh-users/zsh-autosuggestions
zinit light zsh-users/zsh-syntax-highlighting
zinit light zsh-users/zsh-completions
### End of Zinit configuration
`
			if err := appendToFile(zshrcPath, zinitConfig); err != nil {
				return fmt.Errorf("failed to update .zshrc: %w", err)
			}
		}
	}

	return nil
}

func installOhMyZsh(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("oh-my-zsh only works with zsh")
	}

	// Check if oh-my-zsh is already installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	ohmyzshPath := fmt.Sprintf("%s/.oh-my-zsh", home)
	if _, err := os.Stat(ohmyzshPath); err == nil {
		fmt.Println("✅ Oh My Zsh is already installed")
		return nil
	}

	// Install oh-my-zsh
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

	// Check if bash-it is already installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	bashitPath := fmt.Sprintf("%s/.bash_it", home)
	if _, err := os.Stat(bashitPath); err == nil {
		fmt.Println("✅ Bash-it is already installed")
		return nil
	}

	// Install bash-it
	fmt.Println("🔄 Installing Bash-it...")
	cloneCmd := fmt.Sprintf("git clone --depth=1 https://github.com/Bash-it/bash-it.git %s", bashitPath)
	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone bash-it: %w", err)
	}

	// Run the install script
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

	// Check if oh-my-bash is already installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	ohMyBashPath := fmt.Sprintf("%s/.oh-my-bash", home)
	if _, err := os.Stat(ohMyBashPath); err == nil {
		fmt.Println("✅ Oh My Bash is already installed")
		return nil
	}

	// Install oh-my-bash
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

	// Install fisher
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

	// Install oh-my-fish
	fmt.Println("🔄 Installing Oh My Fish...")
	installCmd := `curl -L https://get.oh-my.fish | fish`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// Prompt installation functions

func installStarship(shell string) error {
	// Check if starship is already installed
	if isCommandAvailable("starship") {
		fmt.Println("✅ Starship is already installed")
		return nil
	}

	// Install starship
	fmt.Println("🔄 Installing Starship...")
	installCmd := `curl -sS https://starship.rs/install.sh | sh -s -- -y`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install starship: %w", err)
	}

	// Add to shell config if it doesn't exist
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configPath string
	var initCmd string

	switch shell {
	case "zsh":
		configPath = fmt.Sprintf("%s/.zshrc", home)
		initCmd = "eval \"$(starship init zsh)\""
	case "bash":
		configPath = fmt.Sprintf("%s/.bashrc", home)
		initCmd = "eval \"$(starship init bash)\""
	case "fish":
		configPath = fmt.Sprintf("%s/.config/fish/config.fish", home)
		initCmd = "starship init fish | source"
		// Ensure config directory exists
		os.MkdirAll(filepath.Dir(configPath), 0755)
	default:
		return fmt.Errorf("unsupported shell for starship: %s", shell)
	}

	// Check if starship init is already in config
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", configPath, err)
		}

		if !strings.Contains(string(data), "starship init") {
			// Add starship init to config
			starshipConfig := fmt.Sprintf("\n### Added by bootstrap-cli\n### Starship prompt\n%s\n### End of Starship configuration\n", initCmd)
			if err := appendToFile(configPath, starshipConfig); err != nil {
				return fmt.Errorf("failed to update %s: %w", configPath, err)
			}
		}
	} else {
		// Create config file
		starshipConfig := fmt.Sprintf("### Added by bootstrap-cli\n### Starship prompt\n%s\n### End of Starship configuration\n", initCmd)
		if err := os.WriteFile(configPath, []byte(starshipConfig), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", configPath, err)
		}
	}

	return nil
}

func installPowerlevel10k(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("powerlevel10k only works with zsh")
	}

	// Check if zinit or oh-my-zsh is installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zinitPath := fmt.Sprintf("%s/.zinit", home)
	ohmyzshPath := fmt.Sprintf("%s/.oh-my-zsh", home)
	
	useZinit := false
	useOhMyZsh := false

	if _, err := os.Stat(zinitPath); err == nil {
		useZinit = true
	} else if _, err := os.Stat(ohmyzshPath); err == nil {
		useOhMyZsh = true
	} else {
		return fmt.Errorf("powerlevel10k requires zinit or oh-my-zsh")
	}

	// Install powerlevel10k
	fmt.Println("🔄 Installing Powerlevel10k...")
	
	// Clone powerlevel10k repository
	p10kPath := ""
	if useZinit {
		// Add to zshrc
		zshrcPath := fmt.Sprintf("%s/.zshrc", home)
		if _, err := os.Stat(zshrcPath); err == nil {
			// Check if p10k is already in zshrc
			data, err := os.ReadFile(zshrcPath)
			if err != nil {
				return fmt.Errorf("failed to read .zshrc: %w", err)
			}

			if !strings.Contains(string(data), "powerlevel10k") {
				// Add p10k to zshrc
				p10kConfig := `
### Added by bootstrap-cli
### Powerlevel10k
zinit ice depth=1; zinit light romkatv/powerlevel10k
### End of Powerlevel10k configuration
`
				if err := appendToFile(zshrcPath, p10kConfig); err != nil {
					return fmt.Errorf("failed to update .zshrc: %w", err)
				}
			}
		}
	} else if useOhMyZsh {
		p10kPath = fmt.Sprintf("%s/custom/themes/powerlevel10k", ohmyzshPath)
		if _, err := os.Stat(p10kPath); err == nil {
			fmt.Println("✅ Powerlevel10k is already installed")
		} else {
			// Clone repository
			cloneCmd := fmt.Sprintf("git clone --depth=1 https://github.com/romkatv/powerlevel10k.git %s", p10kPath)
			cmd := exec.Command("bash", "-c", cloneCmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to clone powerlevel10k: %w", err)
			}

			// Update .zshrc to use powerlevel10k theme
			zshrcPath := fmt.Sprintf("%s/.zshrc", home)
			if _, err := os.Stat(zshrcPath); err == nil {
				// Replace ZSH_THEME line
				updateCmd := fmt.Sprintf(`sed -i 's/ZSH_THEME=".*"/ZSH_THEME="powerlevel10k\/powerlevel10k"/g' %s`, zshrcPath)
				cmd = exec.Command("bash", "-c", updateCmd)
				cmd.Run() // Ignore errors
			}
		}
	}

	// Create p10k.zsh if it doesn't exist
	p10kConfigPath := fmt.Sprintf("%s/.p10k.zsh", home)
	if _, err := os.Stat(p10kConfigPath); err != nil {
		fmt.Println("📝 Creating default p10k.zsh configuration...")
		// Download a sample p10k.zsh
		downloadCmd := fmt.Sprintf("curl -o %s https://raw.githubusercontent.com/romkatv/powerlevel10k/master/config/p10k-lean.zsh", p10kConfigPath)
		cmd := exec.Command("bash", "-c", downloadCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Ignore errors
	}

	return nil
}

func installPure(shell string) error {
	if shell != "zsh" {
		return fmt.Errorf("pure only works with zsh")
	}

	// Check if zinit is installed
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zinitPath := fmt.Sprintf("%s/.zinit", home)
	if _, err := os.Stat(zinitPath); err != nil {
		return fmt.Errorf("pure requires zinit")
	}

	// Add pure to zshrc
	zshrcPath := fmt.Sprintf("%s/.zshrc", home)
	if _, err := os.Stat(zshrcPath); err == nil {
		// Check if pure is already in zshrc
		data, err := os.ReadFile(zshrcPath)
		if err != nil {
			return fmt.Errorf("failed to read .zshrc: %w", err)
		}

		if !strings.Contains(string(data), "pure") {
			// Add pure to zshrc
			pureConfig := `
### Added by bootstrap-cli
### Pure prompt
zinit ice pick"async.zsh" src"pure.zsh"
zinit light sindresorhus/pure
### End of Pure configuration
`
			if err := appendToFile(zshrcPath, pureConfig); err != nil {
				return fmt.Errorf("failed to update .zshrc: %w", err)
			}
		}
	}

	return nil
}

func installOhMyPosh(shell string) error {
	// Install Oh My Posh using the package manager if available
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}

	var installErr error
	prompt, exists := PromptOptions["oh-my-posh"]
	if !exists {
		return fmt.Errorf("oh-my-posh configuration not found")
	}

	switch primaryPM {
	case platform.Homebrew:
		installErr = brewInstall(prompt.BrewPackage)
	case platform.Apt:
		installErr = aptInstall(prompt.AptPackage)
	case platform.Dnf:
		installErr = dnfInstall(prompt.DnfPackage)
	case platform.Pacman:
		installErr = pacmanInstall(prompt.PacmanPackage)
	case platform.Zypper:
		installErr = zypperInstall(prompt.ZypperPackage)
	case platform.Chocolatey:
		installErr = chocoInstall(prompt.ChocoPackage)
	default:
		// Fallback to direct download
		installCmd := `curl -s https://ohmyposh.dev/install.sh | bash -s`
		cmd := exec.Command("bash", "-c", installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		installErr = cmd.Run()
	}

	if installErr != nil {
		return fmt.Errorf("failed to install oh-my-posh: %w", installErr)
	}

	// Add to shell config
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	var configPath string
	var initCmd string

	switch shell {
	case "zsh":
		configPath = fmt.Sprintf("%s/.zshrc", home)
		initCmd = "eval \"$(oh-my-posh init zsh)\""
	case "bash":
		configPath = fmt.Sprintf("%s/.bashrc", home)
		initCmd = "eval \"$(oh-my-posh init bash)\""
	case "fish":
		configPath = fmt.Sprintf("%s/.config/fish/config.fish", home)
		initCmd = "oh-my-posh init fish | source"
		// Ensure config directory exists
		os.MkdirAll(filepath.Dir(configPath), 0755)
	default:
		return fmt.Errorf("unsupported shell for oh-my-posh: %s", shell)
	}

	// Check if oh-my-posh init is already in config
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", configPath, err)
		}

		if !strings.Contains(string(data), "oh-my-posh init") {
			// Add oh-my-posh init to config
			ohMyPoshConfig := fmt.Sprintf("\n### Added by bootstrap-cli\n### Oh My Posh prompt\n%s\n### End of Oh My Posh configuration\n", initCmd)
			if err := appendToFile(configPath, ohMyPoshConfig); err != nil {
				return fmt.Errorf("failed to update %s: %w", configPath, err)
			}
		}
	} else {
		// Create config file
		ohMyPoshConfig := fmt.Sprintf("### Added by bootstrap-cli\n### Oh My Posh prompt\n%s\n### End of Oh My Posh configuration\n", initCmd)
		if err := os.WriteFile(configPath, []byte(ohMyPoshConfig), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", configPath, err)
		}
	}

	// Create a default theme
	themesDir := fmt.Sprintf("%s/.poshthemes", home)
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		os.MkdirAll(themesDir, 0755)
	}

	defaultThemePath := fmt.Sprintf("%s/default.omp.json", themesDir)
	if _, err := os.Stat(defaultThemePath); os.IsNotExist(err) {
		// Download a sample theme
		downloadCmd := fmt.Sprintf("curl -o %s https://raw.githubusercontent.com/JanDeDobbeleer/oh-my-posh/main/themes/jandedobbeleer.omp.json", defaultThemePath)
		cmd := exec.Command("bash", "-c", downloadCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Ignore errors
	}

	return nil
}