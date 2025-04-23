package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModernTools(t *testing.T) {
	// Create a temporary test environment
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	require.NoError(t, os.MkdirAll(homeDir, 0755))
	
	// Setup test environment
	env := testutil.NewTestEnv(t, homeDir)
	defer env.Cleanup()

	// Load modern tools configuration
	loader := config.NewConfigLoader("")
	tools, err := loader.LoadTools()
	require.NoError(t, err)

	var modernTools []*interfaces.Tool
	for _, tool := range tools {
		if tool.Category == "modern" {
			modernTools = append(modernTools, tool)
		}
	}

	// Verify we have all expected modern tools
	expectedTools := []string{"bat", "fzf", "ripgrep", "fd", "lsd"}
	assert.Len(t, modernTools, len(expectedTools), "Should have exactly %d modern tools", len(expectedTools))

	toolMap := make(map[string]*interfaces.Tool)
	for _, tool := range modernTools {
		toolMap[tool.Name] = tool
	}

	for _, name := range expectedTools {
		t.Run(name, func(t *testing.T) {
			tool, exists := toolMap[name]
			require.True(t, exists, "Tool %s should exist", name)
			
			// Test tool configuration
			assert.NotEmpty(t, tool.Description, "Tool should have a description")
			assert.Equal(t, "modern", tool.Category, "Tool should be in modern category")
			assert.Contains(t, tool.Tags, "modern", "Tool should have modern tag")
			
			// Test package names for different package managers
			assert.NotEmpty(t, tool.PackageNames.APT, "Should have APT package name")
			assert.NotEmpty(t, tool.PackageNames.Brew, "Should have Homebrew package name")
			assert.NotEmpty(t, tool.PackageNames.DNF, "Should have DNF package name")
			assert.NotEmpty(t, tool.PackageNames.Pacman, "Should have Pacman package name")

			// Test verification command
			assert.NotEmpty(t, tool.VerifyCommand, "Should have verification command")

			// Test shell configuration
			assert.NotNil(t, tool.ShellConfig, "Should have shell configuration")
			assert.NotEmpty(t, tool.ShellConfig.Aliases, "Should have aliases")

			// Tool-specific tests
			switch name {
			case "bat":
				testBatConfig(t, tool)
			case "fzf":
				testFzfConfig(t, tool)
			case "ripgrep":
				testRipgrepConfig(t, tool)
			case "fd":
				testFdConfig(t, tool)
			case "lsd":
				testLsdConfig(t, tool)
			}
		})
	}
}

func testBatConfig(t *testing.T, tool *interfaces.Tool) {
	assert.Contains(t, tool.ShellConfig.Aliases, "cat", "Should have cat alias")
	assert.Contains(t, tool.ShellConfig.Env, "BAT_THEME", "Should have BAT_THEME env var")
	assert.Contains(t, tool.ShellConfig.Env, "BAT_STYLE", "Should have BAT_STYLE env var")
	
	var hasConfigDir bool
	for _, cmd := range tool.PostInstall {
		if cmd.Command == "mkdir -p ~/.config/bat" {
			hasConfigDir = true
			break
		}
	}
	assert.True(t, hasConfigDir, "Should create config directory")
}

func testFzfConfig(t *testing.T, tool *interfaces.Tool) {
	assert.Contains(t, tool.ShellConfig.Env, "FZF_DEFAULT_OPTS", "Should have FZF_DEFAULT_OPTS env var")
	assert.Contains(t, tool.ShellConfig.Env, "FZF_DEFAULT_COMMAND", "Should have FZF_DEFAULT_COMMAND env var")
	assert.Contains(t, tool.ShellConfig.Functions, "fcd", "Should have fcd function")
	assert.Contains(t, tool.ShellConfig.Functions, "fkill", "Should have fkill function")
}

func testRipgrepConfig(t *testing.T, tool *interfaces.Tool) {
	assert.Contains(t, tool.ShellConfig.Aliases, "grep", "Should have grep alias")
	assert.Contains(t, tool.ShellConfig.Env, "RIPGREP_CONFIG_PATH", "Should have RIPGREP_CONFIG_PATH env var")
	
	// Test ripgrep config file
	var hasConfigFile bool
	for _, file := range tool.Files {
		if file.Destination == "~/.config/ripgrep/config" {
			assert.Contains(t, file.Content, "--smart-case", "Config should include smart-case")
			hasConfigFile = true
			break
		}
	}
	assert.True(t, hasConfigFile, "Should have ripgrep config file")
}

func testFdConfig(t *testing.T, tool *interfaces.Tool) {
	assert.Contains(t, tool.ShellConfig.Aliases, "find", "Should have find alias")
	assert.Contains(t, tool.ShellConfig.Env, "FD_OPTIONS", "Should have FD_OPTIONS env var")
	assert.Contains(t, tool.ShellConfig.Functions, "fdsize", "Should have fdsize function")
	assert.Contains(t, tool.ShellConfig.Functions, "fdnewer", "Should have fdnewer function")
}

func testLsdConfig(t *testing.T, tool *interfaces.Tool) {
	assert.Contains(t, tool.ShellConfig.Aliases, "ls", "Should have ls alias")
	assert.Contains(t, tool.ShellConfig.Aliases, "ll", "Should have ll alias")
	assert.Contains(t, tool.ShellConfig.Aliases, "la", "Should have la alias")
	assert.Contains(t, tool.ShellConfig.Env, "LS_COLORS", "Should have LS_COLORS env var")
} 