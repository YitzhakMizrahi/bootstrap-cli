package interfaces

import "errors"

// Shell represents a shell that can be installed and set as default.
type Shell struct {
	Name             string `yaml:"name"`
	Description      string `yaml:"description"`
	InstallCommands  struct {
		Apt    string `yaml:"apt,omitempty"`
		Brew   string `yaml:"brew,omitempty"`
		Dnf    string `yaml:"dnf,omitempty"`
		Pacman string `yaml:"pacman,omitempty"`
	} `yaml:"install_commands"`
	Path            string `yaml:"path"`
	SetDefaultCommand string `yaml:"set_default_command,omitempty"`
	VerifyCommand   string `yaml:"verify_command,omitempty"`
}

// ShellType represents a shell type
type ShellType string

const (
	// BashShell represents the Bash shell
	BashShell ShellType = "bash"
	// ZshShell represents the Zsh shell
	ZshShell ShellType = "zsh"
	// FishShell represents the Fish shell
	FishShell ShellType = "fish"
)

// Error variables
var (
	ErrHomeDirNotFound   = errors.New("home directory not found")
	ErrUnsupportedShell  = errors.New("unsupported shell type")
)

// ShellInfo contains information about a shell
type ShellInfo struct {
	Current     string   // Current shell
	Available   []string // Available shells on the system
	DefaultPath string   // Default shell path
	Type        string   // Shell type (bash, zsh, fish)
	Path        string   // Full path to the shell executable
	Version     string   // Shell version
	IsDefault   bool     // Whether this is the default shell
	IsAvailable bool     // Whether this shell is available on the system
	ConfigFiles []string // Configuration files for this shell
}

// ShellManager defines the interface for shell management operations
type ShellManager interface {
	// DetectCurrent detects the current user's shell
	DetectCurrent() (*ShellInfo, error)
	// ListAvailable returns a list of available shells
	ListAvailable() ([]*ShellInfo, error)
	// IsInstalled checks if a specific shell is installed
	IsInstalled(shell ShellType) bool
	// GetInfo returns detailed information about a specific shell
	GetInfo(shell ShellType) (*ShellInfo, error)
	// ConfigureShell configures a shell with the specified configuration
	ConfigureShell(config *ShellConfig) error
}

// ShellConfig represents shell configuration
type ShellConfig struct {
	// Aliases are shell aliases to add
	Aliases map[string]string `yaml:"aliases"`
	// Exports are environment variables to export
	Exports map[string]string `yaml:"exports"`
	// Functions are shell functions to add
	Functions map[string]string `yaml:"functions"`
	// Path contains paths to add to PATH
	Path []string `yaml:"path"`
	// Source contains files to source
	Source []string `yaml:"source"`
}

// FileConfig represents a file configuration
type FileConfig struct {
	// Source is the source path of the file
	Source string `yaml:"source"`
	// Destination is the destination path of the file
	Destination string `yaml:"destination"`
	// Content is the content of the file
	Content string `yaml:"content"`
}

// IsValidShell checks if a shell type is supported
func IsValidShell(shell string) bool {
	switch ShellType(shell) {
	case BashShell, ZshShell, FishShell:
		return true
	default:
		return false
	}
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