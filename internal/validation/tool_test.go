package validation

import (
	"strings"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
)

func TestValidateTool(t *testing.T) {
	tests := []struct {
		name    string
		tool    *install.Tool
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tool",
			tool: &install.Tool{
				Name:         "test-tool",
				PackageName:  "test-package",
				Version:      "1.0.0",
				Dependencies: []string{"dep1", "dep2"},
				PostInstall: []install.PostInstallCommand{
					{Command: "echo 'test'", Description: "Test command"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			tool: &install.Tool{
				Name:        "",
				PackageName: "test-package",
			},
			wantErr: true,
			errMsg:  "Name: cannot be empty",
		},
		{
			name: "invalid name characters",
			tool: &install.Tool{
				Name:        "test tool",
				PackageName: "test-package",
			},
			wantErr: true,
			errMsg:  "Name: contains invalid characters",
		},
		{
			name: "empty package name",
			tool: &install.Tool{
				Name:        "test-tool",
				PackageName: "",
			},
			wantErr: true,
			errMsg:  "PackageName: cannot be empty",
		},
		{
			name: "invalid version",
			tool: &install.Tool{
				Name:        "test-tool",
				PackageName: "test-package",
				Version:     "invalid",
			},
			wantErr: true,
			errMsg:  "Version: invalid version format",
		},
		{
			name: "empty dependency",
			tool: &install.Tool{
				Name:         "test-tool",
				PackageName:  "test-package",
				Dependencies: []string{""},
			},
			wantErr: true,
			errMsg:  "Dependencies[0]: cannot be empty",
		},
		{
			name: "empty post-install command",
			tool: &install.Tool{
				Name:        "test-tool",
				PackageName: "test-package",
				PostInstall: []install.PostInstallCommand{
					{Command: "", Description: "Empty command"},
				},
			},
			wantErr: true,
			errMsg:  "PostInstall[0].Command: cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTool(tt.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateTool() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
} 