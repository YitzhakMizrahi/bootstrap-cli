package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/YitzhakMizrahi/bootstrap-cli/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/installer"
	"github.com/YitzhakMizrahi/bootstrap-cli/pkg/platform"
	"github.com/YitzhakMizrahi/bootstrap-cli/prompts"
	"github.com/YitzhakMizrahi/bootstrap-cli/symlink"
)

// upCmd orchestrates init → install → link
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run full setup: init + install + link",
	Long: `bootstrap-cli up runs the complete setup process in one command.

It will:
1. Run the initialization wizard if no config is found
2. Install selected shells, tools, languages, and editors
3. Link dotfiles to the appropriate locations

This is the recommended command for new users.`,
Run: func(cmd *cobra.Command, args []string) {
    // Get flags
    force, _ := cmd.Flags().GetBool("force")
    skipInstall, _ := cmd.Flags().GetBool("skip-install")
    skipLink, _ := cmd.Flags().GetBool("skip-link")
    verbose, _ := cmd.Flags().GetBool("verbose")

    // Detect platform
    detector := platform.NewDetector()
    platformInfo, err := detector.Detect()
    if err != nil {
        fmt.Printf("❌ Failed to detect platform: %v\n", err)
        os.Exit(1)
    }

    if verbose {
        fmt.Printf("🔍 Detected platform: %s (%s)\n", platformInfo.OS, platformInfo.Distribution)
        fmt.Printf("🔍 Available package managers: %v\n", platformInfo.PackageManagers)
    }

    // For apt-based systems, set up a sudo session first and handle service prompts
    if platform.Contains(platformInfo.PackageManagers, platform.Apt) {
        // Obtain sudo early
        if err := installer.GetSudoSession(); err != nil {
            fmt.Printf("⚠️ Could not obtain sudo privileges. Some installations may fail: %v\n", err)
        }
        
        // Prevent service restart prompts during installation
        if err := installer.PreventServicePrompts(); err != nil {
            fmt.Printf("⚠️ Could not configure service prompt handling: %v\n", err)
        }
        // Set up a cleanup for when we're done
        defer installer.CleanupServicePromptConfig()
    }

    // Handle welcome message
    printWelcomeMessage()

    // Load or create configuration
    var cfg config.UserConfig
    
    loadedCfg, err := config.Load()
    if err != nil || force {
        if err != nil {
            fmt.Println("🧠 No config found or config error. Running initialization wizard...")
        } else if force {
            fmt.Println("🔄 Force flag detected, running init wizard...")
        }
        
        cfg = prompts.RunInitWizard()
        if err := config.Save(cfg); err != nil {
            fmt.Printf("❌ Failed to save config: %v\n", err)
            os.Exit(1)
        }
    } else {
        cfg = loadedCfg
        fmt.Println("🔧 Using existing configuration")
        if verbose {
            printConfigSummary(cfg)
        }
        fmt.Println("💡 To reset configuration, use the --force flag or run 'bootstrap-cli init'")
    }

		// Installation phase
		if !skipInstall {
			fmt.Println("\n📦 Starting installation phase...")

			// Install shell if needed
			if cfg.Shell != "" && !detector.IsCommandAvailable(cfg.Shell) {
				fmt.Printf("🐚 Installing %s shell...\n", cfg.Shell)
				err := installer.InstallShell(cfg.Shell)
				if err != nil {
					fmt.Printf("⚠️ Failed to install %s shell: %v\n", cfg.Shell, err)
				} else {
					fmt.Printf("✅ Successfully installed %s shell\n", cfg.Shell)
				}
			} else if cfg.Shell != "" {
				fmt.Printf("✅ %s shell is already installed\n", cfg.Shell)
			}

			// Install plugin manager if selected
			if cfg.PluginManager != "" && cfg.PluginManager != "none" {
				fmt.Printf("🔌 Setting up %s plugin manager...\n", cfg.PluginManager)
				err := installer.InstallPluginManager(cfg.PluginManager, cfg.Shell)
				if err != nil {
					fmt.Printf("⚠️ Failed to install plugin manager: %v\n", err)
				} else {
					fmt.Printf("✅ Successfully set up %s\n", cfg.PluginManager)
				}
			}

			// Install prompt if selected
			if cfg.Prompt != "" && cfg.Prompt != "none" {
				fmt.Printf("💅 Setting up %s prompt...\n", cfg.Prompt)
				err := installer.InstallPrompt(cfg.Prompt, cfg.Shell)
				if err != nil {
					fmt.Printf("⚠️ Failed to install prompt: %v\n", err)
				} else {
					fmt.Printf("✅ Successfully set up %s\n", cfg.Prompt)
				}
			}

			// Install CLI tools
			if len(cfg.CLITools) > 0 {
				fmt.Printf("🛠️ Installing CLI tools: %s\n", strings.Join(cfg.CLITools, ", "))
				if err := installer.InstallCLITools(cfg.CLITools); err != nil {
					fmt.Printf("⚠️ Some tools may have failed to install: %v\n", err)
				}
			} else {
				fmt.Println("ℹ️ No CLI tools selected for installation")
			}

			// Install programming languages
			if len(cfg.Languages) > 0 {
				fmt.Printf("🧪 Installing languages: %s\n", strings.Join(cfg.Languages, ", "))
				if err := installer.InstallLanguages(cfg.Languages, cfg.PackageManagers); err != nil {
					fmt.Printf("⚠️ Some languages may have failed to install: %v\n", err)
				}
			} else {
				fmt.Println("ℹ️ No programming languages selected for installation")
			}

			// Install editors
			if len(cfg.Editors) > 0 {
				fmt.Printf("📝 Setting up editors: %s\n", strings.Join(cfg.Editors, ", "))
				if err := installer.InstallEditors(cfg.Editors); err != nil {
					fmt.Printf("⚠️ Some editors may have failed to install: %v\n", err)
				}
			} else {
				fmt.Println("ℹ️ No editors selected for installation")
			}
		} else {
			fmt.Println("🔄 Skipping installation phase (--skip-install)")
		}

		// Symlink dotfiles
		if !skipLink && cfg.DotfilesPath != "" {
			fmt.Println("\n🔗 Linking dotfiles...")
			err := symlink.LinkDotfiles(cfg)
			if err != nil {
				fmt.Printf("❌ Failed to link dotfiles: %v\n", err)
			}
		} else if !skipLink && cfg.DotfilesPath == "" {
			fmt.Println("\n🔗 Preparing for dotfiles linking...")
			fmt.Println("ℹ️ No dotfiles path configured.")
			
			// Ask for dotfiles path
			fmt.Print("Enter path to your dotfiles repository [~/dotfiles]: ")
			var path string
			fmt.Scanln(&path)
			
			if path == "" {
				// Use default
				path = "~/dotfiles"
			}
			
			// Update the config with the new path
			cfg.DotfilesPath = path
			if err := config.Save(cfg); err != nil {
				fmt.Println("⚠️ Failed to save config with new dotfiles path:", err)
			} else {
				fmt.Println("✅ Updated config with dotfiles path:", path)
			}
			
			// Link the dotfiles
			err := symlink.LinkDotfiles(cfg)
			if err != nil {
				fmt.Printf("❌ Failed to link dotfiles: %v\n", err)
			}
		} else if skipLink {
			fmt.Println("🔄 Skipping dotfile linking (--skip-link)")
		} else {
			fmt.Println("ℹ️ No dotfiles path configured, skipping linking")
		}

		// Final message
		fmt.Println("\n✨ Bootstrap complete!")
		
		// Show next steps if not finishing
		skipFinish, _ := cmd.Flags().GetBool("skip-finish")
		if skipFinish {
			printNextSteps(cfg)
		} else {
			// Use the finisher to handle post-install tasks and shell restart
			if err := installer.FinishInstallation(cfg); err != nil {
				fmt.Printf("⚠️ Could not complete finalization: %v\n", err)
				// Fall back to showing next steps
				printNextSteps(cfg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Add flags
	upCmd.Flags().BoolP("force", "f", false, "Force run the init wizard even if config exists")
	upCmd.Flags().Bool("skip-install", false, "Skip the installation phase")
	upCmd.Flags().Bool("skip-link", false, "Skip the dotfile linking phase")
	upCmd.Flags().Bool("skip-finish", false, "Skip the finalization phase")
}

// printWelcomeMessage displays a welcome message
func printWelcomeMessage() {
	fmt.Println(`
✨ Welcome to bootstrap-cli! ✨

This tool will help you set up your development environment.
Follow the prompts to customize your setup.
`)
}

// printConfigSummary displays a summary of the configuration
func printConfigSummary(cfg config.UserConfig) {
	fmt.Println("\n📋 Configuration Summary:")
	fmt.Printf("  Shell: %s\n", cfg.Shell)
	fmt.Printf("  Plugin Manager: %s\n", cfg.PluginManager)
	fmt.Printf("  Prompt: %s\n", cfg.Prompt)

	if len(cfg.CLITools) > 0 {
		fmt.Printf("  CLI Tools: %s\n", strings.Join(cfg.CLITools, ", "))
	}

	if len(cfg.Languages) > 0 {
		fmt.Printf("  Languages: %s\n", strings.Join(cfg.Languages, ", "))
	}

	if len(cfg.Editors) > 0 {
		fmt.Printf("  Editors: %s\n", strings.Join(cfg.Editors, ", "))
	}

	if cfg.DotfilesPath != "" {
		fmt.Printf("  Dotfiles Path: %s\n", cfg.DotfilesPath)
	}
}

// printNextSteps displays next steps after bootstrapping
func printNextSteps(cfg config.UserConfig) {
	fmt.Println("\n📝 Next Steps:")
	fmt.Println("  1. Restart your terminal to apply shell changes")
	fmt.Println("  2. Run 'bootstrap-cli doctor' to verify your setup")
	fmt.Println("  3. Customize your dotfiles and shell configuration")
}

// containsString checks if a string slice contains a specific value
func containsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}