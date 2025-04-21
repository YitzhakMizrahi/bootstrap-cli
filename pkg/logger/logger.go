package logger

import (
	"io"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Color returns the ANSI color code for the log level
func (l Level) Color() string {
	switch l {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case WARN:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	default:
		return "\033[0m" // Reset
	}
}

// Options represents the configuration options for the logger
type Options struct {
	Level      Level
	Output     io.Writer
	File       *os.File
	MaxSize    int64
	MaxBackups int
	MaxAge     time.Duration
	Fields     map[string]interface{}
}

// Logger defines the interface for logging operations
type Logger interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
	WithFields(fields map[string]interface{}) Logger
	SetLevel(level Level)
	Close() error
} 