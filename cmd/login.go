package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to your Elitecode account",
	Long:  `Log in to your Elitecode account using your email and password.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		email, err := promptForInput("Email", validateEmail)
		if err != nil {
			return err
		}

		password, err := promptForPassword("Password", nil)
		if err != nil {
			return err
		}

		ctx := cmd.Context()
		user, err := auth.SignInWithEmailPassword(ctx, email, password)
		if err != nil {
			return fmt.Errorf("error authenticating: %v", err)
		}

		token, err := firebaseClient.Auth.CustomToken(ctx, user.UID)
		if err != nil {
			return fmt.Errorf("error creating token: %v", err)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory: %v", err)
		}

		configDir := filepath.Join(homeDir, ".elitecode")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}

		config := struct {
			Token string `json:"token"`
			Email string `json:"email"`
			UID   string `json:"uid"`
		}{
			Token: token,
			Email: email,
			UID:   user.UID,
		}

		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling config: %v", err)
		}

		configFile := filepath.Join(configDir, "config.json")
		if err := os.WriteFile(configFile, configJSON, 0600); err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		fmt.Printf("Successfully logged in as %s\n", email)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
