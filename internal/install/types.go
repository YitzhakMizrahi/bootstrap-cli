package install

// ToolCategory represents a category of tools
type ToolCategory struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Priority    int      `yaml:"priority"`
	Tools       []string `yaml:"tools"`
} 