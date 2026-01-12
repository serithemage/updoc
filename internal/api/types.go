package api

import (
	"path/filepath"
	"strings"
)

// API endpoints
const (
	DefaultBaseURL = "https://api.upstage.ai/v1"
	DefaultModel   = "document-parse"
)

// BaseURL is kept for backward compatibility
const BaseURL = DefaultBaseURL

// Supported file extensions
var SupportedExtensions = map[string]bool{
	// Documents
	".pdf":  true,
	".docx": true,
	".pptx": true,
	".xlsx": true,
	".hwp":  true,
	// Images
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".bmp":  true,
	".tiff": true,
	".heic": true,
}

// Element categories
const (
	CategoryHeading1  = "heading1"
	CategoryHeading2  = "heading2"
	CategoryHeading3  = "heading3"
	CategoryHeading4  = "heading4"
	CategoryHeading5  = "heading5"
	CategoryHeading6  = "heading6"
	CategoryParagraph = "paragraph"
	CategoryTable     = "table"
	CategoryFigure    = "figure"
	CategoryChart     = "chart"
	CategoryEquation  = "equation"
	CategoryListItem  = "list_item"
	CategoryHeader    = "header"
	CategoryFooter    = "footer"
	CategoryCaption   = "caption"
)

// ParseRequest represents a document parse request
type ParseRequest struct {
	FilePath         string
	Model            string
	Mode             string // standard, enhanced, auto
	OCR              string // auto, force
	OutputFormats    []string
	ChartRecognition bool
	MergeTables      bool
	Coordinates      bool
}

// NewParseRequest creates a new ParseRequest with default values
func NewParseRequest(filePath string) *ParseRequest {
	return &ParseRequest{
		FilePath:         filePath,
		Model:            DefaultModel,
		Mode:             "standard",
		OCR:              "auto",
		OutputFormats:    []string{"html", "markdown", "text"},
		ChartRecognition: true,
		MergeTables:      false,
		Coordinates:      true,
	}
}

// ParseResponse represents the response from the parse API
type ParseResponse struct {
	API      string    `json:"api"`
	Model    string    `json:"model"`
	Content  Content   `json:"content"`
	Elements []Element `json:"elements"`
	Usage    Usage     `json:"usage"`
}

// Content holds the parsed content in different formats
type Content struct {
	HTML     string `json:"html"`
	Markdown string `json:"markdown"`
	Text     string `json:"text"`
	PDF      string `json:"pdf,omitempty"`
}

// Element represents a parsed document element
type Element struct {
	ID          int          `json:"id"`
	Category    string       `json:"category"`
	Page        int          `json:"page"`
	Content     Content      `json:"content"`
	Coordinates []Coordinate `json:"coordinates,omitempty"`
	Base64      string       `json:"base64_encoding,omitempty"`
}

// Coordinate represents a point in the document
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Usage contains usage information
type Usage struct {
	Pages int `json:"pages"`
}

// AsyncResponse represents the response from async parse request
type AsyncResponse struct {
	RequestID string `json:"request_id"`
}

// StatusResponse represents the status of an async request
type StatusResponse struct {
	RequestID      string `json:"request_id"`
	Status         string `json:"status"` // pending, processing, completed, failed
	Progress       int    `json:"progress"`
	PagesProcessed int    `json:"pages_processed"`
	TotalPages     int    `json:"total_pages"`
	Error          string `json:"error,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details
type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// IsSupportedFile checks if the file extension is supported
func IsSupportedFile(filename string) bool {
	if filename == "" {
		return false
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return SupportedExtensions[ext]
}
