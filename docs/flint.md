# Flint - Contract Building

**Flint** compiles Rust smart contracts to WebAssembly (WASM) for deployment on XRPL. It wraps `cargo` with sensible defaults and a great developer experience.

## Overview

Flint provides:

- Smart defaults for WASM compilation
- Automatic toolchain validation
- Progress indicators and formatted output
- Build optimization options
- Artifact management

## Commands

### `bedrock flint build`

Compiles your Rust smart contract to WASM.

```bash
bedrock flint build
```

**Flags:**

- `--release, -r` - Build with release optimizations (smaller WASM, slower build)
- `--verbose, -v` - Show detailed cargo output

**Examples:**

```bash
# Development build (faster, larger WASM)
bedrock flint build

# Production build (optimized, smaller WASM)
bedrock flint build --release

# Verbose output for debugging
bedrock flint build --verbose
```

**What it does:**

1. Validates Rust toolchain (cargo, rustc)
2. Ensures `wasm32-unknown-unknown` target is installed
3. Reads build configuration from `bedrock.toml`
4. Compiles Rust to WASM using cargo
5. Reports build results (path, size, duration)

**Output (Debug build):**

```
Flint - Building smart contract
   Mode: Debug
   Source: contract/src/lib.rs

 Compiling Rust → WASM...
   Compiling xrpl-contract v0.1.0 (/path/to/contract)
    Finished dev [unoptimized + debuginfo] target(s) in 2.34s

✓ Build completed successfully!

Output:   contract/target/wasm32-unknown-unknown/debug/xrpl_contract.wasm
Size:     1.2 MB
Duration: 2.3s

Tip: Use --release flag for optimized builds
```

**Output (Release build):**

```
Flint - Building smart contract
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

**Requirements:**

- Rust toolchain installed (`cargo`, `rustc`)
- Must be in a bedrock project directory (with `bedrock.toml`)
- Valid `Cargo.toml` in `contract/` directory

---

### `bedrock flint clean`

Removes all build artifacts and compiled files.

```bash
bedrock flint clean
```

**What it does:**

1. Runs `cargo clean` in the contract directory
2. Removes `contract/target/` directory and all contents

**Output:**

```
Cleaning build artifacts...

✓ Build artifacts cleaned!
```

**Use cases:**

- Free up disk space
- Force full rebuild
- Clean before distribution
- Troubleshoot build issues

---

## Configuration

Flint reads build settings from `bedrock.toml`:

```toml
[build]
source = "contract/src/lib.rs"
output = "contract/target/wasm32-unknown-unknown/release"
target = "wasm32-unknown-unknown"
```

### Configuration Options

| Option   | Description                  | Default                                          |
| -------- | ---------------------------- | ------------------------------------------------ |
| `source` | Path to contract source file | `contract/src/lib.rs`                            |
| `output` | Build output directory       | `contract/target/wasm32-unknown-unknown/release` |
| `target` | Rust compilation target      | `wasm32-unknown-unknown`                         |

**Note:** In the MVP, `source` and `target` are informational. Flint always builds from the `contract/` directory using cargo's default behavior.

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
xrpl-wasm-std = { path = "../xrpl-wasm-std" }
xrpl-wasm-macro = { path = "../xrpl-wasm-macro" }

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

### Debug Build (Default)

```bash
bedrock flint build
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

### Release Build

```bash
bedrock flint build --release
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

Flint automatically installs the WASM target if missing:

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
bedrock flint build

# Verify WASM was created
ls -lh contract/target/wasm32-unknown-unknown/debug/*.wasm
```

### Iterative Development

```bash
# Make changes
vim contract/src/lib.rs

# Quick build
bedrock flint build

# Test deployment (when ready)
bedrock slate deploy --network local
```

### Preparing for Production

```bash
# Clean previous builds
bedrock flint clean

# Build optimized version
bedrock flint build --release

# Verify size
ls -lh contract/target/wasm32-unknown-unknown/release/*.wasm

# Deploy to testnet
bedrock slate deploy --network alphanet --save
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

Currently, flint uses cargo defaults. To use custom features, you can run cargo directly:

```bash
cd contract
cargo build --target wasm32-unknown-unknown --release --no-default-features
```

**Future enhancement:** Flint will support feature flags:

```bash
bedrock flint build --release --features debug-mode --no-default-features
```

### Build Scripts

For complex builds, you can create a build script:

```bash
# build.sh
#!/bin/bash
bedrock flint clean
bedrock flint build --release
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
3. Run `bedrock flint build --verbose` to see cargo output

### Build is Very Slow

**Problem:** Release builds are inherently slower.

**Solutions:**

- Use debug builds during development
- Use `cargo build --release` separately if you need incremental builds
- Consider using `sccache` for caching
- Reduce dependencies

---

## Comparison with Direct Cargo

### Using Flint

```bash
bedrock flint build --release
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

**Both are valid!** Use flint for convenience, cargo for advanced needs.

---

## Under the Hood

Flint essentially runs:

```bash
# For debug
cargo build --target wasm32-unknown-unknown

# For release
cargo build --target wasm32-unknown-unknown --release
```

**Additional steps:**

1. Verifies `cargo` and `rustc` exist
2. Runs `rustup target add wasm32-unknown-unknown`
3. Changes to `contract/` directory
4. Executes cargo command
5. Finds the output `.wasm` file
6. Calculates size and duration
7. Formats output

**No magic**, just convenience!

---

## Future Enhancements

Planned features for flint:

- [ ] Watch mode (`--watch`) for auto-rebuild
- [ ] Custom feature flags (`--features`)
- [ ] WASM validation (check exports, imports)
- [ ] Automatic `wasm-opt` integration
- [ ] Size comparison vs. previous build
- [ ] Build caching
- [ ] Parallel builds for multiple contracts
- [ ] Integration with WASM analyzers
- [ ] Custom build hooks/scripts

---

## See Also

- [Basalt - Local Node Management](./basalt.md)
- [Quartz - ABI Generation](./quartz.md)
- [Slate - Contract Deployment](./slate.md)
- [Writing XRPL Smart Contracts](./writing-contracts.md)
- [bedrock.toml Configuration Reference](./config-reference.md)
