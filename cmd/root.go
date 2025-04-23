package cmd

import (
	"fmt"
	"os"

	initcmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/init"
	packagecmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/package"
	"github.com/YitzhakMizrahi/bootstrap-cli/cmd/tools"
	upcmd "github.com/YitzhakMizrahi/bootstrap-cli/cmd/up"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/spf13/cobra"
)

var (
	debug      bool
	logger     *log.Logger
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bootstrap-cli",
	Short: "A CLI tool for bootstrapping development environments",
	Long: `Bootstrap CLI is a comprehensive tool for setting up development environments.
It helps you install and configure:
- Core development tools (git, curl, build tools)
- Modern CLI utilities (bat, ripgrep, fzf)
- Shell configurations and plugins
- Programming language environments
- Dotfiles management`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up logging based on debug flag
		if debug {
			logger = log.New(log.DebugLevel)
		} else {
			logger = log.New(log.InfoLevel)
		}
		
		// Set config path in environment for child processes
		if configPath != "" {
			os.Setenv("BOOTSTRAP_CLI_CONFIG", configPath)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add flags
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "Path to config directory")

	// Add commands
	rootCmd.AddCommand(initcmd.NewInitCmd())
	rootCmd.AddCommand(packagecmd.NewPackageCmd())
	rootCmd.AddCommand(tools.NewToolsCmd())
	rootCmd.AddCommand(upcmd.NewUpCmd())
} 