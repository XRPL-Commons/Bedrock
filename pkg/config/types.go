package config

import "time"

// Config represents the bedrock.toml configuration
type Config struct {
	Project     ProjectConfig            `toml:"project"`
	Build       BuildConfig              `toml:"build"`
	Networks    map[string]NetworkConfig `toml:"networks"`
	Contracts   map[string]ContractConfig `toml:"contracts"`
	Deployments map[string]DeploymentInfo `toml:"deployments"`
	Wallets     WalletsConfig            `toml:"wallets"`
	LocalNode   LocalNodeConfig          `toml:"local_node"`
}

type ProjectConfig struct {
	Name    string   `toml:"name"`
	Version string   `toml:"version"`
	Authors []string `toml:"authors,omitempty"`
}

type BuildConfig struct {
	Source string `toml:"source"`
	Output string `toml:"output"`
	Target string `toml:"target"`
}

type NetworkConfig struct {
	URL       string `toml:"url"`
	NetworkID uint32 `toml:"network_id"`
	FaucetURL string `toml:"faucet_url,omitempty"`
	Explorer  string `toml:"explorer,omitempty"`
}

type ContractConfig struct {
	Source string `toml:"source"`
	ABI    string `toml:"abi"`
}

type DeploymentInfo struct {
	ContractAccount string    `toml:"contract_account"`
	ContractID      string    `toml:"contract_id"`
	TxHash          string    `toml:"tx_hash"`
	DeployedAt      time.Time `toml:"deployed_at"`
	Network         string    `toml:"network"`
}

type WalletsConfig struct {
	Default  string `toml:"default,omitempty"`
	Keystore string `toml:"keystore"`
}

// LocalNodeConfig points to the directory containing rippled config files
type LocalNodeConfig struct {
	ConfigDir      string `toml:"config_dir"`
	DockerImage    string `toml:"docker_image"`
	LedgerInterval int    `toml:"ledger_interval"` // Interval in milliseconds for ledger advancement
}

// DefaultLocalNodeConfig returns default local node configuration
func DefaultLocalNodeConfig() LocalNodeConfig {
	return LocalNodeConfig{
		ConfigDir:      ".bedrock/node-config",
		DockerImage:    "transia/alphanet:latest",
		LedgerInterval: 1000, // 1 second default
	}
}
