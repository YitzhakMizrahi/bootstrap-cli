// Package log provides logging functionality for the bootstrap-cli,
// including adapters for different logging implementations and interfaces.
package log

import (
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// Adapter adapts a Logger to implement the InstallLogger interface
type Adapter struct {
	logger interfaces.Logger
}

// NewAdapter creates a new logger adapter
func NewAdapter(logger interfaces.Logger) *Adapter {
	return &Adapter{
		logger: logger,
	}
}

// CommandStart logs the start of a command execution
func (l *Adapter) CommandStart(cmd string, attempt, maxAttempts int) {
	l.logger.CommandStart(cmd, attempt, maxAttempts)
}

// CommandSuccess logs the successful completion of a command
func (l *Adapter) CommandSuccess(cmd string, duration time.Duration) {
	l.logger.CommandSuccess(cmd, duration)
}

// CommandError logs a command execution error
func (l *Adapter) CommandError(cmd string, err error, attempt, maxAttempts int) {
	l.logger.CommandError(cmd, err, attempt, maxAttempts)
}

// Debug logs a debug message
func (l *Adapter) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

// Info logs an info message
func (l *Adapter) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

// Warn logs a warning message
func (l *Adapter) Warn(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

// Error logs an error message
func (l *Adapter) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
} 