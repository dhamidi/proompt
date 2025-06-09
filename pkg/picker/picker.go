package picker

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

// Pick invokes the external picker tool and returns the selected item
func (p *RealPicker) Pick(items []PickerItem) (PickerItem, error) {
	if len(items) == 0 {
		return PickerItem{}, errors.New("no items to pick from")
	}

	// Create input for picker
	var input strings.Builder
	for _, item := range items {
		fmt.Fprintf(&input, "%s (%s)\n", item.Name, item.Source)
	}

	// Execute picker command
	cmd := exec.Command("sh", "-c", p.Command)
	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		return PickerItem{}, fmt.Errorf("picker command failed: %w", err)
	}

	// Parse the selected line
	selected := strings.TrimSpace(string(output))
	if selected == "" {
		return PickerItem{}, errors.New("no selection made")
	}

	// Find the matching item
	for _, item := range items {
		expectedLine := fmt.Sprintf("%s (%s)", item.Name, item.Source)
		if selected == expectedLine {
			return item, nil
		}
	}

	return PickerItem{}, fmt.Errorf("selected item not found: %s", selected)
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

// Pick returns a predetermined selection for testing
func (p *FakePicker) Pick(items []PickerItem) (PickerItem, error) {
	if len(items) == 0 {
		return PickerItem{}, errors.New("no items to pick from")
	}

	if p.SelectedIndex < 0 || p.SelectedIndex >= len(items) {
		return PickerItem{}, fmt.Errorf("invalid selection index: %d", p.SelectedIndex)
	}

	selected := items[p.SelectedIndex]
	p.Selections = append(p.Selections, selected)
	return selected, nil
}
