package inspector

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xrpl-commons/bedrock/pkg/abi"
)

// WasmInfo contains information about a WASM module
type WasmInfo struct {
	Size      int64    // File size in bytes
	Functions []string // Exported function names
	Imports   []string // Imported function names
	MemPages  int      // Initial memory pages
}

// InspectWasm analyzes a WASM binary file
func InspectWasm(wasmPath string) (*WasmInfo, error) {
	data, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read WASM file: %w", err)
	}

	info := &WasmInfo{
		Size: int64(len(data)),
	}

	// Validate WASM magic number
	if len(data) < 8 || data[0] != 0x00 || data[1] != 0x61 || data[2] != 0x73 || data[3] != 0x6D {
		return nil, fmt.Errorf("invalid WASM file: bad magic number")
	}

	// Parse sections to extract exports and imports
	offset := 8 // Skip magic + version

	for offset < len(data) {
		if offset >= len(data) {
			break
		}

		sectionID := data[offset]
		offset++

		sectionSize, n := readLEB128(data[offset:])
		offset += n

		sectionEnd := offset + int(sectionSize)
		if sectionEnd > len(data) {
			break
		}

		switch sectionID {
		case 2: // Import section
			info.Imports = parseImports(data[offset:sectionEnd])
		case 5: // Memory section
			info.MemPages = parseMemory(data[offset:sectionEnd])
		case 7: // Export section
			info.Functions = parseExports(data[offset:sectionEnd])
		}

		offset = sectionEnd
	}

	return info, nil
}

// InspectABI loads and analyzes an ABI file
func InspectABI(abiPath string) (*abi.ABI, error) {
	data, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}

	var a abi.ABI
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &a, nil
}

func readLEB128(data []byte) (uint32, int) {
	var result uint32
	var shift uint
	var n int

	for n < len(data) {
		b := data[n]
		n++
		result |= uint32(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	return result, n
}

func parseExports(data []byte) []string {
	if len(data) == 0 {
		return nil
	}

	count, n := readLEB128(data)
	offset := n
	var exports []string

	for i := 0; i < int(count) && offset < len(data); i++ {
		nameLen, n := readLEB128(data[offset:])
		offset += n

		if offset+int(nameLen) > len(data) {
			break
		}

		name := string(data[offset : offset+int(nameLen)])
		offset += int(nameLen)

		if offset >= len(data) {
			break
		}

		kind := data[offset]
		offset++

		// Skip index
		_, n = readLEB128(data[offset:])
		offset += n

		if kind == 0 { // Function export
			exports = append(exports, name)
		}
	}

	return exports
}

func parseImports(data []byte) []string {
	if len(data) == 0 {
		return nil
	}

	count, n := readLEB128(data)
	offset := n
	var imports []string

	for i := 0; i < int(count) && offset < len(data); i++ {
		// Module name
		moduleLen, n := readLEB128(data[offset:])
		offset += n
		if offset+int(moduleLen) > len(data) {
			break
		}
		module := string(data[offset : offset+int(moduleLen)])
		offset += int(moduleLen)

		// Field name
		fieldLen, n := readLEB128(data[offset:])
		offset += n
		if offset+int(fieldLen) > len(data) {
			break
		}
		field := string(data[offset : offset+int(fieldLen)])
		offset += int(fieldLen)

		imports = append(imports, fmt.Sprintf("%s.%s", module, field))

		// Skip import descriptor
		if offset >= len(data) {
			break
		}
		kind := data[offset]
		offset++

		switch kind {
		case 0: // Function
			_, n = readLEB128(data[offset:])
			offset += n
		case 1: // Table
			offset += 3 // Simplified
		case 2: // Memory
			offset += 2 // Simplified
		case 3: // Global
			offset += 2 // Simplified
		}
	}

	return imports
}

func parseMemory(data []byte) int {
	if len(data) == 0 {
		return 0
	}

	// Skip flags
	offset := 1
	if offset >= len(data) {
		return 0
	}

	pages, _ := readLEB128(data[offset:])
	return int(pages)
}
