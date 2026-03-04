package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion <bash|zsh|fish|powershell>",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for bedrock.

To load completions:

Bash:
  $ source <(bedrock completion bash)
  # To load on startup, add to ~/.bashrc:
  $ echo 'source <(bedrock completion bash)' >> ~/.bashrc

Zsh:
  $ source <(bedrock completion zsh)
  # To load on startup:
  $ bedrock completion zsh > "${fpath[1]}/_bedrock"

Fish:
  $ bedrock completion fish | source
  # To load on startup:
  $ bedrock completion fish > ~/.config/fish/completions/bedrock.fish

PowerShell:
  PS> bedrock completion powershell | Out-String | Invoke-Expression`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE:      runCompletion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) error {
	shell := args[0]

	switch shell {
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		return rootCmd.GenZshCompletion(os.Stdout)
	case "fish":
		return rootCmd.GenFishCompletion(os.Stdout, true)
	case "powershell":
		return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return fmt.Errorf("unsupported shell: %s (use: bash, zsh, fish, powershell)", shell)
	}
}
