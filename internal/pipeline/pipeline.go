package pipeline

import (
	"fmt"
	"log"
	"time"
)

// InstallationStep represents a single step in the installation process
type InstallationStep struct {
	Name     string
	Action   func() error
	Rollback func() error
	Timeout  time.Duration
}

// InstallationPipeline manages the execution of installation steps
type InstallationPipeline struct {
	Steps   []InstallationStep
	State   *InstallationState
	Logger  *log.Logger
}

// NewInstallationPipeline creates a new installation pipeline
func NewInstallationPipeline(state *InstallationState) *InstallationPipeline {
	return &InstallationPipeline{
		Steps:  make([]InstallationStep, 0),
		State:  state,
		Logger: log.New(log.Writer(), "[Pipeline] ", log.LstdFlags),
	}
}

// AddStep adds a new step to the pipeline
func (p *InstallationPipeline) AddStep(step InstallationStep) {
	p.Steps = append(p.Steps, step)
}

// Execute runs all steps in the pipeline
func (p *InstallationPipeline) Execute() error {
	for i, step := range p.Steps {
		p.Logger.Printf("Executing step %d/%d: %s", i+1, len(p.Steps), step.Name)
		
		// Set default timeout if not specified
		timeout := step.Timeout
		if timeout == 0 {
			timeout = 5 * time.Minute
		}

		// Execute step with timeout
		err := p.executeStepWithTimeout(step, timeout)
		if err != nil {
			p.Logger.Printf("Step failed: %v", err)
			
			// Attempt rollback if available
			if step.Rollback != nil {
				p.Logger.Printf("Attempting rollback for step: %s", step.Name)
				if rollbackErr := step.Rollback(); rollbackErr != nil {
					p.Logger.Printf("Rollback failed: %v", rollbackErr)
				}
			}
			
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
		
		p.Logger.Printf("Step completed successfully: %s", step.Name)
	}
	
	return nil
}

// executeStepWithTimeout runs a step with a timeout
func (p *InstallationPipeline) executeStepWithTimeout(step InstallationStep, timeout time.Duration) error {
	done := make(chan error, 1)
	
	go func() {
		done <- step.Action()
	}()
	
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("step timed out after %v", timeout)
	}
}

// ExecuteWithRetry runs a step with retries
func (p *InstallationPipeline) ExecuteWithRetry(step InstallationStep, maxRetries int) error {
	var lastErr error
	
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			p.Logger.Printf("Retry %d/%d for step: %s", i+1, maxRetries, step.Name)
			time.Sleep(time.Second * time.Duration(i))
		}
		
		err := p.executeStepWithTimeout(step, step.Timeout)
		if err == nil {
			return nil
		}
		
		lastErr = err
		p.Logger.Printf("Attempt %d failed: %v", i+1, err)
	}
	
	return fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// GetProgress returns the current progress of the pipeline
func (p *InstallationPipeline) GetProgress() string {
	return fmt.Sprintf("Steps: %d, State: %s", len(p.Steps), p.State.String())
} 