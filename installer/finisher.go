package installer

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/prompts"
	"github.com/YitzhakMizrahi/bootstrap-cli/types"
)

// FinishInstallation handles post-installation tasks
func FinishInstallation(config types.UserConfig) error {
	fmt.Println("\n✨ Bootstrap process completed successfully!")
	
	// Add any custom configurations needed for installed tools
	if err := configureInstalledTools(config); err != nil {
		fmt.Printf("⚠️ Warning: Some configurations could not be applied: %v\n", err)
	}
	
	// If a shell was installed or selected, make it the default if needed
	if config.Shell != "" {
		// Check if the selected shell is already the default
		currentShell, err := GetCurrentShell()
		if err == nil && currentShell != config.Shell {
			fmt.Printf("🔄 Would you like to switch to %s now? (y/n): ", config.Shell)
			var response string
			fmt.Scanln(&response)
			response = strings.ToLower(strings.TrimSpace(response))
			
			if response == "y" || response == "yes" {
				// Run a countdown and restart
				return restartWithCountdown(config.Shell)
			}
		}
	}
	
	return nil
}

// configureInstalledTools adds any missing configurations for installed tools
func configureInstalledTools(config types.UserConfig) error {
	var errs []string
	
	// Configure shell integrations
	if config.Shell != "" {
		if err := configureFinalShell(config.Shell); err != nil {
			errs = append(errs, fmt.Sprintf("Shell configuration: %v", err))
		}
	}
	
	// Configure prompt if installed
	if config.Prompt != "" && config.Prompt != "none" {
		if err := configureFinalPrompt(config.Prompt, config.Shell); err != nil {
			errs = append(errs, fmt.Sprintf("Prompt configuration: %v", err))
		}
	}
	
	// Configure CLI tools
	for _, tool := range config.CLITools {
		if err := configureFinalTool(tool, config.Shell); err != nil {
			errs = append(errs, fmt.Sprintf("%s configuration: %v", tool, err))
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("configuration errors: %s", strings.Join(errs, "; "))
	}
	
	return nil
}

// configureFinalShell adds any final shell configurations
func configureFinalShell(shell string) error {
	// Create a bootstrap completion message
	comment := "\n# Added by bootstrap-cli\necho \"✨ Shell environment bootstrapped successfully!\"\n"
	
	// Add any shell-specific configurations
	switch shell {
	case "zsh":
		// Enable basic zsh features
		zshConfig := `
# Added by bootstrap-cli - basic zsh configuration
autoload -Uz compinit
compinit

# History configuration
HISTSIZE=10000
SAVEHIST=10000
HISTFILE=~/.zsh_history

# Basic settings
setopt auto_cd
setopt interactive_comments
setopt extended_glob
setopt appendhistory
setopt hist_ignore_dups
setopt hist_ignore_space
`
		return AddToRcFile(shell, zshConfig+comment)
		
	case "bash":
		// Enable basic bash features
		bashConfig := `
# Added by bootstrap-cli - basic bash configuration
# History configuration
HISTSIZE=10000
HISTFILESIZE=10000

# Basic settings
shopt -s checkwinsize
shopt -s histappend
shopt -s globstar 2>/dev/null
`
		return AddToRcFile(shell, bashConfig+comment)
		
	case "fish":
		// Enable basic fish features
		fishConfig := `
# Added by bootstrap-cli - basic fish configuration
# Set history size
set -g fish_history_max_lines 10000

# Basic settings
set -g fish_prompt_pwd_dir_length 0
`
		return AddToRcFile(shell, fishConfig+comment)
	}
	
	return nil
}

// configureFinalPrompt adds any final prompt configurations
func configureFinalPrompt(prompt, shell string) error {
	// Most prompt configurations are handled in their install functions
	return nil
}

// configureFinalTool adds any final tool configurations
func configureFinalTool(tool, shell string) error {
	switch tool {
	case "fzf":
		// Add fzf integration to shell
		var fzfConfig string
		switch shell {
		case "zsh":
			fzfConfig = `
# Added by bootstrap-cli - fzf configuration
[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh
`
		case "bash":
			fzfConfig = `
# Added by bootstrap-cli - fzf configuration
[ -f ~/.fzf.bash ] && source ~/.fzf.bash
`
		case "fish":
			fzfConfig = `
# Added by bootstrap-cli - fzf configuration
if test -f ~/.fzf.fish
    source ~/.fzf.fish
end
`
		}
		
		if fzfConfig != "" {
			return AddToRcFile(shell, fzfConfig)
		}
		
	case "bat":
		// Add bat alias
		var batConfig string
		switch shell {
		case "zsh", "bash":
			batConfig = `
# Added by bootstrap-cli - bat configuration
alias cat="bat --style=plain --paging=never"
`
		case "fish":
			batConfig = `
# Added by bootstrap-cli - bat configuration
alias cat "bat --style=plain --paging=never"
`
		}
		
		if batConfig != "" {
			return AddToRcFile(shell, batConfig)
		}
		
	case "lsd":
		// Add lsd aliases
		var lsdConfig string
		switch shell {
		case "zsh", "bash":
			lsdConfig = `
# Added by bootstrap-cli - lsd configuration
alias ls="lsd"
alias l="lsd -l"
alias la="lsd -la"
alias lt="lsd --tree"
`
		case "fish":
			lsdConfig = `
# Added by bootstrap-cli - lsd configuration
alias ls "lsd"
alias l "lsd -l"
alias la "lsd -la"
alias lt "lsd --tree"
`
		}
		
		if lsdConfig != "" {
			return AddToRcFile(shell, lsdConfig)
		}
	}
	
	return nil
}

// restartWithCountdown shows a countdown and then restarts the shell
func restartWithCountdown(shellName string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("⚠️ Shell restart not supported on Windows. Please restart your shell manually.")
		return nil
	}
	
	// Run a countdown timer before restarting
	message := fmt.Sprintf("Restarting with %s in", shellName)
	err := prompts.RunTimer(5, message, func() {
		// This function runs after the countdown completes
		fmt.Printf("\n🔄 Restarting shell...\n")
		time.Sleep(500 * time.Millisecond)
	})
	
	if err != nil {
		return fmt.Errorf("timer error: %w", err)
	}
	
	// Execute the new shell after the countdown
	if _, err := exec.LookPath(shellName); err != nil {
		return fmt.Errorf("shell '%s' not found in PATH: %w", shellName, err)
	}
	
	fmt.Printf("🚀 Starting %s...\n", shellName)
	return RestartShell(shellName)
} 