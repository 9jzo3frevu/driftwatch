package drift

import "testing"

func dedupeResults() []DriftResult {
	return []DriftResult{
		{Service: "api", Key: "timeout", Declared: "30s", Live: ptrStr("60s")},
		{Service: "api", Key: "replicas", Declared: "3", Live: ptrStr("2")},
		{Service: "worker", Key: "memory", Declared: "512Mi", Live: nil},
	}
}

func TestDeduplicator_NoPreviousHistory(t *testing.T) {
	d := NewDeduplicator(DedupeConfig{Window: 2})
	out := d.Filter(dedupeResults(), nil)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestDeduplicator_RemovesSeen(t *testing.T) {
	d := NewDeduplicator(DedupeConfig{Window: 1})
	prev := dedupeResults()[:2]
	out := d.Filter(dedupeResults(), [][]DriftResult{prev})
	if len(out) != 1 {
		t.Fatalf("expected 1 new result, got %d", len(out))
	}
	if out[0].Key != "memory" {
		t.Errorf("unexpected key %q", out[0].Key)
	}
}

func TestDeduplicator_WindowLimitsLookback(t *testing.T) {
	d := NewDeduplicator(DedupeConfig{Window: 1})
	old := dedupeResults() // all three seen two runs ago
	recent := []DriftResult{} // nothing seen last run
	out := d.Filter(dedupeResults(), [][]DriftResult{old, recent})
	// window=1 so only 'recent' (empty) counts → all three pass through
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestDeduplicator_DefaultWindowFloor(t *testing.T) {
	d := NewDeduplicator(DedupeConfig{Window: 0}) // should default to 1
	if d.cfg.Window != 1 {
		t.Errorf("expected window=1, got %d", d.cfg.Window)
	}
}

func TestDeduplicator_AllNew(t *testing.T) {
	d := NewDeduplicator(DedupeConfig{Window: 3})
	prev := []DriftResult{{Service: "other", Key: "x", Declared: "a", Live: ptrStr("b")}}
	out := d.Filter(dedupeResults(), [][]DriftResult{prev})
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}
