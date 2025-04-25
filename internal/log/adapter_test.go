package log

import (
	"testing"
	"time"
)

type mockLogger struct {
	debugMessages     []string
	infoMessages      []string
	warnMessages      []string
	errorMessages     []string
	commandStarts     []string
	commandSuccesses  []string
	commandErrors     []string
}

func (m *mockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *mockLogger) Info(format string, args ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

func (m *mockLogger) Warn(format string, args ...interface{}) {
	m.warnMessages = append(m.warnMessages, format)
}

func (m *mockLogger) Error(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}

func (m *mockLogger) CommandStart(cmd string, attempt, maxAttempts int) {
	m.commandStarts = append(m.commandStarts, cmd)
}

func (m *mockLogger) CommandSuccess(cmd string, duration time.Duration) {
	m.commandSuccesses = append(m.commandSuccesses, cmd)
}

func (m *mockLogger) CommandError(cmd string, err error, attempt, maxAttempts int) {
	m.commandErrors = append(m.commandErrors, cmd)
}

func TestLogAdapter(t *testing.T) {
	mock := &mockLogger{}
	adapter := NewLogAdapter(mock)

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