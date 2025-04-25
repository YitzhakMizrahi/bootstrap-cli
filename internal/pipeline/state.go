package pipeline

import (
	"fmt"
	"sync"
	"time"
)

// InstallationState tracks the state of the installation process
type InstallationState struct {
	mu            sync.Mutex
	CurrentStep   string
	CompletedSteps []string
	FailedSteps   []string
	RollbackSteps []string
	Status        string
	Error         error
	StartTime     time.Time
	LastUpdated   time.Time
}

// NewInstallationState creates a new installation state
func NewInstallationState() *InstallationState {
	return &InstallationState{
		Status:      "initialized",
		StartTime:   time.Now(),
		LastUpdated: time.Now(),
	}
}

// UpdateState updates the installation state
func (s *InstallationState) UpdateState(step string, status string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.CurrentStep = step
	s.Status = status
	s.Error = err
	s.LastUpdated = time.Now()
	
	switch status {
	case "completed":
		s.CompletedSteps = append(s.CompletedSteps, step)
	case "failed":
		s.FailedSteps = append(s.FailedSteps, step)
	case "rollback":
		s.RollbackSteps = append(s.RollbackSteps, step)
	}
}

// GetProgress returns the current progress as a percentage
func (s *InstallationState) GetProgress() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	totalSteps := len(s.CompletedSteps) + len(s.FailedSteps)
	if totalSteps == 0 {
		return 0
	}
	
	return float64(len(s.CompletedSteps)) / float64(totalSteps) * 100
}

// String returns a string representation of the state
func (s *InstallationState) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return fmt.Sprintf("Status: %s, Current Step: %s, Progress: %.1f%%", 
		s.Status, s.CurrentStep, s.GetProgress())
}

// GetDuration returns the duration of the installation
func (s *InstallationState) GetDuration() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return time.Since(s.StartTime)
}

// IsComplete returns true if the installation is complete
func (s *InstallationState) IsComplete() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.Status == "completed" || s.Status == "failed"
}

// HasFailed returns true if the installation has failed
func (s *InstallationState) HasFailed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.Status == "failed"
}

// GetFailedSteps returns the list of failed steps
func (s *InstallationState) GetFailedSteps() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.FailedSteps
}

// GetRollbackSteps returns the list of steps that were rolled back
func (s *InstallationState) GetRollbackSteps() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	return s.RollbackSteps
} 