package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/yourusername/elitecode/internal/github"
)

var githubLoginCmd = &cobra.Command{
	Use:   "github login",
	Short: "Log in to GitHub",
	Long: `Log in to GitHub using a personal access token.
You can create a token at https://github.com/settings/tokens
The token needs the following scopes:
- repo (Full control of private repositories)
- workflow (Update GitHub Action workflows)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prompt for GitHub token
		prompt := promptui.Prompt{
			Label: "GitHub Token",
			Mask:  '*',
			Validate: func(input string) error {
				if len(input) < 10 {
					return fmt.Errorf("token must be at least 10 characters")
				}
				return nil
			},
		}

		token, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("error getting token: %v", err)
		}

		// Validate token by creating a client
		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("error creating GitHub client: %v", err)
		}

		// Get user information to validate token
		ctx := cmd.Context()
		user, err := client.GetUser(ctx)
		if err != nil {
			return fmt.Errorf("error validating token: %v", err)
		}

		// Get config directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory: %v", err)
		}

		configDir := filepath.Join(homeDir, ".elitecode")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("error creating config directory: %v", err)
		}

		// Read existing config
		configFile := filepath.Join(configDir, "config.json")
		var config map[string]interface{}

		configBytes, err := os.ReadFile(configFile)
		if err == nil {
			if err := json.Unmarshal(configBytes, &config); err != nil {
				return fmt.Errorf("error parsing config file: %v", err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("error reading config file: %v", err)
		}

		if config == nil {
			config = make(map[string]interface{})
		}

		// Update config with GitHub token
		config["github_token"] = token
		config["github_username"] = user.GetLogin()

		// Save config
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling config: %v", err)
		}

		if err := os.WriteFile(configFile, configJSON, 0600); err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		fmt.Printf("Successfully logged in as %s\n", user.GetLogin())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(githubLoginCmd)
}
