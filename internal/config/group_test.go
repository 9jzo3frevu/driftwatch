package config

import (
	"testing"

	"github.com/driftwatch/internal/drift"
)

func TestGroupRaw_Build_Disabled(t *testing.T) {
	r := GroupRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Enabled {
		t.Error("expected disabled")
	}
}

func TestGroupRaw_Build_DefaultBy(t *testing.T) {
	r := GroupRaw{Enabled: true, By: ""}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.By != drift.GroupByService {
		t.Errorf("expected service default, got %q", cfg.By)
	}
}

func TestGroupRaw_Build_ValidValues(t *testing.T) {
	for _, v := range []string{"service", "severity", "key"} {
		r := GroupRaw{Enabled: true, By: v}
		cfg, err := r.Build()
		if err != nil {
			t.Errorf("%s: unexpected error: %v", v, err)
		}
		if string(cfg.By) != v {
			t.Errorf("expected %q, got %q", v, cfg.By)
		}
	}
}

func TestGroupRaw_Build_InvalidBy(t *testing.T) {
	r := GroupRaw{Enabled: true, By: "region"}
	_, err := r.Build()
	if err == nil {
		t.Error("expected error for invalid by value")
	}
}
