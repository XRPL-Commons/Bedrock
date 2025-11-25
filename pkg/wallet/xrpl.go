package wallet

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

// XRPLWallet provides XRPL-specific wallet functionality
type XRPLWallet struct{}

// NewXRPLWallet creates a new XRPL wallet helper
func NewXRPLWallet() *XRPLWallet {
	return &XRPLWallet{}
}

// GenerateWallet creates a new XRPL wallet with random seed
func (x *XRPLWallet) GenerateWallet(name string) (*Wallet, error) {
	// Generate 16 bytes of entropy for the seed
	entropy := make([]byte, 16)
	if _, err := rand.Read(entropy); err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	// Create seed in XRPL format (s + base58)
	seed := "s" + base58.Encode(entropy)

	// Derive address from seed
	address, err := x.SeedToAddress(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address from seed: %w", err)
	}

	return &Wallet{
		Name:    name,
		Address: address,
		Seed:    seed,
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

// SeedToAddress derives an XRPL address from a seed
func (x *XRPLWallet) SeedToAddress(seed string) (string, error) {
	if err := x.ValidateSeed(seed); err != nil {
		return "", err
	}

	// For now, we'll generate a mock address since full XRPL key derivation
	// requires the complete cryptographic implementation
	// In production, this would use proper XRPL key derivation
	seedBytes := []byte(seed)
	hash := sha256.Sum256(seedBytes)
	
	// Take first 20 bytes and create a mock XRPL address
	addressBytes := hash[:20]
	
	// Add XRPL address prefix (0x00) and checksum
	payload := append([]byte{0x00}, addressBytes...)
	checksum := x.calculateChecksum(payload)
	fullAddress := append(payload, checksum...)
	
	return "r" + base58.Encode(fullAddress), nil
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
	if err := x.ValidateSeed(seed); err != nil {
		return nil, fmt.Errorf("invalid seed: %w", err)
	}

	address, err := x.SeedToAddress(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %w", err)
	}

	return &Wallet{
		Name:    name,
		Address: address,
		Seed:    seed,
	}, nil
}

// calculateChecksum calculates a simple checksum for address generation
func (x *XRPLWallet) calculateChecksum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:4]
}