package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		logFunc  func(*Logger, string, ...interface{})
		logLevel string
		message  string
		want     bool
	}{
		{
			name:     "debug message at debug level",
			level:    DebugLevel,
			logFunc:  (*Logger).Debug,
			logLevel: "DEBUG",
			message:  "test debug message",
			want:     true,
		},
		{
			name:     "debug message at info level",
			level:    InfoLevel,
			logFunc:  (*Logger).Debug,
			logLevel: "DEBUG",
			message:  "test debug message",
			want:     false,
		},
		{
			name:     "info message at info level",
			level:    InfoLevel,
			logFunc:  (*Logger).Info,
			logLevel: "INFO",
			message:  "test info message",
			want:     true,
		},
		{
			name:     "warn message at warn level",
			level:    WarnLevel,
			logFunc:  (*Logger).Warn,
			logLevel: "WARN",
			message:  "test warn message",
			want:     true,
		},
		{
			name:     "error message at error level",
			level:    ErrorLevel,
			logFunc:  (*Logger).Error,
			logLevel: "ERROR",
			message:  "test error message",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(tt.level)
			logger.SetOutput(&buf)

			tt.logFunc(logger, tt.message)

			got := buf.String()
			if tt.want {
				if !strings.Contains(got, tt.logLevel) || !strings.Contains(got, tt.message) {
					t.Errorf("Logger output = %q, want containing level %q and message %q", got, tt.logLevel, tt.message)
				}
			} else {
				if got != "" {
					t.Errorf("Logger output = %q, want empty string", got)
				}
			}
		})
	}
} 