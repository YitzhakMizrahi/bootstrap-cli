package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

var (
	selectedTools []*interfaces.Tool
	selectedToolsMutex sync.RWMutex
)

// SetSelectedTools sets the selected tools for installation
func SetSelectedTools(tools []*interfaces.Tool) {
	selectedToolsMutex.Lock()
	defer selectedToolsMutex.Unlock()
	selectedTools = tools
}

// GetSelectedTools gets the selected tools for installation
func GetSelectedTools() []*interfaces.Tool {
	selectedToolsMutex.RLock()
	defer selectedToolsMutex.RUnlock()
	return selectedTools
}

// PackageMapping defines package names for different package managers
type PackageMapping struct {
	Default string
	APT     string
	DNF     string
	Pacman  string
	Brew    string
}

// Tool represents a development tool that can be installed
type Tool struct {
	// Name is the name of the tool
	Name string
	// PackageName is the name of the package in the package manager
	PackageName string
	// PackageNames contains system-specific package names
	PackageNames *PackageMapping
	// Version is the desired version of the tool
	Version string
	// Dependencies is a list of package names that this tool depends on
	Dependencies []string
	// SystemDependencies is a list of system-level dependencies required
	SystemDependencies []string
	// PostInstall is a list of commands to run after installation
	PostInstall []Command
	// VerifyCommand is the command to verify the installation
	VerifyCommand string
	// Description is a brief description of the tool
	Description string
	// Category is the category of the tool (e.g., "Essential", "Modern CLI", "System")
	Category string
	// Tags are keywords for searching and filtering
	Tags []string
	// ConfigFiles is a list of configuration files to create/symlink
	ConfigFiles []ConfigFile
	// ShellConfig is shell-specific configuration to add
	ShellConfig interfaces.ShellConfig
	// RequiresRestart indicates if a shell restart is needed after installation
	RequiresRestart bool
}

// ConfigFile represents a configuration file to be created or symlinked
type ConfigFile struct {
	// Source is the source path or content
	Source string
	// Destination is the target path
	Destination string
	// Type is the type of configuration (file, symlink, or content)
	Type string
	// Mode is the file mode (e.g., "0644")
	Mode string
}

// Command represents a post-installation command
type Command struct {
	Command     string
	Description string
}

// Error represents an installation error
type Error struct {
	Tool    string
	Phase   string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s failed: %s (%v)", e.Tool, e.Phase, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s failed: %s", e.Tool, e.Phase, e.Message)
}

// Result represents the result of an installation attempt
type Result struct {
	Success      bool
	InstalledPkg string
	Error        error
}

// Installer handles tool installation
type Installer struct {
	// PackageManager is the package manager to use
	PackageManager interfaces.PackageManager
	// Logger is the logger to use
	Logger *log.Logger
	// MaxRetries is the maximum number of retries for failed operations
	MaxRetries int
	// RetryDelay is the delay between retries
	RetryDelay time.Duration
}

// NewInstaller creates a new installer with the given package manager
func NewInstaller(pm interfaces.PackageManager) *Installer {
	return &Installer{
		PackageManager: pm,
		Logger:        log.New(log.InfoLevel),
		MaxRetries:    3,
		RetryDelay:    time.Second * 2,
	}
}

// getPackageWithVersion returns the package name with version if specified
func (i *Installer) getPackageWithVersion(pkg, version string) string {
	if version == "" || version == "latest" || version == "stable" {
		return pkg
	}
	
	switch i.PackageManager.GetName() {
	case "apt":
		return fmt.Sprintf("%s=%s", pkg, version)
	case "dnf":
		return fmt.Sprintf("%s-%s", pkg, version)
	case "pacman":
		return fmt.Sprintf("%s=%s", pkg, version)
	case "brew":
		return fmt.Sprintf("%s@%s", pkg, version)
	default:
		return pkg
	}
}

// getSystemPackageName returns the appropriate package name for the current system
func (i *Installer) getSystemPackageName(tool *interfaces.Tool) string {
	if tool == nil {
		return ""
	}

	// Try to get system-specific package name
	if i.PackageManager != nil {
		switch i.PackageManager.GetName() {
		case "apt":
			if tool.PackageNames.APT != "" {
				return tool.PackageNames.APT
			}
		case "dnf":
			if tool.PackageNames.DNF != "" {
				return tool.PackageNames.DNF
			}
		case "pacman":
			if tool.PackageNames.Pacman != "" {
				return tool.PackageNames.Pacman
			}
		case "brew":
			if tool.PackageNames.Brew != "" {
				return tool.PackageNames.Brew
			}
		}
	}

	// Fall back to default package name
	return tool.Name
}

// Install installs a tool
func (i *Installer) Install(tool *interfaces.Tool) error {
	if tool == nil {
		return fmt.Errorf("tool is nil")
	}

	i.Logger.Info("Installing %s...", tool.Name)

	// Get the appropriate package name for the current system
	pkgName := i.getSystemPackageName(tool)
	if pkgName == "" {
		return fmt.Errorf("no package name found for tool %s", tool.Name)
	}

	// Add version if specified
	pkgName = i.getPackageWithVersion(pkgName, tool.Version)

	// Install system dependencies first
	if len(tool.SystemDependencies) > 0 {
		i.Logger.Info("Installing system dependencies for %s...", tool.Name)
		for _, dep := range tool.SystemDependencies {
			err := i.PackageManager.Install(dep)
			if err != nil {
				return fmt.Errorf("failed to install system dependency %s: %v", dep, err)
			}
		}
	}

	// Install dependencies
	if len(tool.Dependencies) > 0 {
		i.Logger.Info("Installing dependencies for %s...", tool.Name)
		for _, dep := range tool.Dependencies {
			err := i.PackageManager.Install(dep.Name)
			if err != nil && !dep.Optional {
				return fmt.Errorf("failed to install dependency %s: %v", dep.Name, err)
			}
		}
	}

	// Install the tool
	err := i.PackageManager.Install(pkgName)
	if err != nil {
		return fmt.Errorf("failed to install %s: %v", tool.Name, err)
	}

	// Run post-install commands
	if len(tool.PostInstall) > 0 {
		i.Logger.Info("Running post-install commands for %s...", tool.Name)
		for _, cmd := range tool.PostInstall {
			if err := i.runCommand(cmd.Command); err != nil {
				return fmt.Errorf("post-install command failed: %v", err)
			}
		}
	}

	// Verify installation
	if tool.VerifyCommand != "" {
		i.Logger.Info("Verifying installation of %s...", tool.Name)
		if err := i.verifyInstallation(tool); err != nil {
			return fmt.Errorf("verification failed: %v", err)
		}
	}

	// Set up config files
	if len(tool.ConfigFiles) > 0 {
		i.Logger.Info("Setting up configuration files for %s...", tool.Name)
		if err := i.setupConfigFiles(tool); err != nil {
			return fmt.Errorf("failed to set up config files: %v", err)
		}
	}

	// Apply shell configuration
	if err := i.applyShellConfig(tool); err != nil {
		return fmt.Errorf("failed to apply shell configuration: %v", err)
	}

	i.Logger.Success("Successfully installed %s", tool.Name)
	return nil
}

// Options represents options for installing tools
type Options struct {
	Logger           *log.Logger
	PackageManager   interfaces.PackageManager
	Tools            []*interfaces.Tool
	SkipVerification bool
	AdditionalPaths  []string
}

// CoreTools installs core tools
func CoreTools(opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options are nil")
	}

	installer := &Installer{
		PackageManager: opts.PackageManager,
		Logger:        opts.Logger,
	}

	for _, tool := range opts.Tools {
		if err := installer.Install(tool); err != nil {
			return fmt.Errorf("failed to install %s: %v", tool.Name, err)
		}
	}

	if !opts.SkipVerification {
		return VerifyCoreTools(opts)
	}

	return nil
}

// VerifyCoreTools verifies core tools are installed correctly
func VerifyCoreTools(opts *Options) error {
	if opts == nil {
		return fmt.Errorf("options are nil")
	}

	installer := &Installer{
		PackageManager: opts.PackageManager,
		Logger:        opts.Logger,
	}

	for _, tool := range opts.Tools {
		if tool.VerifyCommand != "" {
			if err := installer.verifyInstallation(tool); err != nil {
				return fmt.Errorf("verification failed for %s: %v", tool.Name, err)
			}
		}
	}

	return nil
}

// Helper functions

func (i *Installer) installWithRetry(pkg string) error {
	return i.retryOperation(func() error {
		return i.PackageManager.Install(pkg)
	})
}

func (i *Installer) verifyInstallation(tool *interfaces.Tool) error {
	return i.retryOperation(func() error {
		cmd := exec.Command("sh", "-c", tool.VerifyCommand)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("verification command failed: %v", err)
		}
		return nil
	})
}

func (i *Installer) retryOperation(op func() error) error {
	var lastErr error
	for attempt := 0; attempt < i.MaxRetries; attempt++ {
		if err := op(); err != nil {
			lastErr = err
			i.Logger.Warn("Operation failed (attempt %d/%d): %v", attempt+1, i.MaxRetries, err)
			if attempt < i.MaxRetries-1 {
				time.Sleep(i.RetryDelay)
				continue
			}
		} else {
			return nil
		}
	}
	return lastErr
}

func (i *Installer) runCommand(cmd string) error {
	command := exec.Command("sh", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func (i *Installer) setupConfigFiles(tool *interfaces.Tool) error {
	for _, file := range tool.ConfigFiles {
		if err := i.createFile(file.Source, file.Destination, file.Mode); err != nil {
			return fmt.Errorf("failed to set up config files: %v", err)
		}
	}
	return nil
}

func (i *Installer) applyShellConfig(tool *interfaces.Tool) error {
	if tool.ShellConfig.Aliases == nil && tool.ShellConfig.Env == nil && len(tool.ShellConfig.Path) == 0 {
		return nil
	}

	shell, err := i.getCurrentShell()
	if err != nil {
		return fmt.Errorf("failed to detect current shell: %v", err)
	}

	switch {
	case strings.Contains(shell, "zsh"):
		return i.applyZshConfig(tool)
	case strings.Contains(shell, "bash"):
		return i.applyBashConfig(tool)
	case strings.Contains(shell, "fish"):
		return i.applyFishConfig(tool)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

func (i *Installer) createFile(source, destination string, mode string) error {
	// Create parent directories if they don't exist
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create parent directories: %v", err)
	}

	// Read source file
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	// Parse mode string into os.FileMode
	var fileMode os.FileMode = 0644 // Default mode
	if mode != "" {
		parsed, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid file mode %s: %v", mode, err)
		}
		fileMode = os.FileMode(parsed)
	}

	// Write to destination with specified mode
	if err := os.WriteFile(destination, content, fileMode); err != nil {
		return fmt.Errorf("failed to write destination file: %v", err)
	}

	return nil
}

func (i *Installer) getCurrentShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable not set")
	}
	return shell, nil
}

func (i *Installer) applyZshConfig(tool *interfaces.Tool) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".zsh")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create zsh config directory: %v", err)
	}

	configFile := filepath.Join(configDir, fmt.Sprintf("%s.zsh", tool.Name))
	var config strings.Builder

	// Add aliases
	for alias, cmd := range tool.ShellConfig.Aliases {
		config.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, cmd))
	}

	// Add environment variables
	for key, value := range tool.ShellConfig.Env {
		config.WriteString(fmt.Sprintf("export %s='%s'\n", key, value))
	}

	// Add PATH entries
	for _, path := range tool.ShellConfig.Path {
		config.WriteString(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", path))
	}

	// Write the config file
	if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write zsh config: %v", err)
	}

	// Add source line to .zshrc if not already present
	zshrc := filepath.Join(os.Getenv("HOME"), ".zshrc")
	sourceLine := fmt.Sprintf("source %s", configFile)
	
	content, err := os.ReadFile(zshrc)
	if err != nil {
		return fmt.Errorf("failed to read .zshrc: %v", err)
	}

	if !strings.Contains(string(content), sourceLine) {
		f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .zshrc: %v", err)
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("\n# Added by bootstrap-cli\n%s\n", sourceLine)); err != nil {
			return fmt.Errorf("failed to update .zshrc: %v", err)
		}
	}

	return nil
}

func (i *Installer) applyBashConfig(tool *interfaces.Tool) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".bash")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create bash config directory: %v", err)
	}

	configFile := filepath.Join(configDir, fmt.Sprintf("%s.bash", tool.Name))
	var config strings.Builder

	// Add aliases
	for alias, cmd := range tool.ShellConfig.Aliases {
		config.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, cmd))
	}

	// Add environment variables
	for key, value := range tool.ShellConfig.Env {
		config.WriteString(fmt.Sprintf("export %s='%s'\n", key, value))
	}

	// Add PATH entries
	for _, path := range tool.ShellConfig.Path {
		config.WriteString(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", path))
	}

	// Write the config file
	if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write bash config: %v", err)
	}

	// Add source line to .bashrc if not already present
	bashrc := filepath.Join(os.Getenv("HOME"), ".bashrc")
	sourceLine := fmt.Sprintf("source %s", configFile)
	
	content, err := os.ReadFile(bashrc)
	if err != nil {
		return fmt.Errorf("failed to read .bashrc: %v", err)
	}

	if !strings.Contains(string(content), sourceLine) {
		f, err := os.OpenFile(bashrc, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .bashrc: %v", err)
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("\n# Added by bootstrap-cli\n%s\n", sourceLine)); err != nil {
			return fmt.Errorf("failed to update .bashrc: %v", err)
		}
	}

	return nil
}

func (i *Installer) applyFishConfig(tool *interfaces.Tool) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".config/fish/conf.d")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create fish config directory: %v", err)
	}

	configFile := filepath.Join(configDir, fmt.Sprintf("%s.fish", tool.Name))
	var config strings.Builder

	// Add aliases
	for alias, cmd := range tool.ShellConfig.Aliases {
		config.WriteString(fmt.Sprintf("alias %s '%s'\n", alias, cmd))
	}

	// Add environment variables
	for key, value := range tool.ShellConfig.Env {
		config.WriteString(fmt.Sprintf("set -gx %s '%s'\n", key, value))
	}

	// Add PATH entries
	for _, path := range tool.ShellConfig.Path {
		config.WriteString(fmt.Sprintf("fish_add_path '%s'\n", path))
	}

	// Write the config file
	if err := os.WriteFile(configFile, []byte(config.String()), 0644); err != nil {
		return fmt.Errorf("failed to write fish config: %v", err)
	}

	return nil
} 