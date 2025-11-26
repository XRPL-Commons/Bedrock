# ABI Generation

**Bedrock** generates Application Binary Interface (ABI) definitions from annotated Rust smart contract source code.

## Overview

Smart contracts need ABIs to describe their functions, parameters, and types. Bedrock extracts this information from JSDoc-style comments in your Rust source code, similar to how Ethereum's Solidity compiler generates ABIs.

Bedrock's ABI generation provides:

- Automatic ABI extraction from source annotations as part of the `deploy` command.
- Type validation against XRPL type system
- JSON output
- Clear error messages for invalid annotations

## Why ABIs Matter

WASM binaries only contain function **names**, not parameter types or descriptions. To deploy and interact with contracts, you need:

- Function names
- Parameter names and types
- Parameter order
- Return types

Without an ABI, developers would need to manually maintain this information separately from the code (error-prone and not scalable).

## Automatic Generation

Bedrock automatically generates the ABI file (`abi.json`) during the `bedrock deploy` process. You no longer need to run a separate command.

---

## Annotation Syntax

Bedrock uses JSDoc-style annotations to define ABIs directly in Rust source code.

### Basic Structure

```rust
/// @xrpl-function function_name
/// @param param_name TYPE - Description
/// @return TYPE - Description
fn function_name(param_name: RustType) -> RustReturnType {
    // Implementation
}
```

### `@xrpl-function`

Marks a function for ABI inclusion.

**Syntax:**
```rust
/// @xrpl-function function_name
```

**Example:**
```rust
/// @xrpl-function register
fn register(name: Blob, resolver: AccountId, duration: u64) -> i32 {
    // ...
}
```

**Rules:**

- Must immediately precede the function declaration
- Function name must match the actual Rust function name
- One `@xrpl-function` per function

---

### `@param`

Defines a function parameter.

**Syntax:**
```rust
/// @param name TYPE - Description
```

**Example:**
```rust
/// @xrpl-function transfer
/// @param name VL - Domain name to transfer
/// @param new_owner ACCOUNT - New owner address
fn transfer(name: Blob, new_owner: AccountId) -> i32 {
    // ...
}
```

**Rules:**

- Must follow `@xrpl-function`
- Order matters - must match function signature order
- Parameter name should match Rust parameter name
- TYPE must be a valid XRPL type (see Type System below)
- Description is optional but recommended

---

### `@return`

Defines the function's return type.

**Syntax:**
```rust
/// @return TYPE - Description
```

**Example:**
```rust
/// @xrpl-function resolve
/// @param name VL - Domain name
/// @return ACCOUNT - Resolver address
fn resolve(name: Blob) -> AccountId {
    // ...
}
```

**Rules:**

- Optional (omit for void/status-only functions)
- Must follow `@xrpl-function` and any `@param` annotations
- TYPE must be a valid XRPL type
- Description is optional

---

### `@flag`

Sets the parameter flag for subsequent parameters.

**Syntax:**
```rust
/// @flag 0
```

**Example:**
```rust
/// @xrpl-function register
/// @flag 0
/// @param name VL - Domain name (required)
/// @param resolver ACCOUNT - Resolver (required)
/// @flag 1
/// @param duration UINT64 - Duration (optional)
fn register(name: Blob, resolver: AccountId, duration: u64) -> i32 {
    // ...
}
```

**Flag Values:**

- `0` - Required parameter (default)
- `1` - Optional parameter
- Other values reserved for future use

**Rules:**

- Applies to all subsequent `@param` annotations until next `@flag`
- Default is `0` if not specified

---

## Type System

Bedrock validates all types against the XRPL smart contract type system.

### Primitive Integer Types

| Type       | Rust Type    | Range                        | Common Uses                  |
| ---------- | ------------ | ---------------------------- | ---------------------------- |
| `UINT8`    | `u8`         | 0 to 255                     | Small counters, flags        |
| `UINT16`   | `u16`        | 0 to 65,535                  | Transaction types            |
| `UINT32`   | `u32`        | 0 to 4,294,967,295           | Sequence numbers, timestamps |
| `UINT64`   | `u64`        | 0 to 2^64-1                  | Token amounts, XRP drops     |
| `UINT128`  | `u128`       | 0 to 2^128-1                 | Large amounts                |
| `UINT160`  | `[u8; 20]`   | 20 bytes                     | Account IDs (internal)       |
| `UINT192`  | `[u8; 24]`   | 24 bytes                     | Extended identifiers         |
| `UINT256`  | `[u8; 32]`   | 32 bytes                     | Hashes, NFT IDs              |

### Specialized Types

| Type       | Rust Type                | Description                                    |
| ---------- | ------------------------ | ---------------------------------------------- |
| `VL`       | `Vec<u8>` or `&[u8]`     | Variable-length binary data                    |
| `ACCOUNT`  | `AccountID`              | XRPL account identifier (20 bytes)             |
| `AMOUNT`   | `Amount`                 | XRP, IOU, or MPT amount                        |
| `ISSUE`    | `Issue`                  | Currency and issuer pair                       |
| `CURRENCY` | `Currency`               | Currency code (3-letter or 160-bit hex)        |
| `NUMBER`   | `f64`                    | Floating-point number                          |

### Type Validation

Bedrock will **reject** invalid types with helpful error messages during the `deploy` process.

---

## Complete Example

### DNS Contract

**contract/src/lib.rs:**

```rust
use xrpl_wasm_std::core::types::account_id::AccountID;
use xrpl_wasm_std::core::types::blob::Blob;

/// @xrpl-function register
/// @param name VL - Domain name to register
/// @param resolver ACCOUNT - Resolver account address
/// @param duration UINT64 - Registration duration in seconds
fn register(name: Blob, resolver: AccountID, duration: u64) -> i32 {
    // Implementation
    0
}

/// @xrpl-function transfer
/// @param name VL - Domain name to transfer
/// @param new_owner ACCOUNT - New owner address
fn transfer(name: Blob, new_owner: AccountID) -> i32 {
    // Implementation
    0
}

/// @xrpl-function resolve
/// @param name VL - Domain name to resolve
/// @return ACCOUNT - Resolver address
fn resolve(name: Blob) -> AccountID {
    // Implementation
    AccountID([0u8; 20])
}

/// @xrpl-function renew
/// @param name VL - Domain name to renew
/// @flag 0
/// @param duration UINT64 - Extension duration (required)
/// @flag 1
/// @param auto_renew UINT8 - Enable auto-renewal (optional)
fn renew(name: Blob, duration: u64, auto_renew: u8) -> i32 {
    // Implementation
    0
}
```

### Generated ABI

**abi.json:**

```json
{
  "contract_name": "dns-contract",
  "functions": [
    {
      "name": "register",
      "parameters": [
        {
          "name": "name",
          "type": "VL",
          "flag": 0,
          "description": "Domain name to register"
        },
        {
          "name": "resolver",
          "type": "ACCOUNT",
          "flag": 0,
          "description": "Resolver account address"
        },
        {
          "name": "duration",
          "type": "UINT64",
          "flag": 0,
          "description": "Registration duration in seconds"
        }
      ]
    },
    {
      "name": "transfer",
      "parameters": [
        {
          "name": "name",
          "type": "VL",
          "flag": 0,
          "description": "Domain name to transfer"
        },
        {
          "name": "new_owner",
          "type": "ACCOUNT",
          "flag": 0,
          "description": "New owner address"
        }
      ]
    },
    {
      "name": "resolve",
      "parameters": [
        {
          "name": "name",
          "type": "VL",
          "flag": 0,
          "description": "Domain name to resolve"
        }
      ],
      "returns": {
        "type": "ACCOUNT",
        "description": "Resolver address"
      }
    },
    {
      "name": "renew",
      "parameters": [
        {
          "name": "name",
          "type": "VL",
          "flag": 0,
          "description": "Domain name to renew"
        },
        {
          "name": "duration",
          "type": "UINT64",
          "flag": 0,
          "description": "Extension duration (required)"
        },
        {
          "name": "auto_renew",
          "type": "UINT8",
          "flag": 1,
          "description": "Enable auto-renewal (optional)"
        }
      ]
    }
  ]
}
```

---

## Common Workflows

### Initial Development

```bash
# Write contract with annotations
vim contract/src/lib.rs

# Deploy the contract, which will also generate the ABI
bedrock deploy
```

### Iterative Development

```bash
# Modify function signatures or add new functions
vim contract/src/lib.rs

# Build and redeploy
bedrock build --debug
bedrock deploy --network local
```

---

## Troubleshooting

### "invalid type 'X' for parameter 'Y'"

**Problem:** Using a type that's not in the XRPL type system.

**Solution:**

Check the Type System section above. Common mistakes:
- `STRING` → Use `VL` (variable length)
- `ADDRESS` → Use `ACCOUNT`
- `INT` → Use `UINT32`, `UINT64`, etc.

### "function name mismatch"

**Problem:** `@xrpl-function` name doesn't match Rust function name.

```rust
/// @xrpl-function registr  // Typo!
fn register(...) { }
```

**Solution:** Ensure names match exactly:

```rust
/// @xrpl-function register
fn register(...) { }
```

### "no functions found in ABI"

**Problem:** No `@xrpl-function` annotations found.

**Solution:**

1. Ensure you're using `///` (doc comments), not `//`
2. Check that annotations are directly above the function
3. Verify you're in a bedrock project directory

### "duplicate parameter name 'X'"

**Problem:** Two `@param` annotations with the same name in one function.

**Solution:** Each parameter must have a unique name:

```rust
/// @xrpl-function swap
/// @param token_in ACCOUNT
/// @param token_out ACCOUNT  // OK - different name even if same type
```

---

## Advanced Usage

### Multiple Source Files

Bedrock automatically scans all `.rs` files in `contract/src/`:

```
contract/src/
├── lib.rs        // Main contract
├── registry.rs   // Registry functions
└── utils.rs      // Utility functions (no @xrpl-function annotations)
```

Only functions with `@xrpl-function` will be included in the ABI.

### CI/CD Integration

```bash
# In your CI pipeline
bedrock build --release
bedrock deploy

# Verify ABI exists and is valid
test -f abi.json || exit 1
```

### Version Control

**Recommended:** Commit `abi.json` to version control.

**Why?**
- Reviewers can see ABI changes in PRs
- Deployment scripts can rely on committed ABI
- Provides audit trail of interface changes

---

## Comparison with Ethereum

| Feature                  | Ethereum (Solidity)    | XRPL (Bedrock)         |
| ------------------------ | ---------------------- | ----------------------------- |
| **ABI Source**           | Function signatures    | JSDoc annotations             |
| **Generation**           | Automatic (compiler)   | Automatic (`deploy`)  |
| **Type System**          | Solidity types         | XRPL STypes                   |
| **Format**               | JSON                   | JSON            |
| **Parameter Metadata**   | Limited                | Descriptions, flags           |
| **Validation**           | Compile-time           | Deploy-time                 |

**Key Difference:**

Solidity's ABI is auto-generated because the language has rich type information. Rust → WASM loses this, so Bedrock uses annotations as the source of truth.

---

## Design Decisions

### Why Comments Instead of Rust Attributes?

**Considered:**
```rust
#[xrpl_function(name = "register")]
fn register(
    #[xrpl_param(type = "VL")] name: Blob
) { }
```

**Why we chose comments:**

1. **Simpler MVP**: No proc macro crate needed
2. **Language-agnostic parser**: Go can easily parse comments
3. **Familiar**: Developers know JSDoc/Javadoc
4. **No runtime cost**: Comments don't affect WASM output
5. **Future-compatible**: Can migrate to attributes later

### Why Not Parse WASM Binary?

WASM export section only contains function **names**, not types:

```wasm
(export "register" (func $register))  // No type info!
```

To get types, you'd need:
- DWARF debug info (huge, not in release builds)
- Custom sections (not standardized)
- Name section (names only, not types)

**Annotations are the pragmatic solution.**

---

## Future Enhancements

Planned features for ABI generation:

- [ ] Support for struct/enum types in parameters
- [ ] Rust attribute alternative (`#[xrpl_function]`)
- [ ] ABI versioning and compatibility checks
- [ ] TypeScript type generation
- [ ] ABI diff tool (compare versions)
- [ ] Integration with contract testing frameworks
- [ ] Validation against deployed contract

---

## See Also

- [Building Contracts](./building-contracts.md)
- [Deploying and Calling Contracts](./deployment-and-calling.md)
- [Local Node Management](./local-node.md)
- [Wallet Management](./wallet.md)
