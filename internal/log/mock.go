// Package log provides mock implementations for testing
package log

import (
	"time"
)

// MockLogger is a mock implementation of the Logger interface for testing
type MockLogger struct {
	DebugMessages     []string
	InfoMessages      []string
	WarnMessages      []string
	ErrorMessages     []string
	CommandStarts     []string
	CommandSuccesses  []string
	CommandErrors     []string
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		DebugMessages:    make([]string, 0),
		InfoMessages:     make([]string, 0),
		WarnMessages:     make([]string, 0),
		ErrorMessages:    make([]string, 0),
		CommandStarts:    make([]string, 0),
		CommandSuccesses: make([]string, 0),
		CommandErrors:    make([]string, 0),
	}
}

// Debug logs a debug message
func (m *MockLogger) Debug(format string, _ ...interface{}) {
	m.DebugMessages = append(m.DebugMessages, format)
}

// Info logs an info message
func (m *MockLogger) Info(format string, _ ...interface{}) {
	m.InfoMessages = append(m.InfoMessages, format)
}

// Warn logs a warning message
func (m *MockLogger) Warn(format string, _ ...interface{}) {
	m.WarnMessages = append(m.WarnMessages, format)
}

// Error logs an error message
func (m *MockLogger) Error(format string, _ ...interface{}) {
	m.ErrorMessages = append(m.ErrorMessages, format)
}

// CommandStart logs the start of a command execution
func (m *MockLogger) CommandStart(cmd string, _ int, _ int) {
	m.CommandStarts = append(m.CommandStarts, cmd)
}

// CommandSuccess logs the successful completion of a command
func (m *MockLogger) CommandSuccess(cmd string, _ time.Duration) {
	m.CommandSuccesses = append(m.CommandSuccesses, cmd)
}

// CommandError logs a command execution error
func (m *MockLogger) CommandError(cmd string, _ error, _ int, _ int) {
	m.CommandErrors = append(m.CommandErrors, cmd)
}

// SetLevel sets the log level
func (m *MockLogger) SetLevel(_ Level) {} 