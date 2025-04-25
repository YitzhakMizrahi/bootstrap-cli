package log

import (
	"testing"
	"time"
)

// mockLogger is a mock implementation of the Logger interface for testing
type mockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warnMessages    []string
	errorMessages   []string
	commandStarts   []string
	commandSuccesses []string
	commandErrors   []string
}

// newMockLogger creates a new mock logger
func newMockLogger() *mockLogger {
	return &mockLogger{
		debugMessages:   make([]string, 0),
		infoMessages:    make([]string, 0),
		warnMessages:    make([]string, 0),
		errorMessages:   make([]string, 0),
		commandStarts:   make([]string, 0),
		commandSuccesses: make([]string, 0),
		commandErrors:   make([]string, 0),
	}
}

// Debug logs a debug message
func (m *mockLogger) Debug(format string, _ ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

// Info logs an info message
func (m *mockLogger) Info(format string, _ ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

// Warn logs a warning message
func (m *mockLogger) Warn(format string, _ ...interface{}) {
	m.warnMessages = append(m.warnMessages, format)
}

// Error logs an error message
func (m *mockLogger) Error(format string, _ ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}

// CommandStart logs the start of a command execution
func (m *mockLogger) CommandStart(cmd string, _ int, _ int) {
	m.commandStarts = append(m.commandStarts, cmd)
}

// CommandSuccess logs the successful completion of a command
func (m *mockLogger) CommandSuccess(cmd string, _ time.Duration) {
	m.commandSuccesses = append(m.commandSuccesses, cmd)
}

// CommandError logs a command execution error
func (m *mockLogger) CommandError(cmd string, _ error, _ int, _ int) {
	m.commandErrors = append(m.commandErrors, cmd)
}

// SetLevel sets the log level
func (m *mockLogger) SetLevel(_ Level) {}

func TestLogAdapter(t *testing.T) {
	mock := newMockLogger()
	adapter := NewAdapter(mock)

	tests := []struct {
		name     string
		fn       func()
		check    func(*testing.T, *mockLogger)
	}{
		{
			name: "Debug message",
			fn:   func() { adapter.Debug("debug message") },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.debugMessages) != 1 || m.debugMessages[0] != "debug message" {
					t.Errorf("Debug message not properly forwarded")
				}
			},
		},
		{
			name: "Info message",
			fn:   func() { adapter.Info("info message") },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.infoMessages) != 1 || m.infoMessages[0] != "info message" {
					t.Errorf("Info message not properly forwarded")
				}
			},
		},
		{
			name: "Warn message",
			fn:   func() { adapter.Warn("warn message") },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.warnMessages) != 1 || m.warnMessages[0] != "warn message" {
					t.Errorf("Warn message not properly forwarded")
				}
			},
		},
		{
			name: "Error message",
			fn:   func() { adapter.Error("error message") },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.errorMessages) != 1 || m.errorMessages[0] != "error message" {
					t.Errorf("Error message not properly forwarded")
				}
			},
		},
		{
			name: "Command start",
			fn:   func() { adapter.CommandStart("test command", 1, 3) },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.commandStarts) != 1 || m.commandStarts[0] != "test command" {
					t.Errorf("Command start not properly forwarded")
				}
			},
		},
		{
			name: "Command success",
			fn:   func() { adapter.CommandSuccess("test command", time.Second) },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.commandSuccesses) != 1 || m.commandSuccesses[0] != "test command" {
					t.Errorf("Command success not properly forwarded")
				}
			},
		},
		{
			name: "Command error",
			fn:   func() { adapter.CommandError("test command", nil, 1, 3) },
			check: func(t *testing.T, m *mockLogger) {
				if len(m.commandErrors) != 1 || m.commandErrors[0] != "test command" {
					t.Errorf("Command error not properly forwarded")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			tt.check(t, mock)
		})
	}
} 