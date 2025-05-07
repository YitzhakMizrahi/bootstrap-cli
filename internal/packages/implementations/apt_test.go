package implementations

import (
	"os/exec"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestNewAptPackageManager(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}
	if pm == nil {
		t.Fatal("NewAptPackageManager() returned nil")
	}
}

func TestAptPackageManager_Install(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	// Test installing a package that should exist
	err = pm.Install("curl")
	if err != nil {
		t.Errorf("Install() error = %v", err)
	}

	// Test installing a non-existent package
	err = pm.Install("non-existent-package-123456")
	if err == nil {
		t.Error("Install() expected error for non-existent package, got nil")
	}
}

func TestAptPackageManager_IsInstalled(t *testing.T) {
	pm, err := NewAptPackageManager()
	assert.NoError(t, err)

	// Check for a known installed package (e.g., apt itself or bash)
	installed, err := pm.IsInstalled("bash") // Handle error
	assert.NoError(t, err) // Expect no error running the check
	assert.True(t, installed, "Expected 'bash' to be installed")

	// Check for a non-existent package
	installed, err = pm.IsInstalled("nonexistent-package-kjshdfg") // Handle error
	assert.NoError(t, err) // Expect no error running the check, just false result
	assert.False(t, installed, "Expected 'nonexistent-package-kjshdfg' not to be installed")
}

func TestAptPackageManager_Update(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	err = pm.Update()
	if err != nil {
		t.Errorf("Update() error = %v", err)
	}
}

func TestAptPackageManager_Upgrade(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	err = pm.Upgrade()
	if err != nil {
		t.Errorf("Upgrade() error = %v", err)
	}
}

func TestAptPackageManager_GetName(t *testing.T) {
	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}
	assert.Equal(t, string(interfaces.APT), pm.GetName())
}

func TestAptPackageManager_Remove(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	// Test removing a package that should exist
	err = pm.Uninstall("curl")
	if err != nil {
		t.Errorf("Uninstall() error = %v", err)
	}

	// Test removing a non-existent package
	err = pm.Uninstall("non-existent-package-123456")
	if err == nil {
		t.Error("Uninstall() expected error for non-existent package, got nil")
	}
}

func TestAptPackageManager_GetVersion(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	// Test with a package that should be installed
	version, err := pm.GetVersion("apt")
	if err != nil {
		t.Errorf("GetVersion() error = %v", err)
	}
	if version == "" {
		t.Error("GetVersion() returned empty version for apt package")
	}

	// Test with a non-existent package
	version, err = pm.GetVersion("non-existent-package-123456")
	if err == nil {
		t.Error("GetVersion() expected error for non-existent package, got nil")
	}
	if version != "" {
		t.Error("GetVersion() returned non-empty version for non-existent package")
	}
}

func TestAptPackageManager_ListInstalled(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	packages, err := pm.ListInstalled()
	if err != nil {
		t.Errorf("ListInstalled() error = %v", err)
	}
	if len(packages) == 0 {
		t.Error("ListInstalled() returned empty list")
	}
}

func TestAptPackageManager_IsAvailable(t *testing.T) {
	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}
	assert.True(t, pm.IsAvailable())
}

func TestAptPackageManager_InstallMultiple(t *testing.T) {
	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}
	err = pm.Install("package1")
	assert.NoError(t, err)
	err = pm.Install("package2")
	assert.NoError(t, err)
	err = pm.Install("package3")
	assert.NoError(t, err)
}

func TestAptPackageManager_InstallEmpty(t *testing.T) {
	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}
	err = pm.Install("")
	assert.Error(t, err)
}

func TestAptInstall(t *testing.T) {
	// Skip if not on a system with apt-get
	if _, err := exec.LookPath("apt-get"); err != nil {
		t.Skip("apt-get not available, skipping test")
	}

	pm, err := NewAptPackageManager()
	if err != nil {
		t.Fatalf("NewAptPackageManager() error = %v", err)
	}

	// Test installing a package
	err = pm.Install("test-package")
	assert.NoError(t, err)
} 