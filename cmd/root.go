package cmd

import (
	"fmt"
	"os"

	"github.com/YitzhakMizrahi/bootstrap-cli/cmd/tools"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/spf13/cobra"
)

var (
	debug  bool
	logger *log.Logger
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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add persistent flags that carry over to all commands
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	// Add commands
	rootCmd.AddCommand(tools.NewToolsCmd())
} 