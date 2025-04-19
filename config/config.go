// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/types"
	"gopkg.in/yaml.v3"
)

const ConfigPath = "~/.config/bootstrap/config.yaml"

// Save writes the user configuration to config.yaml
func Save(cfg types.UserConfig) error {
	expanded := expandPath(ConfigPath)
	dir := filepath.Dir(expanded)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	f, err := os.Create(expanded)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("💾 Config saved to %s\n", expanded)
	return nil
}

// expandPath replaces ~ with the user’s home directory
func expandPath(p string) string {
	if len(p) > 1 && p[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
