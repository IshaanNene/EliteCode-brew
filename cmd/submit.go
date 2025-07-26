package cmd

import (
	"fmt"
	"github.com/IshaanNene/EliteCode-brew/utils"
	"os"
	"path/filepath"
)

func SubmitSolution() {
	problemPath := filepath.Join(".", utils.GetProblemDirectory())
	if _, err := os.Stat(problemPath); os.IsNotExist(err) {
		fmt.Println("❌ Problem directory not found. Please run `elitecode set_problem` first.")
		return
	}

	lang := utils.GetSelectedLanguage()
	if lang == "" {
		fmt.Println("❌ Selected language not found. Try setting the problem again.")
		return
	}

	fmt.Println("🚀 Submitting your solution...")
	err := utils.RunInDocker(lang, problemPath, "testcases.json", true)
	if err != nil {
		fmt.Printf("❌ Submission failed: %v\n", err)
	}
}
