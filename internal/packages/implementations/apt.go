package implementations

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// APTManager implements the PackageManager interface for APT-based systems
type APTManager struct {
	aptGetPath string
	aptPath    string
}

// NewAptPackageManager creates a new APT package manager instance
func NewAptPackageManager() (interfaces.PackageManager, error) {
	// Check for apt-get
	aptGetPath, err := exec.LookPath("apt-get")
	if err != nil {
		return nil, fmt.Errorf("apt-get is required but not found: %w", err)
	}

	// Check for apt
	aptPath, err := exec.LookPath("apt")
	if err != nil {
		return nil, fmt.Errorf("apt is required but not found: %w", err)
	}

	return &APTManager{
		aptGetPath: aptGetPath,
		aptPath:    aptPath,
	}, nil
}

// Name returns the name of the package manager
func (a *APTManager) Name() string {
	return "apt"
}

// IsAvailable checks if apt-get is available on the system
func (a *APTManager) IsAvailable() bool {
	_, err := exec.LookPath(a.aptGetPath)
	return err == nil
}

// Update updates the package list
func (a *APTManager) Update() error {
	cmd := exec.Command(a.aptGetPath, "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// checkPPAExists checks if a PPA exists before trying to add it
func (a *APTManager) checkPPAExists(ppa string) bool {
	cmd := exec.Command("add-apt-repository", "-n", ppa)
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

// addRepository adds a third-party repository
func (a *APTManager) addRepository(repo string) error {
	// Install prerequisites if needed
	if err := a.installPrerequisites(); err != nil {
		return fmt.Errorf("failed to install prerequisites: %w", err)
	}

	// For PPAs, check if they exist first
	if strings.HasPrefix(repo, "ppa:") {
		if !a.checkPPAExists(repo) {
			return fmt.Errorf("PPA %s does not exist", repo)
		}
	}

	// Add the repository using add-apt-repository
	cmd := exec.Command("add-apt-repository", "-y", repo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add repository: %w", err)
	}

	// Update package list
	return a.Update()
}

// installPrerequisites installs necessary packages for adding repositories
func (a *APTManager) installPrerequisites() error {
	prerequisites := []string{"software-properties-common", "curl", "gnupg"}
	for _, pkg := range prerequisites {
		if !a.IsInstalled(pkg) {
			if err := a.Install(pkg); err != nil {
				return fmt.Errorf("failed to install prerequisite %s: %w", pkg, err)
			}
		}
	}
	return nil
}

// SetupSpecialPackage handles special package installations that require repository setup
func (a *APTManager) SetupSpecialPackage(pkg string) error {
	// No special setup needed for lsd and bat anymore as we're using direct package installation
	// This method is kept for other packages that might need special repository setup
	return nil
}

// Install installs a package using apt
func (a *APTManager) Install(pkg string) error {
	cmd := exec.Command("sudo", "apt-get", "install", "-y", pkg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install package %s: %v\nOutput: %s", pkg, err, output)
	}
	return nil
}

// Remove removes a package
func (a *APTManager) Remove(packageName string) error {
	cmd := exec.Command(a.aptGetPath, "remove", "-y", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsInstalled checks if a package is installed
func (a *APTManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), packageName)
}

// GetVersion returns the version of an installed package
func (a *APTManager) GetVersion(packageName string) (string, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Version}", packageName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// ListInstalled returns a list of installed packages
func (a *APTManager) ListInstalled() ([]string, error) {
	cmd := exec.Command("dpkg-query", "-W", "-f=${Package}\n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

// GetName returns the name of the package manager
func (a *APTManager) GetName() string {
	return string(interfaces.APT)
}

// Upgrade upgrades all packages
func (a *APTManager) Upgrade() error {
	cmd := exec.Command("sudo", a.aptGetPath, "upgrade", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// IsPackageAvailable checks if a package is available in the package manager's repositories
func (a *APTManager) IsPackageAvailable(pkg string) bool {
	cmd := exec.Command("apt-cache", "show", pkg)
	err := cmd.Run()
	return err == nil
}