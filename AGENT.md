# AGENT.md - Development Guide for Proompt

## Project Overview
Proompt is a CLI tool for managing LLM prompts with placeholder substitution. It supports a 4-level hierarchy: directory, project, project-local, and user level prompts. The core functionality is the `pick` command which allows users to select prompts, fill in placeholders via editor, and output the processed prompt.

## Build/Test/Lint Commands
- **Build**: `go build ./...` or `go build -o proompt ./cmd/proompt`
- **Test all**: `go test ./...`
- **Test single**: `go test ./pkg/prompt` (replace with specific package)
- **Lint**: `golint ./...` or `go vet ./...`
- **Format**: `go fmt ./...`
- **Tidy deps**: `go mod tidy`
- **Run**: `./proompt <command>` or `go run ./cmd/proompt <command>`

## Dependencies
- **Go version**: 1.24.3
- **External dependencies**: 
  - `github.com/spf13/cobra` (CLI framework)
  - `github.com/spf13/pflag` (flag parsing)
  - `github.com/inconshreveable/mousetrap` (cobra dependency)

## Project Architecture

### Package Structure
- `cmd/proompt/`: CLI entry point and command definitions
  - `main.go`: Application bootstrap and dependency injection
  - `list.go`, `show.go`, `edit.go`, `rm.go`, `pick.go`: Command implementations
  - `integration_test.go`: End-to-end tests
- `pkg/config/`: Configuration management (environment variables)
- `pkg/filesystem/`: Filesystem abstraction with real and fake implementations
- `pkg/editor/`: Editor invocation abstraction  
- `pkg/picker/`: Selection picker abstraction (fzf integration)
- `pkg/copier/`: Clipboard copy functionality
- `pkg/prompt/`: Core prompt management
  - `prompt.go`: Prompt manager (CRUD operations)
  - `parser.go`: Placeholder parsing and substitution
  - `resolver.go`: Prompt location resolution (4-level hierarchy)

### Prompt Hierarchy (Priority Order)
1. **Directory level**: `./prompts/` (current directory)
2. **Project level**: `<project-root>/prompts/` (shared, committed to git)
3. **Project-local level**: `<project-root>/.git/info/prompts/` (local, git-ignored)
4. **User level**: `$XDG_CONFIG_HOME/proompt/prompts/`

### Key Interfaces
- `prompt.Manager`: Prompt CRUD operations
- `prompt.Parser`: Placeholder parsing/substitution  
- `prompt.LocationResolver`: Prompt location discovery
- `picker.Picker`: Interactive selection interface
- `editor.Editor`: Text editor invocation
- `filesystem.Filesystem`: File system operations

## Environment Variables
- `EDITOR`: Text editor for prompt editing (default: "nano")
- `PROOMPT_PICKER`: Selection picker command (default: "fzf")
- `PROOMPT_COPY_COMMAND`: Copy to clipboard command (default: "pbcopy")

## Code Style Guidelines
- Use standard Go formatting (`gofmt`)
- Package names: lowercase, single word (e.g., `parser`, `picker`)
- Function names: CamelCase for exported, camelCase for unexported
- Error handling: always check errors, wrap with context using `fmt.Errorf`
- Imports: standard library first, then third-party, then local packages
- Types: prefer composition over inheritance, use interfaces sparingly
- Constants: ALL_CAPS with underscores for package-level constants

## Testing Strategy
- Each package has both real and fake implementations for testability
- Use standard `testing` package
- Test files: `*_test.go` in same package
- Table-driven tests preferred for multiple test cases
- Integration tests in `cmd/proompt/integration_test.go`
- Fake implementations enable isolated unit testing

## Placeholder Syntax
- `${VAR}`: Simple placeholder
- `${VAR:-default}`: Placeholder with default value
- `$$`: Escape sequence for literal `$`
- Regex pattern: `\$\{([^}:]+)(?::-([^}]*))?\}`

## CLI Commands
- `proompt list`: List all available prompts with sources
- `proompt show <name>`: Display prompt content  
- `proompt edit [name]`: Edit prompt (uses picker if no name provided)
- `proompt rm [name]`: Remove prompt (uses picker if no name provided)
- `proompt pick`: Core workflow - select prompt, fill placeholders, output result

## Development Notes
- Project is feature-complete based on `docs/steps.md` (all steps marked DONE)
- Comprehensive test coverage across all packages
- Clean separation of concerns with dependency injection
- Uses cobra CLI framework throughout
- File operations abstracted for cross-platform compatibility
