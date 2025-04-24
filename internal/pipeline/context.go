package pipeline

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// InstallationContext holds the context for an installation process
type InstallationContext struct {
	Platform      *Platform
	PackageManager PackageManager
	State         *InstallationState
	Logger        *log.Logger
	Timeout       time.Duration
	RetryCount    int
	RetryDelay    time.Duration
	tools         map[string]*Tool
}

// NewInstallationContext creates a new installation context
func NewInstallationContext(platform *Platform, pkgManager PackageManager) *InstallationContext {
	return &InstallationContext{
		Platform:      platform,
		PackageManager: pkgManager,
		State:         NewInstallationState(),
		Logger:        log.New(log.Writer(), "[Context] ", log.LstdFlags),
		Timeout:       5 * time.Minute,
		RetryCount:    3,
		RetryDelay:    time.Second,
		tools:         make(map[string]*Tool),
	}
}

// GetTool returns a tool by name
func (c *InstallationContext) GetTool(name string) *Tool {
	return c.tools[name]
}

// AddTool adds a tool to the context
func (c *InstallationContext) AddTool(tool *Tool) {
	c.tools[tool.Name] = tool
}

// VerifyInstallation verifies that a tool is properly installed
func (c *InstallationContext) VerifyInstallation(tool *Tool) error {
	if tool.Verify.Command == "" {
		return nil
	}

	cmd := exec.Command("sh", "-c", tool.Verify.Command)
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
		c.Logger.Printf("Executing post-install command: %s", cmd)
		execCmd := exec.Command("sh", "-c", cmd)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("post-install command failed: %w (Output: %s)", err, string(output))
		}
		c.Logger.Printf("Post-install command output: %s", string(output))
	}

	return nil
}

// setupAlias sets up a shell alias
func (c *InstallationContext) setupAlias(alias, command string) error {
	// Implementation depends on the shell being used
	// For now, we'll just log it
	c.Logger.Printf("Setting up alias: %s='%s'", alias, command)
	return nil
}

// setupFunction sets up a shell function
func (c *InstallationContext) setupFunction(name, body string) error {
	// Implementation depends on the shell being used
	// For now, we'll just log it
	c.Logger.Printf("Setting up function: %s() { %s }", name, body)
	return nil
}

// setupEnvVar sets up an environment variable
func (c *InstallationContext) setupEnvVar(key, value string) error {
	// Implementation depends on the shell being used
	// For now, we'll just log it
	c.Logger.Printf("Setting up environment variable: %s='%s'", key, value)
	return nil
}

// ExecutePostInstall executes post-installation commands
func (c *InstallationContext) ExecutePostInstall(tool *Tool) error {
	strategy := tool.GetInstallStrategy(c.Platform)
	if len(strategy.PostInstall) == 0 {
		return nil
	}

	for _, cmd := range strategy.PostInstall {
		c.Logger.Printf("Executing post-install command: %s", cmd)
		execCmd := exec.Command("sh", "-c", cmd)
		output, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("post-install command failed: %w (Output: %s)", err, string(output))
		}
		c.Logger.Printf("Post-install command output: %s", string(output))
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

	// Add paths to PATH if they don't exist
	for _, p := range paths {
		if !strings.Contains(path, p) {
			path = p + ":" + path
		}
	}

	// Set the new PATH
	if err := os.Setenv("PATH", path); err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}

	// Reload shell configuration
	if err := c.reloadShellConfig(); err != nil {
		return fmt.Errorf("failed to reload shell configuration: %w", err)
	}

	return nil
}

// reloadShellConfig reloads the shell configuration
func (c *InstallationContext) reloadShellConfig() error {
	shell := c.Platform.Shell
	switch shell {
	case "bash":
		return exec.Command("source", os.ExpandEnv("$HOME/.bashrc")).Run()
	case "zsh":
		return exec.Command("source", os.ExpandEnv("$HOME/.zshrc")).Run()
	case "fish":
		return exec.Command("source", os.ExpandEnv("$HOME/.config/fish/config.fish")).Run()
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
} 