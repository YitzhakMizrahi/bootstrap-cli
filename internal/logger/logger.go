package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Level represents the logging level
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

// Logger handles logging operations
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new logger
func New(level Level, logPath string) (*Logger, error) {
	var output io.Writer = os.Stdout

	if logPath != "" {
		dir := filepath.Dir(logPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("error creating log directory: %w", err)
		}

		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening log file: %w", err)
		}
		output = file
	}

	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}, nil
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUG {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFO {
		l.logger.Printf("[INFO] "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARN {
		l.logger.Printf("[WARN] "+format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERROR {
		l.logger.Printf("[ERROR] "+format, v...)
	}
} 