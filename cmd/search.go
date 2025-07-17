package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/IshaanNene/EliteCode-brew/internal/problem"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for problems",
	Long: `Search for problems by title, description, or tags.
You can filter by:
- Difficulty (easy, medium, hard, very_hard)
- Tags (comma-separated)
- Status (attempted, solved, bookmarked)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		query := make(map[string]interface{})
		if len(args) > 0 {
			query["text"] = args[0]
		}
		if searchDifficulty != "" {
			query["difficulty"] = models.Difficulty(searchDifficulty)
		}
		if searchTags != "" {
			tags := strings.Split(searchTags, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			query["tags"] = tags
		}

		var userID string
		if searchStatus != "" {
			var err error
			userID, err = getUserID()
			if err != nil {
				return err
			}
		}

		problemService := problem.NewService(firebaseClient.Firestore)
		problems, err := problemService.ListProblems(ctx, query)
		if err != nil {
			return fmt.Errorf("error searching problems: %v", err)
		}

		if searchStatus != "" {
			var filteredProblems []models.Problem
			for _, prob := range problems {
				status, err := problemService.GetUserProblemStatus(ctx, userID, prob.ID)
				if err != nil {
					return fmt.Errorf("error getting problem status: %v", err)
				}
				if status != nil {
					switch searchStatus {
					case "attempted":
						if status.AttemptCount > 0 {
							filteredProblems = append(filteredProblems, prob)
						}
					case "solved":
						if status.SubmissionSummary != nil && status.SubmissionSummary.Solved {
							filteredProblems = append(filteredProblems, prob)
						}
					case "bookmarked":
						if status.Bookmarked {
							filteredProblems = append(filteredProblems, prob)
						}
					}
				}
			}
			problems = filteredProblems
		}

		if len(problems) == 0 {
			fmt.Println("No problems found.")
			return nil
		}

		headers := []string{"ID", "Title", "Difficulty", "Tags", "Acceptance Rate"}
		var rows [][]string

		for _, prob := range problems {
			stats, err := problemService.GetProblemStats(ctx, prob.ID)
			if err != nil {
				return fmt.Errorf("error getting problem stats: %v", err)
			}

			var diffColor *color.Color
			switch prob.Difficulty {
			case models.Easy:
				diffColor = color.New(color.FgGreen)
			case models.Medium:
				diffColor = color.New(color.FgYellow)
			case models.Hard:
				diffColor = color.New(color.FgRed)
			case models.VeryHard:
				diffColor = color.New(color.FgRed, color.Bold)
			}

			row := []string{
				prob.ID,
				prob.Title,
				diffColor.Sprint(prob.Difficulty),
				strings.Join(prob.Tags, ", "),
				fmt.Sprintf("%.1f%%", stats.AcceptanceRate*100),
			}
			rows = append(rows, row)
		}

		fmt.Printf("Found %d problems:\n", len(problems))
		fmt.Print(formatTable(headers, rows))

		return nil
	},
}

var (
	searchDifficulty string
	searchTags       string
	searchStatus     string
)

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVarP(&searchDifficulty, "difficulty", "d", "", "Filter by difficulty (easy, medium, hard, very_hard)")
	searchCmd.Flags().StringVarP(&searchTags, "tags", "t", "", "Filter by tags (comma-separated)")
	searchCmd.Flags().StringVarP(&searchStatus, "status", "s", "", "Filter by status (attempted, solved, bookmarked)")
}
