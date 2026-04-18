package config

import "testing"

func boolPtr(b bool) *bool { return &b }

func TestDedupeRaw_Build_Disabled(t *testing.T) {
	r := DedupeRaw{Enabled: boolPtr(false)}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestDedupeRaw_Build_DefaultWindow(t *testing.T) {
	r := DedupeRaw{}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Error("expected enabled")
	}
	if cfg.Window != 3 {
		t.Errorf("expected default window 3, got %d", cfg.Window)
	}
}

func TestDedupeRaw_Build_CustomWindow(t *testing.T) {
	r := DedupeRaw{Window: 10}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Window != 10 {
		t.Errorf("expected window 10, got %d", cfg.Window)
	}
}

func TestDedupeRaw_Build_WindowTooLarge(t *testing.T) {
	r := DedupeRaw{Window: 200}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for window > 100")
	}
}

func TestDedupeRaw_Build_WindowNegative(t *testing.T) {
	r := DedupeRaw{Window: -1}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}
