package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var HelpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help for EliteCode CLI",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Root().Help()
		if err != nil {
			log.Fatal(err)
		}
	},
}
