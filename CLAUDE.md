# Bedrock CLI Instructions

This file provides instructions for AI assistants working with the Bedrock XRPL smart contract development tool.

## What is Bedrock?

Bedrock is a CLI tool for developing, deploying, and interacting with XRPL smart contracts written in Rust. It compiles contracts to WebAssembly and handles deployment to XRPL networks.

## Command Reference

### Project Initialization

```bash
# Create a new project
bedrock init <project-name>
```

Creates a project structure with `bedrock.toml`, `contract/` directory, and boilerplate Rust code.

### Building Contracts

```bash
# Build in release mode (optimized, smaller WASM)
bedrock build

# Build in debug mode (faster compilation)
bedrock build --debug
```

### Deploying Contracts

```bash
# Deploy to alphanet (default)
bedrock deploy

# Deploy to local node
bedrock deploy --network local

# Deploy with specific wallet
bedrock deploy --wallet <seed>

# Skip auto-build or ABI generation
bedrock deploy --skip-build
bedrock deploy --skip-abi
```

**Important deployment details:**
- Deployment fee: 100 XRP (100000000 drops)
- Auto-generates ABI from Rust annotations
- Creates new wallet if none specified

### Calling Contract Functions

```bash
# Basic call
bedrock call <contract-address> <function-name> --wallet <seed>

# With JSON parameters
bedrock call <contract-address> <function-name> \
  --wallet <seed> \
  --params '{"param1": "value1", "param2": 123}'

# Parameters from file
bedrock call <contract-address> <function-name> \
  --wallet <seed> \
  --params-file params.json

# With custom gas and network
bedrock call <contract-address> <function-name> \
  --wallet <seed> \
  --gas 1000000 \
  --network alphanet
```

**Important call details:**
- Call fee: 1 XRP (1000000 drops)
- Function names are hex-encoded automatically
- `--wallet` is required

### Local Node Management

```bash
# Start local XRPL node (Docker required)
bedrock node start

# Stop the node
bedrock node stop

# Check status
bedrock node status

# View logs
bedrock node logs
```

**Local node endpoints:**
- WebSocket: `ws://localhost:6006`
- Faucet: `http://localhost:8080/faucet`

### Wallet Management

```bash
# Create new wallet
bedrock wallet new <name>
bedrock wallet new <name> --algorithm ed25519

# Import existing wallet
bedrock wallet import <name>

# List wallets
bedrock wallet list

# Export wallet (shows seed)
bedrock wallet export <name>

# Remove wallet
bedrock wallet remove <name>
```

Wallets are encrypted and stored in `~/.config/bedrock/wallets/`.

### Other Commands

```bash
# Request testnet funds
bedrock faucet <address>

# Clean build artifacts
bedrock clean
```

## Network Information

### Alphanet (Testnet)
- WebSocket: `wss://alphanet.nerdnest.xyz`
- Faucet: `https://alphanet.faucet.nerdnest.xyz/accounts`
- Network ID: 21465

### Local Node
- WebSocket: `ws://localhost:6006`
- Faucet: `http://localhost:8080/faucet`

## Smart Contract Development

### Contract Structure

```rust
#![cfg_attr(target_arch = "wasm32", no_std)]

use xrpl_wasm::*;

/// @xrpl-function my_function
/// @param input VL - Input data
/// @return UINT64 - Result code
#[no_mangle]
pub extern "C" fn my_function(input: &[u8]) -> u64 {
    // Implementation
    0
}

#[cfg(target_arch = "wasm32")]
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
```

### ABI Annotation Syntax

```rust
/// @xrpl-function <function_name>
/// @param <name> <TYPE> - <description>
/// @return <TYPE> - <description>
/// @flag 0  (required) or @flag 1 (optional)
```

### Supported Types

| Type | Description |
|------|-------------|
| `UINT8`, `UINT16`, `UINT32`, `UINT64`, `UINT128`, `UINT256` | Unsigned integers |
| `VL` | Variable length bytes/string |
| `ACCOUNT` | XRPL account address |
| `AMOUNT` | XRP or token amount |
| `CURRENCY` | Currency code |
| `ISSUE` | Currency + issuer pair |

## Typical Development Workflow

1. **Initialize project:**
   ```bash
   bedrock init my-contract
   cd my-contract
   ```

2. **Start local node (optional):**
   ```bash
   bedrock node start
   ```

3. **Write contract** in `contract/src/lib.rs`

4. **Build:**
   ```bash
   bedrock build
   ```

5. **Deploy:**
   ```bash
   bedrock deploy --network local
   # Note the contract address and wallet seed from output
   ```

6. **Interact:**
   ```bash
   bedrock call <contract> <function> --wallet <seed> --network local
   ```

7. **Deploy to testnet:**
   ```bash
   bedrock deploy --network alphanet
   ```

## Configuration File (bedrock.toml)

```toml
[project]
name = "my-contract"
version = "0.1.0"

[build]
source = "contract/src/lib.rs"
target = "wasm32-unknown-unknown"

[networks.local]
url = "ws://localhost:6006"
faucet_url = "http://localhost:8080/faucet"

[networks.alphanet]
url = "wss://alphanet.nerdnest.xyz"
faucet_url = "https://alphanet.faucet.nerdnest.xyz/accounts"
```

## Common Issues and Solutions

### Build fails with wasm32 target error
```bash
rustup target add wasm32-unknown-unknown
```

### Node won't start
Ensure Docker is running: `docker ps`

### Module installation fails
Check Node.js version (18+ required): `node --version`

### Contract deployment fails
- Ensure wallet has sufficient XRP (100+ XRP for deployment)
- Check network connectivity to alphanet

## Project File Structure

```
my-contract/
├── bedrock.toml          # Project configuration
├── contract/
│   ├── Cargo.toml        # Rust dependencies
│   └── src/
│       └── lib.rs        # Smart contract code
├── abi.json              # Generated ABI (after deploy)
└── target/               # Build output
    └── wasm32-unknown-unknown/
        └── release/
            └── contract.wasm
```

## Best Practices

1. Always use release mode for production deployments
2. Test on local node before deploying to alphanet
3. Keep wallet seeds secure - use `bedrock wallet` for encrypted storage
4. Include descriptive ABI annotations for all exported functions
5. Optimize contracts with `opt-level = "z"` and `lto = true` in Cargo.toml
