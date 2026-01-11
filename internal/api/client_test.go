package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")

	assert.NotNil(t, client)
	assert.Equal(t, "test-api-key", client.apiKey)
	assert.NotNil(t, client.httpClient)
}

func TestNewClientWithOptions(t *testing.T) {
	customHTTPClient := &http.Client{}

	client := NewClient("test-key",
		WithHTTPClient(customHTTPClient),
		WithBaseURL("https://custom.api.com"),
	)

	assert.Equal(t, "test-key", client.apiKey)
	assert.Equal(t, customHTTPClient, client.httpClient)
	assert.Equal(t, "https://custom.api.com", client.baseURL)
}

func TestClientParse(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer test-api-key")
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")

		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20)
		require.NoError(t, err)

		// Check form values
		assert.Equal(t, "document-parse", r.FormValue("model"))
		assert.Equal(t, "standard", r.FormValue("mode"))

		// Check file
		file, header, err := r.FormFile("document")
		require.NoError(t, err)
		defer func() { _ = file.Close() }()
		assert.Equal(t, "test.pdf", header.Filename)

		// Return mock response
		resp := ParseResponse{
			API:   "document-parse",
			Model: "document-parse-250618",
			Content: Content{
				HTML:     "<h1>Test</h1>",
				Markdown: "# Test",
				Text:     "Test",
			},
			Elements: []Element{
				{
					ID:       1,
					Category: "heading1",
					Page:     1,
					Content: Content{
						HTML:     "<h1>Test</h1>",
						Markdown: "# Test",
						Text:     "Test",
					},
				},
			},
			Usage: Usage{Pages: 1},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create temp file
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := filepath.Join(tmpDir, "test.pdf")
	err = os.WriteFile(testFile, []byte("fake pdf content"), 0644)
	require.NoError(t, err)

	// Create client and make request
	client := NewClient("test-api-key", WithBaseURL(server.URL))

	req := NewParseRequest(testFile)
	resp, err := client.Parse(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "document-parse", resp.API)
	assert.Equal(t, "<h1>Test</h1>", resp.Content.HTML)
	assert.Equal(t, "# Test", resp.Content.Markdown)
	assert.Len(t, resp.Elements, 1)
	assert.Equal(t, 1, resp.Usage.Pages)
}

func TestClientParseFileNotFound(t *testing.T) {
	client := NewClient("test-api-key")

	req := NewParseRequest("/nonexistent/file.pdf")
	_, err := client.Parse(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}

func TestClientParseAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: ErrorDetail{
				Message: "Invalid API key",
				Type:    "authentication_error",
				Code:    "invalid_api_key",
			},
		})
	}))
	defer server.Close()

	// Create temp file
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := filepath.Join(tmpDir, "test.pdf")
	err = os.WriteFile(testFile, []byte("fake pdf content"), 0644)
	require.NoError(t, err)

	client := NewClient("invalid-key", WithBaseURL(server.URL))

	req := NewParseRequest(testFile)
	_, err = client.Parse(context.Background(), req)

	assert.Error(t, err)
	var apiErr *APIError
	assert.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal(t, "Invalid API key", apiErr.Message)
}

func TestClientParseAsync(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.Path, "/async")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AsyncResponse{
			RequestID: "req_abc123",
		})
	}))
	defer server.Close()

	// Create temp file
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := filepath.Join(tmpDir, "test.pdf")
	err = os.WriteFile(testFile, []byte("fake pdf content"), 0644)
	require.NoError(t, err)

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	req := NewParseRequest(testFile)
	resp, err := client.ParseAsync(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, "req_abc123", resp.RequestID)
}

func TestClientGetStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/async/req_abc123")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(StatusResponse{
			RequestID:      "req_abc123",
			Status:         "processing",
			Progress:       50,
			PagesProcessed: 5,
			TotalPages:     10,
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	resp, err := client.GetStatus(context.Background(), "req_abc123")

	require.NoError(t, err)
	assert.Equal(t, "req_abc123", resp.RequestID)
	assert.Equal(t, "processing", resp.Status)
	assert.Equal(t, 50, resp.Progress)
}

func TestClientGetResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/async/req_abc123/result")

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ParseResponse{
			API:   "document-parse",
			Model: "document-parse-250618",
			Content: Content{
				Markdown: "# Result",
			},
			Usage: Usage{Pages: 5},
		})
	}))
	defer server.Close()

	client := NewClient("test-api-key", WithBaseURL(server.URL))

	resp, err := client.GetResult(context.Background(), "req_abc123")

	require.NoError(t, err)
	assert.Equal(t, "# Result", resp.Content.Markdown)
	assert.Equal(t, 5, resp.Usage.Pages)
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 401,
		Message:    "Invalid API key",
		Type:       "authentication_error",
		Code:       "invalid_api_key",
	}

	assert.Contains(t, err.Error(), "401")
	assert.Contains(t, err.Error(), "Invalid API key")
}

func TestBuildMultipartForm(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	testFile := filepath.Join(tmpDir, "test.pdf")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	req := &ParseRequest{
		FilePath:         testFile,
		Model:            "document-parse",
		Mode:             "enhanced",
		OCR:              "force",
		ChartRecognition: true,
		MergeTables:      true,
		Coordinates:      false,
	}

	body, contentType, err := buildMultipartForm(req)
	require.NoError(t, err)
	assert.Contains(t, contentType, "multipart/form-data")

	// Read and verify body contains expected fields
	bodyBytes, err := io.ReadAll(body)
	require.NoError(t, err)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "document-parse")
	assert.Contains(t, bodyStr, "enhanced")
	assert.Contains(t, bodyStr, "force")
	assert.Contains(t, bodyStr, "test.pdf")
}
