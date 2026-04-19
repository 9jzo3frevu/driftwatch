package config

import (
	"testing"
)

func TestRedactRaw_Build_Disabled(t *testing.T) {
	r := RedactRaw{Enabled: false}
	_, ok, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected disabled")
	}
}

func TestRedactRaw_Build_DefaultPatterns(t *testing.T) {
	r := RedactRaw{Enabled: true}
	cfg, ok, err := r.Build()
	if err != nil || !ok {
		t.Fatalf("expected success, got err=%v ok=%v", err, ok)
	}
	if len(cfg.Patterns) == 0 {
		t.Error("expected default patterns")
	}
	if cfg.Mask != "***REDACTED***" {
		t.Errorf("expected default mask, got %q", cfg.Mask)
	}
}

func TestRedactRaw_Build_CustomPatterns(t *testing.T) {
	r := RedactRaw{Enabled: true, Patterns: []string{"pin", "cvv"}}
	cfg, ok, err := r.Build()
	if err != nil || !ok {
		t.Fatalf("unexpected: %v %v", err, ok)
	}
	if len(cfg.Patterns) != 2 || cfg.Patterns[0] != "pin" {
		t.Errorf("unexpected patterns: %v", cfg.Patterns)
	}
}

func TestRedactRaw_Build_CustomMask(t *testing.T) {
	r := RedactRaw{Enabled: true, Mask: "<masked>"}
	cfg, _, _ := r.Build()
	if cfg.Mask != "<masked>" {
		t.Errorf("expected <masked>, got %q", cfg.Mask)
	}
}
