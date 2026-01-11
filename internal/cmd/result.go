package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/serithemage/updoc/internal/api"
	"github.com/spf13/cobra"
)

var resultCmd = &cobra.Command{
	Use:   "result <request-id>",
	Short: "Get async request result",
	Long:  `Get the result of a completed asynchronous parse request.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runResult,
}

func init() {
	resultCmd.Flags().StringP("output", "o", "", "output file path")
	resultCmd.Flags().StringP("format", "f", "", "output format: html, markdown, text")
	resultCmd.Flags().BoolP("wait", "w", false, "wait until completion")
	resultCmd.Flags().IntP("timeout", "t", 300, "wait timeout in seconds")
	resultCmd.Flags().BoolP("json", "j", false, "output as JSON")
	resultCmd.Flags().BoolP("elements-only", "e", false, "output only elements")

	rootCmd.AddCommand(resultCmd)
}

func runResult(cmd *cobra.Command, args []string) error {
	requestID := args[0]

	apiKey := GetAPIKey(cmd)
	if apiKey == "" {
		return fmt.Errorf("API key not set")
	}

	wait, _ := cmd.Flags().GetBool("wait")
	if wait {
		return waitAndGetResult(cmd, apiKey, requestID)
	}

	return getResult(cmd, apiKey, requestID)
}

func getResult(cmd *cobra.Command, apiKey, requestID string) error {
	client := api.NewClient(apiKey)

	// First check status
	status, err := client.GetStatus(context.Background(), requestID)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	if status.Status == "failed" {
		return fmt.Errorf("request failed: %s", status.Error)
	}

	if status.Status != "completed" {
		return fmt.Errorf("request not completed (status: %s). Use --wait to wait for completion", status.Status)
	}

	// Get result
	resp, err := client.GetResult(context.Background(), requestID)
	if err != nil {
		return fmt.Errorf("failed to get result: %w", err)
	}

	return outputResult(cmd, resp)
}

func waitAndGetResult(cmd *cobra.Command, apiKey, requestID string) error {
	client := api.NewClient(apiKey)
	timeout, _ := cmd.Flags().GetInt("timeout")

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	Printf("Waiting for request %s to complete...\n", requestID)

	for {
		status, err := client.GetStatus(context.Background(), requestID)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		if status.Status == "completed" {
			Printf("Request completed!\n")
			resp, err := client.GetResult(context.Background(), requestID)
			if err != nil {
				return fmt.Errorf("failed to get result: %w", err)
			}
			return outputResult(cmd, resp)
		}

		if status.Status == "failed" {
			return fmt.Errorf("request failed: %s", status.Error)
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for completion (status: %s, progress: %d%%)", status.Status, status.Progress)
		}

		Verbosef("Status: %s, Progress: %d%%\n", status.Status, status.Progress)

		<-ticker.C
	}
}
