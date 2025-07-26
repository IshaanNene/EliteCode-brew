package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ProblemSelection struct {
	Problem struct {
		ID string `json:"id"`
	} `json:"problem"`
	Lang string `json:"lang"`
}

func Run() {
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

	problemID := selection.Problem.ID
	language := strings.ToLower(selection.Lang)
	var ext string

	switch language {
	case "cpp":
		ext = "cpp"
	case "python":
		ext = "py"
	case "java":
		ext = "java"
	case "go":
		ext = "go"
	default:
		fmt.Printf("Unsupported language: %s\n", language)
		return
	}

	scriptPath := "scripts/setup_doc_json.sh"
	err = os.Chmod(scriptPath, 0755)
	if err != nil {
		fmt.Printf("Failed to make script executable: %v\n", err)
		return
	}

	folderPath := "Problems" + strings.TrimPrefix(problemID, "problem")
    cmd := exec.Command(scriptPath, folderPath, ext)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		fmt.Printf("Error running Setup_Dir.sh: %v\n", err)
		return
	}

	dockerfile := "Problem/Dockerfile"
	testFile := "Problem/testcases.json"
	fmt.Println("Building Docker container...")
	build := exec.Command("docker", "build", "-f", dockerfile, "-t", "elitecode-runner", "Problem")
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		fmt.Println("Failed to build Docker image.")
		return
	}

	fmt.Println("Running Docker container...")
	run := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", filepath.Join(".", "Problem")),
		"-w", "/app",
		"elitecode-runner",
	)
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	if err := run.Run(); err != nil {
		fmt.Println("Failed to run Docker container.")
	}

	_ = os.Remove(dockerfile)
	_ = os.Remove(testFile)

	fmt.Println("Cleanup done.")
}
