package auth

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new Elitecode account",
	RunE:  runSignup,
}

func runSignup(cmd *cobra.Command, args []string) error {
	var credentials struct {
		Name            string
		Username        string
		Email           string
		Password        string
		ConfirmPassword string
	}

	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Full Name:",
			},
			Validate: survey.Required,
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username:",
			},
			Validate: survey.Required,
		},
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
			},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password:",
			},
			Validate: survey.Required,
		},
		{
			Name: "confirmpassword",
			Prompt: &survey.Password{
				Message: "Confirm Password:",
			},
			Validate: survey.Required,
		},
	}

	err := survey.Ask(questions, &credentials)
	if err != nil {
		return err
	}

	// Validate passwords match
	if credentials.Password != credentials.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Validate email format
	if !strings.Contains(credentials.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	client := api.NewClient()
	user, token, err := client.Signup(credentials.Name, credentials.Username, credentials.Email, credentials.Password)
	if err != nil {
		return fmt.Errorf("signup failed: %v", err)
	}

	// Store auth token and user info
	config := storage.NewConfigManager()
	if err := config.SetAuthToken(token); err != nil {
		return fmt.Errorf("failed to store auth token: %v", err)
	}

	if err := config.SetUser(user); err != nil {
		return fmt.Errorf("failed to store user info: %v", err)
	}

	color.Green("✅ Account created successfully!")
	color.Blue("Welcome to Elitecode, %s!", user.Name)
	
	return nil
}