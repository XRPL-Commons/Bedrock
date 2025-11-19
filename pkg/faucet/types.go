package faucet

// FaucetResult represents the result of a faucet request
type FaucetResult struct {
	TxHash        string `json:"txHash"`
	WalletAddress string `json:"walletAddress"`
	WalletSeed    string `json:"walletSeed"`
	Balance       string `json:"balance"`
	FaucetAmount  string `json:"faucetAmount"`
}

// FaucetConfig holds configuration for requesting from faucet
type FaucetConfig struct {
	FaucetURL     string
	WalletSeed    string
	WalletAddress string
	NetworkURL    string
}
