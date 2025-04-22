package main

import (
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bootstrap-cli",
	Short: "A CLI tool for bootstrapping development environments",
	Long: `A CLI tool for bootstrapping development environments.
It helps you set up your development environment with all the necessary tools and configurations.`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new development environment",
	Long: `Initialize a new development environment by detecting the system
and setting up the necessary tools and configurations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := system.Detect()
		if err != nil {
			return fmt.Errorf("failed to detect system: %w", err)
		}

		fmt.Printf("Detected system:\n")
		fmt.Printf("  OS:      %s\n", info.OS)
		fmt.Printf("  Arch:    %s\n", info.Arch)
		fmt.Printf("  Distro:  %s\n", info.Distro)
		fmt.Printf("  Version: %s\n", info.Version)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
} 