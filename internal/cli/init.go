package cli

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new XRPL smart contract project",
	Long:  `Creates a new XRPL smart contract project with the necessary structure and configuration files.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInit,
}

//go:embed templates/genesis.json
var genesisTemplate string

func init() {
	RootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	color.Cyan("Initializing project: %s\n", projectName)

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create directory structure
	dirs := []string{
		filepath.Join(projectName, "contract", "src"),
		filepath.Join(projectName, ".wallets"),
		filepath.Join(projectName, ".bedrock", "node-config"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create bedrock.toml
	configContent := fmt.Sprintf(`[project]
name = "%s"
version = "0.1.0"
authors = ["Your Name"]

[build]
source = "contract/src/lib.rs"
output = "contract/target/wasm32-unknown-unknown/release"
target = "wasm32-unknown-unknown"

[local_node]
config_dir = ".bedrock/node-config"
docker_image = "transia/alphanet:latest"
ledger_interval = 1000

[networks.local]
url = "ws://localhost:6006"
network_id = 63456
faucet_url = "http://localhost:8080/faucet"

[networks.alphanet]
url = "wss://alphanet.nerdnest.xyz"
network_id = 21465
faucet_url = "https://alphanet.faucet.nerdnest.xyz/accounts"

[contracts.main]
source = "contract/src/lib.rs"
abi = "contract/build/abi.json"

[wallets]
keystore = ".wallets/keystore.json"
`, projectName)

	if err := os.WriteFile(filepath.Join(projectName, "bedrock.toml"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create bedrock.toml: %w", err)
	}

	// Create Cargo.toml
	cargoContent := fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]

[dependencies]
xrpl-wasm-std = { git = "https://github.com/Transia-RnD/craft.git", branch = "dangell/smart-contracts", package = "xrpl-wasm-std" }
xrpl-wasm-macros= { git = "https://github.com/Transia-RnD/craft.git", branch = "dangell/smart-contracts", package = "xrpl-wasm-macros" }

[profile.release]
opt-level = "z"
lto = true
strip = true
panic = "abort"
`, projectName)

	if err := os.WriteFile(filepath.Join(projectName, "contract", "Cargo.toml"), []byte(cargoContent), 0644); err != nil {
		return fmt.Errorf("failed to create Cargo.toml: %w", err)
	}

	// Create lib.rs
	libContent := `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function hello
#[wasm_export]
fn hello() -> i32 {
    let _ = trace("Hello from XRPL Smart Contract!");
    0
}
`

	if err := os.WriteFile(filepath.Join(projectName, "contract", "src", "lib.rs"), []byte(libContent), 0644); err != nil {
		return fmt.Errorf("failed to create lib.rs: %w", err)
	}

	// Write embedded genesis.json template
	genesisPath := filepath.Join(projectName, ".bedrock", "node-config", "genesis.json")
	if err := os.WriteFile(genesisPath, []byte(genesisTemplate), 0644); err != nil {
		return fmt.Errorf("failed to create genesis.json: %w", err)
	}

	// Create .gitignore
	gitignoreContent := `# Build outputs
contract/target/
*.wasm

# Wallets (keep private!)
.wallets/

# Bedrock internal
.bedrock/

# OS files
.DS_Store
`

	if err := os.WriteFile(filepath.Join(projectName, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Create README.md
	readmeContent := fmt.Sprintf("# %s\n\nXRPL Smart Contract project\n\n## Getting Started\n\n### Build the contract\n```bash\nbedrock flint build --release\n```\n\n### Start local node\n```bash\nbedrock basalt start\n```\n\n### Deploy\n```bash\nbedrock slate deploy --network local\n```\n\n## Project Structure\n\n- `contract/` - Smart contract source code\n- `bedrock.toml` - Project configuration\n- `.wallets/` - Local wallet storage (git-ignored)\n", projectName)

	if err := os.WriteFile(filepath.Join(projectName, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Success message
	color.Green("\n✓ Project initialized successfully!\n\n")
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  bedrock flint build      # Build your contract")
	fmt.Println("  bedrock basalt start     # Start local node")
	fmt.Println("  bedrock slate deploy     # Deploy your contract")

	return nil
}
