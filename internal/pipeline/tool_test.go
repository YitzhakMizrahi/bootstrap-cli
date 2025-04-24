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
		Command:        "test-tool --version",
		ExpectedOutput: "v1.0.0",
		BinaryPaths:    []string{"test-tool"},
	}
	tool.SetVerification(verify)

	if tool.Verify.Command != verify.Command {
		t.Errorf("Expected command '%s', got '%s'", verify.Command, tool.Verify.Command)
	}
}

func TestTool_SetInstallation(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	install := InstallStrategy{
		PackageNames: map[string]string{
			"apt":  "test-tool",
			"brew": "test-tool",
		},
		PreInstall:  []string{"echo 'pre-install'"},
		PostInstall: []string{"echo 'post-install'"},
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

func TestTool_CreateInstallationSteps(t *testing.T) {
	tool := NewTool("test-tool", CategoryDevelopment)
	
	// Set up installation strategy
	install := InstallStrategy{
		PackageNames: map[string]string{
			"apt": "test-tool",
		},
		PreInstall:  []string{"echo 'pre-install'"},
		PostInstall: []string{"echo 'post-install'"},
	}
	tool.SetInstallation(install)

	// Set up verification
	verify := VerifyStrategy{
		Command:        "test-tool --version",
		ExpectedOutput: "v1.0.0",
		BinaryPaths:    []string{"test-tool"},
	}
	tool.SetVerification(verify)

	// Create steps
	platform := &Platform{
		OS:             "linux",
		PackageManager: "apt",
	}
	steps := tool.CreateInstallationSteps(platform)

	// Verify number of steps (pre-install + install + post-install + verify)
	expectedSteps := len(install.PreInstall) + 1 + len(install.PostInstall) + 1
	if len(steps) != expectedSteps {
		t.Errorf("Expected %d steps, got %d", expectedSteps, len(steps))
	}

	// Verify step names
	if steps[0].Name != "test-tool-pre-install-0" {
		t.Errorf("Expected first step name 'test-tool-pre-install-0', got '%s'", steps[0].Name)
	}
	if steps[len(steps)-1].Name != "test-tool-verify" {
		t.Errorf("Expected last step name 'test-tool-verify', got '%s'", steps[len(steps)-1].Name)
	}
}

func TestTool_CustomInstallation(t *testing.T) {
	tool := NewTool("custom-tool", CategoryDevelopment)
	
	// Set up custom installation strategy
	install := InstallStrategy{
		CustomInstall: []string{
			"wget https://example.com/custom-tool",
			"chmod +x custom-tool",
			"sudo mv custom-tool /usr/local/bin/",
		},
	}
	tool.SetInstallation(install)

	platform := &Platform{
		OS:             "linux",
		PackageManager: "apt",
	}
	steps := tool.CreateInstallationSteps(platform)

	// Verify number of steps (custom install commands + verify)
	expectedSteps := len(install.CustomInstall) + 1
	if len(steps) != expectedSteps {
		t.Errorf("Expected %d steps, got %d", expectedSteps, len(steps))
	}

	// Verify step names
	if steps[0].Name != "custom-tool-custom-install-0" {
		t.Errorf("Expected first step name 'custom-tool-custom-install-0', got '%s'", steps[0].Name)
	}
} 