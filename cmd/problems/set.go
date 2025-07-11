package problems

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var setCmd = &cobra.Command{
	Use:   "set <problem_id>",
	Short: "Set up a problem for solving",
	Args:  cobra.ExactArgs(1),
	RunE:  runSet,
}

func init() {
	setCmd.Flags().StringP("language", "l", "c", "Programming language (c, cpp, python, java, javascript)")
	setCmd.Flags().BoolP("open", "o", true, "Open in VS Code after setup")
}

func runSet(cmd *cobra.Command, args []string) error {
	problemID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid problem ID: %v", err)
	}

	language, _ := cmd.Flags().GetString("language")
	openInVSCode, _ := cmd.Flags().GetBool("open")

	client := api.NewClient()
	problem, err := client.GetProblem(problemID)
	if err != nil {
		return fmt.Errorf("failed to fetch problem: %v", err)
	}

	// Create problem directory
	problemDir := fmt.Sprintf("%d_%s", problem.ID, utils.Slugify(problem.Title))
	if err := os.MkdirAll(problemDir, 0755); err != nil {
		return fmt.Errorf("failed to create problem directory: %v", err)
	}

	cache := storage.NewCacheManager()
	
	// Download and save problem files
	spinner := utils.NewSpinner("Setting up problem...")
	spinner.Start()

	// Save problem metadata
	if err := cache.SaveProblemMetadata(problem); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to save problem metadata: %v", err)
	}

	// Download test cases
	testCases, err := client.GetTestCases(problemID)
	if err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to fetch test cases: %v", err)
	}

	if err := cache.SaveTestCases(problemID, testCases); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to save test cases: %v", err)
	}

	// Create source file template
	sourceFile := getSourceFileName(language)
	template := getCodeTemplate(language, problem)
	
	sourceFilePath := filepath.Join(problemDir, sourceFile)
	if err := os.WriteFile(sourceFilePath, []byte(template), 0644); err != nil {
		spinner.Stop()
		return fmt.Errorf("failed to create source file: %v", err)
	}

	// Create input and expected output files
	inputFile := filepath.Join(problemDir, "input.txt")
	expectedFile := filepath.Join(problemDir, "expected_output.txt")
	
	if len(testCases) > 0 {
		if err := os.WriteFile(inputFile, []byte(testCases[0].Input), 0644); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create input file: %v", err)
		}
		
		if err := os.WriteFile(expectedFile, []byte(testCases[0].ExpectedOutput), 0644); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create expected output file: %v", err)
		}
	}

	spinner.Stop()
	color.Green("✅ Problem %d set up successfully!", problemID)
	color.Blue("📁 Directory: %s", problemDir)
	color.Blue("📝 Source file: %s", sourceFile)

	// Open in VS Code if requested
	if openInVSCode {
		if err := openInEditor(problemDir); err != nil {
			color.Yellow("⚠️  Could not open in VS Code: %v", err)
		} else {
			color.Green("🚀 Opened in VS Code")
		}
	}

	return nil
}

func getSourceFileName(language string) string {
	switch language {
	case "c":
		return "main.c"
	case "cpp":
		return "main.cpp"
	case "python":
		return "main.py"
	case "java":
		return "Main.java"
	case "javascript":
		return "main.js"
	default:
		return "main.c"
	}
}

func getCodeTemplate(language string, problem *api.Problem) string {
	switch language {
	case "c":
		return fmt.Sprintf(`#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*
 * Problem: %s
 * Difficulty: %s
 * Category: %s
 * 
 * %s
 */

int main() {
    // Your code here
    
    return 0;
}`, problem.Title, problem.Difficulty, problem.Category, problem.Description)
	case "cpp":
		return fmt.Sprintf(`#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

/*
 * Problem: %s
 * Difficulty: %s
 * Category: %s
 * 
 * %s
 */

int main() {
    // Your code here
    
    return 0;
}`, problem.Title, problem.Difficulty, problem.Category, problem.Description)
	case "python":
		return fmt.Sprintf(`"""
Problem: %s
Difficulty: %s
Category: %s

%s
"""

def main():
    # Your code here
    pass

if __name__ == "__main__":
    main()`, problem.Title, problem.Difficulty, problem.Category, problem.Description)
	case "java":
		return fmt.Sprintf(`import java.util.*;
import java.io.*;

/*
 * Problem: %s
 * Difficulty: %s
 * Category: %s
 * 
 * %s
 */

public class Main {
    public static void main(String[] args) {
        // Your code here
    }
}`, problem.Title, problem.Difficulty, problem.Category, problem.Description)
	case "javascript":
		return fmt.Sprintf(`/*
 * Problem: %s
 * Difficulty: %s
 * Category: %s
 * 
 * %s
 */

function main() {
    // Your code here
}

main();`, problem.Title, problem.Difficulty, problem.Category, problem.Description)
	default:
		return ""
	}
}

func openInEditor(dir string) error {
	// Try VS Code first
	cmd := exec.Command("code", dir)
	if err := cmd.Run(); err != nil {
		// Try other editors
		cmd = exec.Command("subl", dir)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("vim", dir)
			return cmd.Run()
		}
	}
	return nil
}