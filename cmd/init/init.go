package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	logger *log.Logger
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new development environment",
		Long: `Guides you through selecting and installing components
for your development environment, including:
- Core development tools
- Modern CLI utilities
- Shell setup (coming soon)
- Dotfiles management (coming soon)`,
		RunE: runInit,
	}
	// Add flags here if needed later
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	logger = log.New(log.InfoLevel)
	if debug, _ := cmd.Flags().GetBool("debug"); debug {
		logger.SetLevel(log.DebugLevel)
	}
	logger.Info("Starting Bootstrap CLI initialization...")

	// --- Phase 1: Tool Selection ---
	logger.Info("Step 1: Select tools to install")
	defaultTools := []string{
		"git", "curl", "wget", "tmux",
		"ripgrep", "bat", "fzf", "exa",
		"htop", "btop", "neofetch",
	}
	selector := ui.NewToolSelector(defaultTools)
	p := tea.NewProgram(selector)

	// Run the interactive UI
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("UI error: %w", err)
	}

	// Check if user quit
	if selectorModel, ok := model.(*ui.ToolSelector); ok {
		selectedTools := selectorModel.GetSelectedTools()
		if len(selectedTools) == 0 && !selectorModel.Finished() { // Check if quit early or finished with no selection
			logger.Info("Initialization cancelled or no tools selected.")
			return nil
		}

		// Convert string slice to []*install.Tool
		var toolObjects []*install.Tool
		for _, name := range selectedTools {
			toolObjects = append(toolObjects, &install.Tool{
				Name:        name,
				PackageName: name,
				Description: "Selected tool", // Basic description
			})
		}

		// Store selected tools for the installer
		tools.SetSelectedTools(toolObjects)

		if len(selectedTools) == 0 {
			logger.Info("No tools selected. Skipping tool installation.")
		} else {
			// --- Phase 2: Installation (Requires Sudo) ---
			logger.Info("Step 2: Installing selected tools...")

			// Check for root privileges before proceeding with installation
			if !system.IsRoot() {
				logger.Warn("Tool installation requires root privileges.")
				logger.Info("Attempting to relaunch with sudo...")
				// TODO: Implement a more robust sudo prompt/handling mechanism
				// For now, just error out and instruct the user.
				fmt.Printf("\nPlease re-run the command with sudo: sudo %s\n\n", strings.Join(os.Args, " "))
				return fmt.Errorf("root privileges required for tool installation")
			}

			// Configure needrestart automatically (best effort)
			if err := configureNeedrestart(); err != nil {
				logger.Warn("Failed to configure needrestart for automatic restarts: %v", err)
			}

			// Detect system info and package manager
			sysInfo, err := system.Detect()
			if err != nil {
				return fmt.Errorf("failed to detect system info: %w", err)
			}
			
			// Use the factory to get the package manager
			f := factory.NewPackageManagerFactory()
			pm, err := f.GetPackageManager()
			if err != nil {
				return fmt.Errorf("failed to detect package manager: %w", err)
			}
			logger.Info("System: %s %s (%s)", sysInfo.Distro, sysInfo.Version, sysInfo.OS)
			logger.Info("Package Manager: %s", pm.Name())

			// Create installation options
			opts := &tools.InstallOptions{
				Logger:           logger,
				PackageManager:   pm,
				Tools:            toolObjects,
				SkipVerification: false, // Default to verify
				AdditionalPaths:  []string{"/usr/bin", "/usr/local/bin", "/opt/homebrew/bin"}, // Common paths
			}

			// Install selected tools (using the internal function)
			if err := tools.InstallCoreTools(opts); err != nil {
				return fmt.Errorf("failed to install selected tools: %w", err)
			}

			logger.Success("Tool installation completed successfully!")
		}
	} else {
		return fmt.Errorf("failed to get tool selector model")
	}

	// --- Placeholder for Future Phases ---
	logger.Info("Shell configuration and dotfiles setup coming soon.")

	logger.Success("Bootstrap CLI initialization finished!")
	return nil
}

// configureNeedrestart sets needrestart to automatic mode (requires root)
func configureNeedrestart() error {
	logger.Debug("Configuring needrestart for automatic restarts...")
	content := []byte(`$nrconf{restart} = 'a';`)
	filePath := "/etc/needrestart/conf.d/50-autorestart.conf"

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		// Don't fail the whole process if this write fails, just log it.
		logger.Warn("Failed to write needrestart config %s: %v", filePath, err)
		return nil // Return nil to indicate non-fatal error
	}
	logger.Debug("Successfully configured needrestart.")
	return nil
}

// Note: The Quitting() method needs to be implemented or accessible
// on the ui.ToolSelector struct for the cancellation check to work correctly.
// This might involve adding a method like `Finished() bool` to ui.ToolSelector
// that returns true if the UI exited normally (Enter on last category) vs. quit (Ctrl+C). 