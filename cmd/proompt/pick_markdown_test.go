package main

import (
	"strings"
	"testing"

	"github.com/dhamidi/proompt/pkg/prompt"
)

func TestGenerateMarkdownPlaceholderFile(t *testing.T) {
	tests := []struct {
		name        string
		placeholders []prompt.Placeholder
		content     string
		checkFunc   func(result string) bool
	}{
		{
			name: "simple placeholders",
			placeholders: []prompt.Placeholder{
				{Name: "NAME", DefaultValue: "World"},
				{Name: "SCORE", DefaultValue: "100"},
			},
			content: "Hello ${NAME}! Your score is ${SCORE}.",
			checkFunc: func(result string) bool {
				return strings.Contains(result, "---\n") &&
					strings.Contains(result, "NAME: World") &&
					strings.Contains(result, "SCORE:") &&
					strings.Contains(result, "100") &&
					strings.Contains(result, "---\nHello ${NAME}! Your score is ${SCORE}.")
			},
		},
		{
			name: "no placeholders",
			placeholders: []prompt.Placeholder{},
			content: "This is a simple prompt.",
			checkFunc: func(result string) bool {
				return result == "---\n---\nThis is a simple prompt."
			},
		},
		{
			name: "placeholders with empty defaults",
			placeholders: []prompt.Placeholder{
				{Name: "VAR1", DefaultValue: ""},
				{Name: "VAR2", DefaultValue: "default"},
			},
			content: "Using ${VAR1} and ${VAR2}.",
			checkFunc: func(result string) bool {
				return strings.Contains(result, "VAR1:") &&
					strings.Contains(result, "VAR2: default") &&
					strings.Contains(result, "Using ${VAR1} and ${VAR2}.")
			},
		},
		{
			name: "multiline content",
			placeholders: []prompt.Placeholder{
				{Name: "TITLE", DefaultValue: "My Title"},
			},
			content: "# ${TITLE}\n\nThis is a multiline\nprompt template.",
			checkFunc: func(result string) bool {
				return strings.Contains(result, "TITLE: My Title") &&
					strings.Contains(result, "# ${TITLE}\n\nThis is a multiline\nprompt template.")
			},
		},
		{
			name: "placeholders with special yaml characters",
			placeholders: []prompt.Placeholder{
				{Name: "MESSAGE", DefaultValue: "Hello: world!"},
				{Name: "LIST", DefaultValue: "- item1\n- item2"},
			},
			content: "Message: ${MESSAGE}\nList: ${LIST}",
			checkFunc: func(result string) bool {
				return strings.Contains(result, "MESSAGE:") &&
					strings.Contains(result, "Hello: world!") &&
					strings.Contains(result, "LIST:") &&
					strings.Contains(result, "item1") &&
					strings.Contains(result, "Message: ${MESSAGE}\nList: ${LIST}")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateMarkdownPlaceholderFile(tt.placeholders, tt.content)
			if !tt.checkFunc(result) {
				t.Errorf("generateMarkdownPlaceholderFile() failed validation. Got: %q", result)
			}
		})
	}
}

func TestParseMarkdownEditedValues(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectedVars map[string]string
		expectedContent string
		expectError bool
	}{
		{
			name: "valid frontmatter and content",
			content: `---
NAME: Alice
SCORE: 95
---
Hello ${NAME}! Your score is ${SCORE}.`,
			expectedVars: map[string]string{
				"NAME":  "Alice",
				"SCORE": "95",
			},
			expectedContent: "Hello ${NAME}! Your score is ${SCORE}.",
			expectError: false,
		},
		{
			name: "empty frontmatter",
			content: `---
---
Simple content without variables.`,
			expectedVars: map[string]string{},
			expectedContent: "Simple content without variables.",
			expectError: false,
		},
		{
			name: "no frontmatter",
			content: "Just plain content",
			expectedVars: map[string]string{},
			expectedContent: "Just plain content",
			expectError: false,
		},
		{
			name: "multiline content",
			content: `---
TITLE: My Document
---
# ${TITLE}

This is a multiline
document with content.`,
			expectedVars: map[string]string{
				"TITLE": "My Document",
			},
			expectedContent: "# ${TITLE}\n\nThis is a multiline\ndocument with content.",
			expectError: false,
		},
		{
			name: "frontmatter with quoted values",
			content: `---
MESSAGE: "Hello: world!"
PATH: "/path/to/file"
---
Message: ${MESSAGE}
Path: ${PATH}`,
			expectedVars: map[string]string{
				"MESSAGE": "Hello: world!",
				"PATH": "/path/to/file",
			},
			expectedContent: "Message: ${MESSAGE}\nPath: ${PATH}",
			expectError: false,
		},
		{
			name: "frontmatter with multiline values",
			content: `---
CONFIG: |
  key1: value1
  key2: value2
TEXT: |-
  Line 1
  Line 2
---
Config: ${CONFIG}
Text: ${TEXT}`,
			expectedVars: map[string]string{
				"CONFIG": "key1: value1\nkey2: value2\n",
				"TEXT": "Line 1\nLine 2",
			},
			expectedContent: "Config: ${CONFIG}\nText: ${TEXT}",
			expectError: false,
		},
		{
			name: "invalid yaml frontmatter",
			content: `---
INVALID: [unclosed
---
Content here`,
			expectedVars: nil,
			expectedContent: "",
			expectError: true,
		},
		{
			name: "empty file",
			content: "",
			expectedVars: map[string]string{},
			expectedContent: "",
			expectError: false,
		},
		{
			name: "only frontmatter delimiter",
			content: `---
VAR: value
`,
			expectedVars: nil,
			expectedContent: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars, content, err := parseMarkdownEditedValues(tt.content)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(vars) != len(tt.expectedVars) {
				t.Errorf("Expected %d variables, got %d", len(tt.expectedVars), len(vars))
			}

			for key, expectedValue := range tt.expectedVars {
				if actualValue, exists := vars[key]; !exists {
					t.Errorf("Expected variable %q not found", key)
				} else if actualValue != expectedValue {
					t.Errorf("Variable %q: expected %q, got %q", key, expectedValue, actualValue)
				}
			}

			if content != tt.expectedContent {
				t.Errorf("Expected content %q, got %q", tt.expectedContent, content)
			}
		})
	}
}

func TestParseMarkdownEditedValuesBackwardCompatibility(t *testing.T) {
	// Test that the function can still handle the old format for compatibility
	// if needed (VAR=value format)
	content := `NAME=Alice
SCORE=95

# Full prompt preview:
# Hello ${NAME}! Your score is ${SCORE}.`

	vars, _, err := parseMarkdownEditedValues(content)
	if err != nil {
		t.Errorf("Parsing old format should not error: %v", err)
	}

	// For this test, we expect the old format to be treated as plain content
	// since we're changing to markdown format only
	if len(vars) != 0 {
		t.Error("Old format should not parse variables in new function")
	}
}
