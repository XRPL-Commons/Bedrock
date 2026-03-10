package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/abi"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/doc"
)

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate contract documentation",
	Long: `Generate Markdown documentation from contract ABI annotations.

Parses @xrpl-function annotations in your Rust source and generates
a function reference with types and descriptions.

Examples:
  bedrock doc
  bedrock doc --output docs/api`,
	RunE: runDoc,
}

var docOutput string

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.Flags().StringVarP(&docOutput, "output", "o", "", "Output directory (default from config or docs/api)")
}

func runDoc(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	outputDir := docOutput
	if outputDir == "" {
		outputDir = cfg.Doc.Output
		if outputDir == "" {
			outputDir = config.DefaultDocConfig().Output
		}
	}

	color.Cyan("Generating documentation\n")
	fmt.Printf("  Source: %s\n", cfg.Build.Source)
	fmt.Printf("  Output: %s\n\n", outputDir)

	// Parse ABI from source
	sourceDir := filepath.Dir(cfg.Build.Source)
	parser := abi.NewParser(sourceDir)
	abiData, err := parser.ParseContract(cfg.Project.Name)
	if err != nil {
		color.Red("Failed to parse contract: %v\n", err)
		return err
	}

	// Generate documentation
	gen := doc.NewGenerator(outputDir)
	path, err := gen.Generate(abiData)
	if err != nil {
		color.Red("Failed to generate docs: %v\n", err)
		return err
	}

	color.Green("Documentation generated: %s\n", path)
	fmt.Printf("  Functions: %d\n", len(abiData.Functions))

	return nil
}
