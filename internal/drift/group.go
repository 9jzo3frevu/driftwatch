package drift

import "sort"

// GroupKey defines how results are grouped.
type GroupKey string

const (
	GroupByService  GroupKey = "service"
	GroupBySeverity GroupKey = "severity"
	GroupByKey      GroupKey = "key"
)

// Group holds a named collection of drift results.
type Group struct {
	Name    string
	Results []Result
}

// GroupResults partitions results by the given key.
func GroupResults(results []Result, by GroupKey) []Group {
	buckets := make(map[string][]Result)
	for _, r := range results {
		var label string
		switch by {
		case GroupBySeverity:
			label = string(severityFor(r))
		case GroupByKey:
			if len(r.Key) > 0 {
				label = r.Key
			} else {
				label = "unknown"
			}
		default: // GroupByService
			if r.Service != "" {
				label = r.Service
			} else {
				label = "default"
			}
		}
		buckets[label] = append(buckets[label], r)
	}

	groups := make([]Group, 0, len(buckets))
	for name, res := range buckets {
		groups = append(groups, Group{Name: name, Results: res})
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}
