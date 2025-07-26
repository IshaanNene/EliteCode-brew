package problems

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/IshaanNene/EliteCode-brew/firebase"
	"github.com/IshaanNene/EliteCode-brew/utils"
	"github.com/manifoldco/promptui"
)

func SetProblem() {
	app := firebase.InitFirebase()
	ctx := context.Background()
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Firestore init error: %v", err)
	}
	defer client.Close()

	iter := client.Collection("problems").Documents(ctx)
	var list []Problem
	var items []string
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var p Problem
		doc.DataTo(&p)
		list = append(list, p)
		items = append(items, fmt.Sprintf("%s (%s)", p.Title, p.Difficulty))
	}

	prompt := promptui.Select{
		Label: "Select Problem",
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}
	sel := list[idx]

	prompt2 := promptui.Select{
		Label: "Select Language",
		Items: sel.LanguagesSupported,
	}
	_, lang, err := prompt2.Run()
	if err != nil {
		fmt.Println("Cancelled")
		return
	}

	// Save to local state
	sel.Lang = lang
	savePath := filepath.Join(os.Getenv("HOME"), ".elitecode", "selected.json")
	os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	data, _ := json.MarshalIndent(sel, "", "  ")
	os.WriteFile(savePath, data, 0644)

	fmt.Printf("✅ Saved: %s (%s) in %s\n", sel.Title, sel.ID, lang)

	// Fetch code from GitHub
	utils.FetchStarterCode(sel.ID)
}
