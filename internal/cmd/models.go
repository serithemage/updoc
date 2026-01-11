package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type modelInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Recommended bool   `json:"recommended"`
}

var availableModels = []modelInfo{
	{
		Name:        "document-parse",
		Description: "기본 모델 (권장, alias)",
		Recommended: true,
	},
	{
		Name:        "document-parse-250618",
		Description: "특정 버전 (2025-06-18)",
		Recommended: false,
	},
	{
		Name:        "document-parse-nightly",
		Description: "최신 테스트 버전",
		Recommended: false,
	},
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models",
	Long:  `List available document parsing models.`,
	Run:   runModels,
}

func init() {
	modelsCmd.Flags().BoolP("json", "j", false, "output as JSON")
	rootCmd.AddCommand(modelsCmd)
}

func runModels(cmd *cobra.Command, args []string) {
	jsonOutput, _ := cmd.Flags().GetBool("json")

	if jsonOutput {
		data, _ := json.MarshalIndent(availableModels, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Println("Available Models:")
	fmt.Println()
	for _, m := range availableModels {
		rec := ""
		if m.Recommended {
			rec = " *"
		}
		fmt.Printf("  %-25s %s%s\n", m.Name, m.Description, rec)
	}
	fmt.Println()
	fmt.Println("* Recommended")
	fmt.Println()
	fmt.Println("Tip: 'document-parse' alias를 사용하면 자동으로 최신 안정 버전이 적용됩니다.")
}
