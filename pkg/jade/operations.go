package jade

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xrpl-commons/bedrock/pkg/adapter"
)

// Operations handles XRPL network operations via the jade.js module
type Operations struct {
	executor *adapter.Executor
	verbose  bool
}

// NewOperations creates a new Operations instance
func NewOperations(verbose bool) (*Operations, error) {
	executor, err := adapter.NewExecutor(verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	return &Operations{
		executor: executor,
		verbose:  verbose,
	}, nil
}

// Config represents the configuration for jade operations
type Config struct {
	Operation  string                 `json:"operation"`
	NetworkURL string                 `json:"network_url"`
	NetworkID  uint32                 `json:"network_id"`
	Params     map[string]interface{} `json:"params"`
	Verbose    bool                   `json:"verbose"`
}

// BalanceResult represents the result of a balance query
type BalanceResult struct {
	Address      string `json:"address"`
	Balance      string `json:"balance"`
	BalanceDrops string `json:"balance_drops"`
	Funded       bool   `json:"funded,omitempty"`
}

// SendResult represents the result of a send operation
type SendResult struct {
	TxHash      string `json:"tx_hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Amount      string `json:"amount"`
	AmountDrops string `json:"amount_drops"`
	Result      string `json:"result"`
	Fee         string `json:"fee"`
	Sequence    int    `json:"sequence"`
	Validated   bool   `json:"validated"`
}

// TxResult represents the result of a transaction query
type TxResult struct {
	Hash            string      `json:"hash"`
	Type            string      `json:"type"`
	Account         string      `json:"account"`
	Result          string      `json:"result"`
	Fee             string      `json:"fee"`
	Sequence        int         `json:"sequence"`
	Validated       bool        `json:"validated"`
	LedgerIndex     int         `json:"ledger_index"`
	Date            int         `json:"date,omitempty"`
	Destination     string      `json:"destination,omitempty"`
	Amount          interface{} `json:"amount,omitempty"`
	DeliveredAmount interface{} `json:"delivered_amount,omitempty"`
	ContractAccount string      `json:"contract_account,omitempty"`
	FunctionName    string      `json:"function_name,omitempty"`
	ReturnValue     string      `json:"return_value,omitempty"`
	WasmSize        int         `json:"wasm_size,omitempty"`
}

// AccountInfoResult represents the result of an account_info query
type AccountInfoResult struct {
	Address           string      `json:"address"`
	Balance           interface{} `json:"balance"`
	BalanceDrops      string      `json:"balance_drops"`
	Sequence          int         `json:"sequence"`
	OwnerCount        int         `json:"owner_count"`
	PreviousTxnID     string      `json:"previous_txn_id"`
	PreviousTxnLgrSeq int         `json:"previous_txn_lgr_seq"`
	Flags             int         `json:"flags"`
	LedgerIndex       int         `json:"ledger_index"`
	Funded            bool        `json:"funded,omitempty"`
	Error             string      `json:"error,omitempty"`
}

// GetBalanceString returns the balance as a string
func (a *AccountInfoResult) GetBalanceString() string {
	switch v := a.Balance.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.6f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ServerInfoResult represents the result of a server_info query
type ServerInfoResult struct {
	BuildVersion     string                 `json:"build_version"`
	CompleteLedgers  string                 `json:"complete_ledgers"`
	HostID           string                 `json:"hostid"`
	NetworkID        int                    `json:"network_id"`
	Peers            int                    `json:"peers"`
	PubkeyNode       string                 `json:"pubkey_node"`
	ServerState      string                 `json:"server_state"`
	Uptime           int                    `json:"uptime"`
	ValidatedLedger  map[string]interface{} `json:"validated_ledger,omitempty"`
}

// GetBalance retrieves the XRP balance for an address
func (o *Operations) GetBalance(networkURL string, networkID uint32, address string) (*BalanceResult, error) {
	config := Config{
		Operation:  "balance",
		NetworkURL: networkURL,
		NetworkID:  networkID,
		Params: map[string]interface{}{
			"address": address,
		},
		Verbose: o.verbose,
	}

	result, err := o.execute(config)
	if err != nil {
		return nil, err
	}

	var balanceResult BalanceResult
	if err := json.Unmarshal(result.Data, &balanceResult); err != nil {
		return nil, fmt.Errorf("failed to parse balance result: %w", err)
	}

	return &balanceResult, nil
}

// Send transfers XRP to a destination address
func (o *Operations) Send(networkURL string, networkID uint32, walletSeed, destination, amount, algorithm string) (*SendResult, error) {
	config := Config{
		Operation:  "send",
		NetworkURL: networkURL,
		NetworkID:  networkID,
		Params: map[string]interface{}{
			"wallet_seed": walletSeed,
			"destination": destination,
			"amount":      amount,
			"algorithm":   algorithm,
		},
		Verbose: o.verbose,
	}

	result, err := o.execute(config)
	if err != nil {
		return nil, err
	}

	var sendResult SendResult
	if err := json.Unmarshal(result.Data, &sendResult); err != nil {
		return nil, fmt.Errorf("failed to parse send result: %w", err)
	}

	return &sendResult, nil
}

// GetTransaction retrieves transaction details by hash
func (o *Operations) GetTransaction(networkURL string, networkID uint32, hash string) (*TxResult, error) {
	config := Config{
		Operation:  "tx",
		NetworkURL: networkURL,
		NetworkID:  networkID,
		Params: map[string]interface{}{
			"hash": hash,
		},
		Verbose: o.verbose,
	}

	result, err := o.execute(config)
	if err != nil {
		return nil, err
	}

	var txResult TxResult
	if err := json.Unmarshal(result.Data, &txResult); err != nil {
		return nil, fmt.Errorf("failed to parse transaction result: %w", err)
	}

	return &txResult, nil
}

// GetAccountInfo retrieves detailed account information
func (o *Operations) GetAccountInfo(networkURL string, networkID uint32, address string) (*AccountInfoResult, error) {
	config := Config{
		Operation:  "account_info",
		NetworkURL: networkURL,
		NetworkID:  networkID,
		Params: map[string]interface{}{
			"address": address,
		},
		Verbose: o.verbose,
	}

	result, err := o.execute(config)
	if err != nil {
		return nil, err
	}

	var accountResult AccountInfoResult
	if err := json.Unmarshal(result.Data, &accountResult); err != nil {
		return nil, fmt.Errorf("failed to parse account info result: %w", err)
	}

	return &accountResult, nil
}

// GetServerInfo retrieves XRPL server information
func (o *Operations) GetServerInfo(networkURL string, networkID uint32) (*ServerInfoResult, error) {
	config := Config{
		Operation:  "server_info",
		NetworkURL: networkURL,
		NetworkID:  networkID,
		Params:     map[string]interface{}{},
		Verbose:    o.verbose,
	}

	result, err := o.execute(config)
	if err != nil {
		return nil, err
	}

	var serverResult ServerInfoResult
	if err := json.Unmarshal(result.Data, &serverResult); err != nil {
		return nil, fmt.Errorf("failed to parse server info result: %w", err)
	}

	return &serverResult, nil
}

// execute runs the jade.js module with the given configuration
func (o *Operations) execute(config Config) (*adapter.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := o.executor.ExecuteModule(ctx, "jade.js", config)
	if err != nil {
		return nil, fmt.Errorf("jade operation failed: %w", err)
	}

	return result, nil
}
