package cmd

import (
    "github.com/IshaanNene/EliteCode-brew/problems"
    "github.com/IshaanNene/EliteCode-brew/utils"
    "fmt"
    "log"
    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
    Use:   "run",
    Short: "Select a language and fetch problem code from GitHub",
    Run: func(cmd *cobra.Command, args []string) {
        selectedProblem := problems.GetSelectedProblem()
        if selectedProblem.ID == "" {
            fmt.Println("No problem selected. Use 'elitecode set_problem' first.")
            return
        }

        prompt := promptui.Select{
            Label: "Select language",
            Items: selectedProblem.LanguagesSupported,
        }

        _, lang, err := prompt.Run()
        if err != nil {
            fmt.Println("Cancelled.")
            return
        }

        fmt.Printf("Language selected: %s\n", lang)

        folder := fmt.Sprintf("Solutions/%s", selectedProblem.ID)
        target := fmt.Sprintf("%s", selectedProblem.ID)

        err = utils.FetchFilesFromGitHub(
            "IshaanNene",
            "AlgoRank",
            "main",
            folder,
            target,
        )
        if err != nil {
            log.Fatalf("Error fetching files: %v", err)
        }

        fmt.Printf("Problem '%s' fetched for language '%s'.\n", selectedProblem.Title, lang)
    },
}
