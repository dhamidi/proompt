package picker

// Picker interface abstracts prompt selection
type Picker interface {
	Pick(items []PickerItem) (PickerItem, error)
}

// PickerItem represents an item that can be selected
type PickerItem struct {
	Name   string
	Source string // "directory", "project", "project-local", "user"
	Path   string
}

// RealPicker uses an external picker tool
type RealPicker struct {
	Command string
}

// NewRealPicker creates a new RealPicker with the given command
func NewRealPicker(command string) *RealPicker {
	return &RealPicker{
		Command: command,
	}
}

// FakePicker simulates picker behavior for testing
type FakePicker struct {
	SelectedIndex int
	Selections    []PickerItem
}

// NewFakePicker creates a new FakePicker
func NewFakePicker() *FakePicker {
	return &FakePicker{
		Selections: make([]PickerItem, 0),
	}
}
