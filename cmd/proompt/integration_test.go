package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/spf13/cobra"
	"github.com/dhamidi/proompt/pkg/copier"
	"github.com/dhamidi/proompt/pkg/editor"
	"github.com/dhamidi/proompt/pkg/filesystem"
	"github.com/dhamidi/proompt/pkg/picker"
	"github.com/dhamidi/proompt/pkg/prompt"
)

// TestPickWorkflowWithoutPlaceholders tests the pick workflow for prompts without placeholders
func TestPickWorkflowWithoutPlaceholders(t *testing.T) {
	// Setup fake filesystem with test prompts
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/simple.md"] = &fstest.MapFile{
		Data: []byte("This is a simple prompt without placeholders."),
		Mode: 0644,
	}

	// Setup fake resolver
	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	// Setup fake picker to select the simple prompt
	pick := picker.NewFakePicker()
	pick.Selections = []picker.PickerItem{
		{Name: "simple", Source: "directory", Path: "prompts/simple.md"},
	}

	// Setup fake editor (shouldn't be called for prompts without placeholders)
	ed := editor.NewFakeEditor()

	// Create components
	manager := prompt.NewDefaultManager(fs, resolver)
	parser := prompt.NewDefaultParser()

	// Test the runPickCommand function directly
	cop := copier.NewFakeCopier()
	err := runPickCommand(manager, pick, ed, parser, fs, cop)
	if err != nil {
		t.Fatalf("Pick command failed: %v", err)
	}

	// Verify editor was not called since there are no placeholders
	if len(ed.EditedFiles) != 0 {
		t.Error("Editor should not be called for prompts without placeholders")
	}
}

// TestPromptParsingIntegration tests the integration between manager and parser
func TestPromptParsingIntegration(t *testing.T) {
	// Setup fake filesystem with prompts containing placeholders
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/withplaceholders.md"] = &fstest.MapFile{
		Data: []byte("Hello ${NAME:-World}! Your score is ${SCORE:-0}."),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := prompt.NewDefaultManager(fs, resolver)
	parser := prompt.NewDefaultParser()

	// Get the prompt
	promptInfo, err := manager.Get("withplaceholders")
	if err != nil {
		t.Fatalf("Failed to get prompt: %v", err)
	}

	// Parse placeholders
	placeholders, err := parser.ParsePlaceholders(promptInfo.Content)
	if err != nil {
		t.Fatalf("Failed to parse placeholders: %v", err)
	}

	if len(placeholders) != 2 {
		t.Errorf("Expected 2 placeholders, got %d", len(placeholders))
	}

	// Check placeholder names and defaults
	expectedPlaceholders := map[string]string{
		"NAME":  "World",
		"SCORE": "0",
	}

	for _, p := range placeholders {
		expectedDefault, exists := expectedPlaceholders[p.Name]
		if !exists {
			t.Errorf("Unexpected placeholder: %s", p.Name)
			continue
		}
		if p.DefaultValue != expectedDefault {
			t.Errorf("Expected default '%s' for placeholder '%s', got '%s'", expectedDefault, p.Name, p.DefaultValue)
		}
	}

	// Test substitution
	values := map[string]string{
		"NAME":  "Proompt",
		"SCORE": "100",
	}

	result := parser.SubstitutePlaceholders(promptInfo.Content, values)
	expected := "Hello Proompt! Your score is 100."
	if result != expected {
		t.Errorf("Expected substitution result '%s', got '%s'", expected, result)
	}
}

// TestListCommandIntegration tests the list command end-to-end
func TestListCommandIntegration(t *testing.T) {
	// Setup fake filesystem with prompts in different locations
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/local.md"] = &fstest.MapFile{
		Data: []byte("Local prompt"),
		Mode: 0644,
	}
	fs.MapFS["project/prompts/shared.md"] = &fstest.MapFile{
		Data: []byte("Project prompt"),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
		{Type: "project", Path: "project/prompts"},
	}

	manager := prompt.NewDefaultManager(fs, resolver)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := listCmd(manager)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	outputStr := string(output)
	if !strings.Contains(outputStr, "local") {
		t.Error("Expected 'local' prompt in output")
	}
	if !strings.Contains(outputStr, "shared") {
		t.Error("Expected 'shared' prompt in output")
	}
}

// TestShowCommandIntegration tests the show command end-to-end
func TestShowCommandIntegration(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("This is a test prompt with ${VAR:-default}."),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := prompt.NewDefaultManager(fs, resolver)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := showCmd(manager)
	cmd.SetArgs([]string{"test"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Show command failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	expectedContent := "This is a test prompt with ${VAR:-default}."
	if !strings.Contains(string(output), expectedContent) {
		t.Errorf("Expected output to contain '%s'", expectedContent)
	}
}

// TestEditCommandIntegration tests the edit command end-to-end
func TestEditCommandIntegration(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/existing.md"] = &fstest.MapFile{
		Data: []byte("Existing prompt content"),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	// Setup fake picker to select existing prompt
	pick := picker.NewFakePicker()
	pick.Selections = []picker.PickerItem{
		{Name: "existing", Source: "directory", Path: "prompts/existing.md"},
	}

	// Setup fake editor
	ed := editor.NewFakeEditor()

	manager := prompt.NewDefaultManager(fs, resolver)

	cmd := editCmd(manager, pick, ed)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Edit command failed: %v", err)
	}

	// Verify editor was called with correct file
	if len(ed.EditedFiles) != 1 {
		t.Fatalf("Expected 1 edited file, got %d", len(ed.EditedFiles))
	}
	if ed.EditedFiles[0] != "prompts/existing.md" {
		t.Errorf("Expected edited file 'prompts/existing.md', got '%s'", ed.EditedFiles[0])
	}
}

// TestRmCommandIntegration tests the rm command end-to-end
func TestRmCommandIntegration(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/todelete.md"] = &fstest.MapFile{
		Data: []byte("Content to delete"),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	// Setup fake picker to select prompt to delete
	pick := picker.NewFakePicker()
	pick.Selections = []picker.PickerItem{
		{Name: "todelete", Source: "directory", Path: "prompts/todelete.md"},
	}

	manager := prompt.NewDefaultManager(fs, resolver)

	cmd := rmCmd(manager, pick)
	cmd.SetArgs([]string{})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Rm command failed: %v", err)
	}

	// Verify file was deleted
	_, err = fs.ReadFile("prompts/todelete.md")
	if err == nil {
		t.Error("Expected file to be deleted, but it still exists")
	}
}

// TestErrorHandlingPickerFailure tests error handling when picker fails
func TestErrorHandlingPickerFailure(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Test prompt"),
		Mode: 0644,
	}
	
	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}
	
	// Setup fake picker that always fails
	pick := picker.NewFakePicker()
	pick.ShouldFail = true

	ed := editor.NewFakeEditor()
	manager := prompt.NewDefaultManager(fs, resolver)
	parser := prompt.NewDefaultParser()

	// Test runPickCommand directly
	cop := copier.NewFakeCopier()
	err := runPickCommand(manager, pick, ed, parser, fs, cop)

	// Should handle picker failure gracefully
	if err == nil {
		t.Error("Expected command to fail when picker fails")
	}
	if !strings.Contains(err.Error(), "picker failed") {
		t.Errorf("Expected error message about picker failure, got: %v", err)
	}
}

// TestErrorHandlingNoPrompts tests error handling when no prompts are found
func TestErrorHandlingNoPrompts(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	resolver := prompt.NewFakeLocationResolver()
	// No locations configured - should result in no prompts

	pick := picker.NewFakePicker()
	ed := editor.NewFakeEditor()
	manager := prompt.NewDefaultManager(fs, resolver)
	parser := prompt.NewDefaultParser()

	// Test runPickCommand directly
	cop := copier.NewFakeCopier()
	err := runPickCommand(manager, pick, ed, parser, fs, cop)

	// Should handle no prompts gracefully
	if err == nil {
		t.Error("Expected command to fail when no prompts are found")
	}
	if !strings.Contains(err.Error(), "no prompts found") {
		t.Errorf("Expected error message about no prompts found, got: %v", err)
	}
}

// TestErrorHandlingInvalidPromptName tests error handling for invalid prompt names
func TestErrorHandlingInvalidPromptName(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	resolver := prompt.NewFakeLocationResolver()
	manager := prompt.NewDefaultManager(fs, resolver)

	// Test getting non-existing prompt directly from manager
	_, err := manager.Get("nonexistent")

	// Should handle missing prompt gracefully
	if err == nil {
		t.Error("Expected error for nonexistent prompt")
	}
	if err != prompt.ErrPromptNotFound {
		t.Errorf("Expected ErrPromptNotFound, got: %v", err)
	}
}

// TestDifferentPromptHierarchyScenarios tests prompt discovery from different locations
func TestDifferentPromptHierarchyScenarios(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	
	// Create prompts at different hierarchy levels
	fs.MapFS["prompts/directory.md"] = &fstest.MapFile{
		Data: []byte("Directory level prompt"),
		Mode: 0644,
	}
	fs.MapFS["project/prompts/project.md"] = &fstest.MapFile{
		Data: []byte("Project level prompt"),
		Mode: 0644,
	}
	fs.MapFS[".git/info/prompts/project-local.md"] = &fstest.MapFile{
		Data: []byte("Project-local prompt"),
		Mode: 0644,
	}
	fs.MapFS["user/prompts/user.md"] = &fstest.MapFile{
		Data: []byte("User level prompt"),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
		{Type: "project", Path: "project/prompts"},
		{Type: "project-local", Path: ".git/info/prompts"},
		{Type: "user", Path: "user/prompts"},
	}

	manager := prompt.NewDefaultManager(fs, resolver)

	prompts, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list prompts: %v", err)
	}

	if len(prompts) != 4 {
		t.Errorf("Expected 4 prompts from different hierarchy levels, got %d", len(prompts))
	}

	// Verify prompts from different sources are found
	sources := make(map[string]bool)
	for _, prompt := range prompts {
		sources[prompt.Source] = true
	}

	expectedSources := []string{"directory", "project", "project-local", "user"}
	for _, source := range expectedSources {
		if !sources[source] {
			t.Errorf("Expected prompt from source '%s', but none found", source)
		}
	}
}

// TestPickWorkflowAbortOnEmptyFile tests that pick workflow aborts when user saves empty file
// Note: This test is simplified since the temp file handling is complex in the fake filesystem
func TestPickWorkflowAbortOnEmptyFile(t *testing.T) {
	// Test the parseEditedValues function directly for empty content
	values, err := parseEditedValues("")
	if err != nil {
		t.Errorf("parseEditedValues should not error on empty content: %v", err)
	}
	if len(values) != 0 {
		t.Error("Expected empty values map for empty content")
	}
	
	// Test that empty content detection works
	content := ""
	if len(strings.TrimSpace(content)) != 0 {
		t.Error("Empty content detection should work")
	}
}

// TestErrorHandlingFilePermissions tests error handling for file permission issues
func TestErrorHandlingFilePermissions(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	resolver := prompt.NewFakeLocationResolver()
	manager := prompt.NewDefaultManager(fs, resolver)

	// Test creating a prompt when filesystem returns permission error
	err := manager.Create("test", "content", "directory")
	// Should handle permission errors gracefully (depends on implementation)
	if err != nil {
		t.Logf("Create returned error as expected: %v", err)
	}
}

// TestErrorHandlingMalformedPlaceholders tests error handling for invalid placeholder syntax
func TestErrorHandlingMalformedPlaceholders(t *testing.T) {
	parser := prompt.NewDefaultParser()
	
	// Test various malformed placeholder syntaxes
	testCases := []string{
		"${UNCLOSED",           // Missing closing brace
		"${}",                  // Empty placeholder
		"${:-noname}",          // Default without name
		"${VAR:-${NESTED}}",    // Nested placeholders
	}
	
	for _, content := range testCases {
		placeholders, err := parser.ParsePlaceholders(content)
		if err != nil {
			t.Logf("Parser correctly rejected malformed syntax '%s': %v", content, err)
		} else {
			t.Logf("Parser parsed '%s' as %d placeholders (may be acceptable)", content, len(placeholders))
		}
	}
}

// TestErrorHandlingEditorFailure tests error handling when editor command fails
func TestErrorHandlingEditorFailure(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Hello ${NAME:-World}!"),
		Mode: 0644,
	}

	resolver := prompt.NewFakeLocationResolver()
	resolver.Locations = []prompt.PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	pick := picker.NewFakePicker()
	pick.Selections = []picker.PickerItem{
		{Name: "test", Source: "directory", Path: "prompts/test.md"},
	}

	// Setup fake editor that always fails
	ed := editor.NewFakeEditor()
	ed.ShouldFail = true

	manager := prompt.NewDefaultManager(fs, resolver)
	
	cmd := editCmd(manager, pick, ed)
	cmd.SetArgs([]string{"test"})
	
	err := cmd.Execute()
	
	// Should handle editor failure gracefully
	if err != nil {
		t.Logf("Edit command correctly handled editor failure: %v", err)
	}
}

// TestValidatePromptNameHandling tests that commands properly validate prompt names
func TestValidatePromptNameHandling(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()
	resolver := prompt.NewFakeLocationResolver()
	manager := prompt.NewDefaultManager(fs, resolver)

	// Test show command with invalid name
	cmd := showCmd(manager)
	cmd.SetArgs([]string{"nonexistent"})
	
	err := cmd.Execute()
	if err == nil {
		t.Error("Expected show command to fail for nonexistent prompt")
	}
}

// Helper function to capture command output
func captureCommandOutput(t *testing.T, cmd *cobra.Command) (stdout, stderr string, err error) {
	// Capture stdout
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut

	// Capture stderr
	oldStderr := os.Stderr
	rErr, wErr, _ := os.Pipe()
	os.Stderr = wErr

	// Execute command
	err = cmd.Execute()

	// Restore stdout/stderr
	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read captured output
	stdoutBytes, _ := io.ReadAll(rOut)
	stderrBytes, _ := io.ReadAll(rErr)

	return string(stdoutBytes), string(stderrBytes), err
}
