package testutil

import (
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// MockLogger implements interfaces.Logger for testing
type MockLogger struct {
	DebugMessages []string
	InfoMessages  []string
	WarnMessages  []string
	ErrorMessages []string
	Commands      []string
}

// NewMockLogger creates a new mock logger
func NewMockLogger() interfaces.Logger {
	return &MockLogger{
		DebugMessages: make([]string, 0),
		InfoMessages:  make([]string, 0),
		WarnMessages:  make([]string, 0),
		ErrorMessages: make([]string, 0),
		Commands:      make([]string, 0),
	}
}

// Debug logs a debug message
func (l *MockLogger) Debug(format string, args ...interface{}) {
	l.DebugMessages = append(l.DebugMessages, format)
}

// Info logs an info message
func (l *MockLogger) Info(format string, args ...interface{}) {
	l.InfoMessages = append(l.InfoMessages, format)
}

// Warn logs a warning message
func (l *MockLogger) Warn(format string, args ...interface{}) {
	l.WarnMessages = append(l.WarnMessages, format)
}

// Error logs an error message
func (l *MockLogger) Error(format string, args ...interface{}) {
	l.ErrorMessages = append(l.ErrorMessages, format)
}

// CommandStart logs the start of a command execution
func (l *MockLogger) CommandStart(cmd string, attempt, maxAttempts int) {
	l.Commands = append(l.Commands, cmd)
}

// CommandSuccess logs successful command execution
func (l *MockLogger) CommandSuccess(cmd string, duration time.Duration) {
	l.Commands = append(l.Commands, cmd)
}

// CommandError logs command execution error
func (l *MockLogger) CommandError(cmd string, err error, attempt, maxAttempts int) {
	l.Commands = append(l.Commands, cmd)
} 