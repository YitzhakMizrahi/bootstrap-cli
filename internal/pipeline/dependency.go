package pipeline

import (
	"fmt"
	"sort"
)

// DependencyType represents the type of dependency
type DependencyType string

const (
	// PackageDependency represents a package manager dependency
	PackageDependency DependencyType = "package"
	// SystemDependency represents a system-level dependency
	SystemDependency DependencyType = "system"
	// FileDependency represents a file dependency
	FileDependency DependencyType = "file"
)

// Dependency represents a dependency with its requirements
type Dependency struct {
	Name         string
	Type         DependencyType
	Version      string
	Optional     bool
	Platform     []string // Supported platforms (e.g., ["linux", "darwin"])
	Alternatives []string // Alternative dependencies that can satisfy this requirement
}

// DependencyGraph manages dependencies and their relationships
type DependencyGraph struct {
	dependencies map[string][]Dependency
	visited      map[string]bool
	temporary    map[string]bool
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		dependencies: make(map[string][]Dependency),
		visited:      make(map[string]bool),
		temporary:    make(map[string]bool),
	}
}

// AddDependency adds a dependency to the graph
func (g *DependencyGraph) AddDependency(name string, deps []Dependency) error {
	g.dependencies[name] = deps
	return nil
}

// GetInstallOrder returns the order in which dependencies should be installed
func (g *DependencyGraph) GetInstallOrder() ([]string, error) {
	var order []string
	g.visited = make(map[string]bool)
	g.temporary = make(map[string]bool)

	// Sort dependencies to ensure deterministic order
	var names []string
	for name := range g.dependencies {
		names = append(names, name)
	}
	sort.Strings(names)

	// Visit each node
	for _, name := range names {
		if !g.visited[name] {
			if err := g.visit(name, &order); err != nil {
				return nil, err
			}
		}
	}

	// The order is already correct (dependencies first, then dependents)
	// No need to reverse
	return order, nil
}

// visit performs a topological sort using depth-first search
func (g *DependencyGraph) visit(name string, order *[]string) error {
	// Check for cyclic dependencies
	if g.temporary[name] {
		return fmt.Errorf("cyclic dependency detected: %s", name)
	}
	if g.visited[name] {
		return nil
	}

	g.temporary[name] = true

	// Visit all dependencies
	for _, dep := range g.dependencies[name] {
		if !dep.Optional {
			if err := g.visit(dep.Name, order); err != nil {
				return err
			}
		}
	}

	g.visited[name] = true
	g.temporary[name] = false
	*order = append(*order, name)

	return nil
}

// ValidateForPlatform validates dependencies for the given platform
func (g *DependencyGraph) ValidateForPlatform(platform *Platform) error {
	for name, deps := range g.dependencies {
		for _, dep := range deps {
			if !dep.Optional && len(dep.Platform) > 0 {
				supported := false
				for _, p := range dep.Platform {
					if p == platform.OS {
						supported = true
						break
					}
				}
				if !supported {
					return fmt.Errorf("dependency %s is not supported on platform %s", dep.Name, platform.OS)
				}
			}
		}
		// Check if the package itself has platform requirements
		if deps, ok := g.dependencies[name]; ok {
			for _, dep := range deps {
				if len(dep.Platform) > 0 {
					supported := false
					for _, p := range dep.Platform {
						if p == platform.OS {
							supported = true
							break
						}
					}
					if !supported {
						return fmt.Errorf("package %s is not supported on platform %s", name, platform.OS)
					}
				}
			}
		}
	}
	return nil
}

// GetDependencies returns all dependencies for a given package
func (g *DependencyGraph) GetDependencies(name string) []Dependency {
	return g.dependencies[name]
}

// HasDependency checks if a package has a specific dependency
func (g *DependencyGraph) HasDependency(pkg, dep string) bool {
	for _, d := range g.dependencies[pkg] {
		if d.Name == dep {
			return true
		}
	}
	return false
}

// GetOptionalDependencies returns all optional dependencies
func (g *DependencyGraph) GetOptionalDependencies() map[string][]Dependency {
	optional := make(map[string][]Dependency)
	for name, deps := range g.dependencies {
		var optDeps []Dependency
		for _, dep := range deps {
			if dep.Optional {
				optDeps = append(optDeps, dep)
			}
		}
		if len(optDeps) > 0 {
			optional[name] = optDeps
		}
	}
	return optional
}

// GetRequiredDependencies returns all required dependencies
func (g *DependencyGraph) GetRequiredDependencies() map[string][]Dependency {
	required := make(map[string][]Dependency)
	for name, deps := range g.dependencies {
		var reqDeps []Dependency
		for _, dep := range deps {
			if !dep.Optional {
				reqDeps = append(reqDeps, dep)
			}
		}
		if len(reqDeps) > 0 {
			required[name] = reqDeps
		}
	}
	return required
} 