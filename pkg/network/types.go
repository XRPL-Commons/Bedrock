package network

// StartOptions configures how to start the local node
type StartOptions struct {
	DockerImage string
	ConfigDir   string
}

// NodeStatus represents the current status of the node
type NodeStatus struct {
	Running     bool
	ContainerID string
	Image       string
	Ports       []string
}
