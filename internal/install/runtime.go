package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/packages"
)

// configureNeedrestart sets needrestart mode (can be 'a' for automatic or 'i' for interactive)
func configureNeedrestart(mode string) error {
	cmd := exec.Command("sudo", "sed", "-i", 
		fmt.Sprintf("s/^#\\$nrconf{restart} = 'i';/\\$nrconf{restart} = '%s';/", mode),
		"/etc/needrestart/needrestart.conf")
	return cmd.Run()
}

// RuntimeInstaller handles language runtime installation
type RuntimeInstaller struct {
	pm     packages.PackageManager
	logger *log.Logger
}

// NewRuntimeInstaller creates a new runtime installer
func NewRuntimeInstaller(pm packages.PackageManager, logger *log.Logger) *RuntimeInstaller {
	return &RuntimeInstaller{
		pm:     pm,
		logger: logger,
	}
}

// Install installs a language runtime
func (r *RuntimeInstaller) Install(runtime string) error {
	// Configure needrestart to automatic mode
	if err := configureNeedrestart("a"); err != nil {
		r.logger.Warn("Failed to configure needrestart: %v", err)
	}
	
	// Defer resetting needrestart to interactive mode
	defer func() {
		if err := configureNeedrestart("i"); err != nil {
			r.logger.Warn("Failed to reset needrestart: %v", err)
		}
	}()

	switch runtime {
	case "Node.js (nvm)":
		return r.installNVM()
	case "Python (pyenv)":
		return r.installPyenv()
	case "Go (goenv)":
		return r.installGoenv()
	case "Rust (rustup)":
		return r.installRustup()
	default:
		return fmt.Errorf("unknown runtime: %s", runtime)
	}
}

func (r *RuntimeInstaller) installNVM() error {
	r.logger.Info("Installing NVM (Node Version Manager)...")
	
	// Download and run the NVM install script
	cmd := exec.Command("bash", "-c", `curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install NVM: %w", err)
	}

	// Add NVM to shell configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nvmInit := `
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion
`

	// Append to .bashrc and .zshrc if they exist
	for _, rc := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rc)
		if _, err := os.Stat(rcPath); err == nil {
			if err := appendToFile(rcPath, nvmInit); err != nil {
				r.logger.Warn("Failed to update %s: %v", rc, err)
			}
		}
	}

	return nil
}

func (r *RuntimeInstaller) installPyenv() error {
	r.logger.Info("Installing pyenv...")

	// Install all pyenv dependencies in a single command
	deps := []string{
		"make", "build-essential", "libssl-dev", "zlib1g-dev",
		"libbz2-dev", "libreadline-dev", "libsqlite3-dev", "wget",
		"curl", "llvm", "libncursesw5-dev", "xz-utils", "tk-dev",
		"libxml2-dev", "libxmlsec1-dev", "libffi-dev", "liblzma-dev",
	}

	// Join all dependencies into a single installation command
	if err := r.pm.Install(strings.Join(deps, " ")); err != nil {
		return fmt.Errorf("failed to install pyenv dependencies: %w", err)
	}

	// Clone pyenv
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	pyenvPath := filepath.Join(homeDir, ".pyenv")
	if err := exec.Command("git", "clone", "https://github.com/pyenv/pyenv.git", pyenvPath).Run(); err != nil {
		return fmt.Errorf("failed to clone pyenv: %w", err)
	}

	// Add pyenv to shell configuration
	pyenvInit := `
export PYENV_ROOT="$HOME/.pyenv"
command -v pyenv >/dev/null || export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
`

	for _, rc := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rc)
		if _, err := os.Stat(rcPath); err == nil {
			if err := appendToFile(rcPath, pyenvInit); err != nil {
				r.logger.Warn("Failed to update %s: %v", rc, err)
			}
		}
	}

	return nil
}

func (r *RuntimeInstaller) installGoenv() error {
	r.logger.Info("Installing goenv...")

	// Clone goenv
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	goenvPath := filepath.Join(homeDir, ".goenv")
	if err := exec.Command("git", "clone", "https://github.com/syndbg/goenv.git", goenvPath).Run(); err != nil {
		return fmt.Errorf("failed to clone goenv: %w", err)
	}

	// Add goenv to shell configuration
	goenvInit := `
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
`

	for _, rc := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rc)
		if _, err := os.Stat(rcPath); err == nil {
			if err := appendToFile(rcPath, goenvInit); err != nil {
				r.logger.Warn("Failed to update %s: %v", rc, err)
			}
		}
	}

	return nil
}

func (r *RuntimeInstaller) installRustup() error {
	r.logger.Info("Installing Rustup...")

	// Download and run the rustup install script
	cmd := exec.Command("bash", "-c", `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Rustup: %w", err)
	}

	// Add Cargo to shell configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	cargoInit := `
export PATH="$HOME/.cargo/bin:$PATH"
. "$HOME/.cargo/env"
`

	for _, rc := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rc)
		if _, err := os.Stat(rcPath); err == nil {
			if err := appendToFile(rcPath, cargoInit); err != nil {
				r.logger.Warn("Failed to update %s: %v", rc, err)
			}
		}
	}

	return nil
}

func appendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Check if content already exists
	existing, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if !strings.Contains(string(existing), content) {
		if _, err := f.WriteString(content); err != nil {
			return err
		}
	}

	return nil
} 