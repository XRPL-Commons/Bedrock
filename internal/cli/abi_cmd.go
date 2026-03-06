package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/abi"
)

var abiCmd = &cobra.Command{
	Use:   "abi <encode|decode|inspect> [args...]",
	Short: "ABI encoding and inspection tools",
	Long: `Tools for working with contract ABIs.

Commands:
  encode <function> <params-json>  - Encode parameters for a function call
  decode <data>                     - Decode return values
  inspect <abi-file>               - Display ABI in human-readable format

Examples:
  bedrock abi inspect abi.json
  bedrock abi encode transfer '{"to":"rAddr...","amount":100}'
  bedrock abi decode 0x00000064`,
	Args: cobra.MinimumNArgs(1),
	RunE: runABI,
}

func init() {
	rootCmd.AddCommand(abiCmd)
}

func runABI(cmd *cobra.Command, args []string) error {
	subcommand := args[0]

	switch subcommand {
	case "inspect":
		if len(args) < 2 {
			return fmt.Errorf("usage: bedrock abi inspect <abi-file>")
		}
		return abiInspect(args[1])
	case "encode":
		if len(args) < 3 {
			return fmt.Errorf("usage: bedrock abi encode <function> <params-json>")
		}
		return abiEncode(args[1], args[2])
	case "decode":
		if len(args) < 2 {
			return fmt.Errorf("usage: bedrock abi decode <hex-data>")
		}
		return abiDecode(args[1])
	default:
		return fmt.Errorf("unknown subcommand: %s (use: encode, decode, inspect)", subcommand)
	}
}

func abiInspect(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read ABI file: %w", err)
	}

	var abiData abi.ABI
	if err := json.Unmarshal(data, &abiData); err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	color.Cyan("Contract: %s\n", abiData.ContractName)
	fmt.Printf("Functions: %d\n\n", len(abiData.Functions))

	for _, fn := range abiData.Functions {
		// Build signature
		var paramTypes []string
		for _, p := range fn.Parameters {
			paramTypes = append(paramTypes, p.Type)
		}

		retStr := ""
		if fn.Returns != nil {
			retStr = fmt.Sprintf(" -> %s", fn.Returns.Type)
		}

		color.Green("  %s(%s)%s\n", fn.Name, strings.Join(paramTypes, ", "), retStr)

		for _, p := range fn.Parameters {
			desc := ""
			if p.Description != "" {
				desc = fmt.Sprintf(" - %s", p.Description)
			}
			fmt.Printf("    %s: %s%s\n", p.Name, p.Type, desc)
		}

		if fn.Returns != nil && fn.Returns.Description != "" {
			fmt.Printf("    returns: %s - %s\n", fn.Returns.Type, fn.Returns.Description)
		}

		fmt.Println()
	}

	return nil
}

func abiEncode(function string, paramsJSON string) error {
	// Load ABI from default path
	data, err := os.ReadFile("abi.json")
	if err != nil {
		return fmt.Errorf("failed to read abi.json: %w (specify ABI path or run 'bedrock build')", err)
	}

	var abiData abi.ABI
	if err := json.Unmarshal(data, &abiData); err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Find the function
	var fn *abi.Function
	for i := range abiData.Functions {
		if abiData.Functions[i].Name == function {
			fn = &abiData.Functions[i]
			break
		}
	}

	if fn == nil {
		return fmt.Errorf("function '%s' not found in ABI", function)
	}

	// Parse parameters
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		return fmt.Errorf("invalid parameters JSON: %w", err)
	}

	// Validate parameters against ABI
	for _, p := range fn.Parameters {
		if _, ok := params[p.Name]; !ok {
			color.Yellow("Warning: parameter '%s' (%s) not provided\n", p.Name, p.Type)
		}
	}

	// Output the encoded representation (function name + parameters)
	encoded := map[string]interface{}{
		"function":   function,
		"parameters": params,
	}

	pretty, _ := json.MarshalIndent(encoded, "", "  ")
	fmt.Println(string(pretty))

	return nil
}

func abiDecode(hexData string) error {
	// Strip 0x prefix if present
	hexData = strings.TrimPrefix(hexData, "0x")

	color.Cyan("Decoding: %s\n\n", hexData)

	// Try to interpret as an integer return value
	var value int64
	if _, err := fmt.Sscanf(hexData, "%x", &value); err == nil {
		fmt.Printf("  As integer: %d\n", value)
	}

	// Show raw hex interpretation
	fmt.Printf("  Hex: 0x%s\n", hexData)
	fmt.Printf("  Bytes: %d\n", len(hexData)/2)

	return nil
}
