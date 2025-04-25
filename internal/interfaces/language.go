package interfaces

// Language represents a programming language runtime
type Language struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	Version     string   `yaml:"version"`
	Installer   string   `yaml:"installer"`
	VerifyCommand string `yaml:"verify_command"`

	// Dependencies required for installation
	Dependencies []struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Optional bool   `yaml:"optional"`
	} `yaml:"dependencies"`

	// System level dependencies
	SystemDependencies []string `yaml:"system_dependencies"`

	// Package management
	PackageNames struct {
		APT    string `yaml:"apt"`
		Brew   string `yaml:"brew"`
		DNF    string `yaml:"dnf"`
		Pacman string `yaml:"pacman"`
	} `yaml:"package_names"`

	// Post-installation steps
	PostInstall []struct {
		Command     string `yaml:"command"`
		Description string `yaml:"description"`
	} `yaml:"post_install"`

	// Shell configuration
	ShellConfig struct {
		Env    map[string]string `yaml:"env"`
		Source []string         `yaml:"source"`
	} `yaml:"shell_config"`
}

// GetPackageName returns the package name for the given package manager
func (l *Language) GetPackageName(packageManager string) string {
	switch packageManager {
	case "apt":
		return l.PackageNames.APT
	case "brew":
		return l.PackageNames.Brew
	case "dnf":
		return l.PackageNames.DNF
	case "pacman":
		return l.PackageNames.Pacman
	default:
		return ""
	}
}

// GetInstaller returns the language version manager to use
func (l *Language) GetInstaller() string {
	return l.Installer
}

// GetVersion returns the desired version of the language
func (l *Language) GetVersion() string {
	return l.Version
}

// ToTool converts a Language to a Tool for installation
func (l *Language) ToTool() *Tool {
	// Convert dependencies
	deps := make([]struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Optional bool   `yaml:"optional,omitempty"`
	}, len(l.Dependencies))
	
	for i, dep := range l.Dependencies {
		deps[i] = struct {
			Name     string `yaml:"name"`
			Type     string `yaml:"type"`
			Optional bool   `yaml:"optional,omitempty"`
		}{
			Name:     dep.Name,
			Type:     dep.Type,
			Optional: dep.Optional,
		}
	}

	return &Tool{
		Name:               l.Name,
		Description:        l.Description,
		Category:          l.Category,
		Tags:              l.Tags,
		Version:           l.Version,
		Dependencies:      deps,
		SystemDependencies: l.SystemDependencies,
		PackageNames:      l.PackageNames,
		VerifyCommand:     l.VerifyCommand,
		PostInstall:       l.PostInstall,
		ShellConfig: struct {
			Aliases   map[string]string `yaml:"aliases,omitempty"`
			Env       map[string]string `yaml:"env,omitempty"`
			Path      []string         `yaml:"path,omitempty"`
			Functions map[string]string `yaml:"functions,omitempty"`
		}{
			Env: l.ShellConfig.Env,
		},
	}
} 