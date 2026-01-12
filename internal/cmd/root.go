package cmd

import (
	"fmt"
	"os"

	"github.com/serithemage/updoc/internal/config"
	"github.com/spf13/cobra"
)

var (
	// Version info (set by ldflags)
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var (
	cfgFile string
	cfg     *config.Config
	verbose bool
	quiet   bool
)

var rootCmd = &cobra.Command{
	Use:   "updoc",
	Short: "Upstage Document Parse API CLI",
	Long: `updoc은 업스테이지 Document Parse API를 CLI로 사용할 수 있는 도구입니다.

다양한 형식의 문서(PDF, Office, HWP, 이미지)를
구조화된 텍스트(HTML, Markdown, Text)로 변환합니다.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
	rootCmd.PersistentFlags().String("api-key", "", "Upstage API key")
	rootCmd.PersistentFlags().String("endpoint", "", "API endpoint URL (for private hosting or AWS Bedrock)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress progress messages")
}

func initConfig() {
	var err error

	// Load config from file
	if cfgFile != "" {
		cfg, err = config.LoadFrom(cfgFile)
	} else {
		cfg, err = config.LoadFrom(config.GetDefaultConfigPath())
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = config.New()
	}

	// Override with environment variables
	cfg.LoadFromEnv()
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	if cfg == nil {
		cfg = config.New()
		cfg.LoadFromEnv()
	}
	return cfg
}

// GetAPIKey returns the API key from flags, env, or config
func GetAPIKey(cmd *cobra.Command) string {
	// 1. Check command flag
	if key, _ := cmd.Flags().GetString("api-key"); key != "" {
		return key
	}

	// 2. Check env (already loaded in config)
	// 3. Check config file
	return GetConfig().APIKey
}

// GetEndpoint returns the API endpoint from flags, env, or config
func GetEndpoint(cmd *cobra.Command) string {
	// 1. Check command flag
	if endpoint, _ := cmd.Flags().GetString("endpoint"); endpoint != "" {
		return endpoint
	}

	// 2. Check env (already loaded in config)
	// 3. Check config file, or return default
	return GetConfig().GetEndpoint()
}

// IsVerbose returns true if verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// IsQuiet returns true if quiet mode is enabled
func IsQuiet() bool {
	return quiet
}

// Printf prints a message if not in quiet mode
func Printf(format string, a ...interface{}) {
	if !quiet {
		fmt.Printf(format, a...)
	}
}

// Verbosef prints a message only in verbose mode
func Verbosef(format string, a ...interface{}) {
	if verbose {
		fmt.Printf("[DEBUG] "+format, a...)
	}
}
