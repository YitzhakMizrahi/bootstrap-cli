package interfaces

// Shell represents a shell type
type Shell string

const (
	// Supported shell types
	Bash Shell = "bash"
	Zsh  Shell = "zsh"
	Fish Shell = "fish"
)

// ShellInfo contains information about a shell
type ShellInfo struct {
	Type        Shell
	Path        string
	Version     string
	IsDefault   bool
	IsAvailable bool
	ConfigFiles []string
}

// ShellManager handles shell detection and operations
type ShellManager interface {
	// DetectCurrent detects the current user's shell
	DetectCurrent() (*ShellInfo, error)
	// ListAvailable returns a list of available shells
	ListAvailable() ([]*ShellInfo, error)
	// IsInstalled checks if a specific shell is installed
	IsInstalled(shell Shell) bool
	// GetInfo returns detailed information about a specific shell
	GetInfo(shell Shell) (*ShellInfo, error)
}

// DotfilesStrategy defines how to handle existing dotfiles
type DotfilesStrategy int

const (
	// MergeWithExisting merges new configurations with existing ones
	MergeWithExisting DotfilesStrategy = iota
	// SkipIfExists skips adding configurations if they already exist
	SkipIfExists
	// ReplaceExisting replaces existing configurations with new ones
	ReplaceExisting
)

// ShellConfigWriter handles shell configuration file management
type ShellConfigWriter interface {
	// WriteConfig writes shell configurations to the appropriate file
	WriteConfig(configs []string, strategy DotfilesStrategy) error
	// AddToPath adds a directory to the PATH environment variable
	AddToPath(path string) error
	// SetEnvVar sets an environment variable
	SetEnvVar(name, value string) error
	// AddAlias adds a shell alias
	AddAlias(name, command string) error
	// HasConfig checks if a configuration exists
	HasConfig(config string) bool
} 