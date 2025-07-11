package problems

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var bookmarkCmd = &cobra.Command{
	Use:   "bookmark <problem_id>",
	Short: "Bookmark a problem",
	Args:  cobra.ExactArgs(1),
	RunE:  runBookmark,
}

func init() {
	bookmarkCmd.Flags().BoolP("remove", "r", false, "Remove bookmark")
}

func runBookmark(cmd *cobra.Command, args []string) error {
	problemID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid problem ID: %v", err)
	}

	remove, _ := cmd.Flags().GetBool("remove")

	client := api.NewClient()
	
	if remove {
		err = client.RemoveBookmark(problemID)
		if err != nil {
			return fmt.Errorf("failed to remove bookmark: %v", err)
		}
		color.Green("✅ Bookmark removed for problem %d", problemID)
	} else {
		err = client.AddBookmark(problemID)
		if err != nil {
			return fmt.Errorf("failed to add bookmark: %v", err)
		}
		color.Green("✅ Problem %d bookmarked", problemID)
	}

	return nil
}