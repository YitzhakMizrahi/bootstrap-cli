package pipeline

import (
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// PackageManagerAdapter adapts interfaces.PackageManager to pipeline.PackageManager
type PackageManagerAdapter struct {
	pm interfaces.PackageManager
}

// NewPackageManagerAdapter creates a new adapter
func NewPackageManagerAdapter(pm interfaces.PackageManager) *PackageManagerAdapter {
	return &PackageManagerAdapter{pm: pm}
}

// Install installs a package
func (a *PackageManagerAdapter) Install(pkg string) error {
	return a.pm.Install(pkg)
}

// Uninstall uninstalls a package
func (a *PackageManagerAdapter) Uninstall(pkg string) error {
	return a.pm.Remove(pkg)
}

// IsInstalled checks if a package is installed
func (a *PackageManagerAdapter) IsInstalled(pkg string) (bool, error) {
	return a.pm.IsInstalled(pkg), nil
}

// Update updates the package list
func (a *PackageManagerAdapter) Update() error {
	return a.pm.Update()
}

// SetupSpecialPackage handles special package installation requirements
func (a *PackageManagerAdapter) SetupSpecialPackage(_ string) error {
	// This is a no-op for the simple package manager
	return nil
}

// IsPackageAvailable checks if a package is available
func (a *PackageManagerAdapter) IsPackageAvailable(_ string) bool {
	// For the simple package manager, we assume all packages are available
	return true
} 