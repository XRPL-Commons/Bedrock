package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/config"
	"github.com/xrpl-bedrock/bedrock/pkg/network"
)

var basaltCmd = &cobra.Command{
	Use:   "basalt",
	Short: "Manage local XRPL network node",
	Long:  `Start, stop, and manage a local XRPL node using Docker.`,
}

var basaltStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the local XRPL node",
	RunE:  runBasaltStart,
}

var basaltStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the local XRPL node",
	RunE:  runBasaltStop,
}

var basaltStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the local XRPL node",
	RunE:  runBasaltStatus,
}

func init() {
	rootCmd.AddCommand(basaltCmd)
	basaltCmd.AddCommand(basaltStartCmd)
	basaltCmd.AddCommand(basaltStopCmd)
	basaltCmd.AddCommand(basaltStatusCmd)
}

func runBasaltStart(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadFromWorkingDir()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'bedrock init' first)", err)
	}

	color.Cyan("  Starting local XRPL node...\n")
	color.White(" Docker image: %s\n", cfg.LocalNode.DockerImage)
	color.White(" Config dir: %s\n", cfg.LocalNode.ConfigDir)

	// Create network manager
	mgr, err := network.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create network manager: %w", err)
	}
	defer mgr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Start the node
	if err := mgr.Start(ctx, network.StartOptions{
		DockerImage: cfg.LocalNode.DockerImage,
		ConfigDir:   cfg.LocalNode.ConfigDir,
	}); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	color.Green("\n✓ Local node started successfully!\n")
	fmt.Println("\nNode endpoints:")
	fmt.Println("  WebSocket: ws://localhost:6006")
	fmt.Println("  RPC:       http://localhost:5005")
	fmt.Println("\nUse 'bedrock basalt status' to check node status")
	fmt.Println("Use 'bedrock basalt stop' to stop the node")

	return nil
}

func runBasaltStop(cmd *cobra.Command, args []string) error {
	color.Cyan("Stopping local XRPL node...\n")

	mgr, err := network.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create network manager: %w", err)
	}
	defer mgr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := mgr.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop node: %w", err)
	}

	color.Green("\n✓ Local node stopped successfully!\n")
	return nil
}

func runBasaltStatus(cmd *cobra.Command, args []string) error {
	mgr, err := network.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create network manager: %w", err)
	}
	defer mgr.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := mgr.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	fmt.Println("Local XRPL Node Status")
	fmt.Println("=" + "==================================")

	if status.Running {
		color.Green("Status:      Running ✓\n")
		fmt.Printf("Container:   %s\n", status.ContainerID)
		fmt.Printf("Image:       %s\n", status.Image)
		fmt.Println("Ports:")
		for _, port := range status.Ports {
			fmt.Printf("  - %s\n", port)
		}
		fmt.Println("\nEndpoints:")
		fmt.Println("  WebSocket: ws://localhost:6006")
		fmt.Println("  RPC:       http://localhost:5005")
	} else {
		color.Yellow("Status:      Stopped\n")
		fmt.Println("\nStart the node with: bedrock basalt start")
	}

	return nil
}
