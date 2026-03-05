package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/config"
	"github.com/xrpl-commons/bedrock/pkg/deployer"
	"github.com/xrpl-commons/bedrock/pkg/wallet"
)

var (
	modifyNetwork   string
	modifyWallet    string
	modifyWasm      string
	modifyABI       string
	modifyAlgorithm string
	modifyFee       string
)

var modifyCmd = &cobra.Command{
	Use:   "modify <contract-account>",
	Short: "Modify a deployed contract",
	Long: `Update a deployed contract's code or ABI via ContractModify transaction.

Examples:
  bedrock modify rContract123... --wallet sXXX... --wasm contract.wasm
  bedrock modify rContract123... --wallet sXXX... --abi abi.json
  bedrock modify rContract123... --wallet sXXX... --wasm contract.wasm --abi abi.json`,
	Args: cobra.ExactArgs(1),
	RunE: runModify,
}

func init() {
	rootCmd.AddCommand(modifyCmd)

	modifyCmd.Flags().StringVarP(&modifyNetwork, "network", "n", "alphanet", "Network")
	modifyCmd.Flags().StringVarP(&modifyWallet, "wallet", "w", "", "Wallet seed or name (required)")
	modifyCmd.Flags().StringVar(&modifyWasm, "wasm", "", "New WASM file path")
	modifyCmd.Flags().StringVar(&modifyABI, "abi", "", "New ABI file path")
	modifyCmd.Flags().StringVar(&modifyAlgorithm, "algorithm", "secp256k1", "Cryptographic algorithm")
	modifyCmd.Flags().StringVar(&modifyFee, "fee", "10000000", "Transaction fee in drops")

	modifyCmd.MarkFlagRequired("wallet")
}

func runModify(cmd *cobra.Command, args []string) error {
	contractAccount := args[0]

	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	networkCfg, ok := cfg.Networks[modifyNetwork]
	if !ok {
		if modifyNetwork == "local" {
			networkCfg = config.NetworkConfig{URL: "ws://localhost:6006", NetworkID: 63456}
		} else {
			return fmt.Errorf("network '%s' not found in config", modifyNetwork)
		}
	}

	if modifyWasm == "" && modifyABI == "" {
		return fmt.Errorf("at least one of --wasm or --abi must be provided")
	}

	color.Cyan("Modifying contract\n")
	fmt.Printf("  Contract: %s\n", contractAccount)
	fmt.Printf("  Network:  %s\n", modifyNetwork)

	resolver, err := wallet.NewWalletResolver()
	if err != nil {
		return fmt.Errorf("failed to initialize wallet resolver: %w", err)
	}

	walletSeed, err := resolver.ResolveWallet(modifyWallet)
	if err != nil {
		return fmt.Errorf("failed to resolve wallet: %w", err)
	}

	d, err := deployer.NewDeployer(false)
	if err != nil {
		return fmt.Errorf("failed to initialize deployer: %w", err)
	}

	ctx := cmd.Context()
	result, err := d.Modify(ctx, deployer.ModifyConfig{
		ContractAccount: contractAccount,
		NetworkURL:      networkCfg.URL,
		WalletSeed:      walletSeed,
		Algorithm:       modifyAlgorithm,
		WasmPath:        modifyWasm,
		ABIPath:         modifyABI,
		Fee:             modifyFee,
	})

	if err != nil {
		color.Red("\n✗ Modification failed: %v\n", err)
		return err
	}

	color.Green("\n✓ Contract modified successfully!\n")
	fmt.Printf("  Transaction Hash: %s\n", result.TxHash)
	fmt.Printf("  Validated: %v\n", result.Validated)

	return nil
}
