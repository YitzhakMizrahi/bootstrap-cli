package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestInstallLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewInstallLogger(true)

	tests := []struct {
		name     string
		fn       func()
		contains []string
	}{
		{
			name: "Debug message with debug enabled",
			fn:   func() { logger.Debug("debug message") },
			contains: []string{
				"[DEBUG]",
				"debug message",
			},
		},
		{
			name: "Info message",
			fn:   func() { logger.Info("info message") },
			contains: []string{
				"[INFO]",
				"info message",
			},
		},
		{
			name: "Warn message",
			fn:   func() { logger.Warn("warn message") },
			contains: []string{
				"[WARN]",
				"warn message",
			},
		},
		{
			name: "Error message",
			fn:   func() { logger.Error("error message") },
			contains: []string{
				"[ERROR]",
				"error message",
			},
		},
		{
			name: "Command start",
			fn:   func() { logger.CommandStart("test command", 1, 3) },
			contains: []string{
				"Executing command (attempt 1/3):",
				"test command",
			},
		},
		{
			name: "Command success",
			fn:   func() { logger.CommandSuccess("test command", time.Second) },
			contains: []string{
				"Command completed successfully in",
				"test command",
			},
		},
		{
			name: "Command error",
			fn:   func() { logger.CommandError("test command", nil, 1, 3) },
			contains: []string{
				"Command failed (attempt 1/3):",
				"test command",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.fn()
			output := buf.String()
			for _, s := range tt.contains {
				if !strings.Contains(output, s) {
					t.Errorf("output %q does not contain %q", output, s)
				}
			}
		})
	}
}

func TestInstallLoggerDebugDisabled(t *testing.T) {
	logger := NewInstallLogger(false)
	if logger.DebugEnabled {
		t.Error("Debug should be disabled")
	}
} 