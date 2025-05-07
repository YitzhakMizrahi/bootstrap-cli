package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// GenerateFontInstallSteps creates pipeline steps for installing a font based on commands.
func GenerateFontInstallSteps(font *interfaces.Font, platform *Platform) []InstallationStep {
	steps := []InstallationStep{}
	if font == nil {
		fmt.Printf("Skipping font installation: Font data is nil.\n")
		return steps
	}

	// Determine target directory based on OS (might be needed by install commands via env var?)
	homeDir, _ := os.UserHomeDir()
	var targetDir string
	if platform.OS == "darwin" {
		targetDir = filepath.Join(homeDir, "Library", "Fonts")
	} else {
		targetDir = filepath.Join(homeDir, ".local", "share", "fonts")
	}

	// Step 1: Ensure target directory exists
	steps = append(steps, InstallationStep{
		Name:        fmt.Sprintf("ensure-font-dir-%s", font.Name),
		Description: fmt.Sprintf("Ensuring font directory exists (%s)", targetDir),
		Action: func(ctx *InstallationContext) error {
			ctx.sendProgress(TaskLog{TaskID: ctx.State.CurrentStep, Line: fmt.Sprintf("Ensuring directory %s exists", targetDir)})
			return os.MkdirAll(targetDir, 0755)
		},
	})

	// Step 2: Run Install Commands
	for i, cmdStr := range font.Install {
		installCmdStr := cmdStr 
		stepName := fmt.Sprintf("install-font-%s-step%d", font.Name, i)
		steps = append(steps, InstallationStep{
			Name:        stepName,
			Description: fmt.Sprintf("Running font install command: %s", installCmdStr),
			Action: func(ctx *InstallationContext) error {
				ctx.sendProgress(TaskLog{TaskID: stepName, Line: fmt.Sprintf("Executing: %s", installCmdStr)})
				cmd := exec.Command("sh", "-c", installCmdStr)
				// TODO: Capture live output -> TaskLog
				output, err := cmd.CombinedOutput()
				if len(output) > 0 {
					ctx.sendProgress(TaskLog{TaskID: stepName, Line: string(output)})
				}
				if err != nil {
					return fmt.Errorf("font install command failed: %w", err)
				}
				return nil
			},
			Timeout: 5 * time.Minute,
		})
	}

	// Step 3: Run Verify Commands
	for i, cmdStr := range font.Verify {
		verifyCmdStr := cmdStr
		stepName := fmt.Sprintf("verify-font-%s-step%d", font.Name, i)
		steps = append(steps, InstallationStep{
			Name:        stepName,
			Description: fmt.Sprintf("Running font verify command: %s", verifyCmdStr),
			Action: func(ctx *InstallationContext) error {
				ctx.sendProgress(TaskLog{TaskID: stepName, Line: fmt.Sprintf("Verifying: %s", verifyCmdStr)})
				cmd := exec.Command("sh", "-c", verifyCmdStr)
				// TODO: Capture live output -> TaskLog
				err := cmd.Run() // CombinedOutput might be better
				if err != nil {
					return fmt.Errorf("font verify command failed: %w", err)
				}
				return nil
			},
			Timeout: 1 * time.Minute,
		})
	}

	return steps
}

// copyFile is likely no longer needed here if using command-based install
// func copyFile(src, dst string) error { ... } 