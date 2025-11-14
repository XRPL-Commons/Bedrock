package network

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	ContainerName = "bedrock-xrpl-node"
)

// Manager handles Docker-based XRPL node
type Manager struct {
	docker *client.Client
}

// NewManager creates a new network manager
func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	return &Manager{docker: cli}, nil
}

// Start starts the local XRPL node
func (m *Manager) Start(ctx context.Context, opts StartOptions) error {
	// Check if container already exists
	existing, err := m.getContainer(ctx)
	if err == nil && existing != nil {
		return fmt.Errorf("node is already running (container: %s)", existing.ID[:12])
	}

	// Pull the Docker image
	if err := m.pullImage(ctx, opts.DockerImage); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	// Resolve absolute path for config directory
	configDir, err := filepath.Abs(opts.ConfigDir)
	if err != nil {
		return fmt.Errorf("failed to resolve config directory: %w", err)
	}

	// Check if genesis.json exists
	genesisPath := filepath.Join(configDir, "genesis.json")
	if _, err := os.Stat(genesisPath); os.IsNotExist(err) {
		return fmt.Errorf("genesis.json not found in %s", configDir)
	}

	// Configure port bindings
	portBindings := nat.PortMap{
		"6006/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "6006"}},
		"5005/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5005"}},
		"51235/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "51235"}},
	}

	// Create container
	resp, err := m.docker.ContainerCreate(ctx,
		&container.Config{
			Image: opts.DockerImage,
			ExposedPorts: nat.PortSet{
				"6006/tcp":  struct{}{},
				"5005/tcp":  struct{}{},
				"51235/tcp": struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: portBindings,
			Binds: []string{
				fmt.Sprintf("%s:/genesis.json:ro", genesisPath),
			},
			AutoRemove: false,
		},
		nil,
		nil,
		ContainerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := m.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// Stop stops the local XRPL node
func (m *Manager) Stop(ctx context.Context) error {
	containerInfo, err := m.getContainer(ctx)
	if err != nil {
		return fmt.Errorf("node is not running")
	}

	// Stop container
	timeout := 10
	if err := m.docker.ContainerStop(ctx, containerInfo.ID, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Remove container
	if err := m.docker.ContainerRemove(ctx, containerInfo.ID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// Status returns the status of the local node
func (m *Manager) Status(ctx context.Context) (*NodeStatus, error) {
	containerInfo, err := m.getContainer(ctx)
	if err != nil {
		return &NodeStatus{Running: false}, nil
	}

	return &NodeStatus{
		Running:     containerInfo.State.Running,
		ContainerID: containerInfo.ID[:12],
		Image:       containerInfo.Config.Image,
		Ports:       formatPorts(containerInfo.NetworkSettings.Ports),
	}, nil
}

// Close closes the Docker client
func (m *Manager) Close() error {
	return m.docker.Close()
}

func (m *Manager) pullImage(ctx context.Context, imageName string) error {
	reader, err := m.docker.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	// Read the pull output (required for pull to complete)
	_, err = io.Copy(io.Discard, reader)
	return err
}

func (m *Manager) getContainer(ctx context.Context) (*types.ContainerJSON, error) {
	containerJSON, err := m.docker.ContainerInspect(ctx, ContainerName)
	if err != nil {
		return nil, err
	}
	return &containerJSON, nil
}

func formatPorts(ports nat.PortMap) []string {
	var result []string
	for port, bindings := range ports {
		for _, binding := range bindings {
			result = append(result, fmt.Sprintf("%s->%s", binding.HostPort, port))
		}
	}
	return result
}
