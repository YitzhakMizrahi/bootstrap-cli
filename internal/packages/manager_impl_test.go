package packages

import (
	"os/exec"
	"testing"
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