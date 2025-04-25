package interfaces

import "time"

// Logger defines the interface for logging operations
type Logger interface {
	// Debug logs a debug message
	Debug(format string, args ...interface{})
	// Info logs an info message
	Info(format string, args ...interface{})
	// Warn logs a warning message
	Warn(format string, args ...interface{})
	// Error logs an error message
	Error(format string, args ...interface{})
	// CommandStart logs the start of a command execution
	CommandStart(cmd string, attempt, maxAttempts int)
	// CommandSuccess logs the successful completion of a command
	CommandSuccess(cmd string, duration time.Duration)
	// CommandError logs a command execution error
	CommandError(cmd string, err error, attempt, maxAttempts int)
} 