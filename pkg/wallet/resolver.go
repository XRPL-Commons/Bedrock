package wallet

import (
	"fmt"
	"strings"
)

// WalletResolver helps resolve wallet names to seeds
type WalletResolver struct {
	manager      *WalletManager
	authProvider *AuthProvider
}

// NewWalletResolver creates a new wallet resolver
func NewWalletResolver() (*WalletResolver, error) {
	manager, err := NewWalletManager()
	if err != nil {
		return nil, err
	}

	authProvider := NewAuthProvider()

	return &WalletResolver{
		manager:      manager,
		authProvider: authProvider,
	}, nil
}

// ResolveWallet resolves a wallet input to a seed
// If input starts with 's', treats it as a raw seed
// Otherwise, treats it as a wallet name and loads from keystore
func (wr *WalletResolver) ResolveWallet(walletInput string) (string, error) {
	if walletInput == "" {
		return "", fmt.Errorf("wallet input cannot be empty")
	}

	// If input starts with 's', treat as raw seed
	if strings.HasPrefix(walletInput, "s") {
		xrplWallet, err := NewXRPLWallet()
		if err != nil {
			return "", fmt.Errorf("failed to create XRPL wallet: %w", err)
		}
		if err := xrplWallet.ValidateSeed(walletInput); err != nil {
			return "", fmt.Errorf("invalid seed format: %w", err)
		}
		return walletInput, nil
	}

	// Otherwise, treat as wallet name
	return wr.resolveWalletName(walletInput)
}

// resolveWalletName resolves a wallet name to its seed
func (wr *WalletResolver) resolveWalletName(walletName string) (string, error) {
	// Check if wallet exists
	if !wr.manager.WalletExists(walletName) {
		return "", fmt.Errorf("wallet '%s' not found. Use 'bedrock jade list' to see available wallets", walletName)
	}

	// Get password and load wallet
	password, err := wr.authProvider.GetPassword("Enter password: ")
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}

	wallet, err := wr.manager.LoadWallet(walletName, password)
	if err != nil {
		return "", fmt.Errorf("failed to load wallet '%s': %w", walletName, err)
	}

	return wallet.Seed, nil
}