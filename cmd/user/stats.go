package user

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show user statistics",
	RunE:  runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
	client := api.NewClient()
	stats, err := client.GetUserStats()
	if err != nil {
		return fmt.Errorf("failed to get user stats: %v", err)
	}

	color.Cyan("📊 Your Statistics")
	fmt.Printf("Problems Solved: %d\n", stats.ProblemsSolved)
	fmt.Printf("Total Submissions: %d\n", stats.TotalSubmissions)
	fmt.Printf("Accuracy: %.1f%%\n", stats.Accuracy)
	fmt.Printf("Current Streak: %d days\n", stats.CurrentStreak)
	fmt.Printf("Longest Streak: %d days\n", stats.LongestStreak)
	fmt.Printf("Last Activity: %s\n", stats.LastActivity.Format("January 2, 2006"))

	// Show difficulty breakdown
	if stats.DifficultyBreakdown != nil {
		fmt.Printf("\n")
		color.Cyan("📈 Difficulty Breakdown")
		
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Difficulty", "Solved", "Total", "Percentage"})
		table.SetBorder(true)

		for difficulty, breakdown := range stats.DifficultyBreakdown {
			percentage := 0.0
			if breakdown.Total > 0 {
				percentage = float64(breakdown.Solved) / float64(breakdown.Total) * 100
			}
			
			table.Append([]string{
				difficulty,
				fmt.Sprintf("%d", breakdown.Solved),
				fmt.Sprintf("%d", breakdown.Total),
				fmt.Sprintf("%.1f%%", percentage),
			})
		}
		
		table.Render()
	}

	// Show category breakdown
	if stats.CategoryBreakdown != nil {
		fmt.Printf("\n")
		color.Cyan("📂 Category Breakdown")
		
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Category", "Solved", "Percentage"})
		table.SetBorder(true)

		for category, count := range stats.CategoryBreakdown {
			percentage := 0.0
			if stats.ProblemsSolved > 0 {
				percentage = float64(count) / float64(stats.ProblemsSolved) * 100
			}
			
			table.Append([]string{
				category,
				fmt.Sprintf("%d", count),
				fmt.Sprintf("%.1f%%", percentage),
			})
		}
		
		table.Render()
	}

	return nil
}