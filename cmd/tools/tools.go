package tools

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
	"github.com/spf13/cobra"
)

var (
	skipVerification bool
	logger          *log.Logger
)

// NewToolsCmd creates a new tools command
func NewToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage development tools",
		Long: `Install and manage development tools.
This command helps you install and verify core development tools
like git, curl, build tools, and modern CLI utilities.`,
	}

	// Add install subcommand
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install core development tools",
		Long: `Install core development tools.
This will install essential tools like:
- Git, cURL, Wget
- Build tools (gcc, make)
- Modern CLI tools (bat, ripgrep, fzf)
- System monitoring (htop)`,
		RunE: runInstall,
	}

	// Add verify subcommand
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify tool installations",
		Long:  `Verify that all core development tools are properly installed and accessible.`,
		RunE:  runVerify,
	}

	// Add flags
	installCmd.Flags().BoolVar(&skipVerification, "skip-verify", false, "Skip verification after installation")

	// Add subcommands
	cmd.AddCommand(installCmd)
	cmd.AddCommand(verifyCmd)

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	logger.Info("Detecting system information...")

	// Detect system info
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Detect package manager
	pm, err := packages.DetectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.Name())

	// Create installation options
	opts := &tools.InstallOptions{
		Logger:           logger,
		PackageManager:   pm,
		SkipVerification: skipVerification,
	}

	// Install core tools
	if err := tools.InstallCoreTools(opts); err != nil {
		return fmt.Errorf("failed to install core tools: %w", err)
	}

	// Run verification unless skipped
	if !skipVerification {
		if err := tools.VerifyCoreTools(opts); err != nil {
			return fmt.Errorf("tool verification failed: %w", err)
		}
	}

	return nil
}

func runVerify(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	logger.Info("Detecting system information...")

	// Detect system info
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Detect package manager
	pm, err := packages.DetectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.Name())

	// Create verification options
	opts := &tools.InstallOptions{
		Logger:         logger,
		PackageManager: pm,
	}

	// Run verification
	if err := tools.VerifyCoreTools(opts); err != nil {
		return fmt.Errorf("tool verification failed: %w", err)
	}

	return nil
} 