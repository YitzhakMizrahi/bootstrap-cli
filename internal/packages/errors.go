package packages

import (
	"fmt"
)

// PackageError represents a generic package management error
type PackageError struct {
	Operation string // The operation that failed (install, uninstall, etc.)
	Package   string // The package name
	Err       error  // The underlying error
}

func (e *PackageError) Error() string {
	return fmt.Sprintf("package %s error during %s: %v", e.Package, e.Operation, e.Err)
}

// Unwrap returns the underlying error
func (e *PackageError) Unwrap() error {
	return e.Err
}

// NewPackageError creates a new PackageError
func NewPackageError(operation, pkg string, err error) error {
	return &PackageError{
		Operation: operation,
		Package:   pkg,
		Err:       err,
	}
}

// SystemNotSupportedError is returned when the system is not supported
type SystemNotSupportedError struct {
	System string
}

func (e *SystemNotSupportedError) Error() string {
	return fmt.Sprintf("system not supported: %s", e.System)
}

// NewSystemNotSupportedError creates a new SystemNotSupportedError
func NewSystemNotSupportedError(system string) error {
	return &SystemNotSupportedError{
		System: system,
	}
}

// PackageNotFoundError is returned when a package is not found
type PackageNotFoundError struct {
	Package string
}

func (e *PackageNotFoundError) Error() string {
	return fmt.Sprintf("package not found: %s", e.Package)
}

// NewPackageNotFoundError creates a new PackageNotFoundError
func NewPackageNotFoundError(pkg string) error {
	return &PackageNotFoundError{
		Package: pkg,
	}
}

// CommandExecutionError is returned when a package manager command fails
type CommandExecutionError struct {
	Command string
	Output  string
	Err     error
}

func (e *CommandExecutionError) Error() string {
	return fmt.Sprintf("command execution failed: %s, output: %s, error: %v", e.Command, e.Output, e.Err)
}

// Unwrap returns the underlying error
func (e *CommandExecutionError) Unwrap() error {
	return e.Err
}

// NewCommandExecutionError creates a new CommandExecutionError
func NewCommandExecutionError(command, output string, err error) error {
	return &CommandExecutionError{
		Command: command,
		Output:  output,
		Err:     err,
	}
} 