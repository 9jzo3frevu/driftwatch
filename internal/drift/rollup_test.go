package drift

import (
	"testing"
)

func rollupResults() []DriftResult {
	return []DriftResult{
		{Key: "db.host", Service: "api", Severity: SeverityHigh},
		{Key: "db.port", Service: "api", Severity: SeverityLow},
		{Key: "cache.ttl", Service: "worker", Severity: SeverityLow},
		{Key: "cache.host", Service: "worker", Severity: SeverityCritical},
		{Key: "auth.secret", Service: "api", Severity: SeverityCritical},
	}
}

func TestRollupResults_ByService(t *testing.T) {
	cfg := RollupConfig{GroupBy: "service", MaxPerGroup: 10}
	groups := RollupResults(rollupResults(), cfg)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Label != "api" {
		t.Errorf("expected first group 'api', got %q", groups[0].Label)
	}
	if groups[0].Count != 3 {
		t.Errorf("expected api count 3, got %d", groups[0].Count)
	}
	if groups[1].Label != "worker" {
		t.Errorf("expected second group 'worker', got %q", groups[1].Label)
	}
}

func TestRollupResults_BySeverity(t *testing.T) {
	cfg := RollupConfig{GroupBy: "severity", MaxPerGroup: 10}
	groups := RollupResults(rollupResults(), cfg)

	labels := make(map[string]int)
	for _, g := range groups {
		labels[g.Label] = g.Count
	}

	if labels[string(SeverityHigh)] != 1 {
		t.Errorf("expected 1 high, got %d", labels[string(SeverityHigh)])
	}
	if labels[string(SeverityLow)] != 2 {
		t.Errorf("expected 2 low, got %d", labels[string(SeverityLow)])
	}
	if labels[string(SeverityCritical)] != 2 {
		t.Errorf("expected 2 critical, got %d", labels[string(SeverityCritical)])
	}
}

func TestRollupResults_ByKeyPrefix(t *testing.T) {
	cfg := RollupConfig{GroupBy: "key_prefix", MaxPerGroup: 10}
	groups := RollupResults(rollupResults(), cfg)

	labels := make(map[string]int)
	for _, g := range groups {
		labels[g.Label] = g.Count
	}

	if labels["db"] != 2 {
		t.Errorf("expected 2 db keys, got %d", labels["db"])
	}
	if labels["cache"] != 2 {
		t.Errorf("expected 2 cache keys, got %d", labels["cache"])
	}
	if labels["auth"] != 1 {
		t.Errorf("expected 1 auth key, got %d", labels["auth"])
	}
}

func TestRollupResults_MaxPerGroup(t *testing.T) {
	cfg := RollupConfig{GroupBy: "service", MaxPerGroup: 2}
	groups := RollupResults(rollupResults(), cfg)

	for _, g := range groups {
		if len(g.Results) > 2 {
			t.Errorf("group %q has %d results, expected <= 2", g.Label, len(g.Results))
		}
		if g.Label == "api" && g.Count != 3 {
			t.Errorf("api total count should be 3, got %d", g.Count)
		}
	}
}

func TestRollupResults_Empty(t *testing.T) {
	groups := RollupResults(nil, DefaultRollupConfig())
	if groups != nil {
		t.Errorf("expected nil for empty input, got %v", groups)
	}
}

func TestRollupResults_DefaultMaxPerGroup(t *testing.T) {
	cfg := RollupConfig{GroupBy: "service", MaxPerGroup: 0}
	groups := RollupResults(rollupResults(), cfg)
	if len(groups) == 0 {
		t.Error("expected groups with zero MaxPerGroup to use default")
	}
}
