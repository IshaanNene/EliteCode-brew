package cmd

import (
    "github.com/IshaanNene/EliteCode-brew/problems"
    "fmt"
    "github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
    Use:   "submit",
    Short: "Submit current solution",
    Run: func(cmd *cobra.Command, args []string) {
        selectedProblem := problems.GetSelectedProblem()
        if selectedProblem.ID == "" {
            fmt.Println("No problem selected.")
            return
        }
        fmt.Printf("Submitted '%s' successfully!\n", selectedProblem.Title)
    },
}
