package auth

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Create a new Elitecode account",
	RunE:  runSignup,
}

func runSignup(cmd *cobra.Command, args []string) error {
	var answers struct {
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
				Message: "Name:",
			},
			Validate: survey.Required,
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username:",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) < 3 {
					return fmt.Errorf("username must be at least 3 characters")
				}
				return nil
			},
		},
		{
			Name: "email",
			Prompt: &survey.Input{
				Message: "Email:",
			},
			Validate: utils.ValidateEmail,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password:",
			},
			Validate: utils.ValidatePassword,
		},
		{
			Name: "confirmPassword",
			Prompt: &survey.Password{
				Message: "Confirm Password:",
			},
		},
	}

	err := survey.Ask(questions, &answers)
	if err != nil {
		return err
	}

	if answers.Password != answers.ConfirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	client := api.NewClient()
	user, err := client.Signup(answers.Name, answers.Username, answers.Email, answers.Password)
	if err != nil {
		return fmt.Errorf("signup failed: %v", err)
	}

	color.Green("✅ Account created successfully!")
	color.Blue("Welcome, %s!", user.Name)
	
	return nil
}