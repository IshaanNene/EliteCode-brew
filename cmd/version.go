package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	Version = "0.1.0"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long: `Show version information about Elitecode CLI.
This includes:
- Version number
- Git commit hash
- Build date
- Go version
- Operating system and architecture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("\n%s\n", color.CyanString("Elitecode CLI"))
		fmt.Println(color.YellowString(strings.Repeat("-", 30)))

		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Built:      %s\n", BuildDate)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
