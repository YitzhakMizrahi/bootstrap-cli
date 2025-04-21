package shell

import (
	plugins "github.com/YitzhakMizrahi/bootstrap-cli/pkg/plugins/shell"
)

// Type represents the shell type
type Type string

const (
	// Bash shell
	Bash Type = "bash"
	// Zsh shell
	Zsh Type = "zsh"
	// Fish shell
	Fish Type = "fish"
)

// PluginManager represents a shell plugin manager
type PluginManager string

const (
	// OhMyZsh plugin manager for Zsh
	OhMyZsh PluginManager = "oh-my-zsh"
	// Antigen plugin manager for Zsh
	Antigen PluginManager = "antigen"
	// Zinit plugin manager for Zsh
	Zinit PluginManager = "zinit"
	// Fisherman plugin manager for Fish
	Fisherman PluginManager = "fisherman"
)

// Config represents shell configuration
type Config struct {
	Type         Type
	PluginMgr    PluginManager
	Plugins      []string
	CustomConfig string
	Aliases      map[string]string
	EnvVars      map[string]string
	Path         []string
}

// ShellInfo contains information about a shell
type ShellInfo struct {
	Name        string
	Description string
	Package     string
}

// PluginManagerInfo contains information about a plugin manager
type PluginManagerInfo struct {
	Name            string
	Description     string
	CompatibleShells []string
	Install         func(shell string) error
}

// PromptInfo contains information about a shell prompt
type PromptInfo struct {
	Name            string
	Description     string
	CompatibleShells []string
	Package         string
	Install         func(shell string) error
	Configure       func(shell string) error
}

var (
	// Shells maps shell names to their information
	shells = map[string]ShellInfo{
		"zsh": {
			Name:        "Zsh",
			Description: "The Z Shell",
			Package:     "zsh",
		},
		"bash": {
			Name:        "Bash",
			Description: "Bourne Again Shell",
			Package:     "bash",
		},
		"fish": {
			Name:        "Fish",
			Description: "Friendly Interactive Shell",
			Package:     "fish",
		},
	}

	// PluginManagers maps plugin manager names to their information
	pluginManagers = map[string]PluginManagerInfo{
		"zinit": {
			Name:            "Zinit",
			Description:     "Fast and feature-rich plugin manager for Zsh",
			CompatibleShells: []string{"zsh"},
			Install:         plugins.InstallZinit,
		},
		"oh-my-zsh": {
			Name:            "Oh My Zsh",
			Description:     "Community-driven framework for Zsh",
			CompatibleShells: []string{"zsh"},
			Install:         plugins.InstallOhMyZsh,
		},
		"bash-it": {
			Name:            "Bash-it",
			Description:     "Community framework for Bash",
			CompatibleShells: []string{"bash"},
			Install:         plugins.InstallBashIt,
		},
		"oh-my-bash": {
			Name:            "Oh My Bash",
			Description:     "Community-driven framework for Bash",
			CompatibleShells: []string{"bash"},
			Install:         plugins.InstallOhMyBash,
		},
		"fisher": {
			Name:            "Fisher",
			Description:     "Plugin manager for Fish",
			CompatibleShells: []string{"fish"},
			Install:         plugins.InstallFisher,
		},
		"oh-my-fish": {
			Name:            "Oh My Fish",
			Description:     "Community Fish framework",
			CompatibleShells: []string{"fish"},
			Install:         plugins.InstallOhMyFish,
		},
	}

	// Prompts maps prompt names to their information
	prompts = map[string]PromptInfo{
		"starship": {
			Name:            "Starship",
			Description:     "Cross-shell customizable prompt",
			CompatibleShells: []string{"zsh", "bash", "fish"},
			Package:         "starship",
			Install:         plugins.InstallStarship,
			Configure:       plugins.ConfigureStarship,
		},
		"powerlevel10k": {
			Name:            "Powerlevel10k",
			Description:     "Fast and customizable Zsh theme",
			CompatibleShells: []string{"zsh"},
			Install:         plugins.InstallPowerlevel10k,
			Configure:       plugins.ConfigurePowerlevel10k,
		},
		"pure": {
			Name:            "Pure",
			Description:     "Pretty, minimal and fast Zsh prompt",
			CompatibleShells: []string{"zsh"},
			Install:         plugins.InstallPure,
			Configure:       plugins.ConfigurePure,
		},
		"oh-my-posh": {
			Name:            "Oh My Posh",
			Description:     "Cross-platform, customizable prompt",
			CompatibleShells: []string{"zsh", "bash", "fish"},
			Package:         "oh-my-posh",
			Install:         plugins.InstallOhMyPosh,
			Configure:       plugins.ConfigureOhMyPosh,
		},
	}
)

// GetShellInfo returns information about a shell
func GetShellInfo(name string) (ShellInfo, bool) {
	info, exists := shells[name]
	return info, exists
}

// GetPluginManagerInfo returns information about a plugin manager
func GetPluginManagerInfo(name string) (PluginManagerInfo, bool) {
	info, exists := pluginManagers[name]
	return info, exists
}

// GetPromptInfo returns information about a prompt
func GetPromptInfo(name string) (PromptInfo, bool) {
	info, exists := prompts[name]
	return info, exists
}

// SupportsShell checks if a plugin manager supports a shell
func (pm PluginManagerInfo) SupportsShell(shell string) bool {
	for _, s := range pm.CompatibleShells {
		if s == shell {
			return true
		}
	}
	return false
}

// SupportsShell checks if a prompt supports a shell
func (p PromptInfo) SupportsShell(shell string) bool {
	for _, s := range p.CompatibleShells {
		if s == shell {
			return true
		}
	}
	return false
}