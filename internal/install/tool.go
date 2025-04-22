package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// Tool represents a development tool that can be installed
type Tool struct {
	// Name is the name of the tool
	Name string
	// PackageName is the name of the package in the package manager
	PackageName string
	// Version is the desired version of the tool
	Version string
	// Dependencies is a list of package names that this tool depends on
	Dependencies []string
	// PostInstall is a list of commands to run after installation
	PostInstall []string
	// VerifyCommand is the command to verify the installation
	VerifyCommand string
}

// Installer handles tool installation
type Installer struct {
	// PackageManager is the package manager to use
	PackageManager packages.PackageManager
}

// NewInstaller creates a new installer with the given package manager
func NewInstaller(pm packages.PackageManager) *Installer {
	return &Installer{
		PackageManager: pm,
	}
}

// Install installs a tool and its dependencies
func (i *Installer) Install(tool *Tool) error {
	// Install dependencies first
	for _, dep := range tool.Dependencies {
		if !i.PackageManager.IsInstalled(dep) {
			if err := i.PackageManager.Install(dep); err != nil {
				return fmt.Errorf("failed to install dependency %s: %w", dep, err)
			}
		}
	}

	// Install the tool
	if !i.PackageManager.IsInstalled(tool.PackageName) {
		if err := i.PackageManager.Install(tool.PackageName); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool.Name, err)
		}
	}

	// Run post-install commands
	for _, cmd := range tool.PostInstall {
		if err := i.runCommand(cmd); err != nil {
			return fmt.Errorf("failed to run post-install command for %s: %w", tool.Name, err)
		}
	}

	// Verify installation
	if tool.VerifyCommand != "" {
		if err := i.runCommand(tool.VerifyCommand); err != nil {
			return fmt.Errorf("failed to verify installation of %s: %w", tool.Name, err)
		}
	}

	return nil
}

// runCommand executes a shell command
func (i *Installer) runCommand(cmd string) error {
	// Create a shell command
	shell := exec.Command("sh", "-c", cmd)
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	// Run the command
	return shell.Run()
} 