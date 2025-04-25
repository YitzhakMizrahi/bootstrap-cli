package interfaces

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

// GetName returns the font name
func (f *Font) GetName() string {
	return f.Name
}

// GetSource returns the font source URL
func (f *Font) GetSource() string {
	return f.Source
}

// GetInstallCommands returns the installation commands
func (f *Font) GetInstallCommands() []string {
	return f.Install
}

// GetVerifyCommands returns the verification commands
func (f *Font) GetVerifyCommands() []string {
	return f.Verify
} 