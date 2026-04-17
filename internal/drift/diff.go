package drift

import "fmt"

// DiffMode controls how two snapshots are compared.
type DiffMode int

const (
	DiffModeExact DiffMode = iota
	DiffModeSubset // only keys present in declared are checked
)

// Diff compares a declared map against a live map and returns DriftResults.
func Diff(declared, live map[string]string, mode DiffMode) []DriftResult {
	var results []DriftResult

	for key, declaredVal := range declared {
		liveVal, ok := live[key]
		if !ok {
			results = append(results, DriftResult{
				Key:      key,
				Expected: ptrStr(declaredVal),
				Actual:   nil,
				Drifted:  true,
				Reason:   fmt.Sprintf("key %q missing from live state", key),
			})
			continue
		}
		if liveVal != declaredVal {
			results = append(results, DriftResult{
				Key:      key,
				Expected: ptrStr(declaredVal),
				Actual:   ptrStr(liveVal),
				Drifted:  true,
				Reason:   fmt.Sprintf("value mismatch for key %q", key),
			})
		}
	}

	if mode == DiffModeExact {
		for key, liveVal := range live {
			if _, ok := declared[key]; !ok {
				results = append(results, DriftResult{
					Key:     key,
					Actual:  ptrStr(liveVal),
					Drifted: true,
					Reason:  fmt.Sprintf("key %q present in live state but not declared", key),
				})
			}
		}
	}

	return results
}
