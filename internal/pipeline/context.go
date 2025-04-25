// Package pipeline provides a flexible and extensible installation pipeline system
// for the bootstrap-cli, managing the installation process, dependencies, and
// post-installation tasks.
package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
)

// InstallationContext holds the context for an installation process
type InstallationContext struct {
	Platform       *Platform
	PackageManager PackageManager
	State         *InstallationState
	Logger        interfaces.Logger
	Timeout       time.Duration
	RetryCount    int
	RetryDelay    time.Duration
	tools         map[string]*Tool
	shellConfig   *shell.Config
	// Add dependency graph
	dependencyGraph *DependencyGraph
	// Track installed tools
	installedTools map[string]bool
}

// NewInstallationContext creates a new installation context
func NewInstallationContext(platform *Platform, pkgManager PackageManager) *InstallationContext {
	logger := log.NewInstallLogger(false)
	return &InstallationContext{
		Platform:       platform,
		PackageManager: pkgManager,
		State:         NewInstallationState(),
		Logger:        logger,
		Timeout:       5 * time.Minute,
		RetryCount:    3,
		RetryDelay:    time.Second,
		tools:         make(map[string]*Tool),
		shellConfig:   shell.NewConfig(platform.Shell, logger),
		dependencyGraph: NewDependencyGraph(),
		installedTools: make(map[string]bool),
	}
}

// GetTool returns a tool by name
func (c *InstallationContext) GetTool(name string) *Tool {
	return c.tools[name]
}

// AddTool adds a tool to the context and its dependencies to the graph
func (c *InstallationContext) AddTool(tool *Tool) {
	c.tools[tool.Name] = tool
	
	// Add tool dependencies to the graph
	if len(tool.Dependencies) > 0 {
		c.dependencyGraph.AddDependency(tool.Name, tool.Dependencies)
	}
}

// VerifyInstallation verifies that a tool is properly installed
func (c *InstallationContext) VerifyInstallation(tool *Tool) error {
	if tool.Verify.Command.Command == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", tool.Verify.Command.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("verification failed: %w (Output: %s)", err, string(output))
	}

	// Check if the command output indicates success
	if tool.Verify.ExpectedOutput != "" && !strings.Contains(string(output), tool.Verify.ExpectedOutput) {
		return fmt.Errorf("verification failed: unexpected output (Output: %s)", string(output))
	}

	// Check binary paths
	for _, path := range tool.Verify.BinaryPaths {
		if _, err := exec.LookPath(path); err != nil {
			return fmt.Errorf("binary not found in PATH: %s", path)
		}
	}

	// Check required files
	for _, file := range tool.Verify.RequiredFiles {
		if _, err := exec.Command("test", "-f", file).Output(); err != nil {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	return nil
}

// SetupEnvironment sets up the environment for a tool
func (c *InstallationContext) SetupEnvironment(tool *Tool) error {
	// Get the installation strategy for the current platform
	strategy := tool.GetInstallStrategy(c.Platform)

	// Execute post-install commands
	for _, cmd := range strategy.PostInstall {
		c.Logger.Info("Executing post-install command: %s", cmd.Command)
		execCmd := exec.Command("sh", "-c", cmd.Command)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			c.Logger.Error("Post-install command failed: %v (Output: %s)", err, string(output))
			return fmt.Errorf("post-install command failed: %w (Output: %s)", err, string(output))
		}
		c.Logger.Info("Post-install command output: %s", string(output))
	}

	return nil
}

// setupAlias sets up a shell alias
func (c *InstallationContext) setupAlias(alias, command string) error {
	c.Logger.Info("Setting up alias: %s='%s'", alias, command)
	c.shellConfig.AddAlias(alias, command)
	return nil
}

// setupFunction sets up a shell function
func (c *InstallationContext) setupFunction(name, body string) error {
	c.Logger.Info("Setting up function: %s() { %s }", name, body)
	c.shellConfig.AddFunction(name, body)
	return nil
}

// setupEnvVar sets up an environment variable
func (c *InstallationContext) setupEnvVar(key, value string) error {
	c.Logger.Info("Setting up environment variable: %s='%s'", key, value)
	c.shellConfig.AddEnvVar(key, value)
	return nil
}

// ExecutePostInstall executes post-installation commands
func (c *InstallationContext) ExecutePostInstall(tool *Tool) error {
	strategy := tool.GetInstallStrategy(c.Platform)
	if len(strategy.PostInstall) == 0 {
		return nil
	}

	for _, cmd := range strategy.PostInstall {
		c.Logger.Info("Executing post-install command: %s", cmd.Command)
		execCmd := exec.Command("sh", "-c", cmd.Command)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			c.Logger.Error("Post-install command failed: %v (Output: %s)", err, string(output))
			return fmt.Errorf("post-install command failed: %w (Output: %s)", err, string(output))
		}
		c.Logger.Info("Post-install command output: %s", string(output))
	}

	return nil
}

// UpdatePath updates the PATH environment variable with installed binary paths
func (c *InstallationContext) UpdatePath() error {
	// Get the current PATH
	path := os.Getenv("PATH")
	if path == "" {
		path = "/usr/local/bin:/usr/bin:/bin"
	}

	// Add common binary paths
	paths := []string{
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/usr/local/go/bin",
		os.ExpandEnv("$HOME/.local/bin"),
		os.ExpandEnv("$HOME/go/bin"),
		os.ExpandEnv("$HOME/.cargo/bin"),
	}

	// Add paths to shell config
	for _, p := range paths {
		if !strings.Contains(path, p) {
			c.shellConfig.AddPath(p)
		}
	}

	// Apply shell configuration
	sourceCmd, err := c.shellConfig.Apply()
	if err != nil {
		return fmt.Errorf("failed to apply shell configuration: %w", err)
	}

	// Execute the source command
	cmd := exec.Command("sh", "-c", sourceCmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload shell configuration: %w (output: %s)", err, string(output))
	}

	return nil
}

// reloadShellConfig reloads the shell configuration
func (c *InstallationContext) reloadShellConfig() error {
	sourceCmd, err := c.shellConfig.Apply()
	if err != nil {
		return fmt.Errorf("failed to apply shell configuration: %w", err)
	}

	cmd := exec.Command("sh", "-c", sourceCmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reload shell configuration: %w (output: %s)", err, string(output))
	}

	return nil
}

// ResolveDependencies resolves and installs all dependencies for a tool
func (c *InstallationContext) ResolveDependencies(tool *Tool) error {
	c.Logger.Info("Resolving dependencies for %s", tool.Name)
	
	// Get all dependencies for the tool
	deps := c.dependencyGraph.GetDependencies(tool.Name)
	if len(deps) == 0 {
		c.Logger.Info("No dependencies found for %s", tool.Name)
		return nil
	}
	
	// Get installation order
	order, err := c.dependencyGraph.GetInstallOrder()
	if err != nil {
		return fmt.Errorf("failed to determine installation order: %w", err)
	}
	
	// Install dependencies in order
	for _, depName := range order {
		// Skip if already installed
		if c.installedTools[depName] {
			c.Logger.Info("Dependency %s already installed, skipping", depName)
			continue
		}
		
		// Find the dependency tool
		depTool, exists := c.tools[depName]
		if !exists {
			return fmt.Errorf("dependency %s not found in available tools", depName)
		}
		
		// Install the dependency
		c.Logger.Info("Installing dependency %s for %s", depName, tool.Name)
		if err := c.installTool(depTool); err != nil {
			return fmt.Errorf("failed to install dependency %s: %w", depName, err)
		}
		
		// Mark as installed
		c.installedTools[depName] = true
	}
	
	return nil
}

// installTool installs a tool using the pipeline
func (c *InstallationContext) installTool(tool *Tool) error {
	// Generate installation steps
	steps := tool.GenerateInstallationSteps(c.Platform, c)
	
	// Execute each step
	for _, step := range steps {
		c.Logger.Info("Executing step: %s", step.Name)
		
		// Execute the step
		if err := step.Action(); err != nil {
			c.Logger.Error("Step failed: %v", err)
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}
		
		c.Logger.Info("Step completed successfully")
	}
	
	return nil
} 