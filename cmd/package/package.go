// Package packagecmd provides command-line functionality for managing system packages,
// including installation, removal, listing, updating, and upgrading packages through
// the system's native package manager.
package packagecmd

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/spf13/cobra"
)

var (
	packageName string
	system      string
	logger      *log.Logger
)

// NewPackageCmd creates the package command
func NewPackageCmd() *cobra.Command {
	packageCmd := &cobra.Command{
		Use:   "package",
		Short: "Manage system packages",
		Long: `Manage system packages using the system's package manager.`,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			logger = log.New(log.InfoLevel)
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				logger.SetLevel(log.DebugLevel)
			}
		},
	}

	// Add subcommands
	packageCmd.AddCommand(newInstallCmd())
	packageCmd.AddCommand(newRemoveCmd())
	packageCmd.AddCommand(newListCmd())
	packageCmd.AddCommand(newUpdateCmd())
	packageCmd.AddCommand(newUpgradeCmd())

	// Add flags
	packageCmd.PersistentFlags().StringVarP(&system, "system", "s", "", "System type (ubuntu, debian, fedora, arch)")
	packageCmd.MarkPersistentFlagRequired("system")

	return packageCmd
}

// newInstallCmd creates the install command
func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install [packages...]",
		Short: "Install packages",
		Long:  `Install packages using the system's package manager.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}

			for _, pkg := range args {
				if err := pm.Install(pkg); err != nil {
					return fmt.Errorf("failed to install package %s: %w", pkg, err)
				}
				logger.Info("Successfully installed %s", pkg)
			}
			return nil
		},
	}
}

// newRemoveCmd creates the remove command
func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [packages...]",
		Short: "Remove packages",
		Long:  `Remove packages using the system's package manager.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}

			for _, pkg := range args {
				if err := pm.Remove(pkg); err != nil {
					return fmt.Errorf("failed to remove package %s: %w", pkg, err)
				}
				logger.Info("Successfully removed %s", pkg)
			}
			return nil
		},
	}
}

// newListCmd creates the list command
func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [packages...]",
		Short: "List installed packages",
		Long:  `List all installed packages using the system's package manager.`,
		RunE:  runList,
	}
}

func runList(_ *cobra.Command, args []string) error {
	logger.Info("Listing installed packages...")

	// Get package manager
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return fmt.Errorf("failed to get package manager: %w", err)
	}

	// List installed packages
	logger.Info("Installed packages:")
	for _, pkg := range args {
		if pm.IsInstalled(pkg) {
			logger.Info("- %s", pkg)
		}
	}

	return nil
}

// newUpdateCmd creates the update command
func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update package list",
		Long:  `Update the package list using the system's package manager.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}

			if err := pm.Update(); err != nil {
				return fmt.Errorf("failed to update package list: %w", err)
			}
			logger.Info("Successfully updated package list")
			return nil
		},
	}
}

// newUpgradeCmd creates the upgrade command
func newUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade all packages",
		Long:  `Upgrade all installed packages using the system's package manager.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}

			if err := pm.Upgrade(); err != nil {
				return fmt.Errorf("failed to upgrade packages: %w", err)
			}
			logger.Info("Successfully upgraded all packages")
			return nil
		},
	}
} 