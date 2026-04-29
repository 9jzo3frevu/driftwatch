package drift

import (
	"testing"
)

func mergeResults() []DriftResult {
	return []DriftResult{
		{Key: "app.port", Service: "api", Severity: "low", Drifted: true},
		{Key: "app.host", Service: "api", Severity: "high", Drifted: true},
		{Key: "db.url", Service: "db", Severity: "critical", Drifted: true},
	}
}

func TestMergeResults_CombinesSets(t *testing.T) {
	a := []DriftResult{{Key: "a", Service: "svc", Severity: "low", Drifted: true}}
	b := []DriftResult{{Key: "b", Service: "svc", Severity: "high", Drifted: true}}

	out := MergeResults(DefaultMergeConfig(), a, b)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestMergeResults_Empty(t *testing.T) {
	out := MergeResults(DefaultMergeConfig())
	if out != nil {
		t.Fatalf("expected nil for empty input, got %v", out)
	}
}

func TestMergeResults_DeduplicatesKeepHighest(t *testing.T) {
	a := []DriftResult{{Key: "x", Service: "svc", Severity: "low", Drifted: true}}
	b := []DriftResult{{Key: "x", Service: "svc", Severity: "critical", Drifted: true}}

	out := MergeResults(DefaultMergeConfig(), a, b)
	if len(out) != 1 {
		t.Fatalf("expected 1 deduplicated result, got %d", len(out))
	}
	if out[0].Severity != "critical" {
		t.Errorf("expected severity 'critical', got %q", out[0].Severity)
	}
}

func TestMergeResults_NoDedupe_KeepsBoth(t *testing.T) {
	cfg := MergeConfig{DeduplicateByKey: false}
	a := []DriftResult{{Key: "x", Service: "svc", Severity: "low", Drifted: true}}
	b := []DriftResult{{Key: "x", Service: "svc", Severity: "critical", Drifted: true}}

	out := MergeResults(cfg, a, b)
	if len(out) != 2 {
		t.Fatalf("expected 2 results without dedup, got %d", len(out))
	}
}

func TestMergeResults_DifferentServices_NotDeduped(t *testing.T) {
	a := []DriftResult{{Key: "x", Service: "svc-a", Severity: "low", Drifted: true}}
	b := []DriftResult{{Key: "x", Service: "svc-b", Severity: "high", Drifted: true}}

	out := MergeResults(DefaultMergeConfig(), a, b)
	if len(out) != 2 {
		t.Fatalf("expected 2 results for different services, got %d", len(out))
	}
}

func TestMergeAndSort_OrdersBySeverityDesc(t *testing.T) {
	sets := [][]DriftResult{mergeResults()}
	out := MergeAndSort(DefaultMergeConfig(), sets...)

	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
	if out[0].Severity != "critical" {
		t.Errorf("first result should be critical, got %q", out[0].Severity)
	}
	if out[1].Severity != "high" {
		t.Errorf("second result should be high, got %q", out[1].Severity)
	}
	if out[2].Severity != "low" {
		t.Errorf("third result should be low, got %q", out[2].Severity)
	}
}

func TestMergeAndSort_TiesSortedByKey(t *testing.T) {
	set := []DriftResult{
		{Key: "z.config", Service: "svc", Severity: "high", Drifted: true},
		{Key: "a.config", Service: "svc", Severity: "high", Drifted: true},
	}
	out := MergeAndSort(DefaultMergeConfig(), set)
	if out[0].Key != "a.config" {
		t.Errorf("expected 'a.config' first, got %q", out[0].Key)
	}
}
