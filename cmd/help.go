package cmd

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Get detailed help for a command",
	Long: `Get detailed help and examples for a command.
Without arguments, shows general usage information.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Printf("\n%s\n", color.CyanString("Elitecode CLI - Your Coding Practice Companion"))
			fmt.Println(strings.Repeat("-", 50))
			fmt.Println("\nAvailable Commands:")

			categories := map[string][]string{
				"Authentication": {
					"login      - Log in to your account",
					"logout     - Log out from your account",
					"github login - Log in with GitHub",
				},
				"Problem Management": {
					"problems   - List available problems",
					"search    - Search for problems",
					"set_problem - Set up a problem for solving",
					"bookmark  - Manage bookmarked problems",
				},
				"Solution Management": {
					"submit    - Submit your solution",
					"run      - Run your solution locally",
					"push     - Push solution to GitHub",
				},
				"Statistics": {
					"stats    - View problem statistics",
					"profile  - View your profile",
					"leaderboard - View rankings",
				},
			}

			for category, commands := range categories {
				fmt.Printf("\n%s:\n", color.YellowString(category))
				for _, cmd := range commands {
					fmt.Printf("  %s\n", cmd)
				}
			}

			fmt.Printf("\nRun '%s' for detailed help.\n", color.GreenString("elitecode help <command>"))
			return nil
		}

		commandName := args[0]
		var examples string
		var usage string

		switch commandName {
		case "login":
			usage = "elitecode login"
			examples = `
Examples:
  # Log in with email and password
  $ elitecode login

  # Log in with GitHub
  $ elitecode github login`

		case "problems":
			usage = "elitecode problems [flags]"
			examples = `
Examples:
  # List all problems
  $ elitecode problems

  # List easy problems
  $ elitecode problems -d easy

  # List problems with specific tags
  $ elitecode problems -t "arrays,sorting"

  # List solved problems
  $ elitecode problems -s solved`

		case "search":
			usage = "elitecode search [query] [flags]"
			examples = `
Examples:
  # Search by text
  $ elitecode search "binary tree"

  # Search by difficulty
  $ elitecode search -d medium

  # Search by tags
  $ elitecode search -t "dynamic programming"

  # Search by status
  $ elitecode search -s attempted`

		case "set_problem":
			usage = "elitecode set_problem <problem_id> [flags]"
			examples = `
Examples:
  # Set up a problem
  $ elitecode set_problem two-sum

  # Set up with specific language
  $ elitecode set_problem two-sum -l python

  # Set up with test cases
  $ elitecode set_problem two-sum --with-tests`

		case "submit":
			usage = "elitecode submit"
			examples = `
Examples:
  # Submit current solution
  $ elitecode submit

  # Submit and create PR
  $ elitecode submit --pr`

		case "run":
			usage = "elitecode run [flags]"
			examples = `
Examples:
  # Run with sample test cases
  $ elitecode run

  # Run with all test cases
  $ elitecode run --all

  # Run with specific test case
  $ elitecode run -t 1`

		case "push":
			usage = "elitecode push [flags]"
			examples = `
Examples:
  # Push solution to GitHub
  $ elitecode push

  # Push and create PR
  $ elitecode push --pr`

		case "stats":
			usage = "elitecode stats <problem_id>"
			examples = `
Examples:
  # View problem statistics
  $ elitecode stats two-sum

  # View with recent submissions
  $ elitecode stats two-sum --recent`

		case "profile":
			usage = "elitecode profile"
			examples = `
Examples:
  # View your profile
  $ elitecode profile`

		case "leaderboard":
			usage = "elitecode leaderboard [flags]"
			examples = `
Examples:
  # View global leaderboard
  $ elitecode leaderboard

  # View problem-specific leaderboard
  $ elitecode leaderboard -p two-sum

  # View monthly rankings
  $ elitecode leaderboard -t month`

		case "bookmark":
			usage = "elitecode bookmark [problem_id]"
			examples = `
Examples:
  # List bookmarked problems
  $ elitecode bookmark

  # Add/remove bookmark
  $ elitecode bookmark two-sum`

		default:
			return fmt.Errorf("unknown command: %s", commandName)
		}

		fmt.Printf("\n%s\n", color.CyanString("Command: "+commandName))
		fmt.Println(strings.Repeat("-", 50))
		fmt.Printf("\nUsage:\n  %s\n", usage)
		fmt.Println(examples)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}
