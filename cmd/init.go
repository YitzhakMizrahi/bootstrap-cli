// cmd/init.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/prompts"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start the interactive setup wizard",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := prompts.RunInitWizard()

		err := config.Save(cfg)
		if err != nil {
			fmt.Println("❌ Failed to save config:", err)
			return
		}

		fmt.Println("✅ Config successfully saved. You’re ready to bootstrap!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
