package tools

import "testing"

func TestCoreTools(t *testing.T) {
	tools := CoreTools()

	// Verify we have tools defined
	if len(tools) == 0 {
		t.Error("CoreTools() returned empty list")
	}

	// Verify each tool has required fields
	for _, tool := range tools {
		if tool.Name == "" {
			t.Error("Tool name is empty")
		}
		if tool.PackageName == "" {
			t.Errorf("Package name is empty for tool %s", tool.Name)
		}
		if tool.PackageNames == nil {
			t.Errorf("PackageNames mapping is nil for tool %s", tool.Name)
		} else {
			// Verify package mappings
			if tool.PackageNames.Default == "" {
				t.Errorf("Default package name is empty for tool %s", tool.Name)
			}
			if tool.PackageNames.APT == "" {
				t.Errorf("APT package name is empty for tool %s", tool.Name)
			}
			if tool.PackageNames.DNF == "" {
				t.Errorf("DNF package name is empty for tool %s", tool.Name)
			}
			if tool.PackageNames.Pacman == "" {
				t.Errorf("Pacman package name is empty for tool %s", tool.Name)
			}
			if tool.PackageNames.Brew == "" {
				t.Errorf("Brew package name is empty for tool %s", tool.Name)
			}
		}
		if tool.VerifyCommand == "" {
			t.Errorf("Verify command is empty for tool %s", tool.Name)
		}
	}

	// Verify specific tools are present
	requiredTools := []string{"Git", "cURL", "Wget", "Build Essential"}
	for _, required := range requiredTools {
		found := false
		for _, tool := range tools {
			if tool.Name == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required tool %s not found in CoreTools()", required)
		}
	}
} 