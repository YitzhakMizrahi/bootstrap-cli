package packages

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// mockExecCommand creates a mock command that returns the given output and error
func mockExecCommand(output string, err error) execCommand {
	return func(name string, arg ...string) *exec.Cmd {
		cmd := exec.Command("echo", output) // Use echo to simulate command output
		if err != nil {
			// Create a failing command
			cmd = exec.Command("false")
		}
		return cmd
	}
}

func TestNewPackageManager(t *testing.T) {
	tests := []struct {
		name    string
		system  string
		wantCmd string
		wantErr bool
	}{
		{
			name:    "ubuntu system",
			system:  "ubuntu",
			wantCmd: "apt-get",
			wantErr: false,
		},
		{
			name:    "debian system",
			system:  "debian",
			wantCmd: "apt-get",
			wantErr: false,
		},
		{
			name:    "fedora system",
			system:  "fedora",
			wantCmd: "dnf",
			wantErr: false,
		},
		{
			name:    "arch system",
			system:  "arch",
			wantCmd: "pacman",
			wantErr: false,
		},
		{
			name:    "unsupported system",
			system:  "windows",
			wantCmd: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.system)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pm.(*packageManager).cmd != tt.wantCmd {
				t.Errorf("NewPackageManager() cmd = %v, want %v", pm.(*packageManager).cmd, tt.wantCmd)
			}
		})
	}
}

func TestPackageManager_Install(t *testing.T) {
	tests := []struct {
		name        string
		system      string
		packageName string
		mockOutput  string
		mockErr     error
		wantErr     bool
	}{
		{
			name:        "successful install",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "Package installed successfully",
			mockErr:     nil,
			wantErr:     false,
		},
		{
			name:        "install failure",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "Installation failed",
			mockErr:     fmt.Errorf("command failed"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.system)
			if err != nil {
				t.Fatalf("Failed to create package manager: %v", err)
			}

			// Set mock command
			pm.(*packageManager).execCmd = mockExecCommand(tt.mockOutput, tt.mockErr)

			err = pm.Install(tt.packageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPackageManager_IsInstalled(t *testing.T) {
	tests := []struct {
		name        string
		system      string
		packageName string
		mockOutput  string
		mockErr     error
		want        bool
		wantErr     bool
	}{
		{
			name:        "package is installed",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "ii test-package 1.0.0",
			mockErr:     nil,
			want:        true,
			wantErr:     false,
		},
		{
			name:        "package is not installed",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "No packages found",
			mockErr:     nil,
			want:        false,
			wantErr:     false,
		},
		{
			name:        "command error",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "",
			mockErr:     fmt.Errorf("command failed"),
			want:        false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.system)
			if err != nil {
				t.Fatalf("Failed to create package manager: %v", err)
			}

			// Set mock command
			pm.(*packageManager).execCmd = mockExecCommand(tt.mockOutput, tt.mockErr)

			got, err := pm.IsInstalled(tt.packageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageManager_ListInstalled(t *testing.T) {
	tests := []struct {
		name       string
		system     string
		mockOutput string
		mockErr    error
		want       []string
		wantErr    bool
	}{
		{
			name:       "list installed packages",
			system:     "ubuntu",
			mockOutput: "pkg1\npkg2\npkg3",
			mockErr:    nil,
			want:       []string{"pkg1", "pkg2", "pkg3"},
			wantErr:    false,
		},
		{
			name:       "no packages installed",
			system:     "ubuntu",
			mockOutput: "",
			mockErr:    nil,
			want:       nil,
			wantErr:    false,
		},
		{
			name:       "command error",
			system:     "ubuntu",
			mockOutput: "",
			mockErr:    fmt.Errorf("command failed"),
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.system)
			if err != nil {
				t.Fatalf("Failed to create package manager: %v", err)
			}

			// Set mock command
			pm.(*packageManager).execCmd = mockExecCommand(tt.mockOutput, tt.mockErr)

			got, err := pm.ListInstalled()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListInstalled() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !stringSliceEqual(got, tt.want) {
				t.Errorf("ListInstalled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageManager_GetVersion(t *testing.T) {
	tests := []struct {
		name        string
		system      string
		packageName string
		mockOutput  string
		mockErr     error
		want        string
		wantErr     bool
	}{
		{
			name:        "get package version",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "Version: 1.0.0",
			mockErr:     nil,
			want:        "1.0.0",
			wantErr:     false,
		},
		{
			name:        "package not found",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "No version found",
			mockErr:     nil,
			want:        "",
			wantErr:     true,
		},
		{
			name:        "command error",
			system:      "ubuntu",
			packageName: "test-package",
			mockOutput:  "",
			mockErr:     fmt.Errorf("command failed"),
			want:        "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.system)
			if err != nil {
				t.Fatalf("Failed to create package manager: %v", err)
			}

			// Set mock command
			pm.(*packageManager).execCmd = mockExecCommand(tt.mockOutput, tt.mockErr)

			got, err := pm.GetVersion(tt.packageName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackageManagerLogging(t *testing.T) {
	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	
	// Create a logger that writes to our buffer
	customLogger := log.New(log.DebugLevel)
	customLogger.SetOutput(&logBuffer)
	
	// Create a package manager with our custom logger
	pm := &packageManager{
		system:  "ubuntu",
		cmd:    "apt-get",
		logger: customLogger,
		execCmd: mockExecCommand("Success", nil),
	}
	
	// Test Install logging
	t.Run("Install Logging", func(t *testing.T) {
		logBuffer.Reset()
		_ = pm.Install("test-package")
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Installing package: test-package") {
			t.Error("Log missing installation start message")
		}
	})
	
	// Test Uninstall logging
	t.Run("Uninstall Logging", func(t *testing.T) {
		logBuffer.Reset()
		_ = pm.Uninstall("test-package")
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Uninstalling package: test-package") {
			t.Error("Log missing uninstallation start message")
		}
	})
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
} 