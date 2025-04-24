package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/implementations"
)

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
	PostInstall []PostInstallCommand
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
	// Permissions are the file permissions (e.g., "0644")
	Permissions string
}

// PostInstallCommand represents a post-installation command
type PostInstallCommand struct {
	Command     string
	Description string
}

// InstallError represents an installation error
type InstallError struct {
	Tool    string
	Phase   string
	Message string
	Err     error
}

func (e *InstallError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s failed: %s (%v)", e.Tool, e.Phase, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s failed: %s", e.Tool, e.Phase, e.Message)
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
func (i *Installer) getSystemPackageName(tool *Tool) string {
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
	return tool.PackageName
}

// InstallResult represents the result of an installation attempt
type InstallResult struct {
	Success      bool
	InstalledPkg string
	Error        error
}

// Install installs a tool
func (i *Installer) Install(tool *Tool) error {
	i.Logger.Info("Starting installation of %s", tool.Name)

	// Get the appropriate package name for the current system
	pkg := i.getSystemPackageName(tool)
	if pkg == "" {
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "package name resolution",
			Message: "no suitable package name found for current system",
		}
	}

	// Add version if specified
	pkgWithVersion := i.getPackageWithVersion(pkg, tool.Version)
	i.Logger.Info("Installing %s version %s", pkg, tool.Version)

	// Install system dependencies first
	if len(tool.SystemDependencies) > 0 {
		i.Logger.Info("Installing system dependencies for %s", tool.Name)
		for _, dep := range tool.SystemDependencies {
			if err := i.installWithRetry(dep).Error; err != nil {
				return &InstallError{
					Tool:    tool.Name,
					Phase:   "system dependency installation",
					Message: fmt.Sprintf("failed to install system dependency %s", dep),
					Err:     err,
				}
			}
		}
	}

	// Install package dependencies
	if len(tool.Dependencies) > 0 {
		i.Logger.Info("Installing dependencies for %s", tool.Name)
		for _, dep := range tool.Dependencies {
			if err := i.installWithRetry(dep).Error; err != nil {
				return &InstallError{
					Tool:    tool.Name,
					Phase:   "dependency installation",
					Message: fmt.Sprintf("failed to install dependency %s", dep),
					Err:     err,
				}
			}
		}
	}

	// Install the main package
	result := i.installWithRetry(pkgWithVersion)
	if result.Error != nil {
		// Clean up dependencies if main package installation failed
		i.cleanup(tool.Dependencies)
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "package installation",
			Message: fmt.Sprintf("failed to install package %s", pkgWithVersion),
			Err:     result.Error,
		}
	}

	// Create configuration files
	if len(tool.ConfigFiles) > 0 {
		i.Logger.Info("Setting up configuration for %s", tool.Name)
		if err := i.setupConfigFiles(tool); err != nil {
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "configuration setup",
				Message: "failed to set up configuration files",
				Err:     err,
			}
		}
	}

	// Run post-install commands
	if len(tool.PostInstall) > 0 {
		i.Logger.Info("Running post-install commands for %s", tool.Name)
		for _, cmd := range tool.PostInstall {
			if err := i.runPostInstall(cmd); err != nil {
				return &InstallError{
					Tool:    tool.Name,
					Phase:   "post-install",
					Message: fmt.Sprintf("failed to run post-install command: %s", cmd.Command),
					Err:     err,
				}
			}
		}
	}

	// Apply shell configuration
	i.Logger.Info("Applying shell configuration for %s", tool.Name)
	if err := i.applyShellConfig(tool); err != nil {
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "shell configuration",
			Message: "failed to apply shell configuration",
			Err:     err,
		}
	}

	// Reload shell configuration
	i.Logger.Info("Reloading shell configuration")
	shell, err := i.getCurrentShell()
	if err == nil {
		reloadCmd := PostInstallCommand{
			Command: fmt.Sprintf("source ~/.%src", shell),
			Description: "Reload shell configuration",
		}
		if shell == "fish" {
			reloadCmd.Command = "source ~/.config/fish/config.fish"
		}
		if err := i.runPostInstall(reloadCmd); err != nil {
			i.Logger.Warn("Failed to reload shell configuration: %v", err)
		}
	}

	// Verify installation
	i.Logger.Info("Verifying installation of %s", tool.Name)
	if err := i.verifyInstallation(tool); err != nil {
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "verification",
			Message: "installation verification failed",
			Err:     err,
		}
	}

	i.Logger.Success("Successfully installed %s", tool.Name)
	return nil
}

// installWithRetry attempts to install a package with retries
func (i *Installer) installWithRetry(pkg string) InstallResult {
	// Check if package is already installed
	if i.PackageManager.IsInstalled(pkg) {
		i.Logger.Info("%s is already installed", pkg)
		return InstallResult{
			Success:      true,
			InstalledPkg: pkg,
			Error:        nil,
		}
	}

	var lastErr error
	var repositorySetupAttempted bool

	for attempt := 0; attempt <= i.MaxRetries; attempt++ {
		if attempt > 0 {
			i.Logger.Debug("Retrying installation of %s (attempt %d/%d)", pkg, attempt, i.MaxRetries)
			time.Sleep(i.RetryDelay)

			// Update package lists before retry
			if err := i.PackageManager.Update(); err != nil {
				i.Logger.Debug("Failed to update package lists: %v", err)
			}
		}

		// Try to install the package
		if err := i.PackageManager.Install(pkg); err != nil {
			lastErr = err
			
			// Check if the error indicates package not found
			if strings.Contains(err.Error(), "Unable to locate package") ||
				strings.Contains(err.Error(), "No matching package") ||
				strings.Contains(err.Error(), "package not found") {
				i.Logger.Info("Package %s not found in repositories", pkg)
				return InstallResult{
					Success:      false,
					InstalledPkg: pkg,
					Error:        fmt.Errorf("package not found: %w", err),
				}
			}
			
			// If this is a special package that needs repository setup and we haven't tried that yet
			if !repositorySetupAttempted {
				// Try to set up repository for the package
				if setupErr := i.setupRepository(pkg); setupErr != nil {
					i.Logger.Debug("Failed to set up repository for %s: %v", pkg, setupErr)
					// Don't retry if repository setup failed
					return InstallResult{
						Success:      false,
						InstalledPkg: pkg,
						Error:        fmt.Errorf("failed to set up repository: %w", setupErr),
					}
				}
				repositorySetupAttempted = true
				// Continue to next attempt after repository setup
				continue
			}
			
			i.Logger.Debug("Installation attempt failed: %v", err)
			continue
		}

		return InstallResult{
			Success:      true,
			InstalledPkg: pkg,
			Error:        nil,
		}
	}

	return InstallResult{
		Success:      false,
		InstalledPkg: pkg,
		Error:        fmt.Errorf("failed to install after %d attempts: %w", i.MaxRetries, lastErr),
	}
}

// setupRepository attempts to set up the repository for a package
func (i *Installer) setupRepository(pkg string) error {
	switch pm := i.PackageManager.(type) {
	case *implementations.APTManager:
		return pm.SetupSpecialPackage(pkg)
	case *implementations.DnfPackageManager:
		return pm.SetupSpecialPackage(pkg)
	case *implementations.PacmanPackageManager:
		return pm.SetupSpecialPackage(pkg)
	default:
		return nil // No repository setup needed for this package manager
	}
}

// cleanup removes installed packages in case of failure
func (i *Installer) cleanup(packages []string) {
	if len(packages) == 0 {
		return
	}

	i.Logger.Info("Cleaning up installed packages: %v", packages)
	for _, pkg := range packages {
		i.Logger.Debug("Removing package: %s", pkg)
		if err := i.PackageManager.Remove(pkg); err != nil {
			i.Logger.Error("Failed to remove package %s: %v", pkg, err)
		}
	}
}

// runPostInstall executes post-installation commands
func (i *Installer) runPostInstall(cmd PostInstallCommand) error {
	if cmd.Command == "" {
		return nil
	}

	command := exec.Command("sh", "-c", cmd.Command)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("post-install command failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// verifyInstallation checks if the tool was installed correctly
func (i *Installer) verifyInstallation(tool *Tool) error {
	if tool.VerifyCommand == "" {
		return nil
	}

	// Wait for PATH to be updated
	time.Sleep(time.Second * 2)

	// Try to find the binary in PATH
	cmd := exec.Command("sh", "-c", tool.VerifyCommand)
	if output, err := cmd.CombinedOutput(); err != nil {
		// If not found in PATH, check common locations
		commonPaths := []string{
			"/usr/bin/",
			"/usr/local/bin/",
			"/opt/homebrew/bin/",
		}

		for _, path := range commonPaths {
			fullPath := filepath.Join(path, tool.Name)
			if _, err := os.Stat(fullPath); err == nil {
				return nil
			}
		}

		return fmt.Errorf("verification failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// retryOperation retries an operation with exponential backoff
func (i *Installer) retryOperation(op func() error) error {
	var lastErr error
	for attempt := 0; attempt < i.MaxRetries; attempt++ {
		if attempt > 0 {
			i.Logger.Debug("Retrying operation (attempt %d/%d)", attempt+1, i.MaxRetries)
			time.Sleep(i.RetryDelay * time.Duration(attempt))
		}

		if err := op(); err != nil {
			lastErr = err
			i.Logger.Warn("Operation failed (attempt %d/%d): %v", attempt+1, i.MaxRetries, err)
			continue
		}

		return nil
	}

	return lastErr
}

// runCommand executes a shell command
func (i *Installer) runCommand(cmd string) error {
	if cmd == "" {
		return nil
	}

	return i.runPostInstall(PostInstallCommand{
		Command: cmd,
		Description: "Run shell command",
	})
}

// setupConfigFiles creates or symlinks configuration files
func (i *Installer) setupConfigFiles(tool *Tool) error {
	for _, config := range tool.ConfigFiles {
		switch config.Type {
		case "file":
			if err := i.createFile(config.Source, config.Destination, config.Permissions); err != nil {
				return fmt.Errorf("failed to create file %s: %w", config.Destination, err)
			}
		case "symlink":
			if err := i.createSymlink(config.Source, config.Destination); err != nil {
				return fmt.Errorf("failed to create symlink %s: %w", config.Destination, err)
			}
		case "content":
			if err := i.writeContent(config.Source, config.Destination, config.Permissions); err != nil {
				return fmt.Errorf("failed to write content to %s: %w", config.Destination, err)
			}
		default:
			return fmt.Errorf("unknown config type: %s", config.Type)
		}
	}
	return nil
}

// applyShellConfig applies shell-specific configuration
func (i *Installer) applyShellConfig(tool *Tool) error {
	if len(tool.ShellConfig.Exports) == 0 && len(tool.ShellConfig.Path) == 0 {
		return nil
	}

	shell, err := i.getCurrentShell()
	if err != nil {
		return err
	}

	switch shell {
	case "zsh":
		return i.applyZshConfig(tool)
	case "bash":
		return i.applyBashConfig(tool)
	case "fish":
		return i.applyFishConfig(tool)
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
}

// createFile creates a new file with the given content and permissions
func (i *Installer) createFile(source, destination, permissions string) error {
	// Read source file
	content, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Parse permissions
	perm, err := strconv.ParseInt(permissions, 8, 32)
	if err != nil {
		return fmt.Errorf("invalid permissions: %w", err)
	}

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(destination, content, os.FileMode(perm)); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// createSymlink creates a symbolic link
func (i *Installer) createSymlink(source, destination string) error {
	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Remove existing symlink or file
	if _, err := os.Lstat(destination); err == nil {
		if err := os.Remove(destination); err != nil {
			return fmt.Errorf("failed to remove existing file: %w", err)
		}
	}

	// Create symlink
	if err := os.Symlink(source, destination); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// writeContent writes content to a file
func (i *Installer) writeContent(content, destination, permissions string) error {
	// Parse permissions
	perm, err := strconv.ParseInt(permissions, 8, 32)
	if err != nil {
		return fmt.Errorf("invalid permissions: %w", err)
	}

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(destination, []byte(content), os.FileMode(perm)); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// getCurrentShell returns the current shell name (zsh, bash, fish)
func (i *Installer) getCurrentShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable not set")
	}

	switch {
	case strings.Contains(shell, "zsh"):
		return "zsh", nil
	case strings.Contains(shell, "bash"):
		return "bash", nil
	case strings.Contains(shell, "fish"):
		return "fish", nil
	default:
		return "", fmt.Errorf("unsupported shell: %s", shell)
	}
}

// applyZshConfig applies Zsh-specific configuration
func (i *Installer) applyZshConfig(tool *Tool) error {
	config := tool.ShellConfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	zshrc := filepath.Join(homeDir, ".zshrc")
	content := []string{}

	// Add environment variables
	if len(config.Exports) > 0 {
		content = append(content, "\n# Environment variables")
		for key, value := range config.Exports {
			content = append(content, fmt.Sprintf("export %s=%s", key, value))
		}
	}

	// Add PATH additions
	if len(config.Path) > 0 {
		content = append(content, "\n# PATH additions")
		for _, path := range config.Path {
			content = append(content, fmt.Sprintf("export PATH=%s:$PATH", path))
		}
	}

	// Read existing .zshrc
	existingContent, err := os.ReadFile(zshrc)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .zshrc: %w", err)
	}

	// Prepare new content
	var newContent strings.Builder
	if len(existingContent) > 0 {
		newContent.Write(existingContent)
		newContent.WriteString("\n\n")
	}

	// Add configuration header
	newContent.WriteString(fmt.Sprintf("# Configuration for %s\n", tool.Name))

	// Add aliases
	if len(config.Exports) > 0 {
		newContent.WriteString("\n# Aliases\n")
		for alias, cmd := range config.Exports {
			newContent.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, cmd))
		}
	}

	// Add functions
	if len(config.Functions) > 0 {
		newContent.WriteString("\n# Functions\n")
		for name, body := range config.Functions {
			newContent.WriteString(fmt.Sprintf("%s() {\n%s\n}\n", name, body))
		}
	}

	// Write updated .zshrc
	if err := os.WriteFile(zshrc, []byte(newContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write .zshrc: %w", err)
	}

	return nil
}

// applyBashConfig applies Bash-specific configuration
func (i *Installer) applyBashConfig(tool *Tool) error {
	config := tool.ShellConfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	bashrc := filepath.Join(homeDir, ".bashrc")
	content := []string{}

	// Add environment variables
	if len(config.Exports) > 0 {
		content = append(content, "\n# Environment variables")
		for key, value := range config.Exports {
			content = append(content, fmt.Sprintf("export %s=%s", key, value))
		}
	}

	// Add PATH additions
	if len(config.Path) > 0 {
		content = append(content, "\n# PATH additions")
		for _, path := range config.Path {
			content = append(content, fmt.Sprintf("export PATH=%s:$PATH", path))
		}
	}

	// Read existing .bashrc
	existingContent, err := os.ReadFile(bashrc)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .bashrc: %w", err)
	}

	// Prepare new content
	var newContent strings.Builder
	if len(existingContent) > 0 {
		newContent.Write(existingContent)
		newContent.WriteString("\n\n")
	}

	// Add configuration header
	newContent.WriteString(fmt.Sprintf("# Configuration for %s\n", tool.Name))

	// Add aliases
	if len(config.Exports) > 0 {
		newContent.WriteString("\n# Aliases\n")
		for alias, cmd := range config.Exports {
			newContent.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, cmd))
		}
	}

	// Add functions
	if len(config.Functions) > 0 {
		newContent.WriteString("\n# Functions\n")
		for name, body := range config.Functions {
			newContent.WriteString(fmt.Sprintf("%s() {\n%s\n}\n", name, body))
		}
	}

	// Write updated .bashrc
	if err := os.WriteFile(bashrc, []byte(newContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write .bashrc: %w", err)
	}

	return nil
}

// applyFishConfig applies Fish-specific configuration
func (i *Installer) applyFishConfig(tool *Tool) error {
	config := tool.ShellConfig
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fishConfig := filepath.Join(homeDir, ".config", "fish", "config.fish")
	content := []string{}

	// Add environment variables
	if len(config.Exports) > 0 {
		content = append(content, "\n# Environment variables")
		for key, value := range config.Exports {
			content = append(content, fmt.Sprintf("set -x %s %s", key, value))
		}
	}

	// Add PATH additions
	if len(config.Path) > 0 {
		content = append(content, "\n# PATH additions")
		for _, path := range config.Path {
			content = append(content, fmt.Sprintf("fish_add_path %s", path))
		}
	}

	// Read existing fish config
	existingContent, err := os.ReadFile(fishConfig)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read fish config: %w", err)
	}

	// Prepare new content
	var newContent strings.Builder
	if len(existingContent) > 0 {
		newContent.Write(existingContent)
		newContent.WriteString("\n\n")
	}

	// Add configuration header
	newContent.WriteString(fmt.Sprintf("# Configuration for %s\n", tool.Name))

	// Add aliases
	if len(config.Exports) > 0 {
		newContent.WriteString("\n# Aliases\n")
		for alias, cmd := range config.Exports {
			newContent.WriteString(fmt.Sprintf("alias %s='%s'\n", alias, cmd))
		}
	}

	// Add functions
	if len(config.Functions) > 0 {
		newContent.WriteString("\n# Functions\n")
		for name, body := range config.Functions {
			newContent.WriteString(fmt.Sprintf("function %s\n%s\nend\n", name, body))
		}
	}

	// Write updated fish config
	if err := os.WriteFile(fishConfig, []byte(newContent.String()), 0644); err != nil {
		return fmt.Errorf("failed to write fish config: %w", err)
	}

	return nil
} 