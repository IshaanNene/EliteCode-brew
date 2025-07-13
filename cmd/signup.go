package cmd

import (
	"EliteCode-brew/internal/auth"
	"fmt"
	"github.com/spf13/cobra"
)

var SignupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Creating a new EliteCode account",
	Run: func(cmd *cobra.Command, args []string) {
		var name, email, username, password string
		fmt.Print("Full Name: ")
		fmt.Scanln(&name)
		fmt.Print("Email: ")
		fmt.Scanln(&email)
		fmt.Print("Username: ")
		fmt.Scanln(&username)
		fmt.Print("Password: ")
		fmt.Scanln(&password)
		err := auth.SignupUser(name, email, username, password)
		if err != nil {
			fmt.Println("❌ Signup failed:", err)
			return
		}
		fmt.Println("✅ Signup successful! You can now login using your credentials.")
	},
}
