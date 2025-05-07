package pipeline

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/cmdexec"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// ToolCategory represents the category of a tool
type ToolCategory string

const (
	// CategoryEssential represents tools that are required for basic system functionality
	CategoryEssential ToolCategory = "essential"
	// CategoryDevelopment represents tools that are used for software development
	CategoryDevelopment ToolCategory = "development"
	// CategoryShell represents tools that enhance the shell experience
	CategoryShell ToolCategory = "shell"
	// CategorySystem represents tools that are used for system management
	CategorySystem ToolCategory = "system"
)

// VerifyStrategy defines how to verify if a tool is installed correctly
type VerifyStrategy struct {
	// Command to run for verification
	Command Command
	// Expected output (substring match)
	ExpectedOutput string
	// Binary paths to check
	BinaryPaths []string
	// Files that should exist
	RequiredFiles []string
}

// InstallationMethod defines how a tool should be installed
type InstallationMethod string

const (
	// PackageManagerInstall uses the system's package manager
	PackageManagerInstall InstallationMethod = "package_manager"
	// BinaryInstall downloads and installs a binary directly
	BinaryInstall InstallationMethod = "binary"
	// CustomInstall uses custom installation commands
	CustomInstall InstallationMethod = "custom"
)

// PackageInfo contains information about a package's availability
type PackageInfo struct {
	// Whether the package is available in repositories
	Available bool
	// The actual package name to use (may differ from requested name)
	PackageName string
	// Version available in repositories
	Version string
	// Error message if package is not available
	Error string
}

// Platform represents the target platform for installation
// This type is defined in platform.go

// InstallationContext provides context for the installation process
// This type is defined in context.go

// PackageManager interface defines methods for package management
// This type is defined in interfaces.go

// Logger interface for installation logging
// This type is defined in interfaces/logger.go

// Tool represents a tool that can be installed
type Tool struct {
	Name        string
	Category    ToolCategory
	Description string
	Version     string
	Homepage    string
	Tags        []string

	// Dependencies required by this tool
	Dependencies []Dependency

	// System dependencies required by this tool
	SystemDependencies []string

	// Installation strategy
	Install InstallStrategy

	// Verification strategy
	Verify VerifyStrategy

	// Platform-specific configuration
	PlatformConfig map[string]InstallStrategy
	
	// Command executor for running commands
	cmdExecutor *cmdexec.CommandExecutor

	// Logger for the tool
	logger interfaces.Logger
}

// NewTool creates a new tool with the given name and category
func NewTool(name string, category ToolCategory) *Tool {
	logger := log.New(log.InfoLevel)
	return &Tool{
		Name:            name,
		Category:        category,
		PlatformConfig:  make(map[string]InstallStrategy),
		cmdExecutor:     cmdexec.NewCommandExecutor(logger),
		logger:          logger,
	}
}

// AddDependency adds a dependency to the tool
func (t *Tool) AddDependency(dep Dependency) {
	t.Dependencies = append(t.Dependencies, dep)
}

// SetVerification sets the verification strategy
func (t *Tool) SetVerification(verify VerifyStrategy) {
	t.Verify = verify
}

// SetInstallation sets the installation strategy
func (t *Tool) SetInstallation(install InstallStrategy) {
	t.Install = install
}

// SetPlatformConfig sets platform-specific installation strategy
func (t *Tool) SetPlatformConfig(platform string, install InstallStrategy) {
	t.PlatformConfig[platform] = install
}

// GetInstallStrategy returns the appropriate installation strategy for the platform
func (t *Tool) GetInstallStrategy(platform *Platform) InstallStrategy {
	// Check for platform-specific config
	if strategy, ok := t.PlatformConfig[platform.OS]; ok {
		return strategy
	}
	return t.Install
}

// checkBinaryPath checks if a binary exists in the PATH
func (t *Tool) checkBinaryPath(path string) (bool, error) {
	_, err := exec.LookPath(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// checkRequiredFile checks if a required file exists
func (t *Tool) checkRequiredFile(file string) (bool, error) {
	cmd := exec.Command("test", "-f", file)
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}

// VerifyInstallation checks if the tool is installed correctly
func (t *Tool) VerifyInstallation(_ *InstallationContext) error {
	// Check binary paths with retries
	for _, path := range t.Verify.BinaryPaths {
		var exists bool
		var err error
		for i := 0; i < t.cmdExecutor.DefaultRetries; i++ {
			exists, err = t.checkBinaryPath(path)
			if err == nil && exists {
				break
			}
			if i < t.cmdExecutor.DefaultRetries-1 {
				t.cmdExecutor.Logger.Debug("Binary path %s not found, retrying in %v...", path, t.cmdExecutor.DefaultDelay)
				time.Sleep(t.cmdExecutor.DefaultDelay)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to check binary path %s: %w", path, err)
		}
		if !exists {
			return fmt.Errorf("binary path %s not found after %d attempts", path, t.cmdExecutor.DefaultRetries)
		}
	}

	// Check required files with retries
	for _, file := range t.Verify.RequiredFiles {
		var exists bool
		var err error
		for i := 0; i < t.cmdExecutor.DefaultRetries; i++ {
			exists, err = t.checkRequiredFile(file)
			if err == nil && exists {
				break
			}
			if i < t.cmdExecutor.DefaultRetries-1 {
				t.cmdExecutor.Logger.Debug("Required file %s not found, retrying in %v...", file, t.cmdExecutor.DefaultDelay)
				time.Sleep(t.cmdExecutor.DefaultDelay)
			}
		}
		if err != nil {
			return fmt.Errorf("failed to check required file %s: %w", file, err)
		}
		if !exists {
			return fmt.Errorf("required file %s not found after %d attempts", file, t.cmdExecutor.DefaultRetries)
		}
	}

	// Execute verification command with timeout
	if t.Verify.Command.Command != "" {
		cmd := exec.Command("sh", "-c", t.Verify.Command.Command)
		if err := t.cmdExecutor.ExecuteWithRetry(cmd, t.cmdExecutor.DefaultRetries, t.cmdExecutor.DefaultDelay); err != nil {
			return fmt.Errorf("verification command failed: %w", err)
		}
	}

	return nil
}

// determineInstallationMethod determines the best installation method for the tool
func (t *Tool) determineInstallationMethod(context *InstallationContext) (InstallationMethod, error) {
	// Get the package name for the current platform
	packageName := t.Install.PackageNames[context.Platform.OS]
	if packageName == "" {
		packageName = t.Install.PackageNames["default"]
	}

	// Check if package is available in repositories
	if !context.PackageManager.IsPackageAvailable(packageName) {
		return "", fmt.Errorf("package %s is not available", packageName)
	}

	// If package is available, use package manager installation
	return PackageManagerInstall, nil
}

// GenerateInstallationSteps generates the steps needed to install a tool.
// If skipDependencyResolution is true, the initial dependency resolution step is omitted.
func (t *Tool) GenerateInstallationSteps(platform *Platform, context *InstallationContext, skipDependencyResolution bool) []InstallationStep {
	var steps []InstallationStep
	
	// First, resolve dependencies unless skipped
	if !skipDependencyResolution {
		steps = append(steps, InstallationStep{
			Name: fmt.Sprintf("%s-resolve-dependencies", t.Name),
			Description: fmt.Sprintf("Resolving dependencies for %s", t.Name),
			Action: func(ctx *InstallationContext) error {
				// Note: This might still be problematic if context.ResolveDependencies assumes
				// it should install ALL dependencies in the graph vs just those for 't'.
				// It might need adjustment if called from the old single Install path.
				return ctx.ResolveDependencies(t)
			},
			Timeout: 5 * time.Minute,
		})
	}
	
	// Determine installation method
	method, err := t.determineInstallationMethod(context)
	if err != nil {
		t.logger.Error("Failed to determine installation method: %v", err)
		return steps
	}
	
	// Get the appropriate installation strategy
	strategy := t.GetInstallStrategy(platform)
	
	// Add pre-install steps
	for i, cmd := range strategy.PreInstall {
		stepName := fmt.Sprintf("%s-pre-install-%d", t.Name, i)
		preCmd := cmd
		steps = append(steps, InstallationStep{
			Name: stepName,
			Description: preCmd.Description,
			Action: func(ctx *InstallationContext) error {
				ctx.Logger.CommandStart(preCmd.Command, 1, 1)
				start := time.Now()
				
				execCmd := exec.Command("sh", "-c", preCmd.Command)
				output, err := execCmd.CombinedOutput()
				
				duration := time.Since(start)
				if err != nil {
					ctx.Logger.CommandError(preCmd.Command, err, 1, 1)
					return fmt.Errorf("pre-install command failed: %w (Output: %s)", err, string(output))
				}
				ctx.Logger.CommandSuccess(preCmd.Command, duration)
				return nil
			},
			Timeout: 5 * time.Minute,
		})
	}
	
	// Add main installation step based on method
	switch method {
	case PackageManagerInstall:
		// Get package name for the platform
		pkgName, ok := strategy.PackageNames[platform.PackageManager]
		if !ok {
			t.logger.Error("No package name defined for %s on %s", t.Name, platform.PackageManager)
			return steps
		}
		
		stepName := fmt.Sprintf("%s-install-package", t.Name)
		steps = append(steps, InstallationStep{
			Name: stepName,
			Description: fmt.Sprintf("Installing %s via %s", pkgName, platform.PackageManager),
			Action: func(ctx *InstallationContext) error {
				var cmdStr string
				switch platform.PackageManager {
				case "apt":
					cmdStr = fmt.Sprintf("sudo apt-get install -y %s", pkgName)
				case "brew":
					cmdStr = fmt.Sprintf("brew install %s", pkgName)
				case "pacman":
					cmdStr = fmt.Sprintf("sudo pacman -S --noconfirm %s", pkgName)
				default:
					return fmt.Errorf("unsupported package manager: %s", platform.PackageManager)
				}
				
				ctx.Logger.CommandStart(cmdStr, 1, 1)
				start := time.Now()
				
				execCmd := exec.Command("sh", "-c", cmdStr)
				output, err := execCmd.CombinedOutput()
				
				duration := time.Since(start)
				if err != nil {
					ctx.Logger.CommandError(cmdStr, err, 1, 1)
					return fmt.Errorf("package installation failed: %w (Output: %s)", err, string(output))
				}
				ctx.Logger.CommandSuccess(cmdStr, duration)
				return nil
			},
			Timeout: 10 * time.Minute,
		})
		
	case BinaryInstall:
		// Binary installation is not directly supported in the current InstallStrategy
		// We'll use custom installation instead
		t.logger.Warn("Binary installation not directly supported, using custom installation")
		fallthrough
		
	case CustomInstall:
		// Custom installation steps
		for i, cmd := range strategy.CustomInstall {
			stepName := fmt.Sprintf("%s-custom-install-%d", t.Name, i)
			customCmd := cmd
			steps = append(steps, InstallationStep{
				Name: stepName,
				Description: customCmd.Description,
				Action: func(ctx *InstallationContext) error {
					ctx.Logger.CommandStart(customCmd.Command, 1, 1)
					start := time.Now()
					
					execCmd := exec.Command("sh", "-c", customCmd.Command)
					output, err := execCmd.CombinedOutput()
					
					duration := time.Since(start)
					if err != nil {
						ctx.Logger.CommandError(customCmd.Command, err, 1, 1)
						return fmt.Errorf("custom installation command failed: %w (Output: %s)", err, string(output))
					}
					ctx.Logger.CommandSuccess(customCmd.Command, duration)
					return nil
				},
				Timeout: 5 * time.Minute,
			})
		}
	}
	
	// Add post-install steps
	for i, cmd := range strategy.PostInstall {
		stepName := fmt.Sprintf("%s-post-install-%d", t.Name, i)
		postCmd := cmd
		steps = append(steps, InstallationStep{
			Name: stepName,
			Description: postCmd.Description,
			Action: func(ctx *InstallationContext) error {
				ctx.Logger.CommandStart(postCmd.Command, 1, 1)
				start := time.Now()
				
				execCmd := exec.Command("sh", "-c", postCmd.Command)
				output, err := execCmd.CombinedOutput()
				
				duration := time.Since(start)
				if err != nil {
					ctx.Logger.CommandError(postCmd.Command, err, 1, 1)
					return fmt.Errorf("post-install command failed: %w (Output: %s)", err, string(output))
				}
				ctx.Logger.CommandSuccess(postCmd.Command, duration)
				return nil
			},
			Timeout: 5 * time.Minute,
		})
	}
	
	// Add verification step
	steps = append(steps, InstallationStep{
		Name: fmt.Sprintf("%s-verify", t.Name),
		Description: fmt.Sprintf("Verifying installation of %s", t.Name),
		Action: func(ctx *InstallationContext) error {
			return t.VerifyInstallation(ctx)
		},
		Timeout: 1 * time.Minute,
	})
	
	return steps
}

// Validate checks if the tool configuration is valid
func (t *Tool) Validate() error {
	// Check required fields
	if t.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if t.Category == "" {
		return fmt.Errorf("tool category cannot be empty")
	}
	if t.Description == "" {
		return fmt.Errorf("tool description cannot be empty")
	}

	// Validate category
	switch t.Category {
	case CategoryEssential, CategoryDevelopment, CategoryShell, CategorySystem:
		// Valid categories
	default:
		return fmt.Errorf("invalid tool category: %s", t.Category)
	}

	// Validate installation strategy
	if err := t.Install.Validate(); err != nil {
		return fmt.Errorf("invalid installation strategy: %w", err)
	}

	// Validate verification strategy
	if err := t.Verify.Validate(); err != nil {
		return fmt.Errorf("invalid verification strategy: %w", err)
	}

	// Validate dependencies
	for i, dep := range t.Dependencies {
		if dep.Name == "" {
			return fmt.Errorf("dependency %d name cannot be empty", i)
		}
		if dep.Version == "" {
			return fmt.Errorf("dependency %d version cannot be empty", i)
		}
	}

	// Validate system dependencies
	for i, dep := range t.SystemDependencies {
		if dep == "" {
			return fmt.Errorf("system dependency %d cannot be empty", i)
		}
	}

	// Validate platform-specific configurations
	for platform, strategy := range t.PlatformConfig {
		if platform == "" {
			return fmt.Errorf("platform name cannot be empty")
		}
		if err := strategy.Validate(); err != nil {
			return fmt.Errorf("invalid platform-specific strategy for %s: %w", platform, err)
		}
	}

	return nil
}

// Validate checks if the verification strategy is valid
func (v *VerifyStrategy) Validate() error {
	// At least one verification method must be specified
	if v.Command.Command == "" && len(v.BinaryPaths) == 0 && len(v.RequiredFiles) == 0 {
		return fmt.Errorf("at least one verification method must be specified")
	}

	// Validate command if specified
	if v.Command.Command != "" {
		if err := v.Command.Validate(); err != nil {
			return fmt.Errorf("invalid verification command: %w", err)
		}
	}

	// Validate binary paths
	for i, path := range v.BinaryPaths {
		if path == "" {
			return fmt.Errorf("binary path %d cannot be empty", i)
		}
	}

	// Validate required files
	for i, file := range v.RequiredFiles {
		if file == "" {
			return fmt.Errorf("required file %d cannot be empty", i)
		}
	}

	return nil
} 