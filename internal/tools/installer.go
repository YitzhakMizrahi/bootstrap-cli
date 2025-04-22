package tools

import (
	"fmt"
	"os/exec"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// InstallOptions configures the tool installation process
type InstallOptions struct {
	// Logger for installation output
	Logger *log.Logger
	// PackageManager to use for installation
	PackageManager packages.PackageManager
	// Tools is the list of tools to install
	Tools []*install.Tool
	// SkipVerification skips the verification step
	SkipVerification bool
}

// InstallCoreTools installs all core development tools
func InstallCoreTools(opts *InstallOptions) error {
	if opts.Logger == nil {
		opts.Logger = log.New(log.InfoLevel)
	}

	if opts.Tools == nil {
		opts.Tools = CoreTools()
	}

	installer := install.NewInstaller(opts.PackageManager)
	installer.Logger = opts.Logger

	opts.Logger.Info("Installing core development tools...")

	for _, tool := range opts.Tools {
		opts.Logger.Info("Installing %s...", tool.Name)
		if err := installer.Install(tool); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool.Name, err)
		}
		opts.Logger.Success("%s installed successfully", tool.Name)
	}

	opts.Logger.Success("All core tools installed successfully!")
	return nil
}

// VerifyCoreTools checks if all core tools are properly installed
func VerifyCoreTools(opts *InstallOptions) error {
	if opts.Logger == nil {
		opts.Logger = log.New(log.InfoLevel)
	}

	if opts.Tools == nil {
		opts.Tools = CoreTools()
	}

	opts.Logger.Info("Verifying core tools installation...")

	for _, tool := range opts.Tools {
		opts.Logger.Debug("Verifying %s...", tool.Name)
		
		// Skip if no verify command
		if tool.VerifyCommand == "" {
			opts.Logger.Warn("No verification command for %s, skipping", tool.Name)
			continue
		}

		cmd := exec.Command("sh", "-c", tool.VerifyCommand)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("verification failed for %s: %w", tool.Name, err)
		}
		opts.Logger.Success("%s verified", tool.Name)
	}

	opts.Logger.Success("All core tools verified successfully!")
	return nil
} 