package github

import (
	"github.com/spf13/cobra"
)

var GithubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub integration commands",
	Long:  "Commands for GitHub repository integration",
}

func init() {
	GithubCmd.AddCommand(pushCmd)
}