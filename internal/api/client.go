package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Client is the Upstage Document Parse API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// ClientOption is a function that configures the client
type ClientOption func(*Client)

// NewClient creates a new API client
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// APIError represents an API error
type APIError struct {
	StatusCode int
	Message    string
	Type       string
	Code       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// Parse sends a synchronous parse request
func (c *Client) Parse(ctx context.Context, req *ParseRequest) (*ParseResponse, error) {
	body, contentType, err := buildMultipartForm(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/document-digitization", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var parseResp ParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&parseResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &parseResp, nil
}

// ParseAsync sends an asynchronous parse request
func (c *Client) ParseAsync(ctx context.Context, req *ParseRequest) (*AsyncResponse, error) {
	body, contentType, err := buildMultipartForm(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/document-digitization/async", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return nil, c.parseError(resp)
	}

	var asyncResp AsyncResponse
	if err := json.NewDecoder(resp.Body).Decode(&asyncResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &asyncResp, nil
}

// GetStatus gets the status of an async request
func (c *Client) GetStatus(ctx context.Context, requestID string) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/document-digitization/async/%s", c.baseURL, requestID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var statusResp StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &statusResp, nil
}

// GetResult gets the result of a completed async request
func (c *Client) GetResult(ctx context.Context, requestID string) (*ParseResponse, error) {
	url := fmt.Sprintf("%s/document-digitization/async/%s/result", c.baseURL, requestID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	var parseResp ParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&parseResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &parseResp, nil
}

// parseError parses an error response from the API
func (c *Client) parseError(resp *http.Response) error {
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    errResp.Error.Message,
		Type:       errResp.Error.Type,
		Code:       errResp.Error.Code,
	}
}

// buildMultipartForm builds a multipart form for the parse request
func buildMultipartForm(req *ParseRequest) (io.Reader, string, error) {
	// Open file
	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file
	filename := filepath.Base(req.FilePath)
	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Add form fields
	_ = writer.WriteField("model", req.Model)
	_ = writer.WriteField("mode", req.Mode)
	_ = writer.WriteField("ocr", req.OCR)
	_ = writer.WriteField("chart_recognition", strconv.FormatBool(req.ChartRecognition))
	_ = writer.WriteField("merge_multipage_tables", strconv.FormatBool(req.MergeTables))
	_ = writer.WriteField("coordinates", strconv.FormatBool(req.Coordinates))

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	return &buf, writer.FormDataContentType(), nil
}
