package cmd

import(
	_"fmt"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
    Use:   "run",
    Short: "run using basic testcases",
    Run: func(cmd *cobra.Command, args []string) {
        
    },
}