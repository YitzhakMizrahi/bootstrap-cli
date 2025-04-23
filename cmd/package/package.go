package packagecmd

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/spf13/cobra"
)

var (
	packageName string
	system      string
)

// NewPackageCmd creates the package command
func NewPackageCmd() *cobra.Command {
	packageCmd := &cobra.Command{
		Use:   "package",
		Short: "Manage system packages",
		Long: `Manage system packages using the appropriate package manager for your system.
This command supports various package managers like apt, dnf, pacman, and more.`,
	}

	// Add subcommands
	packageCmd.AddCommand(newInstallCmd())
	packageCmd.AddCommand(newUninstallCmd())
	packageCmd.AddCommand(newUpdateCmd())
	packageCmd.AddCommand(newListCmd())
	packageCmd.AddCommand(newVersionCmd())

	// Add flags
	packageCmd.PersistentFlags().StringVarP(&system, "system", "s", "", "System type (ubuntu, debian, fedora, arch)")
	packageCmd.MarkPersistentFlagRequired("system")

	return packageCmd
}

// newInstallCmd creates the install command
func newInstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install [package]",
		Short: "Install a package",
		Long:  `Install a package using the system's package manager.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName = args[0]
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}
			return pm.Install(packageName)
		},
	}
}

// newUninstallCmd creates the uninstall command
func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall [package]",
		Short: "Uninstall a package",
		Long:  `Uninstall a package using the system's package manager.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName = args[0]
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}
			return pm.Remove(packageName)
		},
	}
}

// newUpdateCmd creates the update command
func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [package]",
		Short: "Update a package",
		Long:  `Update a package using the system's package manager.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName = args[0]
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}
			return pm.Update()
		},
	}
}

// newListCmd creates the list command
func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed packages",
		Long:  `List all installed packages using the system's package manager.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}
			pkgs, err := pm.ListInstalled()
			if err != nil {
				return fmt.Errorf("failed to list packages: %w", err)
			}
			for _, pkg := range pkgs {
				fmt.Println(pkg)
			}
			return nil
		},
	}
}

// newVersionCmd creates the version command
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version [package]",
		Short: "Get package version",
		Long:  `Get the version of an installed package.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName = args[0]
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to create package manager: %w", err)
			}
			version, err := pm.GetVersion(packageName)
			if err != nil {
				return fmt.Errorf("failed to get package version: %w", err)
			}
			fmt.Println(version)
			return nil
		},
	}
} 