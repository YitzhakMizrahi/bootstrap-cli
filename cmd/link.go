// cmd/link.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/symlink"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Symlink your dotfiles to their correct locations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println("❌ Failed to load config:", err)
			return
		}

		err = symlink.LinkDotfiles(cfg)
		if err != nil {
			fmt.Println("❌ Failed to link dotfiles:", err)
			return
		}

		fmt.Println("✅ Dotfiles linked successfully!")
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
}
