package install

import (
	"strings"
	"testing"
)

func TestValidateTool(t *testing.T) {
	tests := []struct {
		name    string
		tool    *Tool
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid tool",
			tool: &Tool{
				Name:         "test-tool",
				PackageName:  "test-package",
				Version:      "1.0.0",
				Dependencies: []string{"dep1", "dep2"},
				PostInstall: []PostInstallCommand{
					{Command: "echo 'test'", Description: "Test command"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			tool: &Tool{
				Name:        "",
				PackageName: "test-package",
			},
			wantErr: true,
			errMsg:  "Name: cannot be empty",
		},
		{
			name: "invalid name characters",
			tool: &Tool{
				Name:        "test tool",
				PackageName: "test-package",
			},
			wantErr: true,
			errMsg:  "Name: contains invalid characters",
		},
		{
			name: "empty package name",
			tool: &Tool{
				Name:        "test-tool",
				PackageName: "",
			},
			wantErr: true,
			errMsg:  "PackageName: cannot be empty",
		},
		{
			name: "invalid version",
			tool: &Tool{
				Name:        "test-tool",
				PackageName: "test-package",
				Version:     "invalid",
			},
			wantErr: true,
			errMsg:  "Version: invalid version format",
		},
		{
			name: "empty dependency",
			tool: &Tool{
				Name:         "test-tool",
				PackageName:  "test-package",
				Dependencies: []string{""},
			},
			wantErr: true,
			errMsg:  "Dependencies[0]: cannot be empty",
		},
		{
			name: "empty post-install command",
			tool: &Tool{
				Name:        "test-tool",
				PackageName: "test-package",
				PostInstall: []PostInstallCommand{
					{Command: "", Description: "Empty command"},
				},
			},
			wantErr: true,
			errMsg:  "PostInstall[0].Command: cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTool(tt.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateTool() error = %v, want error containing %v", err, tt.errMsg)
			}
		})
	}
} 