package problems

import (
	"context"
	"encoding/json"
	"fmt"
    "strings"
	"log"
	"os"
    "os/exec"
	"path/filepath"
	"github.com/IshaanNene/EliteCode-brew/firebase"
	"github.com/manifoldco/promptui"
)

type ProblemSelection struct {
	Problem struct {
		ID string `json:"id"`
	} `json:"problem"`
	Lang string `json:"lang"`
}


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

	// Select a problem
	problemPrompt := promptui.Select{
		Label: "Select a problem",
		Items: displayList,
		Size:  10,
	}

	problemIndex, _, err := problemPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		return
	}

	selected := problems[problemIndex]
	fmt.Println("\nYou selected:\n")
	PrintProblemCard(selected.Title, selected.ID, selected.Difficulty, selected.Tags, selected.LanguagesSupported)

	// Select a language
	langPrompt := promptui.Select{
		Label: "Select a language",
		Items: selected.LanguagesSupported,
	}

	langIndex, _, err := langPrompt.Run()
	if err != nil {
		fmt.Printf("Language prompt failed: %v\n", err)
		return
	}

	selectedLang := selected.LanguagesSupported[langIndex]

	err = saveSelectedProblem(selected, selectedLang)
	if err != nil {
		fmt.Printf("Failed to save selection: %v\n", err)
		return
	}

	fmt.Printf("Saved selection: %s in %s\n", selected.ID, selectedLang)

	statePath := filepath.Join(os.Getenv("HOME"), ".elitecode", "selected.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		fmt.Printf("Failed to read selected problem: %v\n", err)
		return
	}

	var selection ProblemSelection
	if err := json.Unmarshal(data, &selection); err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
		return
	}

	language := strings.ToLower(selection.Lang)
	if language != "cpp" && language != "python" && language != "java" && language != "go" {
		fmt.Printf("Unsupported language: %s\n", language)
		return
	}


	// ✅ FIX: correct folder name
	folder := "Problems" + strings.TrimPrefix(selected.ID, "problem")

	// Step 2: Setup script
	scriptPath := "scripts/setup_dir.sh"
	err = os.Chmod(scriptPath, 0755)
	if err != nil {
		fmt.Printf("Failed to make script executable: %v\n", err)
		return
	}

	// Step 3: Call script with corrected folder name and language
	cmd := exec.Command("bash", scriptPath, folder, language)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Error running Setup_Dir.sh: %v\n", err)
		return
	}
}


func saveSelectedProblem(p Problem, lang string) error {
	stateDir := filepath.Join(os.Getenv("HOME"), ".elitecode")
	_ = os.MkdirAll(stateDir, 0755)

	statePath := filepath.Join(stateDir, "selected.json")

	data := struct {
		Problem Problem `json:"problem"`
		Lang    string  `json:"lang"`
	}{
		Problem: p,
		Lang:    lang,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, jsonBytes, 0644)
}
