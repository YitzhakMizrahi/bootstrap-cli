package system

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// WriteConfigFile writes content to a configuration file with proper permissions
// Creates parent directories if they don't exist
func WriteConfigFile(path string, content []byte) error {
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write file with restricted permissions
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	return nil
}

// IsRoot checks if the current user has root privileges.
func IsRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		// Consider logging this error, but for simplicity, assume not root if error occurs
		return false
	}
	return currentUser.Uid == "0"
} 