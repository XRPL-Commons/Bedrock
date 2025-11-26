# Building Contracts

**Bedrock** compiles Rust smart contracts to WebAssembly (WASM) for deployment on XRPL. It wraps `cargo` with sensible defaults and a great developer experience.

## Overview

Bedrock's build process provides:

- Smart defaults for WASM compilation
- Automatic toolchain validation
- Progress indicators and formatted output
- Build optimization options
- Artifact management

## Commands

### `bedrock build`

Compiles your Rust smart contract to WASM.

```bash
bedrock build
```

**Flags:**

- `--debug` - Build in debug mode (faster, larger WASM)

**Examples:**

```bash
# Production build (optimized, smaller WASM)
bedrock build

# Development build (faster, larger WASM)
bedrock build --debug
```

**What it does:**

1. Validates Rust toolchain (cargo, rustc)
2. Ensures `wasm32-unknown-unknown` target is installed
3. Reads build configuration from `bedrock.toml`
4. Compiles Rust to WASM using cargo
5. Reports build results (path, size, duration)

**Output (Release build):**

```
Building smart contract...
   Mode: Release (optimized)
   Source: contract/src/lib.rs

 Compiling Rust → WASM...
   Compiling xrpl-contract v0.1.0 (/path/to/contract)
    Finished release [optimized] target(s) in 5.12s

✓ Build completed successfully!

Output:   contract/target/wasm32-unknown-unknown/release/xrpl_contract.wasm
Size:     156.4 KB
Duration: 5.1s

Built with release optimizations (smaller size, slower build)
```

**Output (Debug build):**

```
Building smart contract...
   Mode: Debug
   Source: contract/src/lib.rs

 Compiling Rust → WASM...
   Compiling xrpl-contract v0.1.0 (/path/to/contract)
    Finished dev [unoptimized + debuginfo] target(s) in 2.34s

✓ Build completed successfully!

Output:   contract/target/wasm32-unknown-unknown/debug/xrpl_contract.wasm
Size:     1.2 MB
Duration: 2.3s

Tip: Use default build for optimized builds
```

**Requirements:**

- Rust toolchain installed (`cargo`, `rustc`)
- Must be in a bedrock project directory (with `bedrock.toml`)
- Valid `Cargo.toml` in `contract/` directory

---

## Configuration

Bedrock reads build settings from `bedrock.toml`:

```toml
[build]
source = "contract/src/lib.rs"
target = "wasm32-unknown-unknown"
```

### Configuration Options

| Option   | Description                  | Default                                          |
| -------- | ---------------------------- | ------------------------------------------------ |
| `source` | Path to contract source file | `contract/src/lib.rs`                            |
| `target` | Rust compilation target      | `wasm32-unknown-unknown`                         |

---

## Contract Structure

Your contract should be structured as:

```
my-contract/
├── bedrock.toml           # Project config
├── contract/
│   ├── Cargo.toml         # Rust package manifest
│   ├── src/
│   │   └── lib.rs         # Contract code
│   └── target/            # Build output (gitignored)
│       └── wasm32-unknown-unknown/
│           ├── debug/
│           │   └── *.wasm
│           └── release/
│               └── *.wasm
```

### Cargo.toml

Your `contract/Cargo.toml` should specify:

```toml
[package]
name = "my-contract"
version = "0.1.0"
edition = "2021"

[lib]
crate-type = ["cdylib"]    # Required for WASM

[dependencies]
# Your XRPL dependencies

[profile.release]
opt-level = "z"            # Optimize for size
lto = true                 # Link-time optimization
strip = true               # Strip debug symbols
panic = "abort"            # Smaller panic handler
```

**Key settings for WASM:**

- `crate-type = ["cdylib"]` - Creates dynamic library for WASM
- `opt-level = "z"` - Maximize size optimization
- `lto = true` - Enable link-time optimization
- `strip = true` - Remove debug symbols
- `panic = "abort"` - Use simpler panic handler

---

## Build Modes

### Release Build (Default)

```bash
bedrock build
```

**Characteristics:**

- ❌ Slower compilation (~5-15 seconds)
- ✅ Heavily optimized
- ✅ Small WASM file (50-200 KB)
- ✅ Production-ready
- ❌ Harder to debug

**Use for:**

- Deployment to testnet/mainnet
- Performance testing
- Gas optimization
- Final builds

### Debug Build

```bash
bedrock build --debug
```

**Characteristics:**

- ✅ Fast compilation (~2-5 seconds)
- ✅ Includes debug symbols
- ✅ Better error messages
- ❌ Large WASM file (1-2 MB)
- ❌ Not optimized

**Use for:**

- Rapid development iteration
- Debugging contract logic
- Testing locally

**Size comparison:**

```
Debug:   1.2 MB  (unoptimized)
Release: 156 KB  (~87% smaller)
```

---

## Toolchain Requirements

### Installing Rust

If you don't have Rust installed:

```bash
# Install rustup (Rust installer)
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Restart your shell, then verify
rustc --version
cargo --version
```

### WASM Target

Bedrock automatically installs the WASM target if missing:

```bash
# This happens automatically, but you can also run manually:
rustup target add wasm32-unknown-unknown
```

### Verifying Toolchain

```bash
# Check Rust version
rustc --version
# Should show: rustc 1.70.0 or later

# Check cargo
cargo --version
# Should show: cargo 1.70.0 or later

# Check WASM target
rustup target list | grep wasm32-unknown-unknown
# Should show: wasm32-unknown-unknown (installed)
```

---

## Common Workflows

### Initial Development

```bash
# Initialize project
bedrock init my-contract
cd my-contract

# Edit contract
vim contract/src/lib.rs

# Build and test locally
bedrock build --debug

# Verify WASM was created
ls -lh contract/target/wasm32-unknown-unknown/debug/*.wasm
```

### Iterative Development

```bash
# Make changes
vim contract/src/lib.rs

# Quick build
bedrock build --debug

# Test deployment (when ready)
bedrock deploy --network local
```

### Preparing for Production

```bash
# Build optimized version
bedrock build

# Verify size
ls -lh contract/target/wasm32-unknown-unknown/release/*.wasm

# Deploy to testnet
bedrock deploy --network alphanet
```

---

## Advanced Usage

### Custom Cargo Features

If your `Cargo.toml` defines features:

```toml
[features]
default = ["std"]
std = []
debug-mode = []
```

Currently, bedrock build uses cargo defaults. To use custom features, you can run cargo directly:

```bash
cd contract
cargo build --target wasm32-unknown-unknown --release --no-default-features
```

### Build Scripts

For complex builds, you can create a build script:

```bash
# build.sh
#!/bin/bash
bedrock build
# Post-process WASM
wasm-opt contract/target/wasm32-unknown-unknown/release/*.wasm -Oz -o optimized.wasm
```

### Cargo Configuration

You can customize cargo behavior with `.cargo/config.toml`:

```toml
# contract/.cargo/config.toml
[build]
target = "wasm32-unknown-unknown"

[target.wasm32-unknown-unknown]
rustflags = ["-C", "link-arg=-s"]
```

---

## Optimization Tips

### 1. Minimize Dependencies

```toml
# Bad: Large dependency tree
[dependencies]
serde = { version = "1.0", features = ["derive"] }
regex = "1.0"

# Good: Only what you need
[dependencies]
# Minimal XRPL-specific deps only
```

### 2. Use `opt-level = "z"`

```toml
[profile.release]
opt-level = "z"  # Smaller than "s" or "3"
```

### 3. Enable LTO

```toml
[profile.release]
lto = true  # Link-time optimization
```

### 4. Strip Symbols

```toml
[profile.release]
strip = true  # Remove debug info
```

### 5. Avoid Panics

```rust
// Bad: Uses panic machinery
assert!(condition);

// Good: Handle errors explicitly
if !condition {
    return Err(Error::InvalidInput);
}
```

### Size Comparison

| Configuration                   | WASM Size | Build Time |
| ------------------------------- | --------- | ---------- |
| Debug default                   | 1.2 MB    | 2s         |
| Release `opt-level = "3"`       | 450 KB    | 5s         |
| Release `opt-level = "z"`       | 180 KB    | 6s         |
| Release `opt-level = "z"` + LTO | 156 KB    | 8s         |
| + wasm-opt `-Oz`                | 95 KB     | +2s        |

---

## Troubleshooting

### "cargo not found"

**Problem:** Rust toolchain not installed.

**Solution:**

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source $HOME/.cargo/env
```

### "error: couldn't find library for wasm32-unknown-unknown"

**Problem:** WASM target not installed.

**Solution:**

```bash
rustup target add wasm32-unknown-unknown
```

### "Cargo.toml not found"

**Problem:** Not in a bedrock project or contract directory missing.

**Solution:**

1. Run `bedrock init` to create a project
2. Ensure you're in the project root directory
3. Check that `contract/Cargo.toml` exists

### Build Fails with Dependency Errors

**Problem:** Missing or incompatible dependencies.

**Solution:**

```bash
cd contract
cargo update  # Update dependencies
cargo build --target wasm32-unknown-unknown  # Test build
```

### "no .wasm file found"

**Problem:** Build succeeded but WASM not generated.

**Solution:**

1. Check `Cargo.toml` has `crate-type = ["cdylib"]`
2. Ensure you have a `lib.rs` (not just `main.rs`)
3. Run `bedrock build --debug` to see cargo output

### Build is Very Slow

**Problem:** Release builds are inherently slower.

**Solutions:**

- Use debug builds during development
- Use `cargo build --release` separately if you need incremental builds
- Consider using `sccache` for caching
- Reduce dependencies

---

## Comparison with Direct Cargo

### Using `bedrock build`

```bash
bedrock build
```

✅ Validates toolchain
✅ Pretty output
✅ Reports size/duration
✅ Consistent across team
✅ Integrates with bedrock

### Using Cargo Directly

```bash
cd contract
cargo build --target wasm32-unknown-unknown --release
```

✅ More control
✅ Cargo's native features
❌ More verbose
❌ Manual target specification
❌ No validation

**Both are valid!** Use `bedrock build` for convenience, cargo for advanced needs.

---

## Under the Hood

`bedrock build` essentially runs:

```bash
# For release (default)
cargo build --target wasm32-unknown-unknown --release

# For debug
cargo build --target wasm32-unknown-unknown
```

**Additional steps:**

1. Verifies `cargo` and `rustc` exist
2. Runs `rustup target add wasm32-unknown-unknown`
3. Changes to `contract/` directory
4. Executes cargo command
5. Finds the output `.wasm` file
6. Calculates size and duration
7. Formats output
