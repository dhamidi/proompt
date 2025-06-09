package filesystem

import (
	"os"
	"testing"
	"testing/fstest"
)

func TestFakeFilesystem_ReadOperations(t *testing.T) {
	fs := NewFakeFilesystem()
	
	// Add test files
	fs.MapFS["test.txt"] = &fstest.MapFile{
		Data: []byte("hello world"),
		Mode: 0644,
	}
	fs.MapFS["dir/file.txt"] = &fstest.MapFile{
		Data: []byte("nested file"),
		Mode: 0644,
	}

	t.Run("ReadFile", func(t *testing.T) {
		data, err := fs.ReadFile("test.txt")
		if err != nil {
			t.Errorf("ReadFile() error = %v", err)
			return
		}
		expected := "hello world"
		if string(data) != expected {
			t.Errorf("ReadFile() = %q, want %q", string(data), expected)
		}
	})

	t.Run("ReadFile non-existent", func(t *testing.T) {
		_, err := fs.ReadFile("nonexistent.txt")
		if err == nil {
			t.Error("ReadFile() expected error for non-existent file")
		}
	})

	t.Run("Stat", func(t *testing.T) {
		info, err := fs.Stat("test.txt")
		if err != nil {
			t.Errorf("Stat() error = %v", err)
			return
		}
		if info.Name() != "test.txt" {
			t.Errorf("Stat() name = %q, want %q", info.Name(), "test.txt")
		}
		if info.Mode() != 0644 {
			t.Errorf("Stat() mode = %v, want %v", info.Mode(), 0644)
		}
	})
}

func TestFakeFilesystem_WriteOperations(t *testing.T) {
	fs := NewFakeFilesystem()

	t.Run("WriteFile", func(t *testing.T) {
		data := []byte("test content")
		err := fs.WriteFile("test.txt", data, 0644)
		if err != nil {
			t.Errorf("WriteFile() error = %v", err)
			return
		}

		// Verify file was written
		readData, err := fs.ReadFile("test.txt")
		if err != nil {
			t.Errorf("ReadFile() after WriteFile() error = %v", err)
			return
		}
		if string(readData) != string(data) {
			t.Errorf("ReadFile() after WriteFile() = %q, want %q", string(readData), string(data))
		}
	})

	t.Run("MkdirAll", func(t *testing.T) {
		err := fs.MkdirAll("path/to/dir", 0755)
		if err != nil {
			t.Errorf("MkdirAll() error = %v", err)
		}
		// MkdirAll in FakeFilesystem is a no-op, so no verification needed
	})

	t.Run("Remove existing file", func(t *testing.T) {
		// First create a file
		fs.WriteFile("to_remove.txt", []byte("content"), 0644)
		
		// Verify it exists
		_, err := fs.ReadFile("to_remove.txt")
		if err != nil {
			t.Errorf("File should exist before removal")
			return
		}

		// Remove it
		err = fs.Remove("to_remove.txt")
		if err != nil {
			t.Errorf("Remove() error = %v", err)
			return
		}

		// Verify it's gone
		_, err = fs.ReadFile("to_remove.txt")
		if err == nil {
			t.Error("File should not exist after removal")
		}
	})

	t.Run("Remove non-existent file", func(t *testing.T) {
		err := fs.Remove("nonexistent.txt")
		if err != os.ErrNotExist {
			t.Errorf("Remove() error = %v, want %v", err, os.ErrNotExist)
		}
	})

	t.Run("TempFile", func(t *testing.T) {
		file, err := fs.TempFile("", "temp-*.txt")
		if err != nil {
			t.Errorf("TempFile() error = %v", err)
		}
		// For FakeFilesystem, TempFile returns nil file, which is acceptable for testing
		if file != nil {
			file.Close()
		}
	})
}

func TestFakeFilesystem_Configuration(t *testing.T) {
	fs := NewFakeFilesystem()

	t.Run("Default Getwd", func(t *testing.T) {
		cwd, err := fs.Getwd()
		if err != nil {
			t.Errorf("Getwd() error = %v", err)
			return
		}
		expected := "/"
		if cwd != expected {
			t.Errorf("Getwd() = %q, want %q", cwd, expected)
		}
	})

	t.Run("SetCwd", func(t *testing.T) {
		newCwd := "/custom/path"
		fs.SetCwd(newCwd)
		
		cwd, err := fs.Getwd()
		if err != nil {
			t.Errorf("Getwd() error = %v", err)
			return
		}
		if cwd != newCwd {
			t.Errorf("Getwd() after SetCwd() = %q, want %q", cwd, newCwd)
		}
	})

	t.Run("Default UserConfigDir", func(t *testing.T) {
		dir, err := fs.UserConfigDir()
		if err != nil {
			t.Errorf("UserConfigDir() error = %v", err)
			return
		}
		expected := "/home/user/.config"
		if dir != expected {
			t.Errorf("UserConfigDir() = %q, want %q", dir, expected)
		}
	})

	t.Run("SetUserConfigDir", func(t *testing.T) {
		newDir := "/custom/config"
		fs.SetUserConfigDir(newDir)
		
		dir, err := fs.UserConfigDir()
		if err != nil {
			t.Errorf("UserConfigDir() error = %v", err)
			return
		}
		if dir != newDir {
			t.Errorf("UserConfigDir() after SetUserConfigDir() = %q, want %q", dir, newDir)
		}
	})
}

// TestRealFilesystem tests the interface compliance
func TestRealFilesystem_Interface(t *testing.T) {
	// Test that RealFilesystem implements Filesystem interface
	var _ Filesystem = &RealFilesystem{}
}

// TestFakeFilesystem tests the interface compliance
func TestFakeFilesystem_Interface(t *testing.T) {
	// Test that FakeFilesystem implements Filesystem interface
	var _ Filesystem = &FakeFilesystem{}
}
