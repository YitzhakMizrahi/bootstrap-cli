// types/types.go
package types

// UserConfig holds all the user's bootstrap preferences
// This gets written to ~/.config/bootstrap/config.yaml

type UserConfig struct {
	Shell             string            `yaml:"shell"`
	PluginManager     string            `yaml:"plugin_manager"`
	Prompt            string            `yaml:"prompt"`
	CLITools          []string          `yaml:"cli_tools"`
	Languages         []string          `yaml:"languages"`
	PackageManagers   map[string]string `yaml:"package_managers"`
	DotfilesPath      string            `yaml:"dotfiles_path"`
	UseRelativeLinks  bool              `yaml:"use_relative_links"`
	DevMode           bool              `yaml:"dev_mode"`
	BackupExisting    bool              `yaml:"backup_existing"`
	Editors           []string          `yaml:"editors"`
}

// ToolOption defines metadata for installable tools
// This allows prompt generation and future installer logic

type ToolOption struct {
	Name       string   // "pnpm"
	Category   string   // "pkgmgr", "cli", "prompt"
	Group      string   // Language it belongs to, like "node"
	Default    bool     // Whether it's preselected
	Alternates []string // Alternative tools (like yarn/npm for pnpm)
}
