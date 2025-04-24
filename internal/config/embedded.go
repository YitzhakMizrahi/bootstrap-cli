package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed defaults/tools/schema.yaml
//go:embed defaults/tools/modern/*.yaml
//go:embed defaults/dotfiles/shell/*.yaml
//go:embed defaults/languages/*.yaml
var configFS embed.FS

// ExtractEmbeddedConfigs extracts embedded configurations to the loader's base directory
func (l *ConfigLoader) ExtractEmbeddedConfigs() error {
	// Extract all embedded configs to the base directory
	err := extractDir(configFS, "defaults", l.baseDir)
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

// GetEmbeddedConfigPath returns the path to a specific embedded config file
func (l *ConfigLoader) GetEmbeddedConfigPath(relativePath string) (string, error) {
	fullPath := filepath.Join(l.baseDir, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("embedded config not found: %s", relativePath)
	}

	return fullPath, nil
} 