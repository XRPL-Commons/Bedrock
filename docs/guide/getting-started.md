# Getting Started

This guide walks you through installing Bedrock and creating your first XRPL smart contract project.

## Prerequisites

Before installing Bedrock, ensure you have the following tools:

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.21+ | Building Bedrock from source |
| [Node.js](https://nodejs.org/) | 18+ | XRPL transaction handling |
| [Rust](https://rustup.rs/) | 1.70+ | Compiling smart contracts |
| [Docker](https://www.docker.com/) | Latest | Local XRPL node (optional) |

### Verify Prerequisites

```bash
go version      # Should show 1.21+
node --version  # Should show v18+
rustc --version # Should show 1.70+
cargo --version
```

### Installing Rust

If you don't have Rust installed:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

## Installation

### Install from Source

```bash
# Clone the repository
git clone https://github.com/XRPL-Commons/Bedrock.git
cd Bedrock

# Build and install
go build -o bedrock cmd/bedrock/main.go
sudo mv bedrock /usr/local/bin/

# Verify installation
bedrock --help
```

### First Run

On first run, Bedrock automatically installs JavaScript dependencies:

```
First run detected - installing JavaScript dependencies...
```

This takes ~10-15 seconds and only happens once. Dependencies are cached in `~/.cache/bedrock/modules/`.

## Create Your First Project

### 1. Initialize a Project

```bash
bedrock init my-contract
cd my-contract
```

This creates the following structure:

```
my-contract/
├── bedrock.toml          # Project configuration
├── contract/
│   ├── Cargo.toml        # Rust package manifest
│   └── src/
│       └── lib.rs        # Smart contract boilerplate
```

### 2. Explore the Contract

The generated `contract/src/lib.rs` contains a basic contract:

```rust
#![cfg_attr(target_arch = "wasm32", no_std)]

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
```

### 3. Build the Contract

```bash
bedrock build
```

This compiles your Rust contract to optimized WebAssembly:

```
Building smart contract...
   Mode: Release (optimized)
   Source: contract/src/lib.rs

✓ Build completed successfully!

Output:   contract/target/wasm32-unknown-unknown/release/my_contract.wasm
Size:     156.4 KB
Duration: 5.1s
```

### 4. Start a Local Node (Optional)

For local development, start a Docker-based XRPL node:

```bash
bedrock node start
```

This exposes:
- WebSocket: `ws://localhost:6006`
- Faucet: `http://localhost:8080/faucet`

### 5. Deploy Your Contract

```bash
# Deploy to local node
bedrock deploy --network local

# Or deploy to alphanet (testnet)
bedrock deploy --network alphanet
```

The deploy command automatically:
1. Builds the contract (if needed)
2. Generates the ABI (if needed)
3. Deploys to the network

Save the output - it contains your **wallet seed** and **contract address**.

### 6. Call Your Contract

```bash
bedrock call <contract-address> hello \
  --wallet <seed> \
  --network local
```

## Project Configuration

The `bedrock.toml` file configures your project:

```toml
[project]
name = "my-contract"
version = "0.1.0"

[build]
source = "contract/src/lib.rs"
target = "wasm32-unknown-unknown"

[local_node]
config_dir = ".bedrock/node-config"
docker_image = "transia/alphanet:latest"

[networks.local]
url = "ws://localhost:6006"
network_id = 63456
faucet_url = "http://localhost:8080/faucet"

[networks.alphanet]
url = "wss://alphanet.nerdnest.xyz"
network_id = 21465
faucet_url = "https://alphanet.faucet.nerdnest.xyz/accounts"
```

## Development Workflow

A typical development loop looks like:

```bash
# Terminal 1: Start local node
bedrock node start

# Terminal 2: Develop and test
bedrock build                                         # Build contract
bedrock deploy --network local                        # Deploy to local node
bedrock call <contract> hello --wallet <seed> --network local  # Test

# Make changes to contract/src/lib.rs...
bedrock deploy --network local                        # Redeploy (auto-rebuilds)
```

When ready for testnet:

```bash
bedrock deploy --network alphanet
```

## Next Steps

- [Building Contracts](/guide/building-contracts) - Deep dive into the build system
- [ABI Generation](/guide/abi-generation) - Learn the annotation syntax
- [Deploying & Calling](/guide/deployment-and-calling) - Full deployment guide
- [Local Node](/guide/local-node) - Configure your local environment
- [Wallet Management](/guide/wallet) - Secure wallet handling
