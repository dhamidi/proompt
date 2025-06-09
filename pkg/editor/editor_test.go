package editor

import (
	"testing"
)

func TestFakeEditor_Edit(t *testing.T) {
	editor := NewFakeEditor()

	t.Run("records edited files", func(t *testing.T) {
		filepath := "/path/to/file.txt"
		err := editor.Edit(filepath)
		if err != nil {
			t.Errorf("Edit() error = %v", err)
			return
		}

		if len(editor.EditedFiles) != 1 {
			t.Errorf("EditedFiles length = %d, want 1", len(editor.EditedFiles))
			return
		}

		if editor.EditedFiles[0] != filepath {
			t.Errorf("EditedFiles[0] = %q, want %q", editor.EditedFiles[0], filepath)
		}
	})

	t.Run("records multiple files", func(t *testing.T) {
		editor := NewFakeEditor()
		files := []string{"file1.txt", "file2.txt", "file3.txt"}

		for _, file := range files {
			err := editor.Edit(file)
			if err != nil {
				t.Errorf("Edit(%s) error = %v", file, err)
				return
			}
		}

		if len(editor.EditedFiles) != len(files) {
			t.Errorf("EditedFiles length = %d, want %d", len(editor.EditedFiles), len(files))
			return
		}

		for i, file := range files {
			if editor.EditedFiles[i] != file {
				t.Errorf("EditedFiles[%d] = %q, want %q", i, editor.EditedFiles[i], file)
			}
		}
	})

	t.Run("with WriteContent function", func(t *testing.T) {
		editor := NewFakeEditor()
		contentWritten := false
		expectedContent := []byte("test content")

		editor.WriteContent = func(path string) []byte {
			contentWritten = true
			if path != "/test/file.txt" {
				t.Errorf("WriteContent called with path %q, want %q", path, "/test/file.txt")
			}
			return expectedContent
		}

		err := editor.Edit("/test/file.txt")
		if err != nil {
			t.Errorf("Edit() error = %v", err)
			return
		}

		if !contentWritten {
			t.Error("WriteContent function was not called")
		}

		if len(editor.EditedFiles) != 1 || editor.EditedFiles[0] != "/test/file.txt" {
			t.Errorf("EditedFiles = %v, want [%q]", editor.EditedFiles, "/test/file.txt")
		}
	})
}

func TestRealEditor_Interface(t *testing.T) {
	// Test that RealEditor implements Editor interface
	var _ Editor = &RealEditor{}
}

func TestFakeEditor_Interface(t *testing.T) {
	// Test that FakeEditor implements Editor interface
	var _ Editor = &FakeEditor{}
}

func TestNewRealEditor(t *testing.T) {
	command := "vim"
	editor := NewRealEditor(command)

	if editor.Command != command {
		t.Errorf("NewRealEditor() Command = %q, want %q", editor.Command, command)
	}
}

func TestNewFakeEditor(t *testing.T) {
	editor := NewFakeEditor()

	if editor.EditedFiles == nil {
		t.Error("NewFakeEditor() EditedFiles should not be nil")
	}

	if len(editor.EditedFiles) != 0 {
		t.Errorf("NewFakeEditor() EditedFiles length = %d, want 0", len(editor.EditedFiles))
	}

	if editor.WriteContent != nil {
		t.Error("NewFakeEditor() WriteContent should be nil by default")
	}
}
