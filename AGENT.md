# AGENT.md - Development Guide for Proompt

## Build/Test/Lint Commands
- **Build**: `go build ./...` or `go build -o proompt ./cmd/proompt`
- **Test all**: `go test ./...`
- **Test single**: `go test ./pkg/parser` (replace with specific package)
- **Lint**: `golint ./...` or `go vet ./...`
- **Format**: `go fmt ./...`
- **Tidy deps**: `go mod tidy`

## Code Style Guidelines
- Use standard Go formatting (`gofmt`)
- Package names: lowercase, single word (e.g., `parser`, `picker`)
- Function names: CamelCase for exported, camelCase for unexported
- Error handling: always check errors, wrap with context using `fmt.Errorf`
- Imports: standard library first, then third-party, then local packages
- Types: prefer composition over inheritance, use interfaces sparingly
- Constants: ALL_CAPS with underscores for package-level constants

## Project Structure
- `cmd/proompt/`: main CLI application entry point
- `pkg/`: reusable packages (parser, picker, prompt management)
- `docs/`: project documentation
- Use cobra CLI framework (already in dependencies)

## Testing
- Use standard `testing` package
- Test files: `*_test.go` in same package
- Table-driven tests preferred for multiple test cases
