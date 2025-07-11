package main

import (
	"os"

	"github.com/IshaanNene/EliteCode-brew/cmd/auth"
	"github.com/IshaanNene/EliteCode-brew/cmd/github"
	"github.com/IshaanNene/EliteCode-brew/cmd/problems"
	"github.com/IshaanNene/EliteCode-brew/cmd/system"
	"github.com/IshaanNene/EliteCode-brew/cmd/user"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "elitecode",
		Short: "A CLI tool for competitive programming and coding practice",
		Long: `Elitecode is a comprehensive CLI tool that helps you practice coding problems,
manage your solutions, and track your progress. It integrates with GitHub,
Docker for sandboxed execution, and provides a complete development environment.`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize logging
			utils.InitLogger()
		},
	}

	// Add version info
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
Commit: {{.Commit}}
Date: {{.Date}}
`)

	// Add all command groups
	rootCmd.AddCommand(auth.AuthCmd)
	rootCmd.AddCommand(problems.ProblemsCmd)
	rootCmd.AddCommand(user.UserCmd)
	rootCmd.AddCommand(system.SystemCmd)
	rootCmd.AddCommand(github.GitHubCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		utils.Logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}