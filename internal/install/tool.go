package install

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// PackageMapping defines package names for different package managers
type PackageMapping struct {
	Default string
	APT     string
	DNF     string
	Pacman  string
	Brew    string
}

// Tool represents a development tool that can be installed
type Tool struct {
	// Name is the name of the tool
	Name string
	// PackageName is the name of the package in the package manager
	PackageName string
	// PackageNames contains system-specific package names
	PackageNames *PackageMapping
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
	// Logger is the logger to use
	Logger *log.Logger
	// MaxRetries is the maximum number of retries for failed operations
	MaxRetries int
	// RetryDelay is the delay between retries
	RetryDelay time.Duration
}

// NewInstaller creates a new installer with the given package manager
func NewInstaller(pm packages.PackageManager) *Installer {
	return &Installer{
		PackageManager: pm,
		Logger:        log.New(log.InfoLevel),
		MaxRetries:    3,
		RetryDelay:    time.Second * 2,
	}
}

// getPackageWithVersion returns the package name with version if specified
func (i *Installer) getPackageWithVersion(pkg, version string) string {
	if version == "" || version == "latest" || version == "stable" {
		return pkg
	}
	
	switch i.PackageManager.Name() {
	case "apt":
		return fmt.Sprintf("%s=%s", pkg, version)
	case "dnf":
		return fmt.Sprintf("%s-%s", pkg, version)
	case "pacman":
		return fmt.Sprintf("%s=%s", pkg, version)
	case "brew":
		return fmt.Sprintf("%s@%s", pkg, version)
	default:
		return pkg
	}
}

// getSystemPackageName returns the appropriate package name for the current system
func (i *Installer) getSystemPackageName(tool *Tool) string {
	if tool.PackageNames == nil {
		return tool.PackageName
	}

	switch i.PackageManager.Name() {
	case "apt":
		if tool.PackageNames.APT != "" {
			return tool.PackageNames.APT
		}
	case "dnf":
		if tool.PackageNames.DNF != "" {
			return tool.PackageNames.DNF
		}
	case "pacman":
		if tool.PackageNames.Pacman != "" {
			return tool.PackageNames.Pacman
		}
	case "brew":
		if tool.PackageNames.Brew != "" {
			return tool.PackageNames.Brew
		}
	}

	if tool.PackageNames.Default != "" {
		return tool.PackageNames.Default
	}
	return tool.PackageName
}

// Install installs a tool and its dependencies
func (i *Installer) Install(tool *Tool) error {
	i.Logger.Info("Starting installation of %s", tool.Name)

	// Install dependencies first
	for _, dep := range tool.Dependencies {
		if !i.PackageManager.IsInstalled(dep) {
			i.Logger.Info("Installing dependency %s for %s", dep, tool.Name)
			// Dependencies are always installed with latest version
			if err := i.installWithRetry(dep); err != nil {
				i.Logger.Error("Failed to install dependency %s: %v", dep, err)
				return &InstallError{
					Tool:    tool.Name,
					Phase:   "dependencies",
					Message: fmt.Sprintf("failed to install dependency %s", dep),
					Err:     err,
				}
			}
		} else {
			i.Logger.Debug("Dependency %s is already installed", dep)
		}
	}

	// Get system-specific package name and version
	pkgName := i.getSystemPackageName(tool)
	pkgWithVersion := i.getPackageWithVersion(pkgName, tool.Version)

	// Install the tool
	if !i.PackageManager.IsInstalled(pkgName) {
		i.Logger.Info("Installing %s version %s", tool.Name, tool.Version)
		if err := i.installWithRetry(pkgWithVersion); err != nil {
			i.Logger.Error("Failed to install %s: %v", tool.Name, err)
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "installation",
				Message: "failed to install package",
				Err:     err,
			}
		}
	} else {
		i.Logger.Info("%s is already installed", tool.Name)
	}

	// Run post-install commands
	if err := i.runPostInstall(tool); err != nil {
		i.Logger.Error("Failed to run post-install commands for %s: %v", tool.Name, err)
		return &InstallError{
			Tool:    tool.Name,
			Phase:   "post-install",
			Message: "failed to run post-install commands",
			Err:     err,
		}
	}

	// Verify installation
	if tool.VerifyCommand != "" {
		i.Logger.Info("Verifying installation of %s", tool.Name)
		if err := i.verifyInstallation(tool); err != nil {
			i.Logger.Error("Installation verification failed for %s: %v", tool.Name, err)
			return &InstallError{
				Tool:    tool.Name,
				Phase:   "verification",
				Message: "verification failed",
				Err:     err,
			}
		}
	}

	i.Logger.Info("Successfully installed %s", tool.Name)
	return nil
}

// installWithRetry attempts to install a package with retries
func (i *Installer) installWithRetry(pkg string) error {
	var lastErr error
	for attempt := 0; attempt <= i.MaxRetries; attempt++ {
		if attempt > 0 {
			i.Logger.Debug("Retrying installation of %s (attempt %d/%d)", pkg, attempt, i.MaxRetries)
			time.Sleep(i.RetryDelay)
		}

		// Update package lists before retry
		if attempt > 0 {
			if err := i.PackageManager.Update(); err != nil {
				i.Logger.Debug("Failed to update package lists: %v", err)
			}
		}

		if err := i.PackageManager.Install(pkg); err != nil {
			lastErr = err
			i.Logger.Debug("Installation attempt failed: %v", err)
			continue
		}
		return nil
	}
	return fmt.Errorf("failed to install after %d attempts: %w", i.MaxRetries, lastErr)
}

// runPostInstall executes post-installation commands
func (i *Installer) runPostInstall(tool *Tool) error {
	for _, cmd := range tool.PostInstall {
		i.Logger.Debug("Running post-install command: %s", cmd)
		command := exec.Command("sh", "-c", cmd)
		if output, err := command.CombinedOutput(); err != nil {
			return fmt.Errorf("post-install command failed: %w\nOutput: %s", err, string(output))
		}
	}
	return nil
}

// verifyInstallation checks if the tool was installed correctly
func (i *Installer) verifyInstallation(tool *Tool) error {
	i.Logger.Debug("Verifying installation of %s", tool.Name)
	cmd := exec.Command("sh", "-c", tool.VerifyCommand)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("verification failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

// retryOperation retries an operation with exponential backoff
func (i *Installer) retryOperation(op func() error) error {
	var lastErr error
	for attempt := 0; attempt < i.MaxRetries; attempt++ {
		if attempt > 0 {
			i.Logger.Debug("Retrying operation (attempt %d/%d)", attempt+1, i.MaxRetries)
			time.Sleep(i.RetryDelay * time.Duration(attempt))
		}

		if err := op(); err != nil {
			lastErr = err
			i.Logger.Warn("Operation failed (attempt %d/%d): %v", attempt+1, i.MaxRetries, err)
			continue
		}

		return nil
	}

	return lastErr
}

// runCommand executes a shell command
func (i *Installer) runCommand(cmd string) error {
	i.Logger.Debug("Running command: %s", cmd)
	shell := exec.Command("sh", "-c", cmd)
	shell.Stdout = os.Stdout
	shell.Stderr = os.Stderr

	return shell.Run()
} 