# AGENTS.md - Development Guidelines for AI Agents

## Build/Test Commands
- **Build**: `go build .` or `go build -o anktui`
- **Run**: `go run .` or `go run main.go`
- **Test**: `go test ./...` (all packages) or `go test ./models` (specific package)
- **Single Test**: `go test -run TestFunctionName ./package`
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...` (built-in) or use `golangci-lint run`
- **Deps**: `go mod tidy` to clean up dependencies

## Code Style Guidelines
- **Imports**: Group standard library, then external packages, then local modules
- **Naming**: camelCase for variables/functions, PascalCase for exported types/functions
- **Error Handling**: Always handle errors explicitly, return errors from functions
- **Comments**: Use `//` for comments, document exported functions with proper Go doc format
- **Types**: Use interfaces for abstractions, prefer explicit error returns over panics
- **Line Length**: Prefer 100-120 characters max per line

## Project-Specific Notes
- AnkTUI: Terminal User Interface for Anki-style flashcard studying using Bubbletea
- Config stored in XDG_CONFIG_HOME/anktui/ or ~/.config/anktui/
- Data stored as JSON files in configurable directory (default: ~/.local/share/anktui/)
- Full-screen TUI with centered ASCII art title and responsive navigation