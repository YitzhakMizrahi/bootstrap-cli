package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	resetColor = "\033[0m"
)

// loggerImpl implements the Logger interface
type loggerImpl struct {
	mu       sync.Mutex
	options  Options
	fields   map[string]interface{}
	rotator  *logRotator
}

// logRotator handles log file rotation
type logRotator struct {
	file       *os.File
	maxSize    int64
	maxBackups int
	maxAge     time.Duration
	mu         sync.Mutex
}

// New creates a new logger instance with the given options
func New(opts Options) (Logger, error) {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}

	var rotator *logRotator
	if opts.File != nil {
		rotator = &logRotator{
			file:       opts.File,
			maxSize:    opts.MaxSize,
			maxBackups: opts.MaxBackups,
			maxAge:     opts.MaxAge,
		}
		opts.Output = rotator
	}

	return &loggerImpl{
		options: opts,
		fields:  make(map[string]interface{}),
		rotator: rotator,
	}, nil
}

// log writes a log entry with the given level and message
func (l *loggerImpl) log(level Level, msg string, fields ...map[string]interface{}) {
	if level < l.options.Level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Merge all fields
	allFields := make(map[string]interface{})
	for k, v := range l.fields {
		allFields[k] = v
	}
	for _, f := range fields {
		for k, v := range f {
			allFields[k] = v
		}
	}

	// Create log entry
	entry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level.String(),
		"message":   msg,
		"fields":    allFields,
	}

	// Convert to JSON
	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	// Write to output with color if it's a terminal
	if f, ok := l.options.Output.(*os.File); ok && f == os.Stdout {
		fmt.Fprintf(l.options.Output, "%s%s%s\n", level.Color(), string(jsonEntry), resetColor)
	} else {
		fmt.Fprintf(l.options.Output, "%s\n", string(jsonEntry))
	}
}

func (l *loggerImpl) Debug(msg string, fields ...map[string]interface{}) {
	l.log(DEBUG, msg, fields...)
}

func (l *loggerImpl) Info(msg string, fields ...map[string]interface{}) {
	l.log(INFO, msg, fields...)
}

func (l *loggerImpl) Warn(msg string, fields ...map[string]interface{}) {
	l.log(WARN, msg, fields...)
}

func (l *loggerImpl) Error(msg string, fields ...map[string]interface{}) {
	l.log(ERROR, msg, fields...)
}

func (l *loggerImpl) WithFields(fields map[string]interface{}) Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	newLogger := &loggerImpl{
		options: l.options,
		fields:  make(map[string]interface{}),
		rotator: l.rotator,
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *loggerImpl) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.options.Level = level
}

func (l *loggerImpl) Close() error {
	if l.rotator != nil {
		return l.rotator.Close()
	}
	return nil
}

// Write implements io.Writer interface for log rotation
func (r *logRotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if rotation is needed
	if r.maxSize > 0 {
		info, err := r.file.Stat()
		if err != nil {
			return 0, err
		}

		if info.Size() >= r.maxSize {
			if err := r.rotate(); err != nil {
				return 0, err
			}
		}
	}

	return r.file.Write(p)
}

// rotate performs log rotation
func (r *logRotator) rotate() error {
	// Close current file
	if err := r.file.Close(); err != nil {
		return err
	}

	// Generate new filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	dir := filepath.Dir(r.file.Name())
	base := filepath.Base(r.file.Name())
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	
	newPath := filepath.Join(dir, fmt.Sprintf("%s.%s%s", name, timestamp, ext))

	// Rename current file
	if err := os.Rename(r.file.Name(), newPath); err != nil {
		return err
	}

	// Create new file
	file, err := os.OpenFile(r.file.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	r.file = file

	// Clean up old files if needed
	if r.maxBackups > 0 || r.maxAge > 0 {
		go r.cleanup()
	}

	return nil
}

// cleanup removes old log files based on maxBackups and maxAge
func (r *logRotator) cleanup() {
	dir := filepath.Dir(r.file.Name())
	base := filepath.Base(r.file.Name())
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	pattern := filepath.Join(dir, name+".*"+ext)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	// Sort by modification time
	type fileInfo struct {
		path    string
		modTime time.Time
	}
	var files []fileInfo

	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		files = append(files, fileInfo{match, info.ModTime()})
	}

	// Remove files based on maxAge
	if r.maxAge > 0 {
		cutoff := time.Now().Add(-r.maxAge)
		for _, f := range files {
			if f.modTime.Before(cutoff) {
				os.Remove(f.path)
			}
		}
	}

	// Remove files based on maxBackups
	if r.maxBackups > 0 && len(files) > r.maxBackups {
		// Sort by modification time (oldest first)
		for i := 0; i < len(files)-r.maxBackups; i++ {
			os.Remove(files[i].path)
		}
	}
}

// Close closes the log file
func (r *logRotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Close()
} 