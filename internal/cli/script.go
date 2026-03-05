package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/script"
)

var (
	scriptTemplate string
)

var scriptCmd = &cobra.Command{
	Use:   "script <file>",
	Short: "Run a deployment/interaction script",
	Long: `Execute a multi-step deployment and interaction script.

Scripts are TOML or JSON files defining a sequence of operations:
deploy, call, fund, wait, and assert.

Variables from previous steps can be referenced in subsequent steps
using ${step_name.field} syntax.

Examples:
  bedrock script deploy.toml
  bedrock script --template deploy-and-test`,
	Args: cobra.MaximumNArgs(1),
	RunE: runScript,
}

func init() {
	rootCmd.AddCommand(scriptCmd)

	scriptCmd.Flags().StringVar(&scriptTemplate, "template", "", "Use a built-in template (deploy-and-test)")
}

func runScript(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	var s *script.Script
	var scriptLabel string

	if scriptTemplate != "" {
		var err error
		s, err = script.LoadTemplate(scriptTemplate)
		if err != nil {
			color.Red("Failed to load template: %v\n", err)
			return err
		}
		scriptLabel = fmt.Sprintf("template:%s", scriptTemplate)
	} else if len(args) > 0 {
		var err error
		s, err = script.ParseScript(args[0])
		if err != nil {
			color.Red("Failed to parse script: %v\n", err)
			return err
		}
		scriptLabel = args[0]
	} else {
		return fmt.Errorf("provide a script file or use --template")
	}

	color.Cyan("Running script: %s\n\n", scriptLabel)

	fmt.Printf("  Name: %s\n", s.Name)
	fmt.Printf("  Steps: %d\n", len(s.Steps))
	fmt.Println()

	runner := script.NewRunner(cfg, false)
	ctx := cmd.Context()

	result, err := runner.Run(ctx, s)
	if err != nil {
		color.Red("Script execution failed: %v\n", err)
		return err
	}

	// Display results
	for i, step := range result.Steps {
		if step.Success {
			color.Green("  [%d] %s (%s) -- ok (%v)\n", i+1, step.Name, step.Action, step.Duration)
			if step.Output != nil {
				for k, v := range step.Output {
					if v != "" {
						fmt.Printf("       %s: %s\n", k, v)
					}
				}
			}
		} else {
			color.Red("  [%d] %s (%s) -- FAILED (%v)\n", i+1, step.Name, step.Action, step.Duration)
			if step.Error != "" {
				color.Red("       Error: %s\n", step.Error)
			}
		}
	}

	fmt.Println()
	if result.Success {
		color.Green("Script completed successfully (%v)\n", result.Duration)
	} else {
		color.Red("Script failed (%v)\n", result.Duration)
		return fmt.Errorf("script execution failed")
	}

	return nil
}
