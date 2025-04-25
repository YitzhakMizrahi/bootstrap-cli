package pipeline

import "fmt"

// InstallationError represents an error during the installation process
type InstallationError struct {
	Tool    string
	Step    string
	Message string
	Cause   error
}

func (e *InstallationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("installation error for %s during %s: %s (caused by: %v)", 
			e.Tool, e.Step, e.Message, e.Cause)
	}
	return fmt.Sprintf("installation error for %s during %s: %s", 
		e.Tool, e.Step, e.Message)
}

// Unwrap returns the underlying error
func (e *InstallationError) Unwrap() error {
	return e.Cause
}

// VerificationError represents an error during tool verification
type VerificationError struct {
	Tool    string
	Message string
	Cause   error
}

func (e *VerificationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("verification error for %s: %s (caused by: %v)", 
			e.Tool, e.Message, e.Cause)
	}
	return fmt.Sprintf("verification error for %s: %s", 
		e.Tool, e.Message)
}

// Unwrap returns the underlying error
func (e *VerificationError) Unwrap() error {
	return e.Cause
}

// PackageManagerError represents an error from the package manager
type PackageManagerError struct {
	Operation string
	Package   string
	Message   string
	Cause     error
}

func (e *PackageManagerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("package manager error during %s of %s: %s (caused by: %v)", 
			e.Operation, e.Package, e.Message, e.Cause)
	}
	return fmt.Sprintf("package manager error during %s of %s: %s", 
		e.Operation, e.Package, e.Message)
}

// Unwrap returns the underlying error
func (e *PackageManagerError) Unwrap() error {
	return e.Cause
}

// PlatformError represents an error related to platform compatibility
type PlatformError struct {
	Platform string
	Message  string
	Cause    error
}

func (e *PlatformError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("platform error for %s: %s (caused by: %v)", 
			e.Platform, e.Message, e.Cause)
	}
	return fmt.Sprintf("platform error for %s: %s", 
		e.Platform, e.Message)
}

// Unwrap returns the underlying error
func (e *PlatformError) Unwrap() error {
	return e.Cause
}

// Helper functions to create errors
func NewInstallationError(tool, step, message string, cause error) error {
	return &InstallationError{
		Tool:    tool,
		Step:    step,
		Message: message,
		Cause:   cause,
	}
}

func NewVerificationError(tool, message string, cause error) error {
	return &VerificationError{
		Tool:    tool,
		Message: message,
		Cause:   cause,
	}
}

func NewPackageManagerError(operation, pkg, message string, cause error) error {
	return &PackageManagerError{
		Operation: operation,
		Package:   pkg,
		Message:   message,
		Cause:     cause,
	}
}

func NewPlatformError(platform, message string, cause error) error {
	return &PlatformError{
		Platform: platform,
		Message:  message,
		Cause:    cause,
	}
} 