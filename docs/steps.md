# Implementation Steps for Proompt

## Overview
This document provides detailed implementation steps for building the Proompt CLI tool. Each step focuses on a specific domain function with clear interfaces for testability.

## Phase 1: Core Infrastructure

### Step 1: Project Structure Setup - DONE
Create the basic Go project structure:

```
cmd/proompt/
  main.go                 # CLI entry point
pkg/
  filesystem/
    filesystem.go         # Filesystem interface and implementations
  editor/
    editor.go            # Editor interface and implementations
  picker/
    picker.go            # Picker interface and implementations
  prompt/
    prompt.go            # Prompt management
    parser.go            # Placeholder parsing
    resolver.go          # Prompt location resolution
  config/
    config.go            # Configuration management
```

**Technical Details:**
- Use cobra CLI framework (already in dependencies)
- Each package contains its domain-specific interface
- Implement both fake and real implementations for each interface

### Step 2: Filesystem Interface (`pkg/filesystem/filesystem.go`) - DONE
Leverage Go's standard `fs.FS` with write operations:

```go
// ReadFS combines standard fs interfaces for reading
type ReadFS interface {
    fs.ReadFileFS
    fs.StatFS
    fs.ReadDirFS
}

// WriteFS provides write operations
type WriteFS interface {
    WriteFile(path string, data []byte, perm os.FileMode) error
    MkdirAll(path string, perm os.FileMode) error
    Remove(path string) error
    TempFile(dir, pattern string) (*os.File, error)
}

// Filesystem combines read and write operations
type Filesystem interface {
    ReadFS
    WriteFS
    Getwd() (string, error)
    UserConfigDir() (string, error)
}

// RealFilesystem wraps os.DirFS for reads + standard library for writes
type RealFilesystem struct {
    readFS fs.FS // will be os.DirFS(cwd)
}

// FakeFilesystem uses testing/fstest.MapFS + in-memory writes
type FakeFilesystem struct {
    fstest.MapFS
    cwd string
}
```

**Implementation Notes:**
- Use `os.DirFS(cwd)` for real read operations
- Use `testing/fstest.MapFS` for fake read operations  
- Only custom implement write operations (which fs.FS doesn't provide)
- Replace `Walk` with `fs.WalkDir` function from standard library

### Step 3: Editor Interface (`pkg/editor/editor.go`) - DONE
Abstract editor invocation:

```go
type Editor interface {
    Edit(filepath string) error
}

type RealEditor struct {
    command string // $EDITOR or fallback
}

type FakeEditor struct {
    editedFiles []string
    writeContent func(path string) []byte // for simulating user input
}
```

**Implementation Notes:**
- `RealEditor` uses `os/exec` to invoke `$EDITOR`
- `FakeEditor` records invocations and can simulate file modifications
- Handle `$EDITOR` environment variable with sensible defaults

### Step 4: Picker Interface (`pkg/picker/picker.go`) - DONE
Abstract prompt selection:

```go
type Picker interface {
    Pick(items []PickerItem) (PickerItem, error)
}

type PickerItem struct {
    Name   string
    Source string // "directory", "project", "project-local", "user"
    Path   string
}

type RealPicker struct {
    command string // $PROOMPT_PICKER or "fzf"
}

type FakePicker struct {
    selectedIndex int
    selections    []PickerItem // for simulating user selections
}
```

**Implementation Notes:**
- `RealPicker` invokes external picker (fzf by default)
- `FakePicker` returns predetermined selections for testing
- Handle `$PROOMPT_PICKER` environment variable

## Phase 2: Prompt Management

### Step 5: Prompt Location Resolution (`pkg/prompt/resolver.go`) - DONE
Implement the three-level prompt hierarchy:

```go
type LocationResolver interface {
    GetPromptPaths() ([]PromptLocation, error)
}

type PromptLocation struct {
    Type string // "directory", "project", "project-local", "user"
    Path string
}

type DefaultLocationResolver struct {
    fs Filesystem
}

type FakeLocationResolver struct {
    locations []PromptLocation
}
```

**Implementation Logic:**
1. Directory level: `./prompts/`
2. Project level: Find `.git` or `prompts/` folder upward, use `prompts/` subdirectory
3. Project-local level: Same as project but in `.git/info/prompts/`
4. User level: `$XDG_CONFIG_HOME/proompt/prompts/`

### Step 6: Prompt Parser (`pkg/prompt/parser.go`) - DONE
Parse placeholder syntax from prompt content:

```go
type Parser interface {
    ParsePlaceholders(content string) ([]Placeholder, error)
    SubstitutePlaceholders(content string, values map[string]string) string
}

type Placeholder struct {
    Name         string
    DefaultValue string
    HasDefault   bool
}

type DefaultParser struct{}
type FakeParser struct {
    placeholders []Placeholder
}
```

**Implementation Details:**
- Parse `${VAR}` and `${VAR:-default}` syntax
- Handle `$$` escaping for literal `$`
- Use regex for robust parsing: `\$\{([^}:]+)(?::-([^}]*))?\}`

### Step 7: Prompt Manager (`pkg/prompt/prompt.go`) - DONE
Main prompt management functionality:

```go
type Manager interface {
    List() ([]PromptInfo, error)
    Get(name string) (*PromptInfo, error)
    Create(name, content, location string) error
    Delete(name string) error
    GetAllForPicker() ([]picker.PickerItem, error)
}

type PromptInfo struct {
    Name     string
    Content  string
    Source   string
    Path     string
}

type DefaultManager struct {
    fs       Filesystem
    resolver LocationResolver
}
```

## Phase 3: CLI Commands

### Step 8: List Command (`cmd/proompt/list.go`) - DONE
```go
func listCmd(manager prompt.Manager) *cobra.Command {
    return &cobra.Command{
        Use:   "list",
        Short: "List all available prompts",
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation calls manager.List()
        },
    }
}
```

### Step 9: Show Command (`cmd/proompt/show.go`) - DONE
```go
func showCmd(manager prompt.Manager) *cobra.Command {
    return &cobra.Command{
        Use:   "show <name>",
        Short: "Show a specific prompt",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation calls manager.Get(name)
        },
    }
}
```

### Step 10: Edit Command (`cmd/proompt/edit.go`) - DONE
```go
func editCmd(manager prompt.Manager, picker picker.Picker, editor editor.Editor) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "edit [name]",
        Short: "Edit a prompt",
        Run: func(cmd *cobra.Command, args []string) {
            // If no name provided, use picker
            // Handle --directory, --project, --project-local, --user flags
        },
    }
    
    cmd.Flags().Bool("directory", false, "Create in directory level")
    cmd.Flags().Bool("project", false, "Create in project level")
    cmd.Flags().Bool("project-local", false, "Create in project-local level")
    cmd.Flags().Bool("user", false, "Create in user level")
    
    return cmd
}
```

### Step 11: Remove Command (`cmd/proompt/rm.go`) - DONE
```go
func rmCmd(manager prompt.Manager, picker picker.Picker) *cobra.Command {
    return &cobra.Command{
        Use:   "rm [name]",
        Short: "Remove a prompt",
        Run: func(cmd *cobra.Command, args []string) {
            // If no name provided, use picker
            // Implementation calls manager.Delete(name)
        },
    }
}
```

## Phase 4: Core Pick Functionality

### Step 12: Pick Command Implementation (`cmd/proompt/pick.go`) - DONE
The heart of the application:

```go
func pickCmd(
    manager prompt.Manager,
    picker picker.Picker,
    editor editor.Editor,
    parser prompt.Parser,
    fs filesystem.Filesystem,
) *cobra.Command {
    return &cobra.Command{
        Use:   "pick",
        Short: "Pick and process a prompt",
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation of the 5-step process from spec
        },
    }
}
```

**Implementation Steps:**
1. Get all prompts using `manager.GetAllForPicker()`
2. Use `picker.Pick()` to let user select
3. Parse selected prompt with `parser.ParsePlaceholders()`
4. Create temporary file with placeholders and defaults
5. Invoke `editor.Edit()` on temp file
6. Read back values and substitute with `parser.SubstitutePlaceholders()`
7. Output final prompt to stdout

### Step 13: Temporary File Generation - DONE
Create the placeholder editing experience:

```go
func generatePlaceholderFile(placeholders []Placeholder, originalContent string) string {
    var buf strings.Builder
    buf.WriteString("\n")
    
    for _, p := range placeholders {
        buf.WriteString(fmt.Sprintf("%s=%s\n", p.Name, p.DefaultValue))
    }
    
    buf.WriteString("### Lines starting with # are ignored\n")
    buf.WriteString("# Save empty file to abort\n")
    buf.WriteString("### Full prompt:\n")
    
    for _, line := range strings.Split(originalContent, "\n") {
        buf.WriteString(fmt.Sprintf("# %s\n", line))
    }
    
    return buf.String()
}
```

### Step 14: Value Parsing from Editor - DONE
Parse edited values from the temporary file:

```go
func parseEditedValues(content string) (map[string]string, error) {
    values := make(map[string]string)
    
    for _, line := range strings.Split(content, "\n") {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            values[parts[0]] = parts[1]
        }
    }
    
    return values, nil
}
```

## Phase 5: Configuration and Integration

### Step 15: Configuration Management (`pkg/config/config.go`) - DONE
Handle environment variables and defaults:

```go
type Config struct {
    Editor string
    Picker string
}

func Load() *Config {
    return &Config{
        Editor: getEnv("EDITOR", "nano"),
        Picker: getEnv("PROOMPT_PICKER", "fzf"),
    }
}
```

### Step 16: Main CLI Integration (`cmd/proompt/main.go`) - DONE
Wire everything together:

```go
func main() {
    config := config.Load()
    fs := &filesystem.RealFilesystem{}
    resolver := &prompt.DefaultLocationResolver{Filesystem: fs}
    manager := &prompt.DefaultManager{
        Filesystem: fs,
        Resolver:   resolver,
    }
    picker := &picker.RealPicker{Command: config.Picker}
    editor := &editor.RealEditor{Command: config.Editor}
    parser := &prompt.DefaultParser{}

    rootCmd := &cobra.Command{Use: "proompt"}
    rootCmd.AddCommand(
        listCmd(manager),
        showCmd(manager),
        editCmd(manager, picker, editor),
        rmCmd(manager, picker),
        pickCmd(manager, picker, editor, parser, fs),
    )

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Phase 6: Testing Strategy

### Step 17: Unit Tests for Each Component - DONE
- **Filesystem tests**: Verify fake implementation matches real behavior
- **Parser tests**: Test placeholder parsing and substitution with various edge cases
- **Manager tests**: Test prompt discovery and management using fake filesystem
- **Command tests**: Test CLI commands using all fake implementations

### Step 18: Integration Tests
- Test complete pick workflow end-to-end with fake implementations
- Test error handling (missing editor, picker failures, etc.)
- Test different prompt hierarchy scenarios

### Step 19: Error Handling
- Graceful handling of missing picker/editor commands
- Clear error messages for invalid prompt names
- Proper handling of file permission issues
- Validation of placeholder syntax

## Testing Examples

```go
func TestPickWorkflow(t *testing.T) {
    fs := &filesystem.FakeFilesystem{
        files: map[string][]byte{
            "prompts/test.md": []byte("Hello ${NAME:-World}!"),
        },
    }
    
    picker := &picker.FakePicker{
        selections: []picker.PickerItem{{Name: "test", Path: "prompts/test.md"}},
    }
    
    editor := &editor.FakeEditor{
        writeContent: func(path string) []byte {
            return []byte("NAME=Proompt\n")
        },
    }
    
    // Test the complete workflow
}
```

This implementation plan provides a solid foundation for a junior developer to build the Proompt tool with proper separation of concerns, comprehensive testing, and clean interfaces.
