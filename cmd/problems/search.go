package problems

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search problems by keyword",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().StringP("difficulty", "d", "", "Filter by difficulty (easy, medium, hard)")
	searchCmd.Flags().StringP("category", "c", "", "Filter by category")
	searchCmd.Flags().IntP("limit", "l", 10, "Limit number of results")
}

func runSearch(cmd *cobra.Command, args []string) error {
	keyword := strings.Join(args, " ")
	difficulty, _ := cmd.Flags().GetString("difficulty")
	category, _ := cmd.Flags().GetString("category")
	limit, _ := cmd.Flags().GetInt("limit")

	client := api.NewClient()
	problems, err := client.SearchProblems(keyword, difficulty, category, limit)
	if err != nil {
		return fmt.Errorf("failed to search problems: %v", err)
	}

	if len(problems) == 0 {
		color.Yellow("No problems found matching your search criteria")
		return nil
	}

	color.Green("Found %d problem(s) matching '%s':", len(problems), keyword)
	
	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Difficulty", "Category", "Solved"})
	table.SetBorder(true)

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