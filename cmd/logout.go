package cmd

import (
	"fmt"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
	"github.com/spf13/cobra"
)

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout of your EliteCode account",
	Run: func(cmd *cobra.Command, args []string) {
		err := auth.LogoutUser()
		if err != nil {
			fmt.Println("Logout failed:", err)
			return
		}
		fmt.Println("Logged out successfully.")
	},
}
