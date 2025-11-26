# Local Node Management

**Bedrock** manages local XRPL nodes using Docker, providing a fast development environment for testing smart contracts.

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

**Output:**

```
   Starting local XRPL node...
   Docker image: transia/alphanet:latest
   Config dir: .bedrock/node-config

✓ Local node started successfully!

Node endpoints:
  WebSocket: ws://localhost:6006
  RPC:       http://localhost:5005

Use 'bedrock node status' to check node status
Use 'bedrock node stop' to stop the node
```

**Requirements:**

- Docker must be running
- Must be in a bedrock project directory (with `bedrock.toml`)

**Ports exposed:**

- `6006` - WebSocket (for transactions and subscriptions)
- `5005` - JSON-RPC (for queries)
- `51235` - Peer protocol

---

### `bedrock node stop`

Stops and removes the local XRPL node container.

```bash
bedrock node stop
```

**What it does:**

1. Gracefully stops the running container
2. Removes the container

**Output:**

```
Stopping local XRPL node...

✓ Local node stopped successfully!
```

**Note:** This does **not** delete your genesis configuration or ledger data mounted from `.bedrock/node-config/`.

---

### `bedrock node status`

Shows the current status of the local node.

```bash
bedrock node status
```

**Output (when running):**

```
Local XRPL Node Status
===================================
Status:      Running ✓
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

**Output (when stopped):**

```
Local XRPL Node Status
===================================
Status:      Stopped

Start the node with: bedrock node start
```

---

## Configuration

The local node reads its configuration from `bedrock.toml` in your project root:

```toml
[local_node]
config_dir = ".bedrock/node-config"
docker_image = "transia/alphanet:latest"
```

### Configuration Options

| Option         | Description                                               | Default                   |
| -------------- | --------------------------------------------------------- | ------------------------- |
| `config_dir`   | Directory containing rippled configs (genesis.json, etc.) | `.bedrock/node-config`    |
| `docker_image` | Docker image to use for the node                          | `transia/alphanet:latest` |

### Node Configuration Files

The `config_dir` should contain:

```
.bedrock/node-config/
├── genesis.json          # Genesis ledger state
├── rippled.cfg          # (Optional) Custom rippled config
└── validators.txt       # (Optional) Validator list
```

#### genesis.json

Defines the initial ledger state, including:

- Pre-funded accounts
- Enabled amendments
- Initial balances

**Default genesis account:**

- Address: `rGWrZyQqhTp9Xu7G5Pkayo7bXjH4k4QYpf`
- Balance: 100,000,000,000 XRP (100 billion drops)

**Example:**

```json
{
  "ledger": {
    "accepted": true,
    "accountState": [
      {
        "Account": "rGWrZyQqhTp9Xu7G5Pkayo7bXjH4k4QYpf",
        "Balance": "100000000000000000",
        "Flags": 0,
        "LedgerEntryType": "AccountRoot",
        "OwnerCount": 0,
        "Sequence": 1,
        "index": "2B6AC232AA4C4BE41BF49D2459FA4A0347E1B543A4C92FCEE0821C0201E2E9A8"
      }
    ],
    "ledger_index": "1",
    "parent_hash": "",
    "total_coins": "100000000000000000"
  }
}
```

---

## Docker Image

The local node uses Docker images that contain a pre-configured rippled binary. By default, it uses `transia/alphanet:latest`.

### Using a Custom Image

You can specify a different Docker image in `bedrock.toml`:

```toml
[local_node]
docker_image = "your-registry/custom-rippled:v1.0.0"
```

**Requirements for custom images:**

- Must expose ports 6006 (WebSocket), 5005 (RPC), 51235 (Peer)
- Must accept genesis.json mounted at `/genesis.json`
- Must run rippled on startup

---

## Common Workflows

### Starting a Fresh Development Environment

```bash
# Initialize project
bedrock init my-contract
cd my-contract

# Start local node
bedrock node start

# Verify it's running
bedrock node status

# Use the node
# WebSocket: ws://localhost:6006
```

### Restarting the Node

```bash
# Stop current instance
bedrock node stop

# Start fresh
bedrock node start
```

### Connecting from Your Code

**JavaScript (@transia/xrpl):**

```javascript
const { Client } = require('@transia/xrpl');

const client = new Client('ws://localhost:6006');
await client.connect();
```

**Using the local network config:**

```javascript
// In bedrock.toml:
// [networks.local]
// url = "ws://localhost:6006"
// network_id = 63456

const cfg = loadConfig('bedrock.toml');
const client = new Client(cfg.networks.local.url);
```

---

## Troubleshooting

### "Docker not found" or "failed to create Docker client"

**Problem:** Docker is not installed or not running.

**Solution:**

1. Install Docker Desktop: https://www.docker.com/products/docker-desktop
2. Ensure Docker is running: `docker ps` should work
3. Verify Docker permissions (Linux): Add user to docker group

### "Node is already running"

**Problem:** Container already exists.

**Solution:**

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

**Problem:** Configuration directory is missing or empty.

**Solution:**

1. Ensure you ran `bedrock init` to create the project
2. Check that `.bedrock/node-config/genesis.json` exists
3. Verify `config_dir` in `bedrock.toml` points to the right location

### Port Already in Use

**Problem:** Ports 6006 or 5005 are already bound.

**Solution:**

1. Stop other XRPL nodes or applications using these ports
2. Or modify the Docker port mappings in the future (not yet supported in MVP)

### Container Starts but Node Not Responding

**Problem:** Node is running but not accepting connections.

**Solution:**

1. Check container logs:
   ```bash
   docker logs bedrock-xrpl-node
   ```
2. Wait a few seconds for node initialization
3. Verify genesis.json is valid JSON

---

## Technical Details

### Container Management

The local node functionality manages a Docker container named `bedrock-xrpl-node`:

- Single instance per machine
- Auto-remove is disabled (allows restart)
- Mounts genesis.json as read-only

### Network Mode

The container runs in bridge network mode with port mappings:

```
Host          Container
6006    →     6006      (WebSocket)
5005    →     5005      (JSON-RPC)
51235   →     51235     (Peer)
```

### Volume Mounts

- `.bedrock/node-config/genesis.json` → `/genesis.json` (read-only)

Future versions may support:

- Custom rippled.cfg
- Persistent ledger data
- Log directories

---

## Future Enhancements

Planned features for the local node:

- [ ] Custom port mapping
- [ ] Persistent ledger data (survive restarts)
- [ ] Snapshot/restore functionality
- [ ] Network forking (clone mainnet/testnet state)
- [ ] Multiple nodes (simulate network)
- [ ] Log streaming (`bedrock node logs`)
- [ ] Resource limits (CPU, memory)
- [ ] Health checks and auto-restart

---

## See Also

- [Building Contracts](./building-contracts.md)
- [Deploying and Calling Contracts](./deployment-and-calling.md)
- [ABI Generation](./abi-generation.md)
- [Wallet Management](./wallet.md)
