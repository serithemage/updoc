package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequestDefaults(t *testing.T) {
	req := NewParseRequest("test.pdf")

	assert.Equal(t, "test.pdf", req.FilePath)
	assert.Equal(t, DefaultModel, req.Model)
	assert.Equal(t, "standard", req.Mode)
	assert.Equal(t, "auto", req.OCR)
	assert.True(t, req.ChartRecognition)
	assert.False(t, req.MergeTables)
	assert.True(t, req.Coordinates)
}

func TestParseResponseJSON(t *testing.T) {
	jsonData := `{
		"api": "document-parse",
		"model": "document-parse-250618",
		"content": {
			"html": "<h1>Title</h1>",
			"markdown": "# Title",
			"text": "Title"
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
	}`

	var resp ParseResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, "document-parse", resp.API)
	assert.Equal(t, "document-parse-250618", resp.Model)
	assert.Equal(t, "<h1>Title</h1>", resp.Content.HTML)
	assert.Equal(t, "# Title", resp.Content.Markdown)
	assert.Equal(t, "Title", resp.Content.Text)
	assert.Len(t, resp.Elements, 1)
	assert.Equal(t, 1, resp.Elements[0].ID)
	assert.Equal(t, "heading1", resp.Elements[0].Category)
	assert.Equal(t, 1, resp.Elements[0].Page)
	assert.Len(t, resp.Elements[0].Coordinates, 4)
	assert.Equal(t, 0.1, resp.Elements[0].Coordinates[0].X)
	assert.Equal(t, 10, resp.Usage.Pages)
}

func TestErrorResponseJSON(t *testing.T) {
	jsonData := `{
		"error": {
			"message": "Invalid API key",
			"type": "authentication_error",
			"code": "invalid_api_key"
		}
	}`

	var resp ErrorResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, "Invalid API key", resp.Error.Message)
	assert.Equal(t, "authentication_error", resp.Error.Type)
	assert.Equal(t, "invalid_api_key", resp.Error.Code)
}

func TestAsyncResponseJSON(t *testing.T) {
	jsonData := `{
		"request_id": "req_abc123def456"
	}`

	var resp AsyncResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, "req_abc123def456", resp.RequestID)
}

func TestStatusResponseJSON(t *testing.T) {
	jsonData := `{
		"request_id": "req_abc123",
		"status": "processing",
		"progress": 45,
		"pages_processed": 45,
		"total_pages": 100
	}`

	var resp StatusResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	require.NoError(t, err)

	assert.Equal(t, "req_abc123", resp.RequestID)
	assert.Equal(t, "processing", resp.Status)
	assert.Equal(t, 45, resp.Progress)
	assert.Equal(t, 45, resp.PagesProcessed)
	assert.Equal(t, 100, resp.TotalPages)
}

func TestIsSupportedFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		// Documents
		{"document.pdf", true},
		{"document.PDF", true},
		{"report.docx", true},
		{"slides.pptx", true},
		{"data.xlsx", true},
		{"korean.hwp", true},
		// Images
		{"photo.jpg", true},
		{"photo.jpeg", true},
		{"image.png", true},
		{"scan.bmp", true},
		{"scan.tiff", true},
		{"photo.heic", true},
		// Unsupported
		{"file.txt", false},
		{"code.go", false},
		{"archive.zip", false},
		{"noextension", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsSupportedFile(tt.filename))
		})
	}
}

func TestElementCategories(t *testing.T) {
	// Verify all expected categories exist
	expectedCategories := []string{
		CategoryHeading1,
		CategoryHeading2,
		CategoryHeading3,
		CategoryHeading4,
		CategoryHeading5,
		CategoryHeading6,
		CategoryParagraph,
		CategoryTable,
		CategoryFigure,
		CategoryChart,
		CategoryEquation,
		CategoryListItem,
		CategoryHeader,
		CategoryFooter,
		CategoryCaption,
	}

	for _, cat := range expectedCategories {
		assert.NotEmpty(t, cat)
	}
}
