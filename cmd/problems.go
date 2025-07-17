package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yourusername/elitecode/internal/models"
	"github.com/yourusername/elitecode/internal/problem"
)

var (
	difficultyFilter string
	tagFilter        string
	statusFilter     string
)

var problemsCmd = &cobra.Command{
	Use:   "problems",
	Short: "List available coding problems",
	Long: `List all available coding problems with optional filters.
You can filter by:
- Difficulty (easy, medium, hard, very_hard)
- Tags (comma-separated)
- Status (attempted, solved, bookmarked)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create problem service
		problemService := problem.NewService(firebaseClient.Firestore)

		// Build filters
		filters := make(map[string]interface{})
		if difficultyFilter != "" {
			filters["difficulty"] = models.Difficulty(difficultyFilter)
		}
		if tagFilter != "" {
			tags := strings.Split(tagFilter, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			filters["tags"] = tags
		}

		// Get problems
		ctx := cmd.Context()
		problems, err := problemService.ListProblems(ctx, filters)
		if err != nil {
			return fmt.Errorf("error listing problems: %v", err)
		}

		// Create table
		table := tablewriter.NewWriter(cmd.OutOrStdout())
		table.SetHeader([]string{"ID", "Title", "Difficulty", "Tags", "Status", "Acceptance Rate"})
		table.SetBorder(false)
		table.SetColumnSeparator(" ")

		// Set colors for different difficulties
		difficultyColors := map[models.Difficulty]*color.Color{
			models.Easy:     color.New(color.FgGreen),
			models.Medium:   color.New(color.FgYellow),
			models.Hard:     color.New(color.FgRed),
			models.VeryHard: color.New(color.FgRed, color.Bold),
		}

		// Add problems to table
		for _, p := range problems {
			// Get problem stats
			stats, err := problemService.GetProblemStats(ctx, p.ID)
			if err != nil {
				return fmt.Errorf("error getting problem stats: %v", err)
			}

			// Get user's status for this problem if logged in
			var status string
			if firebaseClient != nil {
				// TODO: Get user ID from config
				userID := "current_user_id"
				userStatus, err := problemService.GetUserProblemStatus(ctx, userID, p.ID)
				if err != nil {
					return fmt.Errorf("error getting user problem status: %v", err)
				}
				if userStatus != nil {
					status = userStatus.Status
				}
			}

			// Format difficulty with color
			diffColor := difficultyColors[p.Difficulty]
			difficulty := diffColor.Sprint(p.Difficulty)

			// Format tags
			tags := strings.Join(p.Tags, ", ")

			// Format acceptance rate
			acceptanceRate := fmt.Sprintf("%.1f%%", stats.AcceptanceRate*100)

			// Add row to table
			row := []string{
				p.ID,
				p.Title,
				difficulty,
				tags,
				status,
				acceptanceRate,
			}
			table.Append(row)
		}

		// Render table
		table.Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(problemsCmd)

	// Add flags
	problemsCmd.Flags().StringVarP(&difficultyFilter, "difficulty", "d", "", "Filter by difficulty (easy, medium, hard, very_hard)")
	problemsCmd.Flags().StringVarP(&tagFilter, "tags", "t", "", "Filter by tags (comma-separated)")
	problemsCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "Filter by status (attempted, solved, bookmarked)")
}
