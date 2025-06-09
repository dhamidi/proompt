package prompt

import (
	"github.com/dhamidi/proompt/pkg/filesystem"
	"github.com/dhamidi/proompt/pkg/picker"
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
