package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/config"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/install"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// ToolCategory represents a category of tools
type ToolCategory struct {
	Name        string
	Description string
	Tools       []*install.Tool
}

// GetToolCategories returns all available tool categories
func GetToolCategories() ([]ToolCategory, error) {
	// Create config loader with the temp directory
	loader := config.NewConfigLoader("")

	// Load all tools
	tools, err := loader.LoadTools()
	if err != nil {
		return nil, fmt.Errorf("error loading tools: %w", err)
	}

	// Group tools by category
	categories := make(map[string][]*install.Tool)
	for _, tool := range tools {
		categories[tool.Category] = append(categories[tool.Category], tool)
	}

	// Create tool categories
	result := []ToolCategory{
		{
			Name:        "Essential Tools",
			Description: "Core development tools required for most workflows",
			Tools:       categories["essential"],
		},
		{
			Name:        "Modern CLI Tools",
			Description: "Modern alternatives to traditional command-line tools",
			Tools:       categories["modern"],
		},
		{
			Name:        "System Tools",
			Description: "System monitoring and management utilities",
			Tools:       categories["system"],
		},
	}

	return result, nil
}

// CoreTool represents a core development tool
type CoreTool struct {
	*install.Tool
	// VerifyCommand is the command to verify the tool is installed correctly
	VerifyCommand string
	// PostInstallCommands are commands to run after installation
	PostInstallCommands []string
}

// CoreTools returns a list of core development tools
func CoreTools() []*CoreTool {
	return []*CoreTool{
		{
			Tool: &install.Tool{
				Name:        "Git",
				PackageName: "git",
				Description: "Version control system",
			},
			VerifyCommand: "git --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Curl",
				PackageName: "curl",
				Description: "Command-line tool for transferring data",
			},
			VerifyCommand: "curl --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Wget",
				PackageName: "wget",
				Description: "Command-line utility for downloading files",
			},
			VerifyCommand: "wget --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Build Essentials",
				PackageName: "build-essential",
				Description: "Basic build tools and libraries",
				PackageNames: &install.PackageMapping{
					Default: "build-essential",
					APT:     "build-essential",
					DNF:     "gcc-c++ make",
					Pacman:  "base-devel",
				},
			},
			VerifyCommand: "gcc --version",
		},
		{
			Tool: &install.Tool{
				Name:        "ZIP",
				PackageName: "zip",
				Description: "Compression and file packaging utility",
			},
			VerifyCommand: "zip --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Unzip",
				PackageName: "unzip",
				Description: "Decompression utility",
			},
			VerifyCommand: "unzip --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Tar",
				PackageName: "tar",
				Description: "Tape archiver",
			},
			VerifyCommand: "tar --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Vim",
				PackageName: "vim",
				Description: "Text editor",
			},
			VerifyCommand: "vim --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Nano",
				PackageName: "nano",
				Description: "Simple text editor",
			},
			VerifyCommand: "nano --version",
		},
		{
			Tool: &install.Tool{
				Name:        "Htop",
				PackageName: "htop",
				Description: "Interactive process viewer",
			},
			VerifyCommand: "htop --version",
		},
	}
}

// InstallEssentialTools installs all essential development tools
func InstallEssentialTools(pm interfaces.PackageManager, logger *log.Logger, skipVerification bool) error {
	logger.Info("Installing essential development tools...")

	tools := CoreTools()
	failed := false

	for _, tool := range tools {
		logger.Info("Installing %s...", tool.Name)
		
		// Check if tool is already installed
		installed := pm.IsInstalled(tool.PackageName)
		if installed {
			logger.Info("%s is already installed", tool.Name)
			continue
		}
		
		// Install the tool
		if err := pm.Install(tool.PackageName); err != nil {
			logger.Error("Failed to install %s: %v", tool.Name, err)
			failed = true
			continue
		}
		
		// Run post-install commands if any
		if len(tool.PostInstallCommands) > 0 {
			logger.Debug("Running post-install commands for %s", tool.Name)
			for _, cmd := range tool.PostInstallCommands {
				if err := runCommand(cmd); err != nil {
					logger.Error("Failed to run post-install command for %s: %v", tool.Name, err)
					failed = true
					break
				}
			}
		}
		
		// Verify installation
		if !skipVerification && tool.VerifyCommand != "" {
			logger.Debug("Verifying %s installation", tool.Name)
			if err := runCommand(tool.VerifyCommand); err != nil {
				logger.Error("Failed to verify %s installation: %v", tool.Name, err)
				failed = true
				continue
			}
		}
		
		logger.Info("Successfully installed %s", tool.Name)
	}

	if failed {
		return fmt.Errorf("one or more essential tools failed to install")
	}

	logger.Info("All essential development tools installed successfully")
	return nil
}

// runCommand executes a shell command
func runCommand(cmd string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	// Create a command with the parts
	command := exec.Command(parts[0], parts[1:]...)
	
	// Capture output
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	
	// Run the command
	err := command.Run()
	if err != nil {
		return fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}
	
	return nil
}

// InstallOptions represents options for tool installation
type InstallOptions struct {
	Logger           *log.Logger
	PackageManager   interfaces.PackageManager
	Tools           []*install.Tool
	SkipVerification bool
	AdditionalPaths []string
}

var selectedTools []*install.Tool

// GetSelectedTools returns the list of tools selected during initialization
func GetSelectedTools() []*install.Tool {
	return selectedTools
}

// SetSelectedTools sets the list of tools to be installed
func SetSelectedTools(tools []*install.Tool) {
	selectedTools = tools
}

// InstallSelectedTools installs a set of selected development tools
func InstallSelectedTools(opts *InstallOptions) error {
	if opts.Logger == nil {
		opts.Logger = log.New(log.InfoLevel)
	}

	for _, tool := range opts.Tools {
		opts.Logger.Info("Installing %s...", tool.Name)
		
		// Create installer for the tool
		installer := install.NewInstaller(opts.PackageManager)
		
		// Install the tool
		if err := installer.Install(tool); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool.Name, err)
		}

		// Verify installation if not skipped
		if !opts.SkipVerification {
			if err := VerifyTool(tool, opts.AdditionalPaths); err != nil {
				return fmt.Errorf("failed to verify %s: %w", tool.Name, err)
			}
		}
	}

	return nil
} 