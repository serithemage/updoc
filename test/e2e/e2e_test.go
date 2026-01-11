//go:build e2e

// Package e2e contains end-to-end tests for the updoc CLI.
// These tests require the UPSTAGE_API_KEY environment variable for API tests.
// Run with: go test -tags=e2e ./test/e2e/...
//
// Go's test cache automatically caches results when source files haven't changed.
// To force re-run: go test -tags=e2e -count=1 ./test/e2e/...
package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	binaryPath  string
	testdataDir string
)

func TestMain(m *testing.M) {
	// Build the binary
	projectRoot := getProjectRoot()
	binaryPath = filepath.Join(projectRoot, "updoc-e2e-test")
	testdataDir = filepath.Join(projectRoot, "test", "testdata")

	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/updoc")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		println("Failed to build binary:", string(output))
		os.Exit(1)
	}

	code := m.Run()

	// Cleanup
	_ = os.Remove(binaryPath)

	os.Exit(code)
}

func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			panic("could not find project root (go.mod)")
		}
		dir = parent
	}
}

func runUpdoc(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func requireAPIKey(t *testing.T) string {
	t.Helper()
	apiKey := os.Getenv("UPSTAGE_API_KEY")
	if apiKey == "" {
		t.Skip("UPSTAGE_API_KEY not set, skipping API test")
	}
	return apiKey
}

// ============================================================
// Offline Tests (no API key required, always cached)
// ============================================================

func TestVersion(t *testing.T) {
	stdout, _, err := runUpdoc(t, "version")
	require.NoError(t, err)
	assert.Contains(t, stdout, "updoc version")
}

func TestVersionShort(t *testing.T) {
	stdout, _, err := runUpdoc(t, "version", "--short")
	require.NoError(t, err)
	assert.NotEmpty(t, strings.TrimSpace(stdout))
}

func TestVersionJSON(t *testing.T) {
	stdout, _, err := runUpdoc(t, "version", "--json")
	require.NoError(t, err)
	assert.Contains(t, stdout, `"version"`)
	assert.Contains(t, stdout, `"go_version"`)
}

func TestHelp(t *testing.T) {
	stdout, _, err := runUpdoc(t, "--help")
	require.NoError(t, err)
	assert.Contains(t, stdout, "updoc")
	assert.Contains(t, stdout, "parse")
	assert.Contains(t, stdout, "config")
}

func TestParseHelp(t *testing.T) {
	stdout, _, err := runUpdoc(t, "parse", "--help")
	require.NoError(t, err)
	assert.Contains(t, stdout, "--format")
	assert.Contains(t, stdout, "--output")
	assert.Contains(t, stdout, "--async")
	assert.Contains(t, stdout, "--output-dir")
	assert.Contains(t, stdout, "--recursive")
}

func TestConfigList(t *testing.T) {
	stdout, _, err := runUpdoc(t, "config", "list")
	require.NoError(t, err)
	assert.Contains(t, stdout, "api-key")
	assert.Contains(t, stdout, "default-format")
}

func TestConfigPath(t *testing.T) {
	stdout, _, err := runUpdoc(t, "config", "path")
	require.NoError(t, err)
	assert.Contains(t, stdout, "config.yaml")
}

func TestConfigSetGet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "updoc-e2e-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Set
	_, _, err = runUpdoc(t, "--config", configPath, "config", "set", "default-format", "html")
	require.NoError(t, err)

	// Get
	stdout, _, err := runUpdoc(t, "--config", configPath, "config", "get", "default-format")
	require.NoError(t, err)
	assert.Contains(t, stdout, "html")
}

func TestModels(t *testing.T) {
	stdout, _, err := runUpdoc(t, "models")
	require.NoError(t, err)
	assert.Contains(t, stdout, "document-parse")
}

func TestParseWithoutAPIKey(t *testing.T) {
	originalKey := os.Getenv("UPSTAGE_API_KEY")
	_ = os.Unsetenv("UPSTAGE_API_KEY")
	defer func() {
		if originalKey != "" {
			_ = os.Setenv("UPSTAGE_API_KEY", originalKey)
		}
	}()

	pdfFile := filepath.Join(testdataDir, "dummy.pdf")
	stdout, stderr, err := runUpdoc(t, "parse", pdfFile)
	assert.Error(t, err)
	combined := stdout + stderr
	assert.Contains(t, combined, "API key not set")
}

func TestParseFileNotFound(t *testing.T) {
	_ = requireAPIKey(t)
	stdout, stderr, err := runUpdoc(t, "parse", "/nonexistent/file.pdf")
	assert.Error(t, err)
	combined := stdout + stderr
	assert.Contains(t, combined, "file not found")
}

func TestParseUnsupportedFormat(t *testing.T) {
	_ = requireAPIKey(t)

	tmpFile, err := os.CreateTemp("", "test-*.xyz")
	require.NoError(t, err)
	defer func() { _ = os.Remove(tmpFile.Name()) }()
	_ = tmpFile.Close()

	stdout, stderr, err := runUpdoc(t, "parse", tmpFile.Name())
	assert.Error(t, err)
	combined := stdout + stderr
	assert.Contains(t, combined, "unsupported file format")
}

// ============================================================
// API Tests (requires UPSTAGE_API_KEY)
// ============================================================

func TestParsePDFMarkdown(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "-f", "markdown")
	require.NoError(t, err)
	assert.NotEmpty(t, stdout)
}

func TestParsePDFHTML(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "test.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "-f", "html")
	require.NoError(t, err)
	assert.NotEmpty(t, stdout)
}

func TestParsePDFText(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "-f", "text")
	require.NoError(t, err)
	assert.NotEmpty(t, stdout)
}

func TestParsePDFJSON(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "--json")
	require.NoError(t, err)
	assert.Contains(t, stdout, `"api"`)
	assert.Contains(t, stdout, `"content"`)
}

func TestParsePDFElementsOnly(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "--elements-only")
	require.NoError(t, err)
	assert.NotEmpty(t, stdout)
}

func TestParsePDFToFile(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "test.pdf") // Use test.pdf which has more content

	tmpDir, err := os.MkdirTemp("", "updoc-e2e-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outputFile := filepath.Join(tmpDir, "output.md")
	stdout, stderr, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "-o", outputFile)
	require.NoError(t, err, "stdout: %s, stderr: %s", stdout, stderr)

	// Verify file was created
	_, err = os.Stat(outputFile)
	require.NoError(t, err, "output file should exist")
}

func TestParseAsync(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "parse", pdfFile, "--async")
	require.NoError(t, err)
	assert.Contains(t, stdout, "Request ID")
}

func TestVerboseMode(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "--verbose", "parse", pdfFile, "-f", "text")
	require.NoError(t, err)
	assert.Contains(t, stdout, "[DEBUG]")
}

func TestQuietMode(t *testing.T) {
	apiKey := requireAPIKey(t)
	pdfFile := filepath.Join(testdataDir, "dummy.pdf")

	stdout, _, err := runUpdoc(t, "--api-key", apiKey, "--quiet", "parse", pdfFile, "-f", "text")
	require.NoError(t, err)
	assert.NotContains(t, stdout, "Parsing")
}
