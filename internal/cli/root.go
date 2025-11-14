package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bedrock",
	Short: "The unshakeable foundation for XRPL smart contracts",
	Long: `BEDROCK - XRPL Smart Contract CLI
The foundation for XRPL smart contract development

Build, deploy, and interact with XRPL smart contracts written in Rust.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
