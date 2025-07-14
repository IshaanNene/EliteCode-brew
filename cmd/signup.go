package cmd

import (
	"fmt"

	"github.com/IshaanNene/EliteCode-brew/internal/auth"
	"github.com/IshaanNene/EliteCode-brew/internal/db"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
	"github.com/spf13/cobra"
)

var SignupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Register a new EliteCode account",
	Run: func(cmd *cobra.Command, args []string) {
		var name, email, username, password string

		fmt.Print("Name: ")
		fmt.Scanln(&name)
		fmt.Print("Email: ")
		fmt.Scanln(&email)
		fmt.Print("Username: ")
		fmt.Scanln(&username)
		fmt.Print("Password: ")
		fmt.Scanln(&password)

		// Hash the password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			fmt.Println("Failed to hash password:", err)
			return
		}

		// Firebase connection
		dbInstance, err := auth.NewFirebaseDB("config/serviceAccountKey.json")
		if err != nil {
			fmt.Println("Failed to connect to Firebase:", err)
			return
		}
		defer dbInstance.Close()

		store := db.NewStore(dbInstance)

		user := &auth.User{
			Name:     name,
			Email:    email,
			Username: username,
			Password: hashedPassword,
		}

		err = store.SaveUser(user)
		if err != nil {
			fmt.Println("Signup failed:", err)
			return
		}

		fmt.Println("Account created! Type 'elitecode login' to continue.")
	},
}
