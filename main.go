package main

import (
	"fmt"
	"os"

	"github.com/IshaanNene/EliteCode-brew/auth"
	"github.com/IshaanNene/EliteCode-brew/problems"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "elitecode",
		Short: "EliteCode CLI - Practice DSA problems from your terminal!",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "signup",
		Short: "Sign up for EliteCode",
		Run: func(cmd *cobra.Command, args []string) {
			auth.Signup()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Login to EliteCode",
		Run: func(cmd *cobra.Command, args []string) {
			auth.Login()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "problems",
		Short: "List available problems",
		Run: func(cmd *cobra.Command, args []string) {
			problems.ListProblems()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "set_problem",
		Short: "Select and set a problem and language",
		Run: func(cmd *cobra.Command, args []string) {
			problems.SetProblem()
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "run",
		Short: "Run your solution against test cases",
		Run: func(cmd *cobra.Command, args []string) {
			cmdArgs := []string{}
			if len(args) > 0 {
				cmdArgs = args
			}
			cmd.RunCode(cmdArgs)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "submit",
		Short: "Submit your solution",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.SubmitCode()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
