package system

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
		failCount: 2,
	}
	retryPM := &retryPackageManager{
		PackageManager: mockPM,
		maxRetries:    3,
		retryDelay:    10 * time.Millisecond,
	}

	// Test Install
	err := retryPM.Install("test-package")
	if err != nil {
		t.Errorf("Install() error = %v", err)
	}
	if mockPM.installCount != 3 {
		t.Errorf("expected 3 install attempts, got %d", mockPM.installCount)
	}

	// Test Uninstall
	err = retryPM.Uninstall("test-package")
	if err != nil {
		t.Errorf("Uninstall() error = %v", err)
	}
	if mockPM.uninstallCount != 3 {
		t.Errorf("expected 3 uninstall attempts, got %d", mockPM.uninstallCount)
	}
}

// mockPackageManager is a test implementation that fails a specified number of times
type mockPackageManager struct {
	failCount      int
	installCount   int
	uninstallCount int
}

func (m *mockPackageManager) Install(pkg string) error {
	m.installCount++
	if m.installCount <= m.failCount {
		return fmt.Errorf("mock install error")
	}
	return nil
}

func (m *mockPackageManager) Uninstall(pkg string) error {
	m.uninstallCount++
	if m.uninstallCount <= m.failCount {
		return fmt.Errorf("mock uninstall error")
	}
	return nil
}

func (m *mockPackageManager) IsInstalled(pkg string) bool {
	return true
} 