package drift

import "sort"

// RankOptions controls how results are ranked.
type RankOptions struct {
	By        string // "severity", "key", "service"
	Descending bool
}

// RankResults returns a sorted copy of results based on RankOptions.
func RankResults(results []DriftResult, opts RankOptions) []DriftResult {
	out := make([]DriftResult, len(results))
	copy(out, results)

	severityOrder := map[string]int{
		"low":      1,
		"medium":   2,
		"high":     3,
		"critical": 4,
	}

	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch opts.By {
		case "severity":
			si := severityOrder[out[i].Severity]
			sj := severityOrder[out[j].Severity]
			less = si < sj
		case "service":
			less = out[i].Service < out[j].Service
		default: // "key"
			less = out[i].Key < out[j].Key
		}
		if opts.Descending {
			return !less
		}
		return less
	})
	return out
}
