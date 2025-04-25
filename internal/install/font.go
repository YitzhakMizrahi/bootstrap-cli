// Package install provides functionality for installing and managing various components
// in the bootstrap-cli, including tools, fonts, and system packages.
package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/interfaces"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/log"
)

// FontInstaller handles font installation
type FontInstaller struct {
	logger *log.Logger
}

// NewFontInstaller creates a new font installer
func NewFontInstaller(logger *log.Logger) *FontInstaller {
	return &FontInstaller{
		logger: logger,
	}
}

// InstallFont installs a font using its configuration
func (f *FontInstaller) InstallFont(font *interfaces.Font) error {
	fontDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "fonts")
	if err := os.MkdirAll(fontDir, 0755); err != nil {
		return fmt.Errorf("failed to create font directory: %w", err)
	}

	// Download the font
	f.logger.Info("Downloading %s...", font.Name)
	downloadPath := filepath.Join(fontDir, filepath.Base(font.Source))
	if err := exec.Command("curl", "-L", "-o", downloadPath, font.Source).Run(); err != nil {
		return fmt.Errorf("failed to download font: %w", err)
	}

	// Extract if it's a zip file
	if filepath.Ext(downloadPath) == ".zip" {
		f.logger.Info("Extracting font files...")
		if err := exec.Command("unzip", "-o", downloadPath, "-d", fontDir).Run(); err != nil {
			return fmt.Errorf("failed to extract font: %w", err)
		}

		// Clean up the zip file
		if err := os.Remove(downloadPath); err != nil {
			f.logger.Warn("Failed to remove font archive: %v", err)
		}
	}

	// Run any additional install commands
	for _, cmd := range font.GetInstallCommands() {
		f.logger.Info("Running install command: %s", cmd)
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return fmt.Errorf("failed to run install command: %w", err)
		}
	}

	// Update font cache
	f.logger.Info("Updating font cache...")
	if err := exec.Command("fc-cache", "-f").Run(); err != nil {
		return fmt.Errorf("failed to update font cache: %w", err)
	}

	// Run verification commands
	for _, cmd := range font.GetVerifyCommands() {
		f.logger.Info("Running verify command: %s", cmd)
		if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
			return fmt.Errorf("failed to verify font installation: %w", err)
		}
	}

	return nil
}

// InstallJetBrainsMono installs JetBrains Mono Nerd Font
func (f *FontInstaller) InstallJetBrainsMono() error {
	font := &interfaces.Font{
		Name:        "JetBrains Mono Nerd Font",
		Description: "A monospace font with programming ligatures and Nerd Font symbols",
		Source:      "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.zip",
		Install:     []string{},
		Verify:      []string{"fc-list | grep -i 'JetBrains Mono'"},
	}
	return f.InstallFont(font)
} 