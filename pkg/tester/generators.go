package tester

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

// ValueGenerator generates random values matching XRPL ABI types
type ValueGenerator struct {
	rng *rand.Rand
}

// NewValueGenerator creates a new generator with the given seed
func NewValueGenerator(seed int64) *ValueGenerator {
	return &ValueGenerator{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Generate produces a random value for the given XRPL type
func (g *ValueGenerator) Generate(typeName string) interface{} {
	switch typeName {
	case "UINT8":
		return g.rng.Intn(256)
	case "UINT16":
		return g.rng.Intn(65536)
	case "UINT32":
		return g.rng.Uint32()
	case "UINT64":
		return g.rng.Uint64()
	case "UINT128":
		return fmt.Sprintf("%032X", g.randomBytes(16))
	case "UINT160":
		return fmt.Sprintf("%040X", g.randomBytes(20))
	case "UINT192":
		return fmt.Sprintf("%048X", g.randomBytes(24))
	case "UINT256":
		return fmt.Sprintf("%064X", g.randomBytes(32))
	case "VL":
		length := g.rng.Intn(64) + 1
		return hex.EncodeToString(g.randomBytes(length))
	case "ACCOUNT":
		return g.generateAccount()
	case "AMOUNT":
		return g.generateAmount()
	case "ISSUE":
		return g.generateIssue()
	case "CURRENCY":
		return g.generateCurrency()
	case "NUMBER":
		return g.rng.Float64() * 1000000
	default:
		return g.rng.Uint64()
	}
}

// Shrink attempts to find a minimal value that still triggers the failure
func (g *ValueGenerator) Shrink(typeName string, value interface{}) interface{} {
	switch typeName {
	case "UINT8", "UINT16", "UINT32", "UINT64":
		return shrinkUint(value)
	case "VL":
		if s, ok := value.(string); ok && len(s) > 2 {
			return s[:len(s)/2]
		}
		return value
	default:
		return value
	}
}

func shrinkUint(v interface{}) interface{} {
	switch val := v.(type) {
	case int:
		if val > 0 {
			return val / 2
		}
		return 0
	case uint32:
		if val > 0 {
			return val / 2
		}
		return uint32(0)
	case uint64:
		if val > 0 {
			return val / 2
		}
		return uint64(0)
	default:
		return v
	}
}

func (g *ValueGenerator) randomBytes(n int) []byte {
	b := make([]byte, n)
	g.rng.Read(b)
	return b
}

func (g *ValueGenerator) generateAccount() string {
	// Generate a plausible XRPL address (starts with 'r')
	// This is not a real valid address but follows the format
	const chars = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"
	b := make([]byte, 25)
	b[0] = 'r'
	for i := 1; i < 25; i++ {
		b[i] = chars[g.rng.Intn(len(chars))]
	}
	return string(b)
}

func (g *ValueGenerator) generateAmount() string {
	// Generate a random XRP amount in drops
	drops := g.rng.Int63n(100_000_000_000) // Up to 100,000 XRP
	return fmt.Sprintf("%d", drops)
}

func (g *ValueGenerator) generateIssue() map[string]string {
	return map[string]string{
		"currency": g.generateCurrency(),
		"issuer":   g.generateAccount(),
	}
}

func (g *ValueGenerator) generateCurrency() string {
	// Generate a random 3-letter currency code
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 3)
	for i := range b {
		b[i] = letters[g.rng.Intn(len(letters))]
	}
	return string(b)
}
