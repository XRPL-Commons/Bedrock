package embedded

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed modules/deploy/deploy.js modules/call/call.js modules/package.json
var ModulesFS embed.FS

var (
	cacheDir   string
	setupOnce  sync.Once
	setupError error
)

const (
	versionFile = ".bedrock-version"
)

// getCacheDir returns the XDG-compliant cache directory for bedrock
func getCacheDir() (string, error) {
	var baseDir string

	if runtime.GOOS == "windows" {
		// Windows: use LOCALAPPDATA
		baseDir = os.Getenv("LOCALAPPDATA")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
	} else {
		// Unix-like: use XDG_CACHE_HOME or default to ~/.cache
		baseDir = os.Getenv("XDG_CACHE_HOME")
		if baseDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			baseDir = filepath.Join(home, ".cache")
		}
	}

	cacheDir := filepath.Join(baseDir, "bedrock", "modules")
	return cacheDir, nil
}

// getModulesVersion returns a hash of the embedded modules to detect changes
func getModulesVersion() (string, error) {
	hasher := sha256.New()

	// Hash package.json
	packageJSON, err := ModulesFS.ReadFile("modules/package.json")
	if err != nil {
		return "", err
	}
	hasher.Write(packageJSON)

	// Hash deploy.js
	deployJS, err := ModulesFS.ReadFile("modules/deploy/deploy.js")
	if err != nil {
		return "", err
	}
	hasher.Write(deployJS)

	// Hash call.js
	callJS, err := ModulesFS.ReadFile("modules/call/call.js")
	if err != nil {
		return "", err
	}
	hasher.Write(callJS)

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// needsReinstall checks if modules need to be installed or reinstalled
func needsReinstall(cacheDir string) (bool, error) {
	// Check if cache directory exists
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		return true, nil
	}

	// Check if node_modules exists
	nodeModulesPath := filepath.Join(cacheDir, "node_modules")
	if _, err := os.Stat(nodeModulesPath); os.IsNotExist(err) {
		return true, nil
	}

	// Check version file
	versionPath := filepath.Join(cacheDir, versionFile)
	cachedVersion, err := os.ReadFile(versionPath)
	if err != nil {
		// No version file or can't read it - reinstall
		return true, nil
	}

	// Get current version
	currentVersion, err := getModulesVersion()
	if err != nil {
		return false, err
	}

	// Compare versions
	return string(cachedVersion) != currentVersion, nil
}

// SetupModules extracts embedded modules to cache and installs dependencies (lazy, once)
// This is called automatically on first use and cached for subsequent calls
func SetupModules() (string, error) {
	setupOnce.Do(func() {
		// Get cache directory
		cache, err := getCacheDir()
		if err != nil {
			setupError = fmt.Errorf("failed to get cache directory: %w", err)
			return
		}

		// Check if reinstall needed
		reinstall, err := needsReinstall(cache)
		if err != nil {
			setupError = fmt.Errorf("failed to check reinstall status: %w", err)
			return
		}

		if reinstall {
			// Clean up old cache if it exists
			if err := os.RemoveAll(cache); err != nil && !os.IsNotExist(err) {
				setupError = fmt.Errorf("failed to clean cache directory: %w", err)
				return
			}

			// Create cache directory
			if err := os.MkdirAll(cache, 0755); err != nil {
				setupError = fmt.Errorf("failed to create cache directory: %w", err)
				return
			}

			// Extract package.json
			packageJSON, err := ModulesFS.ReadFile("modules/package.json")
			if err != nil {
				setupError = fmt.Errorf("failed to read package.json: %w", err)
				return
			}

			packagePath := filepath.Join(cache, "package.json")
			if err := os.WriteFile(packagePath, packageJSON, 0644); err != nil {
				setupError = fmt.Errorf("failed to write package.json: %w", err)
				return
			}

			// Extract deploy.js
			deployJS, err := ModulesFS.ReadFile("modules/deploy/deploy.js")
			if err != nil {
				setupError = fmt.Errorf("failed to read deploy.js: %w", err)
				return
			}

			deployPath := filepath.Join(cache, "deploy.js")
			if err := os.WriteFile(deployPath, deployJS, 0755); err != nil {
				setupError = fmt.Errorf("failed to write deploy.js: %w", err)
				return
			}

			// Extract call.js
			callJS, err := ModulesFS.ReadFile("modules/call/call.js")
			if err != nil {
				setupError = fmt.Errorf("failed to read call.js: %w", err)
				return
			}

			callPath := filepath.Join(cache, "call.js")
			if err := os.WriteFile(callPath, callJS, 0755); err != nil {
				setupError = fmt.Errorf("failed to write call.js: %w", err)
				return
			}

			// Install npm dependencies
			fmt.Println("⚡ First run detected - installing JavaScript dependencies...")
			fmt.Printf("   Cache location: %s\n", cache)

			cmd := exec.Command("npm", "install", "--silent", "--no-progress")
			cmd.Dir = cache
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				setupError = fmt.Errorf("failed to install npm dependencies: %w", err)
				return
			}

			fmt.Println("✓ Dependencies installed successfully")

			// Write version file
			currentVersion, err := getModulesVersion()
			if err != nil {
				setupError = fmt.Errorf("failed to get modules version: %w", err)
				return
			}

			versionPath := filepath.Join(cache, versionFile)
			if err := os.WriteFile(versionPath, []byte(currentVersion), 0644); err != nil {
				setupError = fmt.Errorf("failed to write version file: %w", err)
				return
			}
		}

		cacheDir = cache
	})

	return cacheDir, setupError
}

// GetModulePath returns the path to an extracted module
func GetModulePath(moduleName string) (string, error) {
	dir, err := SetupModules()
	if err != nil {
		return "", err
	}

	modulePath := filepath.Join(dir, moduleName)
	if _, err := os.Stat(modulePath); err != nil {
		return "", fmt.Errorf("module %s not found: %w", moduleName, err)
	}

	return modulePath, nil
}
