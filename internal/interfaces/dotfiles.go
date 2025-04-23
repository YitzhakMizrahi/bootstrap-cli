package interfaces

// DotfilesManager defines the interface for dotfiles management operations
type DotfilesManager interface {
	// Initialize sets up the dotfiles directory structure
	Initialize() error
	// CloneUserRepo clones the user's dotfiles repository
	CloneUserRepo(repoURL string) error
	// ApplyDotfile applies a dotfile configuration
	ApplyDotfile(config *Dotfile) error
	// BackupFile creates a backup of an existing file
	BackupFile(filePath string, suffix string) error
	// CreateSymlink creates a symlink from source to destination
	CreateSymlink(source, destination string) error
	// WriteContentFile writes content to a file
	WriteContentFile(content []byte, destination string) error
	// ApplyShellConfig applies shell-specific configuration
	ApplyShellConfig(config *ShellConfig) error
}

// DotfileOperation represents a dotfile operation type
type DotfileOperation string

const (
	// Create creates a new file
	Create DotfileOperation = "create"
	// Update updates an existing file
	Update DotfileOperation = "update"
	// Delete deletes a file
	Delete DotfileOperation = "delete"
	// Symlink creates a symlink
	Symlink DotfileOperation = "symlink"
)

// DotfileStatus represents the status of a dotfile
type DotfileStatus string

const (
	// Installed means the dotfile is installed
	Installed DotfileStatus = "installed"
	// NotInstalled means the dotfile is not installed
	NotInstalled DotfileStatus = "not_installed"
	// Outdated means the dotfile is installed but outdated
	Outdated DotfileStatus = "outdated"
	// Failed means the dotfile installation failed
	Failed DotfileStatus = "failed"
)

// Dotfile represents a configuration file
type Dotfile struct {
	Name            string   `yaml:"name"`
	Description     string   `yaml:"description"`
	Category        string   `yaml:"category"`
	Tags            []string `yaml:"tags"`
	Files           []DotfileFile `yaml:"files"`
	Dependencies    []string `yaml:"dependencies"`
	ShellConfig     ShellConfig `yaml:"shell_config"`
	PostInstall     []string `yaml:"post_install"`
	RequiresRestart bool     `yaml:"requires_restart"`
	// Fields for centralized management
	SourceRepo      string   `yaml:"source_repo"` // Optional: GitHub repo URL for user's dotfiles
	BaseDir         string   `yaml:"base_dir"` // Base directory for dotfiles (default: ~/.dotfiles)
	SymlinkStrategy SymlinkStrategy `yaml:"symlink_strategy"`
}

// DotfileFile represents a file to be managed
type DotfileFile struct {
	// Source is the source path of the file
	Source string `yaml:"source"`
	// Destination is the destination path of the file
	Destination string `yaml:"destination"`
	// Operation is the operation to perform on the file
	Operation DotfileOperation `yaml:"operation"`
	// Backup determines if a backup should be created
	Backup bool `yaml:"backup"`
	// BackupSuffix is the suffix to use for the backup file
	BackupSuffix string `yaml:"backup_suffix"`
	// Content is the content to write to the file (for Create/Update operations)
	Content string `yaml:"content"`
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