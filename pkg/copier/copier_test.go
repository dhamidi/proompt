package copier

import (
	"testing"
)

func TestFakeCopier(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		shouldFail  bool
		wantError   bool
	}{
		{
			name:       "successful copy",
			content:    "test content",
			shouldFail: false,
			wantError:  false,
		},
		{
			name:       "empty content",
			content:    "",
			shouldFail: false,
			wantError:  false,
		},
		{
			name:       "multiline content",
			content:    "line1\nline2\nline3",
			shouldFail: false,
			wantError:  false,
		},
		{
			name:       "copy failure",
			content:    "test content",
			shouldFail: true,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copier := NewFakeCopier()
			copier.ShouldFail = tt.shouldFail

			err := copier.Copy(tt.content)

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

			if copier.LastCopied() != tt.content {
				t.Errorf("expected last copied %q, got %q", tt.content, copier.LastCopied())
			}

			if copier.CopyCount() != 1 {
				t.Errorf("expected copy count 1, got %d", copier.CopyCount())
			}
		})
	}
}

func TestFakeCopierMultipleCopies(t *testing.T) {
	copier := NewFakeCopier()

	contents := []string{"first", "second", "third"}
	for _, content := range contents {
		err := copier.Copy(content)
		if err != nil {
			t.Errorf("unexpected error copying %q: %v", content, err)
		}
	}

	if copier.CopyCount() != 3 {
		t.Errorf("expected copy count 3, got %d", copier.CopyCount())
	}

	if copier.LastCopied() != "third" {
		t.Errorf("expected last copied %q, got %q", "third", copier.LastCopied())
	}

	if len(copier.CopiedContent) != 3 {
		t.Errorf("expected 3 copied items, got %d", len(copier.CopiedContent))
	}
}

func TestNewFakeCopier(t *testing.T) {
	copier := NewFakeCopier()

	if copier == nil {
		t.Error("expected non-nil copier")
	}

	if copier.CopiedContent == nil {
		t.Error("expected initialized copied content slice")
	}

	if copier.CopyCount() != 0 {
		t.Errorf("expected copy count 0, got %d", copier.CopyCount())
	}

	if copier.LastCopied() != "" {
		t.Errorf("expected empty last copied, got %q", copier.LastCopied())
	}
}

func TestNewRealCopier(t *testing.T) {
	command := "pbcopy"
	copier := NewRealCopier(command)

	if copier == nil {
		t.Error("expected non-nil copier")
	}

	if copier.Command != command {
		t.Errorf("expected command %q, got %q", command, copier.Command)
	}
}

func TestRealCopierEmptyCommand(t *testing.T) {
	copier := NewRealCopier("")
	
	// Should not fail with empty command (no-op)
	err := copier.Copy("test content")
	if err != nil {
		t.Errorf("unexpected error with empty command: %v", err)
	}
}
