package problems

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the current problem to initial state",
	RunE:  runReset,
}

func init() {
	resetCmd.Flags().BoolP("force", "f", false, "Force reset without confirmation")
}

func runReset(cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	// Find source files
	sourceFiles := []string{}
	extensions := []string{".c", ".cpp", ".py", ".java", ".js"}
	
	for _, ext := range extensions {
		matches, err := filepath.Glob("*" + ext)
		if err != nil {
			continue
		}
		sourceFiles = append(sourceFiles, matches...)
	}

	if len(sourceFiles) == 0 {
		color.Yellow("No source files found to reset")
		return nil
	}

	// Confirm reset
	if !force {
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to reset %d source file(s)?", len(sourceFiles)),
		}
		survey.AskOne(prompt, &confirm)
		if !confirm {
			color.Yellow("Reset cancelled")
			return nil
		}
	}

	// Reset each source file
	for _, sourceFile := range sourceFiles {
		if err := resetSourceFile(sourceFile); err != nil {
			color.Red("❌ Failed to reset %s: %v", sourceFile, err)
			continue
		}
		color.Green("✅ Reset %s", sourceFile)
	}

	color.Green("🔄 Reset completed!")
	return nil
}

func resetSourceFile(filename string) error {
	ext := filepath.Ext(filename)
	language := getLanguageFromExtension(ext)
	
	// Create basic template based on language
	var template string
	switch language {
	case "c":
		template = `#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main() {
    // Your code here
    
    return 0;
}`
	case "cpp":
		template = `#include <iostream>
#include <vector>
#include <string>
#include <algorithm>

using namespace std;

int main() {
    // Your code here
    
    return 0;
}`
	case "python":
		template = `def main():
    # Your code here
    pass

if __name__ == "__main__":
    main()`
	case "java":
		template = `import java.util.*;
import java.io.*;

public class Main {
    public static void main(String[] args) {
        // Your code here
    }
}`
	case "javascript":
		template = `function main() {
    // Your code here
}

main();`
	default:
		return fmt.Errorf("unsupported language")
	}

	return os.WriteFile(filename, []byte(template), 0644)
}

func getLanguageFromExtension(ext string) string {
	switch ext {
	case ".c":
		return "c"
	case ".cpp":
		return "cpp"
	case ".py":
		return "python"
	case ".java":
		return "java"
	case ".js":
		return "javascript"
	default:
		return ""
	}
}