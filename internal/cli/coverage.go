package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/tester"
)

var (
	coverageLCOV   bool
	coverageHTML   bool
	coverageOutput string
)

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Run tests with code coverage",
	Long: `Run contract tests with code coverage analysis.

Requires cargo-tarpaulin or cargo-llvm-cov to be installed.

Install with:
  cargo install cargo-tarpaulin
or:
  cargo install cargo-llvm-cov

Examples:
  bedrock coverage
  bedrock coverage --lcov
  bedrock coverage --lcov --output coverage/`,
	RunE: runCoverage,
}

func init() {
	rootCmd.AddCommand(coverageCmd)

	coverageCmd.Flags().BoolVar(&coverageLCOV, "lcov", false, "Generate LCOV format for IDE integration")
	coverageCmd.Flags().BoolVar(&coverageHTML, "html", false, "Generate HTML coverage report")
	coverageCmd.Flags().StringVarP(&coverageOutput, "output", "o", "", "Output directory for coverage files")
}

func runCoverage(cmd *cobra.Command, args []string) error {
	_, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Running tests with code coverage\n")
	fmt.Println()

	ctx := cmd.Context()
	result, err := tester.RunCoverage(ctx, ".", tester.CoverageOptions{
		LCOV:   coverageLCOV,
		HTML:   coverageHTML,
		Output: coverageOutput,
	})

	if err != nil {
		color.Red("Coverage failed: %v\n", err)
		return err
	}

	fmt.Println()
	color.Green("Coverage complete (%v)\n\n", result.Duration)

	if result.Summary != "" {
		fmt.Printf("  %s\n", result.Summary)
	}

	if result.LCOVPath != "" {
		fmt.Printf("  LCOV: %s\n", result.LCOVPath)
	}

	if result.HTMLPath != "" {
		fmt.Printf("  HTML: %s\n", result.HTMLPath)
	}

	return nil
}
