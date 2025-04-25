// Package config provides configuration management functionality for the bootstrap-cli,
// including loading, parsing, and managing configuration files, as well as handling
// embedded default configurations.
package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

// ExtractEmbeddedConfigs extracts embedded configurations to the loader's base directory
func (l *Loader) ExtractEmbeddedConfigs() error {
	// Create the defaults directory if it doesn't exist
	defaultsDir := filepath.Join(l.baseDir, "defaults")
	if err := os.MkdirAll(defaultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create defaults directory: %w", err)
	}

	// Extract all embedded configs to the base directory
	err := extractDir(l.configFS, "defaults", defaultsDir)
	if err != nil {
		return fmt.Errorf("failed to extract embedded configs: %w", err)
	}

	return nil
}

// extractDir recursively extracts files from the embedded filesystem
func extractDir(efs embed.FS, sourceDir, destDir string) error {
	entries, err := efs.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", sourceDir, err)
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Create directory and recurse
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			if err := extractDir(efs, sourcePath, destPath); err != nil {
				return err
			}
		} else {
			// Extract file
			if err := extractFile(efs, sourcePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// extractFile extracts a single file from the embedded filesystem
func extractFile(efs embed.FS, sourcePath, destPath string) error {
	// Read the embedded file
	data, err := efs.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", sourcePath, err)
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
	}

	// Write the file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	return nil
}

// GetEmbeddedConfigPath returns the path to an embedded config file
func (l *Loader) GetEmbeddedConfigPath(relativePath string) (string, error) {
	// Check if the file exists in the embedded filesystem
	if _, err := l.configFS.Open(filepath.Join("defaults", relativePath)); err != nil {
		return "", fmt.Errorf("embedded config not found: %s", relativePath)
	}

	// Return the path in the extracted directory
	return filepath.Join(l.baseDir, "defaults", relativePath), nil
} 