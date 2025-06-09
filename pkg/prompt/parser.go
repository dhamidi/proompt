package prompt

import (
	"regexp"
	"strings"
)

// Parser interface handles placeholder parsing and substitution
type Parser interface {
	ParsePlaceholders(content string) ([]Placeholder, error)
	SubstitutePlaceholders(content string, values map[string]string) string
}

// Placeholder represents a placeholder in a prompt
type Placeholder struct {
	Name         string
	DefaultValue string
	HasDefault   bool
}

// DefaultParser implements placeholder parsing
type DefaultParser struct{}

// NewDefaultParser creates a new DefaultParser
func NewDefaultParser() *DefaultParser {
	return &DefaultParser{}
}

// ParsePlaceholders parses placeholders from content using regex
// Supports ${VAR} and ${VAR:-default} syntax
func (p *DefaultParser) ParsePlaceholders(content string) ([]Placeholder, error) {
	// Regex pattern to match ${VAR} and ${VAR:-default}
	re := regexp.MustCompile(`\$\{([^}:]+)(?::-([^}]*))?\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	var placeholders []Placeholder
	seen := make(map[string]bool)
	
	for _, match := range matches {
		name := match[1]
		if seen[name] {
			continue // Skip duplicates
		}
		seen[name] = true
		
		placeholder := Placeholder{
			Name: name,
		}
		
		// Check if this is a ${VAR:-default} pattern
		if len(match) > 2 && strings.Contains(match[0], ":-") {
			placeholder.HasDefault = true
			placeholder.DefaultValue = match[2]
		}
		
		placeholders = append(placeholders, placeholder)
	}
	
	return placeholders, nil
}

// SubstitutePlaceholders replaces placeholders with provided values
// Handles $$ escaping for literal $
func (p *DefaultParser) SubstitutePlaceholders(content string, values map[string]string) string {
	// First handle literal $$ -> $
	result := strings.ReplaceAll(content, "$$", "\x00LITERAL_DOLLAR\x00")
	
	// Regex pattern to match ${VAR} and ${VAR:-default}
	re := regexp.MustCompile(`\$\{([^}:]+)(?::-([^}]*))?\}`)
	
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		submatch := re.FindStringSubmatch(match)
		name := submatch[1]
		defaultValue := ""
		if len(submatch) > 2 {
			defaultValue = submatch[2]
		}
		
		if value, exists := values[name]; exists {
			return value
		}
		return defaultValue
	})
	
	// Restore literal dollars
	result = strings.ReplaceAll(result, "\x00LITERAL_DOLLAR\x00", "$")
	
	return result
}

// FakeParser simulates parser behavior for testing
type FakeParser struct {
	Placeholders []Placeholder
}

// NewFakeParser creates a new FakeParser
func NewFakeParser() *FakeParser {
	return &FakeParser{
		Placeholders: make([]Placeholder, 0),
	}
}

// ParsePlaceholders returns the predefined placeholders for testing
func (f *FakeParser) ParsePlaceholders(content string) ([]Placeholder, error) {
	return f.Placeholders, nil
}

// SubstitutePlaceholders performs simple string replacement for testing
func (f *FakeParser) SubstitutePlaceholders(content string, values map[string]string) string {
	result := content
	for key, value := range values {
		result = strings.ReplaceAll(result, "${"+key+"}", value)
		result = strings.ReplaceAll(result, "${"+key+":-", value+"}")
	}
	return result
}
