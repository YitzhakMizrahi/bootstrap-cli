package packages

import (
	"fmt"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/factory"
)

// ErrPackageManagerNotFound is returned when no suitable package manager is found
var ErrPackageManagerNotFound = fmt.Errorf("no suitable package manager found")

// GetPackageManager returns the appropriate package manager for the current system
func GetPackageManager() (interfaces.PackageManager, error) {
	f := factory.NewPackageManagerFactory()
	pm, err := f.GetPackageManager()
	if err != nil {
		return nil, err
	}
	return pm, nil
} 