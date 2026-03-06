package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/tester"
)

var (
	snapshotDiff  bool
	snapshotCheck bool
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage gas snapshots",
	Long: `Save and compare computation cost snapshots.

Run integration tests and save gas usage to a snapshot file.
Use --diff to compare against a saved snapshot.
Use --check in CI to fail if gas usage has changed.

Examples:
  bedrock snapshot
  bedrock snapshot --diff
  bedrock snapshot --check`,
	RunE: runSnapshot,
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	snapshotCmd.Flags().BoolVar(&snapshotDiff, "diff", false, "Compare against saved snapshot")
	snapshotCmd.Flags().BoolVar(&snapshotCheck, "check", false, "Check snapshot matches (CI mode, exit non-zero if changed)")
}

func runSnapshot(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	snapshotFile := cfg.Snapshot.File
	if snapshotFile == "" {
		snapshotFile = config.DefaultSnapshotConfig().File
	}

	if snapshotDiff || snapshotCheck {
		return runSnapshotCompare(cmd, cfg, snapshotFile)
	}

	return runSnapshotSave(cmd, cfg, snapshotFile)
}

func runSnapshotSave(cmd *cobra.Command, cfg *config.Config, snapshotFile string) error {
	color.Cyan("Running integration tests for gas snapshot\n")
	fmt.Println()

	runner := tester.NewIntegrationRunner(".", cfg, false)
	ctx := cmd.Context()

	results, err := runner.Run(ctx, tester.TestOptions{GasReport: true})
	if err != nil {
		color.Red("Integration tests failed: %v\n", err)
		return err
	}

	// Collect gas entries
	var entries []tester.GasEntry
	for _, suite := range results {
		for _, test := range suite.Tests {
			if test.GasUsed > 0 {
				entries = append(entries, tester.GasEntry{
					Function: test.Name,
					GasUsed:  test.GasUsed,
				})
			}
		}
	}

	if len(entries) == 0 {
		color.Yellow("No gas data collected from integration tests\n")
		return nil
	}

	if err := tester.SaveSnapshot(snapshotFile, entries); err != nil {
		color.Red("Failed to save snapshot: %v\n", err)
		return err
	}

	color.Green("Snapshot saved to %s (%d entries)\n", snapshotFile, len(entries))
	fmt.Println()
	fmt.Print(tester.FormatGasReportFromEntries(entries))

	return nil
}

func runSnapshotCompare(cmd *cobra.Command, cfg *config.Config, snapshotFile string) error {
	// Load existing snapshot
	oldSnapshot, err := tester.LoadSnapshot(snapshotFile)
	if err != nil {
		color.Red("Failed to load snapshot: %v\n", err)
		return err
	}

	color.Cyan("Running integration tests for snapshot comparison\n")
	fmt.Println()

	runner := tester.NewIntegrationRunner(".", cfg, false)
	ctx := cmd.Context()

	results, err := runner.Run(ctx, tester.TestOptions{GasReport: true})
	if err != nil {
		color.Red("Integration tests failed: %v\n", err)
		return err
	}

	// Build new snapshot from results
	var entries []tester.GasEntry
	for _, suite := range results {
		for _, test := range suite.Tests {
			if test.GasUsed > 0 {
				entries = append(entries, tester.GasEntry{
					Function: test.Name,
					GasUsed:  test.GasUsed,
				})
			}
		}
	}

	newSnapshot := tester.SnapshotFromEntries(entries)
	diff := tester.DiffSnapshots(oldSnapshot, newSnapshot)

	fmt.Println()
	if diff.HasChanges() {
		color.Yellow("Gas snapshot has changes:\n\n")
		fmt.Println(diff.FormatDiff())

		if snapshotCheck {
			return fmt.Errorf("gas snapshot check failed: changes detected")
		}

		color.Yellow("\nRun 'bedrock snapshot' to update the snapshot file\n")
	} else {
		color.Green("Gas snapshot matches -- no changes detected\n")
	}

	return nil
}
