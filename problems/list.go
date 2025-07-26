package problems

import (
    "context"
    "fmt"
    "log"
	"strings"
    "github.com/IshaanNene/EliteCode-brew/firebase"
    _"cloud.google.com/go/firestore"
)

type Problem struct {
    ID                string   `firestore:"id"`
    Title             string   `firestore:"title"`
    Difficulty        string   `firestore:"difficulty"`
    Tags              []string `firestore:"tags"`
    LanguagesSupported []string `firestore:"languages_supported"`
}

func PrintProblemCard(title, id, difficulty string, tags, languages []string) {
    lines := []string{
        fmt.Sprintf("Title     : %s", title),
        fmt.Sprintf("ID        : %s", id),
        fmt.Sprintf("Difficulty: %s", difficulty),
        fmt.Sprintf("Tags      : %s", strings.Join(tags, ", ")),
        fmt.Sprintf("Languages : %s", strings.Join(languages, ", ")),
    }
    maxWidth := 0
    for _, line := range lines {
        if len(line) > maxWidth {
            maxWidth = len(line)
        }
    }
    border := "┌" + strings.Repeat("─", maxWidth+2) + "┐"
    footer := "└" + strings.Repeat("─", maxWidth+2) + "┘"

    fmt.Println(border)
    for _, line := range lines {
        padding := strings.Repeat(" ", maxWidth-len(line))
        fmt.Printf("│ %s%s │\n", line, padding)
    }
    fmt.Println(footer)
    fmt.Println()
}

func ListProblems() {
    app := firebase.InitFirebase()
    ctx := context.Background()

    client, err := app.Firestore(ctx)
    if err != nil {
        log.Fatalf("Failed to get Firestore client: %v", err)
    }
    defer client.Close()

    iter := client.Collection("problems").Documents(ctx)
    fmt.Println("Available Problems:\n")

    for {
        doc, err := iter.Next()
        if err != nil {
            break
        }

        var p Problem
        err = doc.DataTo(&p)
        if err != nil {
            fmt.Println("Error decoding problem:", err)
            continue
        }

		PrintProblemCard(p.Title,p.ID,p.Difficulty,p.Tags,p.LanguagesSupported)
    }
}
