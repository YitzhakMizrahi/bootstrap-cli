package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestGetShellTypeFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected interfaces.ShellType
	}{
		{
			name:     "bash shell",
			path:     "/bin/bash",
			expected: interfaces.BashShell,
		},
		{
			name:     "zsh shell",
			path:     "/usr/bin/zsh",
			expected: interfaces.ZshShell,
		},
		{
			name:     "fish shell",
			path:     "/usr/local/bin/fish",
			expected: interfaces.FishShell,
		},
		{
			name:     "unknown shell",
			path:     "/bin/unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getShellTypeFromPath(tt.path)
			if got != tt.expected {
				t.Errorf("getShellTypeFromPath() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetShellConfigFiles(t *testing.T) {
	tests := []struct {
		name     string
		shell    interfaces.ShellType
		want     []string
		wantErr  bool
	}{
		{
			name: "bash config files",
			shell: interfaces.BashShell,
			want: []string{
				filepath.Join(os.Getenv("HOME"), ".bashrc"),
				filepath.Join(os.Getenv("HOME"), ".bash_profile"),
			},
			wantErr: false,
		},
		{
			name: "zsh config files",
			shell: interfaces.ZshShell,
			want: []string{
				filepath.Join(os.Getenv("HOME"), ".zshrc"),
			},
			wantErr: false,
		},
		{
			name: "fish config files",
			shell: interfaces.FishShell,
			want: []string{
				filepath.Join(os.Getenv("HOME"), ".config/fish/config.fish"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getShellConfigFiles(tt.shell)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetectCurrent(t *testing.T) {
	// Save original SHELL and restore after test
	origShell := os.Getenv("SHELL")
	defer os.Setenv("SHELL", origShell)

	// Get actual shell paths if available
	bashPath, _ := exec.LookPath("bash")
	if bashPath == "" {
		bashPath = "/bin/bash" // fallback
	}
	zshPath, _ := exec.LookPath("zsh")
	if zshPath == "" {
		zshPath = "/usr/bin/zsh" // fallback
	}

	tests := []struct {
		name        string
		shellEnv    string
		wantErr     bool
		wantShell   string
		wantDefault bool
	}{
		{
			name:        "bash shell",
			shellEnv:    bashPath,
			wantErr:     false,
			wantShell:   string(interfaces.BashShell),
			wantDefault: true,
		},
		{
			name:        "zsh shell",
			shellEnv:    zshPath,
			wantErr:     false,
			wantShell:   string(interfaces.ZshShell),
			wantDefault: true,
		},
		{
			name:        "no shell env",
			shellEnv:    "",
			wantErr:     true,
			wantShell:   "",
			wantDefault: false,
		},
		{
			name:        "unknown shell",
			shellEnv:    "/bin/unknown",
			wantErr:     true,
			wantShell:   "",
			wantDefault: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test if shell is not available
			if tt.shellEnv != "" && tt.shellEnv != "/bin/unknown" {
				if _, err := os.Stat(tt.shellEnv); os.IsNotExist(err) {
					t.Skipf("Shell %s not available, skipping test", tt.shellEnv)
					return
				}
			}

			os.Setenv("SHELL", tt.shellEnv)
			
			m, err := NewManager()
			if err != nil {
				t.Fatalf("Failed to create shell manager: %v", err)
			}
			info, err := m.DetectCurrent()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr {
				return
			}
			
			if info.Type != tt.wantShell {
				t.Errorf("DetectCurrent() shell = %v, want %v", info.Type, tt.wantShell)
			}
			
			if info.IsDefault != tt.wantDefault {
				t.Errorf("DetectCurrent() isDefault = %v, want %v", info.IsDefault, tt.wantDefault)
			}
		})
	}
}

func TestListAvailable(t *testing.T) {
	m, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create shell manager: %v", err)
	}
	shells, err := m.ListAvailable()
	if err != nil {
		t.Errorf("ListAvailable() error = %v", err)
		return
	}

	// At least one shell should be available (usually bash)
	if len(shells) == 0 {
		t.Error("ListAvailable() returned no shells, expected at least one")
	}

	// Verify shell info
	for _, info := range shells {
		if info.Path == "" {
			t.Errorf("Shell %s has empty path", info.Type)
		}
		if info.Type == "" {
			t.Error("Shell has empty type")
		}
		if !info.IsAvailable {
			t.Errorf("Shell %s marked as not available", info.Type)
		}
	}
}

func TestDetectCurrentShell(t *testing.T) {
	tests := []struct {
		name     string
		shellEnv string
		want     interfaces.ShellType
	}{
		{
			name:     "bash shell",
			shellEnv: "/bin/bash",
			want:     interfaces.BashShell,
		},
		{
			name:     "zsh shell",
			shellEnv: "/bin/zsh",
			want:     interfaces.ZshShell,
		},
		{
			name:     "fish shell",
			shellEnv: "/usr/bin/fish",
			want:     interfaces.FishShell,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SHELL", tt.shellEnv)
			got := detectCurrentShell()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsShellInstalled(t *testing.T) {
	tests := []struct {
		name     string
		shell    interfaces.ShellType
		shellEnv string
		want     bool
	}{
		{
			name:     "bash installed",
			shell:    interfaces.BashShell,
			shellEnv: "/bin/bash",
			want:     true,
		},
		{
			name:     "zsh installed",
			shell:    interfaces.ZshShell,
			shellEnv: "/bin/zsh",
			want:     true,
		},
		{
			name:     "fish installed",
			shell:    interfaces.FishShell,
			shellEnv: "/usr/bin/fish",
			want:     true,
		},
		{
			name:     "shell not installed",
			shell:    interfaces.BashShell,
			shellEnv: "/bin/nonexistent",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SHELL", tt.shellEnv)
			got := isShellInstalled(tt.shell)
			assert.Equal(t, tt.want, got)
		})
	}
} 