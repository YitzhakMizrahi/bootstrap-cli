// Package validation provides validation utilities for configuration files
package validation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// SchemaValidator handles validation of YAML files against JSON schemas
type SchemaValidator struct {
	schemasDir string
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schemasDir string) *SchemaValidator {
	return &SchemaValidator{
		schemasDir: schemasDir,
	}
}

// ValidateYAML validates a YAML file against its schema
func (v *SchemaValidator) ValidateYAML(yamlPath string, schemaPath string) error {
	// Read and parse YAML file
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Convert YAML to JSON for validation
	var yamlObj interface{}
	if err := yaml.Unmarshal(yamlData, &yamlObj); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	jsonData, err := json.Marshal(yamlObj)
	if err != nil {
		return fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}

	// Load schema
	schemaLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaPath))
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	// Validate
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return fmt.Errorf("validation failed:\n%v", errors)
	}

	return nil
}

// ValidateConfig validates a configuration file against its schema
func (v *SchemaValidator) ValidateConfig(configPath string, configType string) error {
	var schemaFile string
	switch configType {
	case "tool":
		schemaFile = "tool.yaml"
	case "language":
		schemaFile = "language.yaml"
	case "font":
		schemaFile = "font.yaml"
	case "language_manager":
		schemaFile = "language_manager.yaml"
	case "dotfile":
		schemaFile = "dotfile.yaml"
	default:
		return fmt.Errorf("unknown config type: %s", configType)
	}
	schemaPath := filepath.Join(v.schemasDir, schemaFile)
	return v.ValidateYAML(configPath, schemaPath)
} 