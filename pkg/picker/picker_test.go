package picker

import (
	"testing"
)

func TestFakePickerPick(t *testing.T) {
	tests := []struct {
		name          string
		items         []PickerItem
		selectedIndex int
		wantError     bool
		wantName      string
	}{
		{
			name: "valid selection",
			items: []PickerItem{
				{Name: "test1", Source: "directory", Path: "path1"},
				{Name: "test2", Source: "project", Path: "path2"},
			},
			selectedIndex: 0,
			wantError:     false,
			wantName:      "test1",
		},
		{
			name:          "no items",
			items:         []PickerItem{},
			selectedIndex: 0,
			wantError:     true,
		},
		{
			name: "invalid index",
			items: []PickerItem{
				{Name: "test1", Source: "directory", Path: "path1"},
			},
			selectedIndex: 5,
			wantError:     true,
		},
		{
			name: "negative index",
			items: []PickerItem{
				{Name: "test1", Source: "directory", Path: "path1"},
			},
			selectedIndex: -1,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			picker := &FakePicker{
				SelectedIndex: tt.selectedIndex,
				Selections:    make([]PickerItem, 0),
			}

			result, err := picker.Pick(tt.items)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Name != tt.wantName {
				t.Errorf("expected name %s, got %s", tt.wantName, result.Name)
			}

			if len(picker.Selections) != 1 {
				t.Errorf("expected 1 selection recorded, got %d", len(picker.Selections))
			}
		})
	}
}

func TestPickerItem(t *testing.T) {
	item := PickerItem{
		Name:   "test",
		Source: "directory", 
		Path:   "/path/to/test",
	}

	if item.Name != "test" {
		t.Errorf("expected name 'test', got %s", item.Name)
	}
	if item.Source != "directory" {
		t.Errorf("expected source 'directory', got %s", item.Source)
	}
	if item.Path != "/path/to/test" {
		t.Errorf("expected path '/path/to/test', got %s", item.Path)
	}
}

func TestNewFakePicker(t *testing.T) {
	picker := NewFakePicker()
	
	if picker == nil {
		t.Error("expected non-nil picker")
	}
	
	if picker.Selections == nil {
		t.Error("expected initialized selections slice")
	}
	
	if len(picker.Selections) != 0 {
		t.Errorf("expected empty selections slice, got length %d", len(picker.Selections))
	}
}

func TestNewRealPicker(t *testing.T) {
	command := "fzf"
	picker := NewRealPicker(command)
	
	if picker == nil {
		t.Error("expected non-nil picker")
	}
	
	if picker.Command != command {
		t.Errorf("expected command %s, got %s", command, picker.Command)
	}
}
