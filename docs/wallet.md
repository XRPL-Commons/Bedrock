# Wallet Management (Jade)

**Jade** is Bedrock's built-in wallet management tool. It allows you to create, import, and manage XRPL wallets securely.

## Overview

Jade provides a secure way to manage your XRPL wallets by encrypting them and storing them on disk. This means you don't have to expose your wallet seeds in your command-line history or scripts.

## Commands

### `bedrock wallet new <name>`

Creates a new XRPL wallet.

**Usage:**

```bash
bedrock wallet new <name> [flags]
```

**Arguments:**

-   `<name>`: A unique name for your wallet.

**Flags:**

-   `--algorithm`, `-a`: The cryptographic algorithm to use. Can be `secp256k1` or `ed25519`. Defaults to `secp256k1`.

**What it does:**

1.  Generates a new, random XRPL wallet.
2.  Prompts you to enter a password to encrypt the wallet.
3.  Saves the encrypted wallet to `~/.config/bedrock/wallets/<name>.json`.

**Example:**

```bash
bedrock wallet new my-test-wallet
```

---

### `bedrock wallet import <name>`

Imports an existing XRPL wallet from a seed.

**Usage:**

```bash
bedrock wallet import <name> [flags]
```

**Arguments:**

-   `<name>`: A unique name for your wallet.

**Flags:**

-   `--algorithm`, `-a`: The cryptographic algorithm of the wallet you're importing. Can be `secp256k1` or `ed25519`. Defaults to `secp256k1`.

**What it does:**

1.  Prompts you to enter the seed for the wallet you want to import. The input will be hidden for security.
2.  Prompts you to enter a password to encrypt the wallet.
3.  Saves the encrypted wallet to `~/.config/bedrock/wallets/<name>.json`.

**Example:**

```bash
bedrock wallet import my-existing-wallet --algorithm ed25519
```

---

### `bedrock wallet list`

Lists all the wallets you have stored.

**Usage:**

```bash
bedrock wallet list
```

**Example:**

```bash
bedrock wallet list
```

---

### `bedrock wallet export <name>`

Exports the seed and address of a stored wallet.

**Usage:**

```bash
bedrock wallet export <name>
```

**Arguments:**

-   `<name>`: The name of the wallet to export.

**What it does:**

1.  Prompts you to enter the password for the wallet.
2.  Prints the wallet's name, address, and seed to the console.

**Example:**

```bash
bedrock wallet export my-test-wallet
```

---

### `bedrock wallet remove <name>`

Permanently removes a wallet from storage.

**Usage:**

```bash
bedrock wallet remove <name>
```

**Arguments:**

-   `<name>`: The name of the wallet to remove.

**What it does:**

1.  Asks for confirmation before deleting the wallet.
2.  Deletes the wallet's file from `~/.config/bedrock/wallets/`.

**This action cannot be undone. Make sure you have backed up your wallet's seed before removing it.**

**Example:**

```bash
bedrock wallet remove my-test-wallet
```

---

## See Also

- [Local Node Management](./local-node.md)
- [Building Contracts](./building-contracts.md)
- [Deploying and Calling Contracts](./deployment-and-calling.md)
- [ABI Generation](./abi-generation.md)
