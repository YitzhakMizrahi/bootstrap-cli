package pipeline

import (
	"fmt"
	"sync"
)

// InstallationState tracks the state of package installations
type InstallationState struct {
	mu        sync.RWMutex
	Installed map[string]bool
	Failed    map[string]error
	Pending   map[string]bool
}

// NewInstallationState creates a new installation state tracker
func NewInstallationState() *InstallationState {
	return &InstallationState{
		Installed: make(map[string]bool),
		Failed:    make(map[string]error),
		Pending:   make(map[string]bool),
	}
}

// MarkInstalled marks a package as successfully installed
func (s *InstallationState) MarkInstalled(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Installed[name] = true
	delete(s.Pending, name)
	delete(s.Failed, name)
}

// MarkFailed marks a package as failed to install
func (s *InstallationState) MarkFailed(name string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Failed[name] = err
	delete(s.Pending, name)
}

// MarkPending marks a package as pending installation
func (s *InstallationState) MarkPending(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Pending[name] = true
}

// IsInstalled checks if a package is installed
func (s *InstallationState) IsInstalled(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Installed[name]
}

// IsFailed checks if a package failed to install
func (s *InstallationState) IsFailed(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, failed := s.Failed[name]
	return failed
}

// IsPending checks if a package is pending installation
func (s *InstallationState) IsPending(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Pending[name]
}

// GetFailedError returns the error for a failed package
func (s *InstallationState) GetFailedError(name string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Failed[name]
}

// Status returns the current status of a package
func (s *InstallationState) Status(name string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.Installed[name] {
		return "installed"
	}
	if s.Failed[name] != nil {
		return "failed"
	}
	if s.Pending[name] {
		return "pending"
	}
	return "unknown"
}

// String returns a string representation of the installation state
func (s *InstallationState) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprintf("Installed: %d, Failed: %d, Pending: %d",
		len(s.Installed), len(s.Failed), len(s.Pending))
} 