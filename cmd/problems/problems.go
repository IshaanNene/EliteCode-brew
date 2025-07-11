package problems

import (
	"github.com/spf13/cobra"
)

var ProblemsCmd = &cobra.Command{
	Use:   "problems",
	Short: "Problem management commands",
	Long:  "Commands for browsing, solving, and managing coding problems",
}

func init() {
	ProblemsCmd.AddCommand(listCmd)
	ProblemsCmd.AddCommand(setCmd)
	ProblemsCmd.AddCommand(runCmd)
	ProblemsCmd.AddCommand(submitCmd)
	ProblemsCmd.AddCommand(resetCmd)
	ProblemsCmd.AddCommand(searchCmd)
	ProblemsCmd.AddCommand(bookmarkCmd)
	ProblemsCmd.AddCommand(discussCmd)
}