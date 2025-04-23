package interfaces

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a shell command and returns any error
func RunCommand(command string) error {
	// Split the command into parts
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	// Create the command
	cmd := exec.Command(parts[0], parts[1:]...)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		// If there's stderr output, include it in the error
		if stderr.Len() > 0 {
			return fmt.Errorf("%v: %s", err, stderr.String())
		}
		return err
	}

	return nil
} 