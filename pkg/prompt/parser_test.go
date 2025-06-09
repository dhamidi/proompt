package prompt

import (
	"reflect"
	"testing"
)

func TestDefaultParser_ParsePlaceholders(t *testing.T) {
	parser := NewDefaultParser()

	tests := []struct {
		name     string
		content  string
		expected []Placeholder
	}{
		{
			name:     "no placeholders",
			content:  "Hello world!",
			expected: []Placeholder{},
		},
		{
			name:    "simple placeholder",
			content: "Hello ${NAME}!",
			expected: []Placeholder{
				{Name: "NAME", DefaultValue: "", HasDefault: false},
			},
		},
		{
			name:    "placeholder with default",
			content: "Hello ${NAME:-World}!",
			expected: []Placeholder{
				{Name: "NAME", DefaultValue: "World", HasDefault: true},
			},
		},
		{
			name:    "multiple placeholders",
			content: "Hello ${FIRST} ${LAST:-Doe}!",
			expected: []Placeholder{
				{Name: "FIRST", DefaultValue: "", HasDefault: false},
				{Name: "LAST", DefaultValue: "Doe", HasDefault: true},
			},
		},
		{
			name:    "duplicate placeholders",
			content: "Hello ${NAME} and ${NAME} again!",
			expected: []Placeholder{
				{Name: "NAME", DefaultValue: "", HasDefault: false},
			},
		},
		{
			name:    "empty default value",
			content: "Hello ${NAME:-}!",
			expected: []Placeholder{
				{Name: "NAME", DefaultValue: "", HasDefault: true},
			},
		},
		{
			name:    "complex default with spaces",
			content: "Hello ${NAME:-John Doe}!",
			expected: []Placeholder{
				{Name: "NAME", DefaultValue: "John Doe", HasDefault: true},
			},
		},
		{
			name:     "escaped dollar signs",
			content:  "Price: $$100",
			expected: []Placeholder{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			placeholders, err := parser.ParsePlaceholders(tt.content)
			if err != nil {
				t.Errorf("ParsePlaceholders() error = %v", err)
				return
			}
			// Handle nil vs empty slice comparison
			if len(placeholders) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(placeholders, tt.expected) {
				t.Errorf("ParsePlaceholders() = %v, want %v", placeholders, tt.expected)
			}
		})
	}
}

func TestDefaultParser_SubstitutePlaceholders(t *testing.T) {
	parser := NewDefaultParser()

	tests := []struct {
		name     string
		content  string
		values   map[string]string
		expected string
	}{
		{
			name:     "no placeholders",
			content:  "Hello world!",
			values:   map[string]string{},
			expected: "Hello world!",
		},
		{
			name:    "simple substitution",
			content: "Hello ${NAME}!",
			values:  map[string]string{"NAME": "John"},
			expected: "Hello John!",
		},
		{
			name:    "use default value",
			content: "Hello ${NAME:-World}!",
			values:  map[string]string{},
			expected: "Hello World!",
		},
		{
			name:    "override default value",
			content: "Hello ${NAME:-World}!",
			values:  map[string]string{"NAME": "John"},
			expected: "Hello John!",
		},
		{
			name:    "multiple substitutions",
			content: "Hello ${FIRST} ${LAST:-Doe}!",
			values:  map[string]string{"FIRST": "John"},
			expected: "Hello John Doe!",
		},
		{
			name:    "empty value overrides default",
			content: "Hello ${NAME:-World}!",
			values:  map[string]string{"NAME": ""},
			expected: "Hello !",
		},
		{
			name:    "missing value no default",
			content: "Hello ${NAME}!",
			values:  map[string]string{},
			expected: "Hello !",
		},
		{
			name:     "escaped dollar signs",
			content:  "Price: $$100",
			values:   map[string]string{},
			expected: "Price: $100",
		},
		{
			name:    "mixed escaped and placeholders",
			content: "Price: $$${AMOUNT:-50}",
			values:  map[string]string{},
			expected: "Price: $50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.SubstitutePlaceholders(tt.content, tt.values)
			if result != tt.expected {
				t.Errorf("SubstitutePlaceholders() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFakeParser(t *testing.T) {
	fakeParser := NewFakeParser()
	fakeParser.Placeholders = []Placeholder{
		{Name: "TEST", DefaultValue: "default", HasDefault: true},
	}

	// Test ParsePlaceholders
	placeholders, err := fakeParser.ParsePlaceholders("any content")
	if err != nil {
		t.Errorf("FakeParser.ParsePlaceholders() error = %v", err)
		return
	}
	if len(placeholders) != 1 || placeholders[0].Name != "TEST" {
		t.Errorf("FakeParser.ParsePlaceholders() = %v, want placeholder with name TEST", placeholders)
	}

	// Test SubstitutePlaceholders
	content := "Hello ${NAME}!"
	values := map[string]string{"NAME": "World"}
	result := fakeParser.SubstitutePlaceholders(content, values)
	expected := "Hello World!"
	if result != expected {
		t.Errorf("FakeParser.SubstitutePlaceholders() = %q, want %q", result, expected)
	}
}
