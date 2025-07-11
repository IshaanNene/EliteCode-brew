package system

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/IshaanNene/EliteCode-brew/internal/utils"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Elitecode CLI",
	RunE:  runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	spinner := utils.NewSpinner("Checking for updates...")
	spinner.Start()

	// Check if brew is available
	if _, err := exec.LookPath("brew"); err != nil {
		spinner.Stop()
		color.Yellow("⚠️  Homebrew not found. Please update manually.")
		return nil
	}

	// Update using brew
	updateCmd := exec.Command("brew", "upgrade", "elitecode")
	output, err := updateCmd.CombinedOutput()
	
	spinner.Stop()

	if err != nil {
		color.Red("❌ Update failed: %v", err)
		fmt.Println(string(output))
		return nil
	}

	color.Green("✅ Elitecode updated successfully!")
	return nil
}