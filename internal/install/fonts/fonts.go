package fonts

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FontManager handles font installation and management
type FontManager struct {
	InstallPath string
	FontsDir    string
	UserFonts   bool
}

// NewFontManager creates a new FontManager
func NewFontManager(installPath string, userFonts bool) *FontManager {
	fontsDir := filepath.Join(installPath, "fonts")
	if userFonts {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			fontsDir = filepath.Join(homeDir, ".local", "share", "fonts")
		}
	}

	return &FontManager{
		InstallPath: installPath,
		FontsDir:    fontsDir,
		UserFonts:   userFonts,
	}
}

// InstallNerdFonts installs Nerd Fonts
func (f *FontManager) InstallNerdFonts(fonts []string) error {
	// Create fonts directory if it doesn't exist
	if err := os.MkdirAll(f.FontsDir, 0755); err != nil {
		return fmt.Errorf("failed to create fonts directory: %w", err)
	}

	// Install each selected font
	for _, font := range fonts {
		if err := f.installNerdFont(font); err != nil {
			return fmt.Errorf("failed to install %s: %w", font, err)
		}
	}

	// Update font cache
	if err := f.updateFontCache(); err != nil {
		return fmt.Errorf("failed to update font cache: %w", err)
	}

	return nil
}

// installNerdFont installs a specific Nerd Font
func (f *FontManager) installNerdFont(fontName string) error {
	// Check if the font is already installed
	if f.IsFontInstalled(fontName) {
		fmt.Printf("Font %s is already installed, skipping...\n", fontName)
		return nil
	}

	// Check if the font is available
	if !f.IsNerdFontAvailable(fontName) {
		return fmt.Errorf("font %s is not available", fontName)
	}

	// Construct the download URL
	url := fmt.Sprintf("https://github.com/ryanoasis/nerd-fonts/releases/download/v3.1.1/%s.zip", fontName)

	// Create a temporary file for the download
	tempFile, err := os.CreateTemp("", "nerd-font-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Download the font
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download font: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download font: HTTP status %d", resp.StatusCode)
	}

	// Write the downloaded file
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write downloaded font: %w", err)
	}

	// Close the file before unzipping
	tempFile.Close()

	// Create a temporary directory for extraction
	tempDir, err := os.MkdirTemp("", "nerd-font-extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Unzip the font
	cmd := exec.Command("unzip", "-q", tempFile.Name(), "-d", tempDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to unzip font: %w\nOutput: %s", err, string(output))
	}

	// Find the font files (usually .ttf or .otf)
	fontFiles, err := filepath.Glob(filepath.Join(tempDir, "*.ttf"))
	if err != nil {
		return fmt.Errorf("failed to find font files: %w", err)
	}

	// If no .ttf files, try .otf
	if len(fontFiles) == 0 {
		fontFiles, err = filepath.Glob(filepath.Join(tempDir, "*.otf"))
		if err != nil {
			return fmt.Errorf("failed to find font files: %w", err)
		}
	}

	// Copy each font file to the fonts directory
	for _, fontFile := range fontFiles {
		destFile := filepath.Join(f.FontsDir, filepath.Base(fontFile))
		if err := copyFile(fontFile, destFile); err != nil {
			return fmt.Errorf("failed to copy font file %s: %w", fontFile, err)
		}
	}

	fmt.Printf("Successfully installed %s\n", fontName)
	return nil
}

// updateFontCache updates the system font cache
func (f *FontManager) updateFontCache() error {
	// Check if fc-cache is available
	cmd := exec.Command("which", "fc-cache")
	if err := cmd.Run(); err != nil {
		// fc-cache not found, try alternative methods
		return f.updateFontCacheAlternative()
	}

	// Use fc-cache to update the font cache
	cmd = exec.Command("fc-cache", "-f", "-v")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to update font cache: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// updateFontCacheAlternative tries alternative methods to update the font cache
func (f *FontManager) updateFontCacheAlternative() error {
	// Try different commands based on the OS
	commands := [][]string{
		{"mkfontdir", f.FontsDir},
		{"mkfontscale", f.FontsDir},
		{"xset", "+fp", f.FontsDir},
		{"xset", "fp", "reload"},
	}

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		if err := cmd.Run(); err != nil {
			// Just log the error and continue with the next command
			fmt.Printf("Warning: Failed to run %s: %v\n", cmdArgs[0], err)
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// ListAvailableNerdFonts returns a list of available Nerd Fonts
func (f *FontManager) ListAvailableNerdFonts() []string {
	return []string{
		"JetBrainsMono",
		"FiraCode",
		"Hack",
		"SourceCodePro",
		"CascadiaCode",
		"DejaVuSansMono",
		"DroidSansMono",
		"IBMPlexMono",
		"Inconsolata",
		"Meslo",
		"Monoid",
		"Mononoki",
		"Noto",
		"ProFont",
		"RobotoMono",
		"SpaceMono",
		"Terminus",
		"UbuntuMono",
		"VictorMono",
	}
}

// IsNerdFontAvailable checks if a Nerd Font is available for installation
func (f *FontManager) IsNerdFontAvailable(fontName string) bool {
	availableFonts := f.ListAvailableNerdFonts()
	for _, font := range availableFonts {
		if font == fontName {
			return true
		}
	}
	return false
}

// IsFontInstalled checks if a font is already installed
func (f *FontManager) IsFontInstalled(fontName string) bool {
	// Check if the font directory exists
	fontDir := filepath.Join(f.FontsDir, fontName)
	if _, err := os.Stat(fontDir); err == nil {
		return true
	}

	// Check for font files with the font name
	pattern := filepath.Join(f.FontsDir, "*"+fontName+"*.ttf")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		// Try .otf files
		pattern = filepath.Join(f.FontsDir, "*"+fontName+"*.otf")
		matches, err = filepath.Glob(pattern)
		if err != nil || len(matches) == 0 {
			return false
		}
	}

	return true
}

// GetInstalledFonts returns a list of installed fonts
func (f *FontManager) GetInstalledFonts() ([]string, error) {
	var fonts []string

	// Check if the fonts directory exists
	if _, err := os.Stat(f.FontsDir); os.IsNotExist(err) {
		return fonts, nil
	}

	// Get all .ttf and .otf files
	ttfFiles, err := filepath.Glob(filepath.Join(f.FontsDir, "*.ttf"))
	if err != nil {
		return nil, fmt.Errorf("failed to list .ttf files: %w", err)
	}

	otfFiles, err := filepath.Glob(filepath.Join(f.FontsDir, "*.otf"))
	if err != nil {
		return nil, fmt.Errorf("failed to list .otf files: %w", err)
	}

	// Combine the lists
	allFiles := append(ttfFiles, otfFiles...)

	// Extract font names
	for _, file := range allFiles {
		baseName := filepath.Base(file)
		// Remove the extension
		nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		// Remove common suffixes
		nameWithoutSuffix := strings.TrimSuffix(nameWithoutExt, "NerdFont")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Nerd Font")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Mono")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Regular")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Bold")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Italic")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "BoldItalic")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Light")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Medium")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Heavy")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Black")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Thin")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Hairline")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "UltraLight")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "ExtraLight")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "SemiBold")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "DemiBold")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "ExtraBold")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "UltraBold")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Heavy")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "Black")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "UltraBlack")
		nameWithoutSuffix = strings.TrimSuffix(nameWithoutSuffix, "ExtraBlack")

		// Add to the list if not already present
		fontExists := false
		for _, existingFont := range fonts {
			if existingFont == nameWithoutSuffix {
				fontExists = true
				break
			}
		}

		if !fontExists {
			fonts = append(fonts, nameWithoutSuffix)
		}
	}

	return fonts, nil
}

// GetFontDetails returns details about a specific font
func (f *FontManager) GetFontDetails(fontName string) (map[string]interface{}, error) {
	details := make(map[string]interface{})
	
	// Check if the font is installed
	if !f.IsFontInstalled(fontName) {
		return nil, fmt.Errorf("font %s is not installed", fontName)
	}
	
	// Find all font files for this font
	pattern := filepath.Join(f.FontsDir, "*"+fontName+"*.ttf")
	ttfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list .ttf files: %w", err)
	}
	
	pattern = filepath.Join(f.FontsDir, "*"+fontName+"*.otf")
	otfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list .otf files: %w", err)
	}
	
	allFiles := append(ttfFiles, otfFiles...)
	
	// Get file information
	var totalSize int64
	var fileCount int
	
	for _, file := range allFiles {
		fileInfo, err := os.Stat(file)
		if err != nil {
			continue
		}
		
		totalSize += fileInfo.Size()
		fileCount++
	}
	
	details["name"] = fontName
	details["fileCount"] = fileCount
	details["totalSize"] = totalSize
	details["files"] = allFiles
	
	return details, nil
}

// UninstallFont removes a font from the system
func (f *FontManager) UninstallFont(fontName string) error {
	// Check if the font is installed
	if !f.IsFontInstalled(fontName) {
		return fmt.Errorf("font %s is not installed", fontName)
	}
	
	// Find all font files for this font
	pattern := filepath.Join(f.FontsDir, "*"+fontName+"*.ttf")
	ttfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list .ttf files: %w", err)
	}
	
	pattern = filepath.Join(f.FontsDir, "*"+fontName+"*.otf")
	otfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to list .otf files: %w", err)
	}
	
	allFiles := append(ttfFiles, otfFiles...)
	
	// Remove each file
	for _, file := range allFiles {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove font file %s: %w", file, err)
		}
	}
	
	// Update font cache
	if err := f.updateFontCache(); err != nil {
		return fmt.Errorf("failed to update font cache: %w", err)
	}
	
	fmt.Printf("Successfully uninstalled %s\n", fontName)
	return nil
} 