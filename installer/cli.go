package installer

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/YitzhakMizrahi/bootstrap-cli/platform"
)

// Tool represents a CLI tool that can be installed
type Tool struct {
	Name           string
	Description    string
	BrewPackage    string
	AptPackage     string
	DnfPackage     string
	PacmanPackage  string
	ZypperPackage  string
	ChocoPackage   string
	SnapPackage    string
	InstallScript  func() error // Custom installation function
}

// ToolPackages maps tool names to their installation info
var ToolPackages = map[string]Tool{
	"lsd": {
		Name:         "lsd",
		Description:  "The next gen ls command",
		BrewPackage:  "lsd",
		AptPackage:   "lsd",
		DnfPackage:   "lsd",
		PacmanPackage: "lsd",
		ChocoPackage: "lsd",
	},
	"bat": {
		Name:         "bat",
		Description:  "A cat clone with wings",
		BrewPackage:  "bat",
		AptPackage:   "bat",
		DnfPackage:   "bat",
		PacmanPackage: "bat",
		ChocoPackage: "bat",
	},
	"fzf": {
		Name:         "fzf",
		Description:  "Command-line fuzzy finder",
		BrewPackage:  "fzf",
		AptPackage:   "fzf",
		DnfPackage:   "fzf",
		PacmanPackage: "fzf",
		ChocoPackage: "fzf",
	},
	"ripgrep": {
		Name:         "ripgrep",
		Description:  "Fast search tool (rg)",
		BrewPackage:  "ripgrep",
		AptPackage:   "ripgrep",
		DnfPackage:   "ripgrep",
		PacmanPackage: "ripgrep",
		ChocoPackage: "ripgrep",
	},
	"fd": {
		Name:         "fd",
		Description:  "Simple, fast alternative to find",
		BrewPackage:  "fd",
		AptPackage:   "fd-find",
		DnfPackage:   "fd-find",
		PacmanPackage: "fd",
		ChocoPackage: "fd",
	},
	"jq": {
		Name:         "jq",
		Description:  "Lightweight JSON processor",
		BrewPackage:  "jq",
		AptPackage:   "jq",
		DnfPackage:   "jq",
		PacmanPackage: "jq",
		ChocoPackage: "jq",
	},
	"htop": {
		Name:         "htop",
		Description:  "Interactive process viewer",
		BrewPackage:  "htop",
		AptPackage:   "htop",
		DnfPackage:   "htop",
		PacmanPackage: "htop",
		ChocoPackage: "htop",
	},
	"lazygit": {
		Name:         "lazygit",
		Description:  "Simple terminal UI for git",
		BrewPackage:  "lazygit",
		AptPackage:   "lazygit",
		DnfPackage:   "lazygit",
		PacmanPackage: "lazygit",
		ChocoPackage: "lazygit",
	},
	"tmux": {
		Name:         "tmux",
		Description:  "Terminal multiplexer",
		BrewPackage:  "tmux",
		AptPackage:   "tmux",
		DnfPackage:   "tmux",
		PacmanPackage: "tmux",
		ChocoPackage: "tmux",
	},
	"neofetch": {
		Name:         "neofetch",
		Description:  "System info written in bash",
		BrewPackage:  "neofetch",
		AptPackage:   "neofetch",
		DnfPackage:   "neofetch",
		PacmanPackage: "neofetch",
		ChocoPackage: "neofetch",
	},
}

// Language represents a programming language that can be installed
type Language struct {
	Name         string
	Description  string
	InstallCmd   func() error
	Managers     []string // Available package managers for this language
}

// LanguagePackages maps language names to their installation info
var LanguagePackages = map[string]Language{
	"node": {
		Name:        "Node.js",
		Description: "JavaScript runtime",
		InstallCmd:  installNode,
		Managers:    []string{"npm", "yarn", "pnpm"},
	},
	"python": {
		Name:        "Python",
		Description: "Python programming language",
		InstallCmd:  installPython,
		Managers:    []string{"pip", "pipenv", "poetry"},
	},
	"go": {
		Name:        "Go",
		Description: "Go programming language",
		InstallCmd:  installGo,
		Managers:    []string{},
	},
	"rust": {
		Name:        "Rust",
		Description: "Rust programming language",
		InstallCmd:  installRust,
		Managers:    []string{"cargo"},
	},
}

// InstallCLITools installs the selected CLI tools
func InstallCLITools(tools []string) error {
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}
	
	for _, toolName := range tools {
		tool, exists := ToolPackages[toolName]
		if !exists {
			fmt.Printf("⚠️  Unknown tool: %s, skipping...\n", toolName)
			continue
		}
		
		fmt.Printf("📦 Installing %s...\n", tool.Name)
		
		// Check if already installed
		if isCommandAvailable(toolName) {
			fmt.Printf("✅ %s is already installed, skipping\n", tool.Name)
			continue
		}
		
		// Try to install with the appropriate package manager
		var installErr error
		switch primaryPM {
		case platform.Homebrew:
			installErr = brewInstall(tool.BrewPackage)
		case platform.Apt:
			installErr = aptInstall(tool.AptPackage)
		case platform.Dnf:
			installErr = dnfInstall(tool.DnfPackage)
		case platform.Pacman:
			installErr = pacmanInstall(tool.PacmanPackage)
		case platform.Zypper:
			installErr = zypperInstall(tool.ZypperPackage)
		case platform.Chocolatey:
			installErr = chocoInstall(tool.ChocoPackage)
		default:
			installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
		}
		
		if installErr != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", tool.Name, installErr)
			// Try alternative installation methods if available
			if tool.InstallScript != nil {
				fmt.Printf("🔄 Trying alternative installation method for %s...\n", tool.Name)
				if err := tool.InstallScript(); err != nil {
					fmt.Printf("❌ Alternative installation failed: %v\n", err)
				} else {
					fmt.Printf("✅ Installed %s using alternative method\n", tool.Name)
				}
			}
		} else {
			fmt.Printf("✅ Successfully installed %s\n", tool.Name)
		}
	}
	
	return nil
}

// InstallLanguages installs the selected programming languages
func InstallLanguages(languages []string, packageManagers map[string]string) error {
	for _, lang := range languages {
		language, exists := LanguagePackages[lang]
		if !exists {
			fmt.Printf("⚠️  Unknown language: %s, skipping...\n", lang)
			continue
		}
		
		fmt.Printf("🧪 Installing %s...\n", language.Name)
		
		if err := language.InstallCmd(); err != nil {
			fmt.Printf("❌ Failed to install %s: %v\n", language.Name, err)
		} else {
			fmt.Printf("✅ Successfully installed %s\n", language.Name)
			
			// Install selected package manager if applicable
			if mgr, ok := packageManagers[lang]; ok && mgr != "" {
				fmt.Printf("📦 Setting up %s package manager: %s\n", language.Name, mgr)
				if err := installPackageManager(lang, mgr); err != nil {
					fmt.Printf("❌ Failed to set up package manager %s: %v\n", mgr, err)
				}
			}
		}
	}
	
	return nil
}

// Language installations

func installNode() error {
	// Check if nvm exists
	nvmPath := fmt.Sprintf("%s/.nvm", os.Getenv("HOME"))
	if _, err := os.Stat(nvmPath); os.IsNotExist(err) {
		fmt.Println("🔄 Installing NVM (Node Version Manager)...")
		// Install nvm
		installCmd := `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash`
		cmd := exec.Command("bash", "-c", installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install nvm: %w", err)
		}
	}
	
	// Install latest LTS version of Node
	fmt.Println("🔄 Installing Node.js LTS version...")
	cmd := exec.Command("bash", "-c", "source ~/.nvm/nvm.sh && nvm install --lts")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installPython() error {
	// Check if pyenv exists
	pyenvPath := fmt.Sprintf("%s/.pyenv", os.Getenv("HOME"))
	if _, err := os.Stat(pyenvPath); os.IsNotExist(err) {
		fmt.Println("🔄 Installing pyenv...")
		// Install pyenv
		installCmd := `curl https://pyenv.run | bash`
		cmd := exec.Command("bash", "-c", installCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install pyenv: %w", err)
		}
		
		// Add pyenv to shell config
		fmt.Println("🔄 Adding pyenv to shell configuration...")
		shellConfig := fmt.Sprintf("%s/.bashrc", os.Getenv("HOME"))
		if _, err := os.Stat(shellConfig); err == nil {
			appendToFile(shellConfig, "\n# pyenv\nexport PATH=\"$HOME/.pyenv/bin:$PATH\"\neval \"$(pyenv init -)\"\neval \"$(pyenv virtualenv-init -)\"\n")
		}
	}
	
	// Install latest Python
	fmt.Println("🔄 Installing latest Python version...")
	cmd := exec.Command("bash", "-c", "source ~/.pyenv/bin/activate && pyenv install 3.11 && pyenv global 3.11")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installGo() error {
	platformInfo, err := platform.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	primaryPM, err := platform.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}
	
	// Install Go using package manager
	switch primaryPM {
	case platform.Homebrew:
		return brewInstall("go")
	case platform.Apt:
		return aptInstall("golang-go")
	case platform.Dnf:
		return dnfInstall("golang")
	case platform.Pacman:
		return pacmanInstall("go")
	case platform.Zypper:
		return zypperInstall("go")
	case platform.Chocolatey:
		return chocoInstall("golang")
	default:
		// Fallback to manual install for other platforms
		return fmt.Errorf("manual Go installation not implemented for this platform")
	}
}

func installRust() error {
	fmt.Println("🔄 Installing Rust via rustup...")
	installCmd := `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y`
	cmd := exec.Command("bash", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Package manager installations

func installPackageManager(language, manager string) error {
	switch language {
	case "node":
		switch manager {
		case "npm":
			// npm comes with Node.js
			return nil
		case "yarn":
			cmd := exec.Command("bash", "-c", "source ~/.nvm/nvm.sh && npm install -g yarn")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		case "pnpm":
			cmd := exec.Command("bash", "-c", "source ~/.nvm/nvm.sh && npm install -g pnpm")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	case "python":
		switch manager {
		case "pip":
			// pip comes with Python
			return nil
		case "pipenv":
			cmd := exec.Command("bash", "-c", "source ~/.pyenv/bin/activate && pip install pipenv")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		case "poetry":
			cmd := exec.Command("bash", "-c", "source ~/.pyenv/bin/activate && pip install poetry")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}
	
	return fmt.Errorf("unknown package manager %s for language %s", manager, language)
}