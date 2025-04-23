package tools

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
)

// InstallOptions is defined in core.go

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
		if err := installer.Install(tool); err != nil {
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