package network

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Peersyst/xrpl-go/xrpl/rpc"
)

// LedgerAcceptRequest represents a ledger_accept RPC request
// This is an admin command used in standalone mode to manually close ledgers
type LedgerAcceptRequest struct {
	apiVersion int
}

// Method returns the RPC method name
func (r *LedgerAcceptRequest) Method() string {
	return "ledger_accept"
}

// Validate validates the request (always valid for ledger_accept)
func (r *LedgerAcceptRequest) Validate() error {
	return nil
}

// APIVersion returns the API version
func (r *LedgerAcceptRequest) APIVersion() int {
	return r.apiVersion
}

// SetAPIVersion sets the API version
func (r *LedgerAcceptRequest) SetAPIVersion(version int) {
	r.apiVersion = version
}

// LedgerAcceptResponse represents the response from ledger_accept
type LedgerAcceptResponse struct {
	LedgerCurrentIndex uint64 `json:"ledger_current_index"`
	Status             string `json:"status"`
}

// LedgerServiceStatus represents the current status of the ledger service
type LedgerServiceStatus struct {
	Running         bool
	Interval        time.Duration
	LedgersAdvanced uint64
	LastLedgerIndex uint64
	LastError       string
}

// LedgerService manages automatic ledger advancement for local XRPL nodes
// In standalone/local mode, XRPL nodes don't automatically close ledgers.
// This service uses the ledger_accept command to manually advance ledgers.
type LedgerService struct {
	client   *rpc.Client
	interval time.Duration

	mu              sync.RWMutex
	running         bool
	cancel          context.CancelFunc
	ledgersAdvanced uint64
	lastLedgerIndex uint64
	lastError       string
}

// NewLedgerService creates a new ledger service
func NewLedgerService(rpcURL string, interval time.Duration) (*LedgerService, error) {
	cfg, err := rpc.NewClientConfig(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client config: %w", err)
	}

	client := rpc.NewClient(cfg)

	return &LedgerService{
		client:   client,
		interval: interval,
	}, nil
}

// Start begins the ledger advancement loop in the background
func (s *LedgerService) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("ledger service is already running")
	}

	serviceCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.running = true
	s.ledgersAdvanced = 0
	s.lastError = ""
	s.mu.Unlock()

	// Start the background goroutine
	go s.run(serviceCtx)

	return nil
}

// Stop gracefully stops the ledger service
func (s *LedgerService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.running = false
}

// GetStatus returns the current status of the ledger service
func (s *LedgerService) GetStatus() LedgerServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return LedgerServiceStatus{
		Running:         s.running,
		Interval:        s.interval,
		LedgersAdvanced: s.ledgersAdvanced,
		LastLedgerIndex: s.lastLedgerIndex,
		LastError:       s.lastError,
	}
}

// run is the main loop that advances ledgers at the configured interval
func (s *LedgerService) run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Advance first ledger immediately
	s.advanceLedger()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.advanceLedger()
		}
	}
}

// advanceLedger sends a ledger_accept request to the node
func (s *LedgerService) advanceLedger() {
	req := &LedgerAcceptRequest{}

	resp, err := s.client.Request(req)
	if err != nil {
		s.mu.Lock()
		s.lastError = err.Error()
		s.mu.Unlock()
		return
	}

	var result LedgerAcceptResponse
	if err := resp.GetResult(&result); err != nil {
		s.mu.Lock()
		s.lastError = err.Error()
		s.mu.Unlock()
		return
	}

	s.mu.Lock()
	atomic.AddUint64(&s.ledgersAdvanced, 1)
	s.lastLedgerIndex = result.LedgerCurrentIndex
	s.lastError = ""
	s.mu.Unlock()
}

// WaitForReady waits for the node to be ready to accept connections
func (s *LedgerService) WaitForReady(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Try to make a simple request to check if node is ready
		req := &LedgerAcceptRequest{}
		_, err := s.client.Request(req)
		if err == nil {
			return nil
		}

		// Wait a bit before retrying
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for node to be ready")
}
