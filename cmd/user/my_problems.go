package user

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var myProblemsCmd = &cobra.Command{
	Use:   "my-problems",
	Short: "Show your problem history",
	RunE:  runMyProblems,
}

func init() {
	myProblemsCmd.Flags().StringP("status", "s", "", "Filter by status (solved, attempted, bookmarked)")
	myProblemsCmd.Flags().IntP("limit", "l", 20, "Limit number of problems shown")
}

func runMyProblems(cmd *cobra.Command, args []string) error {
	status, _ := cmd.Flags().GetString("status")
	limit, _ := cmd.Flags().GetInt("limit")

	client := api.NewClient()
	problems, err := client.GetUserProblems(status, limit)
	if err != nil {
		return fmt.Errorf("failed to get user problems: %v", err)
	}

	if len(problems) == 0 {
		color.Yellow("No problems found")
		return nil
	}

	var title string
	switch status {
	case "solved":
		title = "✅ Solved Problems"
	case "attempted":
		title = "🔄 Attempted Problems"
	case "bookmarked":
		title = "🔖 Bookmarked Problems"
	default:
		title = "📚 Your Problems"
	}

	color.Cyan(title)
	
	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Difficulty", "Status", "Last Attempt"})
	table.SetBorder(true)

	for _, problem := range problems {
		status := "❌ Not Solved"
		if problem.Solved {
			status = "✅ Solved"
		} else if problem.Attempted {
			status = "🔄 Attempted"
		}
		
		difficulty := problem.Difficulty
		switch difficulty {
		case "easy":
			difficulty = color.GreenString("Easy")
		case "medium":
			difficulty = color.YellowString("Medium")
		case "hard":
			difficulty = color.RedString("Hard")
		}

		lastAttempt := "Never"
		if problem.LastAttempt != nil {
			lastAttempt = problem.LastAttempt.Format("Jan 2, 2006")
		}

		table.Append([]string{
			fmt.Sprintf("%d", problem.ID),
			problem.Title,
			difficulty,
			status,
			lastAttempt,
		})
	}

	table.Render()
	return nil
}