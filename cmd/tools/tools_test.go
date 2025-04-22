package tools

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewToolsCmd(t *testing.T) {
	cmd := NewToolsCmd()

	// Test command structure
	if cmd.Use != "tools" {
		t.Errorf("Expected Use to be 'tools', got %s", cmd.Use)
	}

	// Test subcommands
	subCmds := cmd.Commands()
	if len(subCmds) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(subCmds))
	}

	// Find install and verify commands
	var installCmd, verifyCmd *cobra.Command
	for _, sub := range subCmds {
		switch sub.Use {
		case "install":
			installCmd = sub
		case "verify":
			verifyCmd = sub
		}
	}

	// Test install command
	if installCmd == nil {
		t.Error("Install command not found")
	} else {
		// Test flags
		flags := installCmd.Flags()
		skipVerifyFlag := flags.Lookup("skip-verify")
		if skipVerifyFlag == nil {
			t.Error("skip-verify flag not found")
		}
		if skipVerifyFlag != nil && skipVerifyFlag.DefValue != "false" {
			t.Errorf("Expected skip-verify default value to be false, got %s", skipVerifyFlag.DefValue)
		}
	}

	// Test verify command
	if verifyCmd == nil {
		t.Error("Verify command not found")
	}
}

func TestCommandHelp(t *testing.T) {
	cmd := NewToolsCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	// Test main command help
	if err := cmd.Help(); err != nil {
		t.Errorf("Error getting help: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Manage development tools")) {
		t.Error("Help output missing command description")
	}

	// Test install command help
	buf.Reset()
	installCmd, _, err := cmd.Find([]string{"install"})
	if err != nil {
		t.Errorf("Error finding install command: %v", err)
	}
	if err := installCmd.Help(); err != nil {
		t.Errorf("Error getting install help: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Install core development tools")) {
		t.Error("Help output missing install command description")
	}

	// Test verify command help
	buf.Reset()
	verifyCmd, _, err := cmd.Find([]string{"verify"})
	if err != nil {
		t.Errorf("Error finding verify command: %v", err)
	}
	if err := verifyCmd.Help(); err != nil {
		t.Errorf("Error getting verify help: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("Verify tool installations")) {
		t.Error("Help output missing verify command description")
	}
} 