package network

import "time"

// StartOptions configures how to start the local node
type StartOptions struct {
	DockerImage    string
	ConfigDir      string
	LedgerInterval time.Duration // Interval for ledger advancement
	RPCURL         string        // RPC URL for ledger service
}

// NodeStatus represents the current status of the node
type NodeStatus struct {
	Running     bool
	ContainerID string
	Image       string
	Ports       []string
	// Ledger service status
	LedgerServiceRunning bool
	LedgersAdvanced      uint64
	LastLedgerIndex      uint64
}
