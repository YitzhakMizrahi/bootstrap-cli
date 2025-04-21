package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shellconfig "github.com/YitzhakMizrahi/bootstrap-cli/pkg/config/shell"
	pkgmanager "github.com/YitzhakMizrahi/bootstrap-cli/pkg/core/package"
	"github.com/YitzhakMizrahi/bootstrap-cli/pkg/platform"
)

// DefaultShellManager implements the ShellManagerInterface
type DefaultShellManager struct {
	pkgManager pkgmanager.Manager
	detector   platform.Detector
	config     *Config
	homeDir    string
}

// NewShellManagerWithDeps creates a new shell manager with the given dependencies
func NewShellManagerWithDeps(pkgManager pkgmanager.Manager, detector platform.Detector, cfg *Config) (ShellManagerInterface, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	return &DefaultShellManager{
		pkgManager: pkgManager,
		detector:   detector,
		config:     cfg,
		homeDir:    homeDir,
	}, nil
}

// Install installs a shell
func (m *DefaultShellManager) Install(shell string) error {
	// Check if shell is already installed
	if m.IsShellAvailable(Type(shell)) {
		return nil
	}

	// Get shell package name
	shellInfo, exists := shells[shell]
	if !exists {
		return fmt.Errorf("unknown shell: %s", shell)
	}

	// Install shell using package manager
	if err := m.pkgManager.Install(shellInfo.Package); err != nil {
		return fmt.Errorf("failed to install shell %s: %w", shell, err)
	}

	// Add shell to /etc/shells if not already present
	shellPath := ""
	paths := []string{
		"/bin/" + shell,
		"/usr/bin/" + shell,
		"/usr/local/bin/" + shell,
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			shellPath = path
			break
		}
	}
	if shellPath != "" {
		if err := addToShells(shellPath); err != nil {
			return fmt.Errorf("failed to add shell to /etc/shells: %w", err)
		}
	}

	return nil
}

// InstallPrompt installs a shell prompt
func (m *DefaultShellManager) InstallPrompt(prompt, shell string) error {
	// Get prompt information
	promptInfo, exists := prompts[prompt]
	if !exists {
		return fmt.Errorf("prompt %q not found", prompt)
	}

	// Check shell compatibility
	if !promptInfo.SupportsShell(shell) {
		return fmt.Errorf("prompt %q is not compatible with shell %q", prompt, shell)
	}

	// Install the prompt
	if err := promptInfo.Install(shell); err != nil {
		return fmt.Errorf("failed to install prompt %q: %w", prompt, err)
	}

	// Configure the prompt
	if err := promptInfo.Configure(shell); err != nil {
		return fmt.Errorf("failed to configure prompt %q: %w", prompt, err)
	}

	return nil
}

// SetDefault sets the default shell
func (m *DefaultShellManager) SetDefault(shell string) error {
	// Get shell path
	shellPath := ""
	paths := []string{
		"/bin/" + shell,
		"/usr/bin/" + shell,
		"/usr/local/bin/" + shell,
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			shellPath = path
			break
		}
	}
	if shellPath == "" {
		return fmt.Errorf("shell not found: %s", shell)
	}

	// Add shell to /etc/shells if not already present
	shellsFile := "/etc/shells"
	content, err := os.ReadFile(shellsFile)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", shellsFile, err)
	}
	if !strings.Contains(string(content), shellPath) {
		cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo %s >> %s", shellPath, shellsFile))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add shell to %s: %w", shellsFile, err)
		}
	}

	// Change default shell
	cmd := exec.Command("chsh", "-s", shellPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set default shell: %w", err)
	}

	return nil
}

// GetCurrent returns the current shell
func (m *DefaultShellManager) GetCurrent() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable not set")
	}
	return filepath.Base(shell), nil
}

// GenerateConfig generates shell configuration
func (m *DefaultShellManager) GenerateConfig(shell Type, data *shellconfig.Data) error {
	if data == nil {
		return fmt.Errorf("config data is nil")
	}

	// Create a copy of data for template execution
	templateData := *data
	templateData.Path = append(templateData.Path, "$PATH")

	// Generate config based on shell type
	switch shell {
	case Zsh:
		return shellconfig.GenerateZsh(&templateData)
	case Bash:
		return shellconfig.GenerateBash(&templateData)
	case Fish:
		return shellconfig.GenerateFish(&templateData)
	default:
		return fmt.Errorf("unsupported shell type: %s", shell)
	}
}

// RestoreConfig restores shell configuration from a backup
func (m *DefaultShellManager) RestoreConfig(shell Type, path string) error {
	if path == "" {
		return fmt.Errorf("backup path is empty")
	}

	// Get home directory for config data
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	data := shellconfig.New(homeDir)

	// Restore config based on shell type
	switch shell {
	case Zsh:
		return shellconfig.RestoreZsh(data, path)
	case Bash:
		return shellconfig.RestoreBash(data, path)
	case Fish:
		return shellconfig.RestoreFish(data, path)
	default:
		return fmt.Errorf("unsupported shell type: %s", shell)
	}
}

// IsSupportedPluginManager checks if a plugin manager is supported
func IsSupportedPluginManager(pm PluginManager) bool {
	info, exists := pluginManagers[string(pm)]
	return exists && len(info.CompatibleShells) > 0
}

// Helper functions

func addToShells(shellPath string) error {
	// Check if shell is already in /etc/shells
	shellsFile := "/etc/shells"
	content, err := os.ReadFile(shellsFile)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", shellsFile, err)
	}

	if !strings.Contains(string(content), shellPath) {
		// Add shell to /etc/shells
		cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo %s >> %s", shellPath, shellsFile))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add shell to %s: %w", shellsFile, err)
		}
	}

	return nil
}

func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// IsSupportedShell checks if a shell type is supported
func IsSupportedShell(shell Type) bool {
	switch shell {
	case Bash, Zsh, Fish:
		return true
	default:
		return false
	}
}

// GetConfigPath returns the path to the shell configuration file
func GetConfigPath(shell Type) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch shell {
	case Zsh:
		return filepath.Join(homeDir, ".zshrc")
	case Bash:
		return filepath.Join(homeDir, ".bashrc")
	case Fish:
		return filepath.Join(homeDir, ".config", "fish", "config.fish")
	default:
		return ""
	}
}

// Setup configures the shell environment
func (m *DefaultShellManager) Setup(cfg *Config) error {
	// Install shell if not present
	if !m.IsShellAvailable(cfg.Type) {
		// Install shell package
		shellPkg := getShellPackage(cfg.Type)
		if shellPkg == "" {
			return fmt.Errorf("unsupported shell type: %s", cfg.Type)
		}

		// Install shell using package manager
		if err := m.pkgManager.Install(shellPkg); err != nil {
			return fmt.Errorf("failed to install shell: %w", err)
		}
	}

	// Install plugin manager if configured
	if cfg.PluginMgr != "" {
		if err := m.InstallPluginManager(cfg.PluginMgr, string(cfg.Type)); err != nil {
			return fmt.Errorf("failed to install plugin manager: %w", err)
		}
	}

	// Install plugins
	for _, plugin := range cfg.Plugins {
		if err := m.InstallPlugin(plugin); err != nil {
			return fmt.Errorf("failed to install plugin %s: %w", plugin, err)
		}
	}

	return nil
}

// IsShellAvailable checks if a shell is available on the system
func (m *DefaultShellManager) IsShellAvailable(shell Type) bool {
	paths := []string{
		"/bin/" + string(shell),
		"/usr/bin/" + string(shell),
		"/usr/local/bin/" + string(shell),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// InstallPluginManager installs the specified plugin manager
func (m *DefaultShellManager) InstallPluginManager(pm PluginManager, shell string) error {
	// Get plugin manager info
	info, exists := pluginManagers[string(pm)]
	if !exists {
		return &ErrPluginManagerNotSupported{
			Manager: pm,
			Message: "plugin manager not found",
		}
	}

	// Check shell compatibility
	if !info.SupportsShell(shell) {
		return &ErrPluginManagerNotSupported{
			Manager: pm,
			Message: fmt.Sprintf("not compatible with shell %s", shell),
		}
	}

	// Install using the plugin manager's install function
	if err := info.Install(shell); err != nil {
		return &ErrPluginManagerInstallation{
			Manager: pm,
			Err:     err,
		}
	}

	return nil
}

// InstallPlugin installs a shell plugin
func (m *DefaultShellManager) InstallPlugin(plugin string) error {
	if !m.isValidPlugin(plugin) {
		return &ErrPluginNotFound{Plugin: plugin}
	}

	m.config.Plugins = append(m.config.Plugins, plugin)
	return nil
}

// UninstallPlugin removes a shell plugin
func (m *DefaultShellManager) UninstallPlugin(plugin string) error {
	for i, p := range m.config.Plugins {
		if p == plugin {
			m.config.Plugins = append(m.config.Plugins[:i], m.config.Plugins[i+1:]...)
			return nil
		}
	}
	return &ErrPluginNotFound{Plugin: plugin}
}

// ListAvailablePlugins returns a list of available plugins
func (m *DefaultShellManager) ListAvailablePlugins() []string {
	switch m.config.PluginMgr {
	case OhMyZsh:
		return []string{
			"git",
			"docker",
			"kubectl",
			"golang",
			"node",
			"python",
			"rust",
		}
	case Antigen:
		return []string{
			"zsh-users/zsh-autosuggestions",
			"zsh-users/zsh-syntax-highlighting",
		}
	case Zinit:
		return []string{
			"zdharma-continuum/fast-syntax-highlighting",
			"zsh-users/zsh-autosuggestions",
		}
	case Fisherman:
		return []string{
			"pure",
			"z",
			"nvm",
		}
	default:
		return []string{}
	}
}

// BackupConfig creates a backup of the existing shell configuration
func (m *DefaultShellManager) BackupConfig() error {
	configPath := m.getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		backupPath := configPath + ".backup"
		if err := os.Rename(configPath, backupPath); err != nil {
			return &ErrBackup{Err: err}
		}
	}
	return nil
}

// Cleanup performs cleanup operations
func (m *DefaultShellManager) Cleanup() error {
	// Remove temporary files, etc.
	return nil
}

// Helper methods

func (m *DefaultShellManager) getConfigPath() string {
	switch m.config.Type {
	case Zsh:
		return filepath.Join(m.homeDir, ".zshrc")
	case Fish:
		return filepath.Join(m.homeDir, ".config", "fish", "config.fish")
	default:
		return filepath.Join(m.homeDir, ".bashrc")
	}
}

func (m *DefaultShellManager) isValidPlugin(plugin string) bool {
	available := m.ListAvailablePlugins()
	for _, p := range available {
		if p == plugin {
			return true
		}
	}
	return false
}

// AddAlias adds a shell alias
func (m *DefaultShellManager) AddAlias(name, command string) {
	if m.config.Aliases == nil {
		m.config.Aliases = make(map[string]string)
	}
	m.config.Aliases[name] = command
}

// AddEnvVar adds an environment variable
func (m *DefaultShellManager) AddEnvVar(name, value string) {
	if m.config.EnvVars == nil {
		m.config.EnvVars = make(map[string]string)
	}
	m.config.EnvVars[name] = value
}

// AddPath adds a path to PATH
func (m *DefaultShellManager) AddPath(path string) {
	if m.config.Path == nil {
		m.config.Path = make([]string, 0)
	}
	m.config.Path = append(m.config.Path, path)
}

// getShellPackage returns the package name for a shell type
func getShellPackage(shellType Type) string {
	info, exists := shells[string(shellType)]
	if !exists {
		return ""
	}
	return info.Package
}