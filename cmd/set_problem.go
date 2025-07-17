package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/elitecode/internal/problem"
	"github.com/yourusername/elitecode/internal/utils"
)

var (
	language string
)

var setProblemCmd = &cobra.Command{
	Use:   "set_problem [problem_id]",
	Short: "Set up a problem for solving",
	Long: `Set up a coding problem for solving. This will:
1. Download problem files from Firebase Storage
2. Create a local directory structure
3. Set up Docker environment for isolated execution`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		problemID := args[0]

		// Create problem service
		problemService := problem.NewService(firebaseClient.Firestore)

		// Get problem details
		ctx := cmd.Context()
		prob, err := problemService.GetProblem(ctx, problemID)
		if err != nil {
			return fmt.Errorf("error getting problem: %v", err)
		}

		// Validate language
		if language == "" {
			// Use first supported language as default
			if len(prob.SupportedLangs) > 0 {
				language = prob.SupportedLangs[0]
			} else {
				return fmt.Errorf("no supported languages for this problem")
			}
		} else {
			supported := false
			for _, lang := range prob.SupportedLangs {
				if lang == language {
					supported = true
					break
				}
			}
			if !supported {
				return fmt.Errorf("language %s is not supported for this problem. Supported languages: %s",
					language, strings.Join(prob.SupportedLangs, ", "))
			}
		}

		// Create problem directory
		dirName := fmt.Sprintf("%s_%s", problemID, strings.ReplaceAll(strings.ToLower(prob.Title), " ", "_"))
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return fmt.Errorf("error creating problem directory: %v", err)
		}

		// Download problem files from Firebase Storage
		bucket := firebaseClient.Storage.Bucket(prob.StoragePaths.StarterCode)

		// Download starter code
		starterCodePath := filepath.Join(prob.StoragePaths.StarterCode, language, "main."+utils.GetFileExtension(language))
		starterCodeObj := bucket.Object(starterCodePath)
		starterCodeReader, err := starterCodeObj.NewReader(ctx)
		if err != nil {
			return fmt.Errorf("error reading starter code: %v", err)
		}
		defer starterCodeReader.Close()

		starterCodeFile := filepath.Join(dirName, "main."+utils.GetFileExtension(language))
		starterCodeWriter, err := os.Create(starterCodeFile)
		if err != nil {
			return fmt.Errorf("error creating starter code file: %v", err)
		}
		defer starterCodeWriter.Close()

		if _, err := io.Copy(starterCodeWriter, starterCodeReader); err != nil {
			return fmt.Errorf("error writing starter code: %v", err)
		}

		// Download test cases
		testCasesObj := bucket.Object(prob.StoragePaths.TestCases)
		testCasesReader, err := testCasesObj.NewReader(ctx)
		if err != nil {
			return fmt.Errorf("error reading test cases: %v", err)
		}
		defer testCasesReader.Close()

		testCasesFile := filepath.Join(dirName, "test_cases.json")
		testCasesWriter, err := os.Create(testCasesFile)
		if err != nil {
			return fmt.Errorf("error creating test cases file: %v", err)
		}
		defer testCasesWriter.Close()

		if _, err := io.Copy(testCasesWriter, testCasesReader); err != nil {
			return fmt.Errorf("error writing test cases: %v", err)
		}

		// Create metadata file
		metadataFile := filepath.Join(dirName, "metadata.json")
		if err := os.WriteFile(metadataFile, []byte(fmt.Sprintf(`{
			"id": "%s",
			"title": "%s",
			"language": "%s",
			"time_limit": %d,
			"memory_limit": %d
		}`, prob.ID, prob.Title, language, prob.TimeLimit, prob.MemoryLimit)), 0644); err != nil {
			return fmt.Errorf("error writing metadata file: %v", err)
		}

		// Download Dockerfile
		dockerfileObj := bucket.Object(filepath.Join("docker_templates", language+".dockerfile"))
		dockerfileReader, err := dockerfileObj.NewReader(ctx)
		if err != nil {
			return fmt.Errorf("error reading Dockerfile template: %v", err)
		}
		defer dockerfileReader.Close()

		dockerfileWriter, err := os.Create(filepath.Join(dirName, "Dockerfile"))
		if err != nil {
			return fmt.Errorf("error creating Dockerfile: %v", err)
		}
		defer dockerfileWriter.Close()

		if _, err := io.Copy(dockerfileWriter, dockerfileReader); err != nil {
			return fmt.Errorf("error writing Dockerfile: %v", err)
		}

		fmt.Printf("Successfully set up problem %s in directory %s\n", prob.Title, dirName)
		fmt.Printf("You can now start coding in %s\n", starterCodeFile)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(setProblemCmd)

	// Add flags
	setProblemCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language to use")
}
