package cli

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bedrock",
	Short: "The unshakeable foundation for XRPL smart contracts",
	Long: `BEDROCK - XRPL Smart Contract CLI
The foundation for XRPL smart contract development

Build, deploy, and interact with XRPL smart contracts written in Rust.`,
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	// Global flags
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
