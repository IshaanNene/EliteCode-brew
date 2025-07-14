package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/IshaanNene/EliteCode-brew/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "elitecode",
		Short: "EliteCode - Official CLI for Competitive Coding",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to EliteCode CLI! Use 'elitecode help' to see available commands.")
		},
	}

	rootCmd.AddCommand(
		cmd.SignupCmd,
		cmd.LoginCmd,
		cmd.LogoutCmd,
		cmd.WhoamiCmd,
		cmd.HelpCmd,
		// Add remaining command bindings here
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
