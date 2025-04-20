// cmd/link.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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

		// Check if the dotfiles path is provided in the config
		if cfg.DotfilesPath == "" {
			// Ask for dotfiles path
			fmt.Println("ℹ️ No dotfiles path configured.")
			
			path, _ := cmd.Flags().GetString("path")
			if path == "" {
				fmt.Print("Enter path to your dotfiles repository [~/dotfiles]: ")
				fmt.Scanln(&path)
				
				if path == "" {
					// Use default
					path = "~/dotfiles"
				}
			}
			
			// Update the config with the new path
			cfg.DotfilesPath = path
			if err := config.Save(cfg); err != nil {
				fmt.Println("⚠️ Failed to save config with new dotfiles path:", err)
				// Continue anyway with the provided path
			} else {
				fmt.Println("✅ Updated config with dotfiles path:", path)
			}
		}

		err = symlink.LinkDotfiles(cfg)
		if err != nil {
			fmt.Println("❌ Failed to link dotfiles:", err)
			return
		}

		fmt.Println("✅ Dotfiles linked successfully!")
		
		// Suggest restarting shell if they linked shell config files
		home, _ := os.UserHomeDir()
		if _, err := os.Stat(filepath.Join(home, ".zshrc")); err == nil {
			fmt.Println("\n💡 Tip: Run 'source ~/.zshrc' to apply any shell configuration changes")
		} else if _, err := os.Stat(filepath.Join(home, ".bashrc")); err == nil {
			fmt.Println("\n💡 Tip: Run 'source ~/.bashrc' to apply any shell configuration changes")
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)
	
	// Add flag for dotfiles path
	linkCmd.Flags().StringP("path", "p", "", "Path to dotfiles repository")
}
