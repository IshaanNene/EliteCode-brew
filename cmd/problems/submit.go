package problems

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
	"github.com/IshaanNene/EliteCode-brew/internal/docker"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var submitCmd = &cobra.Command{
	Use:   "submit [problem_id]",
	Short: "Submit your solution",
	RunE:  runSubmit,
}

func init() {
	submitCmd.Flags().StringP("language", "l", "", "Programming language (auto-detect if not specified)")
	submitCmd.Flags().BoolP("force", "f", false, "Force submit even if local tests fail")
}

func runSubmit(cmd *cobra.Command, args []string) error {
	language, _ := cmd.Flags().GetString("language")
	force, _ := cmd.Flags().GetBool("force")

	// Get problem ID
	var problemID int
	var err error
	
	if len(args) > 0 {
		problemID, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid problem ID: %v", err)
		}
	} else {
		// Try to get from current directory
		problemID, err = getProblemIDFromDirectory()
		if err != nil {
			return fmt.Errorf("could not determine problem ID: %v", err)
		}
	}

	// Auto-detect language if not specified
	if language == "" {
		detectedLang, err := detectLanguage()
		if err != nil {
			return fmt.Errorf("failed to detect language: %v", err)
		}
		language = detectedLang
	}

	// Find source file
	sourceFile := getSourceFileName(language)
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source file %s not found", sourceFile)
	}

	// Read source code
	sourceCode, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	// Run local tests first
	spinner := utils.NewSpinner("Running local tests...")
	spinner.Start()

	cache := storage.NewCacheManager()
	testCases, err := cache.GetTestCases(problemID)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to get test cases: %v", err)
	}

	runner := docker.NewRunner()
	passedTests := 0
	failedTests := 0

	for i, testCase := range testCases {
		// Create temporary input file
		tempInput := fmt.Sprintf("temp_input_%d.txt", i)
		if err := os.WriteFile(tempInput, []byte(testCase.Input), 0644); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create temporary input file: %v", err)
		}
		defer os.Remove(tempInput)

		// Run test
		output, err := runner.Run(language, sourceFile, tempInput)
		if err != nil {
			failedTests++
			continue
		}

		// Compare output
		if strings.TrimSpace(output) == strings.TrimSpace(testCase.ExpectedOutput) {
			passedTests++
		} else {
			failedTests++
		}
	}

	spinner.Stop()

	// Display test results
	color.Cyan("📊 Test Results:")
	color.Green("✅ Passed: %d", passedTests)
	if failedTests > 0 {
		color.Red("❌ Failed: %d", failedTests)
	}

	// Check if should proceed with submission
	if failedTests > 0 && !force {
		color.Yellow("⚠️  Some tests failed. Use --force to submit anyway.")
		return nil
	}

	// Submit to server
	spinner = utils.NewSpinner("Submitting solution...")
	spinner.Start()

	client := api.NewClient()
	submission, err := client.SubmitSolution(problemID, string(sourceCode), language)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to submit solution: %v", err)
	}

	spinner.Stop()

	// Display submission results
	color.Green("🚀 Solution submitted successfully!")
	color.Blue("📝 Submission ID: %d", submission.ID)
	color.Blue("⏱️  Server execution time: %v", submission.ExecutionTime)
	color.Blue("💾 Memory usage: %s", submission.MemoryUsage)

	if submission.Accepted {
		color.Green("🎉 ACCEPTED!")
	} else {
		color.Red("❌ REJECTED")
		if submission.ErrorMessage != "" {
			color.Red("Error: %s", submission.ErrorMessage)
		}
	}

	return nil
}

func getProblemIDFromDirectory() (int, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	dirName := filepath.Base(pwd)
	// Extract problem ID from directory name (format: "1001_two_sum")
	parts := strings.Split(dirName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid directory format")
	}

	return strconv.Atoi(parts[0])
}