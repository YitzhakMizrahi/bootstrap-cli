package platform

// Platform defines the interface for platform-specific operations
type Platform interface {
	// Package management
	InstallPackage(name string) error
	UninstallPackage(name string) error
	IsPackageInstalled(name string) bool
	GetPackageManager() string
	UpdatePackages() error
	UpgradePackages() error

	// Shell operations
	GetDefaultShell() string
	GetShellPath(shell string) string
	IsShellAvailable(shell string) bool
	
	// Path operations
	GetHomeDir() string
	GetConfigDir() string
	GetCacheDir() string
	
	// System operations
	IsRoot() bool
	ElevatePrivileges() error
	GetOSInfo() *OSInfo
}

// OSInfo contains information about the operating system
type OSInfo struct {
	Type    string // "linux", "darwin", "windows"
	Version string // OS version
	Arch    string // CPU architecture
}

// Provider provides platform-specific implementations
type Provider interface {
	// GetPlatform returns the platform-specific implementation
	GetPlatform() Platform
}

var currentPlatform Platform

// GetCurrentPlatform returns the current platform implementation
func GetCurrentPlatform() Platform {
	return currentPlatform
}

// SetCurrentPlatform sets the current platform implementation
func SetCurrentPlatform(p Platform) {
	currentPlatform = p
} 