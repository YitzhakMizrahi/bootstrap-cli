package install

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// MockPackageManager simulates a package manager for testing
type MockPackageManager struct {
	installed       map[string]string // package -> version
	failureCount    map[string]int
	maxFailures     int
	removedPackages []string
	name           string
}

func NewMockPackageManager(maxFailures int, name string) *MockPackageManager {
	return &MockPackageManager{
		installed:       make(map[string]string),
		failureCount:    make(map[string]int),
		maxFailures:     maxFailures,
		removedPackages: make([]string, 0),
		name:           name,
	}
}

func (m *MockPackageManager) Install(packageName string) error {
	// Extract package name and version
	name, version := packageName, ""
	if strings.Contains(packageName, "=") {
		parts := strings.Split(packageName, "=")
		name, version = parts[0], parts[1]
	} else if strings.Contains(packageName, "@") {
		parts := strings.Split(packageName, "@")
		name, version = parts[0], parts[1]
	}

	if m.failureCount[name] < m.maxFailures {
		m.failureCount[name]++
		return errors.New("simulated install failure")
	}
	m.installed[name] = version
	return nil
}

func (m *MockPackageManager) IsInstalled(pkg string) (bool, error) {
	// Extract package name without version
	name := pkg
	if strings.Contains(pkg, "=") {
		name = strings.Split(pkg, "=")[0]
	} else if strings.Contains(pkg, "@") {
		name = strings.Split(pkg, "@")[0]
	}
	_, exists := m.installed[name]
	return exists, nil
}

func (m *MockPackageManager) IsPackageAvailable(pkg string) bool {
	// Simple mock: assume all packages are available
	return true
}

func (m *MockPackageManager) Uninstall(pkg string) error {
	// Extract package name without version
	name := pkg
	if strings.Contains(pkg, "=") {
		name = strings.Split(pkg, "=")[0]
	} else if strings.Contains(pkg, "@") {
		name = strings.Split(pkg, "@")[0]
	}

	if m.failureCount[name] < m.maxFailures {
		m.failureCount[name]++
		return errors.New("simulated uninstall failure")
	}
	if _, exists := m.installed[name]; exists {
		delete(m.installed, name)
		m.removedPackages = append(m.removedPackages, name)
	}
	return nil
}

func (m *MockPackageManager) IsAvailable() bool {
	return true
}

func (m *MockPackageManager) GetName() string {
	return m.name
}

func (m *MockPackageManager) Update() error {
	return nil
}

func (m *MockPackageManager) Upgrade() error {
	return nil
}

func (m *MockPackageManager) GetVersion(packageName string) (string, error) {
	if version, exists := m.installed[packageName]; exists {
		return version, nil
	}
	return "", fmt.Errorf("package %s not installed", packageName)
}

func (m *MockPackageManager) ListInstalled() ([]string, error) {
	packages := make([]string, 0, len(m.installed))
	for pkg := range m.installed {
		packages = append(packages, pkg)
	}
	return packages, nil
}

func (m *MockPackageManager) SetupSpecialPackage(_ string) error {
	// For testing purposes, we'll simulate that we can set up any package
	return nil
}

func TestInstaller(t *testing.T) {
	tests := []struct {
		name            string
		tool            *interfaces.Tool
		maxRetries      int
		maxFail         int
		pmName          string
		wantErr         bool
		expectedPkgName string
		expectCleanup   bool
		cleanupPackages []string
	}{
		{
			name: "successful install",
			tool: &interfaces.Tool{
				Name:         "test-tool",
				Description:  "Test tool",
				Version:      "1.0.0",
				Dependencies: []struct {
					Name     string `yaml:"name"`
					Type     string `yaml:"type"`
					Optional bool   `yaml:"optional,omitempty"`
				}{
					{Name: "dep1", Type: "package"},
					{Name: "dep2", Type: "package"},
				},
				PostInstall: []struct {
					Command     string `yaml:"command"`
					Description string `yaml:"description"`
				}{
					{
						Command:     "echo 'test'",
						Description: "Test command",
					},
				},
			},
			maxRetries:      3,
			maxFail:         0,
			pmName:          "apt",
			wantErr:         false,
			expectedPkgName: "test-tool=1.0.0",
			expectCleanup:   false,
		},
		{
			name: "retry success",
			tool: &interfaces.Tool{
				Name:         "retry-tool",
				Description:  "Retry tool",
				Dependencies: []struct {
					Name     string `yaml:"name"`
					Type     string `yaml:"type"`
					Optional bool   `yaml:"optional,omitempty"`
				}{
					{Name: "dep1", Type: "package"},
				},
			},
			maxRetries:     3,
			maxFail:        2,
			pmName:         "apt",
			wantErr:        false,
			expectCleanup:  false,
		},
		{
			name: "retry failure",
			tool: &interfaces.Tool{
				Name:         "fail-tool",
				Description:  "Fail tool",
				Dependencies: []struct {
					Name     string `yaml:"name"`
					Type     string `yaml:"type"`
					Optional bool   `yaml:"optional,omitempty"`
				}{
					{Name: "dep1", Type: "package"},
				},
			},
			maxRetries:      3,
			maxFail:        4,
			pmName:         "apt",
			wantErr:        true,
			expectCleanup:  true,
			cleanupPackages: []string{"dep1"},
		},
		{
			name: "system specific package name",
			tool: &interfaces.Tool{
				Name:        "system-tool",
				Description: "System tool",
				PackageNames: struct {
					APT    string `yaml:"apt"`
					Brew   string `yaml:"brew"`
					DNF    string `yaml:"dnf"`
					Pacman string `yaml:"pacman"`
				}{
					APT:    "apt-package",
					DNF:    "dnf-package",
					Pacman: "pacman-package",
					Brew:   "brew-package",
				},
				Version: "2.0.0",
			},
			maxRetries:      3,
			maxFail:         0,
			pmName:          "apt",
			wantErr:         false,
			expectedPkgName: "apt-package=2.0.0",
			expectCleanup:   false,
		},
		{
			name: "homebrew version format",
			tool: &interfaces.Tool{
				Name:        "brew-tool",
				Description: "Brew tool",
				Version:     "3.0.0",
			},
			maxRetries:      3,
			maxFail:         0,
			pmName:          "brew",
			wantErr:         false,
			expectedPkgName: "brew-package@3.0.0",
			expectCleanup:   false,
		},
		{
			name: "post-install failure cleanup",
			tool: &interfaces.Tool{
				Name:         "post-fail-tool",
				Description:  "Post fail tool",
				Dependencies: []struct {
					Name     string `yaml:"name"`
					Type     string `yaml:"type"`
					Optional bool   `yaml:"optional,omitempty"`
				}{
					{Name: "dep1", Type: "package"},
				},
				PostInstall: []struct {
					Command     string `yaml:"command"`
					Description string `yaml:"description"`
				}{
					{
						Command:     "exit 1",
						Description: "Failing command",
					},
				},
			},
			maxRetries:      3,
			maxFail:         0,
			pmName:         "apt",
			wantErr:        true,
			expectCleanup:  true,
			cleanupPackages: []string{"dep1", "post-fail-tool"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock package manager
			mockPM := NewMockPackageManager(tt.maxFail, tt.pmName)

			// Create an installer with the mock package manager
			installer := &Installer{
				PackageManager: mockPM,
				Logger:        log.New(log.InfoLevel),
				MaxRetries:    tt.maxRetries,
				RetryDelay:    time.Millisecond, // Use short delay for tests
			}

			// Install the tool
			err := installer.Install(tt.tool)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if package was installed with correct version
			if !tt.wantErr && tt.expectedPkgName != "" {
				name, version := tt.expectedPkgName, ""
				if strings.Contains(tt.expectedPkgName, "=") {
					parts := strings.Split(tt.expectedPkgName, "=")
					name, version = parts[0], parts[1]
				} else if strings.Contains(tt.expectedPkgName, "@") {
					parts := strings.Split(tt.expectedPkgName, "@")
					name, version = parts[0], parts[1]
				}

				if installedVersion, ok := mockPM.installed[name]; !ok {
					t.Errorf("Expected package %s was not installed", tt.expectedPkgName)
				} else if version != "" && installedVersion != version {
					t.Errorf("Expected package %s to be installed with version %s, got %s", name, version, installedVersion)
				}
			}

			// Check cleanup
			if tt.expectCleanup {
				for _, pkg := range tt.cleanupPackages {
					if _, exists := mockPM.installed[pkg]; exists {
						t.Errorf("Expected package %s to be removed during cleanup", pkg)
					}
				}
				if len(mockPM.installed) > 0 {
					t.Error("Expected packages to be removed during cleanup")
				}
			}
		})
	}
}

func TestTool_Install(t *testing.T) {
	tool := &interfaces.Tool{
		Name:        "test-tool",
		Description: "Test tool",
		PostInstall: []struct {
			Command     string `yaml:"command"`
			Description string `yaml:"description"`
		}{
			{
				Command:     "echo 'test'",
				Description: "Test command",
			},
		},
	}

	// Create a mock package manager
	mockPM := NewMockPackageManager(0, "apt")

	// Create an installer with the mock package manager
	installer := &Installer{
		PackageManager: mockPM,
		Logger:        log.New(log.InfoLevel),
		MaxRetries:    3,
		RetryDelay:    time.Millisecond, // Use short delay for tests
	}

	// Install the tool
	err := installer.Install(tool)
	if err != nil {
		t.Errorf("Install() error = %v, wantErr false", err)
	}

	// Verify the tool was installed
	installed, err := mockPM.IsInstalled(tool.Name)
	if err != nil {
		t.Errorf("IsInstalled() error = %v", err)
	}
	if !installed {
		t.Errorf("Expected package %s to be installed", tool.Name)
	}
} 