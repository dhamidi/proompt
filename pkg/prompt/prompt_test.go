package prompt

import (
	"testing"
	"testing/fstest"

	"github.com/dhamidi/proompt/pkg/filesystem"
)

func TestDefaultManagerList(t *testing.T) {
	// Create fake filesystem with test prompts
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test1.md"] = &fstest.MapFile{
		Data: []byte("Hello ${NAME:-World}!"),
		Mode: 0644,
	}
	fs.MapFS["prompts/test2.txt"] = &fstest.MapFile{
		Data: []byte("Another prompt"),
		Mode: 0644,
	}
	fs.MapFS["prompts/not_a_prompt.go"] = &fstest.MapFile{
		Data: []byte("package main"),
		Mode: 0644,
	}

	// Create fake resolver
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	prompts, err := manager.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(prompts) != 2 {
		t.Errorf("Expected 2 prompts, got %d", len(prompts))
	}

	// Check that non-prompt files are excluded
	for _, prompt := range prompts {
		if prompt.Name == "not_a_prompt" {
			t.Error("Non-prompt file was included in results")
		}
	}

	// Check that test1 prompt is found correctly
	found := false
	for _, prompt := range prompts {
		if prompt.Name == "test1" {
			found = true
			if prompt.Content != "Hello ${NAME:-World}!" {
				t.Errorf("Expected content 'Hello ${NAME:-World}!', got '%s'", prompt.Content)
			}
			if prompt.Source != "directory" {
				t.Errorf("Expected source 'directory', got '%s'", prompt.Source)
			}
		}
	}
	if !found {
		t.Error("test1 prompt was not found")
	}
}

func TestDefaultManagerGet(t *testing.T) {
	// Create fake filesystem with test prompt
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Test content"),
		Mode: 0644,
	}

	// Create fake resolver
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	// Test getting existing prompt
	prompt, err := manager.Get("test")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if prompt.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", prompt.Name)
	}
	if prompt.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got '%s'", prompt.Content)
	}

	// Test getting non-existing prompt
	_, err = manager.Get("nonexistent")
	if err != ErrPromptNotFound {
		t.Errorf("Expected ErrPromptNotFound, got %v", err)
	}
}

func TestDefaultManagerCreate(t *testing.T) {
	fs := filesystem.NewFakeFilesystem()

	// Create fake resolver
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	// Test creating a new prompt
	err := manager.Create("newprompt", "New content", "directory")
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	// Verify the prompt was created
	content, err := fs.ReadFile("prompts/newprompt.md")
	if err != nil {
		t.Fatalf("Failed to read created prompt: %v", err)
	}

	if string(content) != "New content" {
		t.Errorf("Expected content 'New content', got '%s'", string(content))
	}

	// Test creating with invalid location
	err = manager.Create("another", "content", "invalid")
	if err != ErrInvalidLocation {
		t.Errorf("Expected ErrInvalidLocation, got %v", err)
	}
}

func TestDefaultManagerDelete(t *testing.T) {
	// Create fake filesystem with test prompt
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Test content"),
		Mode: 0644,
	}

	// Create fake resolver
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	// Test deleting existing prompt
	err := manager.Delete("test")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify the prompt was deleted
	_, err = fs.ReadFile("prompts/test.md")
	if err == nil {
		t.Error("Prompt file still exists after deletion")
	}

	// Test deleting non-existing prompt
	err = manager.Delete("nonexistent")
	if err != ErrPromptNotFound {
		t.Errorf("Expected ErrPromptNotFound, got %v", err)
	}
}

func TestDefaultManagerGetAllForPicker(t *testing.T) {
	// Create fake filesystem with test prompts
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test1.md"] = &fstest.MapFile{
		Data: []byte("Content 1"),
		Mode: 0644,
	}
	fs.MapFS["prompts/test2.md"] = &fstest.MapFile{
		Data: []byte("Content 2"),
		Mode: 0644,
	}

	// Create fake resolver
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	items, err := manager.GetAllForPicker()
	if err != nil {
		t.Fatalf("GetAllForPicker() failed: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(items))
	}

	// Check that items have correct structure
	for _, item := range items {
		if item.Name == "" {
			t.Error("Item name is empty")
		}
		if item.Source == "" {
			t.Error("Item source is empty")
		}
		if item.Path == "" {
			t.Error("Item path is empty")
		}
	}
}

func TestDefaultManagerListDeduplication(t *testing.T) {
	// Create fake filesystem with test prompts
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Test content"),
		Mode: 0644,
	}

	// Create fake resolver that returns duplicate locations pointing to same directory
	// This simulates the real bug where "directory" and "project" both resolve to "prompts/"
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
		{Type: "project", Path: "prompts"}, // Same path, different type
	}

	manager := NewDefaultManager(fs, resolver)

	prompts, err := manager.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	// Should only return one prompt, not duplicates
	if len(prompts) != 1 {
		t.Errorf("Expected 1 unique prompt after deduplication, got %d", len(prompts))
		for i, prompt := range prompts {
			t.Logf("Prompt %d: Name=%s, Source=%s, Path=%s", i, prompt.Name, prompt.Source, prompt.Path)
		}
	}

	// Verify it's the correct prompt
	if prompts[0].Name != "test" {
		t.Errorf("Expected prompt name 'test', got '%s'", prompts[0].Name)
	}
}

func TestIsPromptFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"test.md", true},
		{"test.txt", true},
		{"test.MD", true},
		{"test.TXT", true},
		{"test.go", false},
		{"test.py", false},
		{"test", false},
	}

	for _, test := range tests {
		result := isPromptFile(test.filename)
		if result != test.expected {
			t.Errorf("isPromptFile(%s) = %v, expected %v", test.filename, result, test.expected)
		}
	}
}

func TestDefaultManagerGetWithSameNamePrompts(t *testing.T) {
	// Create fake filesystem with same-named prompts at different levels
	fs := filesystem.NewFakeFilesystem()
	fs.MapFS["user/prompts/test.md"] = &fstest.MapFile{
		Data: []byte("User level content"),
		Mode: 0644,
	}
	fs.MapFS["project/prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Project level content"),
		Mode: 0644,
	}
	fs.MapFS["prompts/test.md"] = &fstest.MapFile{
		Data: []byte("Directory level content"),
		Mode: 0644,
	}

	// Create fake resolver with correct hierarchy: directory > project > user
	// This simulates what the real resolver does
	resolver := NewFakeLocationResolver()
	resolver.Locations = []PromptLocation{
		{Type: "directory", Path: "prompts"},
		{Type: "project", Path: "project/prompts"},
		{Type: "user", Path: "user/prompts"},
	}

	manager := NewDefaultManager(fs, resolver)

	// Get should return the highest priority prompt (directory level)
	prompt, err := manager.Get("test")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	// Should return the highest priority prompt (directory level)
	if prompt.Content != "Directory level content" {
		t.Errorf("Expected directory level content, got '%s'", prompt.Content)
	}

	if prompt.Source != "directory" {
		t.Errorf("Expected source 'directory', got '%s'", prompt.Source)
	}
}

func TestRemoveExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"test.md", "test"},
		{"test.txt", "test"},
		{"file.name.md", "file.name"},
		{"noextension", "noextension"},
	}

	for _, test := range tests {
		result := removeExtension(test.filename)
		if result != test.expected {
			t.Errorf("removeExtension(%s) = %s, expected %s", test.filename, result, test.expected)
		}
	}
}
