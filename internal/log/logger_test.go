package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
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

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DebugLevel, WithOutput(&buf))

	tests := []struct {
		level   LogLevel
		logFunc func(string, ...interface{})
		prefix  string
	}{
		{DebugLevel, logger.Debug, "DEBUG"},
		{InfoLevel, logger.Info, "INFO"},
		{WarnLevel, logger.Warn, "WARN"},
		{ErrorLevel, logger.Error, "ERROR"},
	}

	for _, tt := range tests {
		buf.Reset()
		tt.logFunc("test message")
		output := buf.String()
		if !strings.Contains(output, "["+tt.prefix+"]") {
			t.Errorf("Expected log level %s, got %s", tt.prefix, output)
		}
	}
}

func TestLogger_Metadata(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, WithOutput(&buf))
	logger.AddMetadata("test", "value")
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "(test=value)") {
		t.Errorf("Expected metadata in output, got %s", output)
	}
}

func TestLogger_Component(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, WithOutput(&buf), WithComponent("test-component"))
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "[test-component]") {
		t.Errorf("Expected component in output, got %s", output)
	}
}

func TestLogger_InstallationSpecific(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DebugLevel, WithOutput(&buf))

	// Test CommandStart
	buf.Reset()
	logger.CommandStart("test-cmd", 1, 3)
	if !strings.Contains(buf.String(), "Starting command execution (attempt 1/3)") {
		t.Errorf("CommandStart output incorrect: %s", buf.String())
	}

	// Test CommandSuccess
	buf.Reset()
	logger.CommandSuccess("test-cmd", time.Second)
	if !strings.Contains(buf.String(), "Command completed successfully") {
		t.Errorf("CommandSuccess output incorrect: %s", buf.String())
	}

	// Test CommandError
	buf.Reset()
	logger.CommandError("test-cmd", nil, 1, 3)
	if !strings.Contains(buf.String(), "Command failed (attempt 1/3)") {
		t.Errorf("CommandError output incorrect: %s", buf.String())
	}

	// Test StepStart
	buf.Reset()
	logger.StepStart("test-step")
	if !strings.Contains(buf.String(), "Starting step: test-step") {
		t.Errorf("StepStart output incorrect: %s", buf.String())
	}

	// Test StepSuccess
	buf.Reset()
	logger.StepSuccess("test-step", time.Second)
	if !strings.Contains(buf.String(), "Step completed successfully") {
		t.Errorf("StepSuccess output incorrect: %s", buf.String())
	}

	// Test StepError
	buf.Reset()
	logger.StepError("test-step", nil)
	if !strings.Contains(buf.String(), "Step failed: test-step") {
		t.Errorf("StepError output incorrect: %s", buf.String())
	}

	// Test SystemInfo
	buf.Reset()
	logger.SystemInfo("linux", "apt", map[string]string{"PATH": "/usr/bin"})
	output := buf.String()
	if !strings.Contains(output, "Platform: linux") ||
		!strings.Contains(output, "Package Manager: apt") ||
		!strings.Contains(output, "PATH=/usr/bin") {
		t.Errorf("SystemInfo output incorrect: %s", output)
	}

	// Test DependencyInfo
	buf.Reset()
	logger.DependencyInfo([]string{"dep1", "dep2"})
	if !strings.Contains(buf.String(), "Dependencies: dep1, dep2") {
		t.Errorf("DependencyInfo output incorrect: %s", buf.String())
	}

	// Test VerificationInfo
	buf.Reset()
	logger.VerificationInfo("test-bin", "1.0.0", []string{"/usr/bin"})
	output = buf.String()
	if !strings.Contains(output, "Binary: test-bin") ||
		!strings.Contains(output, "Version: 1.0.0") ||
		!strings.Contains(output, "Paths: /usr/bin") {
		t.Errorf("VerificationInfo output incorrect: %s", output)
	}
}

func TestLogger_Options(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel,
		WithOutput(&buf),
		WithComponent("test"),
		WithMetadata("key", "value"),
	)

	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "[test]") {
		t.Errorf("Expected component in output, got %s", output)
	}
	if !strings.Contains(output, "(key=value)") {
		t.Errorf("Expected metadata in output, got %s", output)
	}
} 