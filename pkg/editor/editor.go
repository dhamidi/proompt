package editor

import (
	"os"
	"os/exec"
)

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

// Edit invokes the editor on the given file
func (re *RealEditor) Edit(filepath string) error {
	cmd := exec.Command(re.Command, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

// Edit simulates editing by recording the file and optionally modifying it
func (fe *FakeEditor) Edit(filepath string) error {
	fe.EditedFiles = append(fe.EditedFiles, filepath)
	
	// If WriteContent function is provided, simulate user modifications
	if fe.WriteContent != nil {
		content := fe.WriteContent(filepath)
		// For FakeEditor, we would need access to the filesystem to actually write
		// The calling code should handle writing the content back to the filesystem
		_ = content
	}
	
	return nil
}
