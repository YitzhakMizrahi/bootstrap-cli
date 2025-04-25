package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

func TestConfig_GenerateConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "shell-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set TMPDIR for the test
	oldTmpDir := os.Getenv("TMPDIR")
	os.Setenv("TMPDIR", tempDir)
	defer os.Setenv("TMPDIR", oldTmpDir)

	mockLogger := log.NewMockLogger()
	tests := []struct {
		name     string
		shell    string
		setup    func(*Config)
		expected []string
	}{
		{
			name:  "bash_basic_config",
			shell: "bash",
			setup: func(c *Config) {
				c.AddEnvVar("TEST_VAR", "test_value")
				c.AddAlias("ll", "ls -la")
				c.AddFunction("testfunc", "echo 'test'")
				c.AddPath("/test/path")
			},
			expected: []string{
				`export TEST_VAR=test_value`,
				`export PATH=/test/path:$PATH`,
				`alias ll='ls -la'`,
				`testfunc() {`,
				`echo 'test'`,
				`}`,
			},
		},
		{
			name:  "fish_basic_config",
			shell: "fish",
			setup: func(c *Config) {
				c.AddEnvVar("TEST_VAR", "test_value")
				c.AddAlias("ll", "ls -la")
				c.AddFunction("testfunc", "echo 'test'")
				c.AddPath("/test/path")
			},
			expected: []string{
				`set -gx TEST_VAR test_value`,
				`fish_add_path /test/path`,
				`alias ll='ls -la'`,
				`function testfunc`,
				`echo 'test'`,
				`end`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(tt.shell, mockLogger)
			tt.setup(c)

			content, err := c.GenerateConfig()
			if err != nil {
				t.Fatalf("GenerateConfig() error = %v", err)
			}

			for _, expected := range tt.expected {
				if !strings.Contains(content, expected) {
					t.Errorf("GenerateConfig() content does not contain %q", expected)
				}
			}
		})
	}
}

func TestConfig_Apply(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "shell-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set TMPDIR for the test
	oldTmpDir := os.Getenv("TMPDIR")
	os.Setenv("TMPDIR", tempDir)
	defer os.Setenv("TMPDIR", oldTmpDir)

	mockLogger := log.NewMockLogger()

	tests := []struct {
		name          string
		shell         string
		setup         func(*Config)
		expectedCmd   string
		expectError   bool
	}{
		{
			name:  "bash_config",
			shell: "bash",
			setup: func(c *Config) {
				c.AddEnvVar("TEST_VAR", "test_value")
			},
			expectedCmd: ". " + filepath.Join(tempDir, "bootstrap-cli-bash-config"),
			expectError: false,
		},
		{
			name:  "fish_config",
			shell: "fish",
			setup: func(c *Config) {
				c.AddEnvVar("TEST_VAR", "test_value")
			},
			expectedCmd: "source " + filepath.Join(tempDir, "bootstrap-cli-fish-config"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(tt.shell, mockLogger)
			tt.setup(c)

			cmd, err := c.Apply()
			if tt.expectError {
				if err == nil {
					t.Error("Apply() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Apply() unexpected error = %v", err)
			}

			if cmd != tt.expectedCmd {
				t.Errorf("Apply() got = %v, want %v", cmd, tt.expectedCmd)
			}

			// Check if file exists and contains content
			content, err := os.ReadFile(c.GetTempConfigFile())
			if err != nil {
				t.Fatalf("Failed to read temp config file: %v", err)
			}
			if !strings.Contains(string(content), "TEST_VAR") {
				t.Error("Config file does not contain expected content")
			}
		})
	}
} 