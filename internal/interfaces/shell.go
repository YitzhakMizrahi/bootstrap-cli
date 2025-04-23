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

// ShellType represents supported shell types
type ShellType string

const (
	// Shell types
	BashShell ShellType = "bash"
	ZshShell  ShellType = "zsh"
	FishShell ShellType = "fish"
)

// IsValidShell checks if a shell type is supported
func IsValidShell(shell string) bool {
	switch ShellType(shell) {
	case BashShell, ZshShell, FishShell:
		return true
	default:
		return false
	}
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
	// ConfigureShell configures a shell with the specified type
	ConfigureShell(shellType string) error
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

// Dotfile represents a shell configuration file
type Dotfile struct {
	Name            string
	Description     string
	Category        string
	Tags            []string
	Files           []FileConfig
	Dependencies    []string
	ShellConfig     ShellConfig
	PostInstall     []string
	RequiresRestart bool
	// New fields for centralized management
	SourceRepo      string   // Optional: GitHub repo URL for user's dotfiles
	BaseDir         string   // Base directory for dotfiles (default: ~/.dotfiles)
	SymlinkStrategy SymlinkStrategy
}

// SymlinkStrategy defines how to handle dotfile symlinks
type SymlinkStrategy int

const (
	// SymlinkToHome creates symlinks from ~/.dotfiles to ~/
	SymlinkToHome SymlinkStrategy = iota
	// SymlinkToDotfiles creates symlinks from user's repo to ~/.dotfiles
	SymlinkToDotfiles
	// CopyToHome copies files directly to home directory
	CopyToHome
)

// FileConfig represents a file configuration
type FileConfig struct {
	Source         string
	Destination    string
	Type           string
	Permissions    string
	Backup         bool
	BackupSuffix   string
	SymlinkTarget  string // The target path for symlinks (relative to BaseDir)
	CreateParents  bool   // Whether to create parent directories
}

// ShellConfig represents shell-specific configuration
type ShellConfig struct {
	Env           map[string]string
	PathAdditions []string
	InitScripts   []string // Scripts to run on shell initialization
} 