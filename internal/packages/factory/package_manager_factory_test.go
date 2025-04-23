package factory

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPackageManagerFactory(t *testing.T) {
	factory := NewPackageManagerFactory()
	if factory == nil {
		t.Fatal("NewPackageManagerFactory() returned nil")
	}
	if factory.maxRetries != 3 {
		t.Errorf("expected maxRetries to be 3, got %d", factory.maxRetries)
	}
	if factory.retryDelay != 5*time.Second {
		t.Errorf("expected retryDelay to be 5s, got %v", factory.retryDelay)
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

func (m *mockPackageManager) Install(packages ...string) error {
	m.installCount++
	if m.installCount <= m.failCount {
		return fmt.Errorf("mock install error")
	}
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