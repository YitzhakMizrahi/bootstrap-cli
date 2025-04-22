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
	name            string
}

func (m *MockPackageManager) Name() string {
	return m.name
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
		name           string
		tool           *Tool
		maxRetries     int
		maxFail        int
		pmName         string
		wantErr        bool
		expectedPkgName string
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
			maxRetries:     3,
			maxFail:        0,
			pmName:         "apt",
			wantErr:        false,
			expectedPkgName: "test-package=1.0.0",
		},
		{
			name: "retry success",
			tool: &Tool{
				Name:         "retry-tool",
				PackageName:  "retry-package",
				Dependencies: []string{"dep1"},
			},
			maxRetries: 3,
			maxFail:    2,
			pmName:     "apt",
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
			maxFail:    4,
			pmName:     "apt",
			wantErr:    true,
		},
		{
			name: "system specific package name",
			tool: &Tool{
				Name:        "system-tool",
				PackageName: "default-package",
				PackageNames: &PackageMapping{
					Default: "default-package",
					APT:     "apt-package",
					DNF:     "dnf-package",
					Pacman:  "pacman-package",
					Brew:    "brew-package",
				},
				Version: "2.0.0",
			},
			maxRetries:     3,
			maxFail:        0,
			pmName:         "apt",
			wantErr:        false,
			expectedPkgName: "apt-package=2.0.0",
		},
		{
			name: "homebrew version format",
			tool: &Tool{
				Name:        "brew-tool",
				PackageName: "brew-package",
				Version:     "3.0.0",
			},
			maxRetries:     3,
			maxFail:        0,
			pmName:         "brew",
			wantErr:        false,
			expectedPkgName: "brew-package@3.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock package manager
			mockPM := &MockPackageManager{
				installedPackages: make(map[string]bool),
				failCount:        make(map[string]int),
				maxFailures:      tt.maxFail,
				name:            tt.pmName,
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

				// Verify correct package name was used
				if tt.expectedPkgName != "" && !mockPM.IsInstalled(tt.expectedPkgName) {
					t.Errorf("Expected package %s was not installed", tt.expectedPkgName)
				}
			}

			// Check if dependencies were installed
			if !tt.wantErr {
				for _, dep := range tt.tool.Dependencies {
					if !mockPM.IsInstalled(dep) {
						t.Errorf("Dependency %s was not installed", dep)
					}
				}
			}
		})
	}
} 