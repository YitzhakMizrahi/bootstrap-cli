package platform

// OS represents the operating system type
type OS string

const (
	Linux   OS = "linux"
	Darwin  OS = "darwin"
	Windows OS = "windows"
	MacOS   OS = Darwin  // Alias for Darwin
)

// PackageManager represents a package manager type
type PackageManager string

const (
	Apt      PackageManager = "apt"
	Yum      PackageManager = "yum"
	Dnf      PackageManager = "dnf"
	Pacman   PackageManager = "pacman"
	Homebrew PackageManager = "homebrew"
)

// Info contains information about the current platform
type Info struct {
	OS              OS              // The operating system
	Distribution    string          // Linux distribution name (if applicable)
	PackageManagers []PackageManager // Available package managers
}

// Detector interface for platform detection
type Detector interface {
	DetectOS() OS
	DetectPackageManagers() []PackageManager
	DetectShell() string
	Detect() (Info, error)
	GetPrimaryPackageManager(info Info) (PackageManager, error)
}

// Contains checks if a slice contains a value
func Contains(slice []PackageManager, item PackageManager) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 