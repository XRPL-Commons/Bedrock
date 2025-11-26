package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/builder"
	"github.com/xrpl-commons/bedrock/pkg/config"
)

var (
	buildRelease bool
	buildWatch   bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build smart contract",
	Long:  `Build the smart contract to WASM. Defaults to release mode for deployment-ready builds.`,
	RunE:  runBuild,
}

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolVarP(&buildRelease, "release", "r", true, "Build in release mode (optimized)")
	buildCmd.Flags().BoolVarP(&buildWatch, "watch", "w", false, "Watch for changes and rebuild")
}

func runBuild(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	mode := "debug"
	if buildRelease {
		mode = "release"
	}

	color.Cyan("Building smart contract\n")
	fmt.Printf("   Mode: %s\n", mode)
	fmt.Printf("   Source: %s\n", cfg.Build.Source)
	fmt.Println()

	if buildWatch {
		color.Yellow("Watch mode not yet implemented\n")
		return fmt.Errorf("watch mode coming soon")
	}

	// Create builder (project root is current directory)
	b := builder.New(".")

	// Build with options
	ctx := cmd.Context()
	result, err := b.Build(ctx, builder.BuildOptions{
		Release: buildRelease,
		Verbose: false,
	})

	if err != nil {
		color.Red("\n✗ Build failed: %v\n", err)
		return err
	}

	color.Green("\n✓ Build completed successfully!\n")
	fmt.Println()

	// Show output details
	fmt.Printf("   Output: %s\n", result.WasmPath)
	fmt.Printf("   Size: %d bytes\n", result.Size)
	fmt.Printf("   Duration: %v\n", result.Duration)

	return nil
}
