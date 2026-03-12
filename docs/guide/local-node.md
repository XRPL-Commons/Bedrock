# Local Node

Bedrock manages local XRPL nodes using Docker, providing a fast development environment for testing smart contracts.

## Overview

Bedrock's local node functionality wraps Docker to run a local XRPL node (rippled) with:

- Pre-configured genesis ledger
- Pre-funded test accounts
- Fast block times for development
- Isolated environment

## Commands

### `bedrock node start`

Starts a local XRPL node in a Docker container.

```bash
bedrock node start
```

**What it does:**
1. Reads configuration from `bedrock.toml`
2. Pulls the Docker image (if not cached)
3. Mounts the genesis configuration from `.bedrock/node-config/`
4. Starts the rippled container
5. Exposes WebSocket and RPC endpoints

**Ports exposed:**

| Port | Service |
|------|---------|
| `6006` | WebSocket (transactions and subscriptions) |
| `5005` | JSON-RPC (queries) |
| `51235` | Peer protocol |

### `bedrock node stop`

Stops and removes the local XRPL node container.

```bash
bedrock node stop
```

This does **not** delete your genesis configuration or ledger data in `.bedrock/node-config/`.

### `bedrock node status`

Shows the current status of the local node.

```bash
bedrock node status
```

**Output (when running):**

```
Local XRPL Node Status
===================================
Status:      Running
Container:   a1b2c3d4e5f6
Image:       transia/alphanet:latest
Ports:
  - 6006->6006/tcp
  - 5005->5005/tcp
  - 51235->51235/tcp

Endpoints:
  WebSocket: ws://localhost:6006
  RPC:       http://localhost:5005
```

### `bedrock node logs`

View node container logs.

```bash
bedrock node logs
```

## Configuration

The local node reads its configuration from `bedrock.toml`:

```toml
[local_node]
config_dir = ".bedrock/node-config"
docker_image = "transia/alphanet:latest"
```

| Option | Description | Default |
|--------|-------------|---------|
| `config_dir` | Directory containing rippled configs | `.bedrock/node-config` |
| `docker_image` | Docker image to use for the node | `transia/alphanet:latest` |

### Node Configuration Files

The `config_dir` should contain:

```
.bedrock/node-config/
├── genesis.json          # Genesis ledger state
├── rippled.cfg          # (Optional) Custom rippled config
└── validators.txt       # (Optional) Validator list
```

#### genesis.json

Defines the initial ledger state, including pre-funded accounts and enabled amendments.

**Default genesis account:**
- Address: `rGWrZyQqhTp9Xu7G5Pkayo7bXjH4k4QYpf`
- Balance: 100,000,000,000 XRP

### Using a Custom Docker Image

```toml
[local_node]
docker_image = "your-registry/custom-rippled:v1.0.0"
```

**Requirements for custom images:**
- Must expose ports 6006 (WebSocket), 5005 (RPC), 51235 (Peer)
- Must accept genesis.json mounted at `/genesis.json`
- Must run rippled on startup

## Common Workflows

### Starting a Fresh Development Environment

```bash
bedrock init my-contract
cd my-contract
bedrock node start
bedrock node status
# WebSocket: ws://localhost:6006
```

### Restarting the Node

```bash
bedrock node stop
bedrock node start
```

### Connecting from Code

```javascript
const { Client } = require('@transia/xrpl');

const client = new Client('ws://localhost:6006');
await client.connect();
```

## Troubleshooting

### "Docker not found" or "failed to create Docker client"

Docker is not installed or not running:
1. Install [Docker Desktop](https://www.docker.com/products/docker-desktop)
2. Ensure Docker is running: `docker ps`
3. On Linux: Add your user to the docker group

### "Node is already running"

```bash
bedrock node stop
bedrock node start
```

Or manually remove:

```bash
docker stop bedrock-xrpl-node
docker rm bedrock-xrpl-node
```

### "genesis.json not found"

1. Ensure you ran `bedrock init` to create the project
2. Check that `.bedrock/node-config/genesis.json` exists
3. Verify `config_dir` in `bedrock.toml`

### Port Already in Use

Ports 6006 or 5005 are already bound. Stop other applications using these ports.

### Container Starts but Node Not Responding

1. Check container logs: `docker logs bedrock-xrpl-node`
2. Wait a few seconds for node initialization
3. Verify genesis.json is valid JSON
