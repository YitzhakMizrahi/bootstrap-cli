package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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

// InstallJetBrainsMono installs JetBrains Mono Nerd Font
func (f *FontInstaller) InstallJetBrainsMono() error {
	fontDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "fonts")
	if err := os.MkdirAll(fontDir, 0755); err != nil {
		return fmt.Errorf("failed to create font directory: %w", err)
	}

	// Download JetBrains Mono Nerd Font
	downloadURL := "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.zip"
	downloadPath := filepath.Join(fontDir, "JetBrainsMono.zip")

	f.logger.Info("Downloading JetBrains Mono Nerd Font...")
	if err := exec.Command("curl", "-L", "-o", downloadPath, downloadURL).Run(); err != nil {
		return fmt.Errorf("failed to download font: %w", err)
	}

	// Extract the font
	f.logger.Info("Extracting font files...")
	if err := exec.Command("unzip", "-o", downloadPath, "-d", fontDir).Run(); err != nil {
		return fmt.Errorf("failed to extract font: %w", err)
	}

	// Clean up the zip file
	if err := os.Remove(downloadPath); err != nil {
		f.logger.Warn("Failed to remove font archive: %v", err)
	}

	// Update font cache
	f.logger.Info("Updating font cache...")
	if err := exec.Command("fc-cache", "-f").Run(); err != nil {
		return fmt.Errorf("failed to update font cache: %w", err)
	}

	return nil
} 