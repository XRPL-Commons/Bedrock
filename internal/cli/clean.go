package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/embedded"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cached files and force reinstall of dependencies",
	Long: `Clean removes the bedrock cache directory, which contains:
  - Extracted JavaScript modules (deploy.js, call.js, faucet.js)
  - Installed npm dependencies (node_modules)
  - Version tracking file

After cleaning, the next command that requires JS modules will
automatically reinstall all dependencies fresh.

Use this command if you:
  - Experience issues with JavaScript modules
  - Want to force a fresh reinstall after updating bedrock
  - Need to free up disk space`,
	RunE: runClean,
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}

func runClean(cmd *cobra.Command, args []string) error {
	color.Cyan("Cleaning bedrock cache\n")
	fmt.Println()

	// Get and display cache location
	cacheDir, err := embedded.GetCacheDir()
	if err != nil {
		color.Red("✗ Failed to get cache directory: %v\n", err)
		return err
	}

	fmt.Printf("  Cache location: %s\n", cacheDir)
	fmt.Println()

	// Clean the cache
	if err := embedded.CleanCache(); err != nil {
		color.Red("✗ Failed to clean cache: %v\n", err)
		return err
	}

	color.Green("✓ Cache cleaned successfully\n")
	fmt.Println()
	color.White("  JavaScript modules will be reinstalled on next use.\n")

	return nil
}
