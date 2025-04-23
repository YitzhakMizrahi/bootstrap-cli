package system

import (
	"fmt"
	"time"
)

// PackageManagerType represents the type of package manager
type PackageManagerType string

const (
	PackageManagerApt      PackageManagerType = "apt"
	PackageManagerDnf      PackageManagerType = "dnf"
	PackageManagerPacman   PackageManagerType = "pacman"
	PackageManagerHomebrew PackageManagerType = "homebrew"
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
func (f *PackageManagerFactory) GetPackageManager() (PackageManager, error) {
	info, err := Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect system: %w", err)
	}

	var pm PackageManager
	var pmErr error

	switch info.PackageType {
	case PackageManagerApt:
		pm, pmErr = NewAptPackageManager()
	case PackageManagerDnf:
		pm, pmErr = NewDnfPackageManager()
	case PackageManagerPacman:
		pm, pmErr = NewPacmanPackageManager()
	case PackageManagerHomebrew:
		pm, pmErr = NewHomebrewPackageManager()
	default:
		return nil, fmt.Errorf("unsupported package manager type: %s", info.PackageType)
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
	PackageManager
	maxRetries int
	retryDelay time.Duration
}

// Install with retry logic
func (r *retryPackageManager) Install(pkg string) error {
	var lastErr error
	for i := 0; i < r.maxRetries; i++ {
		if err := r.PackageManager.Install(pkg); err != nil {
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

// Uninstall with retry logic
func (r *retryPackageManager) Uninstall(pkg string) error {
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
	return fmt.Errorf("failed to uninstall package after %d retries: %w", r.maxRetries, lastErr)
} 