package fonts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFontManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with system-wide fonts
	manager := NewFontManager(tempDir, false)
	expectedFontsDir := filepath.Join(tempDir, "fonts")
	if manager.FontsDir != expectedFontsDir {
		t.Errorf("Expected FontsDir to be %s, got %s", expectedFontsDir, manager.FontsDir)
	}
	if manager.UserFonts {
		t.Error("Expected UserFonts to be false, got true")
	}

	// Test with user-specific fonts
	manager = NewFontManager(tempDir, true)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		expectedUserFontsDir := filepath.Join(homeDir, ".local", "share", "fonts")
		if manager.FontsDir != expectedUserFontsDir {
			t.Errorf("Expected FontsDir to be %s, got %s", expectedUserFontsDir, manager.FontsDir)
		}
		if !manager.UserFonts {
			t.Error("Expected UserFonts to be true, got false")
		}
	}
}

func TestListAvailableNerdFonts(t *testing.T) {
	manager := NewFontManager("/tmp", false)
	fonts := manager.ListAvailableNerdFonts()

	// Check if the list is not empty
	if len(fonts) == 0 {
		t.Error("Expected a non-empty list of available Nerd Fonts")
	}

	// Check if some common fonts are in the list
	commonFonts := []string{"JetBrainsMono", "FiraCode", "Hack"}
	for _, font := range commonFonts {
		found := false
		for _, availableFont := range fonts {
			if availableFont == font {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s to be in the list of available fonts", font)
		}
	}
}

func TestIsFontInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new FontManager
	manager := NewFontManager(tempDir, false)

	// Create the fonts directory
	if err := os.MkdirAll(manager.FontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}

	// Check if a non-existent font is reported as not installed
	if manager.IsFontInstalled("NonExistentFont") {
		t.Error("Expected NonExistentFont to be reported as not installed")
	}

	// Create a mock font file
	mockFontFile := filepath.Join(manager.FontsDir, "JetBrainsMonoNerdFont-Regular.ttf")
	if err := os.WriteFile(mockFontFile, []byte("mock font data"), 0644); err != nil {
		t.Fatalf("Failed to create mock font file: %v", err)
	}

	// Check if the font is reported as installed
	if !manager.IsFontInstalled("JetBrainsMono") {
		t.Error("Expected JetBrainsMono to be reported as installed")
	}
}

func TestGetInstalledFonts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new FontManager
	manager := NewFontManager(tempDir, false)

	// Create the fonts directory
	if err := os.MkdirAll(manager.FontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}

	// Create mock font files
	mockFonts := []string{
		"JetBrainsMonoNerdFont-Regular.ttf",
		"FiraCodeNerdFont-Bold.ttf",
		"HackNerdFont-Italic.ttf",
	}

	for _, font := range mockFonts {
		fontPath := filepath.Join(manager.FontsDir, font)
		if err := os.WriteFile(fontPath, []byte("mock font data"), 0644); err != nil {
			t.Fatalf("Failed to create mock font file %s: %v", font, err)
		}
	}

	// Get the list of installed fonts
	installedFonts, err := manager.GetInstalledFonts()
	if err != nil {
		t.Fatalf("Failed to get installed fonts: %v", err)
	}

	// Check if the expected fonts are in the list
	expectedFonts := []string{"JetBrainsMono", "FiraCode", "Hack"}
	for _, expectedFont := range expectedFonts {
		found := false
		for _, installedFont := range installedFonts {
			if installedFont == expectedFont {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s to be in the list of installed fonts", expectedFont)
		}
	}
}

// Note: The following test is commented out because it requires downloading fonts,
// which can be slow and may fail in CI environments.

/*
func TestInstallNerdFonts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a new FontManager
	manager := NewFontManager(tempDir, false)

	// Install a Nerd Font
	fontsToInstall := []string{"JetBrainsMono"}
	if err := manager.InstallNerdFonts(fontsToInstall); err != nil {
		t.Fatalf("Failed to install Nerd Fonts: %v", err)
	}

	// Check if the font is installed
	if !manager.IsFontInstalled("JetBrainsMono") {
		t.Error("Expected JetBrainsMono to be installed")
	}
}
*/ 