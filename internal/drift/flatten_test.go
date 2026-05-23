package drift

import (
	"testing"
)

func flattenGroups() map[string][]DriftResult {
	return map[string][]DriftResult{
		"svc-a": {
			{Key: "db.host", Service: "svc-a", Drifted: true, Severity: "high"},
			{Key: "db.port", Service: "svc-a", Drifted: false, Severity: "low"},
		},
		"svc-b": {
			{Key: "app.timeout", Service: "svc-b", Drifted: true, Severity: "critical"},
			{Key: "app.retries", Service: "svc-b", Drifted: false, Severity: "low"},
		},
	}
}

func TestFlattenResults_AllResults(t *testing.T) {
	cfg := DefaultFlattenConfig()
	out := FlattenResults(flattenGroups(), cfg)
	if len(out) != 4 {
		t.Fatalf("expected 4 results, got %d", len(out))
	}
}

func TestFlattenResults_OnlyDrifted(t *testing.T) {
	cfg := DefaultFlattenConfig()
	cfg.OnlyDrifted = true
	out := FlattenResults(flattenGroups(), cfg)
	if len(out) != 2 {
		t.Fatalf("expected 2 drifted results, got %d", len(out))
	}
	for _, r := range out {
		if !r.Drifted {
			t.Errorf("expected only drifted results, got non-drifted key=%s", r.Key)
		}
	}
}

func TestFlattenResults_SortByKey(t *testing.T) {
	cfg := DefaultFlattenConfig()
	out := FlattenResults(flattenGroups(), cfg)
	for i := 1; i < len(out); i++ {
		if out[i].Key < out[i-1].Key {
			t.Errorf("results not sorted by key: %s before %s", out[i-1].Key, out[i].Key)
		}
	}
}

func TestFlattenResults_SortByService_Descending(t *testing.T) {
	cfg := FlattenConfig{SortBy: "service", Descending: true}
	out := FlattenResults(flattenGroups(), cfg)
	if len(out) == 0 {
		t.Fatal("expected results")
	}
	for i := 1; i < len(out); i++ {
		if out[i].Service > out[i-1].Service {
			t.Errorf("results not sorted descending by service")
		}
	}
}

func TestFlattenResults_EmptyGroups(t *testing.T) {
	out := FlattenResults(map[string][]DriftResult{}, DefaultFlattenConfig())
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d", len(out))
	}
}

func TestFlattenResults_NilGroups(t *testing.T) {
	out := FlattenResults(nil, DefaultFlattenConfig())
	if len(out) != 0 {
		t.Fatalf("expected empty slice for nil input, got %d", len(out))
	}
}

func TestFlattenSummary(t *testing.T) {
	results := []DriftResult{
		{Key: "a", Drifted: true},
		{Key: "b", Drifted: false},
		{Key: "c", Drifted: true},
	}
	summary := FlattenSummary(results)
	expected := "3 result(s), 2 drifted"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
