---
description: A CLI tool for developing, deploying, and interacting with XRPL smart contracts written in Rust. Think Foundry, but for XRPL.
---

# Introduction to Bedrock

<DownloadLLMsFullDoc />

Bedrock is a developer tool for building, deploying, and interacting with XRPL smart contracts written in Rust. Think **Foundry**, but for XRPL.

## What is Bedrock?

Bedrock provides a complete CLI workflow for XRPL smart contract development. It compiles Rust contracts to WebAssembly and handles deployment to XRPL networks. It includes:

- **Build System** - Compile Rust smart contracts to optimized WebAssembly
- **Smart Deployment** - Auto-build, ABI generation, and deployment in one command
- **Contract Interaction** - Call deployed contract functions with typed parameters
- **Local Node** - Manage a local XRPL test network via Docker
- **ABI Generation** - Automatic ABI extraction from Rust code annotations
- **Wallet Management** - Encrypted wallet storage with Jade

## Why Use Bedrock?

Building and deploying XRPL smart contracts involves multiple tools and manual steps. Bedrock abstracts away the complexity of:

- **WASM compilation** - Sensible defaults for Rust-to-WASM compilation
- **ABI management** - Auto-generated from code annotations, no manual maintenance
- **Deployment orchestration** - One command to build, generate ABI, and deploy
- **Network configuration** - Pre-configured for local and alphanet environments
- **Wallet security** - Encrypted storage so seeds don't leak into shell history

## Key Features

### Build System

Compile Rust contracts to optimized WebAssembly with a single command. Release builds produce compact WASM files (~156 KB) ready for deployment.

### Smart Deployment

`bedrock deploy` automatically builds your contract, generates the ABI, and deploys to the network. No manual steps required.

### Local Development

Spin up a local XRPL node in Docker for fast iteration. Build, deploy, and test your contracts without waiting for testnet confirmations.

### Automatic ABI Generation

Annotate your Rust functions with JSDoc-style comments and Bedrock extracts the ABI automatically. No separate ABI files to maintain.

### Wallet Management

Create, import, and manage XRPL wallets with AES-256-GCM encryption. Your seeds are never stored in plaintext.

## Architecture Overview

Bedrock uses a hybrid architecture with Go for the CLI and embedded JavaScript modules for XRPL transaction handling:

```
bedrock CLI (Go)
       |
       ├── Build System ──────── cargo (Rust → WASM)
       |
       ├── ABI Generator ─────── Parses Rust annotations
       |
       ├── Deployer ──────────── Embedded JS (deploy.js)
       |
       ├── Caller ────────────── Embedded JS (call.js)
       |
       ├── Local Node ────────── Docker (rippled)
       |
       └── Wallet Manager ────── AES-256-GCM encryption
```

## Supported Networks

| Network | WebSocket | Faucet |
|---------|-----------|--------|
| Local | `ws://localhost:6006` | `http://localhost:8080/faucet` |
| Alphanet | `wss://alphanet.nerdnest.xyz` | `https://alphanet.faucet.nerdnest.xyz/accounts` |

## Requirements

Before installing Bedrock, ensure you have:

- **[Go](https://go.dev/dl/)** (1.21 or later) - For building Bedrock from source
- **[Node.js](https://nodejs.org/)** (18 or later) - For XRPL transaction handling
- **[Rust](https://rustup.rs/)** - For compiling smart contracts
- **[Docker](https://www.docker.com/)** (optional) - For running a local XRPL node

## Next Steps

Ready to get started? Here's the recommended path:

1. **[Getting Started](/guide/getting-started)** - Install Bedrock and create your first project
2. **[Building Contracts](/guide/building-contracts)** - Understand the build system
3. **[ABI Generation](/guide/abi-generation)** - Learn the annotation syntax
4. **[Deploying & Calling](/guide/deployment-and-calling)** - Deploy and interact with contracts
5. **[Local Node](/guide/local-node)** - Set up a local development environment
6. **[Wallet Management](/guide/wallet)** - Manage your XRPL wallets securely
7. **[Commands Reference](/guide/commands-reference)** - Complete CLI reference

## Community & Support

- **GitHub** - [XRPL-Commons/Bedrock](https://github.com/XRPL-Commons/bedrock)
- **Issues** - Report bugs or request features on GitHub
- **XRPL Commons** - [xrpl-commons.org](https://www.xrpl-commons.org)
- **XRPL Docs** - [xrpl.org](https://xrpl.org/)

## License

MIT License - See the [LICENSE](https://github.com/XRPL-Commons/bedrock/blob/main/LICENSE) file for details.
