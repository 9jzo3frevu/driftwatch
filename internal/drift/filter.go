package drift

import "strings"

// Filter controls which drift results are surfaced.
type Filter struct {
	IncludeKeys []string
	ExcludeKeys []string
	MinSeverity string
}

var severityRank = map[string]int{
	"low":      1,
	"medium":   2,
	"high":      3,
	"critical": 4,
}

// Apply returns only the results that pass the filter criteria.
func (f *Filter) Apply(results []DriftResult) []DriftResult {
	var out []DriftResult
	for _, r := range results {
		if !f.keyAllowed(r.Key) {
			continue
		}
		if !f.severityAllowed(r) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func (f *Filter) keyAllowed(key string) bool {
	for _, ex := range f.ExcludeKeys {
		if strings.EqualFold(key, ex) {
			return false
		}
	}
	if len(f.IncludeKeys) == 0 {
		return true
	}
	for _, inc := range f.IncludeKeys {
		if strings.EqualFold(key, inc) {
			return true
		}
	}
	return false
}

func (f *Filter) severityAllowed(r DriftResult) bool {
	if f.MinSeverity == "" {
		return true
	}
	min := severityRank[strings.ToLower(f.MinSeverity)]
	actual := severityRank[strings.ToLower(severityFor(r))]
	return actual >= min
}
