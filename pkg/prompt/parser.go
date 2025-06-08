package prompt

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
