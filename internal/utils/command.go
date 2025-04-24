package utils

import (
	"fmt"
	"os/exec"
	"time"
)

// CommandExecutor provides utilities for executing commands with retries and error handling
type CommandExecutor struct {
	// Default number of retries for commands
	DefaultRetries int
	// Default delay between retries
	DefaultDelay time.Duration
	// Logger for command execution
	Logger *InstallLogger
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(logger *InstallLogger) *CommandExecutor {
	return &CommandExecutor{
		DefaultRetries: 3,
		DefaultDelay:   time.Second * 2,
		Logger:         logger,
	}
}

// ExecuteWithRetry executes a command with retries
func (e *CommandExecutor) ExecuteWithRetry(cmd *exec.Cmd, retries int, delay time.Duration) error {
	if retries <= 0 {
		retries = e.DefaultRetries
	}
	if delay <= 0 {
		delay = e.DefaultDelay
	}

	var err error
	for i := 0; i < retries; i++ {
		start := time.Now()
		e.Logger.CommandStart(cmd.String(), i+1, retries)
		
		err = cmd.Run()
		duration := time.Since(start)
		
		if err == nil {
			e.Logger.CommandSuccess(cmd.String(), duration)
			return nil
		}
		
		e.Logger.CommandError(cmd.String(), err, i+1, retries)
		
		if i < retries-1 {
			e.Logger.Debug("Waiting %v before retry...", delay)
			time.Sleep(delay)
		}
	}
	
	return fmt.Errorf("failed after %d retries: %w", retries, err)
}

// ExecuteWithOutput executes a command and returns its output
func (e *CommandExecutor) ExecuteWithOutput(cmd *exec.Cmd, retries int, delay time.Duration) (string, error) {
	if retries <= 0 {
		retries = e.DefaultRetries
	}
	if delay <= 0 {
		delay = e.DefaultDelay
	}

	var output []byte
	var err error
	
	for i := 0; i < retries; i++ {
		start := time.Now()
		e.Logger.CommandStart(cmd.String(), i+1, retries)
		
		output, err = cmd.Output()
		duration := time.Since(start)
		
		if err == nil {
			e.Logger.CommandSuccess(cmd.String(), duration)
			return string(output), nil
		}
		
		e.Logger.CommandError(cmd.String(), err, i+1, retries)
		
		if i < retries-1 {
			e.Logger.Debug("Waiting %v before retry...", delay)
			time.Sleep(delay)
		}
	}
	
	return "", fmt.Errorf("failed after %d retries: %w", retries, err)
} 