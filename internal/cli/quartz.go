package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/abi"
	"github.com/xrpl-bedrock/bedrock/pkg/config"
)

var quartzCmd = &cobra.Command{
	Use:   "quartz",
	Short: "Generate and manage contract ABIs",
	Long:  `Extract ABI definitions from contract source code annotations.`,
}

var quartzGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ABI from contract source",
	Long: `Parses Rust source code annotations to generate abi.json.

Reads @xrpl-function, @param, and @return annotations from your contract
source code and generates a JSON ABI file compatible with deployment tools.`,
	RunE: runQuartzGenerate,
}

var (
	outputPath   string
	generateJS   bool
	contractName string
)

func init() {
	rootCmd.AddCommand(quartzCmd)
	quartzCmd.AddCommand(quartzGenerateCmd)

	// Generate flags
	quartzGenerateCmd.Flags().StringVarP(&outputPath, "output", "o", "abi.json", "Output file path")
	quartzGenerateCmd.Flags().BoolVarP(&generateJS, "js", "j", false, "Also generate JavaScript module")
	quartzGenerateCmd.Flags().StringVarP(&contractName, "name", "n", "", "Contract name (defaults to project name)")
}

func runQuartzGenerate(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Quartz - Generating contract ABI\n")
	color.White("   Source: %s\n\n", cfg.Build.Source)

	// Determine contract name
	name := contractName
	if name == "" {
		name = cfg.Project.Name
	}

	// Parse contract source
	contractDir := filepath.Join(".", "contract", "src")
	parser := abi.NewParser(contractDir)

	color.White("Parsing Rust source files...\n")
	contractABI, err := parser.ParseContract(name)
	if err != nil {
		color.Red("\n✗ Failed to parse contract: %v\n", err)
		return err
	}

	// Validate ABI
	generator := abi.NewGenerator(".")
	if validationErrors := generator.Validate(contractABI); len(validationErrors) > 0 {
		color.Red("\n✗ ABI validation failed:\n")
		for _, err := range validationErrors {
			fmt.Printf("  - %v\n", err)
		}
		return fmt.Errorf("ABI validation failed with %d error(s)", len(validationErrors))
	}

	color.Green("✓ Found %d function(s)\n\n", len(contractABI.Functions))

	// Display parsed functions
	for _, fn := range contractABI.Functions {
		fmt.Printf("  %s(\n", fn.Name)
		for _, param := range fn.Parameters {
			fmt.Printf("    %s: %s", param.Name, param.Type)
			if param.Description != "" {
				fmt.Printf(" // %s", param.Description)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("  )")
		if fn.Returns != nil {
			fmt.Printf(" -> %s", fn.Returns.Type)
			if fn.Returns.Description != "" {
				fmt.Printf(" // %s", fn.Returns.Description)
			}
		}
		fmt.Printf("\n\n")
	}

	// Generate JSON ABI
	jsonPath, err := generator.Generate(contractABI, outputPath)
	if err != nil {
		color.Red("✗ Failed to generate ABI: %v\n", err)
		return err
	}

	color.Green("✓ Generated ABI: %s\n", jsonPath)

	// Generate JS module if requested
	if generateJS {
		jsPath := outputPath[:len(outputPath)-len(filepath.Ext(outputPath))] + ".js"
		jsPath, err = generator.GenerateJS(contractABI, jsPath)
		if err != nil {
			color.Red("✗ Failed to generate JS module: %v\n", err)
			return err
		}
		color.Green("✓ Generated JS module: %s\n", jsPath)
	}

	color.Yellow("\nTip: Use this ABI with 'bedrock slate deploy' to deploy your contract\n")

	return nil
}
