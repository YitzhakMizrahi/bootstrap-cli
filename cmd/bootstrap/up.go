package main

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/shell"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/tools"
	"github.com/spf13/cobra"
)

// NewUpCmd creates a new up command
func NewUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run the complete bootstrap process",
		Long: `Run the complete bootstrap process including:
- System detection
- Shell detection and configuration
- Core tool installation
- Package management setup`,
		RunE: runUp,
	}
}

func runUp(cmd *cobra.Command, args []string) error {
	// Step 1: Detect system
	sysInfo, err := system.Detect()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}
	fmt.Printf("Detected system: %s %s (%s)\n", sysInfo.OS, sysInfo.Distro, sysInfo.Arch)

	// Step 2: Detect shell
	shellMgr := shell.NewManager()
	shellInfo, err := shellMgr.DetectCurrent()
	if err != nil {
		return fmt.Errorf("failed to detect shell: %w", err)
	}
	fmt.Printf("Detected shell: %s\n", shellInfo.Type)

	// Step 3: Install core tools
	logger := log.New(log.InfoLevel)
	pm, err := packages.NewManager(string(sysInfo.PackageType))
	if err != nil {
		return fmt.Errorf("failed to create package manager: %w", err)
	}
	
	if err := tools.InstallEssentialTools(pm, logger, false); err != nil {
		return fmt.Errorf("failed to install core tools: %w", err)
	}

	fmt.Println("Bootstrap process completed successfully!")
	return nil
} 