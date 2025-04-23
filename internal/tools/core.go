package tools

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// ToolCategory represents a category of tools
type ToolCategory struct {
	Name        string
	Description string
	Tools       []install.Tool
}

// GetToolCategories returns all available tool categories
func GetToolCategories() []ToolCategory {
	return []ToolCategory{
		{
			Name:        "Essential Tools",
			Description: "Core development tools required for most workflows",
			Tools:       essentialTools(),
		},
		{
			Name:        "Modern CLI Tools",
			Description: "Modern alternatives to traditional command-line tools",
			Tools:       modernCliTools(),
		},
		{
			Name:        "System Tools",
			Description: "System monitoring and management utilities",
			Tools:       systemTools(),
		},
	}
}

func essentialTools() []install.Tool {
	return []install.Tool{
		{
			Name:        "git",
			PackageName: "git",
			Description: "Distributed version control system",
			VerifyCommand: "git --version",
		},
		{
			Name:        "curl",
			PackageName: "curl",
			Description: "Command line tool for transferring data with URLs",
			VerifyCommand: "curl --version",
		},
		{
			Name:        "wget",
			PackageName: "wget",
			Description: "Non-interactive network downloader",
			VerifyCommand: "wget --version",
		},
		{
			Name:        "tmux",
			PackageName: "tmux",
			Description: "Terminal multiplexer",
			VerifyCommand: "tmux -V",
		},
	}
}

func modernCliTools() []install.Tool {
	return []install.Tool{
		{
			Name:        "ripgrep",
			PackageName: "ripgrep",
			Description: "Modern grep alternative written in Rust",
			VerifyCommand: "rg --version",
		},
		{
			Name:        "bat",
			PackageName: "bat",
			Description: "Cat clone with syntax highlighting and Git integration",
			VerifyCommand: "bat --version",
		},
		{
			Name:        "fzf",
			PackageName: "fzf",
			Description: "Command-line fuzzy finder",
			VerifyCommand: "fzf --version",
		},
		{
			Name:        "exa",
			PackageName: "exa",
			Description: "Modern replacement for ls",
			VerifyCommand: "exa --version",
		},
	}
}

func systemTools() []install.Tool {
	return []install.Tool{
		{
			Name:        "htop",
			PackageName: "htop",
			Description: "Interactive process viewer",
			VerifyCommand: "htop --version",
		},
		{
			Name:        "btop",
			PackageName: "btop",
			Description: "Resource monitor with additional features",
			VerifyCommand: "btop --version",
		},
		{
			Name:        "neofetch",
			PackageName: "neofetch",
			Description: "System information tool",
			VerifyCommand: "neofetch --version",
		},
	}
}

// CoreTools returns all available tools (for backward compatibility)
func CoreTools() []install.Tool {
	var allTools []install.Tool
	for _, category := range GetToolCategories() {
		allTools = append(allTools, category.Tools...)
	}
	return allTools
}

// InstallOptions contains options for tool installation
type InstallOptions struct {
	Logger           *log.Logger
	PackageManager   packages.PackageManager
	Tools           []*install.Tool
	SkipVerification bool
	AdditionalPaths  []string // Additional paths to search for binaries during verification
}

var selectedTools []*install.Tool

// GetSelectedTools returns the list of tools selected during initialization
func GetSelectedTools() []*install.Tool {
	return selectedTools
}

// SetSelectedTools sets the list of tools to be installed
func SetSelectedTools(tools []*install.Tool) {
	selectedTools = tools
} 