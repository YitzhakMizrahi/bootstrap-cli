package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Level represents a log level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger handles application logging
type Logger struct {
	level     Level
	output    io.Writer
	file      *os.File
	mu        sync.Mutex
}

// New creates a new logger with the specified level and output
func New(level Level, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// SetLevel sets the log level
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput sets the log output
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// SetFile sets the log file
func (l *Logger) SetFile(path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Close existing file if any
	if l.file != nil {
		l.file.Close()
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open file for appending
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	l.file = file
	l.output = io.MultiWriter(os.Stdout, file)
	return nil
}

// Close closes the log file
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log logs a message at the specified level
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	levelStr := "UNKNOWN"
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.output, "[%s] %s\n", levelStr, msg)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
} 