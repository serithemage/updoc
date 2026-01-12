# updoc CLI Manual

[한국어](CLI_MANUAL.ko.md) | [日本語](CLI_MANUAL.ja.md)

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Commands](#commands)
5. [Usage Examples](#usage-examples)
6. [Troubleshooting](#troubleshooting)
7. [API Reference](#api-reference)

---

## Introduction

`updoc` is a CLI tool that enables command-line access to Upstage's Document Parse API. Written in Go and distributed as a single binary, it converts various document formats into structured text (HTML, Markdown, Text).

### Supported Features

| Feature | Description |
|---------|-------------|
| Document Conversion | Convert PDF, Office, HWP to HTML/Markdown/Text |
| OCR | Extract text from scanned documents and images |
| Structure Analysis | Separate elements like headings, paragraphs, tables, figures |
| Layout Recognition | Handle multi-column layouts and complex structures |
| Coordinate Extraction | Provide position information for each element |

---

## Installation

### Requirements

- Operating System: macOS, Linux, Windows
- For building: Go 1.21 or later

### Binary Download

Download the binary for your OS from the [Releases](https://github.com/serithemage/updoc/releases) page.

#### macOS

```bash
# Apple Silicon (M1/M2/M3)
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-arm64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# Intel
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-darwin-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/
```

#### Linux

```bash
# amd64
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-linux-amd64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/

# arm64
curl -L https://github.com/serithemage/updoc/releases/latest/download/updoc-linux-arm64 -o updoc
chmod +x updoc
sudo mv updoc /usr/local/bin/
```

#### Windows

```powershell
# PowerShell
Invoke-WebRequest -Uri https://github.com/serithemage/updoc/releases/latest/download/updoc-windows-amd64.exe -OutFile updoc.exe

# Add to PATH or move to desired location
Move-Item updoc.exe C:\Users\$env:USERNAME\bin\
```

### Homebrew (macOS/Linux)

```bash
brew install serithemage/tap/updoc
```

### Install with Go

```bash
go install github.com/serithemage/updoc@latest
```

### Build from Source

```bash
git clone https://github.com/serithemage/updoc.git
cd updoc

# Build
go build -o updoc ./cmd/updoc

# Install (optional)
sudo mv updoc /usr/local/bin/
```

### Verify Installation

```bash
updoc version
updoc --help
```

---

## Configuration

### API Key Setup

An Upstage API key is required to use the Document Parse API.

#### 1. Get API Key

1. Log in to [Upstage Console](https://console.upstage.ai)
2. Create a new project or select an existing one
3. Generate a new key in the API Keys menu
4. Copy the generated key

#### 2. Key Configuration Methods

**Method A: Environment Variable (Recommended)**

```bash
# Linux/macOS
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (PowerShell)
$env:UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# Windows (CMD)
set UPSTAGE_API_KEY=up_xxxxxxxxxxxxxxxxxxxx
```

Add to shell config file for permanent settings:

```bash
# Add to ~/.bashrc or ~/.zshrc
echo 'export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"' >> ~/.zshrc
source ~/.zshrc
```

**Method B: Config Command**

```bash
updoc config set api-key up_xxxxxxxxxxxxxxxxxxxx
```

**Method C: Command Option**

```bash
updoc parse document.pdf --api-key up_xxxxxxxxxxxxxxxxxxxx
```

### Private Endpoint Configuration

For AWS Bedrock, private hosting, or custom endpoints, configure as follows:

**Method A: Environment Variable**

```bash
export UPSTAGE_API_ENDPOINT="https://your-private-endpoint.com/v1"
```

**Method B: Config Command**

```bash
updoc config set endpoint https://your-private-endpoint.com/v1
```

**Method C: Command Option**

```bash
updoc parse document.pdf --endpoint https://your-private-endpoint.com/v1
```

Priority: Command option > Environment variable > Config file > Default

### Config File

Config file locations:
- Linux/macOS: `~/.config/updoc/config.yaml`
- Windows: `%APPDATA%\updoc\config.yaml`

```yaml
api_key: "up_xxxxxxxxxxxxxxxxxxxx"
endpoint: ""  # Leave empty for default
default_format: markdown
default_mode: standard
default_ocr: auto
output_dir: "./output"
```

### Configuration Management

```bash
# View current settings
updoc config list

# Change settings
updoc config set default-format html
updoc config set default-mode enhanced

# Query settings
updoc config get default-format

# Reset settings
updoc config reset

# Show config file path
updoc config path
```

---

## Commands

### updoc parse

Parse documents and convert to structured text.

```
updoc parse <file> [options]
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `<file>` | Path to document file (required) |

#### Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--format <type>` | `-f` | Output format: html, markdown, text | markdown |
| `--output <path>` | `-o` | Output file path | stdout |
| `--mode <mode>` | `-m` | Parsing mode: standard, enhanced, auto | standard |
| `--model <name>` | | Model name | document-parse |
| `--ocr <type>` | | OCR setting: auto, force | auto |
| `--chart-recognition` | | Convert charts to tables | true |
| `--no-chart-recognition` | | Disable chart conversion | |
| `--merge-tables` | | Merge multi-page tables | false |
| `--coordinates` | | Include coordinate info | true |
| `--no-coordinates` | | Exclude coordinate info | |
| `--elements-only` | `-e` | Output only elements | false |
| `--json` | `-j` | Output as JSON | false |
| `--async` | `-a` | Use async processing | false |
| `--output-dir` | `-d` | Output directory for batch | . |
| `--recursive` | `-r` | Recursive directory traversal | false |
| `--quiet` | `-q` | Suppress progress messages | false |
| `--verbose` | `-v` | Verbose output | false |
| `--api-key <key>` | | Specify API key | env var |
| `--endpoint <url>` | | API endpoint URL | default endpoint |

#### Examples

```bash
# Basic usage
updoc parse report.pdf

# Convert to HTML and save
updoc parse report.pdf -f html -o report.html

# Process complex documents with enhanced mode
updoc parse complex-form.pdf --mode enhanced

# Force OCR on scanned documents
updoc parse scanned.pdf --ocr force

# JSON output
updoc parse document.pdf --json -o result.json

# Batch processing
updoc parse ./documents/*.pdf --output-dir ./results/
```

---

### updoc status

Check the status of async requests.

```
updoc status <request-id> [options]
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `<request-id>` | Async request ID (required) |

#### Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--json` | `-j` | Output as JSON | false |
| `--watch` | `-w` | Real-time status monitoring | false |
| `--interval` | `-i` | Monitoring interval (seconds) | 5 |

#### Examples

```bash
# Check status
updoc status abc123def456

# JSON output
updoc status abc123def456 --json

# Real-time monitoring
updoc status abc123def456 --watch
```

#### Output Example

```
Request ID: abc123def456
Status: processing
Progress: 45%
Pages processed: 45/100
Elapsed time: 1m 23s
```

---

### updoc result

Get the result of async requests.

```
updoc result <request-id> [options]
```

#### Arguments

| Argument | Description |
|----------|-------------|
| `<request-id>` | Async request ID (required) |

#### Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--output <path>` | `-o` | Output file path | stdout |
| `--format <type>` | `-f` | Output format | markdown |
| `--wait` | `-w` | Wait for completion | false |
| `--timeout <sec>` | `-t` | Wait timeout (seconds) | 300 |
| `--json` | `-j` | Output as JSON | false |

#### Examples

```bash
# Get result
updoc result abc123def456 -o output.md

# Wait for completion and get result
updoc result abc123def456 --wait -o output.md

# With timeout
updoc result abc123def456 --wait --timeout 600 -o output.md
```

---

### updoc models

Display available models.

```
updoc models [options]
```

#### Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--json` | `-j` | Output as JSON | false |

#### Output Example

```
Available Models:

  document-parse          Default model (recommended, alias)
  document-parse-250618   Specific version (2025-06-18)
  document-parse-nightly  Latest test version

Tip: Using 'document-parse' alias automatically applies the latest stable version.
```

---

### updoc config

Manage configuration.

```
updoc config <command> [key] [value]
```

#### Subcommands

| Command | Description |
|---------|-------------|
| `list` | Show all settings |
| `get <key>` | Query specific setting |
| `set <key> <value>` | Change setting |
| `reset` | Reset settings |
| `path` | Show config file path |

#### Config Keys

| Key | Description | Values |
|-----|-------------|--------|
| `api-key` | API key | string |
| `endpoint` | API endpoint URL | URL |
| `default-format` | Default output format | html, markdown, text |
| `default-mode` | Default parsing mode | standard, enhanced, auto |
| `default-ocr` | Default OCR setting | auto, force |
| `output-dir` | Default output directory | path |

#### Examples

```bash
# View all settings
updoc config list

# Query specific setting
updoc config get api-key

# Change settings
updoc config set default-format html
updoc config set default-mode enhanced

# Reset settings
updoc config reset
```

---

### updoc version

Display version information.

```
updoc version [options]
```

#### Options

| Option | Short | Description |
|--------|-------|-------------|
| `--short` | `-s` | Output version number only |
| `--json` | `-j` | Output as JSON |

#### Output Example

```
updoc version 1.0.0
  Commit: abc1234
  Built: 2025-01-11T10:00:00Z
  Go version: go1.21.5
  OS/Arch: darwin/arm64
```

---

## Usage Examples

### Basic Workflow

```bash
# 1. Set API key
export UPSTAGE_API_KEY="up_xxxxxxxxxxxxxxxxxxxx"

# 2. Convert PDF to Markdown
updoc parse report.pdf -o report.md

# 3. View result
cat report.md
```

### Processing Scanned Documents

```bash
# OCR process scanned PDF
updoc parse scanned-document.pdf --ocr force --mode enhanced -o output.md
```

### Complex Layout Documents

```bash
# Process document with many tables and charts in enhanced mode
updoc parse financial-report.pdf \
  --mode enhanced \
  --chart-recognition \
  --merge-tables \
  -o report.html \
  -f html
```

### Element-wise Analysis

```bash
# Parse document into elements and output JSON
updoc parse document.pdf --elements-only --json -o elements.json

# Extract only tables with jq
cat elements.json | jq '.elements[] | select(.category == "table")'

# Extract only headings with jq
cat elements.json | jq '.elements[] | select(.category | startswith("heading"))'
```

### Large Document Processing

```bash
# Start async request
updoc parse large-document.pdf --async
# Output: Request ID: req_abc123

# Real-time status monitoring
updoc status req_abc123 --watch

# Get result after completion
updoc result req_abc123 -o result.md

# Or wait for completion
updoc result req_abc123 --wait --timeout 600 -o result.md
```

### Batch Processing

```bash
# Process all PDFs in current directory
updoc parse *.pdf --output-dir ./results/

# Recursively process documents in directory
updoc parse ./documents/ --output-dir ./results/ --recursive

# Fine-grained control with shell script
for file in *.pdf; do
  echo "Processing: $file"
  updoc parse "$file" -o "${file%.pdf}.md" --quiet
done
```

### Pipeline Usage

```bash
# Pipe conversion results to other tools
updoc parse document.pdf | grep -i "important"

# Output to multiple formats simultaneously
updoc parse document.pdf --json | tee result.json | jq -r '.content.markdown' > result.md

# Extract and process specific elements
updoc parse document.pdf --json | jq -r '.elements[] | select(.category == "table") | .content.markdown'
```

### Automation Script Example

```bash
#!/bin/bash
# batch_convert.sh - Batch document conversion script

INPUT_DIR="${1:-.}"
OUTPUT_DIR="${2:-./output}"
FORMAT="${3:-markdown}"

mkdir -p "$OUTPUT_DIR"

find "$INPUT_DIR" -type f \( -name "*.pdf" -o -name "*.docx" -o -name "*.hwp" \) | while read -r file; do
  filename=$(basename "$file")
  output_file="$OUTPUT_DIR/${filename%.*}.${FORMAT}"

  echo "Converting: $filename"
  updoc parse "$file" -f "$FORMAT" -o "$output_file" --quiet

  if [ $? -eq 0 ]; then
    echo "  -> $output_file"
  else
    echo "  -> Failed"
  fi
done

echo "Done!"
```

---

## Troubleshooting

### Common Errors

#### API Key Error

```
Error: Invalid API key
```

Solutions:
1. Verify API key is correct
2. Check if environment variable is set: `echo $UPSTAGE_API_KEY`
3. Remove whitespace around key
4. Check settings: `updoc config get api-key`

#### File Format Error

```
Error: Unsupported file format: .xyz
```

Supported formats:
- Documents: PDF, DOCX, PPTX, XLSX, HWP
- Images: JPEG, PNG, BMP, TIFF, HEIC

#### Page Limit Exceeded

```
Error: Document exceeds maximum page limit (100 pages for sync API)
```

Sync API supports max 100 pages, async API supports max 1,000 pages.
Use `--async` option for large documents:

```bash
updoc parse large-document.pdf --async
```

#### Timeout

```
Error: Request timeout after 120s
```

Solutions:
- Use async mode: `--async`
- Check network status
- Check file size

#### File Access Error

```
Error: Cannot read file: permission denied
```

Solution:
```bash
# Check file permissions
ls -la document.pdf

# Modify permissions (if needed)
chmod 644 document.pdf
```

### Debugging

```bash
# Verbose output
updoc parse document.pdf --verbose

# Check request/response
updoc parse document.pdf --verbose 2>&1 | tee debug.log

# Check configuration
updoc config list
```

### Log Levels

Using the `--verbose` flag outputs the following information:
- API request URL and headers
- Request parameters
- Response status code
- Processing time

---

## API Reference

### Endpoints

| Purpose | URL |
|---------|-----|
| Sync parsing | `POST https://api.upstage.ai/v1/document-digitization` |
| Async parsing | `POST https://api.upstage.ai/v1/document-digitization/async` |
| Status check | `GET https://api.upstage.ai/v1/document-digitization/async/{id}` |

### Authentication

```
Authorization: Bearer <UPSTAGE_API_KEY>
```

### Request Format

`Content-Type: multipart/form-data`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `model` | string | Yes | Model name (document-parse) |
| `document` | file | Yes | Document file |
| `mode` | string | | standard, enhanced, auto |
| `ocr` | string | | auto, force |
| `output_formats` | string | | Output formats |
| `chart_recognition` | boolean | | Chart conversion |
| `merge_multipage_tables` | boolean | | Table merging |
| `coordinates` | boolean | | Include coordinates |

### Response Structure

```json
{
  "api": "document-parse",
  "model": "document-parse-250618",
  "content": {
    "html": "<h1>Title</h1>...",
    "markdown": "# Title\n...",
    "text": "Title\n..."
  },
  "elements": [
    {
      "id": 1,
      "category": "heading1",
      "page": 1,
      "content": {
        "html": "<h1>Title</h1>",
        "markdown": "# Title",
        "text": "Title"
      },
      "coordinates": [
        {"x": 0.1, "y": 0.05},
        {"x": 0.9, "y": 0.05},
        {"x": 0.9, "y": 0.08},
        {"x": 0.1, "y": 0.08}
      ]
    }
  ],
  "usage": {
    "pages": 10
  }
}
```

### Element Categories

| Category | Description |
|----------|-------------|
| `heading1` ~ `heading6` | Heading levels |
| `paragraph` | Paragraph |
| `table` | Table |
| `figure` | Figure |
| `chart` | Chart |
| `equation` | Equation |
| `list_item` | List item |
| `header` | Header |
| `footer` | Footer |
| `caption` | Caption |

---

## Appendix

### Environment Variables

| Variable | Description |
|----------|-------------|
| `UPSTAGE_API_KEY` | API authentication key |
| `UPSTAGE_API_ENDPOINT` | API endpoint URL (for private hosting) |
| `UPDOC_CONFIG_PATH` | Config file path (optional) |
| `UPDOC_LOG_LEVEL` | Log level: debug, info, warn, error |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Argument error |
| 3 | API error |
| 4 | File I/O error |
| 5 | Authentication error |

### Related Links

- [Upstage Console](https://console.upstage.ai)
- [Document Parse Documentation](https://console.upstage.ai/docs/capabilities/document-parse)
- [API Reference](https://console.upstage.ai/api-reference)
- [Upstage Blog](https://upstage.ai/blog)
- [GitHub Repository](https://github.com/serithemage/updoc)

### Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | - | Initial release |
