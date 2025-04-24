package pipeline

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/utils"
)

// ToolCategory represents the category of a tool
type ToolCategory string

const (
	// Essential tools are required for basic system functionality
	CategoryEssential ToolCategory = "essential"
	// Development tools are used for software development
	CategoryDevelopment ToolCategory = "development"
	// Shell tools enhance the shell experience
	CategoryShell ToolCategory = "shell"
	// System tools are used for system management
	CategorySystem ToolCategory = "system"
)

// VerifyStrategy defines how to verify if a tool is installed correctly
type VerifyStrategy struct {
	// Command to run for verification
	Command string
	// Expected output (substring match)
	ExpectedOutput string
	// Binary paths to check
	BinaryPaths []string
	// Files that should exist
	RequiredFiles []string
}

// InstallStrategy defines how to install a tool
type InstallStrategy struct {
	// Package names for different package managers
	PackageNames map[string]string
	// Pre-install commands (e.g., adding PPAs)
	PreInstall []string
	// Post-install commands (e.g., setting up config)
	PostInstall []string
	// Custom installation script
	CustomInstall []string
}

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
	cmdExecutor *utils.CommandExecutor
}

// NewTool creates a new tool configuration
func NewTool(name string, category ToolCategory) *Tool {
	return &Tool{
		Name:           name,
		Category:       category,
		Dependencies:   make([]Dependency, 0),
		PlatformConfig: make(map[string]InstallStrategy),
		cmdExecutor:    utils.NewCommandExecutor(utils.NewDefaultLogger()),
	}
}

// defaultLogger implements the utils.InstallLogger interface
type defaultLogger struct {
	DebugEnabled bool
}

func (l *defaultLogger) CommandStart(cmd string, attempt, maxAttempts int) {
	if maxAttempts > 1 {
		fmt.Printf("Executing command (attempt %d/%d): %s\n", attempt, maxAttempts, cmd)
	} else {
		fmt.Printf("Executing command: %s\n", cmd)
	}
}

func (l *defaultLogger) CommandSuccess(cmd string, duration time.Duration) {
	fmt.Printf("Command completed successfully in %v: %s\n", duration, cmd)
}

func (l *defaultLogger) CommandError(cmd string, err error, attempt, maxAttempts int) {
	fmt.Printf("Command failed (attempt %d/%d): %s\nError: %v\n", attempt, maxAttempts, cmd, err)
}

func (l *defaultLogger) Debug(format string, args ...interface{}) {
	if l.DebugEnabled {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func (l *defaultLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *defaultLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

func (l *defaultLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
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

// VerifyInstallation checks if the tool is installed correctly
func (t *Tool) VerifyInstallation() error {
	// Check binary paths
	for _, path := range t.Verify.BinaryPaths {
		if _, err := exec.LookPath(path); err != nil {
			return fmt.Errorf("binary not found in PATH: %s", path)
		}
	}

	// Check required files
	for _, file := range t.Verify.RequiredFiles {
		if _, err := exec.Command("test", "-f", file).Output(); err != nil {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	// Run verification command if specified
	if t.Verify.Command != "" {
		cmd := exec.Command("sh", "-c", t.Verify.Command)
		output, err := t.cmdExecutor.ExecuteWithOutput(cmd, 0, 0)
		if err != nil {
			return fmt.Errorf("verification command failed: %v", err)
		}

		// Check expected output if specified
		if t.Verify.ExpectedOutput != "" && !strings.Contains(output, t.Verify.ExpectedOutput) {
			return fmt.Errorf("unexpected verification output: %s", output)
		}
	}

	return nil
}

// GenerateInstallationSteps creates installation steps for the pipeline
func (t *Tool) GenerateInstallationSteps(platform *Platform, context *InstallationContext) []InstallationStep {
	strategy := t.GetInstallStrategy(platform)
	var steps []InstallationStep

	// First, handle dependencies
	for _, dep := range t.Dependencies {
		depTool := context.GetTool(dep.Name)
		if depTool == nil {
			context.Logger.Printf("Warning: Dependency %s not found", dep.Name)
			continue
		}
		depSteps := depTool.GenerateInstallationSteps(platform, context)
		steps = append(steps, depSteps...)
	}

	// Add package manager update step with better error handling
	steps = append(steps, InstallationStep{
		Name: fmt.Sprintf("%s-update-package-lists", t.Name),
		Action: func() error {
			var cmd *exec.Cmd
			switch platform.PackageManager {
			case "apt":
				cmd = exec.Command("sudo", "apt-get", "update")
			case "brew":
				cmd = exec.Command("brew", "update")
			case "pacman":
				cmd = exec.Command("sudo", "pacman", "-Sy")
			default:
				return fmt.Errorf("unsupported package manager: %s", platform.PackageManager)
			}
			if err := t.cmdExecutor.ExecuteWithRetry(cmd, context.RetryCount, context.RetryDelay); err != nil {
				return fmt.Errorf("failed to update package lists: %w", err)
			}
			return nil
		},
	})

	// Add system dependencies installation step with validation
	if len(t.SystemDependencies) > 0 {
		steps = append(steps, InstallationStep{
			Name: fmt.Sprintf("%s-system-dependencies", t.Name),
			Action: func() error {
				// Validate package names first
				for _, pkg := range t.SystemDependencies {
					if !context.PackageManager.IsPackageAvailable(pkg) {
						return fmt.Errorf("system dependency %s is not available", pkg)
					}
				}

				var cmd *exec.Cmd
				switch platform.PackageManager {
				case "apt":
					args := append([]string{"apt-get", "install", "-y"}, t.SystemDependencies...)
					cmd = exec.Command("sudo", args...)
				case "brew":
					args := append([]string{"install"}, t.SystemDependencies...)
					cmd = exec.Command("brew", args...)
				case "pacman":
					args := append([]string{"pacman", "-S", "--noconfirm"}, t.SystemDependencies...)
					cmd = exec.Command("sudo", args...)
				default:
					return fmt.Errorf("unsupported package manager: %s", platform.PackageManager)
				}
				if err := t.cmdExecutor.ExecuteWithRetry(cmd, context.RetryCount, context.RetryDelay); err != nil {
					return fmt.Errorf("failed to install system dependencies: %w", err)
				}
				return nil
			},
		})
	}

	// Add pre-install steps with better error context
	for i, cmd := range strategy.PreInstall {
		steps = append(steps, InstallationStep{
			Name: fmt.Sprintf("%s-pre-install-%d", t.Name, i),
			Action: func() error {
				if err := t.cmdExecutor.ExecuteWithRetry(exec.Command("sh", "-c", cmd), context.RetryCount, context.RetryDelay); err != nil {
					return fmt.Errorf("pre-install step %d failed: %w", i, err)
				}
				return nil
			},
		})
	}

	// Add main installation step with package validation
	if pkgName, ok := strategy.PackageNames[platform.PackageManager]; ok {
		steps = append(steps, InstallationStep{
			Name: fmt.Sprintf("%s-install", t.Name),
			Action: func() error {
				// Validate package name
				if !context.PackageManager.IsPackageAvailable(pkgName) {
					return fmt.Errorf("package %s is not available", pkgName)
				}

				var cmd *exec.Cmd
				switch platform.PackageManager {
				case "apt":
					cmd = exec.Command("sudo", "apt-get", "install", "-y", pkgName)
				case "brew":
					cmd = exec.Command("brew", "install", pkgName)
				case "pacman":
					cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", pkgName)
				default:
					return fmt.Errorf("unsupported package manager: %s", platform.PackageManager)
				}
				if err := t.cmdExecutor.ExecuteWithRetry(cmd, context.RetryCount, context.RetryDelay); err != nil {
					return fmt.Errorf("failed to install %s: %w", t.Name, err)
				}
				return nil
			},
		})
	} else if len(strategy.CustomInstall) > 0 {
		// Use custom installation commands with better error context
		for i, cmd := range strategy.CustomInstall {
			steps = append(steps, InstallationStep{
				Name: fmt.Sprintf("%s-custom-install-%d", t.Name, i),
				Action: func() error {
					if err := t.cmdExecutor.ExecuteWithRetry(exec.Command("sh", "-c", cmd), context.RetryCount, context.RetryDelay); err != nil {
						return fmt.Errorf("custom install step %d failed: %w", i, err)
					}
					return nil
				},
			})
		}
	}

	// Add post-install steps with delay and better error context
	for i, cmd := range strategy.PostInstall {
		steps = append(steps, InstallationStep{
			Name: fmt.Sprintf("%s-post-install-%d", t.Name, i),
			Action: func() error {
				// Add a small delay to ensure the installation is complete
				time.Sleep(2 * time.Second)
				if err := t.cmdExecutor.ExecuteWithRetry(exec.Command("sh", "-c", cmd), context.RetryCount, context.RetryDelay); err != nil {
					return fmt.Errorf("post-install step %d failed: %w", i, err)
				}
				return nil
			},
		})
	}

	// Add verification step with delay
	steps = append(steps, InstallationStep{
		Name: fmt.Sprintf("%s-verify", t.Name),
		Action: func() error {
			// Add a delay to ensure all installation steps are complete
			time.Sleep(5 * time.Second)
			if err := t.VerifyInstallation(); err != nil {
				return fmt.Errorf("verification failed: %w", err)
			}
			return nil
		},
	})

	return steps
} 