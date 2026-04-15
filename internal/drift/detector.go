package drift

import (
	"fmt"
	"reflect"
)

// State represents a key-value configuration state for a service.
type State map[string]interface{}

// DriftResult holds the result of comparing two states.
type DriftResult struct {
	ServiceName string
	Drifted     bool
	Changes     []Change
}

// Change describes a single configuration difference.
type Change struct {
	Key      string
	Expected interface{}
	Actual   interface{}
	Kind     ChangeKind
}

// ChangeKind categorizes the type of drift detected.
type ChangeKind string

const (
	ChangeAdded    ChangeKind = "added"
	ChangeRemoved  ChangeKind = "removed"
	ChangeModified ChangeKind = "modified"
)

// Detect compares the declared (expected) state against the live (actual)
// state for a named service and returns a DriftResult.
func Detect(serviceName string, expected, actual State) DriftResult {
	result := DriftResult{
		ServiceName: serviceName,
	}

	// Check for modified or removed keys.
	for key, expVal := range expected {
		actVal, exists := actual[key]
		if !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Expected: expVal,
				Actual:   nil,
				Kind:     ChangeRemoved,
			})
			continue
		}
		if !reflect.DeepEqual(expVal, actVal) {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Expected: expVal,
				Actual:   actVal,
				Kind:     ChangeModified,
			})
		}
	}

	// Check for keys present in actual but not in expected.
	for key, actVal := range actual {
		if _, exists := expected[key]; !exists {
			result.Changes = append(result.Changes, Change{
				Key:      key,
				Expected: nil,
				Actual:   actVal,
				Kind:     ChangeAdded,
			})
		}
	}

	result.Drifted = len(result.Changes) > 0
	return result
}

// Summary returns a human-readable summary of the drift result.
func (r DriftResult) Summary() string {
	if !r.Drifted {
		return fmt.Sprintf("[OK] %s: no drift detected", r.ServiceName)
	}
	return fmt.Sprintf("[DRIFT] %s: %d change(s) detected", r.ServiceName, len(r.Changes))
}
