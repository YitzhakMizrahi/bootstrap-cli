package factory

import (
	"fmt"
	"testing"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

func TestNewPackageManagerFactory(t *testing.T) {
	tests := []struct {
		name     string
		pmType   interfaces.PackageManagerType
		wantType interfaces.PackageManager
		wantErr  bool
	}{
		{
			name:     "test apt",
			pmType:   interfaces.APT,
			wantType: NewMockPackageManager(0, "apt"),
			wantErr:  false,
		},
		{
			name:     "test dnf",
			pmType:   interfaces.DNF,
			wantType: NewMockPackageManager(0, "dnf"),
			wantErr:  false,
		},
		{
			name:     "test pacman",
			pmType:   interfaces.Pacman,
			wantType: NewMockPackageManager(0, "pacman"),
			wantErr:  false,
		},
		{
			name:     "test brew",
			pmType:   interfaces.Homebrew,
			wantType: NewMockPackageManager(0, "brew"),
			wantErr:  false,
		},
		{
			name:     "test unknown",
			pmType:   interfaces.PackageManagerType("unknown"),
			wantType: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewPackageManagerFactory()
			got, err := f.GetPackageManager()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPackageManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.GetName() != tt.wantType.GetName() {
				t.Errorf("GetPackageManager() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestSetRetryConfig(t *testing.T) {
	factory := NewPackageManagerFactory()
	factory.SetRetryConfig(5, 10*time.Second)
	if factory.maxRetries != 5 {
		t.Errorf("expected maxRetries to be 5, got %d", factory.maxRetries)
	}
	if factory.retryDelay != 10*time.Second {
		t.Errorf("expected retryDelay to be 10s, got %v", factory.retryDelay)
	}
}

func TestGetPackageManager(t *testing.T) {
	factory := NewPackageManagerFactory()
	pm, err := factory.GetPackageManager()
	if err != nil {
		t.Fatalf("GetPackageManager() error = %v", err)
	}
	if pm == nil {
		t.Fatal("GetPackageManager() returned nil package manager")
	}

	// Test that we got a retryPackageManager
	if _, ok := pm.(*retryPackageManager); !ok {
		t.Error("GetPackageManager() did not return a retryPackageManager")
	}
}

func TestRetryPackageManager(t *testing.T) {
	// Create a mock package manager that fails twice then succeeds
	mockPM := &mockPackageManager{
		failCount:      2,
		installCount:   0,
		uninstallCount: 0,
	}

	// Create a retry package manager with the mock
	retryPM := &retryPackageManager{
		PackageManager: mockPM,
		maxRetries:    3,
		retryDelay:    10 * time.Millisecond, // Short delay for testing
	}

	// Test Install with retry
	err := retryPM.Install("test-package")
	if err != nil {
		t.Errorf("Install() error = %v", err)
	}
	if mockPM.installCount != 3 {
		t.Errorf("expected installCount to be 3, got %d", mockPM.installCount)
	}

	// Test Remove with retry
	err = retryPM.Remove("test-package")
	if err != nil {
		t.Errorf("Remove() error = %v", err)
	}
	if mockPM.uninstallCount != 3 {
		t.Errorf("expected uninstallCount to be 3, got %d", mockPM.uninstallCount)
	}

	// Test that it fails after max retries
	mockPM.failCount = 10 // More than maxRetries
	err = retryPM.Install("test-package")
	if err == nil {
		t.Error("expected Install() to fail after max retries")
	}
}

// mockPackageManager is a mock implementation of the PackageManager interface
type mockPackageManager struct {
	failCount      int
	installCount   int
	uninstallCount int
}

func (m *mockPackageManager) Install(packageName string) error {
	return nil
}

func (m *mockPackageManager) Remove(pkg string) error {
	m.uninstallCount++
	if m.uninstallCount <= m.failCount {
		return fmt.Errorf("mock remove error")
	}
	return nil
}

func (m *mockPackageManager) Name() string {
	return "mock"
}

func (m *mockPackageManager) IsAvailable() bool {
	return true
}

func (m *mockPackageManager) Update() error {
	return nil
}

func (m *mockPackageManager) IsInstalled(pkg string) bool {
	return false
}

func (m *mockPackageManager) GetVersion(packageName string) (string, error) {
	return "1.0.0", nil
}

func (m *mockPackageManager) ListInstalled() ([]string, error) {
	return []string{}, nil
}

func (m *mockPackageManager) GetName() string {
	return "mock"
}

func (m *mockPackageManager) Upgrade() error {
	return nil
}

func (m *mockPackageManager) SetupSpecialPackage(pkg string) error {
	return nil
} 