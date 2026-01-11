package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

// Default values
const (
	DefaultFormat = "markdown"
	DefaultMode   = "standard"
	DefaultOCR    = "auto"
)

// Environment variable names
const (
	EnvAPIKey     = "UPSTAGE_API_KEY"
	EnvConfigPath = "UPDOC_CONFIG_PATH"
	EnvLogLevel   = "UPDOC_LOG_LEVEL"
)

// Valid values
var (
	ValidFormats = []string{"html", "markdown", "text"}
	ValidModes   = []string{"standard", "enhanced", "auto"}
	ValidOCRs    = []string{"auto", "force"}
)

// Errors
var (
	ErrUnknownKey    = errors.New("unknown configuration key")
	ErrInvalidFormat = errors.New("invalid format: must be html, markdown, or text")
	ErrInvalidMode   = errors.New("invalid mode: must be standard, enhanced, or auto")
	ErrInvalidOCR    = errors.New("invalid ocr: must be auto or force")
)

// Config holds the application configuration
type Config struct {
	APIKey        string `yaml:"api_key"`
	DefaultFormat string `yaml:"default_format"`
	DefaultMode   string `yaml:"default_mode"`
	DefaultOCR    string `yaml:"default_ocr"`
	OutputDir     string `yaml:"output_dir"`
}

// New creates a new Config with default values
func New() *Config {
	return &Config{
		APIKey:        "",
		DefaultFormat: DefaultFormat,
		DefaultMode:   DefaultMode,
		DefaultOCR:    DefaultOCR,
		OutputDir:     "",
	}
}

// Set sets a configuration value by key
func (c *Config) Set(key, value string) error {
	switch key {
	case "api-key":
		c.APIKey = value
	case "default-format":
		if !IsValidFormat(value) {
			return ErrInvalidFormat
		}
		c.DefaultFormat = value
	case "default-mode":
		if !IsValidMode(value) {
			return ErrInvalidMode
		}
		c.DefaultMode = value
	case "default-ocr":
		if !IsValidOCR(value) {
			return ErrInvalidOCR
		}
		c.DefaultOCR = value
	case "output-dir":
		c.OutputDir = value
	default:
		return fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}
	return nil
}

// Get gets a configuration value by key
func (c *Config) Get(key string) (string, error) {
	switch key {
	case "api-key":
		return c.APIKey, nil
	case "default-format":
		return c.DefaultFormat, nil
	case "default-mode":
		return c.DefaultMode, nil
	case "default-ocr":
		return c.DefaultOCR, nil
	case "output-dir":
		return c.OutputDir, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrUnknownKey, key)
	}
}

// Reset resets the configuration to default values
func (c *Config) Reset() {
	c.APIKey = ""
	c.DefaultFormat = DefaultFormat
	c.DefaultMode = DefaultMode
	c.DefaultOCR = DefaultOCR
	c.OutputDir = ""
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv() {
	if apiKey := os.Getenv(EnvAPIKey); apiKey != "" {
		c.APIKey = apiKey
	}
}

// SaveTo saves the configuration to a file
func (c *Config) SaveTo(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write with restricted permissions (0600)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// LoadFrom loads the configuration from a file
func LoadFrom(path string) (*Config, error) {
	cfg := New()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	// Check for environment override
	if path := os.Getenv(EnvConfigPath); path != "" {
		return path
	}

	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	default: // linux, darwin
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			homeDir, _ := os.UserHomeDir()
			configDir = filepath.Join(homeDir, ".config")
		}
	}

	return filepath.Join(configDir, "updoc", "config.yaml")
}

// IsValidFormat checks if the format is valid
func IsValidFormat(format string) bool {
	for _, v := range ValidFormats {
		if v == format {
			return true
		}
	}
	return false
}

// IsValidMode checks if the mode is valid
func IsValidMode(mode string) bool {
	for _, v := range ValidModes {
		if v == mode {
			return true
		}
	}
	return false
}

// IsValidOCR checks if the OCR setting is valid
func IsValidOCR(ocr string) bool {
	for _, v := range ValidOCRs {
		if v == ocr {
			return true
		}
	}
	return false
}

// MaskAPIKey masks the API key for display
func MaskAPIKey(key string) string {
	if len(key) == 0 {
		return ""
	}
	if len(key) <= 4 {
		return "****"[:len(key)]
	}
	// Show last N characters, mask the rest
	visible := len(key) - 4
	if visible > 12 {
		visible = 12
	}
	return "****" + key[len(key)-visible:]
}
