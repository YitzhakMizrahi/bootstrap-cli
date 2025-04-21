package config

// UserConfig represents the user's configuration
type UserConfig struct {
	Shell            string            `yaml:"shell"`
	PluginManager    string            `yaml:"plugin_manager"`
	Prompt           string            `yaml:"prompt"`
	CLITools         []string          `yaml:"cli_tools"`
	Languages        []string          `yaml:"languages"`
	PackageManagers  map[string]string `yaml:"package_managers"`
	UseRelativeLinks bool              `yaml:"use_relative_links"`
	BackupExisting   bool              `yaml:"backup_existing"`
	Editors          []string          `yaml:"editors"`
	DotfilesPath     string            `yaml:"dotfiles_path"`
	DevMode          bool              `yaml:"dev_mode"`
} 