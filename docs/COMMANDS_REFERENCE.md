# Bedrock CLI Commands Reference

Complete reference for all Bedrock CLI commands with detailed options and examples.

## Table of Contents

- [init](#init)
- [build](#build)
- [deploy](#deploy)
- [call](#call)
- [node](#node)
- [wallet](#wallet)
- [faucet](#faucet)
- [clean](#clean)

---

## init

Create a new Bedrock smart contract project.

### Usage

```bash
bedrock init <name>
```

### Arguments

| Argument | Description |
|----------|-------------|
| `name` | Name of the project directory to create |

### What it creates

```
<name>/
├── bedrock.toml          # Project configuration
├── contract/
│   ├── Cargo.toml        # Rust package manifest
│   └── src/
│       └── lib.rs        # Smart contract boilerplate
```

### Example

```bash
bedrock init my-first-contract
cd my-first-contract
```

---

## build

Compile your Rust smart contract to WebAssembly.

### Usage

```bash
bedrock build [flags]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--debug` | `-d` | Build in debug mode (faster build, larger output) |

### Output

The compiled WASM is placed in:
- Release: `target/wasm32-unknown-unknown/release/<name>.wasm`
- Debug: `target/wasm32-unknown-unknown/debug/<name>.wasm`

### Examples

```bash
# Release build (default, optimized)
bedrock build

# Debug build (faster compilation)
bedrock build --debug
```

### Notes

- Release builds are smaller and optimized for deployment
- Debug builds are faster to compile but produce larger WASM files
- Automatically adds `wasm32-unknown-unknown` target if missing

---

## deploy

Deploy a smart contract to an XRPL network.

### Usage

```bash
bedrock deploy [flags]
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--network` | `-n` | Target network (local, alphanet, testnet, mainnet) | `alphanet` |
| `--wallet` | `-w` | Wallet seed for signing transactions | Auto-generated |
| `--skip-build` | | Skip automatic contract rebuild | `false` |
| `--skip-abi` | | Skip ABI generation | `false` |
| `--algorithm` | `-a` | Wallet algorithm (secp256k1, ed25519) | `secp256k1` |

### Smart Deployment

By default, `bedrock deploy` performs:
1. Contract build (release mode)
2. ABI generation from annotations
3. Contract deployment to network

### Output

On success, displays:
- Wallet address and seed (save these!)
- Contract account address
- Transaction hash

### Examples

```bash
# Deploy to alphanet (default)
bedrock deploy

# Deploy to local node
bedrock deploy --network local

# Deploy with existing wallet
bedrock deploy --wallet sEd7...

# Deploy without rebuilding
bedrock deploy --skip-build

# Full example with all options
bedrock deploy \
  --network alphanet \
  --wallet sEd7... \
  --algorithm secp256k1
```

### Transaction Fees

- **ContractCreate fee**: 100 XRP (100,000,000 drops)
- Ensure wallet has sufficient balance before deploying

---

## call

Call a function on a deployed smart contract.

### Usage

```bash
bedrock call <contract> <function> [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `contract` | The contract's XRPL account address (rXXX...) |
| `function` | Name of the function to call |

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--wallet` | `-w` | Wallet seed for signing (required) | - |
| `--network` | `-n` | Target network | `alphanet` |
| `--params` | `-p` | JSON string of function parameters | - |
| `--params-file` | | Path to JSON file with parameters | - |
| `--gas` | `-g` | Computation allowance | `1000000` |
| `--fee` | `-f` | Transaction fee in drops | `1000000` |

### Parameter Passing

Parameters can be passed as:
1. **Inline JSON**: `--params '{"to":"rXXX...","amount":1000}'`
2. **JSON file**: `--params-file params.json`

Parameter names must match the ABI definition.

### Examples

```bash
# Simple call without parameters
bedrock call rContract123... hello --wallet sEd7...

# Call with inline parameters
bedrock call rContract123... transfer \
  --wallet sEd7... \
  --params '{"to":"rRecipient...","amount":1000}'

# Call with parameters from file
bedrock call rContract123... register \
  --wallet sEd7... \
  --params-file ./params.json

# Call with custom gas and fee
bedrock call rContract123... expensive_operation \
  --wallet sEd7... \
  --gas 5000000 \
  --fee 2000000

# Call on local network
bedrock call rContract123... test_function \
  --wallet sEd7... \
  --network local
```

### Transaction Fees

- **ContractCall fee**: 1 XRP (1,000,000 drops) by default
- Increase fee for complex operations

---

## node

Manage a local XRPL development node running in Docker.

### Usage

```bash
bedrock node <command>
```

### Commands

| Command | Description |
|---------|-------------|
| `start` | Start the local XRPL node |
| `stop` | Stop the running node |
| `status` | Check if the node is running |
| `logs` | View node container logs |

### Local Node Endpoints

| Service | URL |
|---------|-----|
| WebSocket | `ws://localhost:6006` |
| Faucet | `http://localhost:8080/faucet` |

### Examples

```bash
# Start local development node
bedrock node start

# Check if it's running
bedrock node status

# View logs
bedrock node logs

# Stop when done
bedrock node stop
```

### Requirements

- Docker must be installed and running
- Port 6006 must be available

---

## wallet

Manage XRPL wallets with encrypted local storage.

### Subcommands

### wallet new

Create a new XRPL wallet.

```bash
bedrock wallet new <name> [flags]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--algorithm` | `-a` | Cryptographic algorithm (secp256k1, ed25519) | `secp256k1` |

**Example:**
```bash
bedrock wallet new my-dev-wallet
bedrock wallet new my-ed-wallet --algorithm ed25519
```

### wallet import

Import an existing wallet from a seed.

```bash
bedrock wallet import <name> [flags]
```

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--algorithm` | `-a` | Wallet's cryptographic algorithm | `secp256k1` |

**Example:**
```bash
bedrock wallet import my-existing-wallet
# You'll be prompted to enter the seed securely
```

### wallet list

List all stored wallets.

```bash
bedrock wallet list
```

### wallet export

Export a wallet's seed and address.

```bash
bedrock wallet export <name>
```

**Example:**
```bash
bedrock wallet export my-dev-wallet
# You'll be prompted for the password
```

### wallet remove

Permanently delete a stored wallet.

```bash
bedrock wallet remove <name>
```

### Storage Location

Wallets are stored encrypted in:
- Linux/macOS: `~/.config/bedrock/wallets/<name>.json`
- Encryption: AES-256-GCM with PBKDF2 key derivation

---

## faucet

Request testnet funds for an XRPL account.

### Usage

```bash
bedrock faucet <address> [flags]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `address` | XRPL account address to fund (rXXX...) |

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--network` | `-n` | Target network | `alphanet` |

### Examples

```bash
# Fund an account on alphanet
bedrock faucet rMyAddress123...

# Fund on local network
bedrock faucet rMyAddress123... --network local
```

### Network Faucets

| Network | Faucet URL |
|---------|------------|
| Local | `http://localhost:8080/faucet` |
| Alphanet | `https://alphanet.faucet.nerdnest.xyz/accounts` |

---

## clean

Remove build artifacts and cached files.

### Usage

```bash
bedrock clean [flags]
```

### What it removes

- `target/` directory (Rust build artifacts)
- Local cache files

### Flags

| Flag | Description |
|------|-------------|
| `--all` | Also clean global module cache (~/.cache/bedrock) |

### Examples

```bash
# Clean project build artifacts
bedrock clean

# Clean everything including global cache
bedrock clean --all
```

---

## Global Flags

These flags are available for all commands:

| Flag | Description |
|------|-------------|
| `--help` | Display help for the command |
| `--version` | Display Bedrock version |

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BEDROCK_CACHE_DIR` | Override default cache directory |
| `BEDROCK_CONFIG_DIR` | Override default config directory |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | Network error |
| 4 | Build error |
