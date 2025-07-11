package problems

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available problems",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringP("category", "c", "", "Filter by category")
	listCmd.Flags().StringP("difficulty", "d", "", "Filter by difficulty (easy, medium, hard)")
	listCmd.Flags().IntP("limit", "l", 20, "Limit number of problems shown")
}

func runList(cmd *cobra.Command, args []string) error {
	category, _ := cmd.Flags().GetString("category")
	difficulty, _ := cmd.Flags().GetString("difficulty")
	limit, _ := cmd.Flags().GetInt("limit")

	client := api.NewClient()
	problems, err := client.GetProblems(category, difficulty, limit)
	if err != nil {
		return fmt.Errorf("failed to fetch problems: %v", err)
	}

	if len(problems) == 0 {
		color.Yellow("No problems found matching your criteria")
		return nil
	}

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Difficulty", "Category", "Solved"})
	table.SetBorder(true)
	table.SetRowSeparator("-")
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")

	for _, problem := range problems {
		solved := "❌"
		if problem.Solved {
			solved = "✅"
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

		table.Append([]string{
			fmt.Sprintf("%d", problem.ID),
			problem.Title,
			difficulty,
			problem.Category,
			solved,
		})
	}

	table.Render()
	return nil
}