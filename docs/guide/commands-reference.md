# Commands Reference

Complete reference for all Bedrock CLI commands with detailed options and examples.

## init

Create a new Bedrock smart contract project.

```bash
bedrock init <name>
```

| Argument | Description |
|----------|-------------|
| `name` | Name of the project directory to create |

**Creates:**

```
<name>/
├── bedrock.toml          # Project configuration
├── contract/
│   ├── Cargo.toml        # Rust package manifest
│   └── src/
│       └── lib.rs        # Smart contract boilerplate
```

```bash
bedrock init my-first-contract
cd my-first-contract
```

## build

Compile your Rust smart contract to WebAssembly.

```bash
bedrock build [flags]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--release` | `-r` | Build in release mode (optimized) | `true` |

```bash
# Release build (default, optimized)
bedrock build

# Debug build (faster compilation)
bedrock build --release=false
```

**Output:**
- Release: `contract/target/wasm32-unknown-unknown/release/<name>.wasm`
- Debug: `contract/target/wasm32-unknown-unknown/debug/<name>.wasm`

## deploy

Deploy a smart contract to an XRPL network.

```bash
bedrock deploy [flags]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--network` | `-n` | Target network (local, alphanet) | `alphanet` |
| `--wallet` | `-w` | Wallet seed for signing | Auto-generated |
| `--skip-build` | | Skip automatic contract rebuild | `false` |
| `--skip-abi` | | Skip ABI generation | `false` |
| `--abi` | `-a` | Path to ABI file | `abi.json` |
| `--algorithm` | | Cryptographic algorithm (secp256k1, ed25519) | `secp256k1` |

**Smart deployment** automatically: builds the contract, generates the ABI, and deploys to the network.

**Transaction fee:** 100 XRP (100,000,000 drops)

```bash
bedrock deploy                          # Deploy to alphanet
bedrock deploy --network local          # Deploy to local node
bedrock deploy --wallet sEd7...         # Use specific wallet
bedrock deploy --skip-build             # Skip rebuild
```

## call

Call a function on a deployed smart contract.

```bash
bedrock call <contract> <function> [flags]
```

| Argument | Description |
|----------|-------------|
| `contract` | The contract's XRPL account address (rXXX...) |
| `function` | Name of the function to call |

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--wallet` | `-w` | Wallet seed for signing (required) | - |
| `--network` | `-n` | Target network | `alphanet` |
| `--params` | `-p` | JSON string of function parameters | - |
| `--params-file` | `-f` | Path to JSON file with parameters | - |
| `--gas` | `-g` | Computation allowance | `1000000` |
| `--fee` | | Transaction fee in drops | `1000000` |
| `--abi` | `-a` | Path to ABI file | `abi.json` |
| `--algorithm` | | Cryptographic algorithm | `secp256k1` |

**Transaction fee:** 1 XRP (1,000,000 drops) by default

```bash
# Simple call
bedrock call rContract... hello --wallet sEd7...

# With inline parameters
bedrock call rContract... transfer \
  --wallet sEd7... \
  --params '{"to":"rRecipient...","amount":1000}'

# With parameters from file
bedrock call rContract... register \
  --wallet sEd7... \
  --params-file params.json

# With custom gas and fee
bedrock call rContract... expensive_op \
  --wallet sEd7... \
  --gas 5000000 \
  --fee 2000000

# On local network
bedrock call rContract... test \
  --wallet sEd7... \
  --network local
```

## node

Manage a local XRPL development node running in Docker.

```bash
bedrock node <command>
```

| Command | Description |
|---------|-------------|
| `start` | Start the local XRPL node |
| `stop` | Stop the running node |
| `status` | Check if the node is running |
| `logs` | View node container logs |

**Local node endpoints:**

| Service | URL |
|---------|-----|
| WebSocket | `ws://localhost:6006` |
| Faucet | `http://localhost:8080/faucet` |

**Requirements:** Docker must be installed and running.

```bash
bedrock node start     # Start local node
bedrock node status    # Check status
bedrock node logs      # View logs
bedrock node stop      # Stop node
```

## jade

Manage XRPL wallets with encrypted local storage.

### jade new

```bash
bedrock jade new <name> [--algorithm secp256k1|ed25519]
```

Creates a new XRPL wallet, encrypts it, and stores it locally.

### jade import

```bash
bedrock jade import <name> [--algorithm secp256k1|ed25519]
```

Imports an existing wallet from a seed. You'll be prompted to enter the seed securely.

### jade list

```bash
bedrock jade list
```

Lists all stored wallets.

### jade export

```bash
bedrock jade export <name>
```

Exports a wallet's seed and address (password required).

### jade remove

```bash
bedrock jade remove <name>
```

Permanently deletes a stored wallet.

**Storage location:** `~/.config/bedrock/wallets/<name>.json`
**Encryption:** AES-256-GCM with PBKDF2 key derivation

## faucet

Request testnet funds from the XRPL faucet.

```bash
bedrock faucet [flags]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--network` | `-n` | Target network | `alphanet` |
| `--wallet` | `-w` | Wallet seed | - |
| `--address` | `-a` | Wallet address | - |
| `--algorithm` | | Cryptographic algorithm | `secp256k1` |

If neither `--wallet` nor `--address` is provided, a new wallet is generated automatically.

```bash
bedrock faucet                         # Generate new wallet and fund it
bedrock faucet --address rMyAddr...    # Fund specific address
bedrock faucet --wallet sEd7...        # Fund existing wallet
bedrock faucet --network local         # Fund on local network
```

## clean

Remove build artifacts and cached files.

```bash
bedrock clean
```

**What it removes:**
- Extracted JavaScript modules (deploy.js, call.js, faucet.js)
- Installed npm dependencies (node_modules)
- Version tracking file

After cleaning, the next command that requires JS modules will automatically reinstall dependencies.

## Global Flags

| Flag | Description |
|------|-------------|
| `--help` | Display help for the command |
| `--version` | Display Bedrock version |
