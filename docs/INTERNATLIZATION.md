# INTERNATIONALIZATION.md

## Purpose
Ensure Bootstrap CLI can be localized and used by developers worldwide, supporting multiple languages and cultural conventions.

### String Externalization
- Move all user-facing strings into a resource file (e.g., `i18n/en.yaml`).
- Use a simple key-based system, e.g.:  
  ```yaml
  welcome.title: "✨ Bootstrap CLI ✨"
  welcome.prompt: "Start setup?"
  ```
- Load locale files at runtime based on an environment variable (e.g., `BOOTSTRAP_CLI_LOCALE`).

### Locale Detection & Fallback
- Detect locale from OS settings (`LANG`, `LC_ALL`) if not explicitly provided.
- Fallback to English (`en`) for missing translations.

### Formatting & Plurals
- Support simple pluralization rules per language in resource files.
- Externalize date/time formats and number separators according to locale.

### Right-to-Left (RTL) Support
- Automatically detect RTL locales (e.g., Arabic, Hebrew).
- Mirror prompt layouts if necessary, ensuring readability.

### Unicode & Encoding
- Enforce UTF-8 encoding for all text input/output.
- Normalize Unicode strings to NFKC form before comparisons.

### Translation Workflow
- Provide a `scripts/extract_strings.go` tool to scan for literal strings and generate `i18n/en.yaml`.
- Encourage community contributions via PRs to `i18n/*.yaml`.

### Testing
- Include CI job to validate that all keys in `i18n/en.yaml` exist in other locale files.
- Use mock locales in unit tests to assert correct string lookup.  

For more, see [Go i18n](https://github.com/nicksnyder/go-i18n).

