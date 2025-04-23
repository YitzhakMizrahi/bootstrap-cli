package interfaces

// PackageManager defines the interface for package management operations
type PackageManager interface {
	// Name returns the name of the package manager
	Name() string
	// IsAvailable checks if the package manager is available on the system
	IsAvailable() bool
	// Install installs the given packages
	Install(packages ...string) error
	// Update updates the package list
	Update() error
	// IsInstalled checks if a package is installed
	IsInstalled(pkg string) bool
	// Remove removes a package
	Remove(pkg string) error
	// GetVersion returns the version of a package
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