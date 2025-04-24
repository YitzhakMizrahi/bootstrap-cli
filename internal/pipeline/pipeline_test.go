package pipeline

import (
	"errors"
	"testing"
	"time"
)

func TestInstallationState(t *testing.T) {
	state := NewInstallationState()

	// Test initial state
	if state.IsInstalled("test") {
		t.Error("Package should not be installed initially")
	}

	// Test marking as pending
	state.MarkPending("test")
	if !state.IsPending("test") {
		t.Error("Package should be pending")
	}

	// Test marking as installed
	state.MarkInstalled("test")
	if !state.IsInstalled("test") {
		t.Error("Package should be installed")
	}
	if state.IsPending("test") {
		t.Error("Package should not be pending after installation")
	}

	// Test marking as failed
	state.MarkFailed("test2", errors.New("test error"))
	if !state.IsFailed("test2") {
		t.Error("Package should be failed")
	}
	if err := state.GetFailedError("test2"); err == nil {
		t.Error("Should return error for failed package")
	}
}

func TestInstallationPipeline(t *testing.T) {
	state := NewInstallationState()
	pipeline := NewInstallationPipeline(state)

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
	state := NewInstallationState()
	pipeline := NewInstallationPipeline(state)

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
	state := NewInstallationState()
	pipeline := NewInstallationPipeline(state)

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
		Timeout: 5 * time.Second,
	}

	err := pipeline.ExecuteWithRetry(retryStep, 3)
	if err != nil {
		t.Errorf("Step should succeed after retries: %v", err)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
} 