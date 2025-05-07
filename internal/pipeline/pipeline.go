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
	Steps        []InstallationStep
	State        *InstallationState
	Logger       *log.Logger
	progressChan chan<- ProgressEvent // Channel to send progress updates (write-only)
}

// NewInstallationPipeline creates a new installation pipeline
// It now accepts a channel for sending progress events.
func NewInstallationPipeline(progressChan chan<- ProgressEvent) *InstallationPipeline {
	return &InstallationPipeline{
		Steps:        make([]InstallationStep, 0),
		State:        NewInstallationState(),
		progressChan: progressChan,
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
	var finalError error
	// startTime := time.Now() // Track start time for duration - Removed as not used for overall pipeline duration event yet

	// Ensure channel is closed when execution finishes (success or failure)
	if p.progressChan != nil {
		defer close(p.progressChan)
	}

	for i, step := range p.Steps {
		stepStartTime := time.Now()
		p.State.UpdateState(step.Name, "running", nil)
		p.sendProgress(TaskStart{TaskID: step.Name, Description: step.Description})
		
		// Execute step with retry
		err := p.executeStepWithRetry(step)
		duration := time.Since(stepStartTime)

		if err != nil {
			p.State.UpdateState(step.Name, "failed", err)
			p.sendProgress(TaskEnd{TaskID: step.Name, Success: false, Error: err, Duration: duration})
			
			// Attempt rollback of completed steps
			rollbackErr := p.rollback(i)
			if rollbackErr != nil {
				finalError = fmt.Errorf("step '%s' failed: %w; rollback also failed: %w", 
					step.Name, err, rollbackErr)
			} else {
				finalError = fmt.Errorf("step '%s' failed: %w; rollback successful", step.Name, err)
			}
			// Send complete message immediately on critical failure + rollback attempt
			p.sendProgress(PipelineComplete{OverallSuccess: false, FinalError: finalError})
			return finalError // Stop pipeline execution
		}
		
		p.State.UpdateState(step.Name, "completed", nil)
		p.sendProgress(TaskEnd{TaskID: step.Name, Success: true, Duration: duration})
	}
	
	p.State.UpdateState("pipeline", "completed", nil)
	// TODO: Maybe add overall duration to PipelineComplete event if needed?
	p.sendProgress(PipelineComplete{OverallSuccess: true, FinalError: nil})
	return nil
}

// executeStepWithRetry executes a step with retry logic
func (p *InstallationPipeline) executeStepWithRetry(step InstallationStep) error {
	var lastErr error
	
	for attempt := 0; attempt <= step.RetryCount; attempt++ {
		if attempt > 0 {
			p.State.UpdateState(step.Name, "retrying", lastErr)
			// TODO: Send a TaskLog or specific Retry message?
			p.sendProgress(TaskLog{TaskID: step.Name, Line: fmt.Sprintf("Retrying (attempt %d/%d)... Error: %v", attempt, step.RetryCount, lastErr)})
			time.Sleep(step.RetryDelay)
		}
		
		// Execute step
		// TODO: Capture stdout/stderr from step.Action() and send as TaskLog events if possible.
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
	
	return fmt.Errorf("failed after %d attempts: %w", step.RetryCount, lastErr)
}

// rollback attempts to roll back completed steps in reverse order
func (p *InstallationPipeline) rollback(lastCompletedIndex int) error {
	var firstRollbackErr error
	p.sendProgress(TaskLog{TaskID: "pipeline", Line: "Attempting rollback..."})

	for i := lastCompletedIndex; i >= 0; i-- {
		step := p.Steps[i]
		stepStartTime := time.Now()
		p.State.UpdateState(step.Name, "rolling_back", nil)
		p.sendProgress(TaskStart{TaskID: step.Name + "-rollback", Description: "Rolling back: " + step.Name})
		
		// Execute rollback action if defined
		var rollbackErr error
		if step.Rollback != nil {
			rollbackErr = step.Rollback() // TODO: Add retry logic to rollback?
		}
		duration := time.Since(stepStartTime)

		if rollbackErr != nil {
			p.State.UpdateState(step.Name, "rollback_failed", rollbackErr)
			p.sendProgress(TaskEnd{TaskID: step.Name + "-rollback", Success: false, Error: rollbackErr, Duration: duration})
			if firstRollbackErr == nil {
				firstRollbackErr = fmt.Errorf("rollback failed for step '%s': %w", step.Name, rollbackErr)
			}
			// Continue trying to rollback other steps even if one fails
		} else {
		p.State.UpdateState(step.Name, "rolled_back", nil)
			p.sendProgress(TaskEnd{TaskID: step.Name + "-rollback", Success: true, Duration: duration})
		}
	}
	
	return firstRollbackErr // Return the first error encountered during rollback
}

// sendProgress sends an event to the progress channel if it's not nil.
func (p *InstallationPipeline) sendProgress(event ProgressEvent) {
	if p.progressChan != nil {
		// Use non-blocking send or buffered channel to prevent pipeline locking if UI isn't reading
		// For simplicity now, using blocking send. Consider buffered channel in New.
		p.progressChan <- event
	}
	if p.Logger != nil { // Also log the event string representation
		p.Logger.Printf("Progress Event: %s", event)
	}
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