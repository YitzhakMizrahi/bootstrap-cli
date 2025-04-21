package shell

import "fmt"

// ErrShellNotSupported indicates that the requested shell is not supported
type ErrShellNotSupported struct {
	Shell Type
}

func (e *ErrShellNotSupported) Error() string {
	return fmt.Sprintf("shell type %q is not supported", e.Shell)
}

// ErrPluginManagerNotSupported indicates that the requested plugin manager is not supported
type ErrPluginManagerNotSupported struct {
	Manager   PluginManager
	Message   string
}

func (e *ErrPluginManagerNotSupported) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("plugin manager %q is not supported: %s", e.Manager, e.Message)
	}
	return fmt.Sprintf("plugin manager %q is not supported", e.Manager)
}

// ErrPluginManagerInstallation indicates a failure during plugin manager installation
type ErrPluginManagerInstallation struct {
	Manager PluginManager
	Err     error
}

func (e *ErrPluginManagerInstallation) Error() string {
	return fmt.Sprintf("failed to install plugin manager %q: %v", e.Manager, e.Err)
}

func (e *ErrPluginManagerInstallation) Unwrap() error {
	return e.Err
}

// ErrConfigGeneration indicates a failure during configuration generation
type ErrConfigGeneration struct {
	Shell   Type
	Message string
	Err     error
}

func (e *ErrConfigGeneration) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("failed to generate configuration for shell %q: %s: %v", e.Shell, e.Message, e.Err)
	}
	return fmt.Sprintf("failed to generate configuration for shell %q: %v", e.Shell, e.Err)
}

func (e *ErrConfigGeneration) Unwrap() error {
	return e.Err
}

// ErrBackup indicates a failure during configuration backup
type ErrBackup struct {
	Path string
	Err  error
}

func (e *ErrBackup) Error() string {
	return fmt.Sprintf("failed to backup configuration at %q: %v", e.Path, e.Err)
}

func (e *ErrBackup) Unwrap() error {
	return e.Err
}

// ErrRestore indicates a failure during configuration restore
type ErrRestore struct {
	Path    string
	Message string
	Err     error
}

func (e *ErrRestore) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("failed to restore configuration from %q: %s: %v", e.Path, e.Message, e.Err)
	}
	return fmt.Sprintf("failed to restore configuration from %q: %v", e.Path, e.Err)
}

func (e *ErrRestore) Unwrap() error {
	return e.Err
}

// ErrInvalidConfig indicates an invalid shell configuration
type ErrInvalidConfig struct {
	Field string
	Value interface{}
}

func (e *ErrInvalidConfig) Error() string {
	return fmt.Sprintf("invalid configuration: field %q has invalid value %v", e.Field, e.Value)
}

// ErrPluginNotFound indicates that a requested plugin is not found
type ErrPluginNotFound struct {
	Plugin string
}

func (e *ErrPluginNotFound) Error() string {
	return fmt.Sprintf("plugin %q not found", e.Plugin)
}

// ErrTemplateNotFound indicates that a requested template is not found
type ErrTemplateNotFound struct {
	Shell Type
}

func (e *ErrTemplateNotFound) Error() string {
	return fmt.Sprintf("template not found for shell %q", e.Shell)
} 