package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for Elitecode CLI.
To load completions:

Bash:
  $ source <(elitecode completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ elitecode completion bash > /etc/bash_completion.d/elitecode
  # macOS:
  $ elitecode completion bash > /usr/local/etc/bash_completion.d/elitecode

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ elitecode completion zsh > "${fpath[1]}/_elitecode"

  # You will need to start a new shell for this setup to take effect.

Fish:
  $ elitecode completion fish | source

  # To load completions for each session, execute once:
  $ elitecode completion fish > ~/.config/fish/completions/elitecode.fish

PowerShell:
  PS> elitecode completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> elitecode completion powershell > elitecode.ps1
  # and source this file from your PowerShell profile.`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return fmt.Errorf("unsupported shell type %q", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
