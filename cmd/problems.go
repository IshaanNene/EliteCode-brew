package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/IshaanNene/EliteCode-brew/internal/problem"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var problemsCmd = &cobra.Command{
	Use:   "problems",
	Short: "List available problems",
	Long:  `List all available coding problems.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		userID, err := getUserID()
		if err != nil {
			return err
		}

		problemService := problem.NewService(firebaseClient.Firestore)

		problems, err := problemService.ListProblems(ctx, nil)
		if err != nil {
			return fmt.Errorf("error getting problems: %v", err)
		}

		if len(problems) == 0 {
			fmt.Println("No problems found.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{"ID", "Title", "Difficulty", "Tags", "Status"})

		for _, prob := range problems {
			status, err := problemService.GetUserProblemStatus(ctx, userID, prob.ID)
			if err != nil {
				return fmt.Errorf("error getting problem status: %v", err)
			}

			var statusStr string
			if status != nil && status.SubmissionSummary != nil && status.SubmissionSummary.Solved {
				statusStr = color.GreenString("Solved")
			} else if status != nil && status.AttemptCount > 0 {
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
			table.Append(row)
		}

		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(problemsCmd)
}
