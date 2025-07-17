package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/IshaanNene/EliteCode-brew/internal/problem"
)

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark [problem_id]",
	Short: "Manage bookmarked problems",
	Long: `Manage your bookmarked problems.
Without arguments, lists all bookmarked problems.
With a problem ID, toggles the bookmark status.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		userID, err := getUserID()
		if err != nil {
			return err
		}

		problemService := problem.NewService(firebaseClient.Firestore)

		if len(args) == 0 {
			statuses, err := problemService.GetUserProblemStatuses(ctx, userID)
			if err != nil {
				return fmt.Errorf("error getting problem statuses: %v", err)
			}

			var bookmarkedProblems []models.Problem
			for _, status := range statuses {
				if status.Bookmarked {
					prob, err := problemService.GetProblem(ctx, status.ProblemID)
					if err != nil {
						return fmt.Errorf("error getting problem %s: %v", status.ProblemID, err)
					}
					bookmarkedProblems = append(bookmarkedProblems, *prob)
				}
			}

			if len(bookmarkedProblems) == 0 {
				fmt.Println("No bookmarked problems.")
				return nil
			}

			fmt.Printf("Found %d bookmarked problems:\n", len(bookmarkedProblems))
			headers := []string{"ID", "Title", "Difficulty", "Tags", "Status"}
			var rows [][]string

			for _, prob := range bookmarkedProblems {
				status, err := problemService.GetUserProblemStatus(ctx, userID, prob.ID)
				if err != nil {
					return fmt.Errorf("error getting problem status: %v", err)
				}

				var statusStr string
				if status.SubmissionSummary != nil && status.SubmissionSummary.Solved {
					statusStr = color.GreenString("Solved")
				} else if status.AttemptCount > 0 {
					statusStr = color.YellowString("Attempted")
				} else {
					statusStr = "Not Started"
				}

				row := []string{
					prob.ID,
					prob.Title,
					formatDifficulty(prob.Difficulty),
					strings.Join(prob.Tags, ", "),
					statusStr,
				}
				rows = append(rows, row)
			}

			fmt.Print(formatTable(headers, rows))
			return nil
		}

		problemID := args[0]

		status, err := problemService.GetUserProblemStatus(ctx, userID, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem status: %v", err)
		}

		prob, err := problemService.GetProblem(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem: %v", err)
		}

		if status == nil {
			status = &models.UserProblemStatus{
				UserID:    userID,
				ProblemID: problemID,
			}
		}

		status.Bookmarked = !status.Bookmarked

		if err := problemService.SaveUserProblemStatus(ctx, status); err != nil {
			return fmt.Errorf("error saving problem status: %v", err)
		}

		if status.Bookmarked {
			fmt.Printf("Bookmarked problem: %s - %s\n", prob.ID, prob.Title)
		} else {
			fmt.Printf("Removed bookmark: %s - %s\n", prob.ID, prob.Title)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(bookmarkCmd)
}
