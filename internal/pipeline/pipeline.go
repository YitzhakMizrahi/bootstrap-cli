package pipeline

import (
	"fmt"
	"log"
	"time"
)

// InstallationStep represents a single step in the installation pipeline
type InstallationStep struct {
	Name        string
	Description string
	Action      func() error
	Rollback    func() error
	Timeout     time.Duration
	RetryCount  int
	RetryDelay  time.Duration
}

// InstallationPipeline represents a sequence of installation steps
type InstallationPipeline struct {
	Steps   []InstallationStep
	State   *InstallationState
	Logger  *log.Logger
}

// NewInstallationPipeline creates a new installation pipeline
func NewInstallationPipeline() *InstallationPipeline {
	return &InstallationPipeline{
		Steps: make([]InstallationStep, 0),
		State: NewInstallationState(),
	}
}

// AddStep adds a step to the pipeline with default values
func (p *InstallationPipeline) AddStep(step InstallationStep) {
	if step.Timeout == 0 {
		step.Timeout = 5 * time.Minute
	}
	if step.RetryCount == 0 {
		step.RetryCount = 3
	}
	if step.RetryDelay == 0 {
		step.RetryDelay = 5 * time.Second
	}
	p.Steps = append(p.Steps, step)
}

// Execute runs all steps in the pipeline
func (p *InstallationPipeline) Execute() error {
	for i, step := range p.Steps {
		p.State.UpdateState(step.Name, "running", nil)
		
		// Execute step with retry
		err := p.executeStepWithRetry(step)
		if err != nil {
			p.State.UpdateState(step.Name, "failed", err)
			
			// Attempt rollback of completed steps
			if err := p.rollback(i); err != nil {
				return fmt.Errorf("step '%s' failed: %v, rollback also failed: %v", 
					step.Name, err, err)
			}
			
			return fmt.Errorf("step '%s' failed: %v", step.Name, err)
		}
		
		p.State.UpdateState(step.Name, "completed", nil)
	}
	
	p.State.UpdateState("pipeline", "completed", nil)
	return nil
}

// executeStepWithRetry executes a step with retry logic
func (p *InstallationPipeline) executeStepWithRetry(step InstallationStep) error {
	var lastErr error
	
	for attempt := 0; attempt <= step.RetryCount; attempt++ {
		if attempt > 0 {
			p.State.UpdateState(step.Name, "retrying", lastErr)
			time.Sleep(step.RetryDelay)
		}
		
		// Execute step
		err := step.Action()
		if err == nil {
			return nil
		}
		
		lastErr = err
		
		// Check if error is retryable
		if !isRetryableError(err) {
			return err
		}
	}
	
	return fmt.Errorf("failed after %d attempts: %v", step.RetryCount, lastErr)
}

// rollback attempts to roll back completed steps in reverse order
func (p *InstallationPipeline) rollback(lastCompletedIndex int) error {
	for i := lastCompletedIndex; i >= 0; i-- {
		step := p.Steps[i]
		p.State.UpdateState(step.Name, "rolling_back", nil)
		
		// Execute rollback action if defined
		if step.Rollback != nil {
			if err := step.Rollback(); err != nil {
				return fmt.Errorf("rollback failed for step '%s': %v", step.Name, err)
			}
		}
		
		p.State.UpdateState(step.Name, "rolled_back", nil)
	}
	
	return nil
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	// TODO: Implement more sophisticated error classification
	// For now, consider network-related errors as retryable
	return err != nil
}

// GetProgress returns the current progress of the pipeline
func (p *InstallationPipeline) GetProgress() string {
	return fmt.Sprintf("Steps: %d, State: %s", len(p.Steps), p.State.String())
}

// GetState returns the current installation state
func (p *InstallationPipeline) GetState() *InstallationState {
	return p.State
} 