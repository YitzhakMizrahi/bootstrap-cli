package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bootstrap-cli",
	Short: "Quickly bootstrap your development environment",
	Long: `bootstrap-cli is a modular, interactive command-line tool designed to 
bootstrap a personalized development environment quickly and consistently 
across platforms.

It allows you to select shells, plugin managers, prompts, editors, CLI tools, 
programming languages, and symlink your dotfiles intelligently.

🔥 Main Commands:
  • bootstrap-cli up     - Run the complete setup process (recommended)
  • bootstrap-cli init   - Run only the configuration wizard
  • bootstrap-cli install - Install selected components
  • bootstrap-cli link   - Link your dotfiles to the right locations

For more details on each command, run: bootstrap-cli [command] --help`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/bootstrap/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without making changes")
	
	// Set custom version template
	template := `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`
	rootCmd.SetVersionTemplate(template)
}

// GetVerbose returns the value of the verbose flag
func GetVerbose() bool {
	return verbose
}

// GetConfigPath returns the config file path
func GetConfigPath() string {
	return cfgFile
}