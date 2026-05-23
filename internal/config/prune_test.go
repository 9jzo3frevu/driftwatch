package config

import (
	"testing"
	"time"
)

func boolPtrPrune(b bool) *bool { return &b }

func TestPruneRaw_Build_Disabled(t *testing.T) {
	f := false
	r := &PruneRaw{Enabled: &f}
	cfg, err := r.Build()
	if err != nil || cfg != nil {
		t.Fatalf("expected nil config when disabled, got %v %v", cfg, err)
	}
}

func TestPruneRaw_Build_Nil(t *testing.T) {
	var r *PruneRaw
	cfg, err := r.Build()
	if err != nil || cfg != nil {
		t.Fatalf("expected nil for nil PruneRaw")
	}
}

func TestPruneRaw_Build_Defaults(t *testing.T) {
	r := &PruneRaw{}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxAge != 24*time.Hour {
		t.Errorf("expected default MaxAge 24h, got %v", cfg.MaxAge)
	}
	if cfg.MaxResults != 500 {
		t.Errorf("expected default MaxResults 500, got %d", cfg.MaxResults)
	}
}

func TestPruneRaw_Build_ValidMaxAge(t *testing.T) {
	r := &PruneRaw{MaxAge: "6h"}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxAge != 6*time.Hour {
		t.Errorf("expected 6h, got %v", cfg.MaxAge)
	}
}

func TestPruneRaw_Build_InvalidMaxAge(t *testing.T) {
	r := &PruneRaw{MaxAge: "notaduration"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for invalid max_age")
	}
}

func TestPruneRaw_Build_NegativeMaxAge(t *testing.T) {
	r := &PruneRaw{MaxAge: "-1h"}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for negative max_age")
	}
}

func TestPruneRaw_Build_OnlyDrifted(t *testing.T) {
	r := &PruneRaw{OnlyDrifted: true, MaxResults: 10}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.OnlyDrifted {
		t.Error("expected OnlyDrifted true")
	}
	if cfg.MaxResults != 10 {
		t.Errorf("expected MaxResults 10, got %d", cfg.MaxResults)
	}
}
