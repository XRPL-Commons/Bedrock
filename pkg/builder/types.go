package builder

import "time"

// BuildOptions configures the build process
type BuildOptions struct {
	Release bool // Use --release flag
	Verbose bool // Show verbose output
}

// BuildResult contains information about the build
type BuildResult struct {
	WasmPath  string        // Path to the compiled WASM file
	Size      int64         // Size in bytes
	Duration  time.Duration // Build time
	Optimized bool          // Whether release optimizations were applied
}
