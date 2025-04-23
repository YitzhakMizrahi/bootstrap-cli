package tools

import (
	"testing"
)

func TestCoreTools(t *testing.T) {
	tools := CoreTools()
	if len(tools) == 0 {
		t.Error("Expected CoreTools to return at least one tool")
	}

	// Check for essential tools
	essentialTools := map[string]bool{
		"git":  false,
		"curl": false,
		"wget": false,
	}

	for _, tool := range tools {
		if _, exists := essentialTools[tool.PackageName]; exists {
			essentialTools[tool.PackageName] = true
		}
	}

	for tool, found := range essentialTools {
		if !found {
			t.Errorf("Essential tool %s not found in CoreTools", tool)
		}
	}
}

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "valid command",
			cmd:     "echo hello",
			wantErr: false,
		},
		{
			name:    "invalid command",
			cmd:     "nonexistentcommand",
			wantErr: true,
		},
		{
			name:    "empty command",
			cmd:     "",
			wantErr: true,
		},
		{
			name:    "command with arguments",
			cmd:     "echo hello world",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runCommand(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} 