package shell

import (
	shellconfig "github.com/YitzhakMizrahi/bootstrap-cli/pkg/config/shell"
)

// ShellManagerInterface defines the interface for shell management
type ShellManagerInterface interface {
	Install(shell string) error
	InstallPrompt(prompt, shell string) error
	SetDefault(shell string) error
	GetCurrent() (string, error)
	GenerateConfig(shell Type, data *shellconfig.Data) error
	RestoreConfig(shell Type, path string) error
	Setup(cfg *Config) error
	InstallPluginManager(pm PluginManager, shell string) error
	InstallPlugin(plugin string) error
	UninstallPlugin(plugin string) error
	ListAvailablePlugins() []string
	BackupConfig() error
	Cleanup() error
	AddAlias(name, command string)
	AddEnvVar(name, value string)
	AddPath(path string)
} 