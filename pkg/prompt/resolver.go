package prompt

// LocationResolver interface handles prompt location resolution
type LocationResolver interface {
	GetPromptPaths() ([]PromptLocation, error)
}

// PromptLocation represents a location where prompts can be found
type PromptLocation struct {
	Type string // "directory", "project", "project-local", "user"
	Path string
}

// DefaultLocationResolver implements the three-level prompt hierarchy
type DefaultLocationResolver struct {
	Filesystem Filesystem
}

// NewDefaultLocationResolver creates a new DefaultLocationResolver
func NewDefaultLocationResolver(fs Filesystem) *DefaultLocationResolver {
	return &DefaultLocationResolver{
		Filesystem: fs,
	}
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
