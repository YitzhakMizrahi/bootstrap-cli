// prompts/init_prompts.go
package prompts

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/YitzhakMizrahi/bootstrap-cli/types"
)

// RunInitWizard handles the interactive flow for bootstrap init
func RunInitWizard() types.UserConfig {
	cfg := types.UserConfig{}

	// Shell
	survey.AskOne(&survey.Select{
		Message: "Which shell do you use?",
		Options: []string{"zsh", "bash", "fish"},
		Default: "zsh",
	}, &cfg.Shell)

	// Plugin Manager
	survey.AskOne(&survey.Select{
		Message: "Which plugin manager do you prefer?",
		Options: []string{"zinit", "none"},
		Default: "zinit",
	}, &cfg.PluginManager)

	// Prompt
	survey.AskOne(&survey.Select{
		Message: "Which prompt do you want to use?",
		Options: []string{"starship", "pure", "none"},
		Default: "starship",
	}, &cfg.Prompt)

	// Automatically include starship if selected as prompt
	if cfg.Prompt == "starship" {
		cfg.CLITools = append(cfg.CLITools, "starship")
	}

	// CLI Tools (excluding prompt tools)
	survey.AskOne(&survey.MultiSelect{
		Message: "Select CLI tools to install:",
		Options: []string{"lsd", "fzf", "bat", "zoxide", "eza", "curlie", "tmux", "lazygit", "yazi"},
		Default: []string{"lsd"},
	}, &cfg.CLITools)

	// Languages
	survey.AskOne(&survey.MultiSelect{
		Message: "Select languages to install:",
		Options: []string{"python", "node", "go", "rust"},
		Default: []string{"python", "node"},
	}, &cfg.Languages)

	// Package managers
	cfg.PackageManagers = map[string]string{}
	for _, lang := range cfg.Languages {
		switch lang {
		case "node":
			var nodePkg string
			survey.AskOne(&survey.Select{
				Message: "Preferred package manager for Node.js:",
				Options: []string{"pnpm", "npm", "yarn"},
				Default: "pnpm",
			}, &nodePkg)
			cfg.PackageManagers["node"] = nodePkg
		case "python":
			cfg.PackageManagers["python"] = "pip"
		case "go":
			cfg.PackageManagers["go"] = "go"
		case "rust":
			cfg.PackageManagers["rust"] = "cargo"
		}
	}

	// Editors
	survey.AskOne(&survey.MultiSelect{
		Message: "Which editors do you want to set up?",
		Options: []string{"vim", "nvim", "nvim (LazyVim)", "nvim (AstroNvim)", "none"},
		Default: []string{"nvim"},
	}, &cfg.Editors)

	// Dotfiles path
	survey.AskOne(&survey.Input{
		Message: "Where is your dotfiles repo?",
		Default: "~/.dotfiles",
	}, &cfg.DotfilesPath)

	// Symlink mode
	survey.AskOne(&survey.Confirm{
		Message: "Use relative symlinks?",
		Default: false,
	}, &cfg.UseRelativeLinks)

	// Dev mode
	survey.AskOne(&survey.Confirm{
		Message: "Enable developer mode (verbose output)?",
		Default: true,
	}, &cfg.DevMode)

	// Backup mode
	survey.AskOne(&survey.Confirm{
		Message: "Backup existing config files before linking?",
		Default: true,
	}, &cfg.BackupExisting)

	fmt.Println("✨ Setup complete. Saving your config...")
	return cfg
}
