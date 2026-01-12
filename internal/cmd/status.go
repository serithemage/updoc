package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/serithemage/updoc/internal/api"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <request-id>",
	Short: "Check async request status",
	Long:  `Check the status of an asynchronous parse request.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStatus,
}

func init() {
	statusCmd.Flags().BoolP("json", "j", false, "output as JSON")
	statusCmd.Flags().BoolP("watch", "w", false, "watch status until completion")
	statusCmd.Flags().IntP("interval", "i", 5, "watch interval in seconds")

	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) error {
	requestID := args[0]

	apiKey := GetAPIKey(cmd)
	if apiKey == "" {
		return fmt.Errorf("API key not set")
	}

	watch, _ := cmd.Flags().GetBool("watch")
	if watch {
		return watchStatus(cmd, apiKey, requestID)
	}

	return checkStatus(cmd, apiKey, requestID)
}

func checkStatus(cmd *cobra.Command, apiKey, requestID string) error {
	client := api.NewClient(apiKey, api.WithBaseURL(GetEndpoint(cmd)))

	resp, err := client.GetStatus(context.Background(), requestID)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		data, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Println(string(data))
		return nil
	}

	printStatus(resp)
	return nil
}

func watchStatus(cmd *cobra.Command, apiKey, requestID string) error {
	client := api.NewClient(apiKey, api.WithBaseURL(GetEndpoint(cmd)))
	interval, _ := cmd.Flags().GetInt("interval")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	// Check immediately first
	resp, err := client.GetStatus(context.Background(), requestID)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	printStatus(resp)

	if resp.Status == "completed" || resp.Status == "failed" {
		return nil
	}

	fmt.Println("\nWatching for updates (Ctrl+C to stop)...")

	for range ticker.C {
		resp, err := client.GetStatus(context.Background(), requestID)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		// Clear previous output and print new status
		fmt.Print("\033[2K\r") // Clear line
		printStatusLine(resp)

		if resp.Status == "completed" {
			fmt.Println("\n\nCompleted! Get result with: updoc result", requestID)
			return nil
		}
		if resp.Status == "failed" {
			fmt.Println("\n\nRequest failed:", resp.Error)
			return fmt.Errorf("request failed: %s", resp.Error)
		}
	}

	return nil
}

func printStatus(resp *api.StatusResponse) {
	fmt.Printf("Request ID: %s\n", resp.RequestID)
	fmt.Printf("Status: %s\n", resp.Status)
	if resp.TotalPages > 0 {
		fmt.Printf("Progress: %d%%\n", resp.Progress)
		fmt.Printf("Pages processed: %d/%d\n", resp.PagesProcessed, resp.TotalPages)
	}
	if resp.Error != "" {
		fmt.Printf("Error: %s\n", resp.Error)
	}
}

func printStatusLine(resp *api.StatusResponse) {
	if resp.TotalPages > 0 {
		fmt.Printf("Status: %s | Progress: %d%% (%d/%d pages)",
			resp.Status, resp.Progress, resp.PagesProcessed, resp.TotalPages)
	} else {
		fmt.Printf("Status: %s", resp.Status)
	}
}
