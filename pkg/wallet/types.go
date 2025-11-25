package wallet

import (
	"time"
)

// Keystore represents the encrypted wallet data stored on disk
type Keystore struct {
	Version     int       `json:"version"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	EncryptedSeed string  `json:"encrypted_seed"`
	Salt        string    `json:"salt"`
	Nonce       string    `json:"nonce"`
}

// WalletInfo represents wallet information without sensitive data
type WalletInfo struct {
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

// Wallet represents a decrypted wallet in memory
type Wallet struct {
	Name    string
	Address string
	Seed    string
}

// WalletManager handles wallet operations
type WalletManager struct {
	walletsDir string
}

const (
	KeystoreVersion = 1
)