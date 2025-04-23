package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

func TestGetShellTypeFromPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected interfaces.Shell
	}{
		{
			name:     "bash shell",
			path:     "/bin/bash",
			expected: interfaces.Bash,
		},
		{
			name:     "zsh shell",
			path:     "/usr/bin/zsh",
			expected: interfaces.Zsh,
		},
		{
			name:     "fish shell",
			path:     "/usr/local/bin/fish",
			expected: interfaces.Fish,
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
	// Save original home and restore after test
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	// Set up a temporary home directory
	tmpHome := t.TempDir()
	os.Setenv("HOME", tmpHome)

	tests := []struct {
		name     string
		shell    interfaces.Shell
		expected []string
	}{
		{
			name:  "bash config files",
			shell: interfaces.Bash,
			expected: []string{
				filepath.Join(tmpHome, ".bashrc"),
				filepath.Join(tmpHome, ".bash_profile"),
				filepath.Join(tmpHome, ".profile"),
			},
		},
		{
			name:  "zsh config files",
			shell: interfaces.Zsh,
			expected: []string{
				filepath.Join(tmpHome, ".zshrc"),
				filepath.Join(tmpHome, ".zprofile"),
				filepath.Join(tmpHome, ".zshenv"),
			},
		},
		{
			name:  "fish config files",
			shell: interfaces.Fish,
			expected: []string{
				filepath.Join(tmpHome, ".config/fish/config.fish"),
			},
		},
		{
			name:     "unknown shell",
			shell:    "unknown",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getShellConfigFiles(tt.shell)
			if len(got) != len(tt.expected) {
				t.Errorf("getShellConfigFiles() returned %d files, want %d", len(got), len(tt.expected))
				return
			}
			for i, file := range got {
				if file != tt.expected[i] {
					t.Errorf("getShellConfigFiles()[%d] = %v, want %v", i, file, tt.expected[i])
				}
			}
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
			wantShell:   string(interfaces.Bash),
			wantDefault: true,
		},
		{
			name:        "zsh shell",
			shellEnv:    zshPath,
			wantErr:     false,
			wantShell:   string(interfaces.Zsh),
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
			
			m := NewDefaultManager()
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
	m := NewDefaultManager()
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
			t.Error("Shell info missing path")
		}
		if info.Version == "" {
			t.Error("Shell info missing version")
		}
		if !info.IsAvailable {
			t.Error("Shell marked as unavailable but returned in available list")
		}
	}
} 