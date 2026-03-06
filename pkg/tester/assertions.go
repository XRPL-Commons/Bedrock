package tester

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/xrpl-commons/bedrock/pkg/caller"
)

// AssertionResult holds the outcome of an assertion check
type AssertionResult struct {
	Passed   bool
	Message  string
	Expected string
	Actual   string
}

// RunAssertions validates a call result against a set of assertions
func RunAssertions(result *caller.CallResult, assertions []Assertion) []AssertionResult {
	var results []AssertionResult
	for _, a := range assertions {
		results = append(results, runAssertion(result, a))
	}
	return results
}

func runAssertion(result *caller.CallResult, a Assertion) AssertionResult {
	switch a.Type {
	case AssertReturnCode:
		return assertReturnCode(result, a)
	case AssertReturnValue:
		return assertReturnValue(result, a)
	case AssertTxSuccess:
		return assertTxSuccess(result)
	case AssertTxFailure:
		return assertTxFailure(result, a)
	case AssertEvent:
		return assertEvent(result, a)
	case AssertState:
		return assertState(result, a)
	case AssertGasBelow:
		return assertGasBelow(result, a)
	default:
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("unknown assertion type: %s", a.Type),
		}
	}
}

func assertReturnCode(result *caller.CallResult, a Assertion) AssertionResult {
	expected, ok := toInt64(a.Expected)
	if !ok {
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("invalid expected return code: %v", a.Expected),
		}
	}

	actual := int64(result.ReturnCode)
	return AssertionResult{
		Passed:   actual == expected,
		Message:  "return code matches",
		Expected: fmt.Sprintf("%d", expected),
		Actual:   fmt.Sprintf("%d", actual),
	}
}

func assertReturnValue(result *caller.CallResult, a Assertion) AssertionResult {
	expected := fmt.Sprintf("%v", a.Expected)
	actual := result.ReturnValue

	return AssertionResult{
		Passed:   strings.EqualFold(actual, expected),
		Message:  "return value matches",
		Expected: expected,
		Actual:   actual,
	}
}

func assertTxSuccess(result *caller.CallResult) AssertionResult {
	passed := result.TransactionResult == "tesSUCCESS"
	return AssertionResult{
		Passed:   passed,
		Message:  "transaction succeeded",
		Expected: "tesSUCCESS",
		Actual:   result.TransactionResult,
	}
}

func assertTxFailure(result *caller.CallResult, a Assertion) AssertionResult {
	if a.Expected != nil {
		expectedCode := fmt.Sprintf("%v", a.Expected)
		return AssertionResult{
			Passed:   result.TransactionResult == expectedCode,
			Message:  fmt.Sprintf("transaction failed with %s", expectedCode),
			Expected: expectedCode,
			Actual:   result.TransactionResult,
		}
	}

	passed := result.TransactionResult != "tesSUCCESS"
	return AssertionResult{
		Passed:   passed,
		Message:  "transaction failed",
		Expected: "not tesSUCCESS",
		Actual:   result.TransactionResult,
	}
}

func assertEvent(result *caller.CallResult, a Assertion) AssertionResult {
	// Events are in the transaction metadata
	if result.Meta == nil {
		return AssertionResult{
			Passed:  false,
			Message: "no transaction metadata available",
		}
	}

	expectedEvent, ok := a.Expected.(map[string]interface{})
	if !ok {
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("invalid event assertion format: %v", a.Expected),
		}
	}

	// Look for ContractEvent entries in metadata
	events := extractEvents(result.Meta)
	eventType, _ := expectedEvent["type"].(string)

	for _, event := range events {
		if eventType != "" {
			if et, ok := event["type"].(string); ok && et == eventType {
				return AssertionResult{
					Passed:  true,
					Message: fmt.Sprintf("event '%s' found", eventType),
				}
			}
		}
	}

	return AssertionResult{
		Passed:   false,
		Message:  fmt.Sprintf("event '%s' not found in %d events", eventType, len(events)),
		Expected: fmt.Sprintf("%v", expectedEvent),
		Actual:   fmt.Sprintf("%d events emitted", len(events)),
	}
}

func assertState(result *caller.CallResult, a Assertion) AssertionResult {
	// State assertions check ContractData in metadata
	if result.Meta == nil {
		return AssertionResult{
			Passed:  false,
			Message: "no transaction metadata available for state check",
		}
	}

	if a.Field == "" {
		return AssertionResult{
			Passed:  false,
			Message: "state assertion requires a 'field' specifier",
		}
	}

	// Look through AffectedNodes in metadata for state changes
	expected := fmt.Sprintf("%v", a.Expected)
	actual := findStateValue(result.Meta, a.Field)

	return AssertionResult{
		Passed:   actual == expected,
		Message:  fmt.Sprintf("state field '%s'", a.Field),
		Expected: expected,
		Actual:   actual,
	}
}

func assertGasBelow(result *caller.CallResult, a Assertion) AssertionResult {
	threshold, ok := toInt64(a.Expected)
	if !ok {
		return AssertionResult{
			Passed:  false,
			Message: fmt.Sprintf("invalid gas threshold: %v", a.Expected),
		}
	}

	return AssertionResult{
		Passed:   result.GasUsed < threshold,
		Message:  fmt.Sprintf("gas used below %d", threshold),
		Expected: fmt.Sprintf("< %d", threshold),
		Actual:   fmt.Sprintf("%d", result.GasUsed),
	}
}

// extractEvents pulls ContractEvent entries from transaction metadata
func extractEvents(meta map[string]interface{}) []map[string]interface{} {
	var events []map[string]interface{}

	// Navigate metadata structure: AffectedNodes or similar
	affectedNodes, ok := meta["AffectedNodes"]
	if !ok {
		return events
	}

	nodes, ok := affectedNodes.([]interface{})
	if !ok {
		return events
	}

	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		for _, v := range nodeMap {
			innerMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			if ledgerEntryType, ok := innerMap["LedgerEntryType"].(string); ok {
				if ledgerEntryType == "ContractEvent" {
					events = append(events, innerMap)
				}
			}
		}
	}

	return events
}

// findStateValue searches metadata for a specific state field value
func findStateValue(meta map[string]interface{}, field string) string {
	// Search through AffectedNodes for ContractData modifications
	affectedNodes, ok := meta["AffectedNodes"]
	if !ok {
		return ""
	}

	nodes, ok := affectedNodes.([]interface{})
	if !ok {
		return ""
	}

	for _, node := range nodes {
		nodeMap, ok := node.(map[string]interface{})
		if !ok {
			continue
		}

		for _, v := range nodeMap {
			innerMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			if ledgerEntryType, ok := innerMap["LedgerEntryType"].(string); ok {
				if ledgerEntryType == "ContractData" {
					if finalFields, ok := innerMap["FinalFields"].(map[string]interface{}); ok {
						if val, ok := finalFields[field]; ok {
							return fmt.Sprintf("%v", val)
						}
					}
					if newFields, ok := innerMap["NewFields"].(map[string]interface{}); ok {
						if val, ok := newFields[field]; ok {
							return fmt.Sprintf("%v", val)
						}
					}
				}
			}
		}
	}

	return ""
}

// toInt64 converts an interface{} to int64
func toInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	case json.Number:
		n, err := val.Int64()
		return n, err == nil
	case string:
		var n int64
		_, err := fmt.Sscanf(val, "%d", &n)
		return n, err == nil
	default:
		return 0, false
	}
}
