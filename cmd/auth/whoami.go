package auth

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
)

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user information",
	RunE:  runWhoami,
}

func runWhoami(cmd *cobra.Command, args []string) error {
	config := storage.NewConfigManager()
	
	// Check if user is logged in
	token, err := config.GetAuthToken()
	if err != nil || token == "" {
		color.Red("❌ Not logged in")
		return fmt.Errorf("please login first using 'elitecode auth login'")
	}

	user, err := config.GetUser()
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	// Get fresh user stats from API
	client := api.NewClient()
	stats, err := client.GetUserStats()
	if err != nil {
		color.Yellow("⚠️  Could not fetch latest stats")
	}

	// Display user information
	color.Cyan("👤 User Information")
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("Joined: %s\n", user.CreatedAt.Format("January 2, 2006"))

	if stats != nil {
		fmt.Printf("\n")
		color.Cyan("📊 Statistics")
		fmt.Printf("Problems Solved: %d\n", stats.ProblemsSolved)
		fmt.Printf("Accuracy: %.1f%%\n", stats.Accuracy)
		fmt.Printf("Current Streak: %d days\n", stats.CurrentStreak)
		fmt.Printf("Last Activity: %s\n", stats.LastActivity.Format("January 2, 2006"))
	}

	return nil
}