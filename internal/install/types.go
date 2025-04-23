package install

// Font represents a font configuration
type Font struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	Source      string   `yaml:"source"`
	Install     []string `yaml:"install"`
	Verify      []string `yaml:"verify"`
}

// Language represents a programming language configuration
type Language struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	Version     string   `yaml:"version"`
	Install     []string `yaml:"install"`
	Verify      []string `yaml:"verify"`
	Env         map[string]string `yaml:"env"`
} 