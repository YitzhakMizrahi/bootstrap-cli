package tools

import (
	"fmt"
	"os/user"

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

// isRoot checks if the current user has root privileges
func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		return false
	}
	return currentUser.Uid == "0"
}

// NewToolsCmd creates a new tools command that is meant to be used internally
// by the init command, not directly by users
func NewToolsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tools",
		Short:  "Manage development tools",
		Hidden: true, // Hide from user-facing CLI
		Long: `Install and manage development tools.
This command is used internally by the init command to install selected tools.`,
	}

	// Add install subcommand
	installCmd := &cobra.Command{
		Use:    "install",
		Short:  "Install selected development tools",
		Hidden: true,
		RunE:   runInstall,
	}

	// Add verify subcommand
	verifyCmd := &cobra.Command{
		Use:    "verify",
		Short:  "Verify tool installations",
		Hidden: true,
		RunE:   runVerify,
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

	// Detect package manager
	pm, err := packages.DetectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %w", err)
	}

	logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
	logger.Info("Package Manager: %s", pm.Name())

	// Get selected tools from the init context
	selectedTools := tools.GetSelectedTools()
	if len(selectedTools) == 0 {
		logger.Info("No tools selected, skipping installation...")
		return nil
	}

	// Create installation options
	opts := &tools.InstallOptions{
		Logger:           logger,
		PackageManager:   pm,
		Tools:           selectedTools,
		SkipVerification: skipVerification,
		// Add PATH to binary locations for verification
		AdditionalPaths: []string{"/usr/bin", "/usr/local/bin"},
	}

	// Install selected tools
	if err := tools.InstallCoreTools(opts); err != nil {
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