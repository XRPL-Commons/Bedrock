package deployer

// DeploymentResult represents the result of a contract deployment
type DeploymentResult struct {
	TxHash          string                 `json:"txHash"`
	WalletAddress   string                 `json:"walletAddress"`
	WalletSeed      string                 `json:"walletSeed"`
	ContractAccount string                 `json:"contractAccount"`
	ContractIndex   string                 `json:"contractIndex"`
	Validated       bool                   `json:"validated"`
	Meta            map[string]interface{} `json:"meta"`
}

// DeploymentConfig holds configuration for deploying a contract
type DeploymentConfig struct {
	WasmPath   string
	ABIPath    string
	NetworkURL string
	WalletSeed string
	Algorithm  string
	FaucetURL  string
	Fee        string
}
