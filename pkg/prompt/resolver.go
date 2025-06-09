package prompt

import (
	"errors"
	"path/filepath"

	"github.com/dhamidi/proompt/pkg/filesystem"
)

// LocationResolver interface handles prompt location resolution
type LocationResolver interface {
	GetPromptPaths() ([]PromptLocation, error)
}

// PromptLocation represents a location where prompts can be found
type PromptLocation struct {
	Type string // "directory", "project", "project-local", "user"
	Path string
}

// DefaultLocationResolver implements the four-level prompt hierarchy
type DefaultLocationResolver struct {
	Filesystem filesystem.Filesystem
}

// NewDefaultLocationResolver creates a new DefaultLocationResolver
func NewDefaultLocationResolver(fs filesystem.Filesystem) *DefaultLocationResolver {
	return &DefaultLocationResolver{
		Filesystem: fs,
	}
}

// GetPromptPaths returns all prompt locations in order of precedence
func (r *DefaultLocationResolver) GetPromptPaths() ([]PromptLocation, error) {
	var locations []PromptLocation

	// 1. Directory level: ./prompts/
	if info, err := r.Filesystem.Stat("prompts"); err == nil && info.IsDir() {
		locations = append(locations, PromptLocation{
			Type: "directory",
			Path: "prompts",
		})
	}

	// 2. Project level: Find .git or prompts/ folder upward, use prompts/ subdirectory
	if projectPath, err := r.findProjectRoot(); err == nil {
		promptsPath := filepath.Join(projectPath, "prompts")
		if info, err := r.Filesystem.Stat(promptsPath); err == nil && info.IsDir() {
			locations = append(locations, PromptLocation{
				Type: "project",
				Path: promptsPath,
			})
		}
	}

	// 3. Project-local level: Same as project but in .git/info/prompts/
	if projectPath, err := r.findProjectRoot(); err == nil {
		gitInfoPrompts := filepath.Join(projectPath, ".git", "info", "prompts")
		if info, err := r.Filesystem.Stat(gitInfoPrompts); err == nil && info.IsDir() {
			locations = append(locations, PromptLocation{
				Type: "project-local",
				Path: gitInfoPrompts,
			})
		}
	}

	// 4. User level: $XDG_CONFIG_HOME/proompt/prompts/
	if userConfigDir, err := r.Filesystem.UserConfigDir(); err == nil {
		userPromptsPath := filepath.Join(userConfigDir, "proompt", "prompts")
		locations = append(locations, PromptLocation{
			Type: "user",
			Path: userPromptsPath,
		})
	}

	return locations, nil
}

// findProjectRoot searches upward for .git directory or prompts/ folder
func (r *DefaultLocationResolver) findProjectRoot() (string, error) {
	cwd, err := r.Filesystem.Getwd()
	if err != nil {
		return "", err
	}

	current := cwd
	for {
		// Check for .git directory
		gitPath := filepath.Join(current, ".git")
		if info, err := r.Filesystem.Stat(gitPath); err == nil && info.IsDir() {
			return current, nil
		}

		// Check for prompts directory
		promptsPath := filepath.Join(current, "prompts")
		if info, err := r.Filesystem.Stat(promptsPath); err == nil && info.IsDir() {
			return current, nil
		}

		// Move up one directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached root without finding project markers
			break
		}
		current = parent
	}

	return "", errors.New("project root not found")
}

// FakeLocationResolver simulates location resolution for testing
type FakeLocationResolver struct {
	Locations []PromptLocation
}

// NewFakeLocationResolver creates a new FakeLocationResolver
func NewFakeLocationResolver() *FakeLocationResolver {
	return &FakeLocationResolver{
		Locations: make([]PromptLocation, 0),
	}
}

// GetPromptPaths returns the predetermined locations for testing
func (r *FakeLocationResolver) GetPromptPaths() ([]PromptLocation, error) {
	return r.Locations, nil
}
