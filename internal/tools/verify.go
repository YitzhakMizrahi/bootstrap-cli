package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
)

// VerifyTool checks if a tool is properly installed and accessible
func VerifyTool(tool *interfaces.Tool, additionalPaths []string) error {
	// Check if binary exists in PATH or additional paths
	paths := append(filepath.SplitList(os.Getenv("PATH")), additionalPaths...)
	
	var binaryPath string
	for _, path := range paths {
		fullPath := filepath.Join(path, tool.Name)
		if _, err := os.Stat(fullPath); err == nil {
			binaryPath = fullPath
			break
		}
	}

	if binaryPath == "" {
		return fmt.Errorf("binary not found in PATH or additional paths")
	}

	// Try to execute the tool with --version or -v
	cmd := exec.Command(binaryPath, "--version")
	if err := cmd.Run(); err != nil {
		// Try -v if --version fails
		cmd = exec.Command(binaryPath, "-v")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}
	}

	return nil
}

// VerifyCoreTools verifies that all core tools are properly installed
func VerifyCoreTools(opts *InstallOptions) error {
	for _, tool := range opts.Tools {
		if err := VerifyTool(tool, opts.AdditionalPaths); err != nil {
			return fmt.Errorf("%s: verification failed: %w", tool.Name, err)
		}
		opts.Logger.Info("Successfully verified %s", tool.Name)
	}
	return nil
} 