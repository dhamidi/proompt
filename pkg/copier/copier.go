package copier

import (
	"fmt"
	"os/exec"
	"strings"
)

// Copier interface abstracts copying content to clipboard
type Copier interface {
	Copy(content string) error
}

// RealCopier uses an external copy command
type RealCopier struct {
	Command string
}

// NewRealCopier creates a new RealCopier with the given command
func NewRealCopier(command string) *RealCopier {
	return &RealCopier{
		Command: command,
	}
}

// Copy executes the copy command with the provided content
func (c *RealCopier) Copy(content string) error {
	if c.Command == "" {
		return nil // No-op if no command is configured
	}

	cmd := exec.Command("sh", "-c", c.Command)
	cmd.Stdin = strings.NewReader(content)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("copy command failed: %w", err)
	}

	return nil
}

// FakeCopier simulates copy behavior for testing
type FakeCopier struct {
	CopiedContent []string
	ShouldFail    bool
}

// NewFakeCopier creates a new FakeCopier
func NewFakeCopier() *FakeCopier {
	return &FakeCopier{
		CopiedContent: make([]string, 0),
	}
}

// Copy records the content for testing verification
func (c *FakeCopier) Copy(content string) error {
	if c.ShouldFail {
		return fmt.Errorf("copy failed")
	}

	c.CopiedContent = append(c.CopiedContent, content)
	return nil
}

// LastCopied returns the most recently copied content
func (c *FakeCopier) LastCopied() string {
	if len(c.CopiedContent) == 0 {
		return ""
	}
	return c.CopiedContent[len(c.CopiedContent)-1]
}

// CopyCount returns the number of copy operations performed
func (c *FakeCopier) CopyCount() int {
	return len(c.CopiedContent)
}
