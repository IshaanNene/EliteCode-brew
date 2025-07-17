package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/yourusername/elitecode/internal/models"
	"github.com/yourusername/elitecode/internal/problem"
)

var statsCmd = &cobra.Command{
	Use:   "stats [problem_id]",
	Short: "View problem statistics",
	Long: `View detailed statistics for a problem.
This includes:
- Acceptance rate
- Average completion time
- Memory usage
- Recent submissions
- Top solutions`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		problemID := args[0]

		// Get problem service
		problemService := problem.NewService(firebaseClient.Firestore)

		// Get problem details
		prob, err := problemService.GetProblem(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem: %v", err)
		}

		// Get problem stats
		stats, err := problemService.GetProblemStats(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem stats: %v", err)
		}

		// Get recent submissions
		submissions, err := problemService.GetRecentSubmissions(ctx, problemID, 10)
		if err != nil {
			return fmt.Errorf("error getting recent submissions: %v", err)
		}

		// Get top solutions
		topSolutions, err := problemService.GetTopSolutions(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting top solutions: %v", err)
		}

		// Print problem information
		fmt.Printf("\nProblem: %s\n", prob.Title)
		fmt.Printf("Difficulty: %s\n", formatDifficulty(prob.Difficulty))
		fmt.Printf("Tags: %s\n", strings.Join(prob.Tags, ", "))
		fmt.Println()

		// Print statistics
		fmt.Println("Statistics:")
		fmt.Printf("Total Submissions: %d\n", stats.TotalSubmissions)
		fmt.Printf("Accepted: %d (%.1f%%)\n", stats.AcceptedCount, stats.AcceptanceRate*100)
		fmt.Printf("Average Time: %dms\n", stats.AverageTime)
		fmt.Printf("Average Memory: %.2fMB\n", float64(stats.AverageMemory)/1024)
		fmt.Printf("Last Updated: %s\n", stats.UpdatedAt.Format(time.RFC822))
		fmt.Println()

		// Print recent submissions
		fmt.Println("Recent Submissions:")
		submissionHeaders := []string{"User", "Language", "Status", "Time", "Memory", "Submitted"}
		var submissionRows [][]string

		for _, sub := range submissions {
			status := string(sub.Status)
			switch sub.Status {
			case models.StatusAccepted:
				status = color.GreenString(status)
			case models.StatusRejected:
				status = color.RedString(status)
			case models.StatusError:
				status = color.YellowString(status)
			}

			row := []string{
				sub.UserID,
				sub.Language,
				status,
				fmt.Sprintf("%dms", sub.ExecutionTime),
				fmt.Sprintf("%.2fMB", float64(sub.MemoryUsed)/1024),
				sub.SubmittedAt.Format("2006-01-02 15:04:05"),
			}
			submissionRows = append(submissionRows, row)
		}
		fmt.Print(formatTable(submissionHeaders, submissionRows))
		fmt.Println()

		// Print top solutions
		fmt.Println("Top Solutions:")
		solutionHeaders := []string{"User", "Language", "Time", "Memory", "Submitted"}
		var solutionRows [][]string

		for _, sol := range topSolutions {
			row := []string{
				sol.UserID,
				sol.Language,
				fmt.Sprintf("%dms", sol.ExecutionTime),
				fmt.Sprintf("%.2fMB", float64(sol.MemoryUsed)/1024),
				sol.SubmittedAt.Format("2006-01-02 15:04:05"),
			}
			solutionRows = append(solutionRows, row)
		}
		fmt.Print(formatTable(solutionHeaders, solutionRows))

		return nil
	},
}

func formatDifficulty(diff models.Difficulty) string {
	var diffColor *color.Color
	switch diff {
	case models.Easy:
		diffColor = color.New(color.FgGreen)
	case models.Medium:
		diffColor = color.New(color.FgYellow)
	case models.Hard:
		diffColor = color.New(color.FgRed)
	case models.VeryHard:
		diffColor = color.New(color.FgRed, color.Bold)
	}
	return diffColor.Sprint(diff)
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
