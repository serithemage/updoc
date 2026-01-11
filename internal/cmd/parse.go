package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/serithemage/updoc/internal/api"
	"github.com/serithemage/updoc/internal/output"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse <file>",
	Short: "Parse a document",
	Long: `Parse a document and convert it to structured text.

Supported file formats:
  Documents: PDF, DOCX, PPTX, XLSX, HWP
  Images: JPEG, PNG, BMP, TIFF, HEIC`,
	Args: cobra.ExactArgs(1),
	RunE: runParse,
}

func init() {
	parseCmd.Flags().StringP("format", "f", "", "output format: html, markdown, text (default from config or markdown)")
	parseCmd.Flags().StringP("output", "o", "", "output file path (default: stdout)")
	parseCmd.Flags().StringP("mode", "m", "", "parsing mode: standard, enhanced, auto (default from config or standard)")
	parseCmd.Flags().String("ocr", "", "OCR setting: auto, force (default from config or auto)")
	parseCmd.Flags().String("model", api.DefaultModel, "model to use")
	parseCmd.Flags().Bool("chart-recognition", true, "convert charts to tables")
	parseCmd.Flags().Bool("merge-tables", false, "merge multi-page tables")
	parseCmd.Flags().Bool("coordinates", true, "include coordinate information")
	parseCmd.Flags().BoolP("elements-only", "e", false, "output only elements")
	parseCmd.Flags().BoolP("json", "j", false, "output as JSON")
	parseCmd.Flags().BoolP("async", "a", false, "use async processing")

	rootCmd.AddCommand(parseCmd)
}

func runParse(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Validate file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Validate file format
	if !api.IsSupportedFile(filePath) {
		return fmt.Errorf("unsupported file format: %s (supported: pdf, docx, pptx, xlsx, hwp, jpg, jpeg, png, bmp, tiff, heic)", filepath.Ext(filePath))
	}

	// Get API key
	apiKey := GetAPIKey(cmd)
	if apiKey == "" {
		return fmt.Errorf("API key not set. Set it with 'updoc config set api-key <your-key>' or UPSTAGE_API_KEY environment variable")
	}

	// Build parse request
	req := api.NewParseRequest(filePath)

	// Apply flags
	if model, _ := cmd.Flags().GetString("model"); model != "" {
		req.Model = model
	}

	mode := getStringFlagOrConfig(cmd, "mode", GetConfig().DefaultMode)
	req.Mode = mode

	ocr := getStringFlagOrConfig(cmd, "ocr", GetConfig().DefaultOCR)
	req.OCR = ocr

	req.ChartRecognition, _ = cmd.Flags().GetBool("chart-recognition")
	req.MergeTables, _ = cmd.Flags().GetBool("merge-tables")
	req.Coordinates, _ = cmd.Flags().GetBool("coordinates")

	// Check async mode
	async, _ := cmd.Flags().GetBool("async")
	if async {
		return runParseAsync(cmd, apiKey, req)
	}

	return runParseSync(cmd, apiKey, req)
}

func runParseSync(cmd *cobra.Command, apiKey string, req *api.ParseRequest) error {
	client := api.NewClient(apiKey)

	Verbosef("Parsing file: %s\n", req.FilePath)
	Verbosef("Model: %s, Mode: %s, OCR: %s\n", req.Model, req.Mode, req.OCR)

	Printf("Parsing %s...\n", filepath.Base(req.FilePath))

	resp, err := client.Parse(context.Background(), req)
	if err != nil {
		return fmt.Errorf("parse failed: %w", err)
	}

	Verbosef("Parsed %d pages\n", resp.Usage.Pages)

	return outputResult(cmd, resp)
}

func runParseAsync(_ *cobra.Command, apiKey string, req *api.ParseRequest) error {
	client := api.NewClient(apiKey)

	Verbosef("Submitting async parse request for: %s\n", req.FilePath)

	resp, err := client.ParseAsync(context.Background(), req)
	if err != nil {
		return fmt.Errorf("async parse failed: %w", err)
	}

	Printf("Request submitted successfully\n")
	Printf("Request ID: %s\n", resp.RequestID)
	Printf("\n")
	Printf("Check status: updoc status %s\n", resp.RequestID)
	Printf("Get result:   updoc result %s\n", resp.RequestID)

	return nil
}

func outputResult(cmd *cobra.Command, resp *api.ParseResponse) error {
	elementsOnly, _ := cmd.Flags().GetBool("elements-only")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	outputPath, _ := cmd.Flags().GetString("output")

	// Determine format
	format := getStringFlagOrConfig(cmd, "format", GetConfig().DefaultFormat)
	if jsonOutput {
		format = "json"
	}

	// Create formatter
	var formatter output.Formatter
	var err error

	if elementsOnly {
		formatter = &output.ElementsOnlyFormatter{OutputFormat: format}
	} else {
		formatter, err = output.NewFormatter(format)
		if err != nil {
			return err
		}
	}

	// Format output
	result, err := formatter.Format(resp)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	// Write output
	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		Printf("Output written to: %s\n", outputPath)
	} else {
		fmt.Println(result)
	}

	return nil
}

func getStringFlagOrConfig(cmd *cobra.Command, flag, defaultValue string) string {
	if value, _ := cmd.Flags().GetString(flag); value != "" {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	return ""
}
