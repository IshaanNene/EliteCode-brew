package cmd

import (
    "elitecode/problems"
    "fmt"

    "github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
    Use:   "submit",
    Short: "Submit current solution",
    Run: func(cmd *cobra.Command, args []string) {
        selectedProblem := problems.GetSelectedProblem()
        if selectedProblem.ID == "" {
            fmt.Println("❌ No problem selected.")
            return
        }

        // You can later push to remote / Firebase etc.
        fmt.Printf("📤 Submitted '%s' successfully!\n", selectedProblem.Title)
    },
}
