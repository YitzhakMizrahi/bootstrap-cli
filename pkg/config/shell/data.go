package shell

// Data represents shell configuration data
type Data struct {
	HomeDir  string
	Path     []string
	Aliases  map[string]string
	EnvVars  map[string]string
	Plugins  []string
	Custom   string
}

// New creates a new configuration data instance
func New(homeDir string) *Data {
	return &Data{
		HomeDir:  homeDir,
		Path:     make([]string, 0),
		Aliases:  make(map[string]string),
		EnvVars:  make(map[string]string),
		Plugins:  make([]string, 0),
		Custom:   "",
	}
} 