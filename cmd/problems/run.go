package problems

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/docker"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run your solution",
	RunE:  runRun,
}

func init() {
	runCmd.Flags().StringP("language", "l", "", "Programming language (auto-detect if not specified)")
	runCmd.Flags().StringP("input", "i", "input.txt", "Input file")
	runCmd.Flags().BoolP("time", "t", true, "Show execution time")
}

func runRun(cmd *cobra.Command, args []string) error {
	language, _ := cmd.Flags().GetString("language")
	inputFile, _ := cmd.Flags().GetString("input")
	showTime, _ := cmd.Flags().GetBool("time")

	// Auto-detect language if not specified
	if language == "" {
		detectedLang, err := detectLanguage()
		if err != nil {
			return fmt.Errorf("failed to detect language: %v", err)
		}
		language = detectedLang
	}

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file %s not found", inputFile)
	}

	// Find source file
	sourceFile := getSourceFileName(language)
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source file %s not found", sourceFile)
	}

	spinner := utils.NewSpinner("Running your solution...")
	spinner.Start()

	// Run in Docker container
	runner := docker.NewRunner()
	
	start := time.Now()
	output, err := runner.Run(language, sourceFile, inputFile)
	duration := time.Since(start)
	
	spinner.Stop()

	if err != nil {
		color.Red("❌ Execution failed: %v", err)
		return nil
	}

	color.Green("✅ Execution completed")
	
	if showTime {
		color.Blue("⏱️  Execution time: %v", duration)
	}

	color.Cyan("📤 Output:")
	fmt.Println(output)

	return nil
}

func detectLanguage() (string, error) {
	// Check for common source files
	extensions := map[string]string{
		".c":   "c",
		".cpp": "cpp",
		".py":  "python",
		".java": "java",
		".js":  "javascript",
	}

	for ext, lang := range extensions {
		matches, err := filepath.Glob("*" + ext)
		if err != nil {
			continue
		}
		if len(matches) > 0 {
			return lang, nil
		}
	}

	return "", fmt.Errorf("no source files found")
}