package cmd

import (
	"fmt"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login into your EliteCode account",
	Run: func(cmd *cobra.Command, args []string) {
		var username, password string
		fmt.Print("Username: ")
		fmt.Scanln(&username)
		fmt.Print("Password: ")
		fmt.Scanln(&password)
		err := auth.LoginUser(username, password)
		if err != nil {
			fmt.Println("Login failed:", err)
			return
		}
		fmt.Println("Login successful! type elitecode help to continue")
	},
}
