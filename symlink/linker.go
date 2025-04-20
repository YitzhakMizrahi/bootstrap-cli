package symlink

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/types"
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
func LinkDotfiles(config types.UserConfig) error {
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
		return fmt.Errorf("dotfiles directory not found: %s", dotfilesPath)
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
func discoverDotfiles(dotfilesPath string, config types.UserConfig) ([]DotfileMapping, error) {
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
func createSymlink(mapping DotfileMapping, config types.UserConfig) error {
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