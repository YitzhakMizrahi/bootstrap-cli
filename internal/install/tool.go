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

// InstallError represents an installation error
type InstallError struct {
	Tool    string
	Phase   string
	Message string
	Err     error
}

func (e *InstallError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s failed: %s (%v)", e.Tool, e.Phase, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s failed: %s", e.Tool, e.Phase, e.Message)
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
	// Validate the tool configuration
	if err := validateTool(tool); err != nil {
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "validation",
			Message: "invalid tool configuration",
			Err:     err,
		}
	}

	// Install dependencies first
	for _, dep := range tool.Dependencies {
		if !i.PackageManager.IsInstalled(dep) {
			if err := i.PackageManager.Install(dep); err != nil {
				return &InstallError{
					Tool:    tool.Name,
					Phase:   "dependencies",
					Message: fmt.Sprintf("failed to install dependency %s", dep),
					Err:     err,
				}
			}
		}
	}

	// Install the tool
	if !i.PackageManager.IsInstalled(tool.PackageName) {
		if err := i.PackageManager.Install(tool.PackageName); err != nil {
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "installation",
				Message: "failed to install package",
				Err:     err,
			}
		}
	}

	// Run post-install commands
	for _, cmd := range tool.PostInstall {
		if err := i.runCommand(cmd); err != nil {
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "post-install",
				Message: fmt.Sprintf("command failed: %s", cmd),
				Err:     err,
			}
		}
	}

	// Verify installation
	if tool.VerifyCommand != "" {
		if err := i.runCommand(tool.VerifyCommand); err != nil {
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "verification",
				Message: "verification command failed",
				Err:     err,
			}
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