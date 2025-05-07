package pipeline

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// GenerateLanguageInstallSteps creates pipeline steps for installing a language.
func GenerateLanguageInstallSteps(lang *interfaces.Language, context *InstallationContext) []InstallationStep {
	steps := []InstallationStep{}
	if lang == nil {
		fmt.Printf("Skipping language installation: Language data is nil.\n")
		return steps
	}

	// TODO: Determine installation strategy (e.g., use version manager like pyenv/nvm if specified and available, otherwise use system PM)
	// This logic needs access to the InstallationContext to check for installed tools (version managers) and system PM.
	
	// --- Placeholder: Simple system package manager install --- 
	// This assumes the language name directly maps to a package name.
	// A real implementation would use lang.Installer, lang.Version, lang.PackageNames etc.
	pkgName := lang.Name // Very naive assumption
	pkgManagerName := context.Platform.PackageManager
	var installCmdStr string
	switch pkgManagerName {
	case "apt":
		installCmdStr = fmt.Sprintf("sudo apt-get install -y %s", pkgName)
	case "brew":
		installCmdStr = fmt.Sprintf("brew install %s", pkgName)
	case "dnf":
		installCmdStr = fmt.Sprintf("sudo dnf install -y %s", pkgName)
	case "pacman":
		installCmdStr = fmt.Sprintf("sudo pacman -S --noconfirm %s", pkgName)
	default:
		// Return an error step? Log a warning?
		fmt.Printf("Unsupported package manager '%s' for language %s install\n", pkgManagerName, lang.Name)
		return steps
	}

	steps = append(steps, InstallationStep{
		Name:        fmt.Sprintf("install-lang-%s", lang.Name),
		Description: fmt.Sprintf("Installing language %s using %s", lang.Name, pkgManagerName),
		Action: func(ctx *InstallationContext) error {
			// TODO: Add logging via ctx.Logger or ctx.sendProgress
			cmd := exec.Command("sh", "-c", installCmdStr)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		},
		Timeout: 5 * time.Minute,
	})
	// --- End Placeholder ---

	// TODO: Add verification steps based on lang.Verify

	return steps
} 