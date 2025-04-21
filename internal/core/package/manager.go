package pkgmanager

import (
	"fmt"
	"os/exec"
	"strings"
)

// PackageManager defines the interface for package management operations
type PackageManager interface {
	// Name returns the name of the package manager
	Name() string
	
	// IsAvailable checks if the package manager is available on the system
	IsAvailable() bool
	
	// Install installs a package
	Install(packageName string) error
	
	// Uninstall removes a package
	Uninstall(packageName string) error
	
	// Update updates a package
	Update(packageName string) error
	
	// UpdateAll updates all packages
	UpdateAll() error
	
	// IsInstalled checks if a package is installed
	IsInstalled(packageName string) bool
	
	// Search searches for packages
	Search(query string) ([]string, error)
}

// NewPackageManager creates a new package manager based on the system
func NewPackageManager(name string) (PackageManager, error) {
	switch strings.ToLower(name) {
	case "apt":
		return &AptManager{}, nil
	case "dnf":
		return &DnfManager{}, nil
	case "pacman":
		return &PacmanManager{}, nil
	case "brew":
		return &BrewManager{}, nil
	case "choco":
		return &ChocoManager{}, nil
	default:
		return nil, fmt.Errorf("unsupported package manager: %s", name)
	}
}

// AptManager implements PackageManager for apt (Debian/Ubuntu)
type AptManager struct{}

func (m *AptManager) Name() string {
	return "apt"
}

func (m *AptManager) IsAvailable() bool {
	_, err := exec.LookPath("apt-get")
	return err == nil
}

func (m *AptManager) Install(packageName string) error {
	cmd := exec.Command("apt-get", "install", "-y", packageName)
	return cmd.Run()
}

func (m *AptManager) Uninstall(packageName string) error {
	cmd := exec.Command("apt-get", "remove", "-y", packageName)
	return cmd.Run()
}

func (m *AptManager) Update(packageName string) error {
	cmd := exec.Command("apt-get", "install", "--only-upgrade", "-y", packageName)
	return cmd.Run()
}

func (m *AptManager) UpdateAll() error {
	cmd := exec.Command("apt-get", "upgrade", "-y")
	return cmd.Run()
}

func (m *AptManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("dpkg", "-l", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), packageName)
}

func (m *AptManager) Search(query string) ([]string, error) {
	cmd := exec.Command("apt-cache", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	packages := make([]string, 0, len(lines))
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " - ", 2)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}
	
	return packages, nil
}

// DnfManager implements PackageManager for dnf (Fedora)
type DnfManager struct{}

func (m *DnfManager) Name() string {
	return "dnf"
}

func (m *DnfManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

func (m *DnfManager) Install(packageName string) error {
	cmd := exec.Command("dnf", "install", "-y", packageName)
	return cmd.Run()
}

func (m *DnfManager) Uninstall(packageName string) error {
	cmd := exec.Command("dnf", "remove", "-y", packageName)
	return cmd.Run()
}

func (m *DnfManager) Update(packageName string) error {
	cmd := exec.Command("dnf", "update", "-y", packageName)
	return cmd.Run()
}

func (m *DnfManager) UpdateAll() error {
	cmd := exec.Command("dnf", "update", "-y")
	return cmd.Run()
}

func (m *DnfManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("dnf", "list", "installed", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), packageName)
}

func (m *DnfManager) Search(query string) ([]string, error) {
	cmd := exec.Command("dnf", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	packages := make([]string, 0, len(lines))
	
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Last metadata expiration") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}
	
	return packages, nil
}

// PacmanManager implements PackageManager for pacman (Arch Linux)
type PacmanManager struct{}

func (m *PacmanManager) Name() string {
	return "pacman"
}

func (m *PacmanManager) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

func (m *PacmanManager) Install(packageName string) error {
	cmd := exec.Command("pacman", "-S", "--noconfirm", packageName)
	return cmd.Run()
}

func (m *PacmanManager) Uninstall(packageName string) error {
	cmd := exec.Command("pacman", "-R", "--noconfirm", packageName)
	return cmd.Run()
}

func (m *PacmanManager) Update(packageName string) error {
	cmd := exec.Command("pacman", "-S", "--noconfirm", packageName)
	return cmd.Run()
}

func (m *PacmanManager) UpdateAll() error {
	cmd := exec.Command("pacman", "-Syu", "--noconfirm")
	return cmd.Run()
}

func (m *PacmanManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("pacman", "-Q", packageName)
	err := cmd.Run()
	return err == nil
}

func (m *PacmanManager) Search(query string) ([]string, error) {
	cmd := exec.Command("pacman", "-Ss", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	packages := make([]string, 0, len(lines))
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}
	
	return packages, nil
}

// BrewManager implements PackageManager for Homebrew (macOS)
type BrewManager struct{}

func (m *BrewManager) Name() string {
	return "brew"
}

func (m *BrewManager) IsAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (m *BrewManager) Install(packageName string) error {
	cmd := exec.Command("brew", "install", packageName)
	return cmd.Run()
}

func (m *BrewManager) Uninstall(packageName string) error {
	cmd := exec.Command("brew", "uninstall", packageName)
	return cmd.Run()
}

func (m *BrewManager) Update(packageName string) error {
	cmd := exec.Command("brew", "upgrade", packageName)
	return cmd.Run()
}

func (m *BrewManager) UpdateAll() error {
	cmd := exec.Command("brew", "upgrade")
	return cmd.Run()
}

func (m *BrewManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("brew", "list", packageName)
	err := cmd.Run()
	return err == nil
}

func (m *BrewManager) Search(query string) ([]string, error) {
	cmd := exec.Command("brew", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	packages := make([]string, 0, len(lines))
	
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "==>") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}
	
	return packages, nil
}

// ChocoManager implements PackageManager for Chocolatey (Windows)
type ChocoManager struct{}

func (m *ChocoManager) Name() string {
	return "choco"
}

func (m *ChocoManager) IsAvailable() bool {
	_, err := exec.LookPath("choco")
	return err == nil
}

func (m *ChocoManager) Install(packageName string) error {
	cmd := exec.Command("choco", "install", "-y", packageName)
	return cmd.Run()
}

func (m *ChocoManager) Uninstall(packageName string) error {
	cmd := exec.Command("choco", "uninstall", "-y", packageName)
	return cmd.Run()
}

func (m *ChocoManager) Update(packageName string) error {
	cmd := exec.Command("choco", "upgrade", "-y", packageName)
	return cmd.Run()
}

func (m *ChocoManager) UpdateAll() error {
	cmd := exec.Command("choco", "upgrade", "-y", "all")
	return cmd.Run()
}

func (m *ChocoManager) IsInstalled(packageName string) bool {
	cmd := exec.Command("choco", "list", "--local-only", packageName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), packageName)
}

func (m *ChocoManager) Search(query string) ([]string, error) {
	cmd := exec.Command("choco", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	packages := make([]string, 0, len(lines))
	
	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Chocolatey") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			packages = append(packages, parts[0])
		}
	}
	
	return packages, nil
} 