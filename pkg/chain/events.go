package chain

import (
	"context"
	"encoding/json"
	"fmt"
)

// EventHistory represents event_history RPC response
type EventHistory struct {
	Events      []ContractEvent `json:"events"`
	LedgerIndex int64           `json:"ledger_current_index"`
	Status      string          `json:"status"`
	Marker      interface{}     `json:"marker,omitempty"`
}

// ContractEvent represents a single contract event
type ContractEvent struct {
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	LedgerIndex int64           `json:"ledger_index"`
	TxHash      string          `json:"tx_hash"`
	Contract    string          `json:"contract"`
	Timestamp   int64           `json:"close_time"`
}

// GetEventHistory queries contract event history
func (c *Client) GetEventHistory(ctx context.Context, account string, opts EventQueryOpts) (*EventHistory, error) {
	params := map[string]interface{}{
		"account": account,
	}

	if opts.EventType != "" {
		params["event_type"] = opts.EventType
	}
	if opts.FromLedger > 0 {
		params["ledger_index_min"] = opts.FromLedger
	}
	if opts.ToLedger > 0 {
		params["ledger_index_max"] = opts.ToLedger
	}
	if opts.Limit > 0 {
		params["limit"] = opts.Limit
	}

	var result EventHistory
	if err := c.CallTyped(ctx, &result, "event_history", params); err != nil {
		return nil, fmt.Errorf("event_history failed: %w", err)
	}

	return &result, nil
}

// EventQueryOpts configures event history queries
type EventQueryOpts struct {
	EventType  string
	FromLedger int64
	ToLedger   int64
	Limit      int
}
