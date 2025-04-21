package types

// UserConfig holds all the user's bootstrap preferences
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
type ToolOption struct {
	Name       string   // Tool name (e.g. "pnpm")
	Category   string   // Tool category (e.g. "pkgmgr", "cli", "prompt")
	Group      string   // Language/platform it belongs to (e.g. "node")
	Default    bool     // Whether it's preselected
	Alternates []string // Alternative tools (e.g. yarn/npm for pnpm)
} 