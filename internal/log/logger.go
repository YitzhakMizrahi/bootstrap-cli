package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	// DebugLevel logs everything
	DebugLevel Level = iota
	// InfoLevel logs info and above
	InfoLevel
	// WarnLevel logs warnings and above
	WarnLevel
	// ErrorLevel logs only errors
	ErrorLevel
)

// Logger handles logging operations
type Logger struct {
	level  Level
	output io.Writer
}

// New creates a new logger with the specified level
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		output: os.Stdout,
	}
}

// SetOutput sets the output writer for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.output = w
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log("DEBUG", format, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log("INFO", format, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log("WARN", format, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log("ERROR", format, args...)
	}
}

// log formats and writes a log message
func (l *Logger) log(level, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] %-5s %s\n", timestamp, level, msg)
} 