package pipeline

import (
	"fmt"
	"time"
)

// ProgressEvent is an interface that all progress messages should satisfy (optional but good practice).
type ProgressEvent interface {
	IsProgressEvent()
}

// TaskStart indicates a specific installation step has begun.
type TaskStart struct {
	TaskID      string // Unique identifier for the task/step (e.g., step.Name)
	Description string // User-friendly description (e.g., "Installing git...")
}
func (TaskStart) IsProgressEvent() {}

// TaskProgress indicates progress within a potentially long-running task.
type TaskProgress struct {
	TaskID  string  // Unique identifier for the task/step
	Percent float64 // Progress percentage (0.0 to 100.0), -1 if indeterminate
	Message string  // Optional message (e.g., "Downloading file MB/Total MB")
}
func (TaskProgress) IsProgressEvent() {}

// TaskLog provides a log line related to a specific task.
type TaskLog struct {
	TaskID string // Unique identifier for the task/step
	Line   string // The log line content
}
func (TaskLog) IsProgressEvent() {}

// TaskEnd indicates a specific installation step has finished.
type TaskEnd struct {
	TaskID   string        // Unique identifier for the task/step
	Success  bool          // Whether the step succeeded
	Error    error         // Error message if Success is false
	Duration time.Duration // How long the step took
}
func (TaskEnd) IsProgressEvent() {}

// PipelineComplete indicates the entire installation sequence has finished.
type PipelineComplete struct {
	OverallSuccess bool  // Whether all steps succeeded (or rollback completed)
	FinalError     error // Any critical error that stopped the pipeline or occurred during rollback
}
func (PipelineComplete) IsProgressEvent() {}

// Helper function to format error for messages (avoids nil pointer issues)
func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// String methods for nice printing (optional, mainly for debugging)
func (e TaskStart) String() string {
	return fmt.Sprintf("START [%s]: %s", e.TaskID, e.Description)
}
func (e TaskProgress) String() string {
	if e.Percent >= 0 {
		return fmt.Sprintf("PROG  [%s]: %.1f%% %s", e.TaskID, e.Percent, e.Message)
	}
	return fmt.Sprintf("PROG  [%s]: %s", e.TaskID, e.Message)
}
func (e TaskLog) String() string {
	return fmt.Sprintf("LOG   [%s]: %s", e.TaskID, e.Line)
}
func (e TaskEnd) String() string {
	if e.Success {
		return fmt.Sprintf("END   [%s]: OK (%.2fs)", e.TaskID, e.Duration.Seconds())
	}
	return fmt.Sprintf("END   [%s]: FAILED (%.2fs) - %s", e.TaskID, e.Duration.Seconds(), errorString(e.Error))
}
func (e PipelineComplete) String() string {
	if e.OverallSuccess {
		return "PIPELINE COMPLETE: SUCCESS"
	}
	return fmt.Sprintf("PIPELINE COMPLETE: FAILED - %s", errorString(e.FinalError))
} 