package config

import (
	"testing"

	"github.com/driftwatch/internal/drift"
)

func TestTransformRaw_Build_Disabled(t *testing.T) {
	r := TransformRaw{Enabled: false}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil config when disabled")
	}
}

func TestTransformRaw_Build_MissingPrefix(t *testing.T) {
	r := TransformRaw{
		Enabled: true,
		Rules:   []TransformRuleRaw{{Op: "upper"}},
	}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestTransformRaw_Build_InvalidOp(t *testing.T) {
	r := TransformRaw{
		Enabled: true,
		Rules:   []TransformRuleRaw{{Prefix: "app.", Op: "rot13"}},
	}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestTransformRaw_Build_Valid(t *testing.T) {
	r := TransformRaw{
		Enabled: true,
		Rules: []TransformRuleRaw{
			{Prefix: "app.", Op: "prepend", Arg: "svc-"},
			{Prefix: "db.", Op: "upper"},
		},
	}
	cfg, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Op != drift.TransformOpPrefix {
		t.Errorf("expected prepend op, got %q", cfg.Rules[0].Op)
	}
	if cfg.Rules[0].Arg != "svc-" {
		t.Errorf("expected arg svc-, got %q", cfg.Rules[0].Arg)
	}
}

func TestBuildTransformRules_PropagatesError(t *testing.T) {
	_, err := BuildTransformRules([]TransformRuleRaw{
		{Prefix: "x.", Op: "invalid"},
	})
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestBuildTransformRules_Valid(t *testing.T) {
	rules, err := BuildTransformRules([]TransformRuleRaw{
		{Prefix: "svc.", Op: "lower"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Op != drift.TransformOpLower {
		t.Errorf("expected lower op")
	}
}
