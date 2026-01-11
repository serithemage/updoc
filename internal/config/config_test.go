package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfg := New()

	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, DefaultFormat, cfg.DefaultFormat)
	assert.Equal(t, DefaultMode, cfg.DefaultMode)
	assert.Equal(t, DefaultOCR, cfg.DefaultOCR)
	assert.Equal(t, "", cfg.OutputDir)
}

func TestConfigDefaults(t *testing.T) {
	assert.Equal(t, "markdown", DefaultFormat)
	assert.Equal(t, "standard", DefaultMode)
	assert.Equal(t, "auto", DefaultOCR)
}

func TestConfigSet(t *testing.T) {
	cfg := New()

	tests := []struct {
		key      string
		value    string
		expected string
		getFunc  func() string
		wantErr  bool
	}{
		{
			key:      "api-key",
			value:    "test-api-key",
			expected: "test-api-key",
			getFunc:  func() string { return cfg.APIKey },
			wantErr:  false,
		},
		{
			key:      "default-format",
			value:    "html",
			expected: "html",
			getFunc:  func() string { return cfg.DefaultFormat },
			wantErr:  false,
		},
		{
			key:      "default-format",
			value:    "invalid",
			expected: "html", // should remain unchanged
			getFunc:  func() string { return cfg.DefaultFormat },
			wantErr:  true,
		},
		{
			key:      "default-mode",
			value:    "enhanced",
			expected: "enhanced",
			getFunc:  func() string { return cfg.DefaultMode },
			wantErr:  false,
		},
		{
			key:      "default-mode",
			value:    "invalid",
			expected: "enhanced", // should remain unchanged
			getFunc:  func() string { return cfg.DefaultMode },
			wantErr:  true,
		},
		{
			key:      "default-ocr",
			value:    "force",
			expected: "force",
			getFunc:  func() string { return cfg.DefaultOCR },
			wantErr:  false,
		},
		{
			key:      "default-ocr",
			value:    "invalid",
			expected: "force", // should remain unchanged
			getFunc:  func() string { return cfg.DefaultOCR },
			wantErr:  true,
		},
		{
			key:      "output-dir",
			value:    "/tmp/output",
			expected: "/tmp/output",
			getFunc:  func() string { return cfg.OutputDir },
			wantErr:  false,
		},
		{
			key:      "unknown-key",
			value:    "value",
			expected: "",
			getFunc:  func() string { return "" },
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.key+"="+tt.value, func(t *testing.T) {
			err := cfg.Set(tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, tt.getFunc())
		})
	}
}

func TestConfigGet(t *testing.T) {
	cfg := New()
	cfg.APIKey = "test-key"
	cfg.DefaultFormat = "html"
	cfg.DefaultMode = "enhanced"
	cfg.DefaultOCR = "force"
	cfg.OutputDir = "/tmp"

	tests := []struct {
		key      string
		expected string
		wantErr  bool
	}{
		{"api-key", "test-key", false},
		{"default-format", "html", false},
		{"default-mode", "enhanced", false},
		{"default-ocr", "force", false},
		{"output-dir", "/tmp", false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			val, err := cfg.Get(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

func TestConfigValidateFormat(t *testing.T) {
	tests := []struct {
		format  string
		isValid bool
	}{
		{"html", true},
		{"markdown", true},
		{"text", true},
		{"json", false},
		{"pdf", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			assert.Equal(t, tt.isValid, IsValidFormat(tt.format))
		})
	}
}

func TestConfigValidateMode(t *testing.T) {
	tests := []struct {
		mode    string
		isValid bool
	}{
		{"standard", true},
		{"enhanced", true},
		{"auto", true},
		{"fast", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.mode, func(t *testing.T) {
			assert.Equal(t, tt.isValid, IsValidMode(tt.mode))
		})
	}
}

func TestConfigValidateOCR(t *testing.T) {
	tests := []struct {
		ocr     string
		isValid bool
	}{
		{"auto", true},
		{"force", true},
		{"off", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ocr, func(t *testing.T) {
			assert.Equal(t, tt.isValid, IsValidOCR(tt.ocr))
		})
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create and save config
	cfg := New()
	cfg.APIKey = "test-api-key-12345"
	cfg.DefaultFormat = "html"
	cfg.DefaultMode = "enhanced"
	cfg.DefaultOCR = "force"
	cfg.OutputDir = "/custom/output"

	err = cfg.SaveTo(configPath)
	require.NoError(t, err)

	// Verify file exists and has correct permissions
	info, err := os.Stat(configPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Load config
	loaded, err := LoadFrom(configPath)
	require.NoError(t, err)

	assert.Equal(t, cfg.APIKey, loaded.APIKey)
	assert.Equal(t, cfg.DefaultFormat, loaded.DefaultFormat)
	assert.Equal(t, cfg.DefaultMode, loaded.DefaultMode)
	assert.Equal(t, cfg.DefaultOCR, loaded.DefaultOCR)
	assert.Equal(t, cfg.OutputDir, loaded.OutputDir)
}

func TestConfigLoadNonExistent(t *testing.T) {
	cfg, err := LoadFrom("/nonexistent/path/config.yaml")
	assert.NoError(t, err) // Should return default config, not error
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, DefaultFormat, cfg.DefaultFormat)
}

func TestConfigAPIKeyFromEnv(t *testing.T) {
	// Set environment variable
	originalKey := os.Getenv(EnvAPIKey)
	defer func() { _ = os.Setenv(EnvAPIKey, originalKey) }()

	_ = os.Setenv(EnvAPIKey, "env-api-key")

	cfg := New()
	cfg.LoadFromEnv()

	assert.Equal(t, "env-api-key", cfg.APIKey)
}

func TestConfigAPIKeyPriority(t *testing.T) {
	// Env should take priority over config file
	tmpDir, err := os.MkdirTemp("", "updoc-test-*")
	require.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save config with file-based key
	cfg := New()
	cfg.APIKey = "file-api-key"
	err = cfg.SaveTo(configPath)
	require.NoError(t, err)

	// Set env variable
	originalKey := os.Getenv(EnvAPIKey)
	defer func() { _ = os.Setenv(EnvAPIKey, originalKey) }()
	_ = os.Setenv(EnvAPIKey, "env-api-key")

	// Load and apply env
	loaded, err := LoadFrom(configPath)
	require.NoError(t, err)
	loaded.LoadFromEnv()

	assert.Equal(t, "env-api-key", loaded.APIKey)
}

func TestGetConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "updoc")
	assert.Contains(t, path, "config.yaml")
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"", ""},
		{"abc", "***"},
		{"abcd", "****"},
		{"abcdefgh", "****efgh"},
		{"up_1234567890abcdef", "****567890abcdef"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			assert.Equal(t, tt.expected, MaskAPIKey(tt.key))
		})
	}
}

func TestConfigReset(t *testing.T) {
	cfg := New()
	cfg.APIKey = "test-key"
	cfg.DefaultFormat = "html"
	cfg.DefaultMode = "enhanced"
	cfg.DefaultOCR = "force"
	cfg.OutputDir = "/custom"

	cfg.Reset()

	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, DefaultFormat, cfg.DefaultFormat)
	assert.Equal(t, DefaultMode, cfg.DefaultMode)
	assert.Equal(t, DefaultOCR, cfg.DefaultOCR)
	assert.Equal(t, "", cfg.OutputDir)
}
