# Deploying and Calling Contracts

This document covers the `bedrock deploy` and `bedrock call` commands, which are used to deploy and interact with your smart contracts.

## `bedrock deploy`

The `deploy` command is a smart command that handles the entire deployment process.

### Smart Deployment

`bedrock deploy` automatically performs the following steps:

1.  **Builds the contract**: Compiles your Rust contract to WASM in release mode if it hasn't been built yet.
2.  **Generates the ABI**: Extracts the ABI from your Rust source code annotations if `abi.json` is missing.
3.  **Deploys to the network**: Deploys the compiled WASM to the specified network.

### Usage

```bash
bedrock deploy [flags]
```

### Options

| Flag          | Description                                    | Default    |
| ------------- | ---------------------------------------------- | ---------- |
| `--network`   | The network to deploy to (e.g., `local`, `alphanet`) | `alphanet` |
| `--skip-build`| Skip the build step                            | `false`    |
| `--skip-abi`  | Skip the ABI generation step                   | `false`    |
| `--wallet`    | The seed of the wallet to use for deployment   |            |

### Examples

**Deploy to a local node:**

```bash
bedrock deploy --network local
```

**Deploy to the alphanet (testnet):**

```bash
bedrock deploy --network alphanet
```

**Deploy using a specific wallet:**

```bash
bedrock deploy --wallet sXXX...
```

**Deploy without rebuilding the contract:**

```bash
bedrock deploy --skip-build
```

---

## `bedrock call`

The `call` command allows you to interact with the functions of a deployed smart contract.

### Usage

```bash
bedrock call <contract> <function> [flags]
```

### Arguments

-   `<contract>`: The address of the contract to call.
-   `<function>`: The name of the function to call.

### Options

| Flag            | Description                                  | Default    |
| --------------- | -------------------------------------------- | ---------- |
| `--network`     | The network the contract is deployed on        | `alphanet` |
| `--params`      | JSON string of parameters to pass to the function |            |
| `--params-file` | Path to a JSON file containing the parameters    |            |
| `--wallet`      | The seed of the wallet to use for the call (required) |            |
| `--gas`         | The amount of gas to provide for the transaction |            |
| `--fee`         | The transaction fee (in drops)                 |            |

### Examples

**Call a function without parameters:**

```bash
bedrock call rContract123... hello --wallet sXXX...
```

**Call a function with parameters:**

```bash
bedrock call rContract123... transfer \
  --wallet sXXX... \
  --params '{"to":"rRecipient...","amount":1000}'
```

**Call a function with parameters from a file:**

```bash
bedrock call rContract123... register \
  --wallet sXXX... \
  --params-file params.json
```
