package utils

import (
	"fmt"
	"time"
)

// InstallLogger provides logging functionality for installation operations
type InstallLogger struct {
	// Whether to show debug messages
	DebugEnabled bool
}

// CommandStart logs the start of a command execution
func (l *InstallLogger) CommandStart(cmd string, attempt, maxAttempts int) {
	if maxAttempts > 1 {
		fmt.Printf("Executing command (attempt %d/%d): %s\n", attempt, maxAttempts, cmd)
	} else {
		fmt.Printf("Executing command: %s\n", cmd)
	}
}

// CommandSuccess logs the successful completion of a command
func (l *InstallLogger) CommandSuccess(cmd string, duration time.Duration) {
	fmt.Printf("Command completed successfully in %v: %s\n", duration, cmd)
}

// CommandError logs a command execution error
func (l *InstallLogger) CommandError(cmd string, err error, attempt, maxAttempts int) {
	fmt.Printf("Command failed (attempt %d/%d): %s\nError: %v\n", attempt, maxAttempts, cmd, err)
}

// Debug logs a debug message
func (l *InstallLogger) Debug(format string, args ...interface{}) {
	if l.DebugEnabled {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Info logs an info message
func (l *InstallLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

// Warn logs a warning message
func (l *InstallLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

// Error logs an error message
func (l *InstallLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

// DefaultLogger is a simple implementation of InstallLogger
type DefaultLogger struct {
	DebugEnabled bool
}

// NewDefaultLogger creates a new default logger
func NewDefaultLogger() *InstallLogger {
	return &InstallLogger{
		DebugEnabled: false,
	}
}

// CommandStart logs the start of a command execution
func (l *DefaultLogger) CommandStart(cmd string, attempt, maxAttempts int) {
	if maxAttempts > 1 {
		fmt.Printf("Executing command (attempt %d/%d): %s\n", attempt, maxAttempts, cmd)
	} else {
		fmt.Printf("Executing command: %s\n", cmd)
	}
}

// CommandSuccess logs the successful completion of a command
func (l *DefaultLogger) CommandSuccess(cmd string, duration time.Duration) {
	fmt.Printf("Command completed successfully in %v: %s\n", duration, cmd)
}

// CommandError logs a command execution error
func (l *DefaultLogger) CommandError(cmd string, err error, attempt, maxAttempts int) {
	fmt.Printf("Command failed (attempt %d/%d): %s\nError: %v\n", attempt, maxAttempts, cmd, err)
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(format string, args ...interface{}) {
	if l.DebugEnabled {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Info logs an info message
func (l *DefaultLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(format string, args ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", args...)
}

// Error logs an error message
func (l *DefaultLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
} 