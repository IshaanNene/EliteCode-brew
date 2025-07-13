package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/IshaanNene/EliteCode-brew/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "elitecode",
		Short: "EliteCode - Official CLI for Competitive Coding",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to EliteCode CLI! Use help to see available commands.")
		},
	}

	rootCmd.AddCommand(
		cmd.SignupCmd,
		cmd.LoginCmd,
		cmd.SetProblemCmd,
		cmd.RunCmd,
		cmd.SubmitCmd,
		cmd.PushCmd,
		cmd.ResetCmd,
		cmd.LogoutCmd,
		cmd.HelpCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
