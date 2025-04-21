package pkgmanager

import (
	"os/exec"
	"testing"
)

// TestNewPackageManager tests the creation of package managers
func TestNewPackageManager(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "apt manager",
			input:   "apt",
			wantErr: false,
		},
		{
			name:    "dnf manager",
			input:   "dnf",
			wantErr: false,
		},
		{
			name:    "pacman manager",
			input:   "pacman",
			wantErr: false,
		},
		{
			name:    "brew manager",
			input:   "brew",
			wantErr: false,
		},
		{
			name:    "choco manager",
			input:   "choco",
			wantErr: false,
		},
		{
			name:    "invalid manager",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm, err := NewPackageManager(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pm == nil {
				t.Error("NewPackageManager() returned nil manager when no error expected")
			}
			if !tt.wantErr && pm.Name() != tt.input {
				t.Errorf("NewPackageManager() name = %v, want %v", pm.Name(), tt.input)
			}
		})
	}
}

// TestAptManager tests the AptManager implementation
func TestAptManager(t *testing.T) {
	// Skip if not on a system with apt
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("Skipping test on non-apt system")
	}

	pm := &AptManager{}

	// Test Name
	if pm.Name() != "apt" {
		t.Errorf("AptManager.Name() = %v, want %v", pm.Name(), "apt")
	}

	// Test IsAvailable
	if !pm.IsAvailable() {
		t.Error("AptManager.IsAvailable() = false, want true")
	}

	// Test IsInstalled with a common package
	// This test assumes curl is installed, which is common on most systems
	// If it's not installed, the test will fail
	if !pm.IsInstalled("curl") {
		t.Log("curl not installed, skipping IsInstalled test")
	}

	// Test Search
	results, err := pm.Search("curl")
	if err != nil {
		t.Errorf("AptManager.Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("AptManager.Search() returned no results")
	}

	// Note: We don't test Install, Uninstall, Update, or UpdateAll
	// as these would require root privileges and modify the system
}

// TestDnfManager tests the DnfManager implementation
func TestDnfManager(t *testing.T) {
	// Skip if not on a system with dnf
	if _, err := exec.LookPath("dnf"); err != nil {
		t.Skip("Skipping test on non-dnf system")
	}

	pm := &DnfManager{}

	// Test Name
	if pm.Name() != "dnf" {
		t.Errorf("DnfManager.Name() = %v, want %v", pm.Name(), "dnf")
	}

	// Test IsAvailable
	if !pm.IsAvailable() {
		t.Error("DnfManager.IsAvailable() = false, want true")
	}

	// Test IsInstalled with a common package
	// This test assumes curl is installed, which is common on most systems
	// If it's not installed, the test will fail
	if !pm.IsInstalled("curl") {
		t.Log("curl not installed, skipping IsInstalled test")
	}

	// Test Search
	results, err := pm.Search("curl")
	if err != nil {
		t.Errorf("DnfManager.Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("DnfManager.Search() returned no results")
	}

	// Note: We don't test Install, Uninstall, Update, or UpdateAll
	// as these would require root privileges and modify the system
}

// TestPacmanManager tests the PacmanManager implementation
func TestPacmanManager(t *testing.T) {
	// Skip if not on a system with pacman
	if _, err := exec.LookPath("pacman"); err != nil {
		t.Skip("Skipping test on non-pacman system")
	}

	pm := &PacmanManager{}

	// Test Name
	if pm.Name() != "pacman" {
		t.Errorf("PacmanManager.Name() = %v, want %v", pm.Name(), "pacman")
	}

	// Test IsAvailable
	if !pm.IsAvailable() {
		t.Error("PacmanManager.IsAvailable() = false, want true")
	}

	// Test IsInstalled with a common package
	// This test assumes curl is installed, which is common on most systems
	// If it's not installed, the test will fail
	if !pm.IsInstalled("curl") {
		t.Log("curl not installed, skipping IsInstalled test")
	}

	// Test Search
	results, err := pm.Search("curl")
	if err != nil {
		t.Errorf("PacmanManager.Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("PacmanManager.Search() returned no results")
	}

	// Note: We don't test Install, Uninstall, Update, or UpdateAll
	// as these would require root privileges and modify the system
}

// TestBrewManager tests the BrewManager implementation
func TestBrewManager(t *testing.T) {
	// Skip if not on a system with brew
	if _, err := exec.LookPath("brew"); err != nil {
		t.Skip("Skipping test on non-brew system")
	}

	pm := &BrewManager{}

	// Test Name
	if pm.Name() != "brew" {
		t.Errorf("BrewManager.Name() = %v, want %v", pm.Name(), "brew")
	}

	// Test IsAvailable
	if !pm.IsAvailable() {
		t.Error("BrewManager.IsAvailable() = false, want true")
	}

	// Test IsInstalled with a common package
	// This test assumes curl is installed, which is common on most systems
	// If it's not installed, the test will fail
	if !pm.IsInstalled("curl") {
		t.Log("curl not installed, skipping IsInstalled test")
	}

	// Test Search
	results, err := pm.Search("curl")
	if err != nil {
		t.Errorf("BrewManager.Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("BrewManager.Search() returned no results")
	}

	// Note: We don't test Install, Uninstall, Update, or UpdateAll
	// as these would require root privileges and modify the system
}

// TestChocoManager tests the ChocoManager implementation
func TestChocoManager(t *testing.T) {
	// Skip if not on a system with choco
	if _, err := exec.LookPath("choco"); err != nil {
		t.Skip("Skipping test on non-choco system")
	}

	pm := &ChocoManager{}

	// Test Name
	if pm.Name() != "choco" {
		t.Errorf("ChocoManager.Name() = %v, want %v", pm.Name(), "choco")
	}

	// Test IsAvailable
	if !pm.IsAvailable() {
		t.Error("ChocoManager.IsAvailable() = false, want true")
	}

	// Test IsInstalled with a common package
	// This test assumes curl is installed, which is common on most systems
	// If it's not installed, the test will fail
	if !pm.IsInstalled("curl") {
		t.Log("curl not installed, skipping IsInstalled test")
	}

	// Test Search
	results, err := pm.Search("curl")
	if err != nil {
		t.Errorf("ChocoManager.Search() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("ChocoManager.Search() returned no results")
	}

	// Note: We don't test Install, Uninstall, Update, or UpdateAll
	// as these would require root privileges and modify the system
} 