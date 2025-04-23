package internal

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

// ExtractEmbeddedConfigs extracts embedded configurations to a temporary directory
func ExtractEmbeddedConfigs(destDir string) error {
	// Create the base config directory
	configDir := filepath.Join(destDir, "defaults", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Extract all embedded configs
	err := extractDir(configFiles, "defaults/config", destDir)
	if err != nil {
		return fmt.Errorf("failed to extract embedded configs: %w", err)
	}

	return nil
}

// extractDir recursively extracts files from the embedded filesystem
func extractDir(efs embed.FS, sourceDir, destDir string) error {
	entries, err := efs.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Create directory and recurse
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
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
	data, err := efs.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(destPath, data, 0644)
} 