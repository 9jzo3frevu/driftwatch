package drift

import "sort"

// MergeConfig controls how result sets are merged.
type MergeConfig struct {
	// DeduplicateByKey removes duplicate results with the same Key+Service pair,
	// keeping the entry with the highest severity.
	DeduplicateByKey bool
}

// DefaultMergeConfig returns a MergeConfig with sensible defaults.
func DefaultMergeConfig() MergeConfig {
	return MergeConfig{
		DeduplicateByKey: true,
	}
}

// MergeResults combines multiple slices of DriftResult into a single slice.
// When cfg.DeduplicateByKey is true, duplicate Key+Service pairs are collapsed,
// retaining the entry with the highest severity.
func MergeResults(cfg MergeConfig, sets ...[]DriftResult) []DriftResult {
	if len(sets) == 0 {
		return nil
	}

	total := 0
	for _, s := range sets {
		total += len(s)
	}

	merged := make([]DriftResult, 0, total)
	for _, s := range sets {
		merged = append(merged, s...)
	}

	if !cfg.DeduplicateByKey {
		return merged
	}

	type dedupKey struct {
		Key     string
		Service string
	}

	best := make(map[dedupKey]DriftResult, len(merged))
	order := make([]dedupKey, 0, len(merged))

	for _, r := range merged {
		k := dedupKey{Key: r.Key, Service: r.Service}
		if existing, ok := best[k]; !ok {
			best[k] = r
			order = append(order, k)
		} else if severityRank(r.Severity) > severityRank(existing.Severity) {
			best[k] = r
		}
	}

	out := make([]DriftResult, 0, len(order))
	for _, k := range order {
		out = append(out, best[k])
	}
	return out
}

// severityRank maps a severity string to a numeric rank for comparison.
func severityRank(s string) int {
	switch s {
	case "critical":
		return 3
	case "high":
		return 2
	case "medium":
		return 1
	case "low":
		return 0
	default:
		return -1
	}
}

// MergeAndSort merges result sets and sorts the output by severity (descending)
// then by key (ascending).
func MergeAndSort(cfg MergeConfig, sets ...[]DriftResult) []DriftResult {
	out := MergeResults(cfg, sets...)
	sort.SliceStable(out, func(i, j int) bool {
		ri, rj := severityRank(out[i].Severity), severityRank(out[j].Severity)
		if ri != rj {
			return ri > rj
		}
		return out[i].Key < out[j].Key
	})
	return out
}
