# Logging Package

The logging package provides a comprehensive logging system for the Bootstrap CLI application. It implements a flexible, structured logging approach with support for different log levels, command execution logging, and installation-specific logging.

## Components

### Core Interface
- Defined in `internal/interfaces/logger.go`
- Provides the base contract for all logger implementations
- Includes methods for different log levels and command execution logging

### Main Logger
- Located in `logger.go`
- Implements the core `Logger` interface
- Provides structured logging with metadata support
- Supports different log levels (Debug, Info, Warn, Error)
- Includes command execution logging

### Installation Logger
- Located in `install.go`
- Specialized logger for installation operations
- Provides installation-specific logging format
- Includes debug mode toggle
- Implements the `Logger` interface

### Logger Adapter
- Located in `adapter.go`
- Allows adapting any `Logger` implementation to work with installation-specific code
- Useful for bridging between different logger implementations
- Implements the `Logger` interface

### Mock Logger
- Located in `mock.go`
- Used for testing
- Captures all log messages for verification
- Implements the `Logger` interface

## Usage

### Basic Logging
```go
logger := log.New(log.InfoLevel)
logger.Info("Starting installation")
logger.Debug("Debug information")
logger.Warn("Warning message")
logger.Error("Error occurred")
```

### Installation Logging
```go
logger := log.NewInstallLogger(true) // true enables debug mode
logger.Info("Installing package")
logger.CommandStart("apt-get install package", 1, 3)
logger.CommandSuccess("apt-get install package", time.Second)
```

### Using the Adapter
```go
baseLogger := log.New(log.InfoLevel)
adapter := log.NewLogAdapter(baseLogger)
// Use adapter where an installation logger is expected
```

## Best Practices

1. Use the appropriate logger for your context:
   - Use `New()` for general application logging
   - Use `NewInstallLogger()` for installation-specific logging
   - Use `NewLogAdapter()` when you need to adapt between logger types

2. Enable debug logging only when needed:
   - Debug logs can be verbose
   - Use debug mode for troubleshooting

3. Use structured logging for better log analysis:
   - Include relevant context in log messages
   - Use consistent log levels appropriately

4. For testing:
   - Use the mock logger to verify logging behavior
   - Capture and verify log messages in tests 