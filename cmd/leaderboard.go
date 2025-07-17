package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
	"github.com/IshaanNene/EliteCode-brew/internal/problem"
)

var (
	timeRange     string
	problemFilter string
)

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "View global and problem-specific rankings",
	Long: `View the leaderboard of top performers.
You can filter by:
- Time range (all, year, month, week)
- Problem ID (for problem-specific rankings)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		problemService := problem.NewService(firebaseClient.Firestore)

		var startTime time.Time
		switch timeRange {
		case "year":
			startTime = time.Now().AddDate(-1, 0, 0)
		case "month":
			startTime = time.Now().AddDate(0, -1, 0)
		case "week":
			startTime = time.Now().AddDate(0, 0, -7)
		}

		var rankings []models.UserRanking
		var err error
		if problemFilter != "" {
			rankings, err = problemService.GetProblemRankings(ctx, problemFilter, startTime)
			if err != nil {
				return fmt.Errorf("error getting problem rankings: %v", err)
			}

			prob, err := problemService.GetProblem(ctx, problemFilter)
			if err != nil {
				return fmt.Errorf("error getting problem: %v", err)
			}

			fmt.Printf("\nLeaderboard for %s\n", prob.Title)
			fmt.Printf("Difficulty: %s\n", formatDifficulty(prob.Difficulty))
			fmt.Printf("Tags: %s\n", strings.Join(prob.Tags, ", "))
		} else {
			rankings, err = problemService.GetGlobalRankings(ctx, startTime)
			if err != nil {
				return fmt.Errorf("error getting global rankings: %v", err)
			}

			fmt.Printf("\nGlobal Leaderboard")
			if timeRange != "all" {
				fmt.Printf(" (%s)", timeRange)
			}
			fmt.Println()
		}

		if len(rankings) == 0 {
			fmt.Println("No rankings found.")
			return nil
		}

		headers := []string{"Rank", "User", "Score", "Problems Solved", "Avg Time", "Avg Memory"}
		var rows [][]string

		for i, ranking := range rankings {
			rank := fmt.Sprintf("%d", i+1)
			switch i {
			case 0:
				rank = "🥇 " + rank
			case 1:
				rank = "🥈 " + rank
			case 2:
				rank = "🥉 " + rank
			}

			avgTime := float64(ranking.TotalTime) / float64(ranking.ProblemsSolved)
			avgMemory := float64(ranking.TotalMemory) / float64(ranking.ProblemsSolved) / 1024 // Convert to MB

			row := []string{
				rank,
				ranking.Username,
				fmt.Sprintf("%.0f", ranking.Score),
				fmt.Sprintf("%d", ranking.ProblemsSolved),
				fmt.Sprintf("%.1fms", avgTime),
				fmt.Sprintf("%.2fMB", avgMemory),
			}

			userID, _ := getUserID()
			if ranking.UserID == userID {
				for i := range row {
					row[i] = color.CyanString(row[i])
				}
			}

			rows = append(rows, row)
		}

		fmt.Print(formatTable(headers, rows))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(leaderboardCmd)

	leaderboardCmd.Flags().StringVarP(&timeRange, "time", "t", "all", "Time range (all, year, month, week)")
	leaderboardCmd.Flags().StringVarP(&problemFilter, "problem", "p", "", "Filter by problem ID")
}
