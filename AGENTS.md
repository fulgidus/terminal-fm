# Agent Guidelines for Terminal.FM

## Build/Test Commands
```bash
go build -o terminal-fm ./cmd/server              # Build binary
go run ./cmd/server --dev --port 2222             # Run dev server
go test ./...                                     # Run all tests
go test -v -run TestPlayerVolume ./pkg/services/player  # Run single test
go test -race -coverprofile=coverage.out ./...    # Test with race detector + coverage
golangci-lint run                                 # Lint code
gofmt -s -w . && goimports -w .                   # Format code
```

## Code Style (Go 1.21+)
- **Formatting**: Use `gofmt`/`goimports` (tabs, not spaces)
- **Naming**: Exported=`CamelCase`, unexported=`camelCase`, acronyms=`HTTPClient`
- **Imports**: Standard lib → external → internal, grouped with blank lines
- **Errors**: Return errors (don't panic), wrap with `fmt.Errorf("context: %w", err)`
- **Comments**: Exported names need full sentence: `// Station represents...`
- **Testing**: Table-driven tests in `*_test.go`, mock external dependencies
- **Types**: Prefer interfaces for flexibility, use pointers for large structs

## Project Structure
`cmd/server/` → entry point | `pkg/` → reusable packages (ssh, ui, services, storage, i18n) | `internal/` → private code | `docs/` → documentation
