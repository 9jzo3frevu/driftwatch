package drift

import (
	"testing"
	"time"
)

func capResults(n int, age time.Duration) []DriftResult {
	now := time.Now()
	out := make([]DriftResult, n)
	for i := 0; i < n; i++ {
		out[i] = DriftResult{
			Key:        ptrStr("key"),
			Drifted:    true,
			DetectedAt: now.Add(-age * time.Duration(i+1)),
		}
	}
	return out
}

func TestCap_Disabled_ReturnsAll(t *testing.T) {
	cfg := CapConfig{Enabled: false, MaxResults: 2}
	results := capResults(5, time.Minute)
	got := Cap(results, cfg)
	if len(got) != 5 {
		t.Fatalf("expected 5, got %d", len(got))
	}
}

func TestCap_MaxResults(t *testing.T) {
	cfg := CapConfig{Enabled: true, MaxResults: 3}
	results := capResults(10, time.Minute)
	got := Cap(results, cfg)
	if len(got) != 3 {
		t.Fatalf("expected 3, got %d", len(got))
	}
}

func TestCap_MaxAge_RemovesOld(t *testing.T) {
	cfg := CapConfig{Enabled: true, MaxResults: 100, MaxAge: 30 * time.Minute}
	// first 2 results are within 30 min, rest are older
	results := capResults(5, 10*time.Minute)
	// results[0] -> 10m old, results[1] -> 20m old, results[2] -> 30m old (boundary), results[3] -> 40m, results[4] -> 50m
	got := Cap(results, cfg)
	if len(got) != 2 {
		t.Fatalf("expected 2 within age window, got %d", len(got))
	}
}

func TestCap_Empty_ReturnsEmpty(t *testing.T) {
	cfg := DefaultCapConfig()
	got := Cap(nil, cfg)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCap_DefaultConfig(t *testing.T) {
	cfg := DefaultCapConfig()
	if !cfg.Enabled {
		t.Fatal("expected default config to be enabled")
	}
	if cfg.MaxResults != 500 {
		t.Fatalf("expected MaxResults 500, got %d", cfg.MaxResults)
	}
	if cfg.MaxAge != 24*time.Hour {
		t.Fatalf("expected MaxAge 24h, got %v", cfg.MaxAge)
	}
}

func TestCap_MaxResultsAndAge_BothApply(t *testing.T) {
	cfg := CapConfig{Enabled: true, MaxResults: 2, MaxAge: 25 * time.Minute}
	results := capResults(10, 5*time.Minute)
	// ages: 5m, 10m, 15m, 20m, 25m(boundary excluded), ...
	// within age: indices 0,1,2,3 (5,10,15,20 min) -> 4 results, but capped at 2
	got := Cap(results, cfg)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
