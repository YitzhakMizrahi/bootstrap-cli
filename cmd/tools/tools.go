package tools

import (
	"fmt"
	"os/user"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	skipVerification bool
	logger          *log.Logger
)

// isRoot checks if the current user has root privileges
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

// NewToolsCmd creates the tools command
func NewToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Manage development tools",
		Long: `Manage development tools.
This command is used internally by the init command to install selected tools.
It provides functionality for installing and verifying development tools.`,
	}

	cmd.AddCommand(newInstallCmd())
	cmd.AddCommand(newVerifyCmd())

	return cmd
}

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install core development tools",
		Long: `Install core development tools.
This command is used internally by the init command to install selected tools.`,
		RunE: runInstall,
	}

	// Add flags
	cmd.Flags().BoolVar(&skipVerification, "skip-verify", false, "Skip verification after installation")

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}

	// Configure needrestart to run in automatic mode
	if err := configureNeedrestart(); err != nil {
		logger.Debug("Failed to configure needrestart: %v", err)
		// Non-fatal error, continue
	}

	// Detect system info
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Use the factory to get the package manager
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.GetName())

	// Get selected tools
	selectedTools := install.GetSelectedTools()
	if len(selectedTools) == 0 {
		logger.Info("No tools selected for installation.")
		return nil
	}

	// Create installation options
	opts := &install.InstallOptions{
		Logger:           logger,
		PackageManager:   pm,
		Tools:            selectedTools,
		SkipVerification: skipVerification,
		// Add PATH to binary locations for verification
		AdditionalPaths: []string{"/usr/bin", "/usr/local/bin"},
	}

	// Install selected tools
	if err := install.InstallCoreTools(opts); err != nil {
		return fmt.Errorf("failed to install core tools: %w", err)
	}

	return nil
}

// configureNeedrestart sets needrestart to automatic mode
func configureNeedrestart() error {
	// Create or update /etc/needrestart/conf.d/50-autorestart.conf
	content := []byte(`$nrconf{restart} = 'a';`)
	return system.WriteConfigFile("/etc/needrestart/conf.d/50-autorestart.conf", content)
}

func newVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify core development tools",
		Long: `Verify core development tools.
This command is used internally by the init command to verify selected tools.`,
		RunE: runVerify,
	}

	return cmd
}

func runVerify(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Detecting system information...")

	// Detect system info
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to detect system info: %w", err)
	}

	// Use the factory to get the package manager
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.GetName())

	// Create verification options
	opts := &install.InstallOptions{
		Logger:         logger,
		PackageManager: pm,
	}

	// Run verification
	if err := install.VerifyCoreTools(opts); err != nil {
		return fmt.Errorf("tool verification failed: %w", err)
	}

	return nil
} 