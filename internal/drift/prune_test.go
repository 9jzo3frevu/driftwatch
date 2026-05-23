package drift

import (
	"testing"
	"time"
)

func pruneResults(n int, drifted bool, age time.Duration) []DriftResult {
	now := time.Now()
	out := make([]DriftResult, n)
	for i := range out {
		out[i] = DriftResult{
			Key:         "key",
			Drifted:     drifted,
			DetectedAt:  now.Add(-age),
		}
	}
	return out
}

func TestPruner_OnlyDrifted_RemovesNonDrifted(t *testing.T) {
	cfg := PruneConfig{OnlyDrifted: true}
	p := NewPruner(cfg)
	input := []DriftResult{
		{Key: "a", Drifted: true},
		{Key: "b", Drifted: false},
		{Key: "c", Drifted: true},
	}
	out := p.Prune(input)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestPruner_MaxAge_RemovesOld(t *testing.T) {
	cfg := PruneConfig{MaxAge: time.Hour}
	p := NewPruner(cfg)
	now := time.Now()
	input := []DriftResult{
		{Key: "fresh", Drifted: true, DetectedAt: now.Add(-30 * time.Minute)},
		{Key: "stale", Drifted: true, DetectedAt: now.Add(-2 * time.Hour)},
	}
	out := p.Prune(input)
	if len(out) != 1 || out[0].Key != "fresh" {
		t.Fatalf("expected only fresh result, got %v", out)
	}
}

func TestPruner_MaxResults_Caps(t *testing.T) {
	cfg := PruneConfig{MaxResults: 3}
	p := NewPruner(cfg)
	input := pruneResults(10, true, 0)
	out := p.Prune(input)
	if len(out) != 3 {
		t.Fatalf("expected 3, got %d", len(out))
	}
}

func TestPruner_NoConfig_KeepsAll(t *testing.T) {
	p := NewPruner(PruneConfig{})
	input := pruneResults(5, false, 0)
	out := p.Prune(input)
	if len(out) != 5 {
		t.Fatalf("expected 5, got %d", len(out))
	}
}

func TestPruner_ZeroDetectedAt_SkipsAgeCheck(t *testing.T) {
	cfg := PruneConfig{MaxAge: time.Minute}
	p := NewPruner(cfg)
	input := []DriftResult{{Key: "no-time", Drifted: true}}
	out := p.Prune(input)
	if len(out) != 1 {
		t.Fatalf("expected result with zero DetectedAt to be kept, got %d", len(out))
	}
}

func TestDefaultPruneConfig(t *testing.T) {
	cfg := DefaultPruneConfig()
	if cfg.MaxAge != 24*time.Hour {
		t.Errorf("unexpected MaxAge: %v", cfg.MaxAge)
	}
	if cfg.MaxResults != 500 {
		t.Errorf("unexpected MaxResults: %d", cfg.MaxResults)
	}
}
