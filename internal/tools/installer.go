package tools

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// InstallOptions is defined in core.go

// convertTool converts an interfaces.Tool to an install.Tool
func convertTool(tool *interfaces.Tool) *install.Tool {
	return &install.Tool{
		Name:        tool.Name,
		Description: tool.Description,
		Category:    tool.Category,
		Tags:        tool.Tags,
		PackageName: tool.PackageName,
		PackageNames: &install.PackageMapping{
			APT:     tool.PackageNames.APT,
			DNF:     tool.PackageNames.DNF,
			Pacman:  tool.PackageNames.Pacman,
			Brew:    tool.PackageNames.Brew,
		},
		Version:       tool.Version,
		VerifyCommand: tool.VerifyCommand,
		ShellConfig: interfaces.ShellConfig{
			Aliases:   tool.ShellConfig.Aliases,
			Functions: tool.ShellConfig.Functions,
			Exports:   tool.ShellConfig.Env,
		},
	}
}

// InstallCoreTools installs the core development tools
func InstallCoreTools(opts *InstallOptions) error {
	// Create installer, logger is handled internally by the installer now
	installer := install.NewInstaller(opts.PackageManager)
	installer.Logger = opts.Logger // Assign logger if provided

	// Set custom retry settings if needed
	// installer.MaxRetries = 5
	// installer.RetryDelay = 5 * time.Second

	opts.Logger.Info("Installing core development tools...")

	failed := false
	for _, tool := range opts.Tools {
		opts.Logger.Info("Installing %s...", tool.Name)
		if err := installer.Install(convertTool(tool)); err != nil {
			opts.Logger.Error("Failed to install %s: %v", tool.Name, err)
			failed = true
			// Continue installing other tools even if one fails
			continue
		}
		opts.Logger.Info("Successfully installed %s", tool.Name)
		fmt.Printf("✓ %s installed successfully\n", tool.Name)
	}

	if failed {
		return fmt.Errorf("one or more tools failed to install")
	}

	// Run verification unless skipped
	if !opts.SkipVerification {
		opts.Logger.Info("Verifying tool installations...")
		if err := VerifyCoreTools(opts); err != nil {
			return fmt.Errorf("tool verification failed: %w", err)
		}
	}

	opts.Logger.Info("All selected tools installed and verified successfully.")
	return nil
}

// VerifyCoreTools is defined in verify.go
// func VerifyCoreTools(opts *InstallOptions) error { ... } 