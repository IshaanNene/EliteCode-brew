package github

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/github"
	"github.com/IshaanNene/EliteCode-brew/internal/storage"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push solution to GitHub",
	RunE:  runPush,
}

func init() {
	pushCmd.Flags().StringP("repo", "r", "", "GitHub repository (format: owner/repo)")
	pushCmd.Flags().StringP("branch", "b", "main", "Branch to push to")
	pushCmd.Flags().StringP("message", "m", "", "Commit message")
	pushCmd.Flags().BoolP("create-repo", "c", false, "Create repository if it doesn't exist")
}

func runPush(cmd *cobra.Command, args []string) error {
	repo, _ := cmd.Flags().GetString("repo")
	branch, _ := cmd.Flags().GetString("branch")
	message, _ := cmd.Flags().GetString("message")
	createRepo, _ := cmd.Flags().GetBool("create-repo")

	// Get problem ID from current directory
	problemID, err := getProblemIDFromDirectory()
	if err != nil {
		return fmt.Errorf("could not determine problem ID: %v", err)
	}

	// Get problem info
	cache := storage.NewCacheManager()
	problem, err := cache.GetProblemMetadata(problemID)
	if err != nil {
		return fmt.Errorf("could not get problem metadata: %v", err)
	}

	// Set default repository name if not provided
	if repo == "" {
		config := storage.NewConfigManager()
		user, err := config.GetUser()
		if err != nil {
			return fmt.Errorf("could not get user info: %v", err)
		}
		repo = fmt.Sprintf("%s/elitecode-solutions", user.Username)
	}

	// Set default commit message if not provided
	if message == "" {
		message = fmt.Sprintf("Add solution for problem %d: %s", problemID, problem.Title)
	}

	// Find source files
	sourceFiles, err := findSourceFiles()
	if err != nil {
		return fmt.Errorf("could not find source files: %v", err)
	}

	if len(sourceFiles) == 0 {
		return fmt.Errorf("no source files found")
	}

	// Get GitHub token
	token, err := getGitHubToken()
	if err != nil {
		return fmt.Errorf("GitHub token required: %v", err)
	}

	spinner := utils.NewSpinner("Pushing to GitHub...")
	spinner.Start()

	// Initialize GitHub client
	client := github.NewClient(token)

	// Create repository if requested
	if createRepo {
		if err := client.CreateRepository(repo); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to create repository: %v", err)
		}
	}

	// Create problem directory structure
	problemDir := fmt.Sprintf("problems/%d_%s", problemID, utils.Slugify(problem.Title))
	
	// Push files
	for _, sourceFile := range sourceFiles {
		content, err := os.ReadFile(sourceFile)
		if err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to read source file %s: %v", sourceFile, err)
		}

		remotePath := filepath.Join(problemDir, sourceFile)
		if err := client.PushFile(repo, remotePath, content, message, branch); err != nil {
			spinner.Stop()
			return fmt.Errorf("failed to push file %s: %v", sourceFile, err)
		}
	}

	// Push README if it doesn't exist
	readmePath := filepath.Join(problemDir, "README.md")
	readme := generateReadme(problem)
	if err := client.PushFile(repo, readmePath, []byte(readme), message, branch); err != nil {
		spinner.Stop()
		color.Yellow("⚠️  Could not push README: %v", err)
	}

	spinner.Stop()
	color.Green("✅ Successfully pushed to GitHub!")
	color.Blue("📁 Repository: https://github.com/%s", repo)
	color.Blue("📝 Files pushed: %s", strings.Join(sourceFiles, ", "))

	return nil
}

func getProblemIDFromDirectory() (int, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	dirName := filepath.Base(pwd)
	parts := strings.Split(dirName, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid directory format")
	}

	return strconv.Atoi(parts[0])