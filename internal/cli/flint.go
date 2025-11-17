package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/builder"
	"github.com/xrpl-bedrock/bedrock/pkg/config"
)

var flintCmd = &cobra.Command{
	Use:   "flint",
	Short: "Build and compile smart contracts",
	Long:  `Compile Rust smart contracts to WASM using cargo.`,
}

var flintBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the smart contract",
	Long:  `Compiles the Rust smart contract to WASM format.`,
	RunE:  runFlintBuild,
}

var flintCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean build artifacts",
	Long:  `Removes all build artifacts and compiled files.`,
	RunE:  runFlintClean,
}

var (
	buildRelease bool
	buildVerbose bool
)

func init() {
	rootCmd.AddCommand(flintCmd)
	flintCmd.AddCommand(flintBuildCmd)
	flintCmd.AddCommand(flintCleanCmd)

	// Build flags
	flintBuildCmd.Flags().BoolVarP(&buildRelease, "release", "r", false, "Build with release optimizations")
	flintBuildCmd.Flags().BoolVarP(&buildVerbose, "verbose", "v", false, "Show verbose build output")
}

func runFlintBuild(cmd *cobra.Command, args []string) error {
	// Load configuration to verify we're in a project
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Flint - Building smart contract\n")
	if buildRelease {
		color.White("   Mode: Release (optimized)\n")
	} else {
		color.White("   Mode: Debug\n")
	}
	color.White("   Source: %s\n\n", cfg.Build.Source)

	// Create spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Compiling Rust → WASM..."
	s.Start()

	// Create builder
	b := builder.New(".")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Build
	result, err := b.Build(ctx, builder.BuildOptions{
		Release: buildRelease,
		Verbose: buildVerbose,
	})

	s.Stop()

	if err != nil {
		color.Red("\n✗ Build failed: %v\n", err)
		return err
	}

	// Success output
	color.Green("\n✓ Build completed successfully!\n\n")
	fmt.Printf("Output:   %s\n", result.WasmPath)
	fmt.Printf("Size:     %s\n", formatBytes(result.Size))
	fmt.Printf("Duration: %s\n", result.Duration.Round(time.Millisecond))

	if result.Optimized {
		color.Yellow("\nBuilt with release optimizations (smaller size, slower build)\n")
	} else {
		color.Yellow("\nTip: Use --release flag for optimized builds\n")
	}

	return nil
}

func runFlintClean(cmd *cobra.Command, args []string) error {
	// Load configuration to verify we're in a project
	if _, err := config.LoadFromWorkingDir(); err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Cleaning build artifacts...\n")

	b := builder.New(".")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := b.Clean(ctx); err != nil {
		color.Red("\n✗ Clean failed: %v\n", err)
		return err
	}

	color.Green("\n✓ Build artifacts cleaned!\n")
	return nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
