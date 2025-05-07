package pipeline

import (
	"errors"
	"testing"
	"time"
)

func TestInstallationState(t *testing.T) {
	state := NewInstallationState()

	// Test initial state
	if state.Status != "initialized" {
		t.Error("Initial state should be 'initialized'")
	}

	// Test state updates
	state.UpdateState("test", "running", nil)
	if state.CurrentStep != "test" || state.Status != "running" {
		t.Error("State should be updated correctly")
	}

	// Test error state
	testErr := errors.New("test error")
	state.UpdateState("test", "failed", testErr)
	if state.Error != testErr {
		t.Error("Error should be set correctly")
	}
}

func TestInstallationPipeline(t *testing.T) {
	// Create a dummy channel for testing
	dummyChan := make(chan ProgressEvent, 10) // Buffered to avoid blocking if not read
	// Close the channel when the test finishes to prevent leaks if pipeline runs in goroutine (though not here)
	defer close(dummyChan)
	
	pipeline := NewInstallationPipeline(dummyChan)

	// Test successful step
	successStep := InstallationStep{
		Name: "success",
		Action: func() error {
			return nil
		},
	}
	pipeline.AddStep(successStep)

	// Test failing step
	failStep := InstallationStep{
		Name: "fail",
		Action: func() error {
			return errors.New("step failed")
		},
		Rollback: func() error {
			return nil
		},
	}
	pipeline.AddStep(failStep)

	// Execute pipeline
	err := pipeline.Execute()
	if err == nil {
		t.Error("Pipeline should fail due to failing step")
	}
}

func TestPlatformDetection(t *testing.T) {
	platform, err := DetectPlatform()
	if err != nil {
		t.Errorf("Failed to detect platform: %v", err)
	}

	if platform.OS == "" {
		t.Error("OS should not be empty")
	}

	if platform.Arch == "" {
		t.Error("Arch should not be empty")
	}

	if platform.PackageManager == "" {
		t.Error("PackageManager should not be empty")
	}

	if platform.Shell == "" {
		t.Error("Shell should not be empty")
	}
}

func TestPipelineTimeout(t *testing.T) {
	dummyChan := make(chan ProgressEvent, 10)
	defer close(dummyChan)
	pipeline := NewInstallationPipeline(dummyChan)

	// Test step with timeout
	timeoutStep := InstallationStep{
		Name: "timeout",
		Action: func() error {
			time.Sleep(2 * time.Second)
			return nil
		},
		Timeout: 1 * time.Second,
	}
	pipeline.AddStep(timeoutStep)

	err := pipeline.Execute()
	if err == nil {
		t.Error("Step should timeout")
	}
}

func TestPipelineRetry(t *testing.T) {
	dummyChan := make(chan ProgressEvent, 10)
	defer close(dummyChan)
	pipeline := NewInstallationPipeline(dummyChan)

	attempts := 0
	retryStep := InstallationStep{
		Name: "retry",
		Action: func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary failure")
			}
			return nil
		},
		RetryCount: 3,
		RetryDelay: 100 * time.Millisecond,
	}
	pipeline.AddStep(retryStep)

	err := pipeline.Execute()
	if err != nil {
		t.Errorf("Step should succeed after retries: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
} 