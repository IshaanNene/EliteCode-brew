package cmd

import (
	"fmt"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up for a new Elitecode account",
	Long: `Create a new Elitecode account. You will be prompted for:
- Name
- Username
- Email
- Password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get user input
		name, err := promptForInput("Name", validateName)
		if err != nil {
			return err
		}

		username, err := promptForInput("Username", validateUsername)
		if err != nil {
			return err
		}

		email, err := promptForInput("Email", validateEmail)
		if err != nil {
			return err
		}

		password, err := promptForPassword("Password", validatePassword)
		if err != nil {
			return err
		}

		_, err = promptForPassword("Confirm Password", func(input string) error {
			if input != password {
				return fmt.Errorf("passwords do not match")
			}
			return nil
		})
		if err != nil {
			return err
		}

		// Create user in Firebase
		ctx := cmd.Context()
		params := (&auth.UserToCreate{}).
			Email(email).
			Password(password).
			DisplayName(name)

		user, err := firebaseClient.Auth.CreateUser(ctx, params)
		if err != nil {
			return fmt.Errorf("error creating user: %v", err)
		}

		// Store additional user data in Firestore
		userDoc := map[string]interface{}{
			"username":    username,
			"email":       email,
			"displayName": name,
			"createdAt":   user.UserMetadata.CreationTimestamp,
		}

		_, err = firebaseClient.Firestore.Collection("users").Doc(user.UID).Set(ctx, userDoc)
		if err != nil {
			// Try to delete the created user if Firestore update fails
			if delErr := firebaseClient.Auth.DeleteUser(ctx, user.UID); delErr != nil {
				return fmt.Errorf("error creating user profile and failed to rollback: %v, %v", err, delErr)
			}
			return fmt.Errorf("error creating user profile: %v", err)
		}

		fmt.Printf("Successfully created account for %s (%s)\n", name, email)
		fmt.Println("You can now log in using the 'elitecode login' command")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(signupCmd)
}

func promptForInput(label string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	return prompt.Run()
}

func promptForPassword(label string, validate func(string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Mask:     '*',
		Validate: validate,
	}
	return prompt.Run()
}

func validateName(input string) error {
	if len(strings.TrimSpace(input)) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	return nil
}

func validateUsername(input string) error {
	if len(input) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if !strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-") {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}
	return nil
}

func validateEmail(input string) error {
	if !strings.Contains(input, "@") || !strings.Contains(input, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func validatePassword(input string) error {
	if len(input) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if !strings.ContainsAny(input, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !strings.ContainsAny(input, "0123456789") {
		return fmt.Errorf("password must contain at least one number")
	}
	if !strings.ContainsAny(input, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}
