package system

import (
	"github.com/spf13/cobra"
)

var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "System management commands",
	Long:  "Commands for system initialization and maintenance",
}

func init() {
	SystemCmd.AddCommand(initCmd)
	SystemCmd.AddCommand(updateCmd)
}