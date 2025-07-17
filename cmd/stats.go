package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/IshaanNene/EliteCode-brew/internal/problem"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats [problem_id]",
	Short: "Show problem statistics",
	Long:  `Show statistics for a specific problem.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		problemID := args[0]

		problemService := problem.NewService(firebaseClient.Firestore)

		prob, err := problemService.GetProblem(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem: %v", err)
		}

		stats, err := problemService.GetProblemStats(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem stats: %v", err)
		}

		recentSubmissions, err := problemService.GetRecentSubmissions(ctx, problemID, 5)
		if err != nil {
			return fmt.Errorf("error getting recent submissions: %v", err)
		}

		topSolutions, err := problemService.GetTopSolutions(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting top solutions: %v", err)
		}

		fmt.Printf("\nProblem Information:\n")
		fmt.Printf("ID: %s\n", prob.ID)
		fmt.Printf("Title: %s\n", prob.Title)
		fmt.Printf("Difficulty: %s\n", formatDifficulty(prob.Difficulty))
		fmt.Printf("Tags: %s\n", strings.Join(prob.Tags, ", "))

		fmt.Printf("\nStatistics:\n")
		fmt.Printf("Total Submissions: %d\n", stats.TotalSubmissions)
		fmt.Printf("Accepted Submissions: %d\n", stats.AcceptedCount)
		fmt.Printf("Acceptance Rate: %.1f%%\n", stats.AcceptanceRate*100)
		fmt.Printf("Average Time: %d ms\n", stats.AverageTime)
		fmt.Printf("Average Memory: %.2f MB\n", float64(stats.AverageMemory)/1024)

		fmt.Printf("\nRecent Submissions:\n")
		for _, sub := range recentSubmissions {
			var statusColor *color.Color
			switch sub.Status {
			case models.StatusAccepted:
				statusColor = color.New(color.FgGreen)
			case models.StatusRejected:
				statusColor = color.New(color.FgRed)
			default:
				statusColor = color.New(color.FgWhite)
			}

			fmt.Printf("User: %s, Status: %s, Time: %d ms, Memory: %.2f MB, Language: %s, Submitted: %s\n",
				sub.UserID,
				statusColor.Sprint(sub.Status),
				sub.ExecutionTime,
				float64(sub.MemoryUsed)/1024,
				sub.Language,
				sub.SubmittedAt.Format(time.RFC3339),
			)
		}

		fmt.Printf("\nTop Solutions:\n")
		for i, sub := range topSolutions {
			fmt.Printf("#%d - User: %s, Time: %d ms, Memory: %.2f MB, Language: %s\n",
				i+1,
				sub.UserID,
				sub.ExecutionTime,
				float64(sub.MemoryUsed)/1024,
				sub.Language,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
