package packages

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

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
	// Skip if not running in test environment
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	pm, err := NewPackageManager("ubuntu")
	if err != nil {
		t.Fatalf("Failed to create package manager: %v", err)
	}

	// Test installing a package
	err = pm.Install("test-package")
	if err != nil {
		// We expect an error since we're not actually installing packages
		// Just verify it's the expected type of error
		if _, ok := err.(*exec.ExitError); !ok {
			t.Errorf("Install() error = %v, want *exec.ExitError", err)
		}
	}
}

func TestPackageManager_IsInstalled(t *testing.T) {
	// Skip if not running in test environment
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	pm, err := NewPackageManager("ubuntu")
	if err != nil {
		t.Fatalf("Failed to create package manager: %v", err)
	}

	// Test checking if a package is installed
	installed, err := pm.IsInstalled("test-package")
	if err != nil {
		// We expect an error since we're not actually checking packages
		// Just verify it's the expected type of error
		if _, ok := err.(*exec.ExitError); !ok {
			t.Errorf("IsInstalled() error = %v, want *exec.ExitError", err)
		}
	}
	if installed {
		t.Error("IsInstalled() returned true for non-existent package")
	}
}

func TestPackageManager_ListInstalled(t *testing.T) {
	// Skip if not running in test environment
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	pm, err := NewPackageManager("ubuntu")
	if err != nil {
		t.Fatalf("Failed to create package manager: %v", err)
	}

	// Test listing installed packages
	packages, err := pm.ListInstalled()
	if err != nil {
		// We expect an error since we're not actually listing packages
		// Just verify it's the expected type of error
		if _, ok := err.(*exec.ExitError); !ok {
			t.Errorf("ListInstalled() error = %v, want *exec.ExitError", err)
		}
	}
	if packages != nil {
		t.Error("ListInstalled() returned non-nil packages list")
	}
}

func TestPackageManager_GetVersion(t *testing.T) {
	// Skip if not running in test environment
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	pm, err := NewPackageManager("ubuntu")
	if err != nil {
		t.Fatalf("Failed to create package manager: %v", err)
	}

	// Test getting package version
	version, err := pm.GetVersion("test-package")
	if err != nil {
		// We expect an error since we're not actually getting versions
		// Just verify it's the expected type of error
		if _, ok := err.(*exec.ExitError); !ok {
			t.Errorf("GetVersion() error = %v, want *exec.ExitError", err)
		}
	}
	if version != "" {
		t.Error("GetVersion() returned non-empty version for non-existent package")
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
		system: "ubuntu",
		cmd:    "apt-get",
		logger: customLogger,
	}
	
	// Test Install logging
	t.Run("Install Logging", func(t *testing.T) {
		logBuffer.Reset()
		_ = pm.Install("test-package")
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Installing package: test-package") {
			t.Error("Log missing installation start message")
		}
		if !strings.Contains(logOutput, "Failed to install package test-package") {
			t.Error("Log missing installation failure message")
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
		if !strings.Contains(logOutput, "Failed to uninstall package test-package") {
			t.Error("Log missing uninstallation failure message")
		}
	})
	
	// Test IsInstalled logging
	t.Run("IsInstalled Logging", func(t *testing.T) {
		logBuffer.Reset()
		_, _ = pm.IsInstalled("test-package")
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Checking if package is installed: test-package") {
			t.Error("Log missing check installation message")
		}
		if !strings.Contains(logOutput, "Failed to check if package test-package is installed") {
			t.Error("Log missing check installation failure message")
		}
	})
	
	// Test ListInstalled logging
	t.Run("ListInstalled Logging", func(t *testing.T) {
		logBuffer.Reset()
		_, _ = pm.ListInstalled()
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Listing installed packages") {
			t.Error("Log missing list packages message")
		}
		if !strings.Contains(logOutput, "Failed to list installed packages") {
			t.Error("Log missing list packages failure message")
		}
	})
	
	// Test GetVersion logging
	t.Run("GetVersion Logging", func(t *testing.T) {
		logBuffer.Reset()
		_, _ = pm.GetVersion("test-package")
		
		logOutput := logBuffer.String()
		if !strings.Contains(logOutput, "Getting version for package: test-package") {
			t.Error("Log missing get version message")
		}
		if !strings.Contains(logOutput, "Failed to get version for package test-package") {
			t.Error("Log missing get version failure message")
		}
	})
} 