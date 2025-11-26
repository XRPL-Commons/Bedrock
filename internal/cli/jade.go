package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xrpl-commons/bedrock/pkg/wallet"
)

var (
	jadeAlgorithm string
)

var jadeCmd = &cobra.Command{
	Use:   "jade",
	Short: "Wallet management for XRPL",
	Long: `Jade - Wallet management for XRPL smart contracts

Create, import, and manage XRPL wallets securely with encryption.
Never expose your seeds in command line arguments again.`,
}

var jadeNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new XRPL wallet",
	Long: `Create a new XRPL wallet with a randomly generated seed.

The wallet will be encrypted and stored securely on disk.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		walletName := args[0]

		// Initialize components
		walletManager, err := wallet.NewWalletManager()
		if err != nil {
			fmt.Printf("Error: Failed to initialize wallet manager: %v\n", err)
			os.Exit(1)
		}

		xrplWallet, err := wallet.NewXRPLWallet()
		if err != nil {
			fmt.Printf("Error: Failed to initialize XRPL wallet: %v\n", err)
			os.Exit(1)
		}
		authProvider := wallet.NewAuthProvider()

		// Generate new wallet
		newWallet, err := xrplWallet.GenerateWalletWithAlgorithm(walletName, jadeAlgorithm)
		if err != nil {
			fmt.Printf("Error: Failed to generate wallet: %v\n", err)
			os.Exit(1)
		}

		// Get password for encryption
		password, err := authProvider.GetPassword("Enter password to encrypt wallet: ")
		if err != nil {
			fmt.Printf("Error: Failed to read password: %v\n", err)
			os.Exit(1)
		}

		if len(password) == 0 {
			fmt.Println("Error: Password cannot be empty")
			os.Exit(1)
		}

		// Save encrypted wallet
		if err := walletManager.SaveWallet(newWallet, password); err != nil {
			fmt.Printf("Error: Failed to save wallet: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Wallet '%s' created successfully\n", walletName)
		fmt.Printf("   Address: %s\n", newWallet.Address)
		fmt.Printf("   Stored at: ~/.config/bedrock/wallets/%s.json\n", walletName)
	},
}

var jadeImportCmd = &cobra.Command{
	Use:   "import <name>",
	Short: "Import an existing XRPL wallet from seed",
	Long: `Import an existing XRPL wallet using its seed.

The seed input will be hidden for security.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		walletName := args[0]

		// Initialize components
		walletManager, err := wallet.NewWalletManager()
		if err != nil {
			fmt.Printf("Error: Failed to initialize wallet manager: %v\n", err)
			os.Exit(1)
		}

		xrplWallet, err := wallet.NewXRPLWallet()
		if err != nil {
			fmt.Printf("Error: Failed to initialize XRPL wallet: %v\n", err)
			os.Exit(1)
		}
		authProvider := wallet.NewAuthProvider()

		// Get seed (hidden input)
		seed, err := authProvider.GetPassword("Enter seed: ")
		if err != nil {
			fmt.Printf("Error: Failed to read seed: %v\n", err)
			os.Exit(1)
		}

		if len(seed) == 0 {
			fmt.Println("Error: Seed cannot be empty")
			os.Exit(1)
		}

		// Import wallet
		importedWallet, err := xrplWallet.ImportWalletWithAlgorithm(walletName, seed, jadeAlgorithm)
		if err != nil {
			fmt.Printf("Error: Failed to import wallet: %v\n", err)
			os.Exit(1)
		}

		// Get password for encryption
		password, err := authProvider.GetPassword("Enter password to encrypt wallet: ")
		if err != nil {
			fmt.Printf("Error: Failed to read password: %v\n", err)
			os.Exit(1)
		}

		if len(password) == 0 {
			fmt.Println("Error: Password cannot be empty")
			os.Exit(1)
		}

		// Save encrypted wallet
		if err := walletManager.SaveWallet(importedWallet, password); err != nil {
			fmt.Printf("Error: Failed to save wallet: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Wallet '%s' imported successfully\n", walletName)
		fmt.Printf("   Address: %s\n", importedWallet.Address)
		fmt.Printf("   Stored at: ~/.config/bedrock/wallets/%s.json\n", walletName)
	},
}

var jadeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored wallets",
	Long:  `List all wallets stored in the keystore with their addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		walletManager, err := wallet.NewWalletManager()
		if err != nil {
			fmt.Printf("Error: Failed to initialize wallet manager: %v\n", err)
			os.Exit(1)
		}

		wallets, err := walletManager.ListWallets()
		if err != nil {
			fmt.Printf("Error: Failed to list wallets: %v\n", err)
			os.Exit(1)
		}

		if len(wallets) == 0 {
			fmt.Println("No wallets found")
			return
		}

		fmt.Printf("Found %d wallet(s):\n\n", len(wallets))
		for _, w := range wallets {
			fmt.Printf("Name:    %s\n", w.Name)
			fmt.Printf("Address: %s\n", w.Address)
			fmt.Printf("Created: %s\n", w.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println()
		}
	},
}

var jadeExportCmd = &cobra.Command{
	Use:   "export <name>",
	Short: "Export wallet seed and address",
	Long: `Export the seed and address of a stored wallet.

Requires password authentication.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		walletName := args[0]

		// Initialize components
		walletManager, err := wallet.NewWalletManager()
		if err != nil {
			fmt.Printf("Error: Failed to initialize wallet manager: %v\n", err)
			os.Exit(1)
		}

		authProvider := wallet.NewAuthProvider()

		// Check if wallet exists
		if !walletManager.WalletExists(walletName) {
			fmt.Printf("Error: Wallet '%s' not found\n", walletName)
			os.Exit(1)
		}

		// Get password for wallet decryption
		password, err := authProvider.GetPassword("Enter password: ")
		if err != nil {
			fmt.Printf("Error: Failed to read password: %v\n", err)
			os.Exit(1)
		}

		// Load wallet
		loadedWallet, err := walletManager.LoadWallet(walletName, password)
		if err != nil {
			fmt.Printf("Error: Failed to load wallet: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Wallet: %s\n", loadedWallet.Name)
		fmt.Printf("Address: %s\n", loadedWallet.Address)
		fmt.Printf("Seed: %s\n", loadedWallet.Seed)
	},
}

var jadeRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a wallet from storage",
	Long: `Permanently remove a wallet from storage.

This action cannot be undone. Make sure you have backed up the seed if needed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		walletName := args[0]

		walletManager, err := wallet.NewWalletManager()
		if err != nil {
			fmt.Printf("Error: Failed to initialize wallet manager: %v\n", err)
			os.Exit(1)
		}

		// Check if wallet exists
		if !walletManager.WalletExists(walletName) {
			fmt.Printf("Error: Wallet '%s' not found\n", walletName)
			os.Exit(1)
		}

		// Confirm removal
		fmt.Printf("Are you sure you want to remove wallet '%s'? (y/N): ", walletName)
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Operation cancelled")
			return
		}

		// Remove wallet
		if err := walletManager.RemoveWallet(walletName); err != nil {
			fmt.Printf("Error: Failed to remove wallet: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Wallet '%s' removed successfully\n", walletName)
	},
}

func init() {
	// Add subcommands to jade
	jadeCmd.AddCommand(jadeNewCmd)
	jadeCmd.AddCommand(jadeImportCmd)
	jadeCmd.AddCommand(jadeListCmd)
	jadeCmd.AddCommand(jadeExportCmd)
	jadeCmd.AddCommand(jadeRemoveCmd)

	// Add flags to commands that need algorithm selection
	jadeNewCmd.Flags().StringVarP(&jadeAlgorithm, "algorithm", "a", "secp256k1", "Cryptographic algorithm (secp256k1, ed25519)")
	jadeImportCmd.Flags().StringVarP(&jadeAlgorithm, "algorithm", "a", "secp256k1", "Cryptographic algorithm (secp256k1, ed25519)")

	// Add jade to root command
	RootCmd.AddCommand(jadeCmd)
}
