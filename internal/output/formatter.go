package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/serithemage/updoc/internal/api"
)

// Formatter formats ParseResponse to a specific output format
type Formatter interface {
	Format(resp *api.ParseResponse) (string, error)
}

// NewFormatter creates a new formatter for the given format
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "html":
		return &HTMLFormatter{}, nil
	case "markdown":
		return &MarkdownFormatter{}, nil
	case "text":
		return &TextFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// HTMLFormatter outputs HTML content
type HTMLFormatter struct{}

func (f *HTMLFormatter) Format(resp *api.ParseResponse) (string, error) {
	return resp.Content.HTML, nil
}

// MarkdownFormatter outputs Markdown content
type MarkdownFormatter struct{}

func (f *MarkdownFormatter) Format(resp *api.ParseResponse) (string, error) {
	return resp.Content.Markdown, nil
}

// TextFormatter outputs plain text content
type TextFormatter struct{}

func (f *TextFormatter) Format(resp *api.ParseResponse) (string, error) {
	return resp.Content.Text, nil
}

// JSONFormatter outputs the full response as JSON
type JSONFormatter struct{}

func (f *JSONFormatter) Format(resp *api.ParseResponse) (string, error) {
	data, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(data), nil
}

// ElementsOnlyFormatter outputs only the elements
type ElementsOnlyFormatter struct {
	OutputFormat string // markdown, text, json
}

func (f *ElementsOnlyFormatter) Format(resp *api.ParseResponse) (string, error) {
	if f.OutputFormat == "json" {
		return f.formatJSON(resp)
	}
	return f.formatText(resp)
}

func (f *ElementsOnlyFormatter) formatJSON(resp *api.ParseResponse) (string, error) {
	output := struct {
		Elements []api.Element `json:"elements"`
	}{
		Elements: resp.Elements,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal elements: %w", err)
	}
	return string(data), nil
}

func (f *ElementsOnlyFormatter) formatText(resp *api.ParseResponse) (string, error) {
	var sb strings.Builder

	for _, elem := range resp.Elements {
		sb.WriteString(fmt.Sprintf("[%d] %s (page %d)\n", elem.ID, elem.Category, elem.Page))

		content := ""
		switch f.OutputFormat {
		case "markdown":
			content = elem.Content.Markdown
		case "html":
			content = elem.Content.HTML
		default:
			content = elem.Content.Text
			if content == "" {
				content = elem.Content.Markdown
			}
		}

		if content != "" {
			sb.WriteString(content)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n"), nil
}
