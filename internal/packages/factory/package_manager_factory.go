package factory

import (
	"fmt"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/detector"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages/implementations"
)

// PackageManagerFactory creates package managers based on system type
type PackageManagerFactory struct {
	maxRetries int
	retryDelay time.Duration
}

// NewPackageManagerFactory creates a new package manager factory
func NewPackageManagerFactory() *PackageManagerFactory {
	return &PackageManagerFactory{
		maxRetries: 3,
		retryDelay: 5 * time.Second,
	}
}

// SetRetryConfig configures retry behavior
func (f *PackageManagerFactory) SetRetryConfig(maxRetries int, retryDelay time.Duration) {
	f.maxRetries = maxRetries
	f.retryDelay = retryDelay
}

// GetPackageManager returns the appropriate package manager for the current system
func (f *PackageManagerFactory) GetPackageManager() (interfaces.PackageManager, error) {
	pmType, err := detector.DetectPackageManager()
	if err != nil {
		return nil, fmt.Errorf("failed to detect package manager: %w", err)
	}

	var pm interfaces.PackageManager
	var pmErr error

	switch pmType {
	case interfaces.APT:
		pm, pmErr = implementations.NewAptPackageManager()
	case interfaces.DNF:
		pm, pmErr = implementations.NewDnfPackageManager()
	case interfaces.Pacman:
		pm, pmErr = implementations.NewPacmanPackageManager()
	case interfaces.Homebrew:
		pm, pmErr = implementations.NewHomebrewPackageManager()
	default:
		return nil, fmt.Errorf("unsupported package manager type: %s", pmType)
	}

	if pmErr != nil {
		return nil, fmt.Errorf("failed to create package manager: %w", pmErr)
	}

	return &retryPackageManager{
		PackageManager: pm,
		maxRetries:    f.maxRetries,
		retryDelay:    f.retryDelay,
	}, nil
}

// retryPackageManager wraps a PackageManager with retry logic
type retryPackageManager struct {
	interfaces.PackageManager
	maxRetries int
	retryDelay time.Duration
}

// Install installs a package with retries
func (r *retryPackageManager) Install(packageName string) error {
	var lastErr error
	for i := 0; i < r.maxRetries; i++ {
		if err := r.PackageManager.Install(packageName); err != nil {
			lastErr = err
			if i < r.maxRetries-1 {
				time.Sleep(r.retryDelay)
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to install package after %d retries: %w", r.maxRetries, lastErr)
}

// Remove with retry logic
func (r *retryPackageManager) Remove(pkg string) error {
	var lastErr error
	for i := 0; i < r.maxRetries; i++ {
		if err := r.PackageManager.Uninstall(pkg); err != nil {
			lastErr = err
			if i < r.maxRetries-1 {
				time.Sleep(r.retryDelay)
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("failed to remove package after %d retries: %w", r.maxRetries, lastErr)
} 