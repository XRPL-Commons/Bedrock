package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/xrpl-bedrock/bedrock/pkg/adapter"
)

// XRPLWallet provides XRPL-specific wallet functionality
type XRPLWallet struct {
	executor *adapter.Executor
}

// NewXRPLWallet creates a new XRPL wallet helper
func NewXRPLWallet() (*XRPLWallet, error) {
	executor, err := adapter.NewExecutor(false)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	return &XRPLWallet{
		executor: executor,
	}, nil
}

// GenerateWallet creates a new XRPL wallet with random seed using proper XRPL cryptography
func (x *XRPLWallet) GenerateWallet(name string) (*Wallet, error) {
	return x.GenerateWalletWithAlgorithm(name, "secp256k1")
}

// GenerateWalletWithAlgorithm creates a new XRPL wallet with specified algorithm
func (x *XRPLWallet) GenerateWalletWithAlgorithm(name, algorithm string) (*Wallet, error) {
	// Use JavaScript module for proper XRPL wallet generation
	jsConfig := map[string]interface{}{
		"action":    "generate_wallet",
		"algorithm": algorithm,
	}

	ctx := context.Background()
	result, err := x.executor.ExecuteModule(ctx, "wallet.js", jsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet: %w", err)
	}

	var walletResult struct {
		Success   bool   `json:"success"`
		Error     string `json:"error,omitempty"`
		Seed      string `json:"seed"`
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
	}

	if err := json.Unmarshal(result.Data, &walletResult); err != nil {
		return nil, fmt.Errorf("failed to parse wallet result: %w", err)
	}

	if !walletResult.Success {
		return nil, fmt.Errorf("wallet generation failed: %s", walletResult.Error)
	}

	return &Wallet{
		Name:    name,
		Address: walletResult.Address,
		Seed:    walletResult.Seed,
	}, nil
}

// ValidateSeed checks if a seed is in valid XRPL format
func (x *XRPLWallet) ValidateSeed(seed string) error {
	if len(seed) == 0 {
		return fmt.Errorf("seed cannot be empty")
	}

	if !strings.HasPrefix(seed, "s") {
		return fmt.Errorf("XRPL seed must start with 's'")
	}

	// Check if the rest is valid base58
	seedData := seed[1:]
	if len(seedData) < 20 {
		return fmt.Errorf("seed too short")
	}

	// Try to decode base58
	decoded := base58.Decode(seedData)
	if len(decoded) == 0 {
		return fmt.Errorf("invalid base58 encoding in seed")
	}

	return nil
}

// SeedToAddress derives an XRPL address from a seed using proper XRPL cryptography
func (x *XRPLWallet) SeedToAddress(seed string) (string, error) {
	return x.SeedToAddressWithAlgorithm(seed, "secp256k1")
}

// SeedToAddressWithAlgorithm derives an XRPL address from a seed using specified algorithm
func (x *XRPLWallet) SeedToAddressWithAlgorithm(seed, algorithm string) (string, error) {
	if err := x.ValidateSeed(seed); err != nil {
		return "", err
	}

	// Use JavaScript module for proper XRPL address derivation
	jsConfig := map[string]interface{}{
		"action":    "derive_address",
		"seed":      seed,
		"algorithm": algorithm,
	}

	ctx := context.Background()
	result, err := x.executor.ExecuteModule(ctx, "wallet.js", jsConfig)
	if err != nil {
		return "", fmt.Errorf("failed to derive address: %w", err)
	}

	var walletResult struct {
		Success   bool   `json:"success"`
		Error     string `json:"error,omitempty"`
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
	}

	if err := json.Unmarshal(result.Data, &walletResult); err != nil {
		return "", fmt.Errorf("failed to parse wallet result: %w", err)
	}

	if !walletResult.Success {
		return "", fmt.Errorf("wallet derivation failed: %s", walletResult.Error)
	}

	return walletResult.Address, nil
}

// ValidateAddress checks if an address is in valid XRPL format
func (x *XRPLWallet) ValidateAddress(address string) error {
	if len(address) == 0 {
		return fmt.Errorf("address cannot be empty")
	}

	if !strings.HasPrefix(address, "r") {
		return fmt.Errorf("XRPL address must start with 'r'")
	}

	// Check address format (r + base58)
	addressRegex := regexp.MustCompile(`^r[1-9A-HJ-NP-Za-km-z]{24,34}$`)
	if !addressRegex.MatchString(address) {
		return fmt.Errorf("invalid XRPL address format")
	}

	return nil
}

// ImportWallet creates a wallet from an existing seed
func (x *XRPLWallet) ImportWallet(name, seed string) (*Wallet, error) {
	return x.ImportWalletWithAlgorithm(name, seed, "secp256k1")
}

// ImportWalletWithAlgorithm creates a wallet from an existing seed with specified algorithm
func (x *XRPLWallet) ImportWalletWithAlgorithm(name, seed, algorithm string) (*Wallet, error) {
	if err := x.ValidateSeed(seed); err != nil {
		return nil, fmt.Errorf("invalid seed: %w", err)
	}

	address, err := x.SeedToAddressWithAlgorithm(seed, algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %w", err)
	}

	return &Wallet{
		Name:    name,
		Address: address,
		Seed:    seed,
	}, nil
}

