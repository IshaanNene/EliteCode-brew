package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your account",
	Long: `Log out from your account.
This will:
1. Remove your authentication token
2. Clear your GitHub token (if any)
3. Remove any cached data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory: %v", err)
		}

		configDir := filepath.Join(homeDir, ".elitecode")
		configFile := filepath.Join(configDir, "config.json")

		configBytes, err := os.ReadFile(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Already logged out.")
				return nil
			}
			return fmt.Errorf("error reading config file: %v", err)
		}

		var config map[string]interface{}
		if err := json.Unmarshal(configBytes, &config); err != nil {
			return fmt.Errorf("error parsing config file: %v", err)
		}

		if _, ok := config["uid"]; !ok {
			fmt.Println("Already logged out.")
			return nil
		}

		config = make(map[string]interface{})

		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling config: %v", err)
		}

		if err := os.WriteFile(configFile, configJSON, 0600); err != nil {
			return fmt.Errorf("error writing config file: %v", err)
		}

		cacheDir := filepath.Join(configDir, "cache")
		if err := os.RemoveAll(cacheDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("error removing cache directory: %v", err)
		}

		fmt.Println("Successfully logged out.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
