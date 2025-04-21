package symlink

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
)

// DotfileMapping defines a mapping between a source file in the dotfiles repo
// and the target location in the user's home directory
type DotfileMapping struct {
	Source     string // Path relative to the dotfiles repo
	Target     string // Target path (usually in home directory)
	SourceAbs  string // Absolute path to source
	TargetAbs  string // Absolute path to target
	IsDir      bool   // Whether the source is a directory
	Exists     bool   // Whether the target already exists
	IsSymlink  bool   // Whether the target is already a symlink
	NeedsLink  bool   // Whether we need to create/update the symlink
}

// Common dotfile path mappings
var commonDotfileMappings = map[string]string{
	".zshrc":        "~/.zshrc",
	".bashrc":       "~/.bashrc",
	".bash_profile": "~/.bash_profile",
	".profile":      "~/.profile",
	".tmux.conf":    "~/.tmux.conf",
	".vimrc":        "~/.vimrc",
	"nvim":          "~/.config/nvim",
	"starship.toml": "~/.config/starship.toml",
	".gitconfig":    "~/.gitconfig",
	".p10k.zsh":     "~/.p10k.zsh",
	"alacritty":     "~/.config/alacritty",
	"kitty":         "~/.config/kitty",
}

// LinkDotfiles creates symbolic links from the dotfiles repository to the appropriate
// locations in the user's home directory
func LinkDotfiles(config config.UserConfig) error {
	// Expand the dotfiles path if it starts with ~
	dotfilesPath := config.DotfilesPath
	if strings.HasPrefix(dotfilesPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		dotfilesPath = filepath.Join(home, dotfilesPath[2:])
	}

	// Check if dotfiles directory exists
	if _, err := os.Stat(dotfilesPath); os.IsNotExist(err) {
		// Dotfiles directory doesn't exist, prompt user to create it
		fmt.Printf("📂 Dotfiles directory not found: %s\n", dotfilesPath)
		fmt.Printf("Would you like to create it with basic template files? (y/n): ")
		
		var response string
		fmt.Scanln(&response)
		
		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			// Create directory
			if err := os.MkdirAll(dotfilesPath, 0755); err != nil {
				return fmt.Errorf("failed to create dotfiles directory: %w", err)
			}
			
			// Create basic template files
			if err := createDotfilesTemplates(dotfilesPath); err != nil {
				return fmt.Errorf("failed to create template files: %w", err)
			}
			
			fmt.Printf("✅ Created dotfiles directory at %s with template files\n", dotfilesPath)
			
			// Ask if the user wants to initialize a git repository
			fmt.Print("Would you like to initialize a Git repository for your dotfiles? (y/n): ")
			fmt.Scanln(&response)
			
			if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
				if err := initGitRepo(dotfilesPath); err != nil {
					fmt.Printf("⚠️ Failed to initialize Git repository: %v\n", err)
				} else {
					fmt.Println("✅ Initialized Git repository for your dotfiles")
				}
			}
		} else {
			return fmt.Errorf("dotfiles directory not found: %s", dotfilesPath)
		}
	}

	// Discover dotfiles and their mappings
	mappings, err := discoverDotfiles(dotfilesPath, config)
	if err != nil {
		return fmt.Errorf("failed to discover dotfiles: %w", err)
	}

	if len(mappings) == 0 {
		fmt.Println("⚠️  No dotfiles found to link")
		return nil
	}

	// Create symlinks
	for _, mapping := range mappings {
		if err := createSymlink(mapping, config); err != nil {
			fmt.Printf("❌ Failed to link %s: %v\n", mapping.Source, err)
		} else {
			fmt.Printf("✅ Linked %s -> %s\n", mapping.Target, mapping.Source)
		}
	}

	return nil
}

// discoverDotfiles finds dotfiles in the repository and creates mappings
func discoverDotfiles(dotfilesPath string, config config.UserConfig) ([]DotfileMapping, error) {
	var mappings []DotfileMapping
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// Walk through the dotfiles directory
	err = filepath.WalkDir(dotfilesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == dotfilesPath {
			return nil
		}

		// Skip hidden directories (e.g., .git)
		if d.IsDir() && strings.HasPrefix(filepath.Base(path), ".") {
			return filepath.SkipDir
		}

		// Get the relative path from the dotfiles directory
		relPath, err := filepath.Rel(dotfilesPath, path)
		if err != nil {
			return err
		}

		// Skip if path contains ".git" segment
		if strings.Contains(relPath, ".git") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if we have a mapping for this file
		targetPath, found := findTargetPath(relPath, home)
		if !found {
			// Skip if we don't know where to link this file
			return nil
		}

		// Create the mapping
		targetExists := false
		targetIsSymlink := false
		if info, err := os.Lstat(targetPath); err == nil {
			targetExists = true
			targetIsSymlink = info.Mode()&fs.ModeSymlink != 0
		}

		mapping := DotfileMapping{
			Source:     relPath,
			Target:     targetPath,
			SourceAbs:  path,
			TargetAbs:  targetPath,
			IsDir:      d.IsDir(),
			Exists:     targetExists,
			IsSymlink:  targetIsSymlink,
			NeedsLink:  true, // Assume we need to link by default
		}

		// Check if we need to create/update the symlink
		if targetIsSymlink {
			// If it's already a symlink, check if it points to the right place
			linkDest, err := os.Readlink(targetPath)
			if err == nil && (linkDest == path || 
				(config.UseRelativeLinks && isEqualRelative(linkDest, path, targetPath))) {
				mapping.NeedsLink = false
			}
		}

		mappings = append(mappings, mapping)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return mappings, nil
}

// findTargetPath determines the appropriate target path for a dotfile
func findTargetPath(relPath string, home string) (string, bool) {
	// Check if we have a direct mapping
	if targetPath, found := commonDotfileMappings[relPath]; found {
		// Expand home directory
		if strings.HasPrefix(targetPath, "~/") {
			targetPath = filepath.Join(home, targetPath[2:])
		}
		return targetPath, true
	}

	// For files in config directory
	if strings.HasPrefix(relPath, ".config/") {
		return filepath.Join(home, relPath), true
	}

	// For dotfiles in the root of the repo
	if strings.HasPrefix(filepath.Base(relPath), ".") {
		return filepath.Join(home, filepath.Base(relPath)), true
	}

	// We don't know where to link this file
	return "", false
}

// createSymlink creates a symbolic link from target to source
func createSymlink(mapping DotfileMapping, config config.UserConfig) error {
	// If target exists and is not a symlink, back it up if configured to do so
	if mapping.Exists && !mapping.IsSymlink && config.BackupExisting {
		backupPath := mapping.TargetAbs + ".bak." + time.Now().Format("20060102150405")
		fmt.Printf("📦 Backing up %s to %s\n", mapping.Target, backupPath)
		if err := os.Rename(mapping.TargetAbs, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
		// Update mapping since we've moved the target
		mapping.Exists = false
	} else if mapping.Exists && !mapping.IsSymlink && !config.BackupExisting {
		// If we're not backing up, remove the existing file
		fmt.Printf("🗑️  Removing existing %s\n", mapping.Target)
		if err := os.RemoveAll(mapping.TargetAbs); err != nil {
			return fmt.Errorf("failed to remove existing file: %w", err)
		}
		mapping.Exists = false
	} else if mapping.Exists && mapping.IsSymlink {
		// Remove existing symlink
		fmt.Printf("🔄 Updating symlink for %s\n", mapping.Target)
		if err := os.Remove(mapping.TargetAbs); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %w", err)
		}
		mapping.Exists = false
	}

	// Create parent directories if they don't exist
	targetDir := filepath.Dir(mapping.TargetAbs)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %w", err)
	}

	// Create the symlink
	var linkPath string
	if config.UseRelativeLinks {
		// Calculate relative path from target to source
		relPath, err := filepath.Rel(filepath.Dir(mapping.TargetAbs), mapping.SourceAbs)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}
		linkPath = relPath
	} else {
		// Use absolute path
		linkPath = mapping.SourceAbs
	}

	if config.DevMode {
		fmt.Printf("🔍 Creating symlink: %s -> %s\n", mapping.TargetAbs, linkPath)
	}

	return os.Symlink(linkPath, mapping.TargetAbs)
}

// isEqualRelative checks if two paths are equal when considering relative paths
func isEqualRelative(existingLink, sourcePath, targetPath string) bool {
	// Calculate where the existing link would point to
	targetDir := filepath.Dir(targetPath)
	resolvedPath := filepath.Join(targetDir, existingLink)
	
	// Normalize both paths for comparison
	sourcePath, _ = filepath.Abs(sourcePath)
	resolvedPath, _ = filepath.Abs(resolvedPath)
	
	return resolvedPath == sourcePath
}

// createDotfilesTemplates creates basic template files in the dotfiles directory
func createDotfilesTemplates(dotfilesPath string) error {
	// Create template files with some sensible defaults
	templates := map[string]string{
		".zshrc": `# ~/.zshrc - ZSH Configuration
# Created by bootstrap-cli

# Set some basic options
setopt autocd
setopt extendedglob
setopt nomatch
setopt notify

# Basic history configuration
HISTFILE=~/.histfile
HISTSIZE=1000
SAVEHIST=1000
setopt appendhistory

# Basic prompt
PS1="%n@%m:%~$ "

# Add paths
export PATH=$HOME/bin:$PATH

# Aliases
alias ls='ls --color=auto'
alias la='ls -la'
alias ll='ls -l'
alias grep='grep --color=auto'

# Load any local config if it exists
[[ -f ~/.zshrc.local ]] && source ~/.zshrc.local
`,
		".bashrc": `# ~/.bashrc - Bash Configuration
# Created by bootstrap-cli

# If not running interactively, don't do anything
[[ $- != *i* ]] && return

# Basic prompt
PS1='[\u@\h \W]\$ '

# Add paths
export PATH=$HOME/bin:$PATH

# Aliases
alias ls='ls --color=auto'
alias la='ls -la'
alias ll='ls -l'
alias grep='grep --color=auto'

# Load any local config if it exists
[[ -f ~/.bashrc.local ]] && source ~/.bashrc.local
`,
		".gitconfig": `# ~/.gitconfig - Git Configuration
# Created by bootstrap-cli

[user]
	name = Your Name
	email = your.email@example.com

[core]
	editor = vim
	autocrlf = input

[color]
	ui = auto

[alias]
	st = status
	ci = commit
	co = checkout
	br = branch
	lg = log --graph --pretty=format:'%Cred%h%Creset -%C(yellow)%d%Creset %s %Cgreen(%cr) %C(bold blue)<%an>%Creset' --abbrev-commit
`,
		"README.md": `# Dotfiles

This is your personal dotfiles repository, created by bootstrap-cli.

## Structure

- '.zshrc' - ZSH shell configuration
- '.bashrc' - Bash shell configuration
- '.gitconfig' - Git configuration

## Usage

These files are symlinked to their appropriate locations by bootstrap-cli.
You can edit them here and the changes will be reflected in your system.

## Adding More Files

Add more configuration files to this repository and run:

` + "```" + `
bootstrap-cli link
` + "```" + `

to symlink them to your home directory.
`,
	}

	// Create each template file
	for filename, content := range templates {
		filePath := filepath.Join(dotfilesPath, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create template file %s: %w", filename, err)
		}
	}

	return nil
}

// initGitRepo initializes a Git repository in the dotfiles directory
func initGitRepo(dotfilesPath string) error {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed, skipping repository initialization")
	}

	// Initialize the repository
	cmd := exec.Command("git", "init")
	cmd.Dir = dotfilesPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create a .gitignore file
	gitignore := `# bootstrap-cli gitignore
.DS_Store
*.swp
*~
.history/
`
	if err := os.WriteFile(filepath.Join(dotfilesPath, ".gitignore"), []byte(gitignore), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Stage all files
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = dotfilesPath
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create initial commit
	cmd = exec.Command("git", "commit", "-m", "Initial dotfiles commit by bootstrap-cli")
	cmd.Dir = dotfilesPath
	cmd.Env = append(os.Environ(), 
		"GIT_AUTHOR_NAME=bootstrap-cli",
		"GIT_AUTHOR_EMAIL=bootstrap-cli@example.com",
		"GIT_COMMITTER_NAME=bootstrap-cli",
		"GIT_COMMITTER_EMAIL=bootstrap-cli@example.com")
		
	// Ignore errors from commit as it might fail if user needs to configure git first
	cmd.Run()

	// Provide instructions for setting up a remote repository
	fmt.Println("\n💡 To push your dotfiles to GitHub or other remote repositories:")
	fmt.Println("   1. Create a new repository on GitHub")
	fmt.Println("   2. Run the following commands in your dotfiles directory:")
	fmt.Printf("      cd %s\n", dotfilesPath)
	fmt.Println("      git remote add origin https://github.com/yourusername/dotfiles.git")
	fmt.Println("      git push -u origin main")

	return nil
}