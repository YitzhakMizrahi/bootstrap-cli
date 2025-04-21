package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Config holds application configuration
type Config struct {
	LogLevel    string `json:"log_level"`
	LogFile     string `json:"log_file"`
	PluginDir   string `json:"plugin_dir"`
	DataDir     string `json:"data_dir"`
	MaxPlugins  int    `json:"max_plugins"`
	AutoReload  bool   `json:"auto_reload"`
	mu          sync.RWMutex
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:    "info",
		LogFile:     "logs/app.log",
		PluginDir:   "plugins",
		DataDir:     "data",
		MaxPlugins:  10,
		AutoReload:  true,
	}
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// Save saves configuration to a file
func (c *Config) Save(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Get returns a configuration value
func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch key {
	case "log_level":
		return c.LogLevel
	case "log_file":
		return c.LogFile
	case "plugin_dir":
		return c.PluginDir
	case "data_dir":
		return c.DataDir
	case "max_plugins":
		return c.MaxPlugins
	case "auto_reload":
		return c.AutoReload
	default:
		return nil
	}
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch key {
	case "log_level":
		if str, ok := value.(string); ok {
			c.LogLevel = str
		} else {
			return fmt.Errorf("invalid value type for log_level")
		}
	case "log_file":
		if str, ok := value.(string); ok {
			c.LogFile = str
		} else {
			return fmt.Errorf("invalid value type for log_file")
		}
	case "plugin_dir":
		if str, ok := value.(string); ok {
			c.PluginDir = str
		} else {
			return fmt.Errorf("invalid value type for plugin_dir")
		}
	case "data_dir":
		if str, ok := value.(string); ok {
			c.DataDir = str
		} else {
			return fmt.Errorf("invalid value type for data_dir")
		}
	case "max_plugins":
		if num, ok := value.(int); ok {
			c.MaxPlugins = num
		} else {
			return fmt.Errorf("invalid value type for max_plugins")
		}
	case "auto_reload":
		if b, ok := value.(bool); ok {
			c.AutoReload = b
		} else {
			return fmt.Errorf("invalid value type for auto_reload")
		}
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return nil
} 