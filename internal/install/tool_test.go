package install

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// MockPackageManager implements PackageManager for testing
type MockPackageManager struct {
	installedPackages map[string]bool
	failCount        map[string]int
	maxFailures      int
}

func (m *MockPackageManager) Name() string {
	return "mock"
}

func (m *MockPackageManager) IsAvailable() bool {
	return true
}

func (m *MockPackageManager) Install(packages ...string) error {
	for _, pkg := range packages {
		if m.failCount[pkg] < m.maxFailures {
			m.failCount[pkg]++
			return errors.New("simulated failure")
		}
		m.installedPackages[pkg] = true
	}
	return nil
}

func (m *MockPackageManager) Update() error {
	return nil
}

func (m *MockPackageManager) IsInstalled(pkg string) bool {
	return m.installedPackages[pkg]
}

func TestInstaller(t *testing.T) {
	tests := []struct {
		name       string
		tool       *Tool
		maxRetries int
		maxFail    int
		wantErr    bool
	}{
		{
			name: "successful install",
			tool: &Tool{
				Name:         "test-tool",
				PackageName:  "test-package",
				Version:      "1.0.0",
				Dependencies: []string{"dep1", "dep2"},
				PostInstall:  []string{"echo 'test'"},
			},
			maxRetries: 3,
			maxFail:    0,
			wantErr:    false,
		},
		{
			name: "retry success",
			tool: &Tool{
				Name:         "retry-tool",
				PackageName:  "retry-package",
				Dependencies: []string{"dep1"},
			},
			maxRetries: 3,
			maxFail:    2, // Will succeed on third try
			wantErr:    false,
		},
		{
			name: "retry failure",
			tool: &Tool{
				Name:         "fail-tool",
				PackageName:  "fail-package",
				Dependencies: []string{"dep1"},
			},
			maxRetries: 3,
			maxFail:    4, // Will fail all attempts (3 retries + 1 initial attempt)
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock package manager
			mockPM := &MockPackageManager{
				installedPackages: make(map[string]bool),
				failCount:        make(map[string]int),
				maxFailures:      tt.maxFail,
			}

			// Create a buffer for logging output
			var logBuf bytes.Buffer
			logger := log.New(log.DebugLevel)
			logger.SetOutput(&logBuf)

			// Create an installer with custom settings
			installer := &Installer{
				PackageManager: mockPM,
				Logger:        logger,
				MaxRetries:    tt.maxRetries,
				RetryDelay:    time.Millisecond, // Use short delay for tests
			}

			// Install the tool
			err := installer.Install(tt.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check log output
			logOutput := logBuf.String()
			if tt.wantErr {
				if !strings.Contains(logOutput, "ERROR") {
					t.Error("Expected error log message not found")
				}
			} else {
				if !strings.Contains(logOutput, "Successfully installed") {
					t.Error("Expected success log message not found")
				}
			}

			// Check if dependencies and package were installed
			if !tt.wantErr {
				for _, dep := range tt.tool.Dependencies {
					if !mockPM.IsInstalled(dep) {
						t.Errorf("Dependency %s was not installed", dep)
					}
				}
				if !mockPM.IsInstalled(tt.tool.PackageName) {
					t.Errorf("Tool %s was not installed", tt.tool.PackageName)
				}
			}
		})
	}
} 