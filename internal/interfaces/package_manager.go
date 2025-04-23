package interfaces

// PackageManager represents a system package manager
type PackageManager interface {
	// Install installs a package
	Install(packageName string) error
	
	// IsInstalled checks if a package is installed
	IsInstalled(packageName string) bool
	
	// GetName returns the name of the package manager (apt, brew, dnf, pacman)
	GetName() string
	
	// IsAvailable checks if the package manager is available on the system
	IsAvailable() bool
	
	// Update updates the package list
	Update() error
	
	// Upgrade upgrades all packages
	Upgrade() error

	// Remove removes a package
	Remove(packageName string) error

	// GetVersion returns the version of an installed package
	GetVersion(packageName string) (string, error)

	// ListInstalled returns a list of installed packages
	ListInstalled() ([]string, error)
}

// PackageManagerType represents the type of package manager
type PackageManagerType string

const (
	// APT package manager (Debian/Ubuntu)
	APT PackageManagerType = "apt"
	// DNF package manager (Fedora)
	DNF PackageManagerType = "dnf"
	// Pacman package manager (Arch)
	Pacman PackageManagerType = "pacman"
	// Homebrew package manager (macOS)
	Homebrew PackageManagerType = "brew"
) 