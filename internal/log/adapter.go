package log

import (
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// LogAdapter adapts a Logger to implement the InstallLogger interface
type LogAdapter struct {
	logger interfaces.Logger
}

// NewLogAdapter creates a new logger adapter
func NewLogAdapter(logger interfaces.Logger) *LogAdapter {
	return &LogAdapter{
		logger: logger,
	}
}

// CommandStart logs the start of a command execution
func (l *LogAdapter) CommandStart(cmd string, attempt, maxAttempts int) {
	l.logger.CommandStart(cmd, attempt, maxAttempts)
}

// CommandSuccess logs the successful completion of a command
func (l *LogAdapter) CommandSuccess(cmd string, duration time.Duration) {
	l.logger.CommandSuccess(cmd, duration)
}

// CommandError logs a command execution error
func (l *LogAdapter) CommandError(cmd string, err error, attempt, maxAttempts int) {
	l.logger.CommandError(cmd, err, attempt, maxAttempts)
}

// Debug logs a debug message
func (l *LogAdapter) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

// Info logs an info message
func (l *LogAdapter) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

// Warn logs a warning message
func (l *LogAdapter) Warn(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

// Error logs an error message
func (l *LogAdapter) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
} 