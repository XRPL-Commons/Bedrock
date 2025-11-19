package cli

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xrpl-bedrock/bedrock/pkg/network"
)

var nodeCmd = &cobra.Command{
	Use:   "node <start|stop|status|logs>",
	Short: "Manage local XRPL node",
	Long: `Manage a local XRPL test node for development.

The node runs in a Docker container and provides a local network
for testing your smart contracts before deploying to testnet/mainnet.

Commands:
  start   - Start the local node
  stop    - Stop the local node
  status  - Check if node is running
  logs    - View node logs`,
	Args: cobra.ExactArgs(1),
	RunE: runNode,
}

func init() {
	rootCmd.AddCommand(nodeCmd)
}

func runNode(cmd *cobra.Command, args []string) error {
	subcommand := args[0]

	manager, err := network.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create network manager: %w", err)
	}

	ctx := cmd.Context()

	switch subcommand {
	case "start":
		return nodeStart(ctx, manager)
	case "stop":
		return nodeStop(ctx, manager)
	case "status":
		return nodeStatus(ctx, manager)
	case "logs":
		return nodeLogs(manager)
	default:
		return fmt.Errorf("unknown subcommand: %s (use: start, stop, status, logs)", subcommand)
	}
}

func nodeStart(ctx context.Context, manager *network.Manager) error {
	color.Cyan("Starting local XRPL node\n")
	fmt.Println()

	if err := manager.Start(ctx, network.StartOptions{}); err != nil {
		color.Red("✗ Failed to start node: %v\n", err)
		return err
	}

	fmt.Println()
	color.Green("✓ Local node started successfully!\n")
	fmt.Println()
	color.Cyan("Connection Details:\n")
	fmt.Println("  WebSocket URL: ws://localhost:6006")
	fmt.Println("  Faucet URL: http://localhost:8080/faucet")
	fmt.Println()
	color.Yellow("💡 Tips:\n")
	color.Yellow("   • Deploy with: bedrock deploy --network local\n")
	color.Yellow("   • View logs with: bedrock node logs\n")
	color.Yellow("   • Stop with: bedrock node stop\n")

	return nil
}

func nodeStop(ctx context.Context, manager *network.Manager) error {
	color.Cyan("Stopping local XRPL node\n")
	fmt.Println()

	if err := manager.Stop(ctx); err != nil {
		color.Red("✗ Failed to stop node: %v\n", err)
		return err
	}

	fmt.Println()
	color.Green("✓ Local node stopped\n")

	return nil
}

func nodeStatus(ctx context.Context, manager *network.Manager) error {
	color.Cyan("Checking local node status\n")
	fmt.Println()

	status, err := manager.Status(ctx)
	if err != nil {
		color.Red("✗ Failed to check status: %v\n", err)
		return err
	}

	if status.Running {
		color.Green("✓ Local node is running\n")
		fmt.Println()
		fmt.Println("  WebSocket URL: ws://localhost:6006")
		fmt.Println("  Faucet URL: http://localhost:8080/faucet")
	} else {
		color.Yellow("⊙ Local node is not running\n")
		fmt.Println()
		color.White("  Start with: bedrock node start\n")
	}

	return nil
}

func nodeLogs(manager *network.Manager) error {
	color.Cyan("Fetching local node logs (not yet implemented)\n")
	fmt.Println()

	color.Yellow("Log viewing coming soon. Use 'docker logs bedrock-xrpl-node' for now.\n")

	return nil
}
