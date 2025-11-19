package cli

import (
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/config"
	"github.com/xrpl-bedrock/bedrock/pkg/deployer"
)

var slateCmd = &cobra.Command{
	Use:   "slate",
	Short: "Deploy and manage contract deployments",
	Long:  `Write contract deployments to the ledger in layers.`,
}

var slateDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy contract to network",
	Long: `Deploys a compiled WASM contract to an XRPL network.

Uses the generated ABI to include parameter definitions in the deployment.
The contract must be built first with 'bedrock flint build'.`,
	RunE: runSlateDeploy,
}

var (
	deployNetwork    string
	deployWasmPath   string
	deployABIPath    string
	deployWalletSeed string
	deployFaucet     string
)

func init() {
	rootCmd.AddCommand(slateCmd)
	slateCmd.AddCommand(slateDeployCmd)

	// Deploy flags
	slateDeployCmd.Flags().StringVarP(&deployNetwork, "network", "n", "alphanet", "Network to deploy to (alphanet, testnet, mainnet, or local)")
	slateDeployCmd.Flags().StringVarP(&deployWasmPath, "wasm", "w", "", "Path to WASM file (defaults to build output)")
	slateDeployCmd.Flags().StringVarP(&deployABIPath, "abi", "a", "abi.json", "Path to ABI file")
	slateDeployCmd.Flags().StringVarP(&deployWalletSeed, "wallet", "", "", "Wallet seed to use (generates new if not provided)")
	slateDeployCmd.Flags().StringVarP(&deployFaucet, "faucet", "f", "", "Faucet URL for funding wallet")
}

func runSlateDeploy(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("Slate - Deploying smart contract\n")

	// Determine WASM path
	wasmPath := deployWasmPath
	if wasmPath == "" {
		// Default to build output (always use release for deployment)
		buildMode := "release"
		// Build.Source is the source file (e.g., contract/src/lib.rs)
		// We need the contract directory, which is the parent of src
		contractDir := filepath.Dir(filepath.Dir(cfg.Build.Source))
		wasmPath = filepath.Join(contractDir, "target", "wasm32-unknown-unknown", buildMode, cfg.Project.Name+".wasm")
	}

	// Get network configuration
	networkCfg, ok := cfg.Networks[deployNetwork]
	if !ok {
		// Try local_node if network is "local"
		if deployNetwork == "local" {
			networkCfg = config.NetworkConfig{
				URL: "ws://localhost:6006",
			}
		} else {
			return fmt.Errorf("network '%s' not found in config", deployNetwork)
		}
	}

	// Determine faucet URL
	faucetURL := deployFaucet
	if faucetURL == "" && networkCfg.FaucetURL != "" {
		faucetURL = networkCfg.FaucetURL
	}

	color.White("   Network: %s\n", deployNetwork)
	color.White("   URL: %s\n", networkCfg.URL)
	color.White("   WASM: %s\n", wasmPath)
	color.White("   ABI: %s\n", deployABIPath)
	if deployWalletSeed != "" {
		color.White("   Wallet: Using provided seed\n")
	} else {
		color.White("   Wallet: Generating new wallet\n")
	}
	fmt.Println()

	// Create deployer
	verbose := false // TODO: add verbose flag
	d, err := deployer.NewDeployer(verbose)
	if err != nil {
		color.Red("\n✗ Failed to initialize deployer: %v\n", err)
		return err
	}

	// Deploy contract
	color.White("Executing deployment...\n\n")

	ctx := cmd.Context()
	result, err := d.Deploy(ctx, deployer.DeploymentConfig{
		WasmPath:   wasmPath,
		ABIPath:    deployABIPath,
		NetworkURL: networkCfg.URL,
		WalletSeed: deployWalletSeed,
		FaucetURL:  faucetURL,
	})

	if err != nil {
		color.Red("\n✗ Deployment failed: %v\n", err)
		return err
	}

	// Display results
	fmt.Println()
	color.Green("✓ Contract deployed successfully!\n")
	fmt.Println()
	color.Cyan("Deployment Details:\n")
	fmt.Printf("  Transaction Hash: %s\n", result.TxHash)
	fmt.Printf("  Wallet Address: %s\n", result.WalletAddress)
	fmt.Printf("  Wallet Seed: %s\n", result.WalletSeed)
	if result.ContractAccount != "" {
		fmt.Printf("  Contract Account: %s\n", result.ContractAccount)
	}

	fmt.Println()
	color.Yellow("Tip: Save the wallet seed to interact with the contract later\n")
	color.Yellow("Tip: Use 'bedrock crystal call' to call contract functions\n")

	// TODO: Save deployment info to bedrock.toml [deployments] section

	return nil
}
