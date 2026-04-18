package drift

import "fmt"

// CompareMode controls how live keys are compared against declared keys.
type CompareMode int

const (
	CompareModeExact  CompareMode = iota // all keys must match exactly
	CompareModeSubset                    // declared keys must exist in live; extra live keys are ignored
)

// CompareOptions configures a Comparison run.
type CompareOptions struct {
	Mode      CompareMode
	Service   string
	Tolerance float64 // allowed fractional numeric drift (0 = exact)
}

// CompareResult holds the outcome of a single key comparison.
type CompareResult struct {
	Key      string
	Declared string
	Live     string
	Drifted  bool
	Reason   string
}

// Compare compares declared vs live key-value maps using the provided options
// and returns a slice of CompareResult for every evaluated key.
func Compare(declared, live map[string]string, opts CompareOptions) []CompareResult {
	var results []CompareResult

	for k, dv := range declared {
		lv, ok := live[k]
		if !ok {
			results = append(results, CompareResult{
				Key:      k,
				Declared: dv,
				Drifted:  true,
				Reason:   "key missing from live",
			})
			continue
		}
		if dv == lv {
			results = append(results, CompareResult{Key: k, Declared: dv, Live: lv})
			continue
		}
		results = append(results, CompareResult{
			Key:      k,
			Declared: dv,
			Live:     lv,
			Drifted:  true,
			Reason:   fmt.Sprintf("value mismatch: declared=%q live=%q", dv, lv),
		})
	}

	if opts.Mode == CompareModeExact {
		for k, lv := range live {
			if _, found := declared[k]; !found {
				results = append(results, CompareResult{
					Key:     k,
					Live:    lv,
					Drifted: true,
					Reason:  "extra key in live not present in declared",
				})
			}
		}
	}

	return results
}
