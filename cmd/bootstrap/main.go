package main

import (
	"fmt"
	"os"

	initCmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/init"
	packageCmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/package"
	toolsCmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/tools"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

// NewRootCmd creates the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap-cli",
		Short: "A tool to bootstrap your development environment",
		Long: `Bootstrap CLI helps you set up your development environment by installing
and configuring tools, shells, dotfiles, and more.`,
		// PersistentPreRun removed for now - logger setup handled in init command
	}

	// Add global debug flag
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd := NewRootCmd()

	// Add subcommands
	rootCmd.AddCommand(NewUpCmd()) // Add up command
	rootCmd.AddCommand(initCmd.NewInitCmd())
	rootCmd.AddCommand(toolsCmd.NewToolsCmd()) // Keep internal tools command hidden
	rootCmd.AddCommand(packageCmd.NewPackageCmd()) // Add package management command

	// Execute the root command
	Execute(rootCmd)
} 