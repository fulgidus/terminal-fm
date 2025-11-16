# Contributing to Terminal.FM

Thank you for considering contributing to Terminal.FM! We welcome contributions from everyone.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Translation Guidelines](#translation-guidelines)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of experience level, background, or identity.

### Expected Behavior

- Be respectful and constructive in communication
- Accept constructive criticism gracefully
- Focus on what's best for the project and community
- Show empathy towards other contributors

### Unacceptable Behavior

- Harassment, discrimination, or offensive comments
- Personal attacks or trolling
- Publishing others' private information
- Any conduct that would be inappropriate in a professional setting

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report:
1. Check the [issue tracker](https://github.com/fulgidus/terminal-fm/issues) to avoid duplicates
2. Verify you're using the latest version
3. Test with a clean environment if possible

When submitting a bug report, include:
- **Description**: Clear summary of the issue
- **Steps to Reproduce**: Detailed steps to trigger the bug
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happens
- **Environment**: OS, Go version, terminal emulator
- **Logs**: Relevant error messages or logs

### Suggesting Features

Feature requests are welcome! Please:
1. Check existing [feature requests](https://github.com/fulgidus/terminal-fm/issues?q=is%3Aissue+label%3Aenhancement)
2. Explain the problem your feature solves
3. Describe your proposed solution
4. Consider potential drawbacks or alternatives

### Contributing Code

We welcome code contributions! Areas where help is especially appreciated:
- Bug fixes
- New features from the roadmap
- Performance improvements
- Test coverage
- Documentation improvements

### Contributing Translations

See [Translation Guidelines](#translation-guidelines) below for details on adding new languages.

## Development Setup

### Prerequisites

- **Go**: 1.21 or higher
- **Git**: Latest stable version
- **mpv** or **ffplay**: For testing audio playback
- **make**: For build automation (optional)

### Setup Steps

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/terminal-fm.git
   cd terminal-fm
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Run Development Server**
   ```bash
   go run ./cmd/server --dev --port 2222
   ```

4. **Connect Locally**
   ```bash
   ssh localhost -p 2222
   ```

5. **Run Tests**
   ```bash
   go test ./...
   ```

6. **Lint Code**
   ```bash
   golangci-lint run
   ```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed development documentation.

## Pull Request Process

### Before Submitting

1. **Create an Issue**: Discuss significant changes before implementation
2. **Branch from `main`**: Use descriptive branch names (`fix/player-crash`, `feat/lastfm-scrobbling`)
3. **Write Tests**: Add tests for new functionality
4. **Update Documentation**: Keep docs in sync with code changes
5. **Follow Coding Standards**: See below

### Submitting

1. **Push Your Branch**
   ```bash
   git push origin feat/your-feature-name
   ```

2. **Open Pull Request**
   - Use clear, descriptive title
   - Reference related issues (`Fixes #123`, `Closes #456`)
   - Provide detailed description of changes
   - Include screenshots for UI changes

3. **Address Review Comments**
   - Respond to all feedback
   - Make requested changes
   - Push updates to the same branch

### Review Process

- At least one maintainer approval required
- All CI checks must pass
- Code coverage should not decrease
- Documentation must be updated

### After Merge

- Your contribution will be included in the next release
- You'll be added to CONTRIBUTORS.md (if not already)
- Delete your feature branch

## Coding Standards

### Go Style Guide

Follow official [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) and [Effective Go](https://go.dev/doc/effective_go).

#### Key Points

**Formatting**
- Use `gofmt` (or `goimports`) before committing
- No tabs vs spaces debate: Go uses tabs

**Naming**
- Exported names: `CamelCase`
- Unexported names: `camelCase`
- Acronyms: `HTTPClient`, not `HttpClient`
- Interface names: `-er` suffix when appropriate (`Player`, `Renderer`)

**Comments**
- Package comment: `// Package x does y.`
- Exported names: Full sentence starting with name
  ```go
  // Station represents a radio station from Radio Browser API.
  type Station struct { ... }
  ```

**Error Handling**
- Return errors, don't panic
- Wrap errors with context: `fmt.Errorf("failed to connect: %w", err)`
- Check errors immediately

**Structure**
```go
// Good
func (p *Player) Play(url string) error {
    if url == "" {
        return fmt.Errorf("empty URL")
    }
    
    if err := p.connect(url); err != nil {
        return fmt.Errorf("play failed: %w", err)
    }
    
    return nil
}

// Bad
func (p *Player) Play(url string) error {
    err := p.connect(url)
    if err != nil {
        panic(err) // Don't panic!
    }
    return nil
}
```

### Project Structure

```
terminal-fm/
├── cmd/
│   └── server/          # Main application entry point
├── pkg/
│   ├── ssh/            # SSH server logic
│   ├── ui/             # Bubbletea TUI components
│   ├── services/       # Business logic
│   ├── storage/        # Database layer
│   └── i18n/           # Internationalization
├── internal/           # Private application code
├── docs/               # Documentation
├── locales/            # Translation files
└── scripts/            # Build and deployment scripts
```

### Testing Standards

**Unit Tests**
- Test file suffix: `_test.go`
- Test function: `func TestFeatureName(t *testing.T)`
- Use table-driven tests when appropriate
- Mock external dependencies

**Example**
```go
func TestPlayerVolume(t *testing.T) {
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
            p := NewPlayer()
            err := p.SetVolume(tt.volume)
            if (err != nil) != tt.wantErr {
                t.Errorf("SetVolume() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Integration Tests**
- Use build tags: `//go:build integration`
- Test actual API interactions (with rate limiting)
- Clean up resources after tests

**Minimum Coverage**
- New code: 80% coverage
- Bug fixes: Add regression test

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style (formatting, missing semicolons, etc.)
- `refactor`: Code restructuring without behavior change
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks (dependencies, build, etc.)
- `ci`: CI/CD configuration changes

### Examples

```
feat(player): add volume control support

Implement volume adjustment using mpv IPC protocol.
Users can now use +/- keys to adjust volume.

Closes #42
```

```
fix(ssh): prevent connection leak on disconnect

Properly close SSH sessions when user disconnects
abruptly. Fixes memory leak reported in #78.
```

```
docs(api): add Radio Browser API examples

Include code examples for common API queries.
```

### Rules

- Use imperative mood ("add feature", not "added feature")
- Don't capitalize first letter of subject
- No period at end of subject
- Limit subject line to 72 characters
- Separate subject from body with blank line
- Wrap body at 72 characters
- Reference issues in footer

## Translation Guidelines

### Adding a New Language

1. **Check Existing Support**
   Ensure the language isn't already in `pkg/i18n/locales/`

2. **Create Translation File**
   ```bash
   cp pkg/i18n/locales/active.en.toml pkg/i18n/locales/active.{LOCALE}.toml
   ```
   Example: `active.es.toml` for Spanish

3. **Translate Strings**
   Keep keys unchanged, translate only values:
   ```toml
   # English
   [welcome]
   other = "Welcome to Terminal.FM!"
   
   # Spanish
   [welcome]
   other = "¡Bienvenido a Terminal.FM!"
   ```

4. **Update Language Registry**
   Add your language to `pkg/i18n/i18n.go`:
   ```go
   var supportedLanguages = []string{
       "en-US",
       "it-IT",
       "es-ES", // Your addition
   }
   ```

5. **Test Your Translation**
   ```bash
   LANG=es_ES go run ./cmd/server --dev
   ssh localhost -p 2222
   ```

6. **Submit Pull Request**
   Title: `feat(i18n): add {language} translation`

### Translation Best Practices

- **Consistency**: Use consistent terminology throughout
- **Context**: Consider the UI context when translating
- **Length**: Translations may be longer; test UI rendering
- **Pluralization**: Use i18n plural forms correctly
  ```toml
  [station_count]
  one = "{{.Count}} station"
  other = "{{.Count}} stations"
  ```
- **Formatting**: Preserve placeholders like `{{.Name}}`
- **Cultural Sensitivity**: Avoid idioms that don't translate well

### Translation Checklist

- [ ] All strings translated
- [ ] Placeholders preserved
- [ ] Plural forms handled
- [ ] UI tested with translation
- [ ] No encoding issues
- [ ] Language added to settings menu

See [docs/I18N.md](docs/I18N.md) for detailed internationalization documentation.

## Getting Help

- **Documentation**: Check [docs/](docs/) folder
- **Discussions**: Use [GitHub Discussions](https://github.com/fulgidus/terminal-fm/discussions)
- **Issues**: Search existing issues or create new one
- **Discord**: (Coming soon)

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub contributors page

Thank you for contributing to Terminal.FM! 🎵
