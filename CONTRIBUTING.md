# Contributing to updoc

[한국어](CONTRIBUTING.ko.md) | [日本語](CONTRIBUTING.ja.md)

Thank you for your interest in contributing to updoc! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Commit Message Convention](#commit-message-convention)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project follows a Code of Conduct that all contributors are expected to uphold. Please be respectful and constructive in your interactions.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile commands)
- golangci-lint (for linting)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/updoc.git
   cd updoc
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/serithemage/updoc.git
   ```

## Development Setup

### Quick Setup

```bash
# Install development dependencies and set up Git hooks
make dev-setup
```

### Manual Setup

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build the project
go build -o updoc ./cmd/updoc

# Run tests
go test ./...
```

### Project Structure

```
updoc/
├── cmd/updoc/           # Application entry point
├── internal/
│   ├── api/             # Upstage API client
│   ├── cmd/             # CLI command implementations
│   ├── config/          # Configuration management
│   └── output/          # Output formatters
├── test/e2e/            # End-to-end tests
├── docs/                # Documentation
├── .github/             # GitHub templates and workflows
└── Makefile             # Build automation
```

## How to Contribute

### Reporting Bugs

Before creating a bug report:
1. Check existing issues to avoid duplicates
2. Collect relevant information:
   - updoc version (`updoc version`)
   - Operating system and version
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages or logs

Use the [Bug Report template](.github/ISSUE_TEMPLATE/bug_report.md) when creating an issue.

### Suggesting Features

We welcome feature suggestions! Before submitting:
1. Check if the feature has already been requested
2. Consider if it aligns with the project's goals
3. Provide clear use cases

Use the [Feature Request template](.github/ISSUE_TEMPLATE/feature_request.md) when creating an issue.

### Contributing Code

1. **Find an issue** to work on, or create one for discussion
2. **Comment** on the issue to let others know you're working on it
3. **Create a branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/bug-description
   ```
4. **Make your changes** following our coding standards
5. **Write tests** for new functionality
6. **Run tests and linting**:
   ```bash
   make test
   make lint
   ```
7. **Commit** your changes following our commit convention
8. **Push** to your fork and create a Pull Request

## Pull Request Process

### Before Submitting

- [ ] Code compiles without errors
- [ ] All tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation is updated if needed
- [ ] Commit messages follow convention

### PR Guidelines

1. **Title**: Use a clear, descriptive title
2. **Description**: Explain what changes you made and why
3. **Link Issues**: Reference related issues (e.g., "Fixes #123")
4. **Keep it Focused**: One PR should address one concern
5. **Be Responsive**: Address review feedback promptly

### Review Process

1. Automated checks must pass (CI/CD)
2. At least one maintainer approval is required
3. All discussions must be resolved
4. Branch must be up-to-date with `main`

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Guidelines

- Keep functions small and focused
- Write descriptive variable and function names
- Add comments for exported functions and complex logic
- Handle errors explicitly, don't ignore them
- Avoid global state

### Example

```go
// ParseDocument parses a document file and returns structured content.
// It returns an error if the file format is not supported or if parsing fails.
func ParseDocument(filePath string, opts ...Option) (*Result, error) {
    if filePath == "" {
        return nil, errors.New("file path cannot be empty")
    }

    // ... implementation
}
```

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `style` | Code style (formatting, etc.) |
| `refactor` | Code refactoring |
| `test` | Adding or updating tests |
| `chore` | Maintenance tasks |
| `perf` | Performance improvements |
| `ci` | CI/CD changes |

### Examples

```bash
feat(parse): add support for HWPX file format

fix(config): resolve API key loading from environment

docs(readme): update installation instructions

test(api): add unit tests for async parsing
```

## Testing

### Running Tests

```bash
# Unit tests
make test

# E2E tests (requires UPSTAGE_API_KEY)
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# All tests with coverage
go test -cover ./...
```

### Writing Tests

- Place tests in `*_test.go` files
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for meaningful coverage, not just high numbers

### Example Test

```go
func TestParseRequest_Validate(t *testing.T) {
    tests := []struct {
        name    string
        req     *ParseRequest
        wantErr bool
    }{
        {
            name:    "valid request",
            req:     &ParseRequest{FilePath: "test.pdf"},
            wantErr: false,
        },
        {
            name:    "empty file path",
            req:     &ParseRequest{FilePath: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

### When to Update Docs

- Adding new features or commands
- Changing existing behavior
- Adding new configuration options
- Fixing unclear or incorrect documentation

### Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | Project overview and quick start |
| `docs/CLI_MANUAL.md` | Detailed CLI reference |
| `CONTRIBUTING.md` | Contribution guidelines |

### Multilingual Support

This project maintains documentation in English, Korean, and Japanese. When updating documentation:

1. Update the English version first
2. Use `/translate-docs` to sync translations, or
3. Manually update translations maintaining consistency

## Questions?

- Open a [Discussion](https://github.com/serithemage/updoc/discussions) for questions
- Check existing issues and discussions first
- Be clear and provide context

Thank you for contributing to updoc!
