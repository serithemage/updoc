package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/serithemage/updoc/internal/api"
	"github.com/serithemage/updoc/internal/output"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:   "parse <file|directory|pattern>",
	Short: "Parse a document or multiple documents",
	Long: `Parse a document and convert it to structured text.

Supported file formats:
  Documents: PDF, DOCX, PPTX, XLSX, HWP
  Images: JPEG, PNG, BMP, TIFF, HEIC

Batch processing:
  Parse multiple files using glob patterns, directories, or --output-dir option.

Examples:
  # Single file
  updoc parse document.pdf

  # Glob pattern
  updoc parse "*.pdf" --output-dir ./results/

  # Directory (non-recursive)
  updoc parse ./documents/ --output-dir ./results/

  # Directory (recursive)
  updoc parse ./documents/ --output-dir ./results/ --recursive`,
	Args: cobra.ExactArgs(1),
	RunE: runParse,
}

func init() {
	parseCmd.Flags().StringP("format", "f", "", "output format: html, markdown, text (default from config or markdown)")
	parseCmd.Flags().StringP("output", "o", "", "output file path (default: stdout)")
	parseCmd.Flags().StringP("output-dir", "d", "", "output directory for batch processing")
	parseCmd.Flags().BoolP("recursive", "r", false, "process directories recursively")
	parseCmd.Flags().StringP("mode", "m", "", "parsing mode: standard, enhanced, auto (default from config or standard)")
	parseCmd.Flags().String("ocr", "", "OCR setting: auto, force (default from config or auto)")
	parseCmd.Flags().String("model", api.DefaultModel, "model to use")
	parseCmd.Flags().Bool("chart-recognition", true, "convert charts to tables")
	parseCmd.Flags().Bool("no-chart-recognition", false, "disable chart recognition")
	parseCmd.Flags().Bool("merge-tables", false, "merge multi-page tables")
	parseCmd.Flags().Bool("coordinates", true, "include coordinate information")
	parseCmd.Flags().Bool("no-coordinates", false, "exclude coordinate information")
	parseCmd.Flags().BoolP("elements-only", "e", false, "output only elements")
	parseCmd.Flags().BoolP("json", "j", false, "output as JSON")
	parseCmd.Flags().BoolP("async", "a", false, "use async processing")

	rootCmd.AddCommand(parseCmd)
}

func runParse(cmd *cobra.Command, args []string) error {
	inputPath := args[0]
	outputDir, _ := cmd.Flags().GetString("output-dir")
	recursive, _ := cmd.Flags().GetBool("recursive")

	// Get API key
	apiKey := GetAPIKey(cmd)
	if apiKey == "" {
		return fmt.Errorf("API key not set. Set it with 'updoc config set api-key <your-key>' or UPSTAGE_API_KEY environment variable")
	}

	// Collect files to process
	files, err := collectFiles(inputPath, recursive)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no supported files found matching: %s", inputPath)
	}

	// Single file mode
	if len(files) == 1 && outputDir == "" {
		return processSingleFile(cmd, apiKey, files[0])
	}

	// Batch mode requires output-dir
	if outputDir == "" {
		return fmt.Errorf("--output-dir is required for batch processing (multiple files)")
	}

	// Create output directory if needed
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return processBatch(cmd, apiKey, files, outputDir)
}

func collectFiles(inputPath string, recursive bool) ([]string, error) {
	var files []string

	// Check if it's a glob pattern
	if strings.ContainsAny(inputPath, "*?[") {
		matches, err := filepath.Glob(inputPath)
		if err != nil {
			return nil, fmt.Errorf("invalid glob pattern: %w", err)
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			if !info.IsDir() && api.IsSupportedFile(match) {
				files = append(files, match)
			}
		}
		return files, nil
	}

	// Check if path exists
	info, err := os.Stat(inputPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", inputPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// Single file
	if !info.IsDir() {
		if !api.IsSupportedFile(inputPath) {
			return nil, fmt.Errorf("unsupported file format: %s", filepath.Ext(inputPath))
		}
		return []string{inputPath}, nil
	}

	// Directory
	if recursive {
		walkErr := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && api.IsSupportedFile(path) {
				files = append(files, path)
			}
			return nil
		})
		if walkErr != nil {
			return nil, fmt.Errorf("failed to scan directory: %w", walkErr)
		}
	} else {
		entries, readErr := os.ReadDir(inputPath)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read directory: %w", readErr)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				path := filepath.Join(inputPath, entry.Name())
				if api.IsSupportedFile(path) {
					files = append(files, path)
				}
			}
		}
	}

	return files, nil
}

func processSingleFile(cmd *cobra.Command, apiKey string, filePath string) error {
	req := buildParseRequest(cmd, filePath)

	async, _ := cmd.Flags().GetBool("async")
	if async {
		return runParseAsync(apiKey, req)
	}

	return runParseSync(cmd, apiKey, req)
}

func processBatch(cmd *cobra.Command, apiKey string, files []string, outputDir string) error {
	format := getStringFlagOrConfig(cmd, "format", GetConfig().DefaultFormat)
	if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
		format = "json"
	}

	ext := getExtensionForFormat(format)
	client := api.NewClient(apiKey)

	var successCount, failCount int
	var failedFiles []string

	Printf("Processing %d files...\n\n", len(files))

	for _, filePath := range files {
		baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
		outputPath := filepath.Join(outputDir, baseName+ext)

		Printf("Processing: %s... ", filepath.Base(filePath))

		req := buildParseRequest(cmd, filePath)
		resp, err := client.Parse(context.Background(), req)
		if err != nil {
			Printf("failed (%v)\n", err)
			failCount++
			failedFiles = append(failedFiles, filePath)
			continue
		}

		result, err := formatResult(cmd, resp)
		if err != nil {
			Printf("failed (%v)\n", err)
			failCount++
			failedFiles = append(failedFiles, filePath)
			continue
		}

		if err := os.WriteFile(outputPath, []byte(result), 0644); err != nil {
			Printf("failed (%v)\n", err)
			failCount++
			failedFiles = append(failedFiles, filePath)
			continue
		}

		Printf("done -> %s\n", outputPath)
		successCount++
	}

	// Print summary
	Printf("\nSummary:\n")
	Printf("  Total:   %d\n", len(files))
	Printf("  Success: %d\n", successCount)
	Printf("  Failed:  %d\n", failCount)

	if len(failedFiles) > 0 {
		Printf("\nFailed files:\n")
		for _, f := range failedFiles {
			Printf("  - %s\n", f)
		}
		return fmt.Errorf("%d files failed to process", failCount)
	}

	return nil
}

func buildParseRequest(cmd *cobra.Command, filePath string) *api.ParseRequest {
	req := api.NewParseRequest(filePath)

	if model, _ := cmd.Flags().GetString("model"); model != "" {
		req.Model = model
	}

	req.Mode = getStringFlagOrConfig(cmd, "mode", GetConfig().DefaultMode)
	req.OCR = getStringFlagOrConfig(cmd, "ocr", GetConfig().DefaultOCR)

	req.ChartRecognition, _ = cmd.Flags().GetBool("chart-recognition")
	if noChart, _ := cmd.Flags().GetBool("no-chart-recognition"); noChart {
		req.ChartRecognition = false
	}

	req.MergeTables, _ = cmd.Flags().GetBool("merge-tables")

	req.Coordinates, _ = cmd.Flags().GetBool("coordinates")
	if noCoords, _ := cmd.Flags().GetBool("no-coordinates"); noCoords {
		req.Coordinates = false
	}

	return req
}

func getExtensionForFormat(format string) string {
	switch format {
	case "html":
		return ".html"
	case "text":
		return ".txt"
	case "json":
		return ".json"
	default:
		return ".md"
	}
}

func formatResult(cmd *cobra.Command, resp *api.ParseResponse) (string, error) {
	elementsOnly, _ := cmd.Flags().GetBool("elements-only")
	jsonOutput, _ := cmd.Flags().GetBool("json")

	format := getStringFlagOrConfig(cmd, "format", GetConfig().DefaultFormat)
	if jsonOutput {
		format = "json"
	}

	var formatter output.Formatter
	var err error

	if elementsOnly {
		formatter = &output.ElementsOnlyFormatter{OutputFormat: format}
	} else {
		formatter, err = output.NewFormatter(format)
		if err != nil {
			return "", err
		}
	}

	return formatter.Format(resp)
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

func runParseAsync(apiKey string, req *api.ParseRequest) error {
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
	outputPath, _ := cmd.Flags().GetString("output")

	result, err := formatResult(cmd, resp)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

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
