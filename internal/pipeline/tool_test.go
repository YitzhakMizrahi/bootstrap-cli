package pipeline

import (
	"testing"
)

func TestTool_NewTool(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	if tool.Name != "test-tool" {
		t.Errorf("Expected tool name 'test-tool', got '%s'", tool.Name)
	}
	if tool.Category != CategoryDevelopment {
		t.Errorf("Expected category Development, got '%s'", tool.Category)
	}
}

func TestTool_AddDependency(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	dep := Dependency{
		Name:     "dep1",
		Type:     SystemDependency,
		Optional: false,
	}
	tool.AddDependency(dep)

	if len(tool.Dependencies) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(tool.Dependencies))
	}
	if tool.Dependencies[0].Name != "dep1" {
		t.Errorf("Expected dependency name 'dep1', got '%s'", tool.Dependencies[0].Name)
	}
}

func TestTool_SetVerification(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	verify := VerifyStrategy{
		Command: Command{
			Command:     "test-tool --version",
			Description: "Check test-tool version",
		},
		ExpectedOutput: "v1.0.0",
		BinaryPaths:    []string{"test-tool"},
	}
	tool.SetVerification(verify)

	if tool.Verify.Command.Command != verify.Command.Command {
		t.Errorf("Expected command '%s', got '%s'", verify.Command.Command, tool.Verify.Command.Command)
	}
}

func TestTool_SetInstallation(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	install := InstallStrategy{
		PackageNames: map[string]string{
			"apt":  "test-tool",
			"brew": "test-tool",
		},
		PreInstall: []Command{
			{
				Command:     "echo 'pre-install'",
				Description: "Pre-install command",
			},
		},
		PostInstall: []Command{
			{
				Command:     "echo 'post-install'",
				Description: "Post-install command",
			},
		},
	}
	tool.SetInstallation(install)

	if tool.Install.PackageNames["apt"] != "test-tool" {
		t.Errorf("Expected apt package name 'test-tool', got '%s'", tool.Install.PackageNames["apt"])
	}
}

func TestTool_GetInstallStrategy(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	
	// Set default strategy
	defaultStrategy := InstallStrategy{
		PackageNames: map[string]string{
			"apt": "test-tool",
		},
	}
	tool.SetInstallation(defaultStrategy)

	// Set platform-specific strategy
	linuxStrategy := InstallStrategy{
		PackageNames: map[string]string{
			"apt": "test-tool-linux",
		},
	}
	tool.SetPlatformConfig("linux", linuxStrategy)

	// Test platform-specific strategy
	platform := &Platform{OS: "linux"}
	strategy := tool.GetInstallStrategy(platform)
	if strategy.PackageNames["apt"] != "test-tool-linux" {
		t.Errorf("Expected linux package name 'test-tool-linux', got '%s'", strategy.PackageNames["apt"])
	}

	// Test default strategy
	platform.OS = "darwin"
	strategy = tool.GetInstallStrategy(platform)
	if strategy.PackageNames["apt"] != "test-tool" {
		t.Errorf("Expected default package name 'test-tool', got '%s'", strategy.PackageNames["apt"])
	}
}

func TestTool_GenerateInstallationSteps(t *testing.T) {
	// Create a tool with installation strategy
	tool := NewTool("test-tool", CategoryEssential)
	install := InstallStrategy{
		PreInstall: []Command{
			{
				Command:     "echo 'pre-install'",
				Description: "Pre-install command",
			},
		},
		PostInstall: []Command{
			{
				Command:     "echo 'post-install'",
				Description: "Post-install command",
			},
		},
		PackageNames: map[string]string{
			"apt": "test-tool",
		},
	}
	tool.SetInstallation(install)

	// Create steps
	platform := &Platform{
		OS:             "linux",
		PackageManager: "apt",
	}
	// Create a dummy channel for the context
	dummyChan := make(chan ProgressEvent, 1)
	defer close(dummyChan)
	context := NewInstallationContext(platform, nil, dummyChan) // Pass dummy channel
	steps := tool.GenerateInstallationSteps(platform, context, false)

	// Verify number of steps (pre-install + install + post-install + verify)
	expectedSteps := len(install.PreInstall) + 1 + len(install.PostInstall) + 1
	if len(steps) != expectedSteps {
		t.Errorf("Expected %d steps, got %d", expectedSteps, len(steps))
	}

	// Test generating installation steps
	steps = tool.GenerateInstallationSteps(context.Platform, context, false)
	if len(steps) == 0 {
		t.Error("Expected installation steps to be generated")
	}

	// Verify that the steps include installation and verification
}

func TestTool_CustomInstallation(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	
	// Set up custom installation
	install := InstallStrategy{
		CustomInstall: []Command{
			{
				Command:     "echo 'custom install step 1'",
				Description: "Custom install step 1",
			},
			{
				Command:     "echo 'custom install step 2'",
				Description: "Custom install step 2",
			},
		},
	}
	tool.SetInstallation(install)
	
	// Create steps
	platform := &Platform{
		OS:             "linux",
		PackageManager: "apt",
	}
	// Create a dummy channel for the context
	dummyChan := make(chan ProgressEvent, 1)
	defer close(dummyChan)
	context := NewInstallationContext(platform, nil, dummyChan) // Pass dummy channel
	steps := tool.GenerateInstallationSteps(platform, context, false)
	
	// Verify number of steps (custom install steps + verify)
	expectedSteps := len(install.CustomInstall) + 1
	if len(steps) != expectedSteps {
		t.Errorf("Expected %d steps, got %d", expectedSteps, len(steps))
	}
} 