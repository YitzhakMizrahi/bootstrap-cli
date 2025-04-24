package pipeline

import (
	"testing"
)

func TestDependencyGraph_AddDependency(t *testing.T) {
	graph := NewDependencyGraph()

	// Test adding a simple dependency
	deps := []Dependency{
		{
			Name:     "unzip",
			Type:     SystemDependency,
			Optional: false,
		},
	}
	err := graph.AddDependency("font-installer", deps)
	if err != nil {
		t.Errorf("Failed to add dependency: %v", err)
	}

	// Verify dependency was added
	if !graph.HasDependency("font-installer", "unzip") {
		t.Error("Expected dependency 'unzip' not found")
	}
}

func TestDependencyGraph_GetInstallOrder(t *testing.T) {
	graph := NewDependencyGraph()

	// Create a dependency chain: C -> B -> A
	graph.AddDependency("A", nil)
	graph.AddDependency("B", []Dependency{{Name: "A", Type: PackageDependency}})
	graph.AddDependency("C", []Dependency{{Name: "B", Type: PackageDependency}})

	order, err := graph.GetInstallOrder()
	if err != nil {
		t.Errorf("Failed to get install order: %v", err)
	}

	// Check order (should be A, B, C)
	expected := []string{"A", "B", "C"}
	if len(order) != len(expected) {
		t.Errorf("Expected %d items in order, got %d", len(expected), len(order))
	}
	for i, pkg := range expected {
		if order[i] != pkg {
			t.Errorf("Expected package %s at position %d, got %s", pkg, i, order[i])
		}
	}
}

func TestDependencyGraph_CyclicDependencies(t *testing.T) {
	graph := NewDependencyGraph()

	// Create a cyclic dependency: A -> B -> A
	graph.AddDependency("A", []Dependency{{Name: "B", Type: PackageDependency}})
	graph.AddDependency("B", []Dependency{{Name: "A", Type: PackageDependency}})

	_, err := graph.GetInstallOrder()
	if err == nil {
		t.Error("Expected error for cyclic dependency, got nil")
	}
}

func TestDependencyGraph_PlatformValidation(t *testing.T) {
	graph := NewDependencyGraph()
	platform := &Platform{OS: "linux"}

	// Add platform-specific dependency
	graph.AddDependency("test-pkg", []Dependency{
		{
			Name:     "linux-only",
			Platform: []string{"linux"},
		},
	})

	// Test valid platform
	err := graph.ValidateForPlatform(platform)
	if err != nil {
		t.Errorf("Validation failed for supported platform: %v", err)
	}

	// Test unsupported platform
	platform.OS = "windows"
	err = graph.ValidateForPlatform(platform)
	if err == nil {
		t.Error("Expected error for unsupported platform, got nil")
	}
}

func TestDependencyGraph_OptionalDependencies(t *testing.T) {
	graph := NewDependencyGraph()

	// Add package with optional and required dependencies
	graph.AddDependency("test-pkg", []Dependency{
		{
			Name:     "required-dep",
			Optional: false,
		},
		{
			Name:     "optional-dep",
			Optional: true,
		},
	})

	// Test getting optional dependencies
	optDeps := graph.GetOptionalDependencies()
	if len(optDeps["test-pkg"]) != 1 {
		t.Errorf("Expected 1 optional dependency, got %d", len(optDeps["test-pkg"]))
	}
	if optDeps["test-pkg"][0].Name != "optional-dep" {
		t.Errorf("Expected optional dependency 'optional-dep', got '%s'", optDeps["test-pkg"][0].Name)
	}

	// Test getting required dependencies
	reqDeps := graph.GetRequiredDependencies()
	if len(reqDeps["test-pkg"]) != 1 {
		t.Errorf("Expected 1 required dependency, got %d", len(reqDeps["test-pkg"]))
	}
	if reqDeps["test-pkg"][0].Name != "required-dep" {
		t.Errorf("Expected required dependency 'required-dep', got '%s'", reqDeps["test-pkg"][0].Name)
	}
}

func TestDependencyGraph_AlternativeDependencies(t *testing.T) {
	graph := NewDependencyGraph()

	// Add package with alternative dependencies
	graph.AddDependency("test-pkg", []Dependency{
		{
			Name:         "primary-dep",
			Alternatives: []string{"alt-dep1", "alt-dep2"},
		},
	})

	deps := graph.GetDependencies("test-pkg")
	if len(deps) != 1 {
		t.Errorf("Expected 1 dependency, got %d", len(deps))
	}
	if len(deps[0].Alternatives) != 2 {
		t.Errorf("Expected 2 alternatives, got %d", len(deps[0].Alternatives))
	}
	if deps[0].Alternatives[0] != "alt-dep1" || deps[0].Alternatives[1] != "alt-dep2" {
		t.Error("Alternative dependencies not stored correctly")
	}
} 