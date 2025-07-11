package user

import (
	"github.com/spf13/cobra"
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "User management commands",
	Long:  "Commands for managing user profile and statistics",
}

func init() {
	UserCmd.AddCommand(statsCmd)
	UserCmd.AddCommand(myProblemsCmd)
}