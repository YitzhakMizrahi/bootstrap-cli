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
	Logger Logger
}

// Logger interface for command execution logging
type Logger interface {
	Printf(format string, args ...interface{})
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(logger Logger) *CommandExecutor {
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
		e.Logger.Printf("Executing command: %s (attempt %d/%d)", cmd.String(), i+1, retries)
		
		err = cmd.Run()
		if err == nil {
			e.Logger.Printf("Command executed successfully: %s", cmd.String())
			return nil
		}
		
		e.Logger.Printf("Command failed: %s (error: %v)", cmd.String(), err)
		
		if i < retries-1 {
			e.Logger.Printf("Retrying in %v...", delay)
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
		e.Logger.Printf("Executing command: %s (attempt %d/%d)", cmd.String(), i+1, retries)
		
		output, err = cmd.Output()
		if err == nil {
			e.Logger.Printf("Command executed successfully: %s", cmd.String())
			return string(output), nil
		}
		
		e.Logger.Printf("Command failed: %s (error: %v)", cmd.String(), err)
		
		if i < retries-1 {
			e.Logger.Printf("Retrying in %v...", delay)
			time.Sleep(delay)
		}
	}
	
	return "", fmt.Errorf("failed after %d retries: %w", retries, err)
} 