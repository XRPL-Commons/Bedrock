package wallet

import (
	"fmt"
	"strings"

	"github.com/Peersyst/xrpl-go/keypairs"
	"github.com/Peersyst/xrpl-go/pkg/crypto"
	"github.com/Peersyst/xrpl-go/xrpl/interfaces"
	xrplwallet "github.com/Peersyst/xrpl-go/xrpl/wallet"
)

// XRPLWallet provides XRPL-specific wallet functionality using native Go implementation
type XRPLWallet struct{}

// NewXRPLWallet creates a new XRPL wallet helper
func NewXRPLWallet() (*XRPLWallet, error) {
	return &XRPLWallet{}, nil
}

// GenerateWallet creates a new XRPL wallet with random seed using secp256k1
func (x *XRPLWallet) GenerateWallet(name string) (*Wallet, error) {
	return x.GenerateWalletWithAlgorithm(name, "secp256k1")
}

// GenerateWalletWithAlgorithm creates a new XRPL wallet with specified algorithm
func (x *XRPLWallet) GenerateWalletWithAlgorithm(name, algorithm string) (*Wallet, error) {
	// Get the crypto implementation based on algorithm
	alg, err := getAlgorithm(algorithm)
	if err != nil {
		return nil, err
	}

	// Generate a new wallet using xrpl-go
	w, err := xrplwallet.New(alg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate wallet: %w", err)
	}

	return &Wallet{
		Name:    name,
		Address: string(w.ClassicAddress),
		Seed:    w.Seed,
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

	// Try to derive keypair to validate the seed
	_, _, err := keypairs.DeriveKeypair(seed, false)
	if err != nil {
		return fmt.Errorf("invalid seed: %w", err)
	}

	return nil
}

// SeedToAddress derives an XRPL address from a seed using secp256k1
func (x *XRPLWallet) SeedToAddress(seed string) (string, error) {
	return x.SeedToAddressWithAlgorithm(seed, "secp256k1")
}

// SeedToAddressWithAlgorithm derives an XRPL address from a seed using specified algorithm
func (x *XRPLWallet) SeedToAddressWithAlgorithm(seed, algorithm string) (string, error) {
	if err := x.ValidateSeed(seed); err != nil {
		return "", err
	}

	// Derive wallet from seed using xrpl-go
	w, err := xrplwallet.FromSeed(seed, "")
	if err != nil {
		return "", fmt.Errorf("failed to derive address: %w", err)
	}

	return string(w.ClassicAddress), nil
}

// ValidateAddress checks if an address is in valid XRPL format
func (x *XRPLWallet) ValidateAddress(address string) error {
	if len(address) == 0 {
		return fmt.Errorf("address cannot be empty")
	}

	if !strings.HasPrefix(address, "r") {
		return fmt.Errorf("XRPL address must start with 'r'")
	}

	// Additional validation could be added here using xrpl-go's address codec
	// For now, basic format check is sufficient
	if len(address) < 25 || len(address) > 35 {
		return fmt.Errorf("invalid XRPL address length")
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

// getAlgorithm converts algorithm string to xrpl-go CryptoImplementation
func getAlgorithm(algorithm string) (interfaces.CryptoImplementation, error) {
	switch strings.ToLower(algorithm) {
	case "secp256k1", "":
		return crypto.SECP256K1(), nil
	case "ed25519":
		return crypto.ED25519(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s (supported: secp256k1, ed25519)", algorithm)
	}
}
