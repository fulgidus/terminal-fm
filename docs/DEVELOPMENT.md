# Development Guide

This guide will help you set up a local development environment for Terminal.FM.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Building and Running](#building-and-running)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Debugging](#debugging)
- [Common Tasks](#common-tasks)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### Required Software

- **Go**: 1.21 or higher
  ```bash
  # Check Go version
  go version
  ```

- **Git**: Latest stable
  ```bash
  git --version
  ```

- **mpv or ffplay**: For audio playback testing
  ```bash
  # Ubuntu/Debian
  sudo apt install mpv
  
  # macOS
  brew install mpv
  
  # Verify
  mpv --version
  ```

### Recommended Tools

- **golangci-lint**: For linting
  ```bash
  # Install
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  
  # Verify
  golangci-lint --version
  ```

- **Air**: For live reloading (optional)
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

- **SQLite**: For database inspection
  ```bash
  # Ubuntu/Debian
  sudo apt install sqlite3
  
  # macOS
  brew install sqlite
  ```

## Environment Setup

### 1. Clone Repository

```bash
git clone https://github.com/fulgidus/terminal-fm.git
cd terminal-fm
```

### 2. Install Dependencies

```bash
# Download Go modules
go mod download

# Verify dependencies
go mod verify
```

### 3. Set Up Development Environment

```bash
# Create development directories
mkdir -p tmp
mkdir -p logs

# Create development database
touch tmp/dev.db
```

### 4. Configure Environment Variables (Optional)

```bash
# Create .env file for development
cat > .env << EOF
TERMINAL_FM_PORT=2222
TERMINAL_FM_HOST=localhost
TERMINAL_FM_DB_PATH=./tmp/dev.db
TERMINAL_FM_LOG_LEVEL=debug
EOF
```

## Building and Running

### Quick Start

```bash
# Build and run in development mode
go run ./cmd/server --dev --port 2222

# In another terminal, connect
ssh localhost -p 2222
```

### Build Binary

```bash
# Development build (with debug symbols)
go build -o terminal-fm ./cmd/server

# Production build (optimized)
go build -ldflags="-s -w" -o terminal-fm ./cmd/server

# Run
./terminal-fm --port 2222 --dev
```

### Using Air (Live Reload)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with live reload
air

# Air will watch for file changes and automatically rebuild
```

**Air Configuration** (`.air.toml`):
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/terminal-fm ./cmd/server"
  bin = "tmp/terminal-fm --dev --port 2222"
  include_ext = ["go", "toml"]
  exclude_dir = ["tmp", "vendor", "docs"]
  delay = 1000
```

## Project Structure

```
terminal-fm/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── pkg/
│   ├── ssh/                  # SSH server implementation
│   │   ├── server.go
│   │   ├── auth.go
│   │   └── session.go
│   ├── ui/                   # TUI components (Bubbletea)
│   │   ├── model.go
│   │   ├── update.go
│   │   ├── view.go
│   │   └── components/
│   │       ├── browser.go
│   │       ├── player.go
│   │       └── ...
│   ├── services/             # Business logic
│   │   ├── radiobrowser/    # Radio Browser API client
│   │   ├── player/          # Audio player control
│   │   └── userprefs/       # User preferences
│   ├── storage/              # Database layer
│   │   ├── db.go
│   │   ├── models.go
│   │   └── migrations.go
│   └── i18n/                 # Internationalization
│       ├── i18n.go
│       └── locales/
│           ├── active.en.toml
│           └── active.it.toml
├── internal/                 # Private application code
├── docs/                     # Documentation
├── scripts/                  # Build and utility scripts
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Build automation
└── README.md                 # Project documentation
```

## Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feat/your-feature-name
```

### 2. Make Changes

Edit code, following [Go coding standards](../CONTRIBUTING.md#coding-standards).

### 3. Test Locally

```bash
# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run specific package tests
go test ./pkg/services/radiobrowser/...
```

### 4. Lint Code

```bash
# Run linter
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

### 5. Format Code

```bash
# Format all files
gofmt -s -w .

# Or use goimports (also organizes imports)
goimports -w .
```

### 6. Commit Changes

Follow [commit message guidelines](../CONTRIBUTING.md#commit-message-guidelines):

```bash
git add .
git commit -m "feat(player): add volume control"
```

### 7. Push and Create PR

```bash
git push origin feat/your-feature-name
```

Then create a Pull Request on GitHub.

## Testing

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestPlayerVolume ./pkg/services/player

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

```bash
# Run integration tests (requires build tag)
go test -tags=integration ./...

# Skip integration tests in regular test runs
go test -short ./...
```

### Test Structure

```go
// unit test example
func TestPlayerSetVolume(t *testing.T) {
    tests := []struct {
        name    string
        volume  int
        wantErr bool
    }{
        {"valid volume", 50, false},
        {"volume too high", 150, true},
        {"volume negative", -10, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := player.New()
            err := p.SetVolume(tt.volume)
            if (err != nil) != tt.wantErr {
                t.Errorf("SetVolume() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Debugging

### Debug Logging

```bash
# Run with debug logging
go run ./cmd/server --dev --log-level debug

# View logs in real-time
tail -f logs/terminal-fm.log
```

### Using Delve Debugger

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug ./cmd/server -- --dev --port 2222

# Set breakpoint
(dlv) break pkg/ui/model.go:42

# Continue execution
(dlv) continue

# Print variable
(dlv) print variableName

# Step through code
(dlv) next
(dlv) step
```

### VS Code Debug Configuration

`.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/server",
            "args": ["--dev", "--port", "2222"],
            "env": {
                "TERMINAL_FM_LOG_LEVEL": "debug"
            }
        }
    ]
}
```

### Profile Performance

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=. ./pkg/services/radiobrowser
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=. ./pkg/services/radiobrowser
go tool pprof mem.prof

# Analyze with pprof web UI
go tool pprof -http=:8080 cpu.prof
```

## Common Tasks

### Add New Dependency

```bash
# Add dependency
go get github.com/some/package@latest

# Update go.mod and go.sum
go mod tidy

# Verify
go mod verify
```

### Update Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/charmbracelet/bubbletea

# Tidy up
go mod tidy
```

### Generate Code

```bash
# If you add code generation (stringer, etc.)
go generate ./...
```

### Database Migrations

```bash
# Create new migration
cat > pkg/storage/migrations/003_add_history.sql << EOF
CREATE TABLE listening_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    station_uuid TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
EOF

# Run migrations
go run ./cmd/server migrate
```

### Build for Multiple Platforms

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o terminal-fm-linux-amd64 ./cmd/server

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o terminal-fm-darwin-amd64 ./cmd/server

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o terminal-fm-darwin-arm64 ./cmd/server

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o terminal-fm-windows-amd64.exe ./cmd/server
```

### Makefile Tasks

Create `Makefile`:
```makefile
.PHONY: build test lint run clean

build:
	go build -o terminal-fm ./cmd/server

build-prod:
	go build -ldflags="-s -w" -o terminal-fm ./cmd/server

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	goimports -w .

run:
	go run ./cmd/server --dev --port 2222

clean:
	rm -f terminal-fm
	rm -f coverage.out coverage.html
	rm -rf tmp/

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/cosmtrek/air@latest
```

Usage:
```bash
make build
make test
make lint
make run
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 2222
sudo lsof -i :2222

# Kill process
sudo kill -9 <PID>

# Or use different port
go run ./cmd/server --dev --port 2223
```

### Module Issues

```bash
# Clean module cache
go clean -modcache

# Re-download modules
go mod download

# Verify modules
go mod verify
```

### Build Errors

```bash
# Clean build cache
go clean -cache

# Rebuild
go build ./cmd/server
```

### Database Locked

```bash
# Stop any running instances
pkill terminal-fm

# Remove lock
rm -f tmp/dev.db-wal tmp/dev.db-shm

# Restart
go run ./cmd/server --dev
```

### SSH Connection Issues

```bash
# Check if server is listening
netstat -an | grep 2222

# Test SSH connection
ssh -v localhost -p 2222

# Check SSH logs
tail -f logs/ssh.log
```

## Development Tips

### Hot Reload Setup

Use Air for automatic rebuild on file changes:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run
air
```

### Mock External Dependencies

When testing, mock Radio Browser API:

```go
type MockRadioBrowserClient struct {
    SearchFunc func(params SearchParams) ([]Station, error)
}

func (m *MockRadioBrowserClient) Search(params SearchParams) ([]Station, error) {
    return m.SearchFunc(params)
}
```

### Use Table-Driven Tests

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name string
        input int
        want string
    }{
        {"case 1", 1, "one"},
        {"case 2", 2, "two"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Convert(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Debug Bubbletea TUI

Bubbletea captures stdout, so use logging:

```go
import "log"

// Create log file
f, _ := os.Create("debug.log")
log.SetOutput(f)

// Log debug info
log.Printf("Current state: %+v", m)
```

---

For deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).

For contributing guidelines, see [../CONTRIBUTING.md](../CONTRIBUTING.md).

For architecture details, see [ARCHITECTURE.md](ARCHITECTURE.md).
