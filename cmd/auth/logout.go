package auth

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your Elitecode account",
	RunE:  runLogout,
}

func runLogout(cmd *cobra.Command, args []string) error {
	config := storage.NewConfigManager()
	
	if err := config.ClearAuth(); err != nil {
		return err
	}

	color.Green("✅ Logged out successfully!")
	return nil
}