package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Peersyst/xrpl-go/xrpl/rpc"
	"github.com/spf13/cobra"
)

// LedgerAcceptRequest represents a ledger_accept RPC request
type LedgerAcceptRequest struct {
	apiVersion int
}

func (r *LedgerAcceptRequest) Method() string      { return "ledger_accept" }
func (r *LedgerAcceptRequest) Validate() error     { return nil }
func (r *LedgerAcceptRequest) APIVersion() int     { return r.apiVersion }
func (r *LedgerAcceptRequest) SetAPIVersion(v int) { r.apiVersion = v }

var ledgerDaemonCmd = &cobra.Command{
	Use:    "_ledger-daemon",
	Short:  "Internal: Run ledger advancement daemon",
	Hidden: true,
	RunE:   runLedgerDaemon,
}

var (
	daemonRPCURL   string
	daemonInterval int
)

func init() {
	RootCmd.AddCommand(ledgerDaemonCmd)
	ledgerDaemonCmd.Flags().StringVar(&daemonRPCURL, "rpc-url", "http://localhost:5005", "RPC URL")
	ledgerDaemonCmd.Flags().IntVar(&daemonInterval, "interval", 1000, "Interval in milliseconds")
}

func runLedgerDaemon(cmd *cobra.Command, args []string) error {
	interval := time.Duration(daemonInterval) * time.Millisecond

	cfg, err := rpc.NewClientConfig(daemonRPCURL)
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}

	client := rpc.NewClient(cfg)

	// Wait for node to be ready
	fmt.Fprintf(os.Stderr, "Waiting for node to be ready...\n")
	ready := false
	for i := 0; i < 60; i++ { // 30 second timeout
		req := &LedgerAcceptRequest{}
		_, err := client.Request(req)
		if err == nil {
			ready = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		return fmt.Errorf("timeout waiting for node to be ready")
	}

	fmt.Fprintf(os.Stderr, "Ledger daemon started (interval: %v)\n", interval)

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Fprintf(os.Stderr, "\nLedger daemon stopping...\n")
		cancel()
	}()

	// Run the ledger advancement loop
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ledgerCount := uint64(0)

	// Advance first ledger immediately
	advanceLedger(client, &ledgerCount)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(os.Stderr, "Ledger daemon stopped (%d ledgers advanced)\n", ledgerCount)
			return nil
		case <-ticker.C:
			advanceLedger(client, &ledgerCount)
		}
	}
}

func advanceLedger(client *rpc.Client, count *uint64) {
	req := &LedgerAcceptRequest{}
	_, err := client.Request(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to advance ledger: %v\n", err)
		return
	}
	*count++
}
