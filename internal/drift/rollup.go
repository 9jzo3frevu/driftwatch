package drift

import (
	"fmt"
	"sort"
	"strings"
)

// RollupConfig controls how drift results are rolled up into aggregate summaries.
type RollupConfig struct {
	// GroupBy determines the field to aggregate on: "service", "severity", or "key_prefix".
	GroupBy string
	// MaxPerGroup caps the number of individual results included per group.
	MaxPerGroup int
}

// RollupGroup represents an aggregated collection of drift results under a common label.
type RollupGroup struct {
	Label   string
	Count   int
	Results []DriftResult
}

// DefaultRollupConfig returns a RollupConfig with sensible defaults.
func DefaultRollupConfig() RollupConfig {
	return RollupConfig{
		GroupBy:     "service",
		MaxPerGroup: 10,
	}
}

// RollupResults aggregates drift results into groups according to cfg.
// Results within each group are sorted by key ascending.
func RollupResults(results []DriftResult, cfg RollupConfig) []RollupGroup {
	if len(results) == 0 {
		return nil
	}

	if cfg.MaxPerGroup <= 0 {
		cfg.MaxPerGroup = DefaultRollupConfig().MaxPerGroup
	}

	buckets := make(map[string][]DriftResult)
	for _, r := range results {
		label := labelFor(r, cfg.GroupBy)
		buckets[label] = append(buckets[label], r)
	}

	groups := make([]RollupGroup, 0, len(buckets))
	for label, items := range buckets {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Key < items[j].Key
		})
		capped := items
		if len(capped) > cfg.MaxPerGroup {
			capped = capped[:cfg.MaxPerGroup]
		}
		groups = append(groups, RollupGroup{
			Label:   label,
			Count:   len(items),
			Results: capped,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Label < groups[j].Label
	})

	return groups
}

func labelFor(r DriftResult, groupBy string) string {
	switch strings.ToLower(groupBy) {
	case "severity":
		return string(r.Severity)
	case "key_prefix":
		parts := strings.SplitN(r.Key, ".", 2)
		return parts[0]
	default:
		if r.Service != "" {
			return r.Service
		}
		return fmt.Sprintf("unknown")
	}
}
