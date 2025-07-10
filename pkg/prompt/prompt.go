package prompt

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/dhamidi/proompt/pkg/filesystem"
	"github.com/dhamidi/proompt/pkg/picker"
)

var (
	// ErrPromptNotFound is returned when a prompt cannot be found
	ErrPromptNotFound = errors.New("prompt not found")
	// ErrInvalidLocation is returned when an invalid location is specified
	ErrInvalidLocation = errors.New("invalid location")
)

// Manager interface handles prompt management
type Manager interface {
	List() ([]PromptInfo, error)
	Get(name string) (*PromptInfo, error)
	Create(name, content, location string) error
	Delete(name string) error
	GetAllForPicker() ([]picker.PickerItem, error)
}

// PromptInfo contains information about a prompt
type PromptInfo struct {
	Name    string
	Content string
	Source  string
	Path    string
}

// DefaultManager implements prompt management
type DefaultManager struct {
	Filesystem filesystem.Filesystem
	Resolver   LocationResolver
}

// NewDefaultManager creates a new DefaultManager
func NewDefaultManager(fs filesystem.Filesystem, resolver LocationResolver) *DefaultManager {
	return &DefaultManager{
		Filesystem: fs,
		Resolver:   resolver,
	}
}

// List returns all available prompts from all locations
func (m *DefaultManager) List() ([]PromptInfo, error) {
	locations, err := m.Resolver.GetPromptPaths()
	if err != nil {
		return nil, err
	}

	var prompts []PromptInfo
	seenPaths := make(map[string]bool) // Track absolute paths to avoid duplicates
	seenNames := make(map[string]bool) // Track prompt names to respect hierarchy
	
	for _, location := range locations {
		files, err := m.Filesystem.ReadDir(location.Path)
		if err != nil {
			continue // Skip locations that can't be read
		}

		for _, file := range files {
			if !file.IsDir() && isPromptFile(file.Name()) {
				fullPath := location.Path + "/" + file.Name()
				
				// Convert to absolute path for deduplication
				absPath, err := filepath.Abs(fullPath)
				if err != nil {
					absPath = fullPath // Fallback to original path
				}
				
				// Skip if we've already seen this file path
				if seenPaths[absPath] {
					continue
				}
				seenPaths[absPath] = true
				
				promptName := removeExtension(file.Name())
				
				// Skip if we've already seen this prompt name (hierarchy respect)
				if seenNames[promptName] {
					continue
				}
				seenNames[promptName] = true
				
				content, err := m.Filesystem.ReadFile(fullPath)
				if err != nil {
					continue // Skip files that can't be read
				}

				prompts = append(prompts, PromptInfo{
					Name:    promptName,
					Content: string(content),
					Source:  location.Type,
					Path:    fullPath,
				})
			}
		}
	}

	return prompts, nil
}

// Get returns a specific prompt by name
func (m *DefaultManager) Get(name string) (*PromptInfo, error) {
	prompts, err := m.List()
	if err != nil {
		return nil, err
	}

	for _, prompt := range prompts {
		if prompt.Name == name {
			return &prompt, nil
		}
	}

	return nil, ErrPromptNotFound
}

// Create creates a new prompt at the specified location
func (m *DefaultManager) Create(name, content, location string) error {
	locations, err := m.Resolver.GetPromptPaths()
	if err != nil {
		return err
	}

	var targetPath string
	for _, loc := range locations {
		if loc.Type == location {
			targetPath = loc.Path
			break
		}
	}

	if targetPath == "" {
		return ErrInvalidLocation
	}

	// Ensure the directory exists
	err = m.Filesystem.MkdirAll(targetPath, 0755)
	if err != nil {
		return err
	}

	filename := name + ".md"
	filepath := targetPath + "/" + filename
	
	return m.Filesystem.WriteFile(filepath, []byte(content), 0644)
}

// Delete removes a prompt by name
func (m *DefaultManager) Delete(name string) error {
	prompt, err := m.Get(name)
	if err != nil {
		return err
	}

	return m.Filesystem.Remove(prompt.Path)
}

// GetAllForPicker returns all prompts formatted for picker interface
func (m *DefaultManager) GetAllForPicker() ([]picker.PickerItem, error) {
	prompts, err := m.List()
	if err != nil {
		return nil, err
	}

	var items []picker.PickerItem
	for _, prompt := range prompts {
		items = append(items, picker.PickerItem{
			Name:   prompt.Name,
			Source: prompt.Source,
			Path:   prompt.Path,
		})
	}

	return items, nil
}

// isPromptFile checks if a file is a valid prompt file based on extension
func isPromptFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".txt"
}

// removeExtension removes the file extension from a filename
func removeExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
