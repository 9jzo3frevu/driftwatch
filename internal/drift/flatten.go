package drift

import (
	"fmt"
	"sort"
)

// FlattenConfig controls how results are flattened into a single ordered slice.
type FlattenConfig struct {
	// SortBy controls the field used for ordering: "key", "service", "severity".
	SortBy string
	// Descending reverses the sort order when true.
	Descending bool
	// OnlyDrifted drops non-drifted results when true.
	OnlyDrifted bool
}

// DefaultFlattenConfig returns a FlattenConfig with sensible defaults.
func DefaultFlattenConfig() FlattenConfig {
	return FlattenConfig{
		SortBy:      "key",
		Descending:  false,
		OnlyDrifted: false,
	}
}

// FlattenResults merges grouped or nested result sets into a single flat slice
// and applies optional sorting and filtering.
func FlattenResults(groups map[string][]DriftResult, cfg FlattenConfig) []DriftResult {
	if len(groups) == 0 {
		return []DriftResult{}
	}

	// Collect keys in deterministic order.
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var flat []DriftResult
	for _, k := range keys {
		for _, r := range groups[k] {
			if cfg.OnlyDrifted && !r.Drifted {
				continue
			}
			flat = append(flat, r)
		}
	}

	if len(flat) == 0 {
		return []DriftResult{}
	}

	sortFlat(flat, cfg.SortBy, cfg.Descending)
	return flat
}

func sortFlat(results []DriftResult, by string, desc bool) {
	sort.SliceStable(results, func(i, j int) bool {
		var less bool
		switch by {
		case "service":
			less = results[i].Service < results[j].Service
		case "severity":
			less = severityRank(results[i].Severity) < severityRank(results[j].Severity)
		default: // "key"
			less = results[i].Key < results[j].Key
		}
		if desc {
			return !less
		}
		return less
	})
}

// FlattenSummary returns a human-readable summary line for a flattened set.
func FlattenSummary(results []DriftResult) string {
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}
	return fmt.Sprintf("%d result(s), %d drifted", len(results), drifted)
}
