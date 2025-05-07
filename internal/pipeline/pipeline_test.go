package pipeline

import (
	"errors"
	"testing"
	"time"
	// Import for interfaces.PackageManager
)

// --- Fake PackageManager for testing ---
type fakePM struct {}
func (f *fakePM) Install(pkg string) error             { return nil }
func (f *fakePM) Uninstall(pkg string) error           { return nil }
func (f *fakePM) IsInstalled(pkg string) (bool, error) { return false, nil }
func (f *fakePM) Update() error                      { return nil }
func (f *fakePM) SetupSpecialPackage(pkg string) error { return nil }
func (f *fakePM) IsPackageAvailable(pkg string) bool { return true }
func (f *fakePM) GetName() string                    { return "fake" }
// Ensure it satisfies the pipeline.PackageManager interface (defined in interfaces.go of this package)
var _ PackageManager = (*fakePM)(nil)
// It ALSO needs to satisfy interfaces.PackageManager if used elsewhere expecting that.
// var _ interfaces.PackageManager = (*fakePM)(nil) // Add if needed and update methods

// --- Helper to create context --- 
func newTestContext(t *testing.T) (*InstallationContext, chan ProgressEvent) {
	progChan := make(chan ProgressEvent, 10) 
	platform := &Platform{OS: "linux", Arch: "amd64", PackageManager: "apt", Shell: "bash"}
	
	// Use the package-level fakePM
	var pm PackageManager = &fakePM{}
	
	context := NewInstallationContext(platform, pm, progChan)
	return context, progChan
}

// --- Tests --- 

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
	ctx, progChan := newTestContext(t)
	defer close(progChan)
	
	pipeline := NewInstallationPipeline(ctx) // Pass context

	// Test successful step
	successStep := InstallationStep{
		Name: "success",
		Action: func(ctx *InstallationContext) error {
			return nil
		},
	}
	pipeline.AddStep(successStep)

	// Test failing step
	failStep := InstallationStep{
		Name: "fail",
		Action: func(ctx *InstallationContext) error {
			return errors.New("step failed")
		},
		Rollback: func(ctx *InstallationContext) error {
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
	ctx, progChan := newTestContext(t)
	defer close(progChan)
	pipeline := NewInstallationPipeline(ctx) // Pass context

	// Test step with timeout
	timeoutStep := InstallationStep{
		Name: "timeout",
		Action: func(ctx *InstallationContext) error {
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
	ctx, progChan := newTestContext(t)
	defer close(progChan)
	pipeline := NewInstallationPipeline(ctx) // Pass context

	attempts := 0
	retryStep := InstallationStep{
		Name: "retry",
		Action: func(ctx *InstallationContext) error {
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