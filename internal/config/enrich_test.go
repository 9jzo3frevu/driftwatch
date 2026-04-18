package config

import (
	"testing"
)

func TestEnrichRuleRaw_Build_MissingPrefix(t *testing.T) {
	r := EnrichRuleRaw{Metadata: map[string]string{"team": "ops"}}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for missing prefix")
	}
}

func TestEnrichRuleRaw_Build_NoMetadata(t *testing.T) {
	r := EnrichRuleRaw{Prefix: "db."}
	_, err := r.Build()
	if err == nil {
		t.Fatal("expected error for empty metadata")
	}
}

func TestEnrichRuleRaw_Build_Valid(t *testing.T) {
	r := EnrichRuleRaw{Prefix: "db.", Metadata: map[string]string{"team": "data"}}
	rule, err := r.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rule.Prefix != "db." {
		t.Errorf("expected prefix db., got %s", rule.Prefix)
	}
	if rule.Metadata["team"] != "data" {
		t.Errorf("expected team=data")
	}
}

func TestBuildEnrichRules_Valid(t *testing.T) {
	raw := []EnrichRuleRaw{
		{Prefix: "db.", Metadata: map[string]string{"tier": "critical"}},
		{Prefix: "cache.", Metadata: map[string]string{"tier": "low"}},
	}
	cfg, err := BuildEnrichRules(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(cfg.Rules))
	}
}

func TestBuildEnrichRules_PropagatesError(t *testing.T) {
	raw := []EnrichRuleRaw{
		{Prefix: "db.", Metadata: map[string]string{"tier": "critical"}},
		{Prefix: "", Metadata: map[string]string{"tier": "low"}},
	}
	_, err := BuildEnrichRules(raw)
	if err == nil {
		t.Fatal("expected error from invalid rule")
	}
}
