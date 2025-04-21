package fonts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFontManager(t *testing.T) {
	// Test with system-wide fonts
	fm := NewFontManager("/usr/local", false)
	if fm.InstallPath != "/usr/local" {
		t.Errorf("Expected InstallPath to be /usr/local, got %s", fm.InstallPath)
	}
	if fm.FontsDir != "/usr/local/fonts" {
		t.Errorf("Expected FontsDir to be /usr/local/fonts, got %s", fm.FontsDir)
	}
	if fm.UserFonts {
		t.Error("Expected UserFonts to be false")
	}

	// Test with user fonts
	fm = NewFontManager("/usr/local", true)
	if fm.InstallPath != "/usr/local" {
		t.Errorf("Expected InstallPath to be /usr/local, got %s", fm.InstallPath)
	}
	if !fm.UserFonts {
		t.Error("Expected UserFonts to be true")
	}
}

func TestListAvailableNerdFonts(t *testing.T) {
	fm := NewFontManager("/usr/local", false)
	fonts := fm.ListAvailableNerdFonts()

	if len(fonts) == 0 {
		t.Error("Expected non-empty list of fonts")
	}

	// Check for some common fonts
	commonFonts := map[string]bool{
		"JetBrainsMono": false,
		"FiraCode":      false,
		"Hack":          false,
	}

	for _, font := range fonts {
		if _, exists := commonFonts[font]; exists {
			commonFonts[font] = true
		}
	}

	for font, found := range commonFonts {
		if !found {
			t.Errorf("Expected to find %s in available fonts", font)
		}
	}
}

func TestIsNerdFontAvailable(t *testing.T) {
	fm := NewFontManager("/usr/local", false)

	// Test with available font
	if !fm.IsNerdFontAvailable("JetBrainsMono") {
		t.Error("Expected JetBrainsMono to be available")
	}

	// Test with unavailable font
	if fm.IsNerdFontAvailable("NonExistentFont") {
		t.Error("Expected NonExistentFont to be unavailable")
	}
}

func TestIsFontInstalled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir, false)

	// Test with non-existent font
	if fm.IsFontInstalled("NonExistentFont") {
		t.Error("Expected non-existent font to be reported as not installed")
	}

	// Create a mock font file
	fontPath := filepath.Join(fm.FontsDir, "TestFont.ttf")
	if err := os.MkdirAll(fm.FontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}
	if err := os.WriteFile(fontPath, []byte("mock font data"), 0644); err != nil {
		t.Fatalf("Failed to create mock font file: %v", err)
	}

	// Test with existing font
	if !fm.IsFontInstalled("TestFont") {
		t.Error("Expected existing font to be reported as installed")
	}
}

func TestGetInstalledFonts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir, false)

	// Create mock font files
	fontsDir := fm.FontsDir
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}

	mockFonts := []string{
		"TestFont1.ttf",
		"TestFont2.ttf",
		"TestFont3.otf",
	}

	for _, font := range mockFonts {
		fontPath := filepath.Join(fontsDir, font)
		if err := os.WriteFile(fontPath, []byte("mock font data"), 0644); err != nil {
			t.Fatalf("Failed to create mock font file %s: %v", font, err)
		}
	}

	// Get installed fonts
	fonts, err := fm.GetInstalledFonts()
	if err != nil {
		t.Fatalf("Failed to get installed fonts: %v", err)
	}

	// Check if all mock fonts are in the list
	for _, mockFont := range mockFonts {
		fontName := mockFont[:len(mockFont)-4] // Remove extension
		found := false
		for _, font := range fonts {
			if font == fontName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find %s in installed fonts", fontName)
		}
	}
}

func TestGetFontDetails(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir, false)

	// Test with non-existent font
	_, err = fm.GetFontDetails("NonExistentFont")
	if err == nil {
		t.Error("Expected error when getting details of non-existent font")
	}

	// Create mock font files
	fontsDir := fm.FontsDir
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}

	mockFonts := []string{
		"TestFont1.ttf",
		"TestFont1-Bold.ttf",
		"TestFont1-Italic.ttf",
	}

	for _, font := range mockFonts {
		fontPath := filepath.Join(fontsDir, font)
		if err := os.WriteFile(fontPath, []byte("mock font data"), 0644); err != nil {
			t.Fatalf("Failed to create mock font file %s: %v", font, err)
		}
	}

	// Get font details
	details, err := fm.GetFontDetails("TestFont1")
	if err != nil {
		t.Fatalf("Failed to get font details: %v", err)
	}

	// Check details
	if details["name"] != "TestFont1" {
		t.Errorf("Expected name to be TestFont1, got %v", details["name"])
	}
	if details["fileCount"] != 3 {
		t.Errorf("Expected fileCount to be 3, got %v", details["fileCount"])
	}
	if details["totalSize"] != int64(len("mock font data")*3) {
		t.Errorf("Expected totalSize to be %d, got %v", len("mock font data")*3, details["totalSize"])
	}
	if len(details["files"].([]string)) != 3 {
		t.Errorf("Expected 3 files, got %d", len(details["files"].([]string)))
	}
}

func TestUninstallFont(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir, false)

	// Test with non-existent font
	err = fm.UninstallFont("NonExistentFont")
	if err == nil {
		t.Error("Expected error when uninstalling non-existent font")
	}

	// Create mock font files
	fontsDir := fm.FontsDir
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		t.Fatalf("Failed to create fonts directory: %v", err)
	}

	mockFonts := []string{
		"TestFont1.ttf",
		"TestFont1-Bold.ttf",
		"TestFont1-Italic.ttf",
	}

	for _, font := range mockFonts {
		fontPath := filepath.Join(fontsDir, font)
		if err := os.WriteFile(fontPath, []byte("mock font data"), 0644); err != nil {
			t.Fatalf("Failed to create mock font file %s: %v", font, err)
		}
	}

	// Uninstall font
	err = fm.UninstallFont("TestFont1")
	if err != nil {
		t.Fatalf("Failed to uninstall font: %v", err)
	}

	// Check if font files are removed
	for _, font := range mockFonts {
		fontPath := filepath.Join(fontsDir, font)
		if _, err := os.Stat(fontPath); err == nil {
			t.Errorf("Font file %s still exists after uninstallation", font)
		}
	}
}

// TestInstallNerdFonts is commented out because it requires downloading fonts,
// which can be slow and may fail in CI environments.
/*
func TestInstallNerdFonts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir, false)

	// Install a test font
	err = fm.InstallNerdFonts([]string{"JetBrainsMono"})
	if err != nil {
		t.Fatalf("Failed to install font: %v", err)
	}

	// Check if the font is installed
	if !fm.IsFontInstalled("JetBrainsMono") {
		t.Error("Expected JetBrainsMono to be installed")
	}
}
*/ 