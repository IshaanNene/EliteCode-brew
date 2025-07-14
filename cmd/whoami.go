package cmd

import (
	"fmt"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
	"github.com/spf13/cobra"
)

var WhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show currently logged-in user",
	Run: func(cmd *cobra.Command, args []string) {
		user, err := auth.GetCurrentUser()
		if err != nil {
			fmt.Println("Not logged in.")
			return
		}
		fmt.Printf("Logged in as: %s\n", user.Username)
	},
}
