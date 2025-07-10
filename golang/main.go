package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/elitecode/cmd/auth"
	"github.com/yourusername/elitecode/cmd/problems"
	"github.com/yourusername/elitecode/cmd/user"
	"github.com/yourusername/elitecode/cmd/system"
	"github.com/yourusername/elitecode/cmd/github"
)

var rootCmd = &cobra.Command{
	Use:   "elitecode",
	Short: "CLI tool for competitive programming and coding challenges",
	Long:  "Elitecode is a comprehensive CLI tool for solving coding challenges, managing submissions, and tracking progress.",
}

func main() {
	// Add command groups
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(problems.ProblemsCmd)
	rootCmd.AddCommand(user.UserCmd)
	rootCmd.AddCommand(system.SystemCmd)
	rootCmd.AddCommand(github.GithubCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}