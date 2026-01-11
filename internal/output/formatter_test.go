package output

import (
	"testing"

	"github.com/serithemage/updoc/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		format      string
		expectedErr bool
	}{
		{"html", false},
		{"markdown", false},
		{"text", false},
		{"json", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			f, err := NewFormatter(tt.format)
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, f)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, f)
			}
		})
	}
}

func TestHTMLFormatter(t *testing.T) {
	resp := &api.ParseResponse{
		Content: api.Content{
			HTML: "<h1>Test Title</h1><p>Test content</p>",
		},
	}

	f, err := NewFormatter("html")
	require.NoError(t, err)

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Equal(t, "<h1>Test Title</h1><p>Test content</p>", result)
}

func TestMarkdownFormatter(t *testing.T) {
	resp := &api.ParseResponse{
		Content: api.Content{
			Markdown: "# Test Title\n\nTest content",
		},
	}

	f, err := NewFormatter("markdown")
	require.NoError(t, err)

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Equal(t, "# Test Title\n\nTest content", result)
}

func TestTextFormatter(t *testing.T) {
	resp := &api.ParseResponse{
		Content: api.Content{
			Text: "Test Title\n\nTest content",
		},
	}

	f, err := NewFormatter("text")
	require.NoError(t, err)

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Equal(t, "Test Title\n\nTest content", result)
}

func TestJSONFormatter(t *testing.T) {
	resp := &api.ParseResponse{
		API:   "document-parse",
		Model: "document-parse-250618",
		Content: api.Content{
			HTML:     "<h1>Test</h1>",
			Markdown: "# Test",
			Text:     "Test",
		},
		Elements: []api.Element{
			{
				ID:       1,
				Category: "heading1",
				Page:     1,
			},
		},
		Usage: api.Usage{Pages: 1},
	}

	f, err := NewFormatter("json")
	require.NoError(t, err)

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Contains(t, result, "document-parse")
	assert.Contains(t, result, "heading1")
	assert.Contains(t, result, `"pages": 1`)
}

func TestElementsOnlyFormatter(t *testing.T) {
	resp := &api.ParseResponse{
		Elements: []api.Element{
			{
				ID:       1,
				Category: "heading1",
				Page:     1,
				Content: api.Content{
					Markdown: "# Title",
				},
			},
			{
				ID:       2,
				Category: "paragraph",
				Page:     1,
				Content: api.Content{
					Markdown: "Content text",
				},
			},
		},
	}

	f := &ElementsOnlyFormatter{OutputFormat: "markdown"}

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Contains(t, result, "[1] heading1 (page 1)")
	assert.Contains(t, result, "# Title")
	assert.Contains(t, result, "[2] paragraph (page 1)")
	assert.Contains(t, result, "Content text")
}

func TestElementsOnlyFormatterJSON(t *testing.T) {
	resp := &api.ParseResponse{
		Elements: []api.Element{
			{
				ID:       1,
				Category: "heading1",
				Page:     1,
			},
		},
	}

	f := &ElementsOnlyFormatter{OutputFormat: "json"}

	result, err := f.Format(resp)
	require.NoError(t, err)

	assert.Contains(t, result, `"elements"`)
	assert.Contains(t, result, "heading1")
}

func TestFormatterWithEmptyContent(t *testing.T) {
	resp := &api.ParseResponse{
		Content: api.Content{},
	}

	formats := []string{"html", "markdown", "text"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			f, err := NewFormatter(format)
			require.NoError(t, err)

			result, err := f.Format(resp)
			require.NoError(t, err)
			assert.Empty(t, result)
		})
	}
}
