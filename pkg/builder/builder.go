package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Builder handles contract compilation
type Builder struct {
	projectRoot string
}

// New creates a new Builder
func New(projectRoot string) *Builder {
	return &Builder{projectRoot: projectRoot}
}

// Build compiles the Rust contract to WASM
func (b *Builder) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	// Verify cargo is installed
	if err := b.verifyToolchain(); err != nil {
		return nil, err
	}

	// Ensure wasm32 target is installed
	if err := b.ensureWasmTarget(ctx); err != nil {
		return nil, fmt.Errorf("failed to add wasm32 target: %w", err)
	}

	contractDir := filepath.Join(b.projectRoot, "contract")

	// Check if Cargo.toml exists
	if _, err := os.Stat(filepath.Join(contractDir, "Cargo.toml")); os.IsNotExist(err) {
		return nil, fmt.Errorf("Cargo.toml not found in %s", contractDir)
	}

	// Build command
	args := []string{"build", "--target", "wasm32-unknown-unknown"}
	if opts.Release {
		args = append(args, "--release")
	}
	if opts.Verbose {
		args = append(args, "--verbose")
	}

	cmd := exec.CommandContext(ctx, "cargo", args...)
	cmd.Dir = contractDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	startTime := time.Now()

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cargo build failed: %w", err)
	}

	duration := time.Since(startTime)

	// Determine output path
	buildType := "debug"
	if opts.Release {
		buildType = "release"
	}

	wasmPath := filepath.Join(contractDir, "target", "wasm32-unknown-unknown", buildType)

	// Find the .wasm file
	wasmFile, size, err := b.findWasmFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find WASM output: %w", err)
	}

	return &BuildResult{
		WasmPath:  filepath.Join(wasmPath, wasmFile),
		Size:      size,
		Duration:  duration,
		Optimized: opts.Release,
	}, nil
}

// Clean removes build artifacts
func (b *Builder) Clean(ctx context.Context) error {
	contractDir := filepath.Join(b.projectRoot, "contract")

	cmd := exec.CommandContext(ctx, "cargo", "clean")
	cmd.Dir = contractDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cargo clean failed: %w", err)
	}

	return nil
}

// verifyToolchain checks if cargo and rustc are installed
func (b *Builder) verifyToolchain() error {
	if _, err := exec.LookPath("cargo"); err != nil {
		return fmt.Errorf("cargo not found: please install Rust from https://rustup.rs")
	}
	if _, err := exec.LookPath("rustc"); err != nil {
		return fmt.Errorf("rustc not found: please install Rust from https://rustup.rs")
	}
	return nil
}

// ensureWasmTarget adds wasm32-unknown-unknown target if not present
func (b *Builder) ensureWasmTarget(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "rustup", "target", "add", "wasm32-unknown-unknown")
	if err := cmd.Run(); err != nil {
		// Not fatal if rustup fails (target might already be installed)
		return nil
	}
	return nil
}

// findWasmFile finds the first .wasm file in the directory
func (b *Builder) findWasmFile(dir string) (string, int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".wasm" {
			info, err := entry.Info()
			if err != nil {
				return "", 0, err
			}
			return entry.Name(), info.Size(), nil
		}
	}

	return "", 0, fmt.Errorf("no .wasm file found in %s", dir)
}
