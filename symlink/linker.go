// symlink/linker.go
package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/types"
)

var shellFiles = map[string]string{
	"zsh": ".zshrc",
	"bash": ".bashrc",
	"bash_profile": ".bash_profile",
	"profile": ".profile",
}

var promptFiles = map[string]string{
	"starship": ".config/starship.toml",
}

var cliToolFiles = map[string]map[string]string{
	"lsd": {
		".config/lsd/config.yaml": ".config/lsd/config.yaml",
	},
	"lazygit": {
		".config/lazygit/config.yml": ".config/lazygit/config.yml",
	},
	"yazi": {
		".config/yazi/config.toml": ".config/yazi/config.toml",
	},
	"nvim": {
		".config/nvim/init.vim": ".config/nvim/init.vim",
		".config/nvim/init.lua": ".config/nvim/init.lua",
	},
	"zoxide": {},
	"fzf": {},
	"bat": {},
	"eza": {},
	"curlie": {},
	"tmux": {},
}

var miscFiles = map[string]string{
	".gitconfig": ".gitconfig",
	".gitignore_global": ".gitignore_global",
	".ssh/config": ".ssh/config",
	".config/pyenv/config.toml": ".config/pyenv/config.toml",
	".config/nvm/default-packages": ".config/nvm/default-packages",
}

func LinkDotfiles(cfg types.UserConfig) error {
	dotfilesPath := expandPath(cfg.DotfilesPath)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	targets := make(map[string]string)

	// Shell
	if f, ok := shellFiles[cfg.Shell]; ok {
		targets[f] = f
	}

	// Prompt
	if f, ok := promptFiles[cfg.Prompt]; ok {
		targets[f] = f
	}

	// CLI tools
	for _, tool := range cfg.CLITools {
		if fileGroup, ok := cliToolFiles[tool]; ok {
			for src, dest := range fileGroup {
				targets[src] = dest
			}
		}
	}

	// Misc always-on files
	for src, dest := range miscFiles {
		targets[src] = dest
	}

	for srcRel, destRel := range targets {
		srcPath := filepath.Join(dotfilesPath, srcRel)
		destPath := filepath.Join(homeDir, destRel)

		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("⚠️  Skipping missing source: %s\n", srcPath)
			continue
		}

		if _, err := os.Lstat(destPath); err == nil {
			if cfg.BackupExisting {
				backupPath := destPath + ".bak"
				fmt.Printf("📦 Backing up %s → %s\n", destPath, backupPath)
				os.Rename(destPath, backupPath)
			} else {
				fmt.Printf("🔁 Overwriting existing: %s\n", destPath)
				os.Remove(destPath)
			}
		}

		var linkErr error
		if cfg.UseRelativeLinks {
			relSrc, err := filepath.Rel(filepath.Dir(destPath), srcPath)
			if err != nil {
				return fmt.Errorf("failed to compute relative path: %w", err)
			}
			linkErr = os.Symlink(relSrc, destPath)
		} else {
			linkErr = os.Symlink(srcPath, destPath)
		}

		if linkErr != nil {
			fmt.Printf("❌ Failed to link %s → %s: %v\n", srcPath, destPath, linkErr)
		} else {
			fmt.Printf("🔗 Linked %s → %s\n", destPath, srcPath)
		}
	}

	return nil
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return p
}
