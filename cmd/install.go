// cmd/install.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/installer"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install selected tools, languages, and editors",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("❌ Failed to load config:", err)
			return
		}

		fmt.Println("📦 Installing tools...")
		installer.InstallCLITools(cfg.CLITools)

		fmt.Println("🧪 Installing languages...")
		installer.InstallLanguages(cfg.Languages, cfg.PackageManagers)

		fmt.Println("📝 Setting up editors...")
		installer.InstallEditors(cfg.Editors)

		fmt.Println("✅ Installation complete.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
