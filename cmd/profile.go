package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/models"
)

func formatTable(headers []string, rows [][]string) string {
	var sb strings.Builder
	maxLengths := make([]int, len(headers))

	for i, header := range headers {
		maxLengths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > maxLengths[i] {
				maxLengths[i] = len(cell)
			}
		}
	}

	for i, header := range headers {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%-*s", maxLengths[i], header))
	}
	sb.WriteString("\n")

	for i, length := range maxLengths {
		if i > 0 {
			sb.WriteString("  ")
		}
		sb.WriteString(strings.Repeat("-", length))
	}
	sb.WriteString("\n")

	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString(fmt.Sprintf("%-*s", maxLengths[i], cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "View your profile and statistics",
	Long: `View your profile information and coding statistics.
This includes:
- Solved problems by difficulty
- Success rate
- Average completion time
- Language distribution
- Recent activity`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		userID, err := getUserID()
		if err != nil {
			return err
		}

		userDoc, err := firebaseClient.Firestore.Collection("users").Doc(userID).Get(ctx)
		if err != nil {
			return fmt.Errorf("error getting user profile: %v", err)
		}

		var user struct {
			Username    string    `firestore:"username"`
			Email       string    `firestore:"email"`
			DisplayName string    `firestore:"display_name"`
			CreatedAt   time.Time `firestore:"created_at"`
		}
		if err := userDoc.DataTo(&user); err != nil {
			return fmt.Errorf("error parsing user profile: %v", err)
		}

		summariesSnap, err := firebaseClient.Firestore.Collection("users").Doc(userID).Collection("submission_summaries").Documents(ctx).GetAll()
		if err != nil {
			return fmt.Errorf("error getting submission summaries: %v", err)
		}

		var summaries []models.SubmissionSummary
		for _, doc := range summariesSnap {
			var summary models.SubmissionSummary
			if err := doc.DataTo(&summary); err != nil {
				return fmt.Errorf("error parsing submission summary: %v", err)
			}
			summaries = append(summaries, summary)
		}

		problems := make(map[string]models.Problem)
		for _, summary := range summaries {
			if _, ok := problems[summary.ProblemID]; !ok {
				doc, err := firebaseClient.Firestore.Collection("problems").Doc(summary.ProblemID).Get(ctx)
				if err != nil {
					return fmt.Errorf("error getting problem %s: %v", summary.ProblemID, err)
				}
				var problem models.Problem
				if err := doc.DataTo(&problem); err != nil {
					return fmt.Errorf("error parsing problem %s: %v", summary.ProblemID, err)
				}
				problems[summary.ProblemID] = problem
			}
		}

		stats := calculateStats(summaries, problems)

		fmt.Printf("\nProfile: %s (%s)\n", user.DisplayName, user.Username)
		fmt.Printf("Member since: %s\n", user.CreatedAt.Format("January 2, 2006"))
		fmt.Println()

		fmt.Println("Problem Statistics:")
		fmt.Printf("Total Solved: %d\n", stats.TotalSolved)
		fmt.Printf("Success Rate: %.1f%%\n", stats.SuccessRate*100)
		fmt.Println()

		fmt.Println("Difficulty Distribution:")
		difficultyHeaders := []string{"Difficulty", "Solved", "Total", "Percentage"}
		var difficultyRows [][]string

		difficulties := []models.Difficulty{models.Easy, models.Medium, models.Hard, models.VeryHard}
		for _, diff := range difficulties {
			solved := stats.SolvedByDifficulty[diff]
			total := stats.TotalByDifficulty[diff]
			percentage := 0.0
			if total > 0 {
				percentage = float64(solved) / float64(total) * 100
			}

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

			row := []string{
				diffColor.Sprint(diff),
				fmt.Sprintf("%d", solved),
				fmt.Sprintf("%d", total),
				fmt.Sprintf("%.1f%%", percentage),
			}
			difficultyRows = append(difficultyRows, row)
		}
		fmt.Print(formatTable(difficultyHeaders, difficultyRows))
		fmt.Println()

		fmt.Println("Language Distribution:")
		languageHeaders := []string{"Language", "Problems Solved", "Percentage"}
		var languageRows [][]string

		languages := make([]string, 0, len(stats.SolvedByLanguage))
		for lang := range stats.SolvedByLanguage {
			languages = append(languages, lang)
		}
		sort.Slice(languages, func(i, j int) bool {
			return stats.SolvedByLanguage[languages[i]] > stats.SolvedByLanguage[languages[j]]
		})

		for _, lang := range languages {
			solved := stats.SolvedByLanguage[lang]
			percentage := float64(solved) / float64(stats.TotalSolved) * 100
			row := []string{
				strings.Title(lang),
				fmt.Sprintf("%d", solved),
				fmt.Sprintf("%.1f%%", percentage),
			}
			languageRows = append(languageRows, row)
		}
		fmt.Print(formatTable(languageHeaders, languageRows))
		fmt.Println()

		fmt.Println("Recent Activity:")
		activityHeaders := []string{"Date", "Problem", "Status", "Language", "Time", "Memory"}
		var activityRows [][]string

		sort.Slice(summaries, func(i, j int) bool {
			return summaries[i].LastAttemptedAt.After(summaries[j].LastAttemptedAt)
		})

		for i, summary := range summaries {
			if i >= 10 { // Show only last 10 activities
				break
			}

			problem := problems[summary.ProblemID]
			status := "Attempted"
			if summary.Solved {
				status = "Solved"
			}

			row := []string{
				summary.LastAttemptedAt.Format("2006-01-02"),
				problem.Title,
				status,
				strings.Title(summary.Language),
				fmt.Sprintf("%dms", summary.BestTime),
				fmt.Sprintf("%.2fMB", float64(summary.BestMemory)/1024),
			}
			activityRows = append(activityRows, row)
		}
		fmt.Print(formatTable(activityHeaders, activityRows))

		return nil
	},
}

type statistics struct {
	TotalSolved        int
	TotalAttempted     int
	SuccessRate        float64
	SolvedByDifficulty map[models.Difficulty]int
	TotalByDifficulty  map[models.Difficulty]int
	SolvedByLanguage   map[string]int
	AverageTime        int64
	AverageMemory      int64
}

func calculateStats(summaries []models.SubmissionSummary, problems map[string]models.Problem) statistics {
	stats := statistics{
		SolvedByDifficulty: make(map[models.Difficulty]int),
		TotalByDifficulty:  make(map[models.Difficulty]int),
		SolvedByLanguage:   make(map[string]int),
	}

	for _, summary := range summaries {
		problem := problems[summary.ProblemID]
		stats.TotalAttempted++

		if summary.Solved {
			stats.TotalSolved++
			stats.SolvedByDifficulty[problem.Difficulty]++
			stats.SolvedByLanguage[summary.Language]++
			stats.AverageTime += summary.BestTime
			stats.AverageMemory += summary.BestMemory
		}
		stats.TotalByDifficulty[problem.Difficulty]++
	}

	if stats.TotalAttempted > 0 {
		stats.SuccessRate = float64(stats.TotalSolved) / float64(stats.TotalAttempted)
	}

	if stats.TotalSolved > 0 {
		stats.AverageTime /= int64(stats.TotalSolved)
		stats.AverageMemory /= int64(stats.TotalSolved)
	}

	return stats
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
