package problems

import (
    "context"
    "github.com/IshaanNene/EliteCode-brew/firebase"
    "fmt"
    "github.com/manifoldco/promptui"
    "log"
)

func SetProblem() {
    app := firebase.InitFirebase()
    ctx := context.Background()

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Fatalf("Failed to initialize Firestore client: %v", err)
    }
    defer client.Close()

    iter := client.Collection("problems").Documents(ctx)

    var problems []Problem
    var displayList []string

    for {
        doc, err := iter.Next()
        if err != nil {
            break
        }

        var p Problem
        err = doc.DataTo(&p)
        if err != nil {
            fmt.Println("Skipping a problem due to error:", err)
            continue
        }

        problems = append(problems, p)
        displayList = append(displayList, fmt.Sprintf("%s (%s)", p.Title, p.Difficulty))
    }

    if len(problems) == 0 {
        fmt.Println("No problems found.")
        return
    }

    // Use PromptUI to select problem
    prompt := promptui.Select{
        Label: "Select a problem",
        Items: displayList,
        Size:  10,
    }

    index, _, err := prompt.Run()
    if err != nil {
        fmt.Printf("Prompt failed: %v\n", err)
        return
    }

    selected := problems[index]
    fmt.Println("\nYou selected:\n")
    PrintProblemCard(selected.Title, selected.ID, selected.Difficulty, selected.Tags, selected.LanguagesSupported)
}
