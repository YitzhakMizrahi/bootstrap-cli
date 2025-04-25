// Package validation provides tools for validating various components of the bootstrap-cli,
// including tool configurations, package names, and version strings.
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
)

var (
	// namePattern defines valid characters for tool names
	namePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_\.]+$`)
	// versionPattern defines valid version strings
	versionPattern = regexp.MustCompile(`^(latest|stable|\d+\.\d+(\.\d+)?(-[a-zA-Z0-9]+)?)$`)
)

// Error represents a validation error
type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateTool validates a tool configuration
func ValidateTool(tool *install.Tool) error {
	var errors []string

	// Validate Name
	if tool.Name == "" {
		errors = append(errors, (&Error{
			Field:   "Name",
			Message: "cannot be empty",
		}).Error())
	} else if !namePattern.MatchString(tool.Name) {
		errors = append(errors, (&Error{
			Field:   "Name",
			Message: "contains invalid characters",
		}).Error())
	}

	// Validate PackageName
	if tool.PackageName == "" {
		errors = append(errors, (&Error{
			Field:   "PackageName",
			Message: "cannot be empty",
		}).Error())
	} else if !namePattern.MatchString(tool.PackageName) {
		errors = append(errors, (&Error{
			Field:   "PackageName",
			Message: "contains invalid characters",
		}).Error())
	}

	// Validate Version if specified
	if tool.Version != "" && !versionPattern.MatchString(tool.Version) {
		errors = append(errors, (&Error{
			Field:   "Version",
			Message: "invalid version format",
		}).Error())
	}

	// Validate Dependencies
	for i, dep := range tool.Dependencies {
		if dep == "" {
			errors = append(errors, (&Error{
				Field:   fmt.Sprintf("Dependencies[%d]", i),
				Message: "cannot be empty",
			}).Error())
		} else if !namePattern.MatchString(dep) {
			errors = append(errors, (&Error{
				Field:   fmt.Sprintf("Dependencies[%d]", i),
				Message: "contains invalid characters",
			}).Error())
		}
	}

	// Validate PostInstall commands
	for i, cmd := range tool.PostInstall {
		if cmd.Command == "" {
			errors = append(errors, (&Error{
				Field:   fmt.Sprintf("PostInstall[%d].Command", i),
				Message: "cannot be empty",
			}).Error())
		}
		// Description is optional, so we don't validate it
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
} 