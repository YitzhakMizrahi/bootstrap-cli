// Package shell provides shell configuration and management functionality for the bootstrap-cli,
// including shell detection, configuration writing, and environment setup.
package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// mockPackageManager implements interfaces.PackageManager for testing
type mockPackageManager struct{}

func (m *mockPackageManager) GetName() string                     { return "mock" }
func (m *mockPackageManager) IsAvailable() bool                   { return true }
func (m *mockPackageManager) Install(_ string) error             { return nil }
func (m *mockPackageManager) Update() error                      { return nil }
func (m *mockPackageManager) Upgrade() error                     { return nil }
func (m *mockPackageManager) IsInstalled(_ string) bool          { return true }
func (m *mockPackageManager) Remove(_ string) error              { return nil }
func (m *mockPackageManager) GetVersion(_ string) (string, error) { return "", nil }
func (m *mockPackageManager) ListInstalled() ([]string, error)   { return nil, nil }
func (m *mockPackageManager) SetupSpecialPackage(_ string) error { return nil }

// testConfigWriter creates a DefaultConfigWriter for testing
func testConfigWriter(t *testing.T, shell interfaces.ShellType) (*DefaultConfigWriter, string, func()) {
	// Create a temporary directory for test config files
	tmpDir, err := os.MkdirTemp("", "shell-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create a cleanup function
	cleanup := func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to cleanup temp dir: %v", err)
		}
	}

	writer := &DefaultConfigWriter{
		logger: log.New(log.InfoLevel),
		shell:  shell,
		pm:     &mockPackageManager{},
	}

	// Override os.UserHomeDir during test
	os.Setenv("HOME", tmpDir)

	return writer, tmpDir, cleanup
}

func TestWriteConfig(t *testing.T) {
	tests := []struct {
		name     string
		shell    interfaces.ShellType
		configs  []string
		strategy interfaces.DotfilesStrategy
		existing string // existing content in config file
		want     string // expected content after write
		wantErr  bool
	}{
		{
			name:     "write to empty file",
			shell:    interfaces.BashShell,
			configs:  []string{"export PATH=/usr/local/bin:$PATH"},
			strategy: interfaces.MergeWithExisting,
			want:     "export PATH=/usr/local/bin:$PATH\n",
		},
		{
			name:     "merge with existing",
			shell:    interfaces.ZshShell,
			configs:  []string{"export GOPATH=$HOME/go"},
			strategy: interfaces.MergeWithExisting,
			existing: "export PATH=/usr/local/bin:$PATH\n",
			want:     "export PATH=/usr/local/bin:$PATH\nexport GOPATH=$HOME/go\n",
		},
		{
			name:     "skip if exists",
			shell:    interfaces.BashShell,
			configs:  []string{"export PATH=/usr/local/bin:$PATH"},
			strategy: interfaces.SkipIfExists,
			existing: "export PATH=/usr/local/bin:$PATH\n",
			want:     "export PATH=/usr/local/bin:$PATH\n",
		},
		{
			name:     "replace existing",
			shell:    interfaces.FishShell,
			configs:  []string{"set -gx PATH /usr/local/bin $PATH"},
			strategy: interfaces.ReplaceExisting,
			existing: "set -gx PATH /bin $PATH\n",
			want:     "set -gx PATH /usr/local/bin $PATH\n",
		},
		{
			name:     "multiple configs",
			shell:    interfaces.BashShell,
			configs:  []string{"export A=1", "export B=2"},
			strategy: interfaces.MergeWithExisting,
			want:     "export A=1\nexport B=2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, _, cleanup := testConfigWriter(t, tt.shell)
			defer cleanup()

			// Create config file with existing content if any
			configFile := writer.getConfigFile()
			if tt.existing != "" {
				err := os.MkdirAll(filepath.Dir(configFile), 0755)
				if err != nil {
					t.Fatalf("Failed to create config dir: %v", err)
				}
				err = os.WriteFile(configFile, []byte(tt.existing), 0644)
				if err != nil {
					t.Fatalf("Failed to write existing config: %v", err)
				}
			}

			err := writer.WriteConfig(tt.configs, tt.strategy)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := os.ReadFile(configFile)
				if err != nil {
					t.Fatalf("Failed to read config file: %v", err)
				}
				if string(got) != tt.want {
					t.Errorf("WriteConfig() got = %q, want %q", string(got), tt.want)
				}
			}
		})
	}
}

func TestAddToPath(t *testing.T) {
	writer, _, cleanup := testConfigWriter(t, interfaces.BashShell)
	defer cleanup()

	err := writer.AddToPath("/test/bin")
	if err != nil {
		t.Errorf("AddToPath() error = %v", err)
		return
	}

	configFile := writer.getConfigFile()
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	want := "export PATH=/test/bin:$PATH\n"
	if string(content) != want {
		t.Errorf("AddToPath() got = %q, want %q", string(content), want)
	}
}

func TestSetEnvVar(t *testing.T) {
	writer, _, cleanup := testConfigWriter(t, interfaces.BashShell)
	defer cleanup()

	err := writer.SetEnvVar("TESTVAR", "value")
	if err != nil {
		t.Errorf("SetEnvVar() error = %v", err)
		return
	}

	configFile := writer.getConfigFile()
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	want := "export TESTVAR=value\n"
	if string(content) != want {
		t.Errorf("SetEnvVar() got = %q, want %q", string(content), want)
	}
}

func TestAddAlias(t *testing.T) {
	writer, _, cleanup := testConfigWriter(t, interfaces.BashShell)
	defer cleanup()

	err := writer.AddAlias("ll", "ls -la")
	if err != nil {
		t.Errorf("AddAlias() error = %v", err)
		return
	}

	configFile := writer.getConfigFile()
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	want := "alias ll='ls -la'\n"
	if string(content) != want {
		t.Errorf("AddAlias() got = %q, want %q", string(content), want)
	}
}

func TestHasConfig(t *testing.T) {
	tests := []struct {
		name     string
		shell    interfaces.ShellType
		existing string
		config   string
		want     bool
	}{
		{
			name:     "config exists",
			shell:    interfaces.BashShell,
			existing: "export PATH=/usr/local/bin:$PATH\n",
			config:   "export PATH=/usr/local/bin:$PATH",
			want:     true,
		},
		{
			name:     "config does not exist",
			shell:    interfaces.ZshShell,
			existing: "export PATH=/usr/local/bin:$PATH\n",
			config:   "export GOPATH=$HOME/go",
			want:     false,
		},
		{
			name:     "empty file",
			shell:    interfaces.FishShell,
			config:   "set -gx PATH /usr/local/bin $PATH",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, _, cleanup := testConfigWriter(t, tt.shell)
			defer cleanup()

			if tt.existing != "" {
				configFile := writer.getConfigFile()
				err := os.MkdirAll(filepath.Dir(configFile), 0755)
				if err != nil {
					t.Fatalf("Failed to create config dir: %v", err)
				}
				err = os.WriteFile(configFile, []byte(tt.existing), 0644)
				if err != nil {
					t.Fatalf("Failed to write existing config: %v", err)
				}
			}

			if got := writer.HasConfig(tt.config); got != tt.want {
				t.Errorf("HasConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigFile(t *testing.T) {
	tests := []struct {
		name     string
		shell    interfaces.ShellType
		wantPath string
	}{
		{
			name:     "bash config",
			shell:    interfaces.BashShell,
			wantPath: ".bashrc",
		},
		{
			name:     "zsh config",
			shell:    interfaces.ZshShell,
			wantPath: ".zshrc",
		},
		{
			name:     "fish config",
			shell:    interfaces.FishShell,
			wantPath: filepath.Join(".config", "fish", "config.fish"),
		},
		{
			name:     "unknown shell",
			shell:    "unknown",
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, tmpDir, cleanup := testConfigWriter(t, tt.shell)
			defer cleanup()

			got := writer.getConfigFile()
			expectedPath := filepath.Join(tmpDir, tt.wantPath)

			if tt.wantPath == "" {
				if got != "" {
					t.Errorf("getConfigFile() = %q, want empty string", got)
				}
			} else if !strings.HasSuffix(got, expectedPath) {
				t.Errorf("getConfigFile() = %q, want path ending with %q", got, expectedPath)
			}
		})
	}
} 