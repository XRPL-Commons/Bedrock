package chain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is a Go-native XRPL JSON-RPC client
type Client struct {
	httpURL    string
	httpClient *http.Client
}

// RPCRequest represents a JSON-RPC request
type RPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params,omitempty"`
}

// RPCResponse represents a JSON-RPC response
type RPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *RPCError       `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new XRPL RPC client from a WebSocket or HTTP URL
func NewClient(url string) *Client {
	httpURL := toHTTPURL(url)
	return &Client{
		httpURL: httpURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Call sends a JSON-RPC request and returns the raw result
func (c *Client) Call(ctx context.Context, method string, params ...interface{}) (json.RawMessage, error) {
	req := RPCRequest{
		Method: method,
	}
	if len(params) > 0 {
		req.Params = params
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.httpURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// CallTyped sends an RPC request and unmarshals the result into the target
func (c *Client) CallTyped(ctx context.Context, target interface{}, method string, params ...interface{}) error {
	raw, err := c.Call(ctx, method, params...)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, target); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}
	return nil
}

// toHTTPURL converts a WebSocket URL to an HTTP URL for JSON-RPC
func toHTTPURL(url string) string {
	if strings.HasPrefix(url, "ws://") {
		// ws://localhost:6006 -> http://localhost:5005
		url = strings.Replace(url, "ws://", "http://", 1)
		url = strings.Replace(url, ":6006", ":5005", 1)
	} else if strings.HasPrefix(url, "wss://") {
		url = strings.Replace(url, "wss://", "https://", 1)
	}
	return url
}
