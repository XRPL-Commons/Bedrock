package jade

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/Peersyst/xrpl-go/xrpl/queries/account"
	"github.com/Peersyst/xrpl-go/xrpl/queries/common"
	"github.com/Peersyst/xrpl-go/xrpl/queries/server"
	"github.com/Peersyst/xrpl-go/xrpl/queries/transactions"
	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	rpctypes "github.com/Peersyst/xrpl-go/xrpl/rpc/types"
	"github.com/Peersyst/xrpl-go/xrpl/transaction"
	"github.com/Peersyst/xrpl-go/xrpl/transaction/types"
	"github.com/Peersyst/xrpl-go/xrpl/wallet"
)

// Operations handles XRPL network operations using pure Go
type Operations struct {
	verbose bool
}

// NewOperations creates a new Operations instance
func NewOperations(verbose bool) (*Operations, error) {
	return &Operations{
		verbose: verbose,
	}, nil
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
	BuildVersion    string                 `json:"build_version"`
	CompleteLedgers string                 `json:"complete_ledgers"`
	HostID          string                 `json:"hostid"`
	NetworkID       int                    `json:"network_id"`
	Peers           int                    `json:"peers"`
	PubkeyNode      string                 `json:"pubkey_node"`
	ServerState     string                 `json:"server_state"`
	Uptime          int                    `json:"uptime"`
	ValidatedLedger map[string]interface{} `json:"validated_ledger,omitempty"`
}

// createRPCClient creates an RPC client for the given network URL
func createRPCClient(networkURL string) (*rpc.Client, error) {
	// Convert WebSocket URL to HTTP URL for RPC
	rpcURL := networkURL
	if strings.HasPrefix(networkURL, "ws://") {
		rpcURL = strings.Replace(networkURL, "ws://", "http://", 1)
	} else if strings.HasPrefix(networkURL, "wss://") {
		rpcURL = strings.Replace(networkURL, "wss://", "https://", 1)
	}

	cfg, err := rpc.NewClientConfig(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client config: %w", err)
	}

	return rpc.NewClient(cfg), nil
}

// dropsToXRP converts drops to XRP string using integer arithmetic for precision
func dropsToXRP(drops string) string {
	dropsBig, ok := new(big.Int).SetString(drops, 10)
	if !ok {
		return "0"
	}

	million := big.NewInt(1_000_000)
	whole := new(big.Int).Div(dropsBig, million)
	remainder := new(big.Int).Mod(dropsBig, million)

	if remainder.Sign() == 0 {
		return whole.String()
	}

	remainderStr := fmt.Sprintf("%06d", remainder.Int64())
	remainderStr = strings.TrimRight(remainderStr, "0")
	return fmt.Sprintf("%s.%s", whole.String(), remainderStr)
}

// xrpToDrops converts XRP string to drops string using big.Float for precision
func xrpToDrops(xrp string) (string, error) {
	xrpBig, _, err := big.ParseFloat(xrp, 10, 128, big.ToNearestEven)
	if err != nil {
		return "", fmt.Errorf("invalid XRP amount: %w", err)
	}

	if xrpBig.Sign() <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}

	multiplier := big.NewFloat(1_000_000)
	dropsBig := new(big.Float).Mul(xrpBig, multiplier)

	dropsInt, _ := dropsBig.Int(nil)
	return dropsInt.String(), nil
}

// GetBalance retrieves the XRP balance for an address
func (o *Operations) GetBalance(networkURL string, address string) (*BalanceResult, error) {
	client, err := createRPCClient(networkURL)
	if err != nil {
		return nil, err
	}

	if o.verbose {
		fmt.Printf("Getting balance for %s from %s...\n", address, networkURL)
	}

	req := &account.InfoRequest{
		Account: types.Address(address),
	}

	resp, err := client.Request(req)
	if err != nil {
		// Check if account not found
		if strings.Contains(err.Error(), "actNotFound") || strings.Contains(err.Error(), "Account not found") {
			return &BalanceResult{
				Address:      address,
				Balance:      "0",
				BalanceDrops: "0",
				Funded:       false,
			}, nil
		}
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	var infoResp account.InfoResponse
	if err := resp.GetResult(&infoResp); err != nil {
		// Check if account not found in response
		if strings.Contains(err.Error(), "actNotFound") {
			return &BalanceResult{
				Address:      address,
				Balance:      "0",
				BalanceDrops: "0",
				Funded:       false,
			}, nil
		}
		return nil, fmt.Errorf("failed to parse account info: %w", err)
	}

	balanceDrops := infoResp.AccountData.Balance.String()
	balanceXRP := dropsToXRP(balanceDrops)

	return &BalanceResult{
		Address:      address,
		Balance:      balanceXRP,
		BalanceDrops: balanceDrops,
		Funded:       true,
	}, nil
}

// Send transfers XRP to a destination address
func (o *Operations) Send(networkURL string, walletSeed, destination, amount, algorithm string) (*SendResult, error) {
	client, err := createRPCClient(networkURL)
	if err != nil {
		return nil, err
	}

	// Create wallet from seed
	w, err := wallet.FromSeed(walletSeed, algorithm)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet from seed: %w", err)
	}

	if o.verbose {
		fmt.Printf("Sending %s XRP from %s to %s...\n", amount, w.ClassicAddress, destination)
	}

	// Convert XRP to drops
	amountDrops, err := xrpToDrops(amount)
	if err != nil {
		return nil, err
	}

	// Create payment transaction
	payment := transaction.FlatTransaction{
		"TransactionType": "Payment",
		"Account":         string(w.ClassicAddress),
		"Destination":     destination,
		"Amount":          amountDrops,
	}

	// Submit and wait for validation
	opts := &rpctypes.SubmitOptions{
		Autofill: true,
		Wallet:   &w,
		FailHard: false,
	}

	txResp, err := client.SubmitTxAndWait(payment, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to submit transaction: %w", err)
	}

	// Extract result from meta
	result := "tesSUCCESS"
	if meta, ok := txResp.Meta.(map[string]interface{}); ok {
		if txResult, ok := meta["TransactionResult"].(string); ok {
			result = txResult
		}
	}

	// Get fee from transaction
	fee := ""
	if feeVal, ok := txResp.Tx["Fee"]; ok {
		fee = fmt.Sprintf("%v", feeVal)
	}

	// Get sequence from transaction
	sequence := 0
	if seqVal, ok := txResp.Tx["Sequence"]; ok {
		switch v := seqVal.(type) {
		case float64:
			sequence = int(v)
		case int:
			sequence = v
		}
	}

	return &SendResult{
		TxHash:      string(txResp.Hash),
		From:        string(w.ClassicAddress),
		To:          destination,
		Amount:      amount,
		AmountDrops: amountDrops,
		Result:      result,
		Fee:         fee,
		Sequence:    sequence,
		Validated:   txResp.Validated,
	}, nil
}

// GetTransaction retrieves transaction details by hash
func (o *Operations) GetTransaction(networkURL string, hash string) (*TxResult, error) {
	client, err := createRPCClient(networkURL)
	if err != nil {
		return nil, err
	}

	if o.verbose {
		fmt.Printf("Fetching transaction %s...\n", hash)
	}

	req := &transactions.TxRequest{
		Transaction: hash,
		Binary:      false,
	}

	resp, err := client.Request(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var txResp transactions.TxResponse
	if err := resp.GetResult(&txResp); err != nil {
		return nil, fmt.Errorf("failed to parse transaction: %w", err)
	}

	// Extract transaction type
	txType := ""
	if typeVal, ok := txResp.Tx["TransactionType"]; ok {
		txType = fmt.Sprintf("%v", typeVal)
	}

	// Extract account
	txAccount := ""
	if accVal, ok := txResp.Tx["Account"]; ok {
		txAccount = fmt.Sprintf("%v", accVal)
	}

	// Extract fee
	fee := ""
	if feeVal, ok := txResp.Tx["Fee"]; ok {
		fee = fmt.Sprintf("%v", feeVal)
	}

	// Extract sequence
	sequence := 0
	if seqVal, ok := txResp.Tx["Sequence"]; ok {
		switch v := seqVal.(type) {
		case float64:
			sequence = int(v)
		case int:
			sequence = v
		}
	}

	// Extract result from meta
	result := "tesSUCCESS"
	if meta, ok := txResp.Meta.(map[string]interface{}); ok {
		if txResult, ok := meta["TransactionResult"].(string); ok {
			result = txResult
		}
	}

	txResult := &TxResult{
		Hash:        string(txResp.Hash),
		Type:        txType,
		Account:     txAccount,
		Result:      result,
		Fee:         fee,
		Sequence:    sequence,
		Validated:   txResp.Validated,
		LedgerIndex: int(txResp.LedgerIndex),
		Date:        int(txResp.Date),
	}

	// Add type-specific fields
	if txType == "Payment" {
		if destVal, ok := txResp.Tx["Destination"]; ok {
			txResult.Destination = fmt.Sprintf("%v", destVal)
		}
		if amtVal, ok := txResp.Tx["Amount"]; ok {
			txResult.Amount = amtVal
		}
		if meta, ok := txResp.Meta.(map[string]interface{}); ok {
			if delivered, ok := meta["delivered_amount"]; ok {
				txResult.DeliveredAmount = delivered
			}
		}
	} else if txType == "ContractCreate" {
		if meta, ok := txResp.Meta.(map[string]interface{}); ok {
			if contractAcc, ok := meta["ContractAccount"].(string); ok {
				txResult.ContractAccount = contractAcc
			}
		}
		if wasmHex, ok := txResp.Tx["WasmHex"].(string); ok {
			txResult.WasmSize = len(wasmHex) / 2
		}
	} else if txType == "ContractCall" {
		if contractAcc, ok := txResp.Tx["ContractAccount"].(string); ok {
			txResult.ContractAccount = contractAcc
		}
		if funcName, ok := txResp.Tx["FunctionName"].(string); ok {
			txResult.FunctionName = funcName
		}
		if meta, ok := txResp.Meta.(map[string]interface{}); ok {
			if returnVal, ok := meta["HookReturnString"].(string); ok {
				txResult.ReturnValue = returnVal
			}
		}
	}

	return txResult, nil
}

// GetAccountInfo retrieves detailed account information
func (o *Operations) GetAccountInfo(networkURL string, address string) (*AccountInfoResult, error) {
	client, err := createRPCClient(networkURL)
	if err != nil {
		return nil, err
	}

	if o.verbose {
		fmt.Printf("Getting account info for %s...\n", address)
	}

	req := &account.InfoRequest{
		Account:     types.Address(address),
		LedgerIndex: common.Validated,
	}

	resp, err := client.Request(req)
	if err != nil {
		// Check if account not found
		if strings.Contains(err.Error(), "actNotFound") || strings.Contains(err.Error(), "Account not found") {
			return &AccountInfoResult{
				Address: address,
				Funded:  false,
				Error:   "Account not found (not funded)",
			}, nil
		}
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	var infoResp account.InfoResponse
	if err := resp.GetResult(&infoResp); err != nil {
		if strings.Contains(err.Error(), "actNotFound") {
			return &AccountInfoResult{
				Address: address,
				Funded:  false,
				Error:   "Account not found (not funded)",
			}, nil
		}
		return nil, fmt.Errorf("failed to parse account info: %w", err)
	}

	balanceDrops := infoResp.AccountData.Balance.String()
	balanceXRP := dropsToXRP(balanceDrops)

	return &AccountInfoResult{
		Address:           address,
		Balance:           balanceXRP,
		BalanceDrops:      balanceDrops,
		Sequence:          int(infoResp.AccountData.Sequence),
		OwnerCount:        int(infoResp.AccountData.OwnerCount),
		PreviousTxnID:     string(infoResp.AccountData.PreviousTxnID),
		PreviousTxnLgrSeq: int(infoResp.AccountData.PreviousTxnLgrSeq),
		Flags:             int(infoResp.AccountData.Flags),
		LedgerIndex:       int(infoResp.LedgerIndex),
		Funded:            true,
	}, nil
}

// GetServerInfo retrieves XRPL server information
func (o *Operations) GetServerInfo(networkURL string) (*ServerInfoResult, error) {
	client, err := createRPCClient(networkURL)
	if err != nil {
		return nil, err
	}

	if o.verbose {
		fmt.Printf("Getting server info from %s...\n", networkURL)
	}

	req := &server.InfoRequest{}

	resp, err := client.Request(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get server info: %w", err)
	}

	var infoResp server.InfoResponse
	if err := resp.GetResult(&infoResp); err != nil {
		return nil, fmt.Errorf("failed to parse server info: %w", err)
	}

	result := &ServerInfoResult{
		BuildVersion:    infoResp.Info.BuildVersion,
		CompleteLedgers: infoResp.Info.CompleteLedgers,
		HostID:          infoResp.Info.HostID,
		NetworkID:       int(infoResp.Info.NetworkID),
		Peers:           int(infoResp.Info.Peers),
		PubkeyNode:      infoResp.Info.PubkeyNode,
		ServerState:     infoResp.Info.ServerState,
		Uptime:          int(infoResp.Info.Uptime),
	}

	// Add validated ledger info if available
	if infoResp.Info.ValidatedLedger.Seq > 0 {
		result.ValidatedLedger = map[string]interface{}{
			"hash": string(infoResp.Info.ValidatedLedger.Hash),
			"seq":  float64(infoResp.Info.ValidatedLedger.Seq),
			"age":  float64(infoResp.Info.ValidatedLedger.Age),
		}
	}

	return result, nil
}
