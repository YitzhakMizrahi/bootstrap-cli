package interfaces

// LanguageManager represents a language version manager
type LanguageManager struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Languages   []string `yaml:"languages"`  // List of supported languages
	Version     string   `yaml:"version"`
	
	// Package management
	PackageName string `yaml:"package_name"` // Package name for the default package manager
	PackageNames struct {
		APT    string `yaml:"apt"`
		Brew   string `yaml:"brew"`
		DNF    string `yaml:"dnf"`
		Pacman string `yaml:"pacman"`
	} `yaml:"package_names"`

	// Installation details
	Dependencies []string `yaml:"dependencies"`
	PostInstall []struct {
		Command     string `yaml:"command"`
		Description string `yaml:"description"`
	} `yaml:"post_install"`

	// Shell configuration
	ShellConfig struct {
		Exports   map[string]string `yaml:"exports"`
		Path      []string         `yaml:"path"`
		Functions map[string]string `yaml:"functions"`
	} `yaml:"shell_config"`
}

// SupportsLanguage checks if this manager supports the given language
func (m *LanguageManager) SupportsLanguage(lang string) bool {
	for _, supported := range m.Languages {
		if supported == lang {
			return true
		}
	}
	return false
} 