package cmd

import (
	"fmt"
	"github.com/IshaanNene/EliteCode-brew/utils"
  "github.com/IshaanNene/EliteCode-brew/problems"
	"os"
	"path/filepath"
)

func RunSolution() {
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

	err := utils.PullStarterCode()
	if err != nil {
		fmt.Printf("❌ Failed to pull starter code: %v\n", err)
		return
	}

	dockerErr := utils.RunInDocker(lang, problemPath, "testcases.json", false)
	if dockerErr != nil {
		fmt.Printf("❌ Docker execution failed: %v\n", dockerErr)
	}
}
