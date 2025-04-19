// installer/cli.go
package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallCLITools(tools []string) {
	for _, tool := range tools {
		fmt.Printf("🔧 Installing %s...\n", tool)
		switch tool {
		case "lsd":
			exec.Command("brew", "install", "lsd").Run()
		case "fzf":
			exec.Command("brew", "install", "fzf").Run()
		case "bat":
			exec.Command("brew", "install", "bat").Run()
		case "zoxide":
			exec.Command("brew", "install", "zoxide").Run()
		case "eza":
			exec.Command("brew", "install", "eza").Run()
		case "curlie":
			exec.Command("brew", "install", "curlie").Run()
		case "tmux":
			exec.Command("brew", "install", "tmux").Run()
		case "lazygit":
			exec.Command("brew", "install", "lazygit").Run()
		case "yazi":
			exec.Command("brew", "install", "yazi").Run()
		case "starship":
			exec.Command("brew", "install", "starship").Run()
		default:
			fmt.Printf("⚠️  Skipping unknown CLI tool: %s\n", tool)
		}
	}
}

func InstallLanguages(langs []string, managers map[string]string) {
	for _, lang := range langs {
		fmt.Printf("📦 Installing %s...\n", lang)
		switch lang {
		case "python":
			exec.Command("brew", "install", "pyenv").Run()
		case "node":
			exec.Command("brew", "install", "nvm").Run()
			if pkg, ok := managers["node"]; ok && pkg != "npm" {
				exec.Command("npm", "install", "-g", pkg).Run()
			}
		case "go":
			exec.Command("brew", "install", "go").Run()
		case "rust":
			exec.Command("brew", "install", "rustup-init").Run()
			exec.Command("rustup-init", "-y").Run()
		default:
			fmt.Printf("⚠️  Unknown language: %s\n", lang)
		}
	}
}

func InstallEditors(editors []string) {
	for _, editor := range editors {
		fmt.Printf("📝 Setting up %s...\n", editor)
		switch editor {
		case "vim":
			exec.Command("brew", "install", "vim").Run()
		case "nvim":
			exec.Command("brew", "install", "neovim").Run()
		case "nvim (LazyVim)":
			exec.Command("brew", "install", "neovim").Run()
			cloneLazyVim()
		case "nvim (AstroNvim)":
			exec.Command("brew", "install", "neovim").Run()
			cloneAstroNvim()
		default:
			fmt.Printf("⚠️  Unknown editor: %s\n", editor)
		}
	}
}

func cloneLazyVim() {
	dest := filepath.Join(os.Getenv("HOME"), ".config", "nvim")
	cmd := exec.Command("git", "clone", "https://github.com/LazyVim/starter", dest)
	cmd.Run()
	fmt.Println("🚀 Cloned LazyVim starter")
}

func cloneAstroNvim() {
	dest := filepath.Join(os.Getenv("HOME"), ".config", "nvim")
	cmd := exec.Command("git", "clone", "https://github.com/AstroNvim/AstroNvim", dest)
	cmd.Run()
	fmt.Println("🌌 Cloned AstroNvim")
}
