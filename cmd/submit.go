package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/docker"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/IshaanNene/EliteCode-brew/internal/problem"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit your solution",
	Long: `Submit your solution for evaluation.
This will:
1. Run your code against all test cases
2. Check for correctness, time limit, and memory usage
3. Update your submission history`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %v", err)
		}

		metadataBytes, err := os.ReadFile(filepath.Join(wd, "metadata.json"))
		if err != nil {
			return fmt.Errorf("error reading metadata file: %v", err)
		}

		var meta struct {
			ID          string `json:"id"`
			Title       string `json:"title"`
			Language    string `json:"language"`
			TimeLimit   int    `json:"time_limit"`
			MemoryLimit int    `json:"memory_limit"`
		}
		if err := json.Unmarshal(metadataBytes, &meta); err != nil {
			return fmt.Errorf("error parsing metadata: %v", err)
		}

		sourceCode, err := os.ReadFile(filepath.Join(wd, "main."+utils.GetFileExtension(meta.Language)))
		if err != nil {
			return fmt.Errorf("error reading source code: %v", err)
		}

		dockerClient, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("error creating Docker client: %v", err)
		}
		defer dockerClient.Close()

		problemService := problem.NewService(firebaseClient.Firestore)
		submissionService := problem.NewSubmissionService(firebaseClient.Firestore, dockerClient, problemService)

		configFile := filepath.Join(os.Getenv("HOME"), ".elitecode", "config.json")
		configBytes, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %v", err)
		}

		var config struct {
			UID string `json:"uid"`
		}
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return fmt.Errorf("error parsing config: %v", err)
		}

		if config.UID == "" {
			return fmt.Errorf("not logged in. Please run 'elitecode login' first")
		}

		fmt.Println("Submitting your solution...")
		startTime := time.Now()

		submission, err := submissionService.Submit(cmd.Context(), config.UID, meta.ID, string(sourceCode), meta.Language)
		if err != nil {
			return fmt.Errorf("error submitting solution: %v", err)
		}

		fmt.Println("\nResults:")
		fmt.Println("----------------------------------------")

		passedTests := 0
		for i, tc := range submission.TestCases {
			fmt.Printf("Test Case %d: ", i+1)

			isHidden := strings.HasPrefix(tc.TestCaseID, "hidden_")

			switch tc.Status {
			case models.StatusAccepted:
				color.Green("✓ Passed")
				passedTests++
			case models.StatusRejected:
				color.Red("✗ Failed")
				if !isHidden {
					fmt.Printf("  Expected: %s\n", tc.ExpectedOutput)
					fmt.Printf("  Got:      %s\n", tc.ActualOutput)
				}
			case models.StatusTimeLimitExceeded:
				color.Yellow("⚠ Time Limit Exceeded")
				fmt.Printf("  Time: %dms (limit: %dms)\n", tc.ExecutionTime, meta.TimeLimit)
			case models.StatusMemoryLimitExceeded:
				color.Yellow("⚠ Memory Limit Exceeded")
				fmt.Printf("  Memory: %.2fMB (limit: %dMB)\n", float64(tc.MemoryUsed)/1024, meta.MemoryLimit)
			case models.StatusError:
				color.Red("✗ Error")
				fmt.Printf("  %s\n", tc.ErrorMessage)
			}
		}

		fmt.Println("----------------------------------------")

		fmt.Printf("\nTest Cases: %d/%d passed\n", passedTests, len(submission.TestCases))
		fmt.Printf("Average Time: %dms\n", submission.ExecutionTime)
		fmt.Printf("Peak Memory: %.2fMB\n", float64(submission.MemoryUsed)/1024)
		fmt.Printf("Total Time: %v\n", time.Since(startTime))

		if submission.Status == models.StatusAccepted {
			color.Green("\nCongratulations! All test cases passed.")
		} else {
			color.Red("\nSome test cases failed. Keep trying!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
