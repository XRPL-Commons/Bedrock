package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/caller"
	"github.com/xrpl-bedrock/bedrock/pkg/config"
)

var crystalCmd = &cobra.Command{
	Use:   "crystal",
	Short: "Inspect and interact with deployed contracts",
	Long:  `Transparent view into contract state and function calls.`,
}

var crystalCallCmd = &cobra.Command{
	Use:   "call <contract-account> <function-name>",
	Short: "Call a contract function",
	Long: `Calls a function on a deployed smart contract.

Uses the ABI to properly format parameters and decode return values.
Requires a funded wallet to pay for transaction fees and gas.`,
	Args: cobra.ExactArgs(2),
	RunE: runCrystalCall,
}

var (
	callNetwork    string
	callABIPath    string
	callWalletSeed string
	callParams     string
	callParamsFile string
	callGas        string
	callFee        string
)

func init() {
	rootCmd.AddCommand(crystalCmd)
	crystalCmd.AddCommand(crystalCallCmd)

	// Call flags
	crystalCallCmd.Flags().StringVarP(&callNetwork, "network", "n", "alphanet", "Network to call on (alphanet, testnet, mainnet, or local)")
	crystalCallCmd.Flags().StringVarP(&callABIPath, "abi", "a", "abi.json", "Path to ABI file")
	crystalCallCmd.Flags().StringVarP(&callWalletSeed, "wallet", "w", "", "Wallet seed (required)")
	crystalCallCmd.Flags().StringVarP(&callParams, "params", "p", "", "Parameters as JSON string")
	crystalCallCmd.Flags().StringVarP(&callParamsFile, "params-file", "f", "", "Parameters from JSON file")
	crystalCallCmd.Flags().StringVarP(&callGas, "gas", "g", "1000000", "Computation allowance")
	crystalCallCmd.Flags().StringVarP(&callFee, "fee", "", "1000000", "Transaction fee in drops")

	crystalCallCmd.MarkFlagRequired("wallet")
}

func runCrystalCall(cmd *cobra.Command, args []string) error {
	contractAccount := args[0]
	functionName := args[1]

	// Load configuration
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Crystal - Calling contract function\n")

	// Get network configuration
	networkCfg, ok := cfg.Networks[callNetwork]
	if !ok {
		if callNetwork == "local" {
			networkCfg = config.NetworkConfig{
				URL: "ws://localhost:6006",
			}
		} else {
			return fmt.Errorf("network '%s' not found in config", callNetwork)
		}
	}

	// Determine parameters
	var parameters string
	if callParamsFile != "" {
		// Load from file
		data, err := os.ReadFile(callParamsFile)
		if err != nil {
			return fmt.Errorf("failed to read params file: %w", err)
		}
		parameters = string(data)
	} else if callParams != "" {
		parameters = callParams
	}

	// Validate parameters is valid JSON if provided
	if parameters != "" {
		var testJSON interface{}
		if err := json.Unmarshal([]byte(parameters), &testJSON); err != nil {
			return fmt.Errorf("invalid JSON parameters: %w", err)
		}
	}

	color.White("   Network: %s\n", callNetwork)
	color.White("   URL: %s\n", networkCfg.URL)
	color.White("   Contract: %s\n", contractAccount)
	color.White("   Function: %s\n", functionName)
	color.White("   ABI: %s\n", callABIPath)
	if parameters != "" {
		color.White("   Parameters: %s\n", parameters)
	}
	fmt.Println()

	// Create caller
	verbose := false // TODO: add verbose flag
	c, err := caller.NewCaller(verbose)
	if err != nil {
		color.Red("\n✗ Failed to initialize caller: %v\n", err)
		return err
	}

	// Parse parameters from JSON string
	var params map[string]interface{}
	if parameters != "" {
		if err := json.Unmarshal([]byte(parameters), &params); err != nil {
			color.Red("\n✗ Failed to parse parameters JSON: %v\n", err)
			return fmt.Errorf("invalid parameters JSON: %w", err)
		}
	}

	// Call contract
	color.White("Executing contract call...\n\n")

	ctx := cmd.Context()
	result, err := c.Call(ctx, caller.CallConfig{
		ContractAccount:      contractAccount,
		FunctionName:         functionName,
		NetworkURL:           networkCfg.URL,
		WalletSeed:           callWalletSeed,
		ABIPath:              callABIPath,
		Parameters:           params,
		ComputationAllowance: callGas,
		Fee:                  callFee,
	})

	if err != nil {
		color.Red("\n✗ Contract call failed: %v\n", err)
		return err
	}

	// Display results
	fmt.Println()
	color.Green("✓ Contract function called successfully!\n")
	fmt.Println()
	color.Cyan("Call Results:\n")
	fmt.Printf("  Transaction Hash: %s\n", result.TxHash)
	fmt.Printf("  Return Code: %d", result.ReturnCode)
	if result.ReturnCode == 0 {
		color.Green(" (SUCCESS)\n")
	} else {
		color.Red(" (ERROR)\n")
	}

	if result.ReturnValue != "" {
		fmt.Printf("  Return Value (hex): %s\n", result.ReturnValue)
		// Try to parse as decimal
		if len(result.ReturnValue) > 0 {
			fmt.Printf("  Return Value (dec): %d\n", parseInt(result.ReturnValue))
		}
	}

	if result.GasUsed > 0 {
		fmt.Printf("  Gas Used: %d\n", result.GasUsed)
	}

	return nil
}

// parseInt parses hex string to int64
func parseInt(hex string) int64 {
	var result int64
	fmt.Sscanf(hex, "%x", &result)
	return result
}
