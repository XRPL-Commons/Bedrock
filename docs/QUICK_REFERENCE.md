# Bedrock Quick Reference

A quick lookup guide for common Bedrock operations.

## Commands at a Glance

| Command | Purpose |
|---------|---------|
| `bedrock init <name>` | Create new project |
| `bedrock build` | Compile contract to WASM |
| `bedrock deploy` | Deploy contract |
| `bedrock call <addr> <fn>` | Call contract function |
| `bedrock node start/stop` | Manage local node |
| `bedrock jade new <name>` | Create wallet |
| `bedrock faucet` | Get testnet funds |
| `bedrock clean` | Clean build artifacts |

## Quick Start

```bash
bedrock init my-app && cd my-app
bedrock node start
bedrock deploy --network local
bedrock call <contract> hello --wallet <seed> --network local
```

## Networks

| Network | WebSocket | Faucet |
|---------|-----------|--------|
| Local | `ws://localhost:6006` | `http://localhost:8080/faucet` |
| Alphanet | `wss://alphanet.nerdnest.xyz` | `https://alphanet.faucet.nerdnest.xyz/accounts` |

## Deploy Options

```bash
bedrock deploy                      # Default (alphanet, auto-build)
bedrock deploy --network local      # Local node
bedrock deploy --wallet <seed>      # Specific wallet
bedrock deploy --skip-build         # Skip rebuild
```

## Call Options

```bash
bedrock call <contract> <function> --wallet <seed>
  --params '{"key":"value"}'        # Inline JSON
  --params-file params.json         # From file
  --gas 1000000                     # Computation limit
  --network alphanet                # Target network
```

## Wallet Commands (Jade)

```bash
bedrock jade new <name>             # Create encrypted wallet
bedrock jade import <name>          # Import from seed
bedrock jade list                   # List all wallets
bedrock jade export <name>          # Show seed
bedrock jade remove <name>          # Delete wallet
```

## ABI Annotations

```rust
/// @xrpl-function function_name
/// @param name TYPE - description
/// @return TYPE - description
/// @flag 0  // required (default)
/// @flag 1  // optional
```

## XRPL Types

| Type | Use for |
|------|---------|
| `UINT8/16/32/64/128/256` | Integers |
| `VL` | Bytes/strings |
| `ACCOUNT` | Addresses |
| `AMOUNT` | XRP/token amounts |
| `CURRENCY` | Currency codes |
| `ISSUE` | Currency+issuer |

## Fees

| Operation | Fee |
|-----------|-----|
| Deploy (ContractCreate) | 100 XRP |
| Call (ContractCall) | 1 XRP |

## Contract Template

```rust
#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function my_func
#[wasm_export]
fn my_func() -> i32 {
    let _ = trace("Hello from XRPL Smart Contract!");
    0
}
```

## File Structure

```
project/
├── bedrock.toml      # Config
├── contract/src/lib.rs  # Contract code
├── abi.json          # Generated ABI
└── target/wasm32.../release/*.wasm
```

## Troubleshooting

| Problem | Solution |
|---------|----------|
| WASM target missing | `rustup target add wasm32-unknown-unknown` |
| Node won't start | Check Docker: `docker ps` |
| Deployment fails | Ensure 100+ XRP balance |
| Modules not found | Check Node.js 18+: `node --version` |

## Useful Paths

| Path | Contents |
|------|----------|
| `~/.config/bedrock/wallets/` | Encrypted wallets |
| `~/.cache/bedrock/modules/` | JS module cache |
