package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// mockShellManager implements Manager for testing
type mockShellManager struct {
	shells map[Shell]*ShellInfo
}

func (m *mockShellManager) DetectCurrent() (*ShellInfo, error) {
	return nil, nil
}

func (m *mockShellManager) ListAvailable() ([]*ShellInfo, error) {
	return nil, nil
}

func (m *mockShellManager) IsInstalled(shell Shell) bool {
	_, ok := m.shells[shell]
	return ok
}

func (m *mockShellManager) GetInfo(shell Shell) (*ShellInfo, error) {
	info, ok := m.shells[shell]
	if !ok {
		return nil, fmt.Errorf("shell %s not found", shell)
	}
	return info, nil
}

func TestConfigWriter(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create mock shell manager
	mockManager := &mockShellManager{
		shells: map[Shell]*ShellInfo{
			Bash: {
				Type:        Bash,
				Path:        "/bin/bash",
				Version:     "5.0.0",
				IsDefault:   true,
				IsAvailable: true,
				ConfigFiles: []string{filepath.Join(tmpDir, ".bashrc")},
			},
			Zsh: {
				Type:        Zsh,
				Path:        "/bin/zsh",
				Version:     "5.8",
				IsDefault:   false,
				IsAvailable: true,
				ConfigFiles: []string{filepath.Join(tmpDir, ".zshrc")},
			},
			Fish: {
				Type:        Fish,
				Path:        "/usr/bin/fish",
				Version:     "3.3.1",
				IsDefault:   false,
				IsAvailable: true,
				ConfigFiles: []string{filepath.Join(tmpDir, ".config/fish/config.fish")},
			},
		},
	}

	// Create config writer
	writer := NewConfigWriter(mockManager)

	// Test writing config
	t.Run("WriteConfig", func(t *testing.T) {
		tests := []struct {
			name     string
			shell    Shell
			config   string
			strategy DotfilesStrategy
			wantErr  bool
		}{
			{
				name:     "bash simple config",
				shell:    Bash,
				config:   "export TEST_VAR=value",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "zsh simple config",
				shell:    Zsh,
				config:   "export TEST_VAR=value",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "fish simple config",
				shell:    Fish,
				config:   "set -gx TEST_VAR value",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "unknown shell",
				shell:    "unknown",
				config:   "test config",
				strategy: ReplaceExisting,
				wantErr:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := writer.WriteConfig(tt.shell, tt.config, tt.strategy)
				if (err != nil) != tt.wantErr {
					t.Errorf("WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				// Check if config was written
				info, _ := mockManager.GetInfo(tt.shell)
				content, err := os.ReadFile(info.ConfigFiles[0])
				if err != nil {
					t.Errorf("Failed to read config file: %v", err)
					return
				}

				if !strings.Contains(string(content), tt.config) {
					t.Errorf("Config file does not contain expected content. Got: %s", string(content))
				}
			})
		}
	})

	// Test adding to PATH
	t.Run("AddToPath", func(t *testing.T) {
		tests := []struct {
			name     string
			shell    Shell
			path     string
			strategy DotfilesStrategy
			wantErr  bool
		}{
			{
				name:     "bash add path",
				shell:    Bash,
				path:     "/usr/local/bin",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "zsh add path",
				shell:    Zsh,
				path:     "/usr/local/bin",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "fish add path",
				shell:    Fish,
				path:     "/usr/local/bin",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := writer.AddToPath(tt.shell, tt.path, tt.strategy)
				if (err != nil) != tt.wantErr {
					t.Errorf("AddToPath() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				// Check if path was added
				info, _ := mockManager.GetInfo(tt.shell)
				content, err := os.ReadFile(info.ConfigFiles[0])
				if err != nil {
					t.Errorf("Failed to read config file: %v", err)
					return
				}

				expectedPath := tt.path
				if tt.shell == Fish {
					expectedPath = "set -gx PATH " + tt.path + " $PATH"
				} else {
					expectedPath = "export PATH=\"" + tt.path + ":$PATH\""
				}

				if !strings.Contains(string(content), expectedPath) {
					t.Errorf("Config file does not contain expected PATH. Got: %s", string(content))
				}
			})
		}
	})

	// Test setting environment variables
	t.Run("SetEnvVar", func(t *testing.T) {
		tests := []struct {
			name     string
			shell    Shell
			key      string
			value    string
			strategy DotfilesStrategy
			wantErr  bool
		}{
			{
				name:     "bash set env var",
				shell:    Bash,
				key:      "EDITOR",
				value:    "vim",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "zsh set env var",
				shell:    Zsh,
				key:      "EDITOR",
				value:    "vim",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "fish set env var",
				shell:    Fish,
				key:      "EDITOR",
				value:    "vim",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := writer.SetEnvVar(tt.shell, tt.key, tt.value, tt.strategy)
				if (err != nil) != tt.wantErr {
					t.Errorf("SetEnvVar() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				// Check if env var was set
				info, _ := mockManager.GetInfo(tt.shell)
				content, err := os.ReadFile(info.ConfigFiles[0])
				if err != nil {
					t.Errorf("Failed to read config file: %v", err)
					return
				}

				expectedVar := ""
				if tt.shell == Fish {
					expectedVar = "set -gx " + tt.key + " " + tt.value
				} else {
					expectedVar = "export " + tt.key + "=\"" + tt.value + "\""
				}

				if !strings.Contains(string(content), expectedVar) {
					t.Errorf("Config file does not contain expected env var. Got: %s", string(content))
				}
			})
		}
	})

	// Test adding aliases
	t.Run("AddAlias", func(t *testing.T) {
		tests := []struct {
			name     string
			shell    Shell
			alias    string
			command  string
			strategy DotfilesStrategy
			wantErr  bool
		}{
			{
				name:     "bash add alias",
				shell:    Bash,
				alias:    "ll",
				command:  "ls -la",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "zsh add alias",
				shell:    Zsh,
				alias:    "ll",
				command:  "ls -la",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "fish add alias",
				shell:    Fish,
				alias:    "ll",
				command:  "ls -la",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := writer.AddAlias(tt.shell, tt.alias, tt.command, tt.strategy)
				if (err != nil) != tt.wantErr {
					t.Errorf("AddAlias() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				// Check if alias was added
				info, _ := mockManager.GetInfo(tt.shell)
				content, err := os.ReadFile(info.ConfigFiles[0])
				if err != nil {
					t.Errorf("Failed to read config file: %v", err)
					return
				}

				expectedAlias := "alias " + tt.alias + "='" + tt.command + "'"
				if !strings.Contains(string(content), expectedAlias) {
					t.Errorf("Config file does not contain expected alias. Got: %s", string(content))
				}
			})
		}
	})

	// Test adding source commands
	t.Run("AddSource", func(t *testing.T) {
		tests := []struct {
			name     string
			shell    Shell
			file     string
			strategy DotfilesStrategy
			wantErr  bool
		}{
			{
				name:     "bash add source",
				shell:    Bash,
				file:     "~/.bash_aliases",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "zsh add source",
				shell:    Zsh,
				file:     "~/.zsh_aliases",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
			{
				name:     "fish add source",
				shell:    Fish,
				file:     "~/.config/fish/functions.fish",
				strategy: ReplaceExisting,
				wantErr:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := writer.AddSource(tt.shell, tt.file, tt.strategy)
				if (err != nil) != tt.wantErr {
					t.Errorf("AddSource() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.wantErr {
					return
				}

				// Check if source was added
				info, _ := mockManager.GetInfo(tt.shell)
				content, err := os.ReadFile(info.ConfigFiles[0])
				if err != nil {
					t.Errorf("Failed to read config file: %v", err)
					return
				}

				expectedSource := "source " + tt.file
				if !strings.Contains(string(content), expectedSource) {
					t.Errorf("Config file does not contain expected source. Got: %s", string(content))
				}
			})
		}
	})

	// Test checking if config exists
	t.Run("HasConfig", func(t *testing.T) {
		// First, write a config
		config := "export TEST_VAR=value"
		err := writer.WriteConfig(Bash, config, ReplaceExisting)
		if err != nil {
			t.Errorf("Failed to write initial config: %v", err)
			return
		}

		tests := []struct {
			name     string
			shell    Shell
			config   string
			expected bool
		}{
			{
				name:     "existing config",
				shell:    Bash,
				config:   config,
				expected: true,
			},
			{
				name:     "non-existing config",
				shell:    Bash,
				config:   "export NON_EXISTENT=value",
				expected: false,
			},
			{
				name:     "similar config",
				shell:    Bash,
				config:   "export TEST_VAR=other_value",
				expected: true, // Should detect similar config
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := writer.HasConfig(tt.shell, tt.config)
				if got != tt.expected {
					t.Errorf("HasConfig() = %v, want %v", got, tt.expected)
				}
			})
		}
	})
} 