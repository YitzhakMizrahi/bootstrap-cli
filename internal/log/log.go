package log

import (
	"io"
	"log"
	"os"
)

// LogLevel defines the level of logging
type LogLevel int

const (
	// DebugLevel logs everything
	DebugLevel LogLevel = iota
	// InfoLevel logs informational messages
	InfoLevel
	// WarnLevel logs warnings
	WarnLevel
	// ErrorLevel logs errors
	ErrorLevel
	// FatalLevel logs errors and exits
	FatalLevel
)

// Logger is a wrapper around the standard log package
type Logger struct {
	logger *log.Logger
	level  LogLevel
}

// New creates a new logger with the specified level
func New(level LogLevel) *Logger {
	return &Logger{
		logger: log.New(os.Stderr, "", log.Ldate|log.Ltime),
		level:  level,
	}
}

// SetOutput sets the output destination for the logger
func (l *Logger) SetOutput(w io.Writer) {
	l.logger.SetOutput(w)
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DebugLevel {
		l.logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.logger.Printf("[INFO]  "+format, v...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WarnLevel {
		l.logger.Printf("[WARN]  "+format, v...)
	}
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ErrorLevel {
		l.logger.Printf("[ERROR] "+format, v...)
	}
}

// Success logs a success message (convenience function, treated as Info)
func (l *Logger) Success(format string, v ...interface{}) {
	if l.level <= InfoLevel {
		l.logger.Printf("[SUCCESS] "+format, v...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.level <= FatalLevel {
		l.logger.Fatalf("[FATAL] "+format, v...)
	}
}

// Printf logs a formatted message regardless of level (useful for direct output)
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
} 