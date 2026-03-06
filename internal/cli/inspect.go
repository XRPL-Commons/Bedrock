package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/inspector"
)

var (
	inspectABI      bool
	inspectWasmInfo bool
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [contract-path]",
	Short: "Inspect contract artifacts",
	Long: `Display detailed information about compiled contract artifacts.

Shows WASM size, function count, ABI details, imports, exports,
and storage layout.

Examples:
  bedrock inspect
  bedrock inspect --abi
  bedrock inspect --wasm-info
  bedrock inspect contract/target/wasm32-unknown-unknown/release/my_contract.wasm`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.Flags().BoolVar(&inspectABI, "abi", false, "Output just the ABI JSON")
	inspectCmd.Flags().BoolVar(&inspectWasmInfo, "wasm-info", false, "Show WASM module details (imports, exports, memory)")
}

func runInspect(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	// Find WASM path
	var wasmPath string
	if len(args) > 0 {
		wasmPath = args[0]
	} else {
		wasmPath = findWasmFile(cfg)
	}

	// ABI-only mode
	if inspectABI {
		abiData, err := inspector.InspectABI("abi.json")
		if err != nil {
			return err
		}
		pretty, _ := json.MarshalIndent(abiData, "", "  ")
		fmt.Println(string(pretty))
		return nil
	}

	color.Cyan("Contract Inspection\n\n")

	// ABI inspection
	abiData, err := inspector.InspectABI("abi.json")
	if err != nil {
		color.Yellow("  ABI: not found (run 'bedrock build' first)\n")
	} else {
		fmt.Printf("  Contract: %s\n", abiData.ContractName)
		fmt.Printf("  Functions: %d\n", len(abiData.Functions))

		for _, fn := range abiData.Functions {
			var paramTypes []string
			for _, p := range fn.Parameters {
				paramTypes = append(paramTypes, p.Type)
			}
			retStr := ""
			if fn.Returns != nil {
				retStr = fmt.Sprintf(" -> %s", fn.Returns.Type)
			}
			fmt.Printf("    %s(%s)%s\n", fn.Name, strings.Join(paramTypes, ", "), retStr)
		}
	}

	// WASM inspection
	if wasmPath != "" {
		if _, err := os.Stat(wasmPath); err == nil {
			fmt.Println()
			wasmInfo, err := inspector.InspectWasm(wasmPath)
			if err != nil {
				color.Yellow("  WASM: failed to inspect: %v\n", err)
			} else {
				fmt.Printf("  WASM Path: %s\n", wasmPath)
				fmt.Printf("  WASM Size: %d bytes (%.1f KB)\n", wasmInfo.Size, float64(wasmInfo.Size)/1024)
				fmt.Printf("  Memory Pages: %d (%d KB)\n", wasmInfo.MemPages, wasmInfo.MemPages*64)
				fmt.Printf("  Exports: %d functions\n", len(wasmInfo.Functions))

				if inspectWasmInfo {
					fmt.Println()
					color.Cyan("  Exported Functions:\n")
					for _, name := range wasmInfo.Functions {
						fmt.Printf("    - %s\n", name)
					}

					if len(wasmInfo.Imports) > 0 {
						fmt.Println()
						color.Cyan("  Imports:\n")
						for _, imp := range wasmInfo.Imports {
							fmt.Printf("    - %s\n", imp)
						}
					}
				}
			}
		} else {
			color.Yellow("  WASM: not found (run 'bedrock build' first)\n")
		}
	}

	return nil
}

func findWasmFile(cfg *config.Config) string {
	// Try config-derived path first
	if cfg.Build.Output != "" {
		pattern := filepath.Join(cfg.Build.Output, "*.wasm")
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}

	// Fallback to common paths
	patterns := []string{
		filepath.Join("contract", "target", "wasm32-unknown-unknown", "release", "*.wasm"),
		filepath.Join("contract", "target", "wasm32-unknown-unknown", "debug", "*.wasm"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}
