# updoc

[![Go Version](https://img.shields.io/github/go-mod/go-version/serithemage/updoc)](https://go.dev/)
[![CI](https://github.com/serithemage/updoc/actions/workflows/ci.yaml/badge.svg)](https://github.com/serithemage/updoc/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/serithemage/updoc)](https://github.com/serithemage/updoc/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

[한국어](README.ko.md) | [日本語](README.ja.md)

A CLI tool for the Upstage Document Parse API.

## Overview

`updoc` is a command-line interface for the Upstage Document Parse API that converts PDF, images, and office documents into structured text (HTML, Markdown, Text). Written in Go, it's distributed as a single binary and supports cross-platform environments.

### Key Features

- Supports various document formats: PDF, DOCX, PPTX, XLSX, HWP, HWPX
- OCR processing for images and scanned documents: JPEG, PNG, BMP, TIFF, HEIC
- Multiple output formats: HTML, Markdown, Text, JSON
- Structured results by element (headings, paragraphs, tables, figures, etc.)
- Batch processing with recursive directory traversal
- Sync/async processing support (up to 1,000 pages)
- Single binary, no external dependencies

## Installation

Requires Go 1.21 or later.

```bash
go install github.com/serithemage/updoc/cmd/updoc@latest
```

Or build from source:

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc
make build
```

## Quick Start

### 1. Set Up API Key

Get your API key from [Upstage Console](https://console.upstage.ai) and set it as an environment variable.

```bash
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"
```

Or use the config command:

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

### Private Endpoint Configuration (Optional)

For AWS Bedrock or private hosting environments, you can configure a custom endpoint.

```bash
# Set via environment variable
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"

# Or use the config command
updoc config set endpoint https://your-private-endpoint.com/v1

# Or specify as a command option
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

### 2. Parse Documents

```bash
# Convert PDF to Markdown (default)
updoc parse document.pdf

# Save output to file
updoc parse document.pdf -o result.md

# Convert to HTML format
updoc parse document.pdf -f html -o result.html
```

## Usage

### Basic Parsing

```bash
# Output to stdout
updoc parse report.pdf

# Save to file
updoc parse report.pdf -o report.md

# Specify output format: markdown (default), html, text, json
updoc parse report.pdf -f html -o report.html
```

### Parsing Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| `standard` | Fast processing (default) | Simple layout documents |
| `enhanced` | Precise analysis | Complex tables, charts, scanned documents |
| `auto` | Automatic selection | Determined by document characteristics |

```bash
# Document with complex tables and charts
updoc parse financial-report.pdf --mode enhanced

# Scanned document (force OCR)
updoc parse scanned.pdf --ocr force --mode enhanced
```

### Batch Processing

```bash
# Process multiple files at once
updoc parse *.pdf --output-dir ./results/

# Process all documents in directory recursively
updoc parse ./documents/ --output-dir ./results/ --recursive

# Process files matching specific pattern
updoc parse ./docs/**/*.pdf --output-dir ./converted/
```

### Advanced Options

```bash
# Convert charts to tables
updoc parse report.pdf --chart-recognition

# Merge multi-page tables
updoc parse spreadsheet.pdf --merge-tables

# Include coordinate information
updoc parse document.pdf --coordinates

# Output only elements (exclude full content)
updoc parse document.pdf --elements-only

# Output full API response as JSON
updoc parse document.pdf --json -o result.json
```

### Async Processing (Large Documents)

Use async API for documents exceeding 100 pages.

```bash
# Start async request
updoc parse large-document.pdf --async
# Output: Request ID: req_abc123def456

# Check status
updoc status req_abc123def456

# Monitor status in real-time
updoc status req_abc123def456 --watch

# Get result
updoc result req_abc123def456 -o output.md

# Wait for completion and get result
updoc result req_abc123def456 --wait -o output.md
```

### Configuration Management

```bash
# View current settings
updoc config list

# Change default output format
updoc config set default-format html

# List available models
updoc models
```

## Command Summary

| Command | Description |
|---------|-------------|
| `updoc parse <file>` | Parse document |
| `updoc status <id>` | Check async request status |
| `updoc result <id>` | Get async request result |
| `updoc config` | Manage configuration |
| `updoc models` | List available models |
| `updoc version` | Show version info |

For detailed options and usage, see the [CLI Manual](docs/CLI_MANUAL.md).

## Supported File Formats

| Category | Formats |
|----------|---------|
| Documents | PDF, DOCX, PPTX, XLSX, HWP, HWPX |
| Images | JPEG, PNG, BMP, TIFF, HEIC |

## API Limits

| Item | Sync API | Async API |
|------|----------|-----------|
| Max pages | 100 | 1,000 |
| Recommended for | Small documents | Large documents, batch processing |

## Contributing

Thank you for contributing to the project! Please read our [Contributing Guide](CONTRIBUTING.md) for detailed information.

### Development Setup

```bash
# Clone repository
git clone https://github.com/serithemage/updoc.git
cd updoc

# Set up development environment (Git hooks, linter installation)
make dev-setup

# Build
make build

# Run tests
make test

# E2E tests (requires API key)
export UPSTAGE_API_KEY="your-api-key"
make test-e2e

# Lint
make lint
```

### How to Contribute

1. Check existing issues or create a new one
2. Fork the repository
3. Create a feature branch (`git checkout -b feature/amazing-feature`)
4. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Create a Pull Request

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Add/modify tests
- `refactor:` Refactoring
- `chore:` Other changes

### Project Structure

```
updoc/
├── cmd/updoc/           # Entry point
├── internal/
│   ├── api/             # Upstage API client
│   ├── cmd/             # CLI command implementation
│   ├── config/          # Configuration management
│   └── output/          # Output formatters
├── test/e2e/            # E2E tests
├── docs/                # Documentation
└── Makefile
```

## License

MIT License

## References

- [CLI Manual](docs/CLI_MANUAL.md)
- [Upstage Document Parse Documentation](https://console.upstage.ai/docs/capabilities/document-parse)
- [Upstage API Reference](https://console.upstage.ai/api-reference)
- [Upstage Console](https://console.upstage.ai)
