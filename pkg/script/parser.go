package script

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Script represents a multi-step deployment/interaction script
type Script struct {
	Name        string            `toml:"name" json:"name"`
	Description string            `toml:"description" json:"description"`
	Network     string            `toml:"network" json:"network"`
	Variables   map[string]string `toml:"variables" json:"variables"`
	Steps       []Step            `toml:"steps" json:"steps"`
}

// Step represents a single operation in a script
type Step struct {
	Name       string                 `toml:"name" json:"name"`
	Action     string                 `toml:"action" json:"action"`
	Config     map[string]interface{} `toml:"config" json:"config"`
	Store      string                 `toml:"store" json:"store"`
	Assertions []StepAssertion        `toml:"assertions" json:"assertions"`
}

// StepAssertion represents a condition to verify after a step
type StepAssertion struct {
	Field    string `toml:"field" json:"field"`
	Operator string `toml:"operator" json:"operator"`
	Value    string `toml:"value" json:"value"`
}

// ParseScript loads and parses a script file
func ParseScript(path string) (*Script, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read script file: %w", err)
	}

	var script Script
	ext := filepath.Ext(path)

	switch ext {
	case ".toml":
		if err := toml.Unmarshal(data, &script); err != nil {
			return nil, fmt.Errorf("failed to parse TOML script: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &script); err != nil {
			return nil, fmt.Errorf("failed to parse JSON script: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported script format: %s (use .toml or .json)", ext)
	}

	if err := validateScript(&script); err != nil {
		return nil, err
	}

	return &script, nil
}

func validateScript(s *Script) error {
	if len(s.Steps) == 0 {
		return fmt.Errorf("script must contain at least one step")
	}

	validActions := map[string]bool{
		"deploy": true,
		"call":   true,
		"fund":   true,
		"wait":   true,
		"assert": true,
	}

	for i, step := range s.Steps {
		if step.Action == "" {
			return fmt.Errorf("step %d: action is required", i)
		}
		if !validActions[step.Action] {
			return fmt.Errorf("step %d: unknown action '%s' (valid: deploy, call, fund, wait, assert)", i, step.Action)
		}
	}

	return nil
}
