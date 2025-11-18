package abi

// ValidXRPLTypes contains all valid XRPL smart contract types
var ValidXRPLTypes = map[string]TypeInfo{
	// Primitive Integer Types
	"UINT8":   {Name: "UINT8", RustType: "u8", Description: "8-bit unsigned integer (0 to 255)"},
	"UINT16":  {Name: "UINT16", RustType: "u16", Description: "16-bit unsigned integer (0 to 65,535)"},
	"UINT32":  {Name: "UINT32", RustType: "u32", Description: "32-bit unsigned integer (0 to 4,294,967,295)"},
	"UINT64":  {Name: "UINT64", RustType: "u64", Description: "64-bit unsigned integer"},
	"UINT128": {Name: "UINT128", RustType: "u128", Description: "128-bit unsigned integer"},
	"UINT160": {Name: "UINT160", RustType: "[u8; 20]", Description: "160-bit unsigned integer (20 bytes)"},
	"UINT192": {Name: "UINT192", RustType: "[u8; 24]", Description: "192-bit unsigned integer (24 bytes)"},
	"UINT256": {Name: "UINT256", RustType: "[u8; 32]", Description: "256-bit unsigned integer (32 bytes)"},

	// Specialized Types
	"VL":       {Name: "VL", RustType: "Vec<u8> or &[u8]", Description: "Variable-length binary data"},
	"ACCOUNT":  {Name: "ACCOUNT", RustType: "AccountID", Description: "XRPL account identifier (20 bytes)"},
	"AMOUNT":   {Name: "AMOUNT", RustType: "Amount", Description: "XRP, IOU, or MPT amount"},
	"ISSUE":    {Name: "ISSUE", RustType: "Issue", Description: "Currency and issuer pair"},
	"CURRENCY": {Name: "CURRENCY", RustType: "Currency", Description: "Currency code (3-letter or 160-bit hex)"},
	"NUMBER":   {Name: "NUMBER", RustType: "f64", Description: "Floating-point number"},
}

// TypeInfo contains information about an XRPL type
type TypeInfo struct {
	Name        string
	RustType    string
	Description string
}

// ABI represents a contract's Application Binary Interface
type ABI struct {
	ContractName string     `json:"contract_name"`
	Functions    []Function `json:"functions"`
}

// Function represents a contract function definition
type Function struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Returns    *ReturnType `json:"returns,omitempty"`
}

// Parameter represents a function parameter
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Flag        int    `json:"flag"`
	Description string `json:"description,omitempty"`
}

// ReturnType represents a function's return type
type ReturnType struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// IsValidType checks if a type name is valid
func IsValidType(typeName string) bool {
	_, ok := ValidXRPLTypes[typeName]
	return ok
}

// GetTypeInfo returns information about a type
func GetTypeInfo(typeName string) (TypeInfo, bool) {
	info, ok := ValidXRPLTypes[typeName]
	return info, ok
}
