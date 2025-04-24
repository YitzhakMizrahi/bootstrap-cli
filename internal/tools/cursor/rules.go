package cursor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule represents a Cursor coding rule
type Rule struct {
	ID          string   `yaml:"id"`
	Category    string   `yaml:"category"`
	Description string   `yaml:"description"`
	Examples    []string `yaml:"examples,omitempty"`
	Pattern     string   `yaml:"pattern,omitempty"`
	FileGlobs   []string `yaml:"file_globs,omitempty"`
}

// RuleSet represents a collection of Cursor rules
type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}

// RuleAnalyzer analyzes code for potential new rules
type RuleAnalyzer struct {
	fset    *token.FileSet
	ruleset *RuleSet
}

// NewRuleAnalyzer creates a new rule analyzer
func NewRuleAnalyzer(rulesPath string) (*RuleAnalyzer, error) {
	data, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("reading rules file: %w", err)
	}

	var ruleset RuleSet
	if err := yaml.Unmarshal(data, &ruleset); err != nil {
		return nil, fmt.Errorf("parsing rules file: %w", err)
	}

	return &RuleAnalyzer{
		fset:    token.NewFileSet(),
		ruleset: &ruleset,
	}, nil
}

// AnalyzeDirectory analyzes a directory for potential new rules
func (ra *RuleAnalyzer) AnalyzeDirectory(dir string) ([]Rule, error) {
	var suggestedRules []Rule

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			rules, err := ra.analyzeFile(path)
			if err != nil {
				return fmt.Errorf("analyzing %s: %w", path, err)
			}
			suggestedRules = append(suggestedRules, rules...)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	return suggestedRules, nil
}

// analyzeFile analyzes a single file for potential new rules
func (ra *RuleAnalyzer) analyzeFile(path string) ([]Rule, error) {
	f, err := parser.ParseFile(ra.fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	var rules []Rule

	// Analyze interfaces
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.InterfaceType); ok {
				rules = append(rules, Rule{
					ID:          fmt.Sprintf("interface_%s", x.Name),
					Category:    "Interfaces",
					Description: fmt.Sprintf("Interface %s should be defined in internal/interfaces/", x.Name),
					FileGlobs:   []string{"internal/interfaces/*.go"},
				})
			}
		}
		return true
	})

	return rules, nil
}

// GenerateMarkdown generates markdown documentation from rules
func (ra *RuleAnalyzer) GenerateMarkdown() string {
	var sb strings.Builder

	sb.WriteString("# 🧭 Cursor Rules\n\n")

	categories := make(map[string][]Rule)
	for _, rule := range ra.ruleset.Rules {
		categories[rule.Category] = append(categories[rule.Category], rule)
	}

	for category, rules := range categories {
		sb.WriteString(fmt.Sprintf("## %s\n\n", category))
		for _, rule := range rules {
			sb.WriteString(fmt.Sprintf("### %s\n", rule.ID))
			sb.WriteString(fmt.Sprintf("%s\n\n", rule.Description))
			if len(rule.Examples) > 0 {
				sb.WriteString("Examples:\n")
				for _, example := range rule.Examples {
					sb.WriteString(fmt.Sprintf("- %s\n", example))
				}
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
} 