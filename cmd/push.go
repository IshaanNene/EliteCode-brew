package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/github"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push your solution to GitHub",
	Long: `Push your solution to GitHub.
This will:
1. Create a new branch for the problem
2. Commit your solution
3. Create a pull request (optional)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("error getting working directory: %v", err)
		}

		metadataBytes, err := os.ReadFile(filepath.Join(wd, "metadata.json"))
		if err != nil {
			return fmt.Errorf("error reading metadata file: %v", err)
		}

		var meta struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Language string `json:"language"`
		}
		if err := json.Unmarshal(metadataBytes, &meta); err != nil {
			return fmt.Errorf("error parsing metadata: %v", err)
		}

		token, err := github.GetGitHubToken()
		if err != nil {
			return err
		}

		client, err := github.NewClient(token)
		if err != nil {
			return fmt.Errorf("error creating GitHub client: %v", err)
		}

		owner, repo, err := github.GetRepositoryInfo()
		if err != nil {
			return err
		}

		client.SetRepository(owner, repo)

		sourceCode, err := os.ReadFile(filepath.Join(wd, "main."+utils.GetFileExtension(meta.Language)))
		if err != nil {
			return fmt.Errorf("error reading source code: %v", err)
		}

		testCases, err := os.ReadFile(filepath.Join(wd, "test_cases.json"))
		if err != nil {
			return fmt.Errorf("error reading test cases: %v", err)
		}

		dockerfile, err := os.ReadFile(filepath.Join(wd, "Dockerfile"))
		if err != nil {
			return fmt.Errorf("error reading Dockerfile: %v", err)
		}

		files := map[string][]byte{
			fmt.Sprintf("%s/main.%s", meta.ID, utils.GetFileExtension(meta.Language)): sourceCode,
			fmt.Sprintf("%s/test_cases.json", meta.ID):                                testCases,
			fmt.Sprintf("%s/Dockerfile", meta.ID):                                     dockerfile,
			fmt.Sprintf("%s/metadata.json", meta.ID):                                  metadataBytes,
		}

		ctx := cmd.Context()
		if err := client.CommitSolution(ctx, meta.ID, meta.Title, files); err != nil {
			return fmt.Errorf("error committing solution: %v", err)
		}

		if createPR {
			branchName := fmt.Sprintf("problem/%s", meta.ID)
			title := fmt.Sprintf("Solution for problem %s: %s", meta.ID, meta.Title)
			body := fmt.Sprintf("Language: %s\n\nThis PR contains my solution for problem %s.", meta.Language, meta.Title)

			pr, err := client.CreatePullRequest(ctx, title, body, branchName, "main")
			if err != nil {
				return fmt.Errorf("error creating pull request: %v", err)
			}

			fmt.Printf("Created pull request #%d: %s\n", pr.GetNumber(), pr.GetHTMLURL())
		} else {
			fmt.Printf("Pushed solution to branch problem/%s\n", meta.ID)
		}

		return nil
	},
}

var createPR bool

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVarP(&createPR, "pr", "p", false, "Create a pull request")
}
