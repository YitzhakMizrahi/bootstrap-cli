package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// Level defines the level of logging
type Level int

const (
	// DebugLevel logs everything
	DebugLevel Level = iota
	// InfoLevel logs informational messages
	InfoLevel
	// WarnLevel logs warnings
	WarnLevel
	// ErrorLevel logs errors
	ErrorLevel
	// FatalLevel logs errors and exits
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger is a wrapper around the standard log package with enhanced functionality
type Logger struct {
	logger    *log.Logger
	level     Level
	metadata  map[string]string
	component string
}

// LoggerOption is a function that configures a Logger
type LoggerOption func(*Logger)

// WithComponent sets the component name for the logger
func WithComponent(component string) LoggerOption {
	return func(l *Logger) {
		l.component = component
	}
}

// WithMetadata sets metadata key-value pairs for the logger
func WithMetadata(key, value string) LoggerOption {
	return func(l *Logger) {
		l.metadata[key] = value
	}
}

// WithOutput sets the output destination for the logger
func WithOutput(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.logger.SetOutput(w)
	}
}

// New creates a new logger with the specified level and options
func New(level Level, opts ...LoggerOption) *Logger {
	l := &Logger{
		logger:    log.New(os.Stderr, "", log.Ldate|log.Ltime),
		level:     level,
		metadata:  make(map[string]string),
		component: "main",
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// SetComponent sets the current component name
func (l *Logger) SetComponent(component string) {
	l.component = component
}

// AddMetadata adds a metadata key-value pair
func (l *Logger) AddMetadata(key, value string) {
	l.metadata[key] = value
}

// formatMessage formats a log message with metadata
func (l *Logger) formatMessage(level Level, format string, args ...interface{}) string {
	// Format the base message
	msg := fmt.Sprintf(format, args...)

	// Build metadata string
	var metadataPairs []string
	for k, v := range l.metadata {
		metadataPairs = append(metadataPairs, fmt.Sprintf("%s=%s", k, v))
	}
	metadata := ""
	if len(metadataPairs) > 0 {
		metadata = fmt.Sprintf(" (%s)", strings.Join(metadataPairs, ", "))
	}

	// Return formatted message with level, component, and metadata
	return fmt.Sprintf("[%s][%s]%s %s", level, l.component, metadata, msg)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		msg := l.formatMessage(DebugLevel, format, v...)
		l.logger.Print(msg)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		msg := l.formatMessage(InfoLevel, format, v...)
		l.logger.Print(msg)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		msg := l.formatMessage(WarnLevel, format, v...)
		l.logger.Print(msg)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		msg := l.formatMessage(ErrorLevel, format, v...)
		l.logger.Print(msg)
	}
}

// Success logs a success message (convenience function, treated as Info)
func (l *Logger) Success(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		msg := l.formatMessage(InfoLevel, "✓ "+format, v...)
		l.logger.Print(msg)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		msg := l.formatMessage(FatalLevel, format, v...)
		l.logger.Fatal(msg)
	}
}

// Printf logs a formatted message regardless of level (useful for direct output)
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Installation-specific logging methods

// CommandStart logs the start of a command execution
func (l *Logger) CommandStart(cmd string, attempt int, maxAttempts int) {
	l.Debug("Starting command execution (attempt %d/%d): %s", attempt, maxAttempts, cmd)
}

// CommandSuccess logs successful command execution
func (l *Logger) CommandSuccess(cmd string, duration time.Duration) {
	l.Debug("Command completed successfully in %v: %s", duration, cmd)
}

// CommandError logs command execution failure
func (l *Logger) CommandError(cmd string, err error, attempt int, maxAttempts int) {
	if attempt < maxAttempts {
		l.Warn("Command failed (attempt %d/%d): %s\nError: %v", attempt, maxAttempts, cmd, err)
	} else {
		l.Error("Command failed (final attempt %d/%d): %s\nError: %v", attempt, maxAttempts, cmd, err)
	}
}

// StepStart logs the start of an installation step
func (l *Logger) StepStart(step string) {
	prevComponent := l.component
	l.SetComponent(step)
	l.Info("Starting step: %s", step)
	l.component = prevComponent
}

// StepSuccess logs successful completion of an installation step
func (l *Logger) StepSuccess(step string, duration time.Duration) {
	l.Success("Step completed successfully in %v: %s", duration, step)
}

// StepError logs installation step failure
func (l *Logger) StepError(step string, err error) {
	l.Error("Step failed: %s\nError: %v", step, err)
}

// SystemInfo logs system information
func (l *Logger) SystemInfo(platform string, packageManager string, env map[string]string) {
	prevComponent := l.component
	l.SetComponent("system")
	l.Info("System Information:")
	l.Info("  Platform: %s", platform)
	l.Info("  Package Manager: %s", packageManager)
	l.Info("  Environment Variables:")
	for k, v := range env {
		l.Info("    %s=%s", k, v)
	}
	l.component = prevComponent
}

// DependencyInfo logs dependency information
func (l *Logger) DependencyInfo(dependencies []string) {
	if len(dependencies) > 0 {
		l.Info("Dependencies: %s", strings.Join(dependencies, ", "))
	} else {
		l.Info("No dependencies required")
	}
}

// VerificationInfo logs verification information
func (l *Logger) VerificationInfo(binary string, version string, paths []string) {
	prevComponent := l.component
	l.SetComponent("verify")
	l.Info("Verification Information:")
	l.Info("  Binary: %s", binary)
	l.Info("  Version: %s", version)
	l.Info("  Paths: %s", strings.Join(paths, ", "))
	l.component = prevComponent
}

// Ensure Logger implements the interfaces.Logger interface
var _ interfaces.Logger = (*Logger)(nil) 