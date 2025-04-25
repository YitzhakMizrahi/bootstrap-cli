package log

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(InfoLevel, WithOutput(&buf))

	tests := []struct {
		name     string
		level    Level
		messages []string
	}{
		{
			name:  "debug level",
			level: DebugLevel,
			messages: []string{
				"debug message",
				"info message",
				"warn message",
				"error message",
			},
		},
		{
			name:  "info level",
			level: InfoLevel,
			messages: []string{
				"info message",
				"warn message",
				"error message",
			},
		},
		{
			name:  "warn level",
			level: WarnLevel,
			messages: []string{
				"warn message",
				"error message",
			},
		},
		{
			name:  "error level",
			level: ErrorLevel,
			messages: []string{
				"error message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			for _, msg := range tt.messages {
				switch tt.level {
				case DebugLevel:
					logger.Debug("%s", msg)
				case InfoLevel:
					logger.Info("%s", msg)
				case WarnLevel:
					logger.Warn("%s", msg)
				case ErrorLevel:
					logger.Error("%s", msg)
				}
			}
			output := buf.String()
			for _, msg := range tt.messages {
				if !strings.Contains(output, msg) {
					t.Errorf("output %q does not contain %q", output, msg)
				}
			}
		})
	}
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DebugLevel, WithOutput(&buf))

	tests := []struct {
		level   Level
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