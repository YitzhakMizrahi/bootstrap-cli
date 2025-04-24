package pipeline

// ShellConfig represents shell configuration for a tool
type ShellConfig struct {
	Aliases   map[string]string
	Functions map[string]string
	Env       map[string]string
}

// PostInstallCommand represents a command to run after installation
type PostInstallCommand struct {
	Command     string
	Description string
}

// PackageManager represents a package manager interface
type PackageManager interface {
	// Install installs a package
	Install(pkg string) error
	// Uninstall uninstalls a package
	Uninstall(pkg string) error
	// IsInstalled checks if a package is installed
	IsInstalled(pkg string) (bool, error)
	// Update updates the package list
	Update() error
	// SetupSpecialPackage handles special package installation requirements
	SetupSpecialPackage(pkg string) error
	// IsPackageAvailable checks if a package is available in the package manager's repositories
	IsPackageAvailable(pkg string) bool
} 