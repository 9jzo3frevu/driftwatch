package drift

import (
	"testing"
)

func aggregateResults() []Result {
	return []Result{
		{Key: "db.host", Service: "api", Severity: SeverityHigh, Drifted: true},
		{Key: "db.port", Service: "api", Severity: SeverityLow, Drifted: false},
		{Key: "cache.host", Service: "worker", Severity: SeverityCritical, Drifted: true},
		{Key: "cache.ttl", Service: "api", Severity: SeverityLow, Drifted: true},
		{Key: "queue.url", Service: "worker", Severity: SeverityHigh, Drifted: false},
	}
}

func TestAggregateResults_ByService(t *testing.T) {
	cfg := AggregateConfig{GroupBy: "service", IncludeSummary: false}
	groups := AggregateResults(aggregateResults(), cfg)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	// api has 3 results, worker has 2 — api should be first
	if groups[0].Label != "api" {
		t.Errorf("expected first group to be 'api', got %q", groups[0].Label)
	}
	if groups[0].Count() != 3 {
		t.Errorf("expected api group count 3, got %d", groups[0].Count())
	}
}

func TestAggregateResults_BySeverity(t *testing.T) {
	cfg := AggregateConfig{GroupBy: "severity", IncludeSummary: false}
	groups := AggregateResults(aggregateResults(), cfg)

	labels := make(map[string]int, len(groups))
	for _, g := range groups {
		labels[g.Label] = g.Count()
	}

	if labels[string(SeverityLow)] != 2 {
		t.Errorf("expected 2 low-severity results, got %d", labels[string(SeverityLow)])
	}
	if labels[string(SeverityHigh)] != 2 {
		t.Errorf("expected 2 high-severity results, got %d", labels[string(SeverityHigh)])
	}
}

func TestAggregateResults_ByPrefix(t *testing.T) {
	cfg := AggregateConfig{GroupBy: "prefix", IncludeSummary: false}
	groups := AggregateResults(aggregateResults(), cfg)

	counts := make(map[string]int)
	for _, g := range groups {
		counts[g.Label] = g.Count()
	}

	if counts["db"] != 2 {
		t.Errorf("expected db prefix count 2, got %d", counts["db"])
	}
	if counts["cache"] != 2 {
		t.Errorf("expected cache prefix count 2, got %d", counts["cache"])
	}
}

func TestAggregateResults_TopN(t *testing.T) {
	cfg := AggregateConfig{GroupBy: "service", TopN: 1, IncludeSummary: false}
	groups := AggregateResults(aggregateResults(), cfg)

	if len(groups) != 1 {
		t.Fatalf("expected 1 group with TopN=1, got %d", len(groups))
	}
}

func TestAggregateResults_IncludeSummary(t *testing.T) {
	cfg := AggregateConfig{GroupBy: "service", IncludeSummary: true}
	groups := AggregateResults(aggregateResults(), cfg)

	for _, g := range groups {
		if g.Summary == "" {
			t.Errorf("expected non-empty summary for group %q", g.Label)
		}
	}
}

func TestAggregateResults_Empty(t *testing.T) {
	cfg := DefaultAggregateConfig()
	groups := AggregateResults([]Result{}, cfg)

	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}
