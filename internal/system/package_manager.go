package system

// PackageManager defines the interface for system package managers
type PackageManager interface {
	// Install installs a package using sudo privileges
	Install(pkg string) error
	
	// IsInstalled checks if a package is installed
	IsInstalled(pkg string) bool
	
	// Uninstall removes a package using sudo privileges
	Uninstall(pkg string) error
}

// GetPackageManager returns the appropriate package manager for the current system
func GetPackageManager() (PackageManager, error) {
	// For now we only support apt
	return NewAptPackageManager()
} 