package drift

import (
	"fmt"
	"sort"
	"strings"
)

// AggregateConfig controls how drift results are aggregated.
type AggregateConfig struct {
	// GroupBy determines the field used for aggregation: "service", "severity", or "prefix".
	GroupBy string
	// TopN limits the output to the N largest groups. 0 means no limit.
	TopN int
	// IncludeSummary attaches a human-readable summary line to each group.
	IncludeSummary bool
}

// AggregateGroup represents a collection of drift results sharing a common label.
type AggregateGroup struct {
	Label   string
	Results []Result
	Summary string
}

// Count returns the number of results in the group.
func (g AggregateGroup) Count() int { return len(g.Results) }

// DefaultAggregateConfig returns a sensible default configuration.
func DefaultAggregateConfig() AggregateConfig {
	return AggregateConfig{
		GroupBy:        "service",
		TopN:           0,
		IncludeSummary: true,
	}
}

// AggregateResults groups the provided results according to cfg and returns
// the groups sorted by descending count.
func AggregateResults(results []Result, cfg AggregateConfig) []AggregateGroup {
	index := make(map[string][]Result)

	for _, r := range results {
		label := aggregateLabel(r, cfg.GroupBy)
		index[label] = append(index[label], r)
	}

	groups := make([]AggregateGroup, 0, len(index))
	for label, rs := range index {
		g := AggregateGroup{Label: label, Results: rs}
		if cfg.IncludeSummary {
			g.Summary = buildSummary(label, rs)
		}
		groups = append(groups, g)
	}

	sort.Slice(groups, func(i, j int) bool {
		if groups[i].Count() != groups[j].Count() {
			return groups[i].Count() > groups[j].Count()
		}
		return groups[i].Label < groups[j].Label
	})

	if cfg.TopN > 0 && len(groups) > cfg.TopN {
		groups = groups[:cfg.TopN]
	}

	return groups
}

func aggregateLabel(r Result, groupBy string) string {
	switch groupBy {
	case "severity":
		return string(r.Severity)
	case "prefix":
		parts := strings.SplitN(r.Key, ".", 2)
		return parts[0]
	default: // "service"
		if r.Service != "" {
			return r.Service
		}
		return "(unknown)"
	}
}

func buildSummary(label string, results []Result) string {
	drifted := 0
	for _, r := range results {
		if r.Drifted {
			drifted++
		}
	}
	return fmt.Sprintf("%s: %d result(s), %d drifted", label, len(results), drifted)
}
