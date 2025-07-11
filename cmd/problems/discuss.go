package problems

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/api"
)

var discussCmd = &cobra.Command{
	Use:   "discuss <problem_id>",
	Short: "Open problem discussion in browser",
	Args:  cobra.ExactArgs(1),
	RunE:  runDiscuss,
}

func runDiscuss(cmd *cobra.Command, args []string) error {
	problemID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid problem ID: %v", err)
	}

	client := api.NewClient()
	discussionURL, err := client.GetDiscussionURL(problemID)
	if err != nil {
		return fmt.Errorf("failed to get discussion URL: %v", err)
	}

	if err := openBrowser(discussionURL); err != nil {
		color.Red("❌ Failed to open browser: %v", err)
		color.Blue("Discussion URL: %s", discussionURL)
		return nil
	}

	color.Green("✅ Opened discussion in browser")
	return nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}