package installer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/types"
)

// FinishInstallation handles final tasks after all installations are complete
func FinishInstallation(config types.UserConfig) error {
	fmt.Println("\n🎉 Successfully completed installation!")
	
	// Configure installed tools with appropriate shell integrations
	if err := configureInstalledTools(config); err != nil {
		fmt.Printf("⚠️ Some configurations could not be applied: %v\n", err)
	}
	
	// Prompt the user to switch to the selected shell if it's not already the default
	if config.Shell != "" {
		if err := configureFinalShell(config.Shell); err != nil {
			return fmt.Errorf("failed to configure shell: %v", err)
		}
		
		// Check if we need to switch to the newly installed shell
		if err := configureFinalShell(config.Shell); err != nil {
			fmt.Printf("⚠️ Failed to configure shell: %v\n", err)
		} else {
			currentShell, err := GetCurrentShell()
			if err != nil {
				fmt.Printf("⚠️ Could not determine current shell: %v\n", err)
			} else if currentShell != config.Shell {
				// Prompt user to restart shell
				fmt.Printf("\n🔄 Your selected shell %s is now configured, but you're currently using %s.\n", config.Shell, currentShell)
				fmt.Printf("Would you like to switch to %s now? (Y/n): ", config.Shell)
				
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))
				
				if input == "" || input == "y" || input == "yes" {
					fmt.Printf("\n🔄 Restarting with %s in %d seconds...\n", config.Shell, 5)
					if err := restartWithCountdown(config.Shell); err != nil {
						return fmt.Errorf("failed to restart shell: %v", err)
					}
				} else {
					fmt.Printf("\nℹ️ You can manually switch to %s with the following command:\n", config.Shell)
					fmt.Printf("   chsh -s $(which %s)\n\n", config.Shell)
				}
			} else {
				fmt.Printf("\n✅ You're already using %s shell.\n", config.Shell)
			}
		}
	}
	
	// If user has dotfiles configured, inform them about the integration
	if config.DotfilesPath != "" {
		fmt.Println("\n📂 Dotfiles Configuration:")
		fmt.Printf("    Your dotfiles are configured from: %s\n", config.DotfilesPath)
		fmt.Println("    Any custom tool configurations have been integrated with your dotfiles.")
		fmt.Println("    You can customize your setup further by editing your dotfiles.")
	} else {
		fmt.Println("\n📝 No dotfiles path was specified. To set up dotfiles in the future, run:")
		fmt.Println("    bootstrap-cli init --force")
	}
	
	fmt.Println("\n🚀 Your development environment is now ready to use!")
	fmt.Println("    Run 'bootstrap-cli --help' to see all available commands.")
	
	return nil
}

// configureInstalledTools adds configuration for tools that were installed
func configureInstalledTools(config types.UserConfig) error {
	// Configure each installed CLI tool
	for _, tool := range config.CLITools {
		fmt.Printf("🔧 Configuring %s...\n", tool)
		if err := configureFinalTool(tool, config.Shell); err != nil {
			fmt.Printf("⚠️ Failed to configure %s: %v\n", tool, err)
		}
	}
	
	// Add shell configuration
	if err := configureFinalShell(config.Shell); err != nil {
		fmt.Printf("⚠️ Failed to configure shell: %v\n", err)
	}
	
	// Add prompt configuration
	if config.Prompt != "" && config.Prompt != "none" {
		if err := configureFinalPrompt(config.Prompt, config.Shell); err != nil {
			fmt.Printf("⚠️ Failed to configure prompt: %v\n", err)
		}
	}
	
	// Check and integrate with user dotfiles if available
	if err := configureUserDotfiles(config); err != nil {
		fmt.Printf("⚠️ Could not integrate with dotfiles: %v\n", err)
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
# ------------------------------------------------

# History configuration
HISTSIZE=10000
SAVEHIST=10000
HISTFILE=~/.zsh_history
setopt appendhistory
setopt hist_ignore_dups
setopt hist_ignore_space
setopt hist_verify

# Basic settings
setopt auto_cd
setopt interactive_comments
setopt extended_glob
setopt promptsubst
setopt nocaseglob
setopt nocheckjobs
setopt numericglobsort

# Autocompletion
autoload -Uz compinit
compinit

zstyle ':completion:*' menu select
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}' # Case insensitive tab completion
zstyle ':completion:*' list-colors "${(s.:.)LS_COLORS}"
zstyle ':completion:*' rehash true
zstyle ':completion:*' accept-exact '*(N)'
zstyle ':completion:*' use-cache on
zstyle ':completion:*' cache-path ~/.zsh/cache

# Key bindings
bindkey -e                              # Use emacs key bindings
bindkey '^[[1;5C' forward-word          # Ctrl+right - forward one word
bindkey '^[[1;5D' backward-word         # Ctrl+left - backward one word
bindkey '^[[3~' delete-char             # Delete key
bindkey '^[[A' history-search-backward  # Up Arrow - search history backwards
bindkey '^[[B' history-search-forward   # Down Arrow - search history forwards

# Environment path setup for common tools
# ---------------------------------------
# Add common directories to PATH if they exist and aren't already in PATH
path_dirs=(
  "$HOME/.local/bin"
  "$HOME/bin"
  "/usr/local/bin"
)

for dir in "${path_dirs[@]}"; do
  if [ -d "$dir" ] && [[ ":$PATH:" != *":$dir:"* ]]; then
    export PATH="$dir:$PATH"
  fi
done

# Cargo (Rust)
if [ -d "$HOME/.cargo/bin" ] && [[ ":$PATH:" != *":$HOME/.cargo/bin:"* ]]; then
  export PATH="$HOME/.cargo/bin:$PATH"
fi

# Go
if [ -d "$HOME/go/bin" ] && [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
  export PATH="$HOME/go/bin:$PATH"
fi

# Node.js via NVM
if [ -d "$HOME/.nvm" ]; then
  export NVM_DIR="$HOME/.nvm"
  [ -s "$NVM_DIR/nvm.sh" ] && source "$NVM_DIR/nvm.sh"
  [ -s "$NVM_DIR/bash_completion" ] && source "$NVM_DIR/bash_completion"
fi

# Python - consider local user packages
if [ -d "$HOME/.local/bin" ] && [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
  export PATH="$HOME/.local/bin:$PATH"
fi

# Colored output for ls (most systems)
export CLICOLOR=1

# If dircolors is available, use it for colored ls output
if command -v dircolors >/dev/null 2>&1; then
  eval "$(dircolors -b)"
  alias ls='ls --color=auto'
fi

# Common aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias h='history'
alias grep='grep --color=auto'
alias fgrep='fgrep --color=auto'
alias egrep='egrep --color=auto'
alias diff='diff --color=auto'
alias ip='ip --color=auto'
alias c='clear'

# Navigation shortcuts
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias .....='cd ../../../..'
alias ~='cd ~'

# Directory shortcuts
alias md='mkdir -p'
alias rd='rmdir'

# General productivity
alias path='echo -e ${PATH//:/\\n}'
alias now='date +"%T"'
alias nowdate='date +"%d-%m-%Y"'

# Enhanced system commands
alias df='df -h'
alias du='du -h'
alias free='free -m'

# Function to make a directory and change into it
function mkcd {
  mkdir -p "$1" && cd "$1"
}
`
		return configureRcFile(shell, zshConfig+comment)
		
	case "bash":
		// Enable basic bash features
		bashConfig := `
# Added by bootstrap-cli - basic bash configuration
# -------------------------------------------------

# History configuration
HISTSIZE=10000
HISTFILESIZE=10000
HISTCONTROL=ignoreboth
shopt -s histappend

# Basic settings
shopt -s checkwinsize
shopt -s globstar 2>/dev/null   # Pattern ** for recursive matches
shopt -s autocd 2>/dev/null     # Change to named directory
shopt -s cdspell 2>/dev/null    # Autocorrect typos in path names
shopt -s dirspell 2>/dev/null   # Autocorrect typos in path names during completion
shopt -s nocaseglob 2>/dev/null # Case-insensitive globbing

# Environment path setup for common tools
# ---------------------------------------
# Add common directories to PATH if they exist and aren't already in PATH
path_dirs=(
  "$HOME/.local/bin"
  "$HOME/bin"
  "/usr/local/bin"
)

for dir in "${path_dirs[@]}"; do
  if [ -d "$dir" ] && [[ ":$PATH:" != *":$dir:"* ]]; then
    export PATH="$dir:$PATH"
  fi
done

# Cargo (Rust)
if [ -d "$HOME/.cargo/bin" ] && [[ ":$PATH:" != *":$HOME/.cargo/bin:"* ]]; then
  export PATH="$HOME/.cargo/bin:$PATH"
fi

# Go
if [ -d "$HOME/go/bin" ] && [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
  export PATH="$HOME/go/bin:$PATH"
fi

# Node.js via NVM
if [ -d "$HOME/.nvm" ]; then
  export NVM_DIR="$HOME/.nvm"
  [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
  [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
fi

# Python - consider local user packages
if [ -d "$HOME/.local/bin" ] && [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
  export PATH="$HOME/.local/bin:$PATH"
fi

# Colored output for ls (most systems)
export CLICOLOR=1

# If dircolors is available, use it for colored ls output
if command -v dircolors >/dev/null 2>&1; then
  eval "$(dircolors -b)"
  alias ls='ls --color=auto'
fi

# Common aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias h='history'
alias grep='grep --color=auto'
alias fgrep='fgrep --color=auto'
alias egrep='egrep --color=auto'
alias diff='diff --color=auto'
alias ip='ip --color=auto'
alias c='clear'

# Navigation shortcuts
alias ..='cd ..'
alias ...='cd ../..'
alias ....='cd ../../..'
alias .....='cd ../../../..'
alias ~='cd ~'

# Directory shortcuts
alias md='mkdir -p'
alias rd='rmdir'

# General productivity
alias path='echo -e ${PATH//:/\\n}'
alias now='date +"%T"'
alias nowdate='date +"%d-%m-%Y"'

# Enhanced system commands
alias df='df -h'
alias du='du -h'
alias free='free -m'

# Function to make a directory and change into it
function mkcd {
  mkdir -p "$1" && cd "$1"
}

# Enable programmable completion features
if ! shopt -oq posix; then
  if [ -f /usr/share/bash-completion/bash_completion ]; then
    . /usr/share/bash-completion/bash_completion
  elif [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
  fi
fi
`
		return configureRcFile(shell, bashConfig+comment)
		
	case "fish":
		// Enable basic fish features
		fishConfig := `
# Added by bootstrap-cli - basic fish configuration
# -------------------------------------------------

# Set history size
set -g fish_history_max_lines 10000

# Basic settings
set -g fish_prompt_pwd_dir_length 0
set -U fish_greeting ""

# Environment path setup for common tools
# ---------------------------------------
# Add common directories to PATH if they exist and aren't already in PATH
set path_dirs $HOME/.local/bin $HOME/bin /usr/local/bin

for dir in $path_dirs
  if test -d $dir; and not contains $dir $PATH
    set -gx PATH $dir $PATH
  end
end

# Cargo (Rust)
if test -d "$HOME/.cargo/bin"; and not contains "$HOME/.cargo/bin" $PATH
  set -gx PATH "$HOME/.cargo/bin" $PATH
end

# Go
if test -d "$HOME/go/bin"; and not contains "$HOME/go/bin" $PATH
  set -gx PATH "$HOME/go/bin" $PATH
end

# Python - consider local user packages
if test -d "$HOME/.local/bin"; and not contains "$HOME/.local/bin" $PATH
  set -gx PATH "$HOME/.local/bin" $PATH
end

# Colored output for ls (most systems)
if command -v dircolors >/dev/null 2>&1
  eval (dircolors -c)
  alias ls='ls --color=auto'
end

# Common aliases
alias ll 'ls -alF'
alias la 'ls -A'
alias l 'ls -CF'
alias h 'history'
alias grep 'grep --color=auto'
alias fgrep 'fgrep --color=auto'
alias egrep 'egrep --color=auto'
alias diff 'diff --color=auto'
alias c 'clear'

# Navigation shortcuts
alias .. 'cd ..'
alias ... 'cd ../..'
alias .... 'cd ../../..'
alias ..... 'cd ../../../..'
alias ~ 'cd ~'

# Directory shortcuts
alias md 'mkdir -p'
alias rd 'rmdir'

# General productivity
alias path 'echo $PATH | tr " " "\\n"'
alias now 'date +"%T"'
alias nowdate 'date +"%d-%m-%Y"'

# Enhanced system commands
alias df 'df -h'
alias du 'du -h'
alias free 'free -m'

# Node.js via NVM (using bass or compatible plugin if available)
if test -d "$HOME/.nvm"
  set -gx NVM_DIR "$HOME/.nvm"
  
  # Check if bass is installed for running bash scripts
  if type -q bass
    bass source "$NVM_DIR/nvm.sh"
  end
end

# Function to make a directory and change into it
function mkcd
  mkdir -p $argv[1]
  cd $argv[1]
end

# Enable fish_config to configure colors, prompts, and completions
if command -v fish_config >/dev/null 2>&1
  # Add fish_config command hint
  echo "Run 'fish_config' to customize your fish shell"
end
`
		return configureRcFile(shell, fishConfig+comment)
	}
	
	return nil
}

// configureFinalPrompt adds any final prompt configurations
func configureFinalPrompt(prompt, shell string) error {
	// For now all prompt-specific configs are done in their installers
	return nil
}

// configureFinalTool adds any final tool configurations
func configureFinalTool(tool, shell string) error {
	// Check if the tool exists before adding its configuration
	if _, err := exec.LookPath(tool); err != nil {
		// Tool not found, skip configuration
		return nil
	}

	// Add tool-specific configurations
	switch tool {
	case "fzf":
		return configureFzf(shell)
	case "bat":
		return configureBat(shell)
	case "lsd":
		return configureLsd(shell)
	case "ripgrep":
		return configureRipgrep(shell)
	case "fd":
		return configureFd(shell)
	case "lazygit":
		return configureLazyGit(shell)
	}
	
	return nil
}

// configureRcFile adds configuration to shell RC files if not already present
func configureRcFile(shell, content string) error {
	var rcFile string
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	switch shell {
	case "zsh":
		rcFile = filepath.Join(home, ".zshrc")
	case "bash":
		rcFile = filepath.Join(home, ".bashrc")
	case "fish":
		rcFile = filepath.Join(home, ".config/fish/config.fish")
		// Create fish config dir if it doesn't exist
		os.MkdirAll(filepath.Join(home, ".config/fish"), 0755)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	
	// Create the file if it doesn't exist
	if _, err := os.Stat(rcFile); os.IsNotExist(err) {
		if err := os.WriteFile(rcFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		return nil
	}

	// If file exists, check if configuration already exists before adding
	existingConfig, err := os.ReadFile(rcFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Check if the configuration or significant parts of it already exist
	if containsConfiguration(string(existingConfig), content) {
		// Configuration already exists, no need to add it again
		return nil
	}

	// Append the new configuration
	f, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()
	
	if _, err := f.WriteString("\n" + content); err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}
	
	return nil
}

// containsConfiguration checks if existing config already contains the important parts of new config
func containsConfiguration(existing, new string) bool {
	// Split both strings into lines for easier comparison
	existingLines := strings.Split(existing, "\n")
	newLines := strings.Split(new, "\n")

	// Extract meaningful lines from new config (skip comments and empty lines)
	var importantNewLines []string
	for _, line := range newLines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		importantNewLines = append(importantNewLines, line)
	}

	// Check if important config elements are already present
	presentCount := 0
	for _, newLine := range importantNewLines {
		for _, existingLine := range existingLines {
			existingLine = strings.TrimSpace(existingLine)
			// Check for exact match or if this is a variable/alias/function definition
			// and the key part (before = or space) matches
			if existingLine == newLine || 
			   (strings.Contains(newLine, "=") && strings.HasPrefix(existingLine, strings.Split(newLine, "=")[0])) ||
			   (strings.Contains(newLine, " ") && strings.HasPrefix(existingLine, strings.Split(newLine, " ")[0])) {
				presentCount++
				break
			}
		}
	}

	// If more than 70% of important lines are already present, consider it existing
	threshold := float64(len(importantNewLines)) * 0.7
	return float64(presentCount) >= threshold
}

// Tool-specific configuration functions

func configureFzf(shell string) error {
	// Base FZF configuration that works for all shells
	switch shell {
	case "zsh", "bash":
		fzfConfig := `
# Added by bootstrap-cli - fzf configuration
if [ -f ~/.fzf.${SHELL##*/}rc ]; then
  source ~/.fzf.${SHELL##*/}rc
elif command -v fzf >/dev/null 2>&1; then
  # FZF keybindings and completion
  if [ -d /usr/share/fzf ]; then
    # For Linux systems
    [ -f /usr/share/fzf/key-bindings.${SHELL##*/} ] && source /usr/share/fzf/key-bindings.${SHELL##*/}
    [ -f /usr/share/fzf/completion.${SHELL##*/} ] && source /usr/share/fzf/completion.${SHELL##*/}
  elif [ -d /usr/local/share/fzf ]; then
    # For macOS systems
    [ -f /usr/local/share/fzf/key-bindings.${SHELL##*/} ] && source /usr/local/share/fzf/key-bindings.${SHELL##*/}
    [ -f /usr/local/share/fzf/completion.${SHELL##*/} ] && source /usr/local/share/fzf/completion.${SHELL##*/}
  elif [ -d "${HOME}/.fzf" ]; then
    # For user-installed fzf
    [ -f ~/.fzf/shell/key-bindings.${SHELL##*/} ] && source ~/.fzf/shell/key-bindings.${SHELL##*/}
    [ -f ~/.fzf/shell/completion.${SHELL##*/} ] && source ~/.fzf/shell/completion.${SHELL##*/}
  fi
  
  # FZF configuration
  export FZF_DEFAULT_COMMAND='find . -type f -not -path "*/\.git/*" -not -path "*/node_modules/*" -not -path "*/\.cache/*"'
  export FZF_DEFAULT_OPTS='--height 40% --layout=reverse --border'
  export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
  export FZF_ALT_C_COMMAND='find . -type d -not -path "*/\.git/*" -not -path "*/node_modules/*" -not -path "*/\.cache/*"'
fi
`
		return configureRcFile(shell, fzfConfig)
		
	case "fish":
		fzfConfig := `
# Added by bootstrap-cli - fzf configuration
if type -q fzf
  # Try to source fzf key bindings
  if test -f /usr/share/fish/vendor_functions.d/fzf_key_bindings.fish
    source /usr/share/fish/vendor_functions.d/fzf_key_bindings.fish
  else if test -f /usr/local/share/fish/vendor_functions.d/fzf_key_bindings.fish
    source /usr/local/share/fish/vendor_functions.d/fzf_key_bindings.fish
  else if test -f ~/.fzf/shell/key-bindings.fish
    source ~/.fzf/shell/key-bindings.fish
  end
  
  if functions -q fzf_key_bindings
    fzf_key_bindings
  end
  
  # FZF configuration
  set -gx FZF_DEFAULT_COMMAND 'find . -type f -not -path "*/\.git/*" -not -path "*/node_modules/*" -not -path "*/\.cache/*"'
  set -gx FZF_DEFAULT_OPTS '--height 40% --layout=reverse --border'
  set -gx FZF_CTRL_T_COMMAND $FZF_DEFAULT_COMMAND
  set -gx FZF_ALT_C_COMMAND 'find . -type d -not -path "*/\.git/*" -not -path "*/node_modules/*" -not -path "*/\.cache/*"'
end
`
		return configureRcFile(shell, fzfConfig)
	}

	return nil
}

func configureBat(shell string) error {
	switch shell {
	case "zsh", "bash":
		batConfig := `
# Added by bootstrap-cli - bat configuration
if command -v bat >/dev/null 2>&1; then
  export BAT_THEME="Monokai Extended"
  export BAT_STYLE="plain"
  alias cat="bat --paging=never"
  alias batl="bat --paging=always"
fi
`
		return configureRcFile(shell, batConfig)
		
	case "fish":
		batConfig := `
# Added by bootstrap-cli - bat configuration
if type -q bat
  set -gx BAT_THEME "Monokai Extended"
  set -gx BAT_STYLE "plain"
  alias cat "bat --paging=never"
  alias batl "bat --paging=always"
end
`
		return configureRcFile(shell, batConfig)
	}

	return nil
}

func configureLsd(shell string) error {
	switch shell {
	case "zsh", "bash":
		lsdConfig := `
# Added by bootstrap-cli - lsd configuration
if command -v lsd >/dev/null 2>&1; then
  alias ls="lsd"
  alias ll="lsd -l"
  alias la="lsd -la"
  alias lt="lsd --tree"
fi
`
		return configureRcFile(shell, lsdConfig)
		
	case "fish":
		lsdConfig := `
# Added by bootstrap-cli - lsd configuration
if type -q lsd
  alias ls "lsd"
  alias ll "lsd -l"
  alias la "lsd -la"
  alias lt "lsd --tree"
end
`
		return configureRcFile(shell, lsdConfig)
	}

	return nil
}

func configureRipgrep(shell string) error {
	switch shell {
	case "zsh", "bash":
		rgConfig := `
# Added by bootstrap-cli - ripgrep configuration
if command -v rg >/dev/null 2>&1; then
  export RIPGREP_CONFIG_PATH="$HOME/.ripgreprc"
  alias rg="rg --smart-case"
  
  # Create basic ripgrep config if it doesn't exist
  if [ ! -f "$HOME/.ripgreprc" ]; then
    echo "# Don't let ripgrep vomit really long lines to my terminal, and show a preview.
--max-columns=150
--max-columns-preview

# Add 'web' type for JS/TS/CSS
--type-add=web:*.{html,css,js,ts,jsx,tsx}

# Using glob patterns to include/exclude files or folders
--glob=!.git/*
--glob=!node_modules/*

# Set the colors.
--colors=line:none
--colors=line:style:bold

# Because who cares about case!?
--smart-case" > "$HOME/.ripgreprc"
  fi
fi
`
		return configureRcFile(shell, rgConfig)
		
	case "fish":
		rgConfig := `
# Added by bootstrap-cli - ripgrep configuration
if type -q rg
  set -gx RIPGREP_CONFIG_PATH "$HOME/.ripgreprc"
  alias rg "rg --smart-case"
  
  # Create basic ripgrep config if it doesn't exist
  if not test -f "$HOME/.ripgreprc"
    echo "# Don't let ripgrep vomit really long lines to my terminal, and show a preview.
--max-columns=150
--max-columns-preview

# Add 'web' type for JS/TS/CSS
--type-add=web:*.{html,css,js,ts,jsx,tsx}

# Using glob patterns to include/exclude files or folders
--glob=!.git/*
--glob=!node_modules/*

# Set the colors.
--colors=line:none
--colors=line:style:bold

# Because who cares about case!?
--smart-case" > "$HOME/.ripgreprc"
  end
end
`
		return configureRcFile(shell, rgConfig)
	}

	return nil
}

func configureFd(shell string) error {
	switch shell {
	case "zsh", "bash":
		fdConfig := `
# Added by bootstrap-cli - fd configuration
if command -v fd >/dev/null 2>&1 || command -v fdfind >/dev/null 2>&1; then
  # Handle different binary names between distributions
  FD_CMD="fd"
  if ! command -v fd >/dev/null 2>&1 && command -v fdfind >/dev/null 2>&1; then
    FD_CMD="fdfind"
    alias fd="fdfind"
  fi
  
  # Use fd for fzf if available
  export FZF_DEFAULT_COMMAND="$FD_CMD --type file --follow --hidden --exclude .git --exclude node_modules"
  export FZF_CTRL_T_COMMAND="$FZF_DEFAULT_COMMAND"
  export FZF_ALT_C_COMMAND="$FD_CMD --type directory --follow --hidden --exclude .git --exclude node_modules"
fi
`
		return configureRcFile(shell, fdConfig)
		
	case "fish":
		fdConfig := `
# Added by bootstrap-cli - fd configuration
set FD_CMD "fd"
if not type -q fd; and type -q fdfind
  set FD_CMD "fdfind"
  alias fd "fdfind"
end

if type -q $FD_CMD
  # Use fd for fzf if available
  set -gx FZF_DEFAULT_COMMAND "$FD_CMD --type file --follow --hidden --exclude .git --exclude node_modules"
  set -gx FZF_CTRL_T_COMMAND "$FZF_DEFAULT_COMMAND"
  set -gx FZF_ALT_C_COMMAND "$FD_CMD --type directory --follow --hidden --exclude .git --exclude node_modules"
end
`
		return configureRcFile(shell, fdConfig)
	}

	return nil
}

func configureLazyGit(shell string) error {
	switch shell {
	case "zsh", "bash":
		lazygitConfig := `
# Added by bootstrap-cli - lazygit configuration
if command -v lazygit >/dev/null 2>&1; then
  alias lg="lazygit"
  
  # Create config dir if needed
  LAZYGIT_CONFIG_DIR="$HOME/.config/lazygit"
  [ ! -d "$LAZYGIT_CONFIG_DIR" ] && mkdir -p "$LAZYGIT_CONFIG_DIR"
fi
`
		return configureRcFile(shell, lazygitConfig)
		
	case "fish":
		lazygitConfig := `
# Added by bootstrap-cli - lazygit configuration
if type -q lazygit
  alias lg "lazygit"
  
  # Create config dir if needed
  set LAZYGIT_CONFIG_DIR "$HOME/.config/lazygit"
  if not test -d "$LAZYGIT_CONFIG_DIR"
    mkdir -p "$LAZYGIT_CONFIG_DIR"
  end
end
`
		return configureRcFile(shell, lazygitConfig)
	}

	return nil
}

// restartWithCountdown shows a countdown and then restarts the shell
func restartWithCountdown(shellName string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("⚠️ Shell restart not supported on Windows. Please restart your shell manually.")
		return nil
	}
	
	// Run a countdown timer before restarting (fixed at 5 seconds)
	countdownSeconds := 5
	message := fmt.Sprintf("Restarting with %s in", shellName)
	
	// Instead of using prompts.RunTimer, we'll create an explicit countdown
	// that gives more time for the user to see what's happening
	for i := countdownSeconds; i > 0; i-- {
		fmt.Printf("\r%s %d...", message, i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\r%s 0...\n", message)
	
	// Give some extra buffer time for terminal output to complete
	fmt.Printf("\n🔄 Restarting shell now...\n")
	time.Sleep(1 * time.Second)
	
	// Determine the absolute path to the shell
	shellPath, err := exec.LookPath(shellName)
	if err != nil {
		return fmt.Errorf("shell '%s' not found in PATH: %w", shellName, err)
	}

	// Instead of using syscall.Exec, create a script that will replace the current process
	// This avoids some race conditions with the terminal
	tempScript, err := os.CreateTemp("", "bootstrap-shell-restart-*.sh")
	if err != nil {
		return fmt.Errorf("failed to create temporary script: %w", err)
	}
	defer os.Remove(tempScript.Name()) // Clean up the script file

	// Write shell restart script
	script := fmt.Sprintf(`#!/bin/sh
# This script is created by bootstrap-cli to restart the shell
echo "🚀 Starting %s..."
sleep 1
exec %s -l
`, shellName, shellPath)

	if _, err := tempScript.WriteString(script); err != nil {
		return fmt.Errorf("failed to write shell script: %w", err)
	}

	if err := tempScript.Close(); err != nil {
		return fmt.Errorf("failed to close shell script: %w", err)
	}

	// Make the script executable
	if err := os.Chmod(tempScript.Name(), 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	// Execute the restart script
	cmd := exec.Command(tempScript.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Display a final message before restarting
	fmt.Printf("✨ Your new shell will start momentarily...\n\n")
	time.Sleep(500 * time.Millisecond)
	
	return cmd.Run()
}

// configureUserDotfiles integrates bootstrap configurations with user dotfiles
func configureUserDotfiles(config types.UserConfig) error {
	// Skip if no dotfiles path is provided
	if config.DotfilesPath == "" {
		return nil
	}
	
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	// Expand dotfiles path if it starts with ~
	dotfilesPath := config.DotfilesPath
	if strings.HasPrefix(dotfilesPath, "~/") {
		dotfilesPath = filepath.Join(home, dotfilesPath[2:])
	}
	
	// Check if dotfiles directory exists
	if _, err := os.Stat(dotfilesPath); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles directory not found: %s", dotfilesPath)
	}
	
	// Get shell config file paths
	var shellRcFiles = map[string]string{
		"zsh":  filepath.Join(home, ".zshrc"),
		"bash": filepath.Join(home, ".bashrc"),
		"fish": filepath.Join(home, ".config/fish/config.fish"),
	}
	
	// Get tool config file paths
	toolConfigPaths := map[string]string{
		"fzf":     filepath.Join(home, ".fzf.zsh"),
		"bat":     filepath.Join(home, ".config/bat/config"),
		"ripgrep": filepath.Join(home, ".ripgreprc"),
		"lazygit": filepath.Join(home, ".config/lazygit/config.yml"),
	}
	
	// First check if the user has dotfiles for the selected shell
	if shellPath, ok := shellRcFiles[config.Shell]; ok {
		userShellFile := filepath.Join(dotfilesPath, filepath.Base(shellPath))
		if _, err := os.Stat(userShellFile); err == nil {
			// User has shell dotfile, check if we need to migrate bootstrap settings
			fmt.Printf("🔍 Found shell configuration in dotfiles: %s\n", userShellFile)
			
			// Add a more detailed analysis in the future
			// For now, we'll just check for common patterns to detect if our tools are configured
			if err := analyzeAndMigrateShellConfig(shellPath, userShellFile, config); err != nil {
				fmt.Printf("⚠️ Warning: Could not analyze shell configuration: %v\n", err)
			}
		}
	}
	
	// Check for tool-specific configuration files in user dotfiles
	for tool, configPath := range toolConfigPaths {
		// Only check for tools the user has selected
		if !containsString(config.CLITools, tool) {
			continue
		}
		
		relPath, err := filepath.Rel(home, configPath)
		if err != nil {
			continue
		}
		
		userToolConfig := filepath.Join(dotfilesPath, relPath)
		if _, err := os.Stat(userToolConfig); err == nil {
			fmt.Printf("🔍 Found configuration for %s in dotfiles: %s\n", tool, userToolConfig)
			// User has this tool configured in their dotfiles
			// We can analyze and merge configurations if needed
		}
	}
	
	return nil
}

// analyzeAndMigrateShellConfig analyzes user shell config and migrates bootstrap settings if needed
func analyzeAndMigrateShellConfig(systemPath, userPath string, config types.UserConfig) error {
	systemContent, err := os.ReadFile(systemPath)
	if err != nil {
		return fmt.Errorf("failed to read system config: %w", err)
	}
	
	userContent, err := os.ReadFile(userPath)
	if err != nil {
		return fmt.Errorf("failed to read user config: %w", err)
	}
	
	// Look for bootstrap markers in the system config
	bootstrapLines := extractBootstrapConfigurations(string(systemContent))
	if len(bootstrapLines) == 0 {
		return nil // No bootstrap configurations found
	}
	
	// Check if user config already has these configurations
	userConfigStr := string(userContent)
	missingConfigs := []string{}
	
	for _, config := range bootstrapLines {
		if !containsConfiguration(userConfigStr, config) {
			missingConfigs = append(missingConfigs, config)
		}
	}
	
	if len(missingConfigs) == 0 {
		fmt.Println("✅ All bootstrap configurations already exist in dotfiles")
		return nil
	}
	
	// Ask user if they want to migrate missing configurations
	fmt.Printf("📝 Found %d bootstrap configurations not in your dotfiles\n", len(missingConfigs))
	fmt.Println("ℹ️ These will be automatically merged during the next bootstrap run")
	
	// In the future, we can add interactive prompting here
	
	return nil
}

// extractBootstrapConfigurations finds all bootstrap-cli configurations in a file
func extractBootstrapConfigurations(content string) []string {
	lines := strings.Split(content, "\n")
	var bootstrapConfigs []string
	var currentConfig strings.Builder
	inBootstrapSection := false
	
	for _, line := range lines {
		if strings.Contains(line, "Added by bootstrap-cli") {
			// Start of a bootstrap config section
			if currentConfig.Len() > 0 {
				bootstrapConfigs = append(bootstrapConfigs, currentConfig.String())
				currentConfig.Reset()
			}
			inBootstrapSection = true
			currentConfig.WriteString(line + "\n")
		} else if inBootstrapSection {
			if len(strings.TrimSpace(line)) == 0 {
				// Empty line could be end of section
				if strings.HasSuffix(currentConfig.String(), "\n\n") {
					// Two newlines in a row - end of section
					inBootstrapSection = false
					bootstrapConfigs = append(bootstrapConfigs, currentConfig.String())
					currentConfig.Reset()
				} else {
					// Just a blank line within the section
					currentConfig.WriteString(line + "\n")
				}
			} else if strings.HasPrefix(line, "#") && !strings.Contains(line, "bootstrap") {
				// Regular comment, still part of the section
				currentConfig.WriteString(line + "\n")
			} else if !strings.HasPrefix(line, "#") {
				// Non-comment line, part of the configuration
				currentConfig.WriteString(line + "\n")
			} else {
				// End of section
				inBootstrapSection = false
				bootstrapConfigs = append(bootstrapConfigs, currentConfig.String())
				currentConfig.Reset()
			}
		}
	}
	
	// Add the last config if we're still in one
	if currentConfig.Len() > 0 {
		bootstrapConfigs = append(bootstrapConfigs, currentConfig.String())
	}
	
	return bootstrapConfigs
}

// containsString checks if a string slice contains a specific value
func containsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
} 