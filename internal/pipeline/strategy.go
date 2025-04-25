package pipeline

import (
	"fmt"
	"regexp"
	"strings"
)

// VersionConstraint defines version requirements for a package
type VersionConstraint struct {
	// Minimum version required
	MinVersion string
	// Maximum version allowed
	MaxVersion string
	// Exact version required (if set, overrides MinVersion and MaxVersion)
	ExactVersion string
	// Version pattern to match (regex)
	Pattern string
}

// Validate checks if a version meets the constraints
func (vc *VersionConstraint) Validate(version string) error {
	if vc.ExactVersion != "" {
		if version != vc.ExactVersion {
			return fmt.Errorf("version %s does not match required version %s", 
				version, vc.ExactVersion)
		}
		return nil
	}

	if vc.Pattern != "" {
		matched, err := regexp.MatchString(vc.Pattern, version)
		if err != nil {
			return fmt.Errorf("invalid version pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("version %s does not match pattern %s", 
				version, vc.Pattern)
		}
		return nil
	}

	// TODO: Implement proper version comparison
	// For now, just do string comparison
	if vc.MinVersion != "" && version < vc.MinVersion {
		return fmt.Errorf("version %s is below minimum required version %s", 
			version, vc.MinVersion)
	}
	if vc.MaxVersion != "" && version > vc.MaxVersion {
		return fmt.Errorf("version %s is above maximum allowed version %s", 
			version, vc.MaxVersion)
	}

	return nil
}

// Command represents a shell command with validation
type Command struct {
	// The command to execute
	Command string
	// Description of what the command does
	Description string
	// Whether the command requires sudo
	RequiresSudo bool
	// Expected exit code (0 if not specified)
	ExpectedExitCode int
	// Timeout in seconds (0 for no timeout)
	Timeout int
	// Retry count (0 for no retries)
	RetryCount int
	// Retry delay in seconds
	RetryDelay int
}

// Validate checks if the command is valid
func (c *Command) Validate() error {
	if c.Command == "" {
		return fmt.Errorf("command cannot be empty")
	}
	if c.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	if c.RetryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}
	if c.RetryDelay < 0 {
		return fmt.Errorf("retry delay cannot be negative")
	}
	return nil
}

// InstallStrategy defines how to install a tool
type InstallStrategy struct {
	// Package names for different package managers
	PackageNames map[string]string
	// Version constraints
	VersionConstraints map[string]*VersionConstraint
	// Pre-install commands
	PreInstall []Command
	// Post-install commands
	PostInstall []Command
	// Custom installation commands
	CustomInstall []Command
	// Rollback commands in case of failure
	Rollback []Command
	// Environment variables to set
	Env map[string]string
	// Files to create/modify
	Files map[string]string
	// Directories to create
	Directories []string
	// Dependencies to install
	Dependencies []string
	// System dependencies to install
	SystemDependencies []string
}

// Validate checks if the installation strategy is valid
func (s *InstallStrategy) Validate() error {
	// Validate package names
	if len(s.PackageNames) == 0 && len(s.CustomInstall) == 0 {
		return fmt.Errorf("either package names or custom install commands must be specified")
	}

	// Validate version constraints
	for pkg, constraint := range s.VersionConstraints {
		if constraint == nil {
			return fmt.Errorf("version constraint for %s cannot be nil", pkg)
		}
	}

	// Validate commands
	for i, cmd := range s.PreInstall {
		if err := cmd.Validate(); err != nil {
			return fmt.Errorf("invalid pre-install command %d: %w", i, err)
		}
	}
	for i, cmd := range s.PostInstall {
		if err := cmd.Validate(); err != nil {
			return fmt.Errorf("invalid post-install command %d: %w", i, err)
		}
	}
	for i, cmd := range s.CustomInstall {
		if err := cmd.Validate(); err != nil {
			return fmt.Errorf("invalid custom install command %d: %w", i, err)
		}
	}
	for i, cmd := range s.Rollback {
		if err := cmd.Validate(); err != nil {
			return fmt.Errorf("invalid rollback command %d: %w", i, err)
		}
	}

	// Validate environment variables
	for key, value := range s.Env {
		if key == "" {
			return fmt.Errorf("environment variable key cannot be empty")
		}
		if value == "" {
			return fmt.Errorf("environment variable value cannot be empty")
		}
	}

	// Validate files
	for path, content := range s.Files {
		if path == "" {
			return fmt.Errorf("file path cannot be empty")
		}
		if content == "" {
			return fmt.Errorf("file content cannot be empty")
		}
	}

	// Validate directories
	for _, dir := range s.Directories {
		if dir == "" {
			return fmt.Errorf("directory path cannot be empty")
		}
	}

	return nil
}

// GetPackageName returns the package name for the given package manager
func (s *InstallStrategy) GetPackageName(pkgManager string) (string, error) {
	if name, ok := s.PackageNames[pkgManager]; ok {
		return name, nil
	}
	if name, ok := s.PackageNames["default"]; ok {
		return name, nil
	}
	return "", fmt.Errorf("no package name found for package manager %s", pkgManager)
}

// GetVersionConstraint returns the version constraint for the given package
func (s *InstallStrategy) GetVersionConstraint(pkg string) *VersionConstraint {
	if constraint, ok := s.VersionConstraints[pkg]; ok {
		return constraint
	}
	return nil
}

// HasCustomInstall returns whether the strategy uses custom installation
func (s *InstallStrategy) HasCustomInstall() bool {
	return len(s.CustomInstall) > 0
}

// HasRollback returns whether the strategy has rollback commands
func (s *InstallStrategy) HasRollback() bool {
	return len(s.Rollback) > 0
}

// String returns a string representation of the strategy
func (s *InstallStrategy) String() string {
	var parts []string

	if len(s.PackageNames) > 0 {
		parts = append(parts, fmt.Sprintf("PackageNames: %v", s.PackageNames))
	}
	if len(s.VersionConstraints) > 0 {
		parts = append(parts, fmt.Sprintf("VersionConstraints: %v", s.VersionConstraints))
	}
	if len(s.PreInstall) > 0 {
		parts = append(parts, fmt.Sprintf("PreInstall: %d commands", len(s.PreInstall)))
	}
	if len(s.PostInstall) > 0 {
		parts = append(parts, fmt.Sprintf("PostInstall: %d commands", len(s.PostInstall)))
	}
	if len(s.CustomInstall) > 0 {
		parts = append(parts, fmt.Sprintf("CustomInstall: %d commands", len(s.CustomInstall)))
	}
	if len(s.Rollback) > 0 {
		parts = append(parts, fmt.Sprintf("Rollback: %d commands", len(s.Rollback)))
	}
	if len(s.Env) > 0 {
		parts = append(parts, fmt.Sprintf("Env: %d variables", len(s.Env)))
	}
	if len(s.Files) > 0 {
		parts = append(parts, fmt.Sprintf("Files: %d files", len(s.Files)))
	}
	if len(s.Directories) > 0 {
		parts = append(parts, fmt.Sprintf("Directories: %d directories", len(s.Directories)))
	}
	if len(s.Dependencies) > 0 {
		parts = append(parts, fmt.Sprintf("Dependencies: %d packages", len(s.Dependencies)))
	}
	if len(s.SystemDependencies) > 0 {
		parts = append(parts, fmt.Sprintf("SystemDependencies: %d packages", len(s.SystemDependencies)))
	}

	return strings.Join(parts, ", ")
} 