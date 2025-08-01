package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dhamidi/proompt/pkg/copier"
	"github.com/dhamidi/proompt/pkg/editor"
	"github.com/dhamidi/proompt/pkg/filesystem"
	"github.com/dhamidi/proompt/pkg/picker"
	"github.com/dhamidi/proompt/pkg/prompt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func pickCmd(
	manager prompt.Manager,
	pick picker.Picker,
	ed editor.Editor,
	parser prompt.Parser,
	fs filesystem.Filesystem,
	cop copier.Copier,
) *cobra.Command {
	return &cobra.Command{
		Use:   "pick",
		Short: "Pick and process a prompt",
		Long:  "Select a prompt, fill in placeholders, and output the final result.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runPickCommand(manager, pick, ed, parser, fs, cop); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}
}

func runPickCommand(
	manager prompt.Manager,
	pick picker.Picker,
	ed editor.Editor,
	parser prompt.Parser,
	fs filesystem.Filesystem,
	cop copier.Copier,
) error {
	// Step 1: Get all prompts using manager.GetAllForPicker()
	items, err := manager.GetAllForPicker()
	if err != nil {
		return fmt.Errorf("failed to get prompts: %w", err)
	}

	if len(items) == 0 {
		return fmt.Errorf("no prompts found")
	}

	// Step 2: Use picker.Pick() to let user select
	selectedItem, err := pick.Pick(items)
	if err != nil {
		return fmt.Errorf("failed to pick prompt: %w", err)
	}

	// Get the full prompt content
	promptInfo, err := manager.Get(selectedItem.Name)
	if err != nil {
		return fmt.Errorf("failed to get prompt content: %w", err)
	}

	// Step 3: Parse selected prompt with parser.ParsePlaceholders()
	placeholders, err := parser.ParsePlaceholders(promptInfo.Content)
	if err != nil {
		return fmt.Errorf("failed to parse placeholders: %w", err)
	}

	// If no placeholders, just output the content directly
	if len(placeholders) == 0 {
		fmt.Print(promptInfo.Content)
		if err := cop.Copy(promptInfo.Content); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to copy to clipboard: %v\n", err)
		}
		return nil
	}

	// Step 4: Create temporary file with placeholders and defaults
	tempContent := generateMarkdownPlaceholderFile(placeholders, promptInfo.Content)

	tempFile, err := fs.TempFile("", "proompt-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func() {
		tempFile.Close()
		fs.Remove(tempFile.Name())
	}()

	if _, err := tempFile.WriteString(tempContent); err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close()

	// Step 5: Invoke editor.Edit() on temp file
	if err := ed.Edit(tempFile.Name()); err != nil {
		return fmt.Errorf("failed to edit file: %w", err)
	}

	// Step 6: Read back values and template content
	editedContent, err := fs.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read edited file: %w", err)
	}

	values, templateContent, err := parseMarkdownEditedValues(string(editedContent))
	if err != nil {
		return fmt.Errorf("failed to parse edited content: %w", err)
	}

	// Check if file was saved empty (abort signal)
	if len(strings.TrimSpace(string(editedContent))) == 0 {
		return fmt.Errorf("operation aborted (empty file)")
	}

	// Step 7: Output final prompt to stdout (use edited template content)
	finalContent := parser.SubstitutePlaceholders(templateContent, values)
	fmt.Print(finalContent)
	if err := cop.Copy(finalContent); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to copy to clipboard: %v\n", err)
	}

	return nil
}

// generatePlaceholderFile creates the placeholder editing experience
func generatePlaceholderFile(placeholders []prompt.Placeholder, originalContent string) string {
	var buf strings.Builder
	buf.WriteString("# Edit the values below and save the file\n")
	buf.WriteString("# Lines starting with # are ignored\n")
	buf.WriteString("# Save empty file to abort\n")
	buf.WriteString("\n")

	for _, p := range placeholders {
		buf.WriteString(fmt.Sprintf("%s=%s\n", p.Name, p.DefaultValue))
	}

	buf.WriteString("\n")
	buf.WriteString("### Full prompt preview:\n")

	for _, line := range strings.Split(originalContent, "\n") {
		buf.WriteString(fmt.Sprintf("# %s\n", line))
	}

	return buf.String()
}

// parseEditedValues parses edited values from the temporary file
func parseEditedValues(content string) (map[string]string, error) {
	values := make(map[string]string)

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			values[parts[0]] = parts[1]
		}
	}

	return values, nil
}

// generateMarkdownPlaceholderFile creates the markdown frontmatter editing experience
func generateMarkdownPlaceholderFile(placeholders []prompt.Placeholder, originalContent string) string {
	var buf strings.Builder
	
	// Write YAML frontmatter
	buf.WriteString("---\n")
	
	if len(placeholders) > 0 {
		// Create map for YAML marshaling
		vars := make(map[string]string)
		for _, p := range placeholders {
			vars[p.Name] = p.DefaultValue
		}
		
		// Marshal to YAML
		yamlBytes, err := yaml.Marshal(vars)
		if err != nil {
			// Fallback to simple format if YAML marshaling fails
			for _, p := range placeholders {
				buf.WriteString(fmt.Sprintf("%s: %s\n", p.Name, p.DefaultValue))
			}
		} else {
			buf.WriteString(string(yamlBytes))
		}
	}
	
	buf.WriteString("---\n")
	buf.WriteString(originalContent)
	
	return buf.String()
}

// parseMarkdownEditedValues parses edited values from markdown with frontmatter
func parseMarkdownEditedValues(content string) (map[string]string, string, error) {
	content = strings.TrimSpace(content)
	
	// Handle empty content
	if content == "" {
		return make(map[string]string), "", nil
	}
	
	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") {
		// No frontmatter, treat entire content as markdown
		return make(map[string]string), content, nil
	}
	
	// Find the end of frontmatter
	lines := strings.Split(content, "\n")
	var frontmatterEnd = -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			frontmatterEnd = i
			break
		}
	}
	
	if frontmatterEnd == -1 {
		return nil, "", fmt.Errorf("unclosed frontmatter delimiter")
	}
	
	// Extract frontmatter and content
	frontmatterLines := lines[1:frontmatterEnd]
	contentLines := lines[frontmatterEnd+1:]
	
	// Parse YAML frontmatter
	vars := make(map[string]string)
	if len(frontmatterLines) > 0 {
		frontmatterYAML := strings.Join(frontmatterLines, "\n")
		if strings.TrimSpace(frontmatterYAML) != "" {
			err := yaml.Unmarshal([]byte(frontmatterYAML), &vars)
			if err != nil {
				return nil, "", fmt.Errorf("invalid YAML frontmatter: %w", err)
			}
		}
	}
	
	// Join content lines
	markdownContent := strings.Join(contentLines, "\n")
	
	return vars, markdownContent, nil
}
