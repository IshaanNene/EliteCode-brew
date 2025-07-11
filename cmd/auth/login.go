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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your Elitecode account",
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	var credentials struct {
		Username string
		Password string
	}

	questions := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username or Email:",
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
	}

	err := survey.Ask(questions, &credentials)
	if err != nil {
		return err
	}

	// Determine if input is email or username
	isEmail := strings.Contains(credentials.Username, "@")
	
	client := api.NewClient()
	var user *api.User
	var token string

	if isEmail {
		user, token, err = client.LoginWithEmail(credentials.Username, credentials.Password)
	} else {
		user, token, err = client.LoginWithUsername(credentials.Username, credentials.Password)
	}

	if err != nil {
		return fmt.Errorf("login failed: %v", err)
	}

	// Store auth token and user info
	config := storage.NewConfigManager()
	if err := config.SetAuthToken(token); err != nil {
		return fmt.Errorf("failed to store auth token: %v", err)
	}

	if err := config.SetUser(user); err != nil {
		return fmt.Errorf("failed to store user info: %v", err)
	}

	color.Green("✅ Login successful!")
	color.Blue("Welcome back, %s!", user.Name)
	
	return nil
}