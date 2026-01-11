package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/serithemage/updoc/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `Manage updoc configuration settings.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg := GetConfig()

		if err := cfg.Set(key, value); err != nil {
			return err
		}

		configPath := config.GetDefaultConfigPath()
		if cfgFile != "" {
			configPath = cfgFile
		}

		if err := cfg.SaveTo(configPath); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg := GetConfig()
		value, err := cfg.Get(key)
		if err != nil {
			return err
		}

		// Mask API key for security
		if key == "api-key" && value != "" {
			fmt.Println(config.MaskAPIKey(value))
		} else {
			fmt.Println(value)
		}

		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetConfig()

		fmt.Println("Configuration:")
		fmt.Println()

		// API Key (masked)
		apiKeyDisplay := "(not set)"
		if cfg.APIKey != "" {
			apiKeyDisplay = config.MaskAPIKey(cfg.APIKey) + " (set)"
		}
		fmt.Printf("  api-key:        %s\n", apiKeyDisplay)
		fmt.Printf("  default-format: %s\n", cfg.DefaultFormat)
		fmt.Printf("  default-mode:   %s\n", cfg.DefaultMode)
		fmt.Printf("  default-ocr:    %s\n", cfg.DefaultOCR)

		outputDir := cfg.OutputDir
		if outputDir == "" {
			outputDir = "(not set)"
		}
		fmt.Printf("  output-dir:     %s\n", outputDir)
		fmt.Println()

		configPath := config.GetDefaultConfigPath()
		if cfgFile != "" {
			configPath = cfgFile
		}
		fmt.Printf("Config file: %s\n", configPath)
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Print("Are you sure you want to reset all configuration? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		configPath := config.GetDefaultConfigPath()
		if cfgFile != "" {
			configPath = cfgFile
		}

		// Remove config file if exists
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove config file: %w", err)
		}

		// Reset in-memory config
		cfg = config.New()

		Printf("Configuration reset to defaults.\n")
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Run: func(cmd *cobra.Command, args []string) {
		configPath := config.GetDefaultConfigPath()
		if cfgFile != "" {
			configPath = cfgFile
		}
		fmt.Println(configPath)
	},
}

func init() {
	configResetCmd.Flags().Bool("force", false, "skip confirmation prompt")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)

	rootCmd.AddCommand(configCmd)
}
