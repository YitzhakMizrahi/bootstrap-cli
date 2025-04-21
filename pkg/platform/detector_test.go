package platform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectLinuxDistribution(t *testing.T) {
	detector := NewDetector()
	
	// Create a temporary os-release file
	tmpDir := t.TempDir()
	osReleasePath := filepath.Join(tmpDir, "os-release")
	osReleaseContent := `ID=ubuntu
VERSION_ID="22.04"
`
	err := os.WriteFile(osReleasePath, []byte(osReleaseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary os-release file: %v", err)
	}
	
	// Test with the temporary file
	distro := detector.detectLinuxDistribution()
	if distro != "ubuntu" {
		t.Errorf("Expected distribution 'ubuntu', got '%s'", distro)
	}
}

func TestIsCommandAvailable(t *testing.T) {
	detector := NewDetector()
	
	// Test with a command that should exist
	if !detector.IsCommandAvailable("ls") {
		t.Error("Expected 'ls' command to be available")
	}
	
	// Test with a non-existent command
	if detector.IsCommandAvailable("nonexistentcommand123") {
		t.Error("Expected non-existent command to not be available")
	}
}

func TestGetPrimaryPackageManager(t *testing.T) {
	detector := NewDetector()
	
	tests := []struct {
		name     string
		info     Info
		expected PackageManager
		wantErr  bool
	}{
		{
			name: "macOS with Homebrew",
			info: Info{
				OS:              MacOS,
				PackageManagers: []PackageManager{Homebrew},
			},
			expected: Homebrew,
			wantErr:  false,
		},
		{
			name: "Linux with multiple package managers",
			info: Info{
				OS:              Linux,
				PackageManagers: []PackageManager{Apt, Homebrew},
			},
			expected: Homebrew,
			wantErr:  false,
		},
		{
			name: "No package managers",
			info: Info{
				OS:              Linux,
				PackageManagers: []PackageManager{},
			},
			expected: "",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detector.GetPrimaryPackageManager(tt.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrimaryPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("GetPrimaryPackageManager() = %v, want %v", got, tt.expected)
			}
		})
	}
} 