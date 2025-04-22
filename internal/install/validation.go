package install

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// namePattern defines valid characters for tool names
	namePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-_\.]+$`)
	// versionPattern defines valid version strings
	versionPattern = regexp.MustCompile(`^(latest|stable|\d+\.\d+(\.\d+)?(-[a-zA-Z0-9]+)?)$`)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// validateTool validates a tool configuration
func validateTool(tool *Tool) error {
	var errors []string

	// Validate Name
	if tool.Name == "" {
		errors = append(errors, (&ValidationError{
			Field:   "Name",
			Message: "cannot be empty",
		}).Error())
	} else if !namePattern.MatchString(tool.Name) {
		errors = append(errors, (&ValidationError{
			Field:   "Name",
			Message: "contains invalid characters",
		}).Error())
	}

	// Validate PackageName
	if tool.PackageName == "" {
		errors = append(errors, (&ValidationError{
			Field:   "PackageName",
			Message: "cannot be empty",
		}).Error())
	} else if !namePattern.MatchString(tool.PackageName) {
		errors = append(errors, (&ValidationError{
			Field:   "PackageName",
			Message: "contains invalid characters",
		}).Error())
	}

	// Validate Version if specified
	if tool.Version != "" && !versionPattern.MatchString(tool.Version) {
		errors = append(errors, (&ValidationError{
			Field:   "Version",
			Message: "invalid version format",
		}).Error())
	}

	// Validate Dependencies
	for i, dep := range tool.Dependencies {
		if dep == "" {
			errors = append(errors, (&ValidationError{
				Field:   fmt.Sprintf("Dependencies[%d]", i),
				Message: "cannot be empty",
			}).Error())
		} else if !namePattern.MatchString(dep) {
			errors = append(errors, (&ValidationError{
				Field:   fmt.Sprintf("Dependencies[%d]", i),
				Message: "contains invalid characters",
			}).Error())
		}
	}

	// Validate PostInstall commands
	for i, cmd := range tool.PostInstall {
		if cmd == "" {
			errors = append(errors, (&ValidationError{
				Field:   fmt.Sprintf("PostInstall[%d]", i),
				Message: "cannot be empty",
			}).Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
} 