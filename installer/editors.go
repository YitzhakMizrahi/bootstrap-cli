package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/pkg/platform"
)

// Editor represents a code editor that can be installed
type Editor struct {
	Name          string
	Description   string
	BrewPackage   string
	AptPackage    string
	DnfPackage    string
	PacmanPackage string
	ChocoPackage  string
	InstallFunc   func() error
}

// EditorOptions maps editor names to their installation info
var EditorOptions = map[string]Editor{
	"neovim": {
		Name:          "NeoVim",
		Description:   "Hyperextensible Vim-based editor",
		BrewPackage:   "neovim",
		AptPackage:    "neovim",
		DnfPackage:    "neovim",
		PacmanPackage: "neovim",
		ChocoPackage:  "neovim",
	},
	"lazyvim": {
		Name:        "LazyVim",
		Description: "Neovim with LazyVim configuration",
		InstallFunc: installLazyVim,
	},
	"astronvim": {
		Name:        "AstroNvim",
		Description: "Neovim with AstroNvim configuration",
		InstallFunc: installAstroNvim,
	},
	"vscode": {
		Name:          "VS Code",
		Description:   "Visual Studio Code",
		BrewPackage:   "visual-studio-code",
		AptPackage:    "code",
		DnfPackage:    "code",
		PacmanPackage: "visual-studio-code-bin",
		ChocoPackage:  "vscode",
	},
}

// InstallEditors installs the selected editors
func InstallEditors(editors []string) error {
	detector := platform.NewDetector()
	platformInfo, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}
	
	primaryPM, err := detector.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}
	
	for _, editorName := range editors {
		editor, exists := EditorOptions[editorName]
		if !exists {
			fmt.Printf("⚠️  Unknown editor: %s, skipping...\n", editorName)
			continue
		}
		
		fmt.Printf("📝 Installing %s...\n", editor.Name)
		
		// Check if custom install function exists
		if editor.InstallFunc != nil {
			if err := editor.InstallFunc(); err != nil {
				fmt.Printf("❌ Failed to install %s: %v\n", editor.Name, err)
			} else {
				fmt.Printf("✅ Successfully installed %s\n", editor.Name)
			}
			continue
		}
		
		// Try to install with the appropriate package manager
		var installErr error
		switch primaryPM {
		case platform.Homebrew:
			installErr = brewInstall(editor.BrewPackage)
		case platform.Apt:
			installErr = aptInstall(editor.AptPackage)
		case platform.Dnf:
			installErr = dnfInstall(editor.DnfPackage)
		case platform.Pacman:
			installErr = pacmanInstall(editor.PacmanPackage)
		default:
			installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
		}
		
		if installErr != nil {
			fmt.Printf("⚠️  Failed to install %s: %v\n", editor.Name, installErr)
		} else {
			fmt.Printf("✅ Successfully installed %s\n", editor.Name)
		}
	}
	
	return nil
}

// installNeovim installs Neovim 
func installNeovim() error {
	detector := platform.NewDetector()
	platformInfo, err := detector.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect platform: %w", err)
	}

	primaryPM, err := detector.GetPrimaryPackageManager(platformInfo)
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}

	// Install neovim directly
	var installErr error
	switch primaryPM {
	case platform.Homebrew:
		installErr = brewInstall("neovim")
	case platform.Apt:
		installErr = aptInstall("neovim")
	case platform.Dnf:
		installErr = dnfInstall("neovim")
	case platform.Pacman:
		installErr = pacmanInstall("neovim")
	default:
		installErr = fmt.Errorf("unsupported package manager: %s", primaryPM)
	}

	return installErr
}

// installLazyVim installs LazyVim
func installLazyVim() error {
	// First ensure neovim is installed
	if !isCommandAvailable("nvim") {
		fmt.Println("🔄 LazyVim requires Neovim, installing...")
		if err := installNeovim(); err != nil {
			return fmt.Errorf("failed to install neovim: %w", err)
		}
	}
	
	// Backup existing config
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	nvimConfigDir := filepath.Join(home, ".config/nvim")
	if _, err := os.Stat(nvimConfigDir); err == nil {
		backupDir := fmt.Sprintf("%s.bak", nvimConfigDir)
		fmt.Printf("📦 Backing up existing Neovim config to %s\n", backupDir)
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf("failed to remove old backup: %w", err)
		}
		if err := os.Rename(nvimConfigDir, backupDir); err != nil {
			return fmt.Errorf("failed to backup existing config: %w", err)
		}
	}
	
	// Clone LazyVim
	fmt.Println("🔄 Installing LazyVim...")
	cloneCmd := fmt.Sprintf("git clone https://github.com/LazyVim/starter %s", nvimConfigDir)
	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone LazyVim: %w", err)
	}
	
	// Remove .git directory to avoid conflicts
	gitDir := filepath.Join(nvimConfigDir, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		fmt.Printf("⚠️  Failed to remove .git directory: %v\n", err)
	}
	
	fmt.Println("✅ LazyVim installed successfully!")
	return nil
}

// installAstroNvim installs AstroNvim
func installAstroNvim() error {
	// First ensure neovim is installed
	if !isCommandAvailable("nvim") {
		fmt.Println("🔄 AstroNvim requires Neovim, installing...")
		if err := installNeovim(); err != nil {
			return fmt.Errorf("failed to install neovim: %w", err)
		}
	}
	
	// Backup existing config
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	nvimConfigDir := filepath.Join(home, ".config/nvim")
	if _, err := os.Stat(nvimConfigDir); err == nil {
		backupDir := fmt.Sprintf("%s.bak", nvimConfigDir)
		fmt.Printf("📦 Backing up existing Neovim config to %s\n", backupDir)
		if err := os.RemoveAll(backupDir); err != nil {
			return fmt.Errorf("failed to remove old backup: %w", err)
		}
		if err := os.Rename(nvimConfigDir, backupDir); err != nil {
			return fmt.Errorf("failed to backup existing config: %w", err)
		}
	}
	
	// Clone AstroNvim
	fmt.Println("🔄 Installing AstroNvim...")
	cloneCmd := fmt.Sprintf("git clone https://github.com/AstroNvim/AstroNvim %s", nvimConfigDir)
	cmd := exec.Command("bash", "-c", cloneCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone AstroNvim: %w", err)
	}
	
	fmt.Println("✅ AstroNvim installed successfully!")
	return nil
}