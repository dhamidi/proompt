package editor

// Editor interface abstracts editor invocation
type Editor interface {
	Edit(filepath string) error
}

// RealEditor invokes the system editor
type RealEditor struct {
	Command string
}

// NewRealEditor creates a new RealEditor with the given command
func NewRealEditor(command string) *RealEditor {
	return &RealEditor{
		Command: command,
	}
}

// FakeEditor simulates editor behavior for testing
type FakeEditor struct {
	EditedFiles  []string
	WriteContent func(path string) []byte
}

// NewFakeEditor creates a new FakeEditor
func NewFakeEditor() *FakeEditor {
	return &FakeEditor{
		EditedFiles: make([]string, 0),
	}
}
