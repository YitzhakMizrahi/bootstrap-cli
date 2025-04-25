package utils

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// CommandExecutor provides utilities for executing commands with retries and error handling
type CommandExecutor struct {
	// Default number of retries for commands
	DefaultRetries int
	// Default delay between retries
	DefaultDelay time.Duration
	// Logger for command execution
	Logger *log.Logger
}

// NewCommandExecutor creates a new command executor
func NewCommandExecutor(logger *log.Logger) *CommandExecutor {
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
		e.Logger.Info("Running command (attempt %d/%d): %s", i+1, retries, cmd.String())
		
		err = cmd.Run()
		duration := time.Since(start)
		
		if err == nil {
			e.Logger.Success("Command completed in %v: %s", duration, cmd.String())
			return nil
		}
		
		e.Logger.Error("Command failed (attempt %d/%d): %s - %v", i+1, retries, cmd.String(), err)
		
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
		e.Logger.Info("Running command (attempt %d/%d): %s", i+1, retries, cmd.String())
		
		output, err = cmd.Output()
		duration := time.Since(start)
		
		if err == nil {
			e.Logger.Success("Command completed in %v: %s", duration, cmd.String())
			return string(output), nil
		}
		
		e.Logger.Error("Command failed (attempt %d/%d): %s - %v", i+1, retries, cmd.String(), err)
		
		if i < retries-1 {
			e.Logger.Debug("Waiting %v before retry...", delay)
			time.Sleep(delay)
		}
	}
	
	return "", fmt.Errorf("failed after %d retries: %w", retries, err)
} 