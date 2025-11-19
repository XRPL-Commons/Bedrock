package faucet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/xrpl-bedrock/bedrock/pkg/adapter"
)

// Faucet handles requesting funds from XRPL faucets
type Faucet struct {
	executor *adapter.Executor
}

// NewFaucet creates a new Faucet instance
func NewFaucet(verbose bool) (*Faucet, error) {
	executor, err := adapter.NewExecutor(verbose)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor: %w", err)
	}

	return &Faucet{
		executor: executor,
	}, nil
}

// Request requests funds from the faucet
func (f *Faucet) Request(ctx context.Context, config FaucetConfig) (*FaucetResult, error) {
	// Prepare config for JS module
	jsConfig := map[string]interface{}{
		"faucet_url":     config.FaucetURL,
		"wallet_seed":    config.WalletSeed,
		"wallet_address": config.WalletAddress,
		"network_url":    config.NetworkURL,
		"verbose":        f.executor != nil, // Use executor's verbose setting
	}

	// Execute faucet module
	result, err := f.executor.ExecuteModule(ctx, "faucet.js", jsConfig)
	if err != nil {
		return nil, err
	}

	// Parse result
	var faucetResult FaucetResult
	if err := json.Unmarshal(result.Data, &faucetResult); err != nil {
		return nil, fmt.Errorf("failed to parse faucet result: %w", err)
	}

	return &faucetResult, nil
}
